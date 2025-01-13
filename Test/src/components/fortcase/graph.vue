<template>
  <div ref="lineGraph" class="line-container"></div>
</template>

<script>
import { arr } from '@/assets/data'
import bus from '@/util/eventBus'
export default {
  data() {
    return {
      cardHeight: window.innerHeight - 95,
      data: arr,
      legend: ['水温', 'pH', '溶解氧', '浊度', '电导率', '氨氮', 'CODcr'],
      myChart: null
    }
  },
  methods: {
    processData() {
      this.data = this.data.map(item => {
        Object.keys(item).forEach(key => {
          if (item[key] === '0') {
            item[key] = ''
          }
        })
        return item
      })
    },
    buildSeries() {
      const series = []
      this.legend.forEach((_, index) => {
        series.push({
          type: 'line',
          symbol: 'none',
          smooth: true,
          sampling: 'lttb',
          lineStyle: {
            width: 1
          },
          name: this.legend[index],
          encode: {
            x: 0,
            y: index + 1
          },
          markLine: {
            label: {
              position: 'middle'
            },
            lineStyle: {
              type: 'solid'
            },
            data: [
              {
                xAxis: arr[arr.length - 1].时间
              }
            ]
          }
        })
      })
      return series
    },
    draw(series) {
      const graph = this.$refs.lineGraph
      this.myChart = this.$echarts.init(graph)
      const options = {
        dataset: {
          source: this.data
        },
        // 提示框组件
        tooltip: {
          // 触发类型：坐标轴触发
          trigger: 'axis'
        },
        title: {
          text: '对未来的观测值预测',
          left: 'center',
          textStyle: {
            fontFamily: 'SimSun',
            fontSize: '16'
          }
        },
        legend: {
          type: 'scroll',
          left: 'center',
          bottom: '15%',
          // 如果series 对象有name 值，则 legend可以不用写data
          // 修改图例组件 文字颜色
          textStyle: {
            color: '#4c9bfd',
            fontSize: 14
          }
        },
        grid: {
          top: '5%',
          left: '1%',
          right: '1%',
          bottom: '25%',
          show: false, // 显示边框
          borderColor: '#012f4a', // 边框颜色
          containLabel: true // 包含刻度文字在内
        },
        // 工具栏
        toolbox: {
          feature: {
            // 保存为图片
            saveAsImage: {}
          }
        },
        xAxis: {
          // 坐标轴类型
          type: 'time',
          // 坐标轴两边的留白策略
          boundaryGap: false
        },
        yAxis: {
          type: 'value',
          min: 'dataMin'
        },
        dataZoom: [
          {
            type: 'inside',
            start: 0,
            end: 100
          },
          {
            height: '10',
            start: 0,
            end: 100
          }
        ],
        series: series
      }
      this.myChart.setOption(options)
      // 当浏览器窗口缩放时，图标同时缩放
      window.addEventListener('resize', () => {
        this.myChart.resize()
      })
    }
  },
  created() {
    bus.$on('showLoading', () => {
      this.myChart.showLoading()
    })
    bus.$on('dataChange', val => {
      this.data = this.data.concat(val.dataset)
      this.myChart.hideLoading()
      this.myChart.setOption({
        dataset: {
          source: this.data
        }
      })
    })
    bus.$on('empty', () => {
      this.data = arr
      this.myChart.setOption({
        dataset: {
          source: this.data
        }
      })
    })
  },
  mounted() {
    this.processData()
    this.draw(this.buildSeries())
    bus.$emit('listMsg', {
      lastDate: this.lastDate
    })
  }
}
</script>

<style lang="less" scoped>
.line-container {
  width: 100%;
  height: 100%;
  box-sizing: border-box;
  padding: 20px;
}
</style>
