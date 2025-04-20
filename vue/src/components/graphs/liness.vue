<template>
  <div class="container">
    <div class="lineGraph" ref="liness"></div>
  </div>
</template>

<script>
import { getActiveMapInfosByStationName } from "@/api/graphs.js";
import bus from "@/util/eventBus";
import { level } from "@/assets/level.js";
export default {
  data() {
    return {
      graph: null,
      myChart: null,
      indexOfTime: null,
      stationName: "",
      levels: level,
      lineData: [],
      options: [],
      levelLabels: ["I级", "II级", "III级", "IV级", "V级"],
      indexOfTime: -1,
    };
  },
  methods: {
    initChart() {
      this.graph = this.$refs.liness;
      this.myChart = this.$echarts.init(this.graph);
    },
    // 构造MyCharts.options.series
    buildSeries() {
      const series = [];
      Object.keys(this.lineData[0]).forEach((key, index) => {
        this.options.some((item) => {
          if (key === item.value) {
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
            return true;
          }
          return false;
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
        let len = this.lineData.length;
        for (let i = len - 1; i >= 0 && i >= len - 10; i--) {
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
          bottom: "15%",
          left: "center",
          // 如果series 对象有name 值，则 legend可以不用写data
          // 修改图例组件 文字颜色
          textStyle: {
            color: "#4c9bfd",
            fontSize: 10,
          },
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
    async initOptions() {
      try {
        const res = await getActiveMapInfosByStationName(
          "列字段映射",
          this.stationName
        );
        if (res.code === 200) {
          this.options = res.data.mapInfos;
          bus.$emit("curOptions", this.options);
          this.draw(this.buildSeries(), this.buildSelected());
        }
      } catch (err) {
        console.error("加载失败", err);
      }
    },
    getLevelIndex(val, thresholds, reverse = false) {
      for (let i = 0; i < thresholds.length; i++) {
        if (reverse ? val >= thresholds[i] : val <= thresholds[i]) {
          return i;
        }
      }
      return 4;
    },
    classify(list) {
      const reverseFields = new Set(["溶解氧"]);
      const classified = list.map((item) => {
        let worstLevel = -1;
        for (const key in this.levels) {
          const val = parseFloat(item[key]);
          if (!isNaN(val)) {
            const thresholds = this.levels[key];
            const isReverse = reverseFields.has(key);
            const index = this.getLevelIndex(val, thresholds, isReverse);
            worstLevel = Math.max(worstLevel, index);
          }
        }
        return {
          水质等级: this.levelLabels[worstLevel] || "错误数据",
        };
      });
      return classified;
    },
    conveyList(list) {
      const classified = this.classify(list);
      const levelCount = {};

      classified.forEach((item) => {
        const level = item.水质等级;
        levelCount[level] = (levelCount[level] || 0) + 1;
      });

      const pieData = this.levelLabels.map(label => ({
        name: label,
        value: levelCount[label] || 0
      }));

      bus.$emit("getDataByField", pieData);
    },
  },
  mounted() {
    this.initChart();
    bus.$on("showLoading", () => {
      this.myChart.showLoading();
    });
    bus.$on("drawLine", (stationName, indexOfTime, data) => {
      this.myChart.hideLoading();
      this.stationName = stationName;
      this.indexOfTime = indexOfTime;
      this.lineData = data;
      this.initOptions();
    });
    bus.$on("requestByfield", (fields) => {
      const filtered = this.lineData.map((item) => {
        const obj = {};
        fields.forEach((field) => {
          if (field.value in item) {
            obj[field.key] = item[field.value];
          }
        });
        return obj;
      });
      this.conveyList(filtered);
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
</style>
