<template>
  <div class="container">
    <div ref="barGraph" :style="`height: ${graphHeight}; width: 100%`"></div>
  </div>
</template>

<script>
import { arr } from '@/assets/level'
export default {
  props: {
    pattern: {
      require: true,
      type: String
    }
  },
  data() {
    return {
      data: arr
    }
  },
  computed: {
    graphHeight() {
      return this.pattern === 'overview'
        ? '100%'
        : `${window.innerHeight - 93}px`
    }
  },
  methods: {
    draw() {
      const graph = this.$refs.barGraph
      const myChart = this.$echarts.init(graph)
      const options = {
        dataset: {
          source: this.data
        },
        // 提示框组件
        tooltip: {
          // 触发类型：坐标轴触发
          trigger: 'axis',
          axisPointer: {
            type: 'shadow'
          }
        },
        // title: {
        //   text: '监测范围内各项指标不同水质等级所占天数',
        //   left: 'center',
        //   top: '3%',
        //   textStyle: {
        //     fontFamily: 'SimSun',
        //     fontSize: '14'
        //   }
        // },
        // 图例组件。通过点击图例组件控制某个系列的显示与否
        legend: {
          type: 'scroll',
          // orient: 'vertical',
          left: 'center',
          bottom: '3%',
          // 如果series 对象有name 值，则 legend可以不用写data
          // 修改图例组件 文字颜色
          textStyle: {
            color: '#4c9bfd',
            fontSize: 10
          }
        },
        grid: {
          top: '3%',
          left: '0',
          right: '0',
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
          }
        },
        xAxis: {
          // 坐标轴类型
          type: 'category',
          // 坐标轴两边的留白策略
          boundaryGap: true
        },
        yAxis: {
          type: 'value',
          name: '天数',
          nameTextStyle: {
            fontSize: 14,
            align: 'right',
            lineHeight: 30
          }
        },
        series: [
          {
            type: 'bar',
            emphasis: {
              focus: 'series'
            },
            itemStyle: {
              borderRadius: [10, 10, 0, 0]
            }
          },
          {
            type: 'bar',
            emphasis: {
              focus: 'series'
            },
            itemStyle: {
              borderRadius: [10, 10, 0, 0]
            }
          },
          {
            type: 'bar',
            emphasis: {
              focus: 'series'
            },
            itemStyle: {
              borderRadius: [10, 10, 0, 0]
            }
          },
          {
            type: 'bar',
            emphasis: {
              focus: 'series'
            },
            itemStyle: {
              borderRadius: [10, 10, 0, 0]
            }
          },
          {
            type: 'bar',
            emphasis: {
              focus: 'series'
            },
            itemStyle: {
              borderRadius: [10, 10, 0, 0]
            }
          },
          {
            type: 'bar',
            emphasis: {
              focus: 'series'
            },
            itemStyle: {
              borderRadius: [10, 10, 0, 0]
            }
          },
          {
            type: 'bar',
            emphasis: {
              focus: 'series'
            },
            itemStyle: {
              borderRadius: [10, 10, 0, 0]
            }
          },
          {
            type: 'bar',
            emphasis: {
              focus: 'series'
            },
            itemStyle: {
              borderRadius: [10, 10, 0, 0]
            }
          },
          {
            type: 'bar',
            emphasis: {
              focus: 'series'
            },
            itemStyle: {
              borderRadius: [10, 10, 0, 0]
            }
          },
          {
            type: 'bar',
            emphasis: {
              focus: 'series'
            },
            itemStyle: {
              borderRadius: [10, 10, 0, 0]
            }
          },
          {
            type: 'bar',
            emphasis: {
              focus: 'series'
            },
            itemStyle: {
              borderRadius: [10, 10, 0, 0]
            }
          },
          {
            type: 'bar',
            emphasis: {
              focus: 'series'
            },
            itemStyle: {
              borderRadius: [10, 10, 0, 0]
            }
          },
          {
            type: 'bar',
            emphasis: {
              focus: 'series'
            },
            itemStyle: {
              borderRadius: [10, 10, 0, 0]
            }
          },
          {
            type: 'bar',
            emphasis: {
              focus: 'series'
            },
            itemStyle: {
              borderRadius: [10, 10, 0, 0]
            }
          },
          {
            type: 'bar',
            emphasis: {
              focus: 'series'
            },
            itemStyle: {
              borderRadius: [10, 10, 0, 0]
            }
          },
          {
            type: 'bar',
            emphasis: {
              focus: 'series'
            },
            itemStyle: {
              borderRadius: [10, 10, 0, 0]
            }
          },
          {
            type: 'bar',
            emphasis: {
              focus: 'series'
            },
            itemStyle: {
              borderRadius: [10, 10, 0, 0]
            }
          },
          {
            type: 'bar',
            emphasis: {
              focus: 'series'
            },
            itemStyle: {
              borderRadius: [10, 10, 0, 0]
            }
          }
        ]
      }
      myChart.setOption(options)
      // 当浏览器窗口缩放时，图标同时缩放
      window.addEventListener('resize', () => {
        myChart.resize()
      })
    }
  },
  mounted() {
    this.draw()
  },
  created() {
    window.addEventListener('resize', () => {
      this.graphHeight = window.innerHeight - 93
    })
  }
}
</script>

<style lang="less" scoped>
.container {
  height: 100%;
  width: 100%;
  padding: 20px;
  box-sizing: border-box;
}
</style>
