<template>
  <div class="ov-container" :style="`height: ${containerHeight}px`"> 
    <div class="row" style="margin: 0.5% 0">
      <Map pattern="overview"></Map>
    </div>
    <span class="title">可视化图表</span>  
    <div class="option-container">
    <el-form
        ref="dataForm"
        :model="formData"
      >
        <el-form-item prop="system">
          <span class="item_label">制度：</span>
          <el-select
            v-model="formData.system"
            placeholder="请选择制度"
            @change="systemChange"
          >
          <el-option label="小时制" value="小时制" />
          <el-option label="月度制" value="月度制" />
          </el-select>
        </el-form-item>
        <el-form-item prop="stationName">
          <span class="item_label">站名：</span>
          <el-select
            v-model="formData.stationName"
            placeholder="请选择站名"
          >
            <el-option
              v-for="item in stationOptions"
              :key="item.name"
              :label="item.name"
              :value="item.name"
            />
          </el-select>
        </el-form-item>
        <el-form-item prop="startTime">
          <span class="item_label">开始日期：</span>
          <el-date-picker
            v-model="formData.startTime"
            type="datetime"
            placeholder="选择日期时间"
            clearable
          >
          </el-date-picker>
        </el-form-item>
        <el-form-item prop="endTime">
          <span class="item_label">结束日期：</span>
          <el-date-picker
            v-model="formData.endTime"
            type="datetime"
            placeholder="选择日期时间"
            clearable
          >
          </el-date-picker>
        </el-form-item>
        <el-button type="primary" @click="updateData()"> 确认 </el-button>
      </el-form>
  </div> 
    <div class="row">
      <div id="lineGraph" class="col card" style="width: 96.66%">
        <LineGraph pattern="overview"></LineGraph>
      </div>
    </div>
    <div class="row">
      <div class="col card" style="width: 55%;">
        <PieGraph pattern="overview"></PieGraph>
      </div>
      <div class="col card" style="width: 40%; display: flex;">
        <Graph></Graph>
      </div>
    </div>
    <div style="height: 1px"></div>
  </div>
</template>

<script>
import {
  getLineData,
  getStationName,
} from "@/api/graphs.js";
import { fullFormatTime } from "@/util/timeFormater";
import bus from "@/util/eventBus";
import LineGraph from "@/components/graphs/liness.vue";
import PieGraph from "@/components/graphs/pie.vue";
import Graph from "@/components/fortcase/graph.vue";
import Map from "@/components/map/map.vue";
export default {
  name: "my-overview",
  data() {
    return {
      containerHeight: window.innerHeight - 57,
      stationNames: [],
      stationOptions: [],
      stationName: "",
      formData: {
        system: "小时制",
        stationName: "",
        startTime: "",
        endTime: "",
      },
    };
  },
  methods:{
    fullFormatTime,
    systemChange() {
      this.stationOptions = this.stationNames.filter(item => item.system === this.formData.system);
      this.formData.stationName = "";
      this.stationName = "";
    },
    async initStationOptions() {
      try {
        const res = await getStationName();
        if (res.code === 200) {
          this.stationNames = res.data.stationNames;
          
          if (this.stationNames.length > 0) {
            this.formData.stationName = this.stationNames[0].name;
            this.stationName = this.formData.stationName;
            this.stationOptions = this.stationNames.filter(item => item.system === this.formData.system);
          }
        }
      } catch (err) {
        console.error("加载失败", err);
      }
    },
    getLineList() {
      let params = "";
      if (this.formData.startTime) {
        params += `start=${this.fullFormatTime(this.formData.startTime)}&`;
      }
      if (this.formData.endTime) {
        params += `end=${this.fullFormatTime(this.formData.endTime)}`;
      }

      bus.$emit("showLoading");

      getLineData(this.stationName, this.formData.system, params)
        .then((res) => {
          if (res.code === 200) {
            this.formData.startTime = res.data.startTime;
            this.formData.endTime = res.data.endTime;
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

              // 计算 data 中 time 字段的索引
              const indexOfTime = Object.keys(data[0]).indexOf("time");
              bus.$emit("drawLine", this.stationName, indexOfTime, data);
              bus.$emit("dataChange", this.formData.startTime, this.formData.endTime,  this.formData.system, this.stationName);
            }
          }
        })
        .catch((err) => {
          this.$message.error(err.message);
          console.log(err);
        });
    },
    async updateData() {
      if (!this.formData.system) {
        this.$message.warning("制度不能为空");
        return;
      }
      if (!this.formData.stationName) {
        this.$message.warning("站名不能为空");
        return;
      }
      if (!this.formData.startTime) {
        this.$message.warning("开始时间不能为空");
        return;
      }
      if (!this.formData.endTime) {
        this.$message.warning("结束时间不能为空");
        return;
      }
      if (new Date(this.formData.startTime) > new Date(this.formData.endTime)) {
        this.$message.warning("结束时间不能小于开始时间");
        return;
      }
      this.stationName = this.formData.stationName;
      this.getLineList();
    },
  },
  async mounted(){
    await this.initStationOptions();
    this.getLineList();
  },
  components: {
    LineGraph,
    PieGraph,
    Map,
    Graph,
  },
};
</script>

<style lang="less" scoped>
.ov-container {
  width: 100%;
  padding-top: 56px;
  .title{
    display: block;
    text-align: center;
    box-sizing: border-box;
    padding-top: 30px;
    padding-bottom: 30px;
    font-size: 36px;
    font-weight: bold;
    color: #606266;
  }

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

  .option-container {
    min-width: 1261px;
    max-width: 96.66%;
    margin: 0px;
    margin-left: calc(3.33% / 2);
    padding: 20px 10px;
    box-sizing: border-box;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    background-color: rgb(248, 249, 251);
    .el-form{
      display: flex;
      height: 40px;

      .el-form-item {
        .item_label {
          color: #606266;
        }

        margin-left: 10px;
      }

      .el-button {
        margin-left: 10px;
      }
    }

    :deep(label::before) {
      display: none;
    }
  }
}
</style>
