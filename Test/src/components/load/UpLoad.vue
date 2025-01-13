<template>
  <el-card>
    <div class="title">
      <img src="../../assets/cloud-upload.png" />
      <h3>上传xlsx文件</h3>
    </div>
    <el-divider></el-divider>
    <div class="body">
      <div>
        <span>选择要上传文件的类型：</span>

        <el-radio v-model="system" label="小时制" border
          >2018-2020年数据（小时制）</el-radio
        >
        <el-radio v-model="system" label="月度制" border
          >汕头练江水站年数据汇总（月度制）</el-radio
        >
      </div>
      <el-divider></el-divider>
      <el-divider></el-divider>
      <el-upload
        ref="uploadElement"
        drag
        action="none"
        multiple
        :show-file-list="false"
        :auto-upload="false"
        :on-change="checkFileType"
      >
        <i class="el-icon-upload"></i>
        <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
      </el-upload>
    </div>
    <div class="dialog">
      <el-dialog
        title="上传失败文件列表"
        :visible.sync="isfailToUpload"
        width="400px"
      >
        <el-table :data="failToUpload">
          <el-table-column
            property="filename"
            label="文件名称"
            width="150"
          ></el-table-column>
          <el-table-column
            property="reason"
            label="失败原因"
            width="150"
          ></el-table-column>
        </el-table>
      </el-dialog>
    </div>
  </el-card>
</template>

<script>
import { upload } from '@/api/file'
export default {
  data() {
    return {
      filelist: [],
      failToUpload: [],
      isfailToUpload: false,
      system: '小时制'
    }
  },
  methods: {
    // 检查用户上传的文件类型，file是当前文件，filelist是所有添加过的文件包括当前文件
    checkFileType(file, filelist) {
      // 截取文件类型
      const fileType = file.name.substring(file.name.lastIndexOf('.') + 1)
      if (fileType !== 'xls' && fileType !== 'xlsx' && fileType !== 'csv') {
        alert('文件格式不正确，请上传.xls、.xlsx或.csv文件！')
        // 如果文件格式不正确，将文件列表中最后一个文件删除
        filelist = filelist.pop()
        return false
      }
      this.file = file
      this.filelist = filelist
      this.uploadFile()
      this.$refs.uploadElement.clearFiles()
    },
    uploadFile() {
      //
      this.filelist.forEach(item => {
        const formData = new FormData()
        formData.append('file', item.raw)
        const config = {
          headers: {
            Authorization: 'Bearer ' + localStorage.getItem('token')
          }
        }
        upload(this.system, formData, config)
          .then(res => {
            if (res.data.code === 200) {
              this.$message.success(`${item.name} 上传成功`)
            } else if (res.data.code === 400) {
              this.$message.warning(`${item.name} ${res.data.msg}，请重新上传`)
              this.isfailToUpload = true
              this.failToUpload.push({
                filename: item.name,
                reason: res.data.msg
              })
            } else {
              this.$message.warning(res.data.msg)
            }
          })
          .catch(err => {
            console.log(err)
            this.$message.error(err.message)
          })
      })

      // 调用 axios，向服务器上传文件
    }
  }
}
</script>

<style lang="less" scoped>
.el-card {
  width: 1022px;
  background-color: #f2f6fc;
  margin: 0;
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
  .body {
    /deep/.el-upload-dragger {
      width: 984px;
    }
    .el-radio {
      margin-left: 50px;
      margin-right: 0;
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
