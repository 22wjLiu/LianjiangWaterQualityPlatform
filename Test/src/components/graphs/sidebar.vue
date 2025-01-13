<template>
  <div class="sidebar">
    <el-menu :index="imgList[0].index" @select="graphChange(imgList[0].index)">
      <el-popover placement="right" title="折线图" trigger="hover">
        <p>水质基本参数均值</p>
        <p>污染物监测指标均值</p>
        <el-divider></el-divider>
        <div class="system">
          <div>
            <el-select
              v-model="system"
              placeholder="小时制"
              @change="systemChange"
            >
              <el-option value="小时制" label="小时制"></el-option>
              <el-option value="月度制" label="月度制"></el-option>
            </el-select>
          </div>
          <div>
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
          </div>
          <el-divider></el-divider>
        </div>
        <div class="archor-container">
          <div v-for="item in labelList" :key="item">
            <el-link type="primary" @click="goAnchor(item)">{{ item }}</el-link>
          </div>
        </div>

        <el-menu-item slot="reference">
          <img
            :src="require(`../../assets/${imgList[0].src}`)"
            :alt="imgList[0].name"
          />
        </el-menu-item>
      </el-popover>
    </el-menu>
    <el-menu :index="imgList[1].index" @select="graphChange(imgList[1].index)">
      <el-popover placement="right" title="柱形图" trigger="hover">
        <el-menu-item slot="reference">
          <img
            :src="require(`../../assets/${imgList[1].src}`)"
            :alt="imgList[1].name"
          />
        </el-menu-item>
        <p>根据《中华人民共和国地表水环境质量标准》</p>
        <p>对各项指标在一段时间内的水质等级进行划分，</p>
        <p>最后计算其时间总和。</p>
      </el-popover>

      <!-- </el-tooltip> -->
    </el-menu>
    <el-menu :index="imgList[2].index" @select="graphChange(imgList[2].index)">
      <el-popover placement="right" title="饼图" trigger="hover">
        <el-menu-item slot="reference">
          <img
            :src="require(`../../assets/${imgList[2].src}`)"
            :alt="imgList[2].name"
          />
        </el-menu-item>
        <p>根据《中华人民共和国地表水环境质量标准》</p>
        <p>对各项指标在一段时间内的水质等级进行划分，</p>
        <p>将不同等级的水质所占时间进行相比。</p>
      </el-popover>
    </el-menu>
  </div>
</template>

<script>
import { getLineData, getNameList } from '@/api/graphs.js'
import bus from '@/util/eventBus'
export default {
  data() {
    return {
      // 请求和渲染图标数据
      indexOfTime: 0,
      timeRange: {
        start: '',
        end: ''
      },
      nameList: [],
      hourName: '海门湾桥闸',
      monthName: '北港河闸',
      system: '小时制',
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
      ],
      controller: null,
      config: {
        signal: null,
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`
        }
      },
      lineData: [],
      labelList: [],
      // 导航跳转数据
      index: 1,
      containerHeight: window.innerHeight - 57,
      isCollapse: true,
      item: null,
      imgList: [
        { src: 'line2.png', name: '折线图', index: 1 },
        { src: 'bar1.png', name: '柱形图', index: 2 },
        { src: 'pie1.png', name: '饼图', index: 3 }
      ]
    }
  },
  methods: {
    // 锚点跳转
    goAnchor(id) {
      const el = document.getElementById(id)
      // document.body.scrollTop = el.offsetTop + 57
      el.scrollIntoView({
        behavior: 'smooth',
        block: 'center'
      })
    },
    graphChange(index) {
      this.index = index
      switch (index) {
        case 1:
          this.imgList[0].src = 'line2.png'
          this.imgList[1].src = 'bar1.png'
          this.imgList[2].src = 'pie1.png'
          break
        case 2:
          this.imgList[0].src = 'line1.png'
          this.imgList[1].src = 'bar2.png'
          this.imgList[2].src = 'pie1.png'
          break
        case 3:
          this.imgList[0].src = 'line1.png'
          this.imgList[1].src = 'bar1.png'
          this.imgList[2].src = 'pie2.png'
          break
      }
      bus.$emit('graphChange', index)
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
        this.labelList.length = 0
        while (flag) {
          if (this.lineData.length === 0 && this.labelList.length === 0) {
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
      this.controller = new AbortController()
      this.config.signal = this.controller.signal
      const name = this.system === '小时制' ? this.hourName : this.monthName
      // 是否已发送 token已过期 的讯息
      let hasWarning = false
      // 循环获取每个数据项的数据
      this.options.forEach(option => {
        const params = 'fields=' + option.label
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

            this.lineData.push(data)
            this.labelList.push(option.label)

            // 计算 data 中 time 字段的索引
            Object.keys(data[0]).forEach((item, index) => {
              if (item === 'time') {
                this.indexOfTime = index
              }
            })

            bus.$emit('lineData', {
              data: this.lineData,
              indexOfTime: this.indexOfTime,
              labelList: this.labelList
            })
            this.timeRange.start = data[0].time
            this.timeRange.end = data[data.length - 1].time
          })
          .catch(err => {
            console.log(err)
            if (!hasWarning) {
              this.$message.error(err.message)
              hasWarning = true
            }
          })
      })
    },
    async getLineListAgain() {
      this.controller.abort()
      bus.$emit('reload', true)
      await this.empty()
      this.getLineList()
    },
    abortController() {
      this.controller.abort()
    }
  },
  created() {
    this.getNameList()
    this.getLineList()
    bus.$on('logOut', () => {
      this.controller.abort()
    })
    window.addEventListener('resize', () => {
      this.containerHeight = window.innerHeight - 57
    })
  }
}
</script>

<style lang="less" scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  justify-content: center;
  position: fixed;
  left: 0;
  top: 20%;
  z-index: 999;
  box-shadow: h-shadow v-shadow blur spread color inset;
}

.item {
  margin: 4px;
}

.system {
  .el-select {
    margin: 5px 0;
  }
}

/deep/.el-popover {
  height: 500px;
  overflow: auto;
}

.archor-container {
  max-height: 300px;

  overflow: auto;
}
</style>
