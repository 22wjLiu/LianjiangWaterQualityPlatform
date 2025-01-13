<template>
  <div class="container">
    <h3>历史预测</h3>
    <el-divider></el-divider>
    <div class="select-container">
      <!-- <span>根据过去</span> -->
      <p>根据过去：</p>
      <el-select v-model="pastOption" placeholder="请选择">
        <el-option :value="1" label="一个月"></el-option>
        <el-option :value="2" label="两个月"></el-option>
        <el-option :value="3" label="三个月"></el-option>
      </el-select>
    </div>
    <div class="select-container">
      <!-- <span>预测未来</span> -->
      <p>预测未来：</p>
      <el-select v-model="futureOption" placeholder="请选择">
        <el-option :value="1" label="一个月"></el-option>
        <el-option :value="2" label="两个月"></el-option>
        <el-option :value="3" label="三个月"></el-option>
      </el-select>
    </div>
    <el-divider></el-divider>
    <el-button type="primary" @click="submit" :disabled="btnDisabled">{{
      btnMsg
    }}</el-button>
  </div>
</template>

<script>
import { arr } from '@/assets/data'
import { forecast } from '@/api/forecast'
import bus from '@/util/eventBus'
export default {
  data() {
    return {
      pastOption: 1,
      futureOption: 1,
      COcr: '',
      conductivity: '',
      ammoniaNitrogen: '',
      btnMsg: '点击预测',
      isloading: false,
      btnDisabled: false,
      datasetList: ['Temperature', 'PH', 'DO', 'Turbidity', 'EC', 'CODcr'],
      resLength: 0,
      dataset: []
    }
  },
  methods: {
    submit() {
      bus.$emit('showLoading')
      this.btnMsg = '预测中'
      this.btnDisabled = true
      // const split = this.pastOption / (this.pastOption + this.futureOption)
      const split = 0.7
      const config = {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`
        }
      }
      let index = 0
      this.datasetList.forEach(item => {
        const params = `split=${split}&dataset=${item}`
        forecast(params, config)
          .then(res => {
            const data = res.data.data.data

            if (index === 0) {
              if (res.data.code === 201) {
                this.$message.warning('登录已过期，请重新登录')
                this.$router.push('/login')
                return
              }
              this.resLength = data.length
              const interval = Math.trunc(
                (24 * 30 * 60 * 60 * this.futureOption) / this.resLength
              )
              let lastDate = this.$moment(arr[arr.length - 1].时间)
              // result = this.lastDate
              for (let i = 0; i < this.resLength; i++) {
                const obj = {
                  时间: '',
                  氨氮: 0
                }
                lastDate = lastDate.add(interval, 'seconds')
                obj.时间 = lastDate.format()
                this.dataset.push(obj)
              }
            }

            switch (item) {
              case 'Temperature':
                this.dataset.forEach((item, index) => {
                  item.水温 = data[index]
                })
                break
              case 'PH':
                this.dataset.forEach((item, index) => {
                  item.pH = data[index]
                })
                break
              case 'DO':
                this.dataset.forEach((item, index) => {
                  item.溶解氧 = data[index]
                })
                break
              case 'Turbidity':
                this.dataset.forEach((item, index) => {
                  item.浊度 = data[index]
                })
                break
              case 'EC':
                this.dataset.forEach((item, index) => {
                  item.电导率 = data[index]
                })
                break
              case 'CODcr':
                this.dataset.forEach((item, index) => {
                  item.CODcr = data[index]
                })
                break
              default:
                break
            }
            if (index === this.datasetList.length - 1) {
              this.btnMsg = '点击预测'
              this.btnDisabled = false
              bus.$emit('dataChange', {
                dataset: this.dataset
              })
            }
            index++
          })
          .catch(err => {
            console.log(err)
            // if (!hasWarning) {
            //   this.$message.error(err.message)
            //   hasWarning = true
            // }
          })
      })
    }
  },
  created() {
    bus.$on('listMsg', val => {
      this.lastDate = val.lastDate
    })
  },
  mounted() {
    this.submit()
  }
}
</script>

<style lang="less" scoped>
.select-container {
  // span {
  //   display: inline-block;
  //   width: 60px;
  //   font-size: 14px;
  //   line-height: 32px;
  // }
  .el-select {
    display: block;
    width: 200px;
    margin-bottom: 10px;
  }
}
</style>
