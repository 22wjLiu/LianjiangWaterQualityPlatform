<template>
  <div>
    <el-menu-item index="/Logout" @click="logOut" style="width: 110px"
      >退出登录</el-menu-item
    >
  </div>
</template>

<script>
import bus from '@/util/eventBus'
export default {
  methods: {
    logOut() {
      bus.$emit('logOut', true)
    }
  }
}
</script>

<style></style>
