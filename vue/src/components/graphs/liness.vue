<template>
  <div class="container">
    <el-button
      style="position: absolute; left: calc(5% / 3); z-index: 999"
      type="primary"
      size="small"
      @click="handleUpdate()"
      >编辑</el-button
    >
    <el-dialog
      title="编辑可视化数据信息表"
      :visible.sync="dialogFormVisible"
      center
    >
      <el-form
        ref="dataForm"
        :rules="rules"
        :model="temp"
        label-position="left"
        label-width="100px"
      >
        <el-form-item label="站名" prop="stationName">
          <el-select
            v-model="temp.stationName"
            placeholder="请选择站名"
            clearable
          >
            <el-option
              v-for="item in stationOptions"
              :key="item.value"
              :label="item.value"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="开始时间" prop="startTime">
          <el-date-picker
            v-model="temp.startTime"
            type="datetime"
            placeholder="选择日期时间"
            clearable
          >
          </el-date-picker>
        </el-form-item>
        <el-form-item label="结束时间" prop="endTime">
          <el-date-picker
            v-model="temp.endTime"
            type="datetime"
            placeholder="选择日期时间"
            clearable
          >
          </el-date-picker>
        </el-form-item>
      </el-form>
      <div slot="footer">
        <el-button @click="dialogFormVisible = false"> 取消 </el-button>
        <el-button type="primary" @click="updateData()"> 确认 </el-button>
      </div>
    </el-dialog>
    <div class="lineGraph" ref="liness"></div>
  </div>
</template>

