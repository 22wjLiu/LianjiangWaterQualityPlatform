<template>
  <div>
    <navigator
      :activeIndex="activeIndex"
      @changeLogType="changeLogType"
    ></navigator>
    <router-view></router-view>
  </div>
</template>

<script>
import navigator from "@/components/backStage/navigator/Navigator";
export default {
  data() {
    return {
      activeIndex: 0,
      options: [
        {
          label: "数据日志",
          value: "data",
        },
        {
          label: "文件日志",
          value: "file",
        },
        {
          label: "映射日志",
          value: "map",
        },
      ],
    };
  },
  components: {
    navigator,
  },
  methods: {
    changeLogType(type) {
      if (type === "data") {
        this.$router.push("/backStage/log/dataLog");
        this.activeIndex = 0;
      } else if (type === "file") {
        this.$router.push("/backStage/log/fileLog");
        this.activeIndex = 1;
      } else if (type === "map") {
        this.$router.push("/backStage/log/mapLog");
        this.activeIndex = 2;
      } else {
        this.$router.push("/backStage/log/dataLog");
        this.activeIndex = 0;
      }
    },
  },
  created() {
    if (this.$route.path === "/backStage/log/dataLog") {
      this.activeIndex = 0;
    } else if (this.$route.path === "/backStage/log/fileLog") {
      this.activeIndex = 1;
    } else if (this.$route.path === "/backStage/log/mapLog") {
      this.activeIndex = 2;
    }
  },
};
</script>
