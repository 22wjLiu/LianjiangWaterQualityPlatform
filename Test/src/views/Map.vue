<template>
  <div class="Map-container" :style="`height: ${containerHeight}px`">
    <Map></Map>
  </div>
</template>

<script>
import Map from '@/components/map/map.vue'
export default {
  data() {
    return {
      containerHeight: window.innerHeight - 57
    }
  },
  components: {
    Map
  },
  created() {
    window.addEventListener('resize', () => {
      this.containerHeight = window.innerHeight - 57
    })
  }
}
</script>

<style lang="less" scoped>
.Map-container {
  width: 100%;
  margin-top: 56px;
}
</style>
