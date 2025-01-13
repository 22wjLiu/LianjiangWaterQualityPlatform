<template>
  <el-card class="box-card" :style="`height: ${cardHeight}px;`">
    <slot name="graph"></slot>
    <slot name="selector"></slot>
  </el-card>
</template>

<script>
export default {
  data() {
    return {
      cardHeight: window.innerHeight - 93
    }
  },
  created() {
    window.addEventListener('resize', () => {
      this.cardHeight = window.innerHeight - 93
    })
  }
}
</script>

<style>
.el-card {
  width: 92%;
  float: right;
  margin: 18px;
}
</style>
