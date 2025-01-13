<template>
  <div class="container" :style="`height: ${containerHeight}px`">
    <el-menu :default-active="activeIndex" router>
      <el-menu-item index="/BackStage/user">
        <i class="el-icon-user"></i>
        <span slot="title">用户管理</span>
      </el-menu-item>
      <el-menu-item index="/BackStage/map">
        <i class="el-icon-tickets"></i>
        <span slot="title">映射管理</span>
      </el-menu-item>
      <el-menu-item index="/BackStage/file">
        <i class="el-icon-folder-opened"></i>
        <span slot="title">文件管理</span>
      </el-menu-item>
      <el-menu-item index="/BackStage/data">
        <i class="el-icon-data-analysis"></i>
        <span slot="title">数据管理</span>
      </el-menu-item>
      <el-menu-item index="/BackStage/log">
        <i class="el-icon-notebook-1"></i>
        <span slot="title">日志记录</span>
      </el-menu-item>
    </el-menu>
    <div>
      <router-view></router-view>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      containerHeight: window.innerHeight - 56,
      activeIndex: '/BackStage/user'
    }
  }
}
</script>

<style lang="less" scoped>
.container {
  position: absolute;
  top: 56px;
}

.el-menu {
  height: 100%;
}
</style>
