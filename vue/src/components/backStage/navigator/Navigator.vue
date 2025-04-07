<template>
  <div class="navigator-container">
    <ul>
      <li @click="changeLogType(0)">数据日志</li>
      <li @click="changeLogType(1)">文件日志</li>
      <li @click="changeLogType(2)">映射日志</li>
      
      <li
        class="nav-slider"
        :style="`left: calc(${activeIndex} * 33.3%); background-color: ${backColor}`"
      ></li>
    </ul>
  </div>
</template>

<style lang="less" scoped>
.navigator-container {
  position: relative;
  min-width: 1300px;
  display: flex;
  justify-content: center;
  margin: 20px 0px;

  ul,
  li {
    list-style: none;
    padding: 0px;
    margin: 0px;
  }

  ul {
    position: relative;
    width: 60%;
    display: flex;
  }

  ul li {
    position: relative;
    margin-bottom: 20px;
    width: 33.3%;
    height: 0px;
    text-align: center;
  }

  ul li:hover {
    border-bottom-color: blue;
  }

  ul .nav-slider {
    position: absolute;
    top: 100%;
    left: 33.3%;
    width: 33.3%;
    height: 3px;
    transition: 0.5s;
  }

  ul li:nth-child(1):hover ~ .nav-slider {
    left: 0px !important;
    background-color: red !important;
  }

  ul li:nth-child(2):hover ~ .nav-slider {
    left: 33.3% !important;
    background-color: #f5a425 !important;
  }

  ul li:nth-child(3):hover ~ .nav-slider {
    left: 66.6% !important;
    background-color: #7df8dd !important;
  }
}
</style>

<script>
export default {
  props: {
    activeIndex: {
      require: true,
      type: Number,
    },
  },
  data() {
    return {
      backColor: "transparent",
    };
  },
  mounted() {
    if (this.activeIndex === 0) {
      this.backColor = "red";
    } else if (this.activeIndex === 1) {
      this.backColor = "#f5a425";
    } else if (this.activeIndex === 2) {
      this.backColor = "#7df8dd";
    }
  },
  methods: {
    changeLogType(index) {
      if (index === 0) {
        this.backColor = "red";
        this.$emit("changeLogType", "data");
      } else if (index === 1) {
        this.backColor = "#f5a425";
        this.$emit("changeLogType", "file");
      } else if (index === 2) {
        this.backColor = "#7df8dd";
        this.$emit("changeLogType", "map");
      }
    },
  },
};
</script>
