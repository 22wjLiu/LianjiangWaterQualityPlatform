<template>
  <div class="welcome-container" :style="'height: ' + containerHeight + 'px;'">
    <h2 class="login-text">LOGIN</h2>
    <!-- 登录卡片 -->
    <el-card>
      <div class="login-container">
        <el-form :model="form" :rule="rules" class="input" label-width="60px">
          <el-form-item prop="email" label="邮箱：">
            <el-input
              v-model.trim="form.email"
              placeholder="请输入邮箱"
              clearable
              class="input-box"
            ></el-input>
          </el-form-item>
          <el-form-item prop="password" label="密码：">
            <el-input
              v-model.trim="form.password"
              placeholder="请输入密码"
              show-password
              class="input-box"
            ></el-input>
          </el-form-item>
        </el-form>
      </div>
      <div class="findPassword">
        <el-link type="primary" @click="forgotPassword">忘记密码</el-link>
      </div>
      <div class="btn">
        <router-link to="/Register">
          <el-button type="primary">注册</el-button>
        </router-link>
        <el-button type="primary" @click="login">登录</el-button>
      </div>
      <!-- 忘记密码 -->
      <el-dialog
        title="找回密码"
        :visible.sync="verifyDialog"
        width="400px"
        :close-on-click-modal="false"
      >
        <el-form ref="form2" :model="form2" :rules="rules2" status-icon>
          <el-form-item label="邮箱：" prop="email">
            <el-input
              v-model.trim="form2.email"
              placeholder="请输入正确的邮箱地址"
            ></el-input>
          </el-form-item>
          <el-form-item label="验证码" prop="verify">
            <el-input
              v-model="form2.verify"
              placeholder="输入验证码"
            ></el-input>
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
        <div class="nextButton">
          <el-button type="primary" @click="nextStep" :loading="loading1"
            >下一步</el-button
          >
        </div>
      </el-dialog>
      <el-dialog
        title="设置新密码"
        :visible.sync="findDialog"
        width="400px"
        :close-on-click-modal="false"
      >
        <el-form ref="form2" :model="form2" :rules="rules2" status-icon>
          <el-form-item label="系统重置密码：" prop="resetPass">
            <el-input
              v-model.trim="form2.resetPass"
              placeholder="请查看邮箱并输入"
              type="password"
            ></el-input>
          </el-form-item>
          <el-form-item label="新密码：" prop="password">
            <el-input
              v-model.trim="form2.password"
              placeholder="长度在6位以上"
              type="password"
            ></el-input>
          </el-form-item>
          <el-form-item label="确认密码" prop="checkPass">
            <el-input
              v-model.trim="form2.checkPass"
              placeholder="请再次输入秘密"
              type="password"
            ></el-input>
          </el-form-item>
        </el-form>
        <div class="pre-submit">
          <el-button type="primary" @click="preStep">上一步</el-button>
          <el-button type="primary" @click="submit" :loading="loading2"
            >提交</el-button
          >
        </div>
      </el-dialog>
    </el-card>
  </div>
</template>

