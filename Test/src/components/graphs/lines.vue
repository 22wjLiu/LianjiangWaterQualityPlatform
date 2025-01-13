<template>
  <div
    class="container"
    ref="lines"
    v-loading="loading"
    element-loading-text="数据量较大，请耐心等待"
    :style="`width: ${containerWidth}`"
  >
    <h3>监测范围内水质基本参数均值&污染物监测指标均值</h3>
    <el-divider></el-divider>
    <LineGraph
      v-for="(item, index) in lineData"
      :key="index"
      :id="labelList[index]"
      :lineData="item"
      :indexOfTime="indexOfTime"
      :labelList="labelList"
      :index="index"
      @drawed="loadingChange"
    ></LineGraph>
  </div>
</template>

<script>
import LineGraph from '@/components/graphs/line.vue'
import bus from '@/util/eventBus'
export default {
  props: {
    pattern: {
      require: true,
      type: String
    }
  },
  data() {
    return {
      lineData: [],
      indexOfTime: '',
      labelList: [],
      loading: true
    }
  },
  computed: {
    containerHeight() {
      if (this.pattern === 'overview') {
        return '100%'
      } else {
        return `${window.innerHeight - 57}px`
      }
    },
    containerWidth() {
      if (this.pattern === 'overview') {
        return '100%'
      } else {
        return `${window.innerWidth - 10}px`
      }
    }
  },
  components: {
    LineGraph
  },
  methods: {
    loadingChange(val) {
      this.loading = val
    }
  },
  created() {
    bus.$on('lineData', val => {
      this.lineData = val.data
      this.indexOfTime = val.indexOfTime
      this.labelList = val.labelList
    })
    bus.$on('reload', val => {
      this.loading = val
    })
    window.addEventListener('resize', () => {
      this.containerHeight = window.innerHeight - 57
    })
  }
}
</script>

<style lang="less" scoped>
.container {
  h3 {
    text-align: center;
  }
  height: 300px;
  margin-top: 20px;
  z-index: 99;
}
</style>
