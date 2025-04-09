<template>
  <el-card>
    <div class="title">
      <img src="../../assets/cloud-download.png" />
      <h3>{{ sort }}文件下载</h3>
    </div>
    <el-divider></el-divider>
    <div class="body">
      <el-select v-model="value" placeholder="请选择">
        <el-option
          v-for="item in list"
          :key="item.id"
          :label="item.file_name"
          :value="item.id"
        ></el-option>
      </el-select>
      <el-button type="primary" @click="download">点击下载</el-button>
    </div>
  </el-card>
</template>

<script>
import { download } from "@/api/file";
export default {
  props: ["sort", "list"],
  data() {
    return {
      value: "",
    };
  },
  methods: {
    download() {
      let params = "";
      if (this.value) {
        params = `?id=${this.value}`;
      }
      download(params)
        .then((res) => {
          // 创建 Blob 对象
          const blob = new Blob([res.data], {
            type: "application/vnd.ms-excel",
          });

          // 提取文件名
          let fileName = "下载文件.xlsx";
          const matchedItem = this.list.find(item => item.id === this.value);
          if (matchedItem) {
            fileName = matchedItem.file_name;
          }

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
  created() {},
};
</script>

<style lang="less" scoped>
.el-card {
  width: 500px;
  background-color: #f2f6fc;
  margin: 18px 10px;
  border: 1px solid #dcdfe6;

  .title {
    display: flex;
    align-items: center;
    padding-bottom: 0;
    height: 24px;
    img {
      width: 24px;
    }
    h3 {
      line-height: 24px;
      padding-left: 10px;
    }
  }
  .el-select {
    width: 70%;
  }
  .el-button {
    float: right;
  }
}
</style>
