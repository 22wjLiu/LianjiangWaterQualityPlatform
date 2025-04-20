<template>
  <div class="container">
    <el-button size="small" type="primary" class="handle-button" @click="handleForecate"
      >预测设置</el-button
    >
    <div ref="lineGraph" class="line-container"></div>
    <el-dialog title="预测设置表" :visible.sync="dialogFormVisible" center>
      <el-form ref="dataForm" label-position="left" label-width="70px">
        <el-form-item label="预测元素">
          <el-select v-model="field" placeholder="请设置预测元素">
            <el-option
              v-for="item in options"
              :key="item.key"
              :label="item.key"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="开始时间">
          <el-date-picker
            v-model="start"
            type="datetime"
            placeholder="选择日期时间"
            clearable
          >
          </el-date-picker>
        </el-form-item>
        <el-form-item label="结束时间">
          <el-date-picker
            v-model="end"
            type="datetime"
            placeholder="选择日期时间"
            clearable
          >
          </el-date-picker>
        </el-form-item>
      </el-form>
      <div slot="footer">
        <el-button @click="dialogFormVisible = false"> 取消 </el-button>
        <el-button type="primary" @click="forecastData()"> 确认 </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { forecast } from "@/api/forecast.js";
import { formatTime } from "@/util/timeFormater";
import bus from "@/util/eventBus";
export default {
  data() {
    return {
      dialogFormVisible: false,
      field: "",
      selected: "",
      start: "",
      end: "",
      system: "",
      stationName: "",
      cardHeight: window.innerHeight - 95,
      data: [],
      options: [],
      myChart: null,
    };
  },
  methods: {
    formatTime,
    initChart() {
      const graph = this.$refs.lineGraph;
      this.myChart = this.$echarts.init(graph);
    },
    buildSeries() {
      const series = [];
      series.push({
          type: "line",
          symbol: "none",
          smooth: true,
          sampling: "lttb",
          lineStyle: {
            width: 1,
          },
          name: "预测值",
          encode: {
            x: "time",
            y: "value",
          },
          markLine: {
            label: {
              position: "middle",
            },
            lineStyle: {
              type: "solid",
            },
          },
        });
      return series;
    },
    draw(series) {
      const options = {
        dataset: {
          source: this.data,
          dimensions: [
            { name: "time", type: "time" }, // 第 1 维：时间
            { name: "value", type: "number" }, // 第 2 维：数值
          ],
        },
        // 提示框组件
        tooltip: {
          // 触发类型：坐标轴触发
          trigger: "axis",
        },
        title: {
          text: `对${this.selected}的未来观测值`,
          left: "center",
          textStyle: {
            fontFamily: "SimSun",
            fontSize: "16",
          },
        },
        legend: {
          type: "scroll",
          left: "center",
          bottom: "15%",
          // 如果series 对象有name 值，则 legend可以不用写data
          // 修改图例组件 文字颜色
          textStyle: {
            color: "#4c9bfd",
            fontSize: 14,
          },
        },
        grid: {
          top: "5%",
          left: "1%",
          right: "1%",
          bottom: "25%",
          show: false, // 显示边框
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
      this.myChart.setOption(options);
    },
    handleForecate() {
      this.dialogFormVisible = true;
    },
    forecastData() {
      if (!this.start) {
        this.$message.warning("开始时间不能为空");
        return;
      }
      if (!this.end) {
        this.$message.warning("结束时间不能为空");
        return;
      }
      if (new Date(this.start) > new Date(this.end)) {
        this.$message.warning("结束时间不能小于开始时间");
        return;
      }
      this.dialogFormVisible = false;
      this.myChart.showLoading({
        text: "正在预测中...",
      });
      const params = `?system=${this.system}&station_name=${this.stationName}&field=${this.field}`;
      forecast(this.formatTime(this.start), this.formatTime(this.end), params)
        .then((res) => {
          if (res.code === 200) {
            this.data = res.data.result.forecast;
            if (this.data) {
              this.myChart.hideLoading();
              this.options.some((item) => {
                if(item.value === this.field) {
                  this.selected = item.key;
                  return true;
                }
                return false;
              })
              this.draw(this.buildSeries());
            }
          }
        })
        .catch((err) => {
          this.noData();
          this.$message.error("预测失败");
          console.log(err.message);
        });
    },
    noData() {
      this.myChart.clear();
      this.myChart.showLoading({
        text: "暂无预测数据",
        showSpinner: false,
        maskColor: "rgba(248, 249, 251, 1)",
      });
    },
  },
  mounted() {
    this.initChart();
    this.noData();
    window.addEventListener("resize", () => {
        this.myChart.resize();
    });
    bus.$on("curOptions", (opts) => {
      this.options = opts;
      if (this.options) {
        this.field = this.options[0].value;
      }
    });
    bus.$on("dataChange", (start, end, system, stationName) => {
      this.start = start;
      this.end = end;
      this.system = system;
      this.stationName = stationName;
    });
  },
};
</script>

<style lang="less" scoped>
.container {
  width: 100%;
  height: 100%;
  .handle-button {
    position: absolute;
    left: calc(98.33% - 79px);
    z-index: 999;
    opacity: 0;
  }

  :deep(.el-dialog) {
    max-width: 500px;
    min-width: 300px;

    .el-select {
      width: 100%;
    }

    .el-input{
      width: 100%;
    }
  }
}

.container:hover .handle-button{
  opacity: 1;
}

.line-container {
  width: 100%;
  height: 100%;
  box-sizing: border-box;
  padding: 20px;
}
</style>
