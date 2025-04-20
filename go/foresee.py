import sys
import json
import warnings
from pathlib import Path

import numpy as np
import pandas as pd

from prophet import Prophet
from statsmodels.tsa.arima.model import ARIMA
from statsmodels.tsa.statespace.sarimax import SARIMAX
from statsmodels.tsa.holtwinters import ExponentialSmoothing
from sklearn.metrics import mean_absolute_error
from sklearn.preprocessing import MinMaxScaler

from tensorflow.keras.models import Sequential
from tensorflow.keras.layers import LSTM, Dense
from tensorflow.keras.callbacks import EarlyStopping

warnings.filterwarnings("ignore")


class WaterQualityPredictor:

    # --------------------------- INIT ---------------------------
    def __init__(self, csv_path: str, target_field: str, *, forecast_freq: str = "H"):
        self.timestamp_col = "time"
        self.target_field = target_field
        self.forecast_freq = forecast_freq.upper()

        self._load_data(csv_path)
        self._preprocess()

        # model containers
        self.prophet_models, self.arima_models = {}, {}
        self.sarima_models, self.ets_models = {}, {}
        self.lstm_models, self.scalers = {}, {}

        self._train_all_models()

    # ------------------------ DATA -----------------------------
    def _load_data(self, path: str):
        df = pd.read_csv(path)
        need = {self.timestamp_col, self.target_field}
        if need - set(df.columns):
            raise ValueError(f"CSV 必须包含列: {need}")
        self.df_raw = df

    def _impute(self, s: pd.Series) -> pd.Series:
        s = s.asfreq(self.forecast_freq)
        if s.isna().all():
            raise ValueError("时间序列全为空，无法填补")

        s = s.interpolate(method="time", limit=6, limit_direction="both").ffill().bfill()
        s = s.fillna(s.median())
        return s

    def _preprocess(self):
        df = self.df_raw[[self.timestamp_col, self.target_field]].copy()
        df.columns = ["ds", "y"]
        df["ds"] = pd.to_datetime(df["ds"], errors="coerce")
        df.dropna(subset=["ds"], inplace=True)
        df.sort_values("ds", inplace=True)
        df["y"] = pd.to_numeric(df["y"], errors="coerce")
        df.dropna(subset=["y"], inplace=True)
        df = df.groupby("ds", as_index=False).mean()
        if df.empty:
            raise ValueError("❌ 无有效数据行")
        df.set_index("ds", inplace=True)
        s_filled = self._impute(df["y"])
        df = s_filled.to_frame(name="y")  
        df.reset_index(inplace=True)
        self.df = df
        self.ts = df.set_index("ds")["y"]
        assert not self.df["y"].isna().any(), "预处理后 y 仍含 NaN"

    # ---------------------- TRAINING ---------------------------
    def _train_prophet(self):
        m = Prophet()
        m.fit(self.df)
        self.prophet_models[self.target_field] = m

    def _train_arima(self):
        import contextlib, os
        with contextlib.redirect_stdout(open(os.devnull, "w")):
            self.arima_models[self.target_field] = ARIMA(self.ts, order=(5, 1, 0)).fit()

    def _train_sarima(self):
        import contextlib, os
        with contextlib.redirect_stdout(open(os.devnull, "w")):
            self.sarima_models[self.target_field] = SARIMAX(
                self.ts, order=(5, 1, 0), seasonal_order=(1, 1, 1, 24)
            ).fit(disp=False)

    def _train_ets(self):
        import contextlib, os
        with contextlib.redirect_stdout(open(os.devnull, "w")):
            self.ets_models[self.target_field] = ExponentialSmoothing(
                self.ts, trend="add", seasonal="add", seasonal_periods=24
            ).fit()

    def _train_lstm(self):
        scaler = MinMaxScaler()
        scaled = scaler.fit_transform(self.ts.values.reshape(-1, 1)).flatten()
        self.scalers[self.target_field] = scaler
        X, y = scaled[:-1].reshape(-1, 1, 1), scaled[1:]
        model = Sequential([LSTM(50, activation="relu", input_shape=(1, 1)), Dense(1)])
        model.compile(optimizer="adam", loss="mse")
        model.fit(X, y, epochs=50, batch_size=32, validation_split=0.2, shuffle=False,
                  callbacks=[EarlyStopping(monitor="val_loss", patience=5, restore_best_weights=True)], verbose=0)
        self.lstm_models[self.target_field] = model

    def _train_all_models(self):
        self._train_prophet(); self._train_arima(); self._train_sarima(); self._train_ets(); self._train_lstm()

    # -------------------- EVALUATION ---------------------------
    def _evaluate(self):
        y = self.ts.values
        mae = {
            "Prophet": mean_absolute_error(y, self.prophet_models[self.target_field].predict(self.df[["ds"]])["yhat"].values),
            "ARIMA": mean_absolute_error(y, self.arima_models[self.target_field].predict(0, len(y)-1)),
            "SARIMA": mean_absolute_error(y, self.sarima_models[self.target_field].predict(0, len(y)-1)),
            "ETS": mean_absolute_error(y, self.ets_models[self.target_field].fittedvalues),
        }
        # LSTM
        scaler = self.scalers[self.target_field]
        lstm_scaled = self.lstm_models[self.target_field].predict(y[:-1].reshape(-1,1,1), verbose=0).flatten()
        mae["LSTM"] = mean_absolute_error(y[1:], scaler.inverse_transform(lstm_scaled.reshape(-1,1)).flatten())
        best = min(mae, key=mae.get)
        return best, mae[best]

    # --------------------- PRED HELPERS ------------------------
    @staticmethod
    def _safe_stat_predict(model, idx, rng):
        before, after = rng[rng <= idx[-1]], rng[rng > idx[-1]]
        out = []
        if not before.empty:
            s, e = idx.get_indexer([before[0]])[0], idx.get_indexer([before[-1]])[0]
            s, e = max(0, s), max(0, e)
            out.extend(model.predict(start=s, end=e))
        if not after.empty:
            out.extend(model.forecast(len(after)))
        return out

    # ----------------------- FORECAST -------------------------
    def forecast_between(self, start: str, end: str):
        best, mae = self._evaluate()
        st, et = pd.to_datetime(start), pd.to_datetime(end)
        if et < st:
            raise ValueError("结束时间必须不早于开始时间")
        rng = pd.date_range(st, et, freq=self.forecast_freq)

        if best == "Prophet":
            yhat = self.prophet_models[self.target_field].predict(pd.DataFrame({"ds": rng}))["yhat"].values.tolist()
        elif best == "ARIMA":
            yhat = self._safe_stat_predict(self.arima_models[self.target_field], self.ts.index, rng)
        elif best == "SARIMA":
            yhat = self._safe_stat_predict(self.sarima_models[self.target_field], self.ts.index, rng)
        elif best == "ETS":
            yhat = self._safe_stat_predict(self.ets_models[self.target_field], self.ts.index, rng)
        else:  # LSTM autoregressive
            scaler = self.scalers[self.target_field]
            hist_scaled = self.lstm_models[self.target_field].predict(self.ts.values[:-1].reshape(-1,1,1), verbose=0).flatten()
            hist = scaler.inverse_transform(hist_scaled.reshape(-1,1)).flatten(); hist = np.insert(hist,0,np.nan)
            series = dict(zip(pd.date_range(self.ts.index[0], periods=len(hist), freq=self.forecast_freq), hist))
            cur_scaled, cur_t = scaler.transform([[self.ts.values[-1]]])[0,0], self.ts.index[-1]
            while cur_t < et:
                cur_t += pd.Timedelta(self.forecast_freq)
                cur_scaled = self.lstm_models[self.target_field].predict(np.array([[cur_scaled]]).reshape(1,1,1), verbose=0)[0,0]
                series[cur_t] = scaler.inverse_transform([[cur_scaled]])[0,0]
            yhat = [series.get(t, np.nan) for t in rng]

        # pad & fill
        yhat = (yhat + [np.nan]*len(rng))[:len(rng)]
        ser = pd.Series(yhat, index=rng).ffill().bfill()
        if ser.isna().all():
            ser[:] = self.ts.iloc[-1]
        return best, mae, pd.DataFrame({"time": rng, "value": ser.tolist()})


# ------------------------------------------------------------------
# CLI
# ------------------------------------------------------------------
if __name__ == "__main__":
    if len(sys.argv) != 6:
        print("Usage: python predictor.py <csv> <field> <start> <end> <freq>", file=sys.stderr)
        sys.exit(1)

    path, field, st, et, fq = sys.argv[1:6]
    try:
        predictor = WaterQualityPredictor(path, field, forecast_freq=fq)
        best_model, mae, df_fc = predictor.forecast_between(st, et)
        resp = {
            "best_model": best_model,
            "mae": float(mae),
            "forecast": json.loads(df_fc.to_json(orient="records", date_format="iso"))
        }
        print(json.dumps(resp, ensure_ascii=False, allow_nan=False))
    except Exception as exc:
        print(str(exc), file=sys.stderr)
        sys.exit(1)
