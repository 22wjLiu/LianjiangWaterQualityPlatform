<template>
  <div class="re-container" :style="'height: ' + containerHeight + 'px;'">
    <h2 class="login-text">REGISTER</h2>
    <el-card class="box-card">
      <el-form
        ref="form"
        :model="form"
        label-width="80px"
        label-position="left"
        :rules="rules"
        status-icon
      >
        <el-form-item label="邮箱：" prop="email">
          <el-input
            v-model.trim="form.email"
            placeholder="请输入正确的邮箱地址"
          ></el-input>
        </el-form-item>
        <el-form-item label="密码：" prop="password">
          <el-input
            v-model.trim="form.password"
            placeholder="长度在6位以上"
            type="password"
          ></el-input>
        </el-form-item>
        <el-form-item label="确认密码" prop="checkPass">
          <el-input
            v-model.trim="form.checkPass"
            placeholder="请再次输入秘密"
            type="password"
          ></el-input>
        </el-form-item>
        <el-form-item label="昵称：" prop="name">
          <el-input
            v-model.trim="form.name"
            placeholder="长度在20位以内"
          ></el-input>
        </el-form-item>
        <el-form-item label="验证码" prop="verify">
          <el-input v-model="form.verify" placeholder="输入验证码"></el-input>
        </el-form-item>
      </el-form>
      <div class="verifyButton">
        <el-button
          type="info"
          @click="getVerify"
          size="mini"
          :disabled="disabled"
          >{{ verifyButton }}</el-button
        >
      </div>
      <div class="submit">
        <el-button type="primary" @click="submit">注册</el-button>
      </div>
    </el-card>
  </div>
</template>

<script>
import { getVerify, submit } from '@/api/register'
import Vue from 'vue'
export default {
  data() {
    const validateEmail = (rule, value, callback) => {
      const reg =
        /^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z0-9]{2,6}$/
      if (value === '') {
        return callback(new Error('请输入邮箱地址'))
      } else {
        if (!reg.test(value)) {
          return callback(new Error('请输入正确的邮箱地址'))
        } else {
          callback()
        }
      }
    }
    const validatePass = (rule, value, callback) => {
      if (value === '') {
        return callback(new Error('请输入密码'))
      } else {
        if (this.form.checkPass !== '') {
          this.$refs.form.validateField('checkPass')
        }
        return callback()
      }
    }
    const validatePass2 = (rule, value, callback) => {
      if (value === '') {
        return callback(new Error('请再次输入密码'))
      } else {
        if (value !== this.form.password) {
          return callback(new Error('两次输入的密码不一致'))
        } else {
          return callback()
        }
      }
    }

    return {
      containerHeight: window.innerHeight - 57,
      verifyStatus: '',
      submitStatus: '',
      disabled: false,
      verifyButton: '获取验证码',
      form: {
        email: '',
        password: '',
        checkPass: '',
        name: '',
        verify: ''
      },
      rules: {
        email: [
          {
            required: true,
            message: '请输入邮箱地址',
            trigger: 'blur'
          },
          {
            validator: validateEmail,
            trigger: 'blur'
          }
        ],
        password: [
          {
            required: true,
            message: '请输入密码',
            trigger: 'blur'
          },
          {
            validator: validatePass,
            trigger: 'blur'
          }
        ],
        checkPass: [
          {
            required: true,
            message: '请再次输入密码',
            trigger: 'blur'
          },
          {
            validator: validatePass2,
            trigger: 'blur'
          }
        ],
        name: [
          {
            required: true,
            message: '请输入昵称',
            trigger: 'blur'
          },
          {
            min: 1,
            max: 20,
            message: '昵称长度在20位以内',
            trigger: 'blur'
          }
        ]
      }
    }
  },
  methods: {
    getVerify() {
      this.$refs.form.validate((boolean, obj) => {
        if (boolean) {
          getVerify(this.form.email)
            .then(res => {
              this.verifyStatus = res.data.code
              if (this.verifyStatus === 200) {
                this.$message.info('验证码获取成功，请查看邮箱')
                this.setTime()
              } else {
                this.$message.warning(res.data.msg)
              }
            })
            .catch(err => {
              this.$message.error(err.message)
            })
        }
      })
    },
    submit() {
      const formData = new FormData()
      formData.append('Email', this.form.email)
      formData.append('Password', this.form.password)
      formData.append('Name', this.form.name)
      formData.append('Verify', this.form.verify)

      submit(formData)
        .then(res => {
          this.submitStatus = res.data.code
          if (this.submitStatus === 200) {
            Vue.prototype.$token = res.data.data.token
            localStorage.setItem('token', res.data.data.token)
            this.$router.push('/Grpah')
          } else {
            this.$message.warning(res.data.msg)
          }
        })
        .catch(err => {
          this.$message.error(err.message)
        })
    },
    setTime() {
      this.disabled = true
      let surplus = 300
      const id = setInterval(() => {
        this.verifyButton = `${surplus--} s`
        if (surplus === 0) {
          this.verifyButton = '重新获取验证码'
          this.disabled = false
          clearInterval(id)
        }
      }, 1000)
    }
  },
  created() {
    window.addEventListener('resize', () => {
      this.containerHeight = window.innerHeight - 57
    })
  }
}
</script>

<style lang="less" scoped>
.re-container {
  position: absolute;
  top: 56px;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  background-color: #f2f6fc;
  background-image: url('../assets/bgc2.jpg');
  background-size: 100% 100%;
  background-position: 0 0;
  background-repeat: no-repeat;
  .login-text {
    margin-top: 5%;
    line-height: 50px;
    text-align: center;
    font-size: 40px;
    font-weight: bolder;
    color: #48889f;
    text-shadow: 2px 2px 4px #000000;
  }
  .el-card {
    width: 400px;
    background-color: rgba(255, 255, 255, 0.7);
  }
  .submit {
    display: flex;
    justify-content: center;
  }
  .verifyButton {
    width: 100%;
    height: 40px;
    margin-top: 10px;
    .el-button {
      float: right;
    }
  }
}
</style>
