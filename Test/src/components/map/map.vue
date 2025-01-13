<template>
  <div class="map-container">
    <!-- 切换地图类型 -->
    <div class="btns">
      <el-radio-group v-model="btnLabel" size="small">
        <el-radio-button label="标准地图"></el-radio-button>
        <el-radio-button label="卫星地图"></el-radio-button>
      </el-radio-group>
    </div>
    <!-- 地图容器 -->
    <div class="map">
      <baidu-map
        class="bm-view"
        :center="{ lng: 116.45173120555876, lat: 23.260721350209432 }"
        :zoom="13"
        scroll-wheel-zoom
        inertial-dragging
        :mapType="mapType"
        @ready="mapReady"
      >
        <!-- 比例尺 -->
        <bm-scale anchor="BMAP_ANCHOR_TOP_RIGHT"></bm-scale>
        <!-- 地图导航
        <bm-navigation anchor="BMAP_ANCHOR_BOTTOM_RIGHT"></bm-navigation> -->
        <!-- 地图概览 -->
        <!-- <bm-overview-map
          anchor="BMAP_ANCHOR_BOTTOM_LEFT"
          :isOpen="true"
          :size="overviewSize"
        ></bm-overview-map> -->
        <!-- 自定义覆盖物 -->
        <MyOVerlay
          v-for="item in pointData"
          :key="item.地点"
          :lng="item.经度"
          :lat="item.纬度"
          :text="item.地点"
          :active="item.active"
          system="小时制"
          :site="item.地点"
          :waterQuulity="item.水质等级"
          :pH="item.pH值"
          :oxygen="item.化学需氧量"
          :ammonianitrogen="item.氨氮"
          :phosphorus="item.总磷"
          :dissolvedOxygen="item.溶解氧"
          @mouseover.native="item.active = true"
          @mouseleave.native="item.active = false"
        >
        </MyOVerlay>
      </baidu-map>
    </div>
  </div>
</template>

<script>
import MyOVerlay from '@/components/map/MyOverlay.vue'
import { mapStyle } from '@/assets/custom_map_config'
import { pointData } from '@/assets/data'
export default {
  props: {
    pattern: {
      require: true,
      type: String
    }
  },
  data() {
    return {
      pointData: pointData,
      activeOverlay: 0,
      btnLabel: '标准地图',
      mapType: 'BMAP_NORMAL_MAP',
      overviewSize: {
        width: 200,
        height: 150
      },
      mapStyle
    }
  },
  watch: {
    btnLabel(newVal) {
      newVal === '标准地图'
        ? (this.mapType = 'BMAP_NORMAL_MAP')
        : (this.mapType = 'BMAP_SATELLITE_MAP')
    },
    'overlayList.active'(newVal) {}
  },
  methods: {
    mapReady({ BMap, map }) {
      map.setMapStyle({
        styleId: 'af09776652e1308ed209179eda981505'
      })
    }
  },
  components: {
    MyOVerlay
  }
}
</script>

<style lang="less" scoped>
.map-container {
  width: 100%;
  height: 100%;
  margin: 0 1.5%;
  box-shadow: 0 2px 12px 0 rgb(0 0 0 / 10%);
  position: static;

  .btns {
    position: absolute;
    z-index: 99;
    top: 77px;
    left: 2%;
  }

  .map {
    height: 100%;
    width: 100%;
    .bm-view {
      width: 100%;
      height: 100%;
    }
  }
}

.map-container:hover {
  box-shadow: 0 2px 12px 0 rgb(0 0 0 / 50%);
}
</style>
