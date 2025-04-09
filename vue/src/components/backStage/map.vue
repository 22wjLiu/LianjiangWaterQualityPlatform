<template>
  <div class="body">
    <!-- 搜索输入区 -->
    <div class="searcher-container">
      <el-input
        placeholder="请输入版本名"
        v-model="searchList[0].value"
        clearable
      >
      </el-input>
      <el-select
        placeholder="请选择文件类型"
        v-model="searchList[1].value"
        clearable
      >
        <el-option
          v-for="item in statusOptions"
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
      <el-button type="primary" icon="el-icon-plus" @click="handleCreate">
        创建
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
      <el-table-column prop="created_at" label="创建时间" align="center">
        <template slot-scope="scope">
          {{ formatTime(scope.row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="更新时间" align="center">
        <template slot-scope="scope">
          {{ formatTime(scope.row.updated_at) }}
        </template>
      </el-table-column>
      <el-table-column prop="version_name" label="映射版本名" align="center">
      </el-table-column>
      <el-table-column label="状态" align="center" width="100">
        <template slot-scope="scope">
          <el-tag v-if="scope.row.active" type="success" size="medium">
            使用中
          </el-tag>
          <el-tag v-else type="danger" size="medium"> 未使用 </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="220">
        <template slot-scope="scope">
          <el-button
            type="success"
            size="small"
            @click="handleChange(scope.row.id)"
            >切换</el-button
          >
          <el-button
            type="primary"
            size="small"
            @click="handleUpdate(scope.row.id)"
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
      <el-button type="danger" size="small" @click="handleDelete()"
        >批量删除</el-button
      >
    </div>

    <el-dialog title="创建映射版本表" :visible.sync="createFormVisible" center>
      <el-form
        ref="createForm"
        :rules="createRules"
        :model="createTemp"
        label-position="left"
        label-width="85px"
      >
        <el-form-item label="版本名" prop="version_name">
          <el-input
            v-model="createTemp.version_name"
            autocomplete="off"
            placeholder="请设置映射版本名"
          >
          </el-input>
        </el-form-item>
        <el-form-item label="是否复制" prop="isCopy">
          <el-select
            placeholder="是否复制当前使用中版本"
            v-model="createTemp.isCopy"
            clearable
          >
            <el-option label="是" value="true"></el-option>
            <el-option label="否" value="false"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="createFormVisible = false"> 取消 </el-button>
        <el-button type="primary" @click="createData()"> 确认 </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { formatTime, dateFullFormatTime } from "@/util/timeFormater.js";
import { getMapVersions, getMapTables, createMapVersion, deleteMapVersion } from "@/api/map.js";
export default {
  data() {
    return {
      totalNum: 0,
      loading: true,
      createFormVisible: false,
      editFormVisible: false,
      createdAt: [],
      tableData: [],
      selection: [],
      tableOptions: [],
      createTemp: {
        version_name: "",
        isCopy: "",
      },
      origin: {},
      searchList: [
        {
          label: "version_name",
          value: "",
        },
        {
          label: "active",
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
      statusOptions: [
        {
          label: "使用中",
          value: 1,
        },
        {
          label: "未使用",
          value: 0,
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
      createRules: {
        version_name: [
          { required: true, message: "映射版本名不能为空", trigger: "blur" },
        ],
        isCopy: [
          { required: true, message: "该选项不能为空", trigger: "blur" },
          { required: true, message: "该选项不能为空", trigger: "change" },
        ],
      },
    };
  },
  methods: {
    formatTime,
    dateFullFormatTime,
    getMapTables() {
      getMapTables()
        .then((res) => {
          let temp;
          res.data.tables.forEach((item) => {
            temp = {
              value: item,
            };
            this.tableOptions.push(temp);
          });
        })
        .catch((err) => {
          this.$message.error(err.message);
        });
    },
    getTableData(params) {
      const start = this.createdAt ? this.createdAt[0] : "";
      const end = this.createdAt ? this.createdAt[1] : "";
      const query = params ? `?${params}` : "";
      this.loading = true;
      getMapVersions(
        this.dateFullFormatTime(start),
        this.dateFullFormatTime(end),
        query
      )
        .then((res) => {
          this.tableData = res.data.mapVersions;
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
      this.$confirm("此操作将永久删除选中映射版本及其相关数据表, 是否继续?", "提示", {
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
      deleteMapVersion(ids)
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
      this.searchList[2].value = val;
      this.handleSearch();
    },
    handleUpdate(id) {
      console.log(id);
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
    handleCreate() {
      this.createTemp = {};
      this.createFormVisible = true;
      this.$nextTick(() => {
        this.$refs.createForm.clearValidate();
      });
    },
    createData() {
      this.$refs.createForm.validate((valid) => {
        if (!valid) return;
        this.$confirm("确认创建吗?", "提示", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning",
        }).then(() => {
          const params = `?isCopy=${this.createTemp.isCopy}`;
          createMapVersion(params, {
            version_name: this.createTemp.version_name,
          })
            .then((res) => {
              if (res.code === 200) {
                this.$message.success(res.msg);
                this.createFormVisible = false;
                this.handleSearch();
              }
            })
            .catch((err) => {
              this.$message.error(err.message);
              console.log(err.message);
            });
        });
      });
    },
    handleChange(id) {
      console.log(id);
    },
  },
  created() {
    this.getMapTables();
    this.handleSearch();
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

.body :deep(.el-dialog) {
  max-width: 500px;
  min-width: 300px;

  .el-select {
    width: 100%;
  }
}
</style>
