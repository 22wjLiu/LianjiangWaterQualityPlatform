<template>
  <bm-overlay ref="customOverlay" class="overlay" pane="labelPane" @draw="draw">
    <el-popover placement="top-start" :title="name" trigger="hover">
      <el-divider></el-divider>

      <div class="text item">水质类别：{{ level }}</div>
      <div class="text item">pH: {{ pH }}</div>
      <div class="text item">化学含氧量：{{ oxygen }}</div>
      <div class="text item">氨氮：{{ ammonianitrogen }}</div>
      <div class="text item">总磷：{{ phosphorus }}</div>
      <div class="text item">溶解氧：{{ dissolvedOxygen }}</div>

      <el-divider></el-divider>
      <el-button type="primary" size="small" @click="getDetails"
        >查看断面数据详情</el-button
      >

      <div slot="reference">
        <div :class="{ sample: true, active: active }">
          {{ level }}
        </div>
        <div class="name">
          {{ name }}
        </div>
      </div>
    </el-popover>
  </bm-overlay>
</template>

<script>
import bus from '@/util/eventBus'
export default {
  props: [
    'text',
    'lng',
    'lat',
    'system',
    'active',
    'site',
    'pH',
    'waterQuulity',
    'oxygen',
    'ammonianitrogen',
    'phosphorus',
    'dissolvedOxygen'
  ],
  data() {
    return {
      id: 0
    }
  },
  computed: {
    name() {
      return this.site.slice(3)
    },
    level() {
      let level = ''
      switch (this.waterQuulity) {
        case 1:
          level = 'Ⅰ'
          break
        case 2:
          level = 'Ⅱ'
          break
        case 3:
          level = 'Ⅲ'
          break
        case 4:
          level = 'Ⅳ'
          break
        case 5:
          level = 'Ⅴ'
          break
        default:
          break
      }
      return level
    },
    color() {
      let color = ''
      switch (this.waterQuulity) {
        case 1:
          color = '#10B981'
          break
        case 2:
          color = '#10B981'
          break
        case 3:
          color = '#10B981'
          break
        case 4:
          color = '#10B981'
          break
        case 5:
          color = '#EF4444'
          break
        default:
          break
      }
      return color
    }
  },
  watch: {
    position: {
      handler() {
        this.$refs.customOverlay.reload()
      },
      deep: true
    }
  },
  methods: {
    getDetails() {
      if (this.$route.path === '/map') {
        this.$router.push('/overview')
      }
      bus.$emit('siteChange', {
        system: this.system,
        siteName: this.name
      })
    },
    draw({ el, BMap, map }) {
      const lng = this.lng
      const lat = this.lat
      const pixel = map.pointToOverlayPixel(new BMap.Point(lng, lat))
      el.style.left = pixel.x - 50 + 'px'
      el.style.top = pixel.y - 40 + 'px'
    }
  }
}
</script>

<style lang="less" scoped>
.overlay {
  width: 100px;
  height: 80px;
  position: absolute;
  display: flex;
  flex-direction: column;
  justify-content: center;
}
.sample {
  position: relative;
  width: 40px;
  height: 40px;
  border-radius: 40px;
  box-shadow: 0 0 2px rgb(0, 0, 0);
  line-height: 40px;
  background: v-bind(color);
  color: #fff;
  text-align: center;
  margin-left: 30px;
  z-index: 9;
}

.active {
  background: #6b7280;
  color: #fff;
  z-index: 999;
}

.text {
  font-size: 14px;
}

.item {
  margin-bottom: 10px;
}

.name {
  width: 100%;
  text-align: center;
  line-height: 40px;
  font-size: 14px;
  font-weight: 700;
  color: #374151;
}

.el-button {
  width: 100%;
}
</style>