<script>
import bus from '../util/eventBus'
import { login, verifyEmail, security, updatepass } from '@/api/login'
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

    const validateResetpass = (rule, value, callback) => {
      if (value === '') {
        return callback(new Error('请输入系统重置的密码'))
      } else {
        return callback()
      }
    }
    const validatePass = (rule, value, callback) => {
      if (value === '') {
        return callback(new Error('请输入密码'))
      } else {
        if (this.form2.checkPass !== '') {
          this.$refs.form2.validateField('checkPass')
        }
        return callback()
      }
    }
    const validatePass2 = (rule, value, callback) => {
      if (value === '') {
        return callback(new Error('请再次输入密码'))
      } else {
        if (value !== this.form2.password) {
          return callback(new Error('两次输入的密码不一致'))
        } else {
          return callback()
        }
      }
    }
    return {
      token: '',
      form: {
        email: '',
        password: ''
      },
      form2: {
        email: '',
        resetPass: '',
        password: '',
        checkPass: '',
        verify: ''
      },
      rules: {
        email: [{ required: true, message: '请输入邮箱', trigger: 'blur' }],
        password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
      },
      rules2: {
        email: [
          { required: true, message: '请输入邮箱地址', trigger: 'blur' },
          { validator: validateEmail, trigger: 'blur' }
        ],
        resetPass: [
          { required: true, message: '请输如重置密码', trigger: 'blur' },
          { validator: validateResetpass, trigger: 'blur' }
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
        ]
      },
      verifyButton: '获取验证码',
      verify: '',
      disabled: false,
      containerHeight: window.innerHeight - 57,
      verifyDialog: false,
      findDialog: false,
      loading1: false,
      loading2: false
    }
  },
  methods: {
    login() {
      const formData = new FormData()
      formData.append('Email', this.form.email)
      formData.append('Password', this.form.password)
      login(formData)
        .then(res => {
          const data = res.data
          if (data.code === 200) {
            this.$message.success(data.msg)
            this.token = data.data.token
            localStorage.setItem('token', data.data.token)
            this.$router.push('/overview')
            bus.$emit('activeIndexChange', {
              item: 'overview',
              isLogin: true
            })
          } else {
            this.$message.warning(res.data.msg)
          }
        })
        .catch(err => {
          this.$message.error(err.message)
        })
    },
    forgotPassword() {
      this.verifyDialog = true
    },
    getVerify() {
      this.$refs.form2.validate((boolean, obj) => {
        if (boolean) {
          verifyEmail(this.form2.email)
            .then(res => {
              if (res.data.code === 200) {
                this.$message.success('验证码获取成功，请查看邮箱')
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
    },
    nextStep() {
      if (this.form2.email !== '' && this.form2.verify !== '') {
        this.loading1 = true
      } else if (this.form2.verify === '') {
        this.$message.warning('验证码不能为空')
        return
      }
      const formData = new FormData()
      formData.append('Email', this.form2.email)
      formData.append('Verify', this.form2.verify)
      security(formData)
        .then(res => {
          if (res.data.code === 200) {
            this.loading1 = false
            this.verifyDialog = false
            this.findDialog = true
            setTimeout(() => {
              this.$message.success('密码已重置，请查看邮箱')
            }, 1000)
          } else {
            this.loading1 = false
            this.$message.warning(res.data.msg)
          }
        })
        .catch(err => {
          this.$message.error(err.message)
        })
    },
    preStep() {
      this.form2.resetPass = ''
      this.form2.password = ''
      this.form2.checkPass = ''
      this.verifyDialog = true
      this.findDialog = false
    },
    submit() {
      this.loading2 = true
      const formData = new FormData()
      formData.append('Email', this.form2.email)
      formData.append('Password', this.form2.resetPass)
      login(formData)
        .then(res => {
          if (res.data.code === 200) {
            console.log('get token')
            console.log(res)
            this.updatePass(res.data.data.token)
          } else {
            this.$message.warning(res.data.msg)
          }
        })
        .catch(err => {
          this.$message.error(err.message)
        })
    },
    updatePass(token) {
      const config = {
        headers: {
          Authorization: 'Bearer ' + token
        }
      }
      const data = {
        first: this.form2.resetPass,
        second: this.form2.password
      }
      updatepass(data, config)
        .then(res => {
          if (res.data.code === 200) {
            this.verifyDialog = false
            this.findDialog = false
            this.$message.success('密码修改成功')
          } else {
            this.$message.warning(res.data.msg)
          }
          this.loading2 = false
        })
        .catch(err => {
          this.loading2 = false
          this.$message.error(err.message)
        })
    }
  },

  mounted() {
    window.addEventListener('resize', () => {
      this.containerHeight = window.innerHeight - 57
    })
  }
}
</script>

<style lang="less" scoped>
.welcome-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  position: absolute;
  top: 56px;
  width: 100%;
  background-image: url('../assets/bgc2.jpg');
  background-size: 100% 100%;
  background-position: 0 0;
  background-repeat: no-repeat;
  .login-text {
    margin-top: 10%;
    line-height: 50px;
    text-align: center;
    font-size: 40px;
    font-weight: bolder;
    color: #48889f;
    text-shadow: 2px 2px 4px #000000;
  }
  .login-container {
    .input-box {
      width: 300px;
    }
  }
}

.el-card {
  max-width: 400px;
  background-color: rgba(255, 255, 255, 0.7);
}

.btn {
  display: flex;
  justify-content: center;
}
.el-button {
  margin: 0 10px;
}

.findPassword {
  height: 30px;
  padding: 0 20px;
}

.el-link {
  float: right;
}

.verifyButton {
  width: 100%;
  height: 40px;
  margin-top: 10px;
  .el-button {
    float: right;
  }
}

.nextButton {
  width: 100%;
  display: flex;
  justify-content: center;
}

.pre-submit {
  width: 100%;
  display: flex;
  justify-content: center;
}
</style>
