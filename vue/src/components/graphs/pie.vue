<template>
  <div class="container" :class="{ flex: showMd }">
    <div :class="{ left: showMd, container: !showMd }">
      <div class="select">
        <!-- <span>的各水质等级所占时间的比例</span> -->
      </div>
      <div
        class="pieGraph"
        ref="pieGraph"
        :style="`height: ${graphHeight};`"
      ></div>
    </div>

    <div
      v-if="showMd"
      class="markdown-body"
      :class="{ right: showMd }"
      :style="`height: ${cardHeight}`"
    >
      <mdfile></mdfile>
    </div>
  </div>
</template>

<script>
import mdfile from "@/assets/mdfile.md";
import bus from "@/util/eventBus";
export default {
  props: {
    pattern: {
      require: true,
      type: String,
    },
  },
  data() {
    return {
      data: [],
      myChart: null,
      curOptions: [],
      keySet: null,
      options: [
        "溶解氧",
        "高锰酸盐指数",
        "CODcr",
        "化学需氧量",
        "五日生化需氧量",
        "氨氮",
        "总磷",
        "总氮",
        "铜",
        "锌",
        "氟化物",
        "氰化物",
        "硒",
        "砷",
        "汞",
        "铬",
        "挥发酚",
        "石油类",
        "阴离子表面活性剂",
        "硫化物",
        "粪大肠菌群"
      ],
    };
  },
  computed: {
    graphHeight() {
      return this.pattern === "overview"
        ? "100%"
        : `${(window.innerHeight - 57) * 0.9}px`;
    },
    cardHeight() {
      return this.pattern === "overview"
        ? "100%"
        : `${window.innerHeight - 93 - 20}px`;
    },
    showMd() {
      let bool = false;
      this.pattern === "overview" ? (bool = false) : (bool = true);
      return bool;
    },
  },
  methods: {
    initChart() {
      const graph = this.$refs.pieGraph;
      this.myChart = this.$echarts.init(graph);
    },
    draw() {
      const options = {
        // 图例组件。通过点击图例组件控制某个系列的显示与否
        legend: {
          type: "scroll",
          right: "center",
          bottom: "3%",
          // 如果series 对象有name 值，则 legend可以不用写data
          // 修改图例组件 文字颜色
          textStyle: {
            color: "#4c9bfd",
            fontSize: 14,
          },
        },
        title: {
          text: "水质等级",
          left: "center",
          top: "10%",
          textStyle: {
            fontFamily: "SimSun",
            fontSize: "16",
          },
        },
        grid: {
          left: "5%",
          right: "5%",
          bottom: "5%",
          top: "5%",
        },
        series: [
          {
            type: "pie",
            data: this.data,
            radius: [0, '60%'],
            center: ["50%", "50%"],
            itemStyle: {
              borderRadius: 5,
            },
            label: {
              show: true,
              backgroundColor: "#F6F8FC",
              formatter: '{b}: ({d}%)',
            },
            labelLine: {
              show: true,
              showAbove: true,
              smooth: true,
              length: 10,
              length2: 10,
            },
            emphasis: {
              label: {
                show: true,
              },
            },
          },
        ],
      };
      this.myChart.setOption(options);
      // 当浏览器窗口缩放时，图标同时缩放
      window.addEventListener("resize", () => {
        this.myChart.resize();
      });
    },
    optionChange() {
      bus.$emit("requestByfield", this.curOptions);
    },
  },
  mounted() {
    this.keySet = new Set(this.options);
    this.initChart();
    bus.$on("showLoading", () => {
      this.myChart.showLoading();
    });
    bus.$on("curOptions", (val) => {
      this.curOptions = val.filter(v => this.keySet.has(v.key));
      if(this.curOptions){
        this.optionChange();
      }
    });
    bus.$on("getDataByField", (val) => {
      this.data = val;
      this.myChart.hideLoading();
      this.draw();
    });
    bus.$emit("ready");
  },
  components: {
    mdfile,
  },
  created() {
    window.addEventListener("resize", () => {
      if (this.pattern === "overview") return;
      this.graphHeight = (window.innerHeight - 57) * 0.9;
      this.cardHeight = window.innerHeight - 93 - 20;
    });
  },
};
</script>

<style lang="less" scoped>
.container {
  width: 100%;
  height: 100%;
  position: relative;
}

.flex {
  display: flex;
}

.left {
  width: 60%;
}
.right {
  width: 40%;
  margin: 18px;
  overflow: auto;
  padding: 10px;
  border-radius: 4px;
  box-shadow: 0 0 10px 0 #888888;
}
.select {
  width: 100%;
  display: flex;
  justify-content: center;
  padding-top: 10px;
  position: absolute;
  top: 0;
  z-index: 9;
  span {
    display: inline-block;
    height: 30px;
    line-height: 30px;
    font-weight: 700;
  }
}
.pieGraph {
  width: 100%;
}
</style>
