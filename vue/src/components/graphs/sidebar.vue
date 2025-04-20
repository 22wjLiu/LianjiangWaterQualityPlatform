<template>
  <div class="sidebar">
    <el-menu :index="imgList[0].index" @select="graphChange(imgList[0].index)">
      <el-popover placement="right" title="水质基本参数均值图" trigger="click">
        <div class="system">
          <div>
            <el-select
              v-model="system"
              placeholder="小时制"
              @change="systemChange"
            >
              <el-option value="小时制" label="小时制"></el-option>
              <el-option value="月度制" label="月度制"></el-option>
            </el-select>
          </div>
          <div>
            <el-select v-model="stationName">
              <el-option
                v-for="item in stationOptions"
                :key="item.name"
                :value="item.name"
                :label="item.name"
              >
              </el-option>
            </el-select>
          </div>
          <div>
            <el-date-picker
            v-model="startTime"
            type="datetime"
            placeholder="选择日期时间"
            clearable
          >
          </el-date-picker>
          </div>
          <div>
            <el-date-picker
            v-model="endTime"
            type="datetime"
            placeholder="选择日期时间"
            clearable
          >
          </el-date-picker>
          </div>
          <el-button type="primary" @click="getLineListAgain"> 查询 </el-button>
          <el-divider></el-divider>
        </div>
        <div class="archor-container">
          <div v-for="item in options" :key="item.value">
            <el-link type="primary" @click="goAnchor(item.value)">{{
              item.key
            }}</el-link>
          </div>
        </div>
        <el-menu-item slot="reference">
          <img
            :src="require(`../../assets/${imgList[0].src}`)"
            :alt="imgList[0].name"
          />
        </el-menu-item>
      </el-popover>
    </el-menu>

    <el-menu :index="imgList[1].index" @select="graphChange(imgList[1].index)">
      <el-popover placement="right" title="饼图" trigger="hover">
        <el-menu-item slot="reference">
          <img
            :src="require(`../../assets/${imgList[1].src}`)"
            :alt="imgList[1].name"
          />
        </el-menu-item>
        <p>根据《中华人民共和国地表水环境质量标准》</p>
        <p>对各项指标在一段时间内的水质等级进行划分，</p>
        <p>将不同等级的水质所占时间进行相比。</p>
      </el-popover>
    </el-menu>
  </div>
</template>

