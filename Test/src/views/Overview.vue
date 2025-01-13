<template>
  <div class="ov-container" :style="`height: ${containerHeight}px`">
    <div class="row" style="margin: 0.5% 0">
      <Map pattern="overview"></Map>
    </div>
    <div class="row">
      <div id="lineGraph" class="col card" style="width: 55%">
        <LineGraph pattern="overview"></LineGraph>
      </div>
      <div class="col card" style="width: 40%">
        <BarGraph pattern="overview"></BarGraph>
      </div>
    </div>
    <div class="row">
      <div class="col card" style="width: 55%; display: flex">
        <Graph></Graph>
        <InputCard v-show="show"></InputCard>
      </div>
      <div class="col card" style="width: 40%">
        <PieGraph pattern="overview"></PieGraph>
      </div>
    </div>
    <div style="height: 1px"></div>
  </div>
</template>

<script>
import LineGraph from '@/components/graphs/liness.vue'
import BarGraph from '@/components/graphs/bar.vue'
import PieGraph from '@/components/graphs/pie.vue'
import Graph from '@/components/fortcase/graph.vue'
import InputCard from '@/components/fortcase/form.vue'
import Map from '@/components/map/map.vue'
export default {
  name: 'my-overview',
  data() {
    return {
      containerHeight: window.innerHeight - 57,
      show: false
    }
  },
  components: {
    LineGraph,
    BarGraph,
    PieGraph,
    Map,
    Graph,
    InputCard
  }
}
</script>

<style lang="less" scoped>
.ov-container {
  width: 100%;
  padding-top: 56px;
  .row {
    width: 100%;
    height: 48%;
    display: flex;
    justify-content: space-evenly;
    box-sizing: border-box;
    margin: 1% 0;
  }
  .col {
    height: 100%;
    background-color: rgb(248, 249, 251);
  }

  .card {
    border-radius: 4px;
    box-shadow: 0 2px 12px 0 rgb(0 0 0 / 10%);
  }

  .card:hover {
    box-shadow: 0 4px 24px 0 rgb(0 0 0 / 50%);
  }
}
</style>
