<template>
  <div class="body">
    <!-- 搜索输入区 -->
    <div class="searcher-container">
      <el-select
        placeholder="请选择映射类型"
        v-model="searchList[2].value"
        clearable
      >
        <el-option
          v-for="item in tableOptions"
          :key="item.value"
          :label="item.label"
          :value="item.value"
        >
        </el-option>
      </el-select>
      <el-input
        placeholder="请输入键名"
        v-model="searchList[0].value"
        clearable
      >
      </el-input>
      <el-input placeholder="请输入值" v-model="searchList[1].value" clearable>
      </el-input>
      <el-button type="primary" icon="el-icon-search" @click="handleSearch">
        搜索
      </el-button>
      <el-button type="primary" icon="el-icon-plus" @click="handleCreate">
        新增版本
      </el-button>
      <el-button type="primary" icon="el-icon-s-order" @click="handleVersion">
        版本管理
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
      <el-table-column prop="table" label="类型" align="center">
      </el-table-column>
      <el-table-column prop="key" label="键" align="center"> </el-table-column>
      <el-table-column prop="value" label="值" align="center">
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
        :current-page="searchList[3].value"
        :page-sizes="[25, 50, 75, 100]"
        :page-size="searchList[4].value"
        layout="total, sizes, prev, pager, next, jumper"
        :total="totalNum"
      >
      </el-pagination>
      <el-button type="danger" size="small" @click="handleMutiDelete()"
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
        <el-form-item label="等级" prop="level">
          <el-select v-model="temp.level" placeholder="请设置等级">
            <el-option
              v-for="item in tableOptions"
              :key="item.value"
              :label="item.value"
              :value="item.value"
            />
          </el-select>
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
import { formatTime } from "@/util/timeFormater.js";
import { getCurMaps, getMapTables } from "@/api/map.js";
export default {
  data() {
    return {
      totalNum: 0,
      loading: true,
      dialogFormVisible: false,
      tableData: [],
      selection: [],
      tableOptions: [],
      temp: {},
      origin: {},
      searchList: [
        {
          label: "key",
          value: "",
        },
        {
          label: "value",
          value: "",
        },
        {
          label: "table",
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
      rules: {
        name: [{ required: true, message: "用户名不能为空", trigger: "blur" }],
        level: [{ required: true, message: "等级不能为空", trigger: "change" }],
      },
    };
  },
  methods: {
    formatTime,
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
      const query = params ? `?${params}` : "";
      this.loading = true;
      getCurMaps(query)
        .then((res) => {
          this.tableData = res.data.maps;
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
    handleCreate() {
      console.log("yes")
    },
    handleVersion() {
      console.log("yes")
    }
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