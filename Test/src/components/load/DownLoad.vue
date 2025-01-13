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
          :key="item.name"
          :label="item.title"
          :value="item.name"
        ></el-option>
      </el-select>
      <el-button type="primary" @click="download">点击下载</el-button>
    </div>
  </el-card>
</template>

<script>
import { download } from '@/api/file'
export default {
  props: ['sort', 'list'],
  data() {
    return {
      value: ''
    }
  },
  methods: {
    download() {
      const system = this.sort === '小时制' ? 'hour' : 'month'
      const path = 'path=/' + system + '/' + this.value
      const filename = '&file=' + this.value
      const config = {
        headers: {
          Authorization: 'Bearer ' + localStorage.getItem('token')
        },
        responseType: 'arraybuffer'
      }
      download(path, filename, config)
        .then(res => {
          // Blob 可以将response的二进制流文件下载到本地
          const blob = new Blob([res.data], {
            type: 'application/vnd.ms-excel'
          })
          // 创建下载的链接
          const href = window.URL.createObjectURL(blob)
          const downloadElement = document.createElement('a')
          downloadElement.href = href
          downloadElement.style.display = 'none'
          // 下载后的文件名
          downloadElement.download = this.value
          document.body.appendChild(downloadElement)
          // 点击下载
          downloadElement.click()
          // 下载完成后移除元素
          document.body.removeChild(downloadElement)
          // 释放掉blob对象
          window.URL.revokeObjectURL(href)
        })
        .catch(err => {
          console.log(err)
          this.$message.error(err.message)
        })
    }
  },
  created() {}
}
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
