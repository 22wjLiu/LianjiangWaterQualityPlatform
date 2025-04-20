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
        placeholder="请选择映射状态"
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
            @click="handleEdit(scope.row.id, scope.row.version_name)"
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
      <div slot="footer">
        <el-button @click="createFormVisible = false"> 取消 </el-button>
        <el-button type="primary" @click="createData()"> 确认 </el-button>
      </div>
    </el-dialog>

    <el-dialog
      :title="`映射版本——${editVersionName}`"
      class="edit-dialog"
      :visible.sync="editFormVisible"
      center
    >
      <div class="edit-table-container">
        <div class="edit-searcher-container">
          <el-select
            placeholder="请选择映射表类型"
            v-model="infosSearchList[0].value"
            clearable
          >
            <el-option
              v-for="item in tableOptions"
              :key="item.value"
              :label="item.value"
              :value="item.value"
            >
            </el-option>
          </el-select>
          <el-input
            placeholder="请输入主键"
            v-model="infosSearchList[1].value"
            clearable
          >
          </el-input>
          <el-input
            placeholder="请输入值"
            v-model="infosSearchList[2].value"
            clearable
          >
          </el-input>
          <el-button
            type="primary"
            icon="el-icon-search"
            @click="handleInfosSearch"
          >
            搜索
          </el-button>
          <el-button type="success" @click="handleCreateMap()">
            新建映射
          </el-button>
          <el-button type="danger" @click="handleEditDelete()">
            批量删除
          </el-button>
        </div>
        <el-table
          :data="editTableData"
          v-loading="editLoading"
          style="width: 100%"
          border
          @selection-change="handleEditSelectionChange"
        >
          <el-table-column type="selection" align="center" width="55">
          </el-table-column>
          <el-table-column prop="table" label="映射类型" align="center">
          </el-table-column>
          <el-table-column prop="key" label="主键" align="center">
          </el-table-column>
          <el-table-column prop="value" label="值" align="center">
          </el-table-column>
          <el-table-column label="操作" align="center" width="180">
            <template slot-scope="scope">
              <el-button
                type="primary"
                size="small"
                @click="handleEditUpdate(scope.row)"
                >编辑</el-button
              >
              <el-button
                type="danger"
                size="small"
                @click="handleEditDelete(scope.row.id)"
                >删除</el-button
              >
            </template>
          </el-table-column>
        </el-table>
        <div class="edit-page-container">
          <el-pagination
            @size-change="handleEditSizeChange"
            @current-change="handleEditCurrentChange"
            :current-page="infosSearchList[3].value"
            :page-sizes="[10, 20, 30, 40]"
            :page-size="infosSearchList[4].value"
            layout="total, sizes, prev, pager, next, jumper"
            :total="editTotalNum"
          >
          </el-pagination>
        </div>
      </div>
    </el-dialog>

    <el-dialog title="创建映射表" :visible.sync="mapFormVisible" center>
      <el-form
        ref="createMapForm"
        :rules="curEditRule"
        :model="editTemp"
        label-position="left"
        label-width="85px"
      >
        <el-form-item label="映射类型" prop="table">
          <el-select
            placeholder="是否复制当前使用中版本"
            v-model="editTemp.table"
            clearable
            @change="handleTableChange"
          >
            <el-option
              v-for="item in createMapOptions"
              :key="item.value"
              :label="item.value"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="主键" prop="key">
          <el-input
            v-model="editTemp.key"
            autocomplete="off"
            placeholder="请设置主键"
          >
          </el-input>
        </el-form-item>
        <el-form-item label="值" prop="value">
          <el-input
            v-model="editTemp.value"
            autocomplete="off"
            placeholder="请设置值"
          >
          </el-input>
        </el-form-item>
        <el-form-item v-if="isMutiLineMap" label="公式" prop="value">
          <el-input
            v-model="editTemp.formula"
            autocomplete="off"
            placeholder="请设置公式"
          >
          </el-input>
        </el-form-item>
      </el-form>
      <div slot="footer">
        <el-button @click="mapFormVisible = false"> 取消 </el-button>
        <el-button type="primary" @click="createMap()"> 确认 </el-button>
      </div>
    </el-dialog>

    <el-dialog title="更新映射表" :visible.sync="updateMapFormVisible" center>
      <el-form
        ref="updateMapForm"
        :rules="curUpdateRule"
        :model="updateTemp"
        label-position="left"
        label-width="85px"
      >
        <el-form-item v-if="isNotFormula" label="主键" prop="key">
          <el-input
            v-model="updateTemp.key"
            autocomplete="off"
            placeholder="请设置主键"
          >
          </el-input>
        </el-form-item>
        <el-form-item label="值" prop="value">
          <el-input
            v-model="updateTemp.value"
            autocomplete="off"
            placeholder="请设置值"
          >
          </el-input>
        </el-form-item>
      </el-form>
      <div slot="footer">
        <el-button @click="updateMapFormVisible = false"> 取消 </el-button>
        <el-button type="primary" @click="updateData()"> 确认 </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { formatTime, dateFullFormatTime } from "@/util/timeFormater.js";
