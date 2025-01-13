<template>
  <div class="container" ref="line"></div>
</template>

<script>
export default {
  data() {
    return {
      graph: null,
      myChart: null
    }
  },
  props: ['lineData', 'indexOfTime', 'labelList', 'index'],
  methods: {
    draw() {
      this.graph = this.$refs.line
      this.myChart = this.$echarts.init(this.graph)
      const options = {
        dataset: [
          {
            source: this.lineData
          }
        ],
        title: {
          text: this.labelList[this.index],
          left: 'center',
          textStyle: {
            fontFamily: 'SimSun'
          }
        },
        // 提示框组件
        tooltip: {
          // 触发类型：坐标轴触发
          trigger: 'axis'
        },
        grid: {
          top: '10%',
          left: '10%',
          right: '3%',
          bottom: '15%',
          show: false, // 显示边框
          borderColor: '#012f4a', // 边框颜色
          containLabel: true // 包含刻度文字在内
        },
        // 工具栏
        toolbox: {
          feature: {
            // 保存为图片
            saveAsImage: {}
          },
          right: '2%'
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
            height: '20',
            start: 0,
            end: 100
          }
        ],
        series: {
          type: 'line',
          // 标记的图形
          symbol: 'none',
          // 采用 Largest-Triangle-Three-Bucket 算法，可以最大程度保证采样后线条的趋势，形状和极值
          sampling: 'lttb',

          lineStyle: {
            width: 1
          },
          name: this.labelList[this.index],
          encode: {
            x: this.indexOfTime,
            y: this.indexOfTime === 0 ? 1 : 0
          },
          areaStyle: {
            color: new this.$echarts.graphic.LinearGradient(1, 0, 0, 0, [
              {
                offset: 0,
                color: `rgb(${Math.floor(
                  Math.random() * 128 + 128
                )}, ${Math.floor(Math.random() * 128 + 128)}, ${Math.floor(
                  Math.random() * 128 + 128
                )})`
              },
              {
                offset: 1,
                color: `rgb(${Math.floor(
                  Math.random() * 128 + 128
                )}, ${Math.floor(Math.random() * 128 + 128)}, ${Math.floor(
                  Math.random() * 128 + 128
                )})`
              }
            ])
          }
        }
      }
      this.myChart.setOption(options, true)
      window.addEventListener('resize', () => {
        this.myChart.resize()
      })
    }
  },
  mounted() {
    this.draw()
    this.$emit('drawed', false)
  },
  activated() {
    this.myChart.resize()
  }
}
</script>

<style lang="less" scoped></style>
