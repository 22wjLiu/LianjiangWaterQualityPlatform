<template>
  <div class="container" ref="liness"></div>
</template>

<script>
import { getLineData, getNameList } from '@/api/graphs.js'
import bus from '@/util/eventBus'
export default {
  data() {
    return {
      graph: null,
      myChart: null,
      lineData: [],
      indexOfTime: -1,
      options: [
        {
          label: '水温',
          value: 'temperature'
        },
        {
          label: 'pH',
          value: 'ph'
        },
        {
          label: '溶解氧',
          value: 'do'
        },
        {
          label: '电导率',
          value: 'ec'
        },
        {
          label: '浊度',
          value: 'turbidity'
        },
        {
          label: '高锰酸盐指数',
          value: 'codmii'
        },
        {
          label: '氨氮',
          value: 'nh_n'
        },
        {
          label: '总磷',
          value: 'tp'
        },
        {
          label: '总氮',
          value: 'tn'
        },
        {
          label: 'CODcr',
          value: 'codcr'
        },
        {
          label: '氰化物',
          value: 'cn'
        },
        {
          label: '挥发酚',
          value: 'volatile_penol'
        },
        {
          label: '六价铬',
          value: 'cr'
        },
        {
          label: '铜',
          value: 'cu'
        },
        {
          label: '锌',
          value: 'zn'
        },
        {
          label: '铅',
          value: 'pb'
        },
        {
          label: '镉',
          value: 'cd'
        },
        {
          label: '阴离子表面活性剂',
          value: 'las'
        },
        {
          label: '硫化物',
          value: 'sox'
        },
        {
          label: '累计流量',
          value: 'cumulative_discharge'
        },
        {
          label: '水流量',
          value: 'water_discharge'
        },
        {
          label: '总累积流量',
          value: 'total_cumulative_flow'
        }
      ],
      hourName: '海门湾桥闸',
      monthName: '北港河闸',
      system: '小时制',
      config: {
        signal: null,
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`
        }
      }
    }
  },
  methods: {
    async getLineListAgain() {
      await this.empty()
      const el = document.getElementById('lineGraph')
      el.scrollIntoView({
        behavior: 'smooth',
        block: 'start'
      })
      this.getLineList()
      this.$message.info('请稍等')
    },
    systemChange() {
      this.getLineListAgain()
    },
    nameChange() {
      this.getLineListAgain()
    },
    empty() {
      return new Promise(resolve => {
        let flag = true
        this.lineData.length = 0
        while (flag) {
          if (this.lineData.length === 0) {
            flag = false
            resolve()
          }
        }
      })
    },

    getNameList() {
      getNameList('站名映射', this.config)
        .then(res => {
          this.nameList = res.data.data.keys
        })
        .catch(err => {
          console.log(err)
        })
    },
    getLineList() {
      const name = this.system === '小时制' ? this.hourName : this.monthName
      let params = ''
      this.options.forEach(item => {
        params += 'fields=' + item.label + '&'
      })
      const last = params.lastIndexOf('&')
      params = params.slice(0, last)

      getLineData(name, this.system, params, this.config)
        .then(res => {
          let data = res.data.data.resultArr

          // 将表中数值为 0 的字段赋值为 ''
          data = data.map(item => {
            Object.keys(item).forEach(key => {
              if (item[key] <= 0) {
                item[key] = ''
              }
            })
            return item
          })

          this.lineData = data

          // 计算 data 中 time 字段的索引
          Object.keys(data[0]).forEach((item, index) => {
            if (item === 'time') {
              this.indexOfTime = index
            }
          })

          this.myChart.hideLoading()
          this.draw(this.buildSeries(), this.buildSelected())
        })
        .catch(err => {
          console.log(err)
        })
    },
    initChart() {
      this.graph = this.$refs.liness
      this.myChart = this.$echarts.init(this.graph)
    },
    // 构造MyCharts.options.series
    buildSeries() {
      const series = []
      Object.keys(this.lineData[0]).forEach((key, index) => {
        this.options.forEach(item => {
          if (item.value === key) {
            series.push({
              type: 'line',
              symbol: 'none',
              smoott: true,
              sampling: 'lttb',
              lineStyle: {
                width: 1
              },
              name: item.label,
              encode: {
                x: this.indexOfTime,
                y: index
              }
            })
          }
        })
      })
      return series
    },
    buildSelected() {
      const selected = Object.create(Object)
      this.options.forEach(item => {
        const label = item.label

        if (
          item.label === '累计流量' ||
          item.label === '总累积流量' ||
          item.label === '电导率' ||
          item.label === '水流量' ||
          item.label === '浊度'
        ) {
          selected[label] = false
        } else {
          selected[label] = true
        }
      })
      return selected
    },
    draw(series, selected) {
      const options = {
        dataset: [
          {
            source: this.lineData
          }
        ],
        legend: {
          type: 'scroll',
          // orient: 'vertical',
          bottom: '15%',
          left: 'center',
          // 如果series 对象有name 值，则 legend可以不用写data
          // 修改图例组件 文字颜色
          textStyle: {
            color: '#4c9bfd',
            fontSize: 10
          },
          // selected: {
          //   pH: true,
          //   溶解氧: true,
          //   水温: true,
          //   浊度: true
          // }
          selected: selected
        },
        title: {
          text: this.hourName,
          left: 'center',
          textStyle: {
            fontFamily: 'SimSun',
            fontSize: '16'
          }
        },
        // 提示框组件
        tooltip: {
          // 触发类型：坐标轴触发
          trigger: 'axis'
        },
        grid: {
          top: '5%',
          left: '1%',
          right: '1%',
          bottom: '25%',
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
      this.myChart.setOption(options, true)
      window.addEventListener('resize', () => {
        this.myChart.resize()
      })
    }
  },
  created() {
    this.getNameList()
    this.getLineList()
  },
  mounted() {
    this.initChart()
    this.draw()
    this.myChart.showLoading()
    bus.$on('siteChange', val => {
      this.myChart.showLoading()
      this.system = val.system
      val.system === '月度制'
        ? (this.monthName = val.siteName)
        : (this.hourName = val.siteName)
      this.getLineListAgain()
    })
    bus.$on('hideLoading', () => {
      this.myChart.hideLoading()
    })
  },
  activated() {
    this.myChart.resize()
  }
}
</script>

<style lang="less" scoped>
.container {
  height: 100%;
  width: 100%;
  box-sizing: border-box;
  padding: 20px;
}
</style>
