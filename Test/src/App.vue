<template>
  <div class="app-container">
    <Header></Header>
    <keep-alive :include="needAlive">
      <router-view v-if="$route.meta.keepAlive"></router-view>
    </keep-alive>
    <router-view v-if="!$route.meta.keepAlive"></router-view>
  </div>
</template>

<script>
import Header from '@/components/header/Header.vue'
import bus from '@/util/eventBus'
export default {
  data() {
    return {
      needAlive: ['my-graph', 'my-overview']
    }
  },
  components: {
    Header
  },
  created() {
    bus.$on('logOut', () => {
      this.needAlive = ''
    })
  }
  // mounted() {
  //   window.addEventListener('beforeunload', () => {
  //     localStorage.removeItem('token')
  //   })
  // }
}
</script>

<style lang="less">
.app-container {
  height: 100%;
  width: 100%;
}

::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 0;
}

::-webkit-scrollbar {
  -webkit-appearance: none;
  width: 10px;
  height: 6px;
}

::-webkit-scrollbar-thumb {
  cursor: pointer;
  border-radius: 10px;
  background: rgba(0, 0, 0, 0.15);
  transition: color 0.2s ease;
}
</style>