<script>
import { getLineData, getFieldsLineData, getStationName } from "@/api/graphs.js";
import { fullFormatTime } from "@/util/timeFormater";
import { level } from "@/assets/level.js";
import bus from "@/util/eventBus";
export default {
  data() {
    return {
      // 请求和渲染图标数据
      indexOfTime: 0,
      nameList: [],
      stationName: "",
      system: "小时制",
      startTime: "",
      endTime: "",
      options: [],
      levelLabels: ["I级", "II级", "III级", "IV级", "V级"],
      stationNames: [],
      stationOptions: [],
      controller: null,
      levels: level,
      lineData: [],
      pieLineData: [],
      // 导航跳转数据
      index: 1,
      containerHeight: window.innerHeight - 57,
      isCollapse: true,
      item: null,
      imgList: [
        { src: "line2.png", name: "折线图", index: 1 },
        { src: "pie1.png", name: "饼图", index: 2 },
      ],
    };
  },
  methods: {
    fullFormatTime,
    // 锚点跳转
    goAnchor(id) {
      const el = document.getElementById(id);
      // document.body.scrollTop = el.offsetTop + 57
      el.scrollIntoView({
        behavior: "smooth",
        block: "center",
      });
    },
    graphChange(index) {
      this.index = index;
      switch (index) {
        case 1:
          this.imgList[0].src = "line2.png";
          this.imgList[1].src = "pie1.png";
          break;
        case 2:
          this.imgList[0].src = "line1.png";
          this.imgList[1].src = "pie2.png";
          break;
      }
      bus.$emit("graphChange",index);
    },
    systemChange() {
      this.stationOptions = this.stationNames.filter(item => item.system === this.system);
      this.stationName = "";
    },
    async initStationOptions() {
      try {
        const res = await getStationName();
        if (res.code === 200) {
          this.stationNames = res.data.stationNames;
          if (this.stationNames.length > 0) {
            this.stationName = this.stationNames[0].name;
            this.stationOptions = this.stationNames.filter(item => item.system === this.system);
          }
        }
      } catch (err) {
        console.error("加载失败", err);
      }
    },
    getLineList() {
      this.controller = new AbortController();
      if (this.stationName == "") {
        this.$message.error("站名不能为空");
      }
      let params = "";
      if (this.startTime) {
        params += `start=${this.fullFormatTime(this.startTime)}&`;
      }
      if (this.endTime) {
        params += `end=${this.fullFormatTime(this.endTime)}`;
      }

      bus.$emit("showLoading");

      getFieldsLineData(this.stationName, this.system, params)
        .then((res) => {
          if (res.code === 200) {
            this.startTime = res.data.startTime;
            this.endTime = res.data.endTime;
            this.options = res.data.mapInfos;
            let data = res.data.resultArr;
            if (data.length) {
              // 将表中数值为 0 的字段赋值为 ''
              data = data.map((row) =>
                row.map((cell) => {
                  for (const key in cell) {
                    if (cell[key] <= 0) {
                      cell[key] = 0;
                    }
                  }
                  return cell;
                })
              );

              this.lineData = data;
              // 计算 data 中 time 字段的索引
              this.indexOfTime = Object.keys(data[0][0]).indexOf("time");
            }
            bus.$emit("lineData", {
              data: this.lineData,
              indexOfTime: this.indexOfTime,
              options: this.options,
            });
          }
        })
        .catch((err) => {
          this.$message.error(err.message);
          console.log(err);
        });

        getLineData(this.stationName, this.system, params)
        .then((res) => {
          if (res.code === 200) {
            let data = res.data.resultArr;
            if (data.length) {
              // 将表中数值为 0 的字段赋值为 ''
              data = data.map((item) => {
                Object.keys(item).forEach((key) => {
                  if (item[key] <= 0) {
                  item[key] = 0;
                  } 
                });
               return item;
              });
              this.pieLineData = data;
            }
          }
        })
        .catch((err) => {
          this.$message.error(err.message);
          console.log(err);
        });
    },
    async getLineListAgain() {
      this.controller.abort();
      bus.$emit("reload", true);
      this.getLineList();
    },
    getLevelIndex(val, thresholds, reverse = false) {
      for (let i = 0; i < thresholds.length; i++) {
        if (reverse ? val >= thresholds[i] : val <= thresholds[i]) {
          return i;
        }
      }
      return 4;
    },
    classify(list) {
      const reverseFields = new Set(["溶解氧"]);
      const classified = list.map((item) => {
        let worstLevel = -1;
        for (const key in this.levels) {
          const val = parseFloat(item[key]);
          if (!isNaN(val)) {
            const thresholds = this.levels[key];
            const isReverse = reverseFields.has(key);
            const index = this.getLevelIndex(val, thresholds, isReverse);
            worstLevel = Math.max(worstLevel, index);
          }
        }
        return {
          水质等级: this.levelLabels[worstLevel] || "错误数据",
        };
      });
      return classified;
    },
    conveyList(list) {
      const classified = this.classify(list);
      const levelCount = {};

      classified.forEach((item) => {
        const level = item.水质等级;
        levelCount[level] = (levelCount[level] || 0) + 1;
      });

      const pieData = this.levelLabels.map(label => ({
        name: label,
        value: levelCount[label] || 0
      }));
      bus.$emit("getDataByField", pieData);
    },
    conveyOptions(){
      bus.$emit("curOptions", this.options);
    }
  },
  mounted() {
    bus.$on("requestByfield", (fields) => {
      const filtered = this.pieLineData.map((item) => {
        const obj = {};
        fields.forEach((field) => {
          if (field.value in item) {
            obj[field.key] = item[field.value];
          }
        });
        return obj;
      });
      this.conveyList(filtered);
    });
    bus.$on("ready", () => {
      this.conveyOptions();
    });
  },
  async created() {
    await this.initStationOptions();
    this.getLineList();
    bus.$on("logOut", () => {
      this.controller.abort();
    });
    window.addEventListener("resize", () => {
      this.containerHeight = window.innerHeight - 57;
    });
  },
};
</script>

<style lang="less" scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  justify-content: center;
  position: fixed;
  left: 0;
  top: 20%;
  z-index: 999;
  box-shadow: h-shadow v-shadow blur spread color inset;
}

.item {
  margin: 4px;
}

.system {
  .el-select {
    margin: 5px 0;
  }

  .el-date-editor {
    margin: 5px 0;
  }
}

:deep(.el-popover) {
  height: 500px;
  overflow: auto;
}

.archor-container {
  max-height: 300px;

  overflow: auto;
}
</style>
