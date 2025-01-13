<template>
  <div class="header-container">
    <el-menu
      :default-active="activeIndex"
      class="el-menu-demo menu-container"
      mode="horizontal"
      background-color="#375577"
      text-color="#C0C4CC"
      active-text-color="#fff"
      router
    >
      <div class="welcome-page-mark">
        <component :is="logStatus"></component>
      </div>
      <div class="items-container">
        <el-menu-item index="/backStage">后台管理</el-menu-item>
        <el-menu-item index="/file">文件</el-menu-item>
        <el-menu-item index="/graph">图表</el-menu-item>
        <el-menu-item index="/overview">首页</el-menu-item>
      </div>
    </el-menu>
  </div>
</template>

<script>
import bus from '../../util/eventBus'
import logoutTag from '@/components/header/LogoutTag'
import welcomeTag from '@/components/header/WelcomeTag'
export default {
  data() {
    return {
      activeIndex: '/welcome',
      logStatus: 'welcomeTag'
    }
  },
  watch: {
    $route: {
      handler(newVal) {
        const path = newVal.path
        if (path === '/graph' || path === '/file' || path === '/overview') {
          this.logStatus = 'logoutTag'
        }
        if (path === '/welcome') {
          this.logStatus = 'welcomeTag'
          this.activeIndex = '/welcome'
        }
        if (path === '/overview') {
          this.activeIndex = '/overview'
        }
        if (path === '/logout') {
          this.$router.push('/welcome')
        }
      }
    }
  },
  created() {
    bus.$on('activeIndexChange', val => {
      this.activeIndex = val.item
      if (val.isLogin) {
        this.logStatus = 'logoutTag'
      }
    })
    bus.$on('logOut', () => {
      localStorage.removeItem('token')
      this.logStatus = 'welcomeTag'
      this.$router.push('/Welcome')
      this.activeIndex = '/Welcome'
    })
    this.activeIndex = location.hash.slice(1)
    if (this.activeIndex !== '/Welcome' && '/Login') {
      this.logStatus = 'logoutTag'
    }
  },
  components: {
    logoutTag,
    welcomeTag
  }
}
</script>

<style lang="less" scoped>
.header-container {
  width: 100%;
  position: fixed;
  top: 0;
  z-index: 999;
  .el-menu,
  .el-menu-item,
  .welcome-page-mark {
    border: 0;
  }
  .menu-container {
    display: flex;
    .welcome-page-mark {
      width: 20%;
    }
    .items-container {
      width: 80%;
      position: relative;
      display: flex;
      flex-direction: row-reverse;
    }
    .el-menu-item {
      height: 56px;
      display: flex;
      justify-content: center;
      width: 120px;
    }
    .is-active {
      height: 56px;
    }
  }
}
</style>
