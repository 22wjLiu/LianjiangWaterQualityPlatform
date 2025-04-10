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
      <el-select
        placeholder="请选择制度"
        v-model="searchList[1].value"
        clearable
      >
        <el-option
          v-for="item in sysOptions"
          :key="item.value"
          :label="item.value"
          :value="item.value"
        >
        </el-option>
      </el-select>
      <el-input
        placeholder="请输入文件名"
        v-model="searchList[2].value"
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
      <el-table-column prop="system" label="制度" align="center" width="100">
      </el-table-column>
      <el-table-column label="文件名" align="center">
        <template slot-scope="scope">
          {{ scope.row.file_name.split(".")[0] }}
        </template>
      </el-table-column>
      <el-table-column prop="file_path" label="文件路径" align="center">
      </el-table-column>
      <el-table-column
        prop="file_type"
        label="文件类型"
        align="center"
        width="100"
      >
      </el-table-column>
      <el-table-column label="操作" align="center" width="220">
        <template slot-scope="scope">
          <el-button
            type="success"
            size="small"
            @click="handleDownload(scope.row.id, scope.row.file_name)"
            >下载</el-button
          >
          <el-button
            type="primary"
            size="small"
            @click="handleUpdate(scope.row.id, scope.row.file_name)"
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

    <el-dialog class="upload-dialog" :visible.sync="uploadVisible" center>
      <div class="upload-container">
        <UpLoad></UpLoad>
      </div>
    </el-dialog>

    <el-dialog title="编辑文件信息表" :visible.sync="editFormVisible" center>
      <el-form
        ref="editForm"
        :rules="rules"
        :model="temp"
        label-position="left"
        label-width="70px"
      >
        <el-form-item label="文件名" prop="fileName">
          <el-input
            v-model="temp.fileName"
            autocomplete="off"
            placeholder="请设置文件名"
          >
          </el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="editFormVisible = false"> 取消 </el-button>
        <el-button type="primary" @click="updateData()"> 确认 </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { formatTime, dateFullFormatTime } from "@/util/timeFormater.js";
import {
  getFileInfos,
  deleteFiles,
  updateFileName,
  download,
} from "@/api/file.js";
import UpLoad from "@/components/load/UpLoad";
export default {
  data() {
    return {
      totalNum: 0,
      curParams: "",
      loading: true,
      uploadVisible: false,
      editFormVisible: false,
      createdAt: [],
      tableData: [],
      selection: [],
      origin: {},
      temp: {
        id: "",
        fileName: "",
      },
      sysOptions: [
        {
          value: "小时制",
        },
        {
          value: "月度制",
        },
      ],
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
          label: "system",
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
        fileName: [
          { required: true, message: "文件名不能为空", trigger: "blur" },
        ],
      },
    };
  },
  methods: {
    formatTime,
    dateFullFormatTime,
    getTableData(params) {
      if (params !== this.curParams) {
        this.searchList[3].value = 1;
        this.curParams = params;
      }
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
      this.$confirm("此操作将永久删除选中文件, 是否继续?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      }).then(() => {
        this.deleteData(id);
      });
    },
    deleteData(id) {
      console.log("yes");
      const ids = [];
      if (id) {
        ids.push(parseInt(id));
      } else {
        this.selection.forEach((item) => {
          ids.push(parseInt(item.id));
        });
      }
      deleteFiles(ids)
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
    handleSizeChange(val) {
      this.searchList[3].value = val;
      this.handleSearch();
    },
    handleCurrentChange(val) {
      this.searchList[4].value = val;
      this.handleSearch();
    },
    handleUpdate(id, fileName) {
      this.origin = {
        fileName: fileName.split(".")[0],
      };
      this.temp = Object.assign({}, this.origin);
      this.origin.id = id;
      this.origin.orginName = fileName;
      this.editFormVisible = true;
      this.$nextTick(() => {
        this.$refs.editForm.clearValidate();
      });
    },
    updateData() {
      let body = {};
      if (this.temp.fileName === this.origin.fileName) {
        body.file_name = "";
      } else
        body.file_name =
          this.temp.fileName + "." + this.origin.orginName.split(".")[1];
      updateFileName(this.origin.id, body)
        .then((res) => {
          if (res.code === 200) {
            this.$message.success(res.msg);
            this.editFormVisible = false;
            this.handleSearch();
          } else {
            this.$message.warning(res.msg);
          }
        })
        .catch((err) => {
          this.$message.error(err.message);
        });
    },
    handleUpload() {
      this.uploadVisible = true;
    },
    handleDownload(id, fileName) {
      let params = "";
      if (id) {
        params = `?id=${id}`;
      }
      download(params)
        .then((res) => {
          // 创建 Blob 对象
          const blob = new Blob([res.data], {
            type: "application/vnd.ms-excel",
          });

          // 创建隐藏的下载链接并点击
          const href = URL.createObjectURL(blob);
          const a = document.createElement("a");
          a.href = href;
          a.download = fileName;
          a.style.display = "none";
          document.body.appendChild(a);
          a.click();
          document.body.removeChild(a);
          URL.revokeObjectURL(href);
        })
        .catch((err) => {
          this.$message.error(err.message);
          console.log(err);
        });
    },
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
  padding: 20px 0px;
  width: 100%;
  min-width: 1300px;
}

.searcher-container {
  margin-bottom: 15px;
  padding-left: max(2.7%, 35px);
  .el-input,
  .el-select {
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

.body .upload-dialog :deep(.el-dialog) {
  max-width: 800px;
  min-width: 550px;

  .el-dialog__body,
  .el-dialog__head,
  .el-dialog__header {
    padding: 0px;
  }
}

.upload-container {
  :deep(.el-card),
  :deep(.el-upload),
  :deep(.el-upload-dragger) {
    width: 100% !important;
  }
}
</style>
