<template>
  <div class="container">
    <el-card>
      <div class="selects">
        <el-select
          v-model="selectedList"
          @change="optionsChange"
          multiple
          placeholder="选择数据项"
        >
          <el-option-group
            v-for="group in options"
            :key="group.label"
            :label="group.label"
          >
            <el-option
              v-for="item in group.options"
              :key="item.value"
              :value="item.value"
              :label="item.label"
            ></el-option>
          </el-option-group>
        </el-select>
      </div>
      <div class="system">
        <el-select v-model="system" placeholder="小时制" @change="systemChange">
          <el-option value="小时制" label="小时制"></el-option>
          <el-option value="月度制" label="月度制"></el-option>
        </el-select>
        <el-select
          v-if="system === '小时制'"
          v-model="hourName"
          @change="nameChange"
        >
          <el-option value="海门湾桥闸" label="海门湾桥闸"></el-option>
          <el-option value="汕头练江水站" label="汕头练江水站"></el-option>
        </el-select>
        <el-select
          v-if="system === '月度制'"
          v-model="monthName"
          @change="nameChange"
        >
          <el-option
            v-for="item in nameList"
            :key="item"
            :value="item"
            :label="item"
          ></el-option>
        </el-select>

        <el-date-picker
          v-model="timeRange.start"
          type="datetime"
          start-placeholder="开始日期"
          value-format="yyyy-MM-dd"
          @change="dateChange"
        >
          >
        </el-date-picker>
        <span>--至--</span>
        <el-date-picker
          v-model="timeRange.end"
          type="datetime"
          end-placeholder="结束日期"
          value-format="yyyy-MM-dd"
          @change="dateChange"
        >
        </el-date-picker>
      </div>
    </el-card>
  </div>
</template>

<script>
import { getLineData, getNameList } from '@/api/graphs.js'
import bus from '@/util/eventBus'
export default {
  data() {
    return {
      indexOfTime: 0,
      timeRange: {
        start: '',
        end: ''
      },
      selectedList: ['temperature', 'ph', 'do', 'ec'],
      selectedLabel: [],
      selectedIndex: [],
      accuracyValue: '小时制',
      nameList: [],
      hourName: '海门湾桥闸',
      monthName: '北港河闸',
      system: '小时制',
      options: [
        {
          label: '水质基本参数',
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
            }
          ]
        },
        {
          label: '污染物监测指标',
          options: [
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
              value: 'volatilePenol'
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
              value: 'cumulativeDischarge'
            },
            {
              label: '水流量',
              value: 'waterDischarge'
            },
            {
              label: '总累积流量',
              value: 'totalCumulativeFlow'
            }
          ]
        }
      ],
      start: '',
      end: '',
      config: {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`
        }
      }
    }
  },
  methods: {
    systemChange() {
      this.initData()
    },
    dateChange() {
      this.initData()
    },
    nameChange() {
      this.initData()
      bus.$emit('nameChange', {})
    },
    optionsChange() {
      this.initData()
      // bus.$emit('optionsChange', {})
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
    initData() {
      // 生成请求参数
      let fields = ''
      this.options.forEach(group => {
        group.options.forEach(item1 => {
          this.selectedList.forEach(item2 => {
            if (item1.value === item2) {
              fields += `fields=${item1.label}&`
            }
          })
        })
      })
      const params = fields.slice(0, fields.length - 1)
      // `&start=${this.timeRange.start}` +
      // `&end=${this.timeRange.end}`
      // 发起请求
      const name = this.system === '小时制' ? this.hourName : this.monthName
      getLineData(name, this.system, params, this.config)
        .then(res => {
          console.log(res.data.code, res.data.msg, res)
          let data = res.data.data.resultArr
          // 将表中数值为 0 的字段赋值为 ''
          data = data.map(item => {
            Object.keys(item).forEach(key => {
              if (item[key] === 0) {
                item[key] = ''
              }
            })
            return item
          })
          this.selectedLabel = []
          Object.keys(data[0]).forEach((item, index) => {
            // 计算 data 中 time 字段的索引
            if (item === 'time') {
              this.indexOfTime = index
            }
            // 生成选中的labels列表
            this.options.forEach(group => {
              group.options.forEach(item1 => {
                if (item1.value === item) {
                  this.selectedLabel.push(item1.label)
                }
              })
            })
          })
          // 计算 selecteedList 中元素在 res.data.data.resultArr[0] 中的顺序
          this.selectedIndex = []
          this.selectedList.forEach(item1 => {
            Object.keys(data[0]).forEach((item, index) => {
              if (item === item1) {
                this.selectedIndex.push(index)
              }
            })
          })

          bus.$emit('lineData', {
            data,
            indexOfTime: this.indexOfTime,
            selectedList: this.selectedList,
            selectedIndex: this.selectedIndex,
            selectedLabel: this.selectedLabel
          })
          this.timeRange.start = data[0].time
          this.timeRange.end = data[data.length - 1].time
        })
        .catch(err => {
          console.log(err)
        })
    }
  },
  created() {
    this.getNameList()
    this.initData()
  }
}
</script>

<style lang="less" scoped>
.container {
  margin: 20px 20px 20px 10%;
}
.selects {
  .el-select {
    width: 100%;
    margin-bottom: 20px;
  }
}
.system {
  .el-select {
    margin-right: 20px;
  }
}
</style>