<script>
import {
  getLineData,
  getStationName,
  getTimeRange,
  getActiveMapInfosByStationName,
} from "@/api/graphs.js";
import { fullFormatTime } from "@/util/timeFormater";
import bus from "@/util/eventBus";
export default {
  data() {
    return {
      graph: null,
      myChart: null,
      dialogFormVisible: false,
      stationName: "",
      temp: {
        stationName: "",
        startTime: "",
        endTime: "",
      },
      startTime: "",
      endTime: "",
      dataMinTime: "",
      dataMaxTime: "",
      lineData: [],
      options: [],
      stationOptions: [],
      indexOfTime: -1,
      system: "小时制",
      rules: {
        stationName: [
          { required: true, message: "站名不能为空", trigger: "change" },
        ],
        startTime: [
          { required: true, message: "开始时间不能为空", trigger: "change" },
        ],
        endTime: [
          { required: true, message: "结束时间不能为空", trigger: "change" },
          {
            validator: (rule, value, callback) => {
              if (!value || !this.temp.startTime) {
                callback();
              } else if (value <= this.temp.startTime) {
                callback(new Error("结束时间必须大于开始时间"));
              } else {
                callback();
              }
            },
            trigger: "change",
          },
        ],
      },
    };
  },
  methods: {
    fullFormatTime,
    handleUpdate() {
      this.dialogFormVisible = true;
      this.temp.stationName = this.stationName;
      this.temp.startTime = this.startTime;
      this.temp.endTime = this.endTime;
      this.$nextTick(() => {
        this.$refs.dataForm.clearValidate();
      });
    },
    async updateData() {
      const valid = await this.$refs.dataForm.validate();

      if (valid) {
        this.getLineList(this.temp.startTime, this.temp.endTime);
      }
    },
    getLineList(start, end) {
      let params = "";
      if (start) {
        params += `start=${this.fullFormatTime(start)}&`;
      }
      if (end) {
        params += `end=${this.fullFormatTime(end)}&`;
      }
      this.options.forEach((item) => {
        params += "fields=" + item.key + "&";
      });
      const last = params.lastIndexOf("&");
      params = params.slice(0, last);

      getLineData(this.stationName, this.system, params)
        .then((res) => {
          if (res.code === 200) {
            this.dialogFormVisible = false;
            this.startTime = res.data.startTime;
            this.endTime = res.data.endTime;
            let data = res.data.resultArr;
            // 将表中数值为 0 的字段赋值为 ''
            data = data.map((item) => {
              Object.keys(item).forEach((key) => {
                if (item[key] <= 0) {
                  item[key] = "";
                }
              });
              return item;
            });
            this.lineData = data;
            // 计算 data 中 time 字段的索引
            this.indexOfTime = Object.keys(data[0]).indexOf("time");

            this.myChart.hideLoading();
            this.draw(this.buildSeries(), this.buildSelected());
          }
        })
        .catch((err) => {
          this.$message.error(err.message);
          console.log(err);
        });
    },
    initChart() {
      this.graph = this.$refs.liness;
      this.myChart = this.$echarts.init(this.graph);
    },
    async initStationOptions() {
      try {
        const res = await getStationName();
        if (res.code === 200) {
          this.stationOptions = res.data.names;
          this.stationName = this.stationOptions[0].value;
        }
      } catch (err) {
        console.error("加载失败", err);
      }
    },
    initTimeRange() {
      getTimeRange(this.stationName, this.system)
        .then((res) => {
          if (res.code === 200) {
            this.dataMinTime = res.data.minTime;
            this.dataMaxTime = res.data.maxTime;
          }
        })
        .catch((err) => {
          console.error("加载失败", err);
        });
    },
    async initOptions() {
      try {
        const res = await getActiveMapInfosByStationName(
          "列字段映射",
          "海门湾桥闸"
        );
        if (res.code === 200) {
          this.options = res.data.mapInfos;
        }
      } catch (err) {
        console.error("加载失败", err);
      }
    },
    // 构造MyCharts.options.series
    buildSeries() {
      const series = [];
      Object.keys(this.lineData[0]).forEach((key, index) => {
        this.options.forEach((item) => {
          if (item.value === key) {
            series.push({
              type: "line",
              symbol: "none",
              smoott: true,
              sampling: "lttb",
              lineStyle: {
                width: 1,
              },
              name: item.key,
              encode: {
                x: this.indexOfTime,
                y: index,
              },
            });
          }
        });
      });
      return series;
    },
    buildSelected() {
      const selected = {};
      this.options.forEach((item) => {
        const key = item.key;
        for (let i = 0; i < this.lineData.length && i < 10; i++) {
          let value = Number(this.lineData[i][item.value]);

          // 如果是 NaN（表示乱码、非法等），按 0 处理
          if (isNaN(value)) {
            value = 0;
          }
          if (i > 0 && !selected[key]) continue;
          selected[key] = value < 100;
        }
      });

      return selected;
    },
    draw(series, selected) {
      const options = {
        dataset: [
          {
            source: this.lineData,
          },
        ],
        legend: {
          type: "scroll",
          // orient: 'vertical',
          bottom: "15%",
          left: "center",
          // 如果series 对象有name 值，则 legend可以不用写data
          // 修改图例组件 文字颜色
          textStyle: {
            color: "#4c9bfd",
            fontSize: 10,
          },
          // selected: {
          //   pH: true,
          //   溶解氧: true,
          //   水温: true,
          //   浊度: true
          // }
          selected: selected,
        },
        title: {
          text: this.stationName,
          left: "center",
          textStyle: {
            fontFamily: "SimSun",
            fontSize: "16",
          },
        },
        // 提示框组件
        tooltip: {
          // 触发类型：坐标轴触发
          trigger: "axis",
        },
        grid: {
          top: "5%",
          left: "1%",
          right: "1%",
          bottom: "25%",
          borderColor: "#012f4a", // 边框颜色
          containLabel: true, // 包含刻度文字在内
        },
        // 工具栏
        toolbox: {
          feature: {
            // 保存为图片
            saveAsImage: {},
          },
        },
        xAxis: {
          // 坐标轴类型
          type: "time",
          // 坐标轴两边的留白策略
          boundaryGap: false,
        },
        yAxis: {
          type: "value",
          min: "dataMin",
        },
        dataZoom: [
          {
            type: "inside",
            start: 0,
            end: 100,
          },
          {
            height: "10",
            start: 0,
            end: 100,
          },
        ],
        series: series,
      };
      this.myChart.setOption(options, true);
      window.addEventListener("resize", () => {
        this.myChart.resize();
      });
    },
  },
  async mounted() {
    await this.initStationOptions();
    this.initTimeRange();
    await this.initOptions();
    this.getLineList();
    this.initChart();
    this.draw();
    this.myChart.showLoading();
    bus.$on("siteChange", (val) => {
      this.myChart.showLoading();
      this.system = val.system;
      val.system === "月度制"
        ? (this.monthName = val.siteName)
        : (this.hourName = val.siteName);
    });
    bus.$on("hideLoading", () => {
      this.myChart.hideLoading();
    });
  },
  activated() {
    this.myChart.resize();
  },
};
</script>

<style lang="less" scoped>
.container {
  height: 100%;
}

.lineGraph {
  height: 100%;
  width: 100%;
  box-sizing: border-box;
  padding: 20px;
}

.container  :deep(.el-dialog) {
  max-width: 500px;
  min-width: 300px;

  .el-select {
    width: 100%;
  }

  .el-input {
    width: 100%;
  }
}
</style>
