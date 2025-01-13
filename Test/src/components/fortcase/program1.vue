<template>
  <div>
    <el-divider></el-divider>
    <h4>输入数据：</h4>
    <el-form :model="form" :rules="rules" class="input">
      <el-form-item prop="temp" label="温度："
        ><el-input v-model.number="form.temp" class="input-box"></el-input
      ></el-form-item>
      <el-form-item prop="pH" label="pH："
        ><el-input v-model.number="form.pH" class="input-box"></el-input
      ></el-form-item>
      <el-form-item prop="DO" label="溶解氧："
        ><el-input v-model.number="form.DO" class="input-box"></el-input
      ></el-form-item>
      <el-form-item prop="turbidity" label="浊度："
        ><el-input v-model.number="form.turbidity" class="input-box"></el-input
      ></el-form-item>
    </el-form>
    <el-divider></el-divider>
    <el-button type="primary" @click="submit" :disabled="btnDisabled">{{
      btnMsg
    }}</el-button>
    <el-divider></el-divider>
    <h4>预测结果：</h4>
    <div v-loading="isloading" element-loading-text="预测中，请耐心等待">
      <el-form label-width="80px" label-position="left">
        <el-form-item label="COcr：">
          <el-input v-model="COcr" disabled></el-input>
        </el-form-item>
        <el-form-item label="电导率：">
          <el-input v-model="conductivity" disabled></el-input>
        </el-form-item>
        <el-form-item label="氨氮：">
          <el-input v-model="ammoniaNitrogen" disabled></el-input>
        </el-form-item>
      </el-form>
    </div>
    <div style="height: 20px"></div>
  </div>
</template>

<script>
import { forecast } from '@/api/forecast'
export default {
  data() {
    return {
      form: {
        temp: '',
        pH: '',
        DO: '',
        turbidity: ''
      },
      COcr: '',
      conductivity: '',
      ammoniaNitrogen: '',
      btnMsg: '点击预测',
      isloading: false,
      btnDisabled: false,
      rules: {
        temp: [{ required: true, message: '请输入', trigger: 'blur' }],
        pH: [{ required: true, message: '请输入', trigger: 'blur' }],
        DO: [{ required: true, message: '请输入', trigger: 'blur' }],
        turbidity: [{ required: true, message: '请输入', trigger: 'blur' }]
      }
    }
  },
  methods: {
    submit() {
      this.btnMsg = '预测中'
      this.isloading = true
      this.btnDisabled = true
      const params = `Temperature=${this.form.temp}&PH=${this.form.pH}&Turbidity=${this.form.turbidity}&DO=${this.form.DO}`
      const config = {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`
        }
      }
      forecast(params, config)
        .then(res => {
          console.log(res)
          this.btnMsg = '点击预测'
          this.isloading = false
          this.btnDisabled = false
          if (res.data.code === 200) {
            const data = res.data.data.data
            this.COcr = data[0]
            this.conductivity = data[1]
            this.ammoniaNitrogen = data[2]
          } else {
            this.$message.warning(res.data.msg)
          }
        })
        .catch(err => {
          this.btnMsg = '点击预测'
          this.isloading = false
          this.btnDisabled = false
          console.log(err)
          this.$message.error(err.message)
        })
    }
  }
}
</script>

<style lang="less" scoped></style>
