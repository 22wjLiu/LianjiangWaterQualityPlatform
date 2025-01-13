<template>
  <div class="graph-container" ref="Graph">
    <Sidebar></Sidebar>
    <keep-alive include="lineGraph">
      <component :is="componentName"></component>
    </keep-alive>
  </div>
</template>

<script>
import Sidebar from '@/components/graphs/sidebar.vue'
import lineGraph from '@/components/graphs/lines.vue'
import barGraph from '@/components/graphs/bar.vue'
import pieGraph from '@/components/graphs/pie.vue'
import bus from '@/util/eventBus'
export default {
  name: 'my-graph',
  data() {
    return {
      componentName: 'lineGraph'
    }
  },
  created() {
    bus.$on('graphChange', val => {
      switch (val) {
        case 1:
          this.componentName = 'lineGraph'
          break
        case 2:
          this.componentName = 'barGraph'
          break
        case 3:
          this.componentName = 'pieGraph'
          break
      }
    })
  },
  components: {
    Sidebar,
    lineGraph,
    barGraph,
    pieGraph
  },
  activated() {
    this.$forceUpdate()
  }
}
</script>

<style lang="less" scoped>
.graph-container {
  position: absolute;
  top: 57px;
  width: 100%;
}
</style>
