<template>
  <div class="body">
    <!-- 搜索输入区 -->
    <div class="searcher-container">
      <el-select
        placeholder="请选择文件类型"
        v-model="searchList[0].value"
        clearable
      >
        <el-option
          v-for="item in typeOptions"
          :key="item.value"
          :label="item.value"
          :value="item.value"
        >
        </el-option>
      </el-select>
      <el-input
        placeholder="请输入文件名"
        v-model="searchList[1].value"
        clearable
      >
      </el-input>
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
      <el-button type="primary" icon="el-icon-plus" @click="handleUpload">
        上传文件
      </el-button>
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
      <el-table-column label="创建时间" align="center" width="150">
        <template slot-scope="scope">
          {{ formatTime(scope.row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="更新时间" align="center" width="150">
        <template slot-scope="scope">
          {{ formatTime(scope.row.updated_at) }}
        </template>
      </el-table-column>
      <el-table-column prop="file_name" label="文件名" align="center">
      </el-table-column>
      <el-table-column prop="file_path" label="文件路径" align="center">
      </el-table-column>
      <el-table-column prop="file_type" label="文件类型" align="center" width="55">
      </el-table-column>
      <el-table-column label="操作" align="center" width="220">
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
        :current-page="searchList[2].value"
        :page-sizes="[25, 50, 75, 100]"
        :page-size="searchList[3].value"
        layout="total, sizes, prev, pager, next, jumper"
        :total="totalNum"
      >
      </el-pagination>
      <el-button type="danger" size="small" @click="handleMutiDelete()"
        >批量删除</el-button
      >
    </div>

    <el-dialog :visible.sync="uploadVisible" center>
      <div class="uploadContainer">
        <UpLoad></UpLoad>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { formatTime, dateFullFormatTime } from "@/util/timeFormater.js";
import { getFileInfos } from "@/api/file.js";
import UpLoad from "@/components/load/UpLoad";
export default {
  data() {
    return {
      totalNum: 0,
      loading: true,
      uploadVisible: false,
      createdAt: [],
      tableData: [],
      selection: [],
      temp: {},
      origin: {},
      typeOptions: [
        {
          value: "xlsx",
        },
        {
          value: "xls",
        },
        {
          value: "css",
        },
      ],
      searchList: [
        {
          label: "fileType",
          value: "",
        },
        {
          label: "fileName",
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
    dateFullFormatTime,
    getTableData(params) {
      const start = this.createdAt ? this.createdAt[0] : "";
      const end = this.createdAt ? this.createdAt[1] : "";
      const query = params ? `?${params}` : "";
      this.loading = true;
      getFileInfos(
        this.dateFullFormatTime(start),
        this.dateFullFormatTime(end),
        query
      )
        .then((res) => {
          this.tableData = res.data.fileInfos;
          this.totalNum = res.data.total;
          this.loading = false;
        })
        .catch((err) => {
          this.$message.error(err.message);
        });
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
    handleDelete(id) {
      deleteUser(id)
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
    handleSelectionChange(selected) {
      this.selection = selected;
    },
    handleMutiDelete() {
      const ids = [];
      this.selection.forEach((item) => {
        ids.push(item.id);
      });
      deleteUsers(ids)
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
      this.searchList[3].value = val;
      this.handleSearch();
    },
    handleCurrentChange(val) {
      this.searchList[4].value = val;
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
    handleUpload(){
      this.uploadVisible = true
    }
  },
  created() {
    this.handleSearch();
  },
  components: {
    UpLoad,
  },
};
</script>

<style lang="less" scoped>
.body {
  margin-top: 20px;
  width: 100%;
  min-width: 1300px;
}

.searcher-container {
  margin-bottom: 15px;
  padding-left: max(2.7%, 35px);
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
    margin-left: 660px;
  }
}

.body :deep(.el-dialog) {
  max-width: 500px;
  min-width: 300px;

  .el-select {
    width: 100%;
  }
}

.body :deep(.el-dialog:first-of-type) {
  max-width: 800px;
  min-width: 550px;

  .el-dialog__body,
  .el-dialog__head {
    padding: 0px;
  }
}

.uploadContainer {
  :deep(.el-card),
  :deep(.el-upload),
  :deep(.el-upload-dragger) {
    width: 100% !important;
  }

  
}
</style>