import {
  getMapVersions,
  getMapTables,
  getMapInfos,
  createMap,
  deleteMap,
  updateMap,
  createMapVersion,
  deleteMapVersion,
  changeMapVersion,
} from "@/api/map.js";
export default {
  data() {
    return {
      totalNum: 0,
      editTotalNum: 0,
      curParams: "",
      curEditParams: "",
      loading: true,
      editLoading: true,
      createFormVisible: false,
      editFormVisible: false,
      mapFormVisible: false,
      updateMapFormVisible: false,
      isMutiLineMap: false,
      isNotFormula: true,
      curEditRule: {},
      curUpdateRule: {},
      editVersionName: "无",
      editId: "",
      createdAt: [],
      tableData: [],
      editTableData: [],
      selection: [],
      editSelection: [],
      tableOptions: [],
      createMapOptions: [],
      createTemp: {
        version_name: "",
        isCopy: "",
      },
      editTemp: {
        table: "",
        key: "",
        value: "",
        formula: "",
      },
      updateTemp: {
        key: "",
        value: "",
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
      infosSearchList: [
        {
          label: "table",
          value: "",
        },
        {
          label: "key",
          value: "",
        },
        {
          label: "value",
          value: "",
        },
        {
          label: "page",
          value: 1,
        },
        {
          label: "pageSize",
          value: 10,
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
      updateRules1: {
        key: [
          { required: true, message: "主键不能为空", trigger: "blur" },
        ],
        value: [
          { required: true, message: "值不能为空", trigger: "blur" },
        ],
      },
      updateRules2: {
        value: [
          { required: true, message: "值不能为空", trigger: "blur" },
        ],
      },
      editRule1: {
        table: [
          { required: true, message: "该选项不能为空", trigger: "blur" },
          { required: true, message: "该选项不能为空", trigger: "change" },
        ],
        key: [{ required: true, message: "主键不能为空", trigger: "blur" }],
        value: [{ required: true, message: "值不能为空", trigger: "blur" }],
      },
      editRule2: {
        table: [
          { required: true, message: "该选项不能为空", trigger: "blur" },
          { required: true, message: "该选项不能为空", trigger: "change" },
        ],
        key: [{ required: true, message: "主键不能为空", trigger: "blur" }],
        value: [{ required: true, message: "值不能为空", trigger: "blur" }],
        formula: [{ required: true, message: "公式不能为空", trigger: "blur" }],
      },
    };
  },
  methods: {
    formatTime,
    dateFullFormatTime,
    getMapTables() {
      getMapTables()
        .then((res) => {
          this.tableOptions = res.data.tables;
          this.createMapOptions = this.tableOptions.filter(item => item.value !== "行字段一对多公式映射");
        })
        .catch((err) => {
          this.$message.error(err.message);
        });
    },
    getTableData(params) {
      if (params !== this.curParams) {
        this.searchList[2].value = 1;
        this.curParams = params;
      }
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
          console.log(err);
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
    handleInfosSearch() {
      let params = "";
      this.infosSearchList.forEach((item) => {
        if (item.value !== "") {
          params += item.label + "=" + item.value + "&";
        }
      });
      const last = params.lastIndexOf("&");
      params = params.slice(0, last);
      this.getEditTableData(params);
    },
    getEditTableData(params) {
      if (params !== this.curEditParams) {
        this.infosSearchList[3].value = 1;
        this.curEditParams = params;
      }
      const query = params ? `?${params}` : "";
      this.editLoading = true;
      getMapInfos(this.editId, query)
        .then((res) => {
          if (res.code == 200) {
            this.editTableData = res.data.mapInfos;
            this.editTotalNum = res.data.total;
            this.editLoading = false;
          }
        })
        .catch((err) => {
          this.$message.error(err.message);
          console.log(err);
        });
    },
    handleDelete(id) {
      this.$confirm(
        "此操作将永久删除选中映射版本及其相关数据表, 是否继续?",
        "提示",
        {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning",
        }
      ).then(() => {
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
    handleEditDelete(id) {
      this.$confirm(
        "此操作将永久删除选中映射并删除数据表对应映射, 是否继续?",
        "提示",
        {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning",
        }
      ).then(() => {
        this.editDeleteData(id);
      });
    },
    editDeleteData(id) {
      let ids = [];
      if (id) {
        ids.push(parseInt(id));
      } else {
        this.selection.forEach((item) => {
          ids.push(parseInt(item.id));
        });
      }
      ids.forEach((item) => {
        if (item.table === "行字段一对多映射") {
          const mathed = this.editTableData.find(
            (data) =>
              data.table === "行字段一对多公式映射" && data.key === item.key
          );
          if (mathed) {
            ids.push(parseInt(mathed.id));
          }
        } else if (item.table === "行字段一对多公式映射") {
          const mathed = this.editTableData.find(
            (data) => data.table === "行字段一对多映射" && data.key === item.key
          );
          if (mathed) {
            ids.push(parseInt(mathed.id));
          }
        }
      });
      ids = [...new Set(ids)];
      deleteMap(this.editId, ids)
        .then((res) => {
          if (res.code == 200) {
            this.handleInfosSearch();
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
    handleEditSelectionChange(selected) {
      this.editSelection = selected;
    },
    handleSizeChange(val) {
      this.searchList[3].value = val;
      this.handleSearch();
    },
    handleEditSizeChange(val) {
      this.infosSearchList[4].value = val;
      this.handleInfosSearch();
    },
    handleCurrentChange(val) {
      this.searchList[2].value = val;
      this.handleSearch();
    },
    handleEditCurrentChange(val) {
      this.infosSearchList[3].value = val;
      this.handleInfosSearch();
    },
    handleEdit(id, name) {
      this.editVersionName = name;
      this.editFormVisible = true;
      this.editId = id;
      this.infosSearchList = [
        {
          label: "table",
          value: "",
        },
        {
          label: "key",
          value: "",
        },
        {
          label: "value",
          value: "",
        },
        {
          label: "page",
          value: 1,
        },
        {
          label: "pageSize",
          value: 10,
        },
      ],
      this.handleInfosSearch();
    },
    handleEditUpdate(row) {
      if (row.table === "行字段一对多公式映射") {
        this.isNotFormula = false;
        this.curUpdateRule = this.updateRules2;
      } else {
        this.isNotFormula = true;
        this.curUpdateRule = this.updateRules1;
      }
      this.updateTemp = {
        table: row.table,
        key: row.key,
        value: row.value,
      };
      this.origin.id = row.id;
      this.updateMapFormVisible = true;
      this.$nextTick(() => {
        this.$refs.updateMapForm.clearValidate();
      });
    },
    updateData() {
      this.$refs.updateMapForm.validate((valid) => {
        if (!valid) return;
        updateMap(this.editId, this.origin.id, this.updateTemp)
        .then((res) => {
          if (res.code === 200) {
            this.$message.success(res.msg);
            this.updateMapFormVisible = false;
            this.handleInfosSearch();
          } else {
            this.$message.warning(res.msg);
          }
        })
        .catch((err) => {
          this.$message.error(err.message);
        });
      });
    },
    handleTableChange(val) {
      if (val === "行字段一对多映射") {
        this.isMutiLineMap = true;
        this.curEditRule = this.editRule2;
      } else {
        this.isMutiLineMap = false;
        this.curEditRule = this.editRule1;
      }
    },
    handleCreateMap() {
      this.mapFormVisible = true;
      this.editTemp = {};
      this.$nextTick(() => {
        this.$refs.createMapForm.clearValidate();
      });
    },
    createMap() {
      this.$refs.createMapForm.validate((valid) => {
        if (!valid) return;
        this.$confirm("确认创建吗?", "提示", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning",
        }).then(() => {
          createMap(this.editId, this.editTemp)
            .then((res) => {
              if (res.code === 200) {
                this.$message.success(res.msg);
                this.mapFormVisible = false;
                this.handleInfosSearch();
              }
            })
            .catch((err) => {
              this.$message.error(err.message);
              console.log(err.message);
            });
        });
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
      this.$confirm("确认要切换当前映射版本吗?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      }).then(() => {
        const params = id ? `?id=${id}` : "";
        changeMapVersion(params)
          .then((res) => {
            if (res.code === 200) {
              this.handleSearch();
              this.$message.success(res.msg);
            }
          })
          .catch((err) => {
            this.$message.error(err.message);
            console.log(err.message);
          });
      });
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

.body .edit-dialog :deep(.el-dialog) {
  max-width: 1300px;
  min-width: 890px;

  .edit-table-container {
    width: 100%;
  }

  .el-table {
    width: 100% !important;
  }

  .edit-searcher-container {
    .el-input {
      margin-bottom: 5px;
      margin-left: 10px;
      width: 150px;
    }

    :nth-child(1) {
      margin-left: 0px;
    }

    .el-button {
      margin-bottom: 5px;
      margin-left: 10px;
    }

    .el-select {
      margin-bottom: 5px;
      width: 200px;

      .el-input {
        width: 100%;
      }
    }
  }

  .edit-page-container {
    display: flex;
    margin-top: 10px;
    .el-pagination {
      display: inline-block;
    }
  }
}
</style>
