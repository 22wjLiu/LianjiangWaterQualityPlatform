<template>
  <div class="body">
    <span style="position: relative; color:red; left: max(2.7%, 35px);">
        需要先选择搜索项（站名、制度、映射版本）后点击"搜索"才会显示数据，且只会在所选时间范围内排序并查询10000条数据
    </span>
    <!-- 搜索输入区 -->
    <div class="searcher-container">
      <el-select
        placeholder="请选择站名"
        v-model="searchList[0].value"
        @change="handleStationChange"
        clearable
      >
        <el-option
          v-for="item in stationNameOptions"
          :key="item.value"
          :label="item.value"
          :value="item.value"
        >
        </el-option>
      </el-select>
      <el-select
        placeholder="请选择制度"
        v-model="searchList[1].value"
        clearable
      >
        <el-option
          v-for="item in systemOptions"
          :key="item.value"
          :label="item.value"
          :value="item.value"
        >
        </el-option>
      </el-select>
      <el-select
        placeholder="请选择映射版本"
        v-model="searchList[2].value"
        clearable
      >
        <el-option
          v-for="item in mapVersionOptions"
          :key="item.value"
          :label="item.label"
          :value="item.value"
        >
        </el-option>
      </el-select>
      <div class="time-picker">
        <el-date-picker
          v-model="createdAt"
          type="daterange"
          align="right"
          unlink-panels
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          :picker-options="pickerOptions"
        >
        </el-date-picker>
      </div>
      <el-button type="primary" icon="el-icon-search" @click="handleSearch">
        搜索
      </el-button>
    </div>

    <!-- 表格 -->
    <el-table
      :data="tableData"
      v-loading="loading"
      style="width: 94.6%; min-width: 1230px; left: max(2.7%, 35px);"
      border
      @selection-change="handleSelectionChange"
    >
    <el-table-column v-if="tableData.length" type="selection" align="center" width="55">
    </el-table-column>
      <el-table-column v-if="tableData.length" fixed="left" label="时间" align="center" width="150">
        <template slot-scope="scope">
          {{ formatTime(scope.row["时间"]) }}
        </template>
      </el-table-column>

      <el-table-column
        v-for="key in tableKeys"
        :key="key"
        :prop="key"
        :label="key"
        align="center">
      </el-table-column>

      <el-table-column v-if="tableData.length" fixed="right" label="操作" align="center" width="150">
        <template slot-scope="scope">
          <el-button
            type="primary"
            size="small"
            @click="handleUpdate(scope.row.id, scope.row.name, scope.row.level)"
            >编辑</el-button
          >
          <el-button
            type="danger"
            size="small"
            @click="handleDelete(scope.row)"
            >删除</el-button
          >
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="page-container">
      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="searchList[3].value"
        :page-sizes="[25, 50, 75, 100]"
        :page-size="searchList[4].value"
        layout="total, sizes, prev, pager, next, jumper"
        :total="totalNum"
      >
      </el-pagination>
      <el-button type="danger" size="small" @click="handleDelete()"
        >批量删除</el-button
      >
    </div>

    <el-dialog title="编辑用户信息表" :visible.sync="dialogFormVisible" center>
      <el-form
        ref="dataForm"
        :rules="rules"
        :model="temp"
        label-position="left"
        label-width="70px"
      >
        <el-form-item label="用户名" prop="name">
          <el-input
            v-model="temp.name"
            autocomplete="off"
            placeholder="请设置用户名"
          >
          </el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false"> 取消 </el-button>
        <el-button type="primary" @click="updateData()"> 确认 </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { formatTime, fullFormatTime, dateFullFormatTime } from "@/util/timeFormater.js";
