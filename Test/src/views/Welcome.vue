<template>
  <div class="welcome-container" :style="'height: ' + containerHeight + 'px;'">
    <div class="title">
      <h1>练江可视化平台</h1>
      <el-divider></el-divider>
      <el-button @click="enterInto"><span>进入</span></el-button>
    </div>
    <div class="content">
      <el-row class="row1" :style="`height: ${rowHeight}px`">
        <el-col :span="12">
          <h4 align="right">上世纪90年代以来，</h4>
          <h4 align="right">随流域内污染问题不断加剧，</h4>
          <h4 align="right">练江成为广东省污染最严重的河流之一。</h4></el-col
        >
        <el-col :span="12" class="rightImg">
          <img src="../assets/barImg.jpg" />
        </el-col>
      </el-row>
      <el-row class="row2" :style="`height: ${rowHeight}px`">
        <el-col class="leftImg" :span="12">
          <img src="../assets/lineImg.jpg" />
        </el-col>
        <el-col :span="12">
          <h3 align="left">在党中央及省委、省政府高度重视和支持下，</h3>
          <h3 align="left">
            汕头市深入践行习近平生态文明思想，坚持科学、依法、精准治污，
          </h3>
        </el-col>
      </el-row>
      <el-row class="row3" :style="`height: ${rowHeight}px`">
        <el-col :span="12">
          <h2 align="right">练江综合整治取得阶段性成效。</h2>
        </el-col>
        <el-col class="rightImg" :span="12">
          <img src="../assets/pieImg.jpg" />
        </el-col>
      </el-row>
    </div>
    <div class="icp-container">
      <p>粤ICP备2021079778号</p>
      <p>粤ICP备2021079778号-1</p>
      <p>@CopyRight STU-ACM</p>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      containerHeight: window.innerHeight - 57,
      rowHeight: (window.innerHeight - 137) / 3
    }
  },
  methods: {
    enterInto() {
      if (localStorage.getItem('token')) {
        this.$router.push('/overview')
      } else {
        this.$router.push('/login')
      }
    }
  },
  mounted() {
    window.addEventListener('resize', () => {
      this.containerHeight = window.innerHeight - 57
    })
  },
  created() {
    window.addEventListener('resize', () => {
      this.containerHeight = window.innerHeight - 57
      this.rowHeight = (window.innerHeight - 137) / 3
    })
  }
}
</script>

<style lang="less" scoped>
.welcome-container {
  position: absolute;
  top: 56px;
  width: 100%;
  color: #f2f6fc;
  background-image: url('../assets/bgc2.jpg');
  background-size: 100% 100%;
  background-position: 0 0;
  background-repeat: no-repeat;
  overflow: auto;

  .title {
    color: #48889f;
    position: absolute;
    top: 20px;
    left: 40px;
    z-index: 9;
    .el-button {
      width: 100px;
    }
  }

  .row1 {
    color: #48889f;

    margin-top: 40px;
  }
  .row2 {
    color: #fff;
  }
  .row3 {
    color: #fff;

    margin-bottom: 40px;
  }
  .leftImg {
    display: flex;
    align-items: flex-end;
  }

  .rightImg {
    display: flex;
    align-items: flex-start;
  }

  .el-col {
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
  }

  img {
    height: auto;
    width: 50%;
    box-shadow: 0 0 5px #5b5353, 0 0 5px #ffffff;
    margin: 0 10%;
  }
  .el-divider {
    margin: 40px 0;
    background-color: #48889f;
  }
  .el-button {
    background-color: #48889f;
    color: #fff;
  }

  .icp-container {
    height: 40px;
    background-color: rgb(0 0 0 / 0);
    display: flex;
    justify-content: center;
    z-index: 999;
    p {
      height: 100%;
      line-height: 100%;
      text-align: center;
      margin: 0 10px;
    }
  }
}
</style>
