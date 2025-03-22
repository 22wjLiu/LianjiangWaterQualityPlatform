import pandas as pd
import numpy as np
import torch as t

Max_iter = 10000000

def Metrics(pre, real):
    from sklearn import metrics

    mse = metrics.mean_squared_error(real, pre)
    rmse = np.sqrt(metrics.mean_squared_error(real, pre))
    mae = metrics.mean_absolute_error(real, pre)

    # MAPE = np.mean(np.abs((pre - real) / real)) * 100
    # SMAPE = 2.0 * np.mean(np.abs(pre - real) / (np.abs(pre) + np.abs(real))) * 100
    return mse, rmse, mae

def Machine_learning_fit(train_x, train_y, test_x, test_y):
    best_mae = 1e9
    idx = 0

    from sklearn.linear_model import LinearRegression
    model1 = LinearRegression().fit(train_x, train_y)
    predict_y = model1.predict(test_x)
    mse = Metrics(predict_y, test_y)
    if(mse[2] < best_mae):
        idx = 1
        best_mae = mse[0]
    # print("算法Linear回归   MSE : %12.6lf,   RMSE : %12.6lf,   MAE : %12.6lf" % Metrics(predict_y, test_y))

    from sklearn.linear_model import Ridge
    model2 = Ridge(max_iter = Max_iter).fit(train_x, train_y)
    predict_y = model2.predict(test_x)
    mse= Metrics(predict_y, test_y)
    if(mse[2] < best_mae):
        idx = 2
        best_mae = mse[0]
    # print("算法Ridge回归    MSE : %12.6lf,   RMSE : %12.6lf,   MAE : %12.6lf" % Metrics(predict_y, test_y))

    from sklearn.linear_model import Lasso
    model5 = Lasso(max_iter = Max_iter).fit(train_x, train_y)
    predict_y = model5.predict(test_x)
    mse= Metrics(predict_y, test_y)
    if(mse[2] < best_mae):
        idx = 5
        best_mae = mse[0]
    # print("算法Lasso回归    MSE : %12.6lf,   RMSE : %12.6lf,   MAE : %12.6lf" % Metrics(predict_y, test_y))

    from sklearn.neighbors import KNeighborsRegressor
    model3 = KNeighborsRegressor(n_neighbors = 4).fit(train_x, train_y)
    predict_y = model3.predict(test_x)
    mse= Metrics(predict_y, test_y)
    if(mse[2] < best_mae):
        idx = 3
        best_mae = mse[0]
    # print("算法KNN回归      MSE : %12.6lf,   RMSE : %12.6lf,   MAE : %12.6lf" % Metrics(predict_y, test_y))

    from sklearn.svm import SVR
    model4 = SVR(max_iter = Max_iter).fit(train_x, train_y.ravel())
    predict_y = model4.predict(test_x)
    mse= Metrics(predict_y, test_y)

    if(mse[2] < best_mae):
        idx = 4
        best_mae = mse[0]
    # print("算法SVR回归      MSE : %12.6lf,   RMSE : %12.6lf,   MAE : %12.6lf" % Metrics(predict_y, test_y))
    
    if idx == 1:
        # print("1")
        return model1
    elif idx == 2:
        # print("2")
        return model2
    elif idx ==  3:
        # print("3")
        return model3
    elif idx ==  4:
        # print("4")
        return model4
    else:
        # print("5")
        return model5

def prediction(dataset, target):
    #string = '水温'
    arr = ['2018-06-01 00:00:00','水温','pH','溶解氧','浊度','电导率','氨氮','CODcr']
    #a = np.array(df[string])

    # 先获取数据，准备扩充
    extend = dataset[arr[1]], dataset[arr[2]], dataset[arr[3]], dataset[arr[4]]
    extend = np.array(extend).T
    #扩充，准备预测

    # y
    temp = dataset[arr[target]]
    temp = np.array(temp)

    # print(extend.shape)
    # print(temp.shape)
    p = int(dataset.shape[0] * 0.8)
    end = int(dataset.shape[0])

    train_x = extend[0 : p]
    train_y = temp[0 : p]


    test_x = extend[p : end]
    test_y = temp[p : end]

    train_x = train_x.reshape(-1, 4)
    train_y = train_y.reshape(-1, 1)

    test_x = test_x.reshape(-1, 4)
    test_y = test_y.reshape(-1, 1)

    # print(train_x.shape, train_y.shape)

    best_model = Machine_learning_fit(train_x, train_y, test_x, test_y)

    return best_model


if __name__ == '__main__':
    df = pd.read_excel("dataset.xlsx", header = 0)
    import argparse

    parser = argparse.ArgumentParser()
    parser.add_argument('--Temperature', type=float, default=0.3)  # 采样率
    parser.add_argument('--PH', type=float, default=0.1)
    parser.add_argument('--Turbidity', type=float, default=0.1)
    parser.add_argument('--DO', type=float, default=0.1)

    args = parser.parse_args()

    inputs = np.array([args.Temperature, args.PH, args.Turbidity, args.DO]).reshape(-1, 4)
    # '电导率','氨氮','CODcr' 
    # print(inputs.shape)
    
    model = prediction(df, 5)
    print(model.predict(inputs).item())
    model = prediction(df, 6) 
    print(model.predict(inputs).item())
    model = prediction(df, 7)
    print(model.predict(inputs).item())
    # print(pre.shape)
    

