<template>
  <div class="body" v-loading="loading" element-loading-text="数据量较大，请耐心等待">
    <div class="container" :style="`width: ${containerWidth}`">
      <h3>监测范围内水质基本参数均值&污染物监测指标均值</h3>
      <el-divider></el-divider>
      <LineGraph
        v-for="(item, index) in lineData"
        :key="item[0].time + index + item[item.length - 1].time"
        :id="options[index].value"
        :lineData="item"
        :indexOfTime="indexOfTime"
        :options="options"
        :index="index"
      ></LineGraph>
    </div>
  </div>
</template>

<script>
import LineGraph from "@/components/graphs/line.vue";
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
      lineData: [],
      indexOfTime: "",
      options: [],
      loading: true,
    };
  },
  computed: {
    containerHeight() {
      if (this.pattern === "overview") {
        return "100%";
      } else {
        return `${window.innerHeight - 57}px`;
      }
    },
    containerWidth() {
      if (this.pattern === "overview") {
        return "100%";
      } else {
        return `${window.innerWidth - 10}px`;
      }
    },
  },
  components: {
    LineGraph,
  },
  created() {
    bus.$on("lineData", (val) => {
      this.loading = false;
      this.lineData = JSON.parse(JSON.stringify(val.data));
      this.indexOfTime = val.indexOfTime;
      this.options = JSON.parse(JSON.stringify(val.options));
    });
    bus.$on("reload", (val) => {
      this.loading = val;
    });
    window.addEventListener("resize", () => {
      this.containerHeight = window.innerHeight - 57;
    });
  },
};
</script>

<style lang="less" scoped>
.body {
  height: 100vh;
}

.container {
  h3 {
    text-align: center;
  }
  height: 300px;
  margin-top: 20px;
  z-index: 99;
}
</style>