import {
  getDataTableInfos,
  getDataBackStage,
  deleteDataBackStage,
} from "@/api/data.js";
export default {
  data() {
    return {
      totalNum: 0,
      curParams: "",
      loading: false,
      dialogFormVisible: false,
      tableData: [],
      tableKeys: [],
      createdAt: [],
      selection: [],
      allOptions: [],
      stationNameOptions: [],
      systemOptions: [],
      mapVersionOptions: [],
      temp: {},
      origin: {},
      curSearchList: {
        stationName: "",
        system: "",
        versionName: "",
        dataTableName: "",
        times: [],
      },
      searchList: [
        {
          name: "站名",
          label: "station_name",
          value: "",
        },
        {
          name: "制度",
          label: "system",
          value: "",
        },
        {
          name: "映射版本",
          label: "map_ver_id",
          value: "",
        },
        {
          label: "page",
          value: 1,
        },
        {
          label: "pageSize",
          value: 25,
        },
      ],
      pickerOptions: {
        shortcuts: [
          {
            text: "最近7天",
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 7);
              picker.$emit("pick", [start, end]);
            },
          },
          {
            text: "最近30天",
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 30);
              picker.$emit("pick", [start, end]);
            },
          },
          {
            text: "最近90天",
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 90);
              picker.$emit("pick", [start, end]);
            },
          },
        ],
      },
      rules: {
        name: [{ required: true, message: "用户名不能为空", trigger: "blur" }],
        level: [{ required: true, message: "等级不能为空", trigger: "change" }],
      },
    };
  },
  methods: {
    formatTime,
    fullFormatTime,
    dateFullFormatTime,
    getStationSystemMapVersion() {
      getDataTableInfos()
        .then((res) => {
          if (res.code === 200) {
            this.allOptions = res.data.tableInfos;
            this.allOptions.forEach((item) => {
              let temp = { value: item.station_name }
              this.stationNameOptions.push(temp);
            });
          }
        })
        .catch((err) => {
          console.log(err.message);
        });
    },
    handleStationChange(val) {
      const matched = this.allOptions.filter(item => item.station_name === val);
      this.systemOptions = [];
      this.mapVersionOptions = [];
      matched.forEach((item) => {
        let temp = { value: item.system }
        this.systemOptions.push(temp);
        temp = {
          label: item.version_name,
          value: item.map_ver_id,
        };
        this.mapVersionOptions.push(temp);
      });
      this.systemOptions = [...new Set(this.systemOptions)];
      const seen = new Set();
      this.mapVersionOptions = this.mapVersionOptions.filter((opt) => {
        if (seen.has(opt.value)) return false;
        seen.add(opt.value);
        return true;
      });
    },
    getTableData(params) {
      if (params !== this.curParams) {
        this.searchList[3].value = 1;
        this.curParams = params;
      }
      const start = this.createdAt ? this.createdAt[0] : "";
      const end = this.createdAt ? this.createdAt[1] : "";
      const query = params ? `?${params}` : "";
      this.loading = true;
      getDataBackStage(
        this.dateFullFormatTime(start),
        this.dateFullFormatTime(end),
        query
      )
        .then((res) => {
          this.tableData = res.data.resultArr;
          this.totalNum = res.data.total;
          if (this.tableData.length) {
            this.tableKeys = Object.keys(this.tableData[0]).filter((item) => {
            return item !== "时间";
           });
           this.curSearchList.versionName = res.data.versionName;
           this.curSearchList.stationName = res.data.stationName;
           this.curSearchList.system = res.data.system;
           this.curSearchList.dataTableName = res.data.dataTableName;
          }
          this.loading = false;
        })
        .catch((err) => {
          this.$message.error(err.message);
          console.log(err.message);
        });
    },
    handleSearch() {
      let params = "";
      let errName = false;
      this.searchList.some(item => {
        if (item.value !== "") {
          params += item.label + "=" + item.value + "&";
        } else {
          errName = item.name;
          return true
        }
        return false
      });
      if (errName) {
        this.$message.error(`${errName}不能为空`)
        return
      }
      const last = params.lastIndexOf("&");
      params = params.slice(0, last);
      this.getTableData(params);
    },
    handleSelectionChange(selected) {
      this.selection = selected;
    },
    handleDelete(row) {
      this.$confirm("此操作将永久删除选中数据, 是否继续?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      }).then(() => {
        this.deleteData(row);
      });
    },
    deleteData(row) {
      const temp = [];
      if (row) {
        temp.push(this.fullFormatTime(row["时间"]));
      } else {
        this.selection.forEach((item) => {
          temp.push(this.fullFormatTime(item["时间"]));
        })
      }
      
      this.curSearchList.times = temp;
      deleteDataBackStage(this.curSearchList)
        .then((res) => {
          if (res.code == 200) {
            this.handleSearch();
            this.$message.success(res.msg);
          } else {
            this.$message.warning(res.msg);
          }
        })
        .catch((err) => {
          this.$message.error(err.message);
        });
    },
    handleSizeChange(val) {
      this.searchList[4].value = val;
      this.handleSearch();
    },
    handleCurrentChange(val) {
      this.searchList[3].value = val;
      this.handleSearch();
    },
    handleUpdate(id, name, level) {
      this.origin = {
        name: name,
        level: level,
      };
      this.temp = Object.assign({}, this.origin);
      this.origin.id = id;
      this.dialogFormVisible = true;
      this.$nextTick(() => {
        this.$refs.dataForm.clearValidate();
      });
    },
    updateData() {
      const body = Object.assign({}, this.temp);
      if (this.temp.name === this.origin.name) body.name = "";
      if (this.temp.level === this.origin.level) body.level = 0;
      updateUser(this.origin.id, body)
        .then((res) => {
          if (res.code === 200) {
            this.$message.success(res.msg);
            this.dialogFormVisible = false;
            this.handleSearch();
          } else {
            this.$message.warning(res.msg);
          }
        })
        .catch((err) => {
          this.$message.error(err.message);
        });
    },
  },
  created() {
    this.getStationSystemMapVersion();
  },
};
</script>

<style lang="less" scoped>
.body {
  padding: 20px 0px;
  width: 100%;
  min-width: 1300px;
}

.searcher-container {
  margin-bottom: 15px;
  padding-left: max(2.7%, 35px);
  .el-select {
    margin: 5px 0;
    margin-left: 10px;
    width: 150px;
  }

  :nth-child(1) {
    margin-left: 0px;
  }

  .time-picker {
    display: inline-block;
  }

  .time-picker,
  .el-button {
    margin: 5px 0;
    margin-left: 10px;
  }
}

.page-container {
  display: flex;
  margin-top: 10px;
  padding-left: max(2.7%, 35px);
  width: 94.6%;
  .el-pagination {
    display: inline-block;
  }

  .el-button {
    width: 100px;
    margin-left: auto;
  }
}

.body :deep(.el-dialog) {
  max-width: 500px;
  min-width: 300px;

  .el-select {
    width: 100%;
  }
}
</style>
