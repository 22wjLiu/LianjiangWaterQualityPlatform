<template>
  <div class="body">
    <!-- 搜索输入区 -->
    <div class="searcher-container">
      <el-input
        placeholder="请输入用户ID"
        v-model="searchList[0].value"
        clearable
      >
      </el-input>
      <el-input
        placeholder="请输入站点"
        v-model="searchList[1].value"
        clearable
      >
      </el-input>
      <el-select
        v-model="searchList[2].value"
        placeholder="请选择制度"
        clearable
      >
        <el-option
          v-for="item in sysOptions"
          :key="item.value"
          :label="item.label"
          :value="item.value"
        >
        </el-option>
      </el-select>
      <el-select
        v-model="searchList[3].value"
        placeholder="请选择操作方式"
        clearable
      >
        <el-option
          v-for="item in options"
          :key="item.label"
          :label="item.label"
          :value="item.label"
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
      <el-button type="primary" icon="el-icon-search" @click="handleSearch"
        >搜索</el-button
      >
    </div>

    <!-- 表格 -->
    <el-table
      :data="tableData"
      v-loading="loading"
      style="width: 94.6%; min-width: 1230px; left: max(2.7%, 35px)"
      border
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" align="center" width="55">
      </el-table-column>
      <el-table-column label="创建时间" align="center">
        <template slot-scope="scope">
          {{ formatTime(scope.row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column prop="user_id" label="用户ID" align="center" width="100">
      </el-table-column>
      <el-table-column prop="station_name" label="站点" align="center">
      </el-table-column>
      <el-table-column prop="system" label="时间制" align="center" width="100">
      </el-table-column>
      <el-table-column label="开始时间" align="center">
        <template slot-scope="scope">
          {{ formatTime(scope.row.start_time) }}
        </template>
      </el-table-column>
      <el-table-column label="结束时间" align="center">
        <template slot-scope="scope">
          {{ formatTime(scope.row.end_time) }}
        </template>
      </el-table-column>
      <el-table-column label="操作方式" align="center" width="100">
        <template slot-scope="scope">
          <el-tag size="medium" :type="getTagType(scope.row.option)">
            {{ scope.row.option }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="100">
        <template slot-scope="scope">
          <el-button
            type="danger"
            size="small"
            @click="handleDelete(scope.row.id)"
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
        :current-page="searchList[4].value"
        :page-sizes="[25, 50, 75, 100]"
        :page-size="searchList[5].value"
        layout="total, sizes, prev, pager, next, jumper"
        :total="totalNum"
      >
      </el-pagination>
      <el-button type="danger" size="small" @click="handleDelete()"
        >批量删除</el-button
      >
    </div>
  </div>
</template>

<script>
import { formatTime, dateFullFormatTime } from "@/util/timeFormater.js";
import { getDataLog, deleteDataLog } from "@/api/history.js";
export default {
  data() {
    return {
      totalNum: 0,
      curParams: "",
      loading: true,
      tableData: [],
      createdAt: [],
      selection: [],
      searchList: [
        {
          label: "id",
          value: "",
        },
        {
          label: "station_name",
          value: "",
        },
        {
          label: "system",
          value: "",
        },
        {
          label: "option",
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
      sysOptions: [
        {
          value: "小时制",
        },
        {
          value: "月度制",
        },
      ],
      options: [
        {
          label: "创建",
          type: "primary",
        },
        {
          label: "删除",
          type: "danger",
        },
        {
          label: "更新(前)",
          type: "warning",
        },
        {
          label: "更新(后)",
          type: "success",
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
    };
  },
  methods: {
    formatTime,
    dateFullFormatTime,
    getTableData(params) {
      if (params !== this.curParams) {
        this.searchList[4].value = 1;
        this.curParams = params;
      }
      const start = this.createdAt ? this.createdAt[0] : "";
      const end = this.createdAt ? this.createdAt[1] : "";
      const query = params ? `?${params}` : "";
      this.loading = true;
      getDataLog(this.dateFullFormatTime(start), this.dateFullFormatTime(end), query)
        .then((res) => {
          this.tableData = res.data.dataHistories;
          this.totalNum = res.data.total;
          this.loading = false;
        })
        .catch((err) => {
          this.$message.error(err.message);
        });
    },
    getTagType(option) {
      let type = "";
      this.options.some((item) => {
        if (item.label === option) {
          type = item.type;
          return true;
        }
        return false;
      });
      return type;
    },
    handleSearch() {
      let params = "";
      this.searchList.forEach((item) => {
        if (item.value !== "") {
          params += item.label + "=" + item.value + "&";
        }
      });
      const last = params.lastIndexOf("&");
      params = params.slice(0, last);
      this.getTableData(params);
    },
    handleSelectionChange(selected) {
      this.selection = selected;
    },
    handleDelete(id) {
      this.$confirm("此操作将永久删除选中日志, 是否继续?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      }).then(() => {
        this.deleteData(id);
      });
    },
    deleteData(id) {
      const ids = [];
      if (id) {
        ids.push(parseInt(id));
      } else {
        this.selection.forEach((item) => {
          ids.push(parseInt(item.id));
        });
      }
      deleteDataLog(ids)
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
      this.searchList[5].value = val;
      this.handleSearch();
    },
    handleCurrentChange(val) {
      this.searchList[4].value = val;
      this.handleSearch();
    },
  },
  created() {
    this.handleSearch();
  },
};
</script>

<style lang="less" scoped>
.body {
  width: 100%;
  min-width: 1300px;
  padding-bottom: 20px;
}

.searcher-container {
  padding-left: max(2.7%, 35px);
  margin-bottom: 15px;
  .el-input {
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

  .el-select {
    margin: 5px 0;
    margin-left: 10px;
    width: 150px;
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
</style>