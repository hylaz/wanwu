<template>
  <div class="add-dialog">
    <el-dialog
      :title="title"
      :visible.sync="dialogVisible"
      width="50%"
      :show-close="false"
      :close-on-click-modal="false"
    >
      <div>
        <el-form
          :model="ruleForm"
          :rules="rules"
          ref="ruleForm"
          label-width="130px"
        >
          <el-form-item :label="$t('tool.server.avatar')" prop="avatar">
            <div class="avatar">
              <el-upload
                class="avatar-uploader"
                action=""
                name="files"
                :show-file-list="false"
                :multiple="false"
                :http-request="handleUploadAvatar"
                :on-error="handleUploadError"
                accept=".png,.jpg,.jpeg"
              >
                <div class="echo-img">
                  <img :src="ruleForm.avatar.path ? basePath + '/user/api/' + ruleForm.avatar.path : defaultAvatar"
                       alt=""/>
                </div>
              </el-upload>
            </div>
          </el-form-item>
          <el-form-item :label="$t('tool.server.name')" prop="name">
            <el-input v-model="ruleForm.name"></el-input>
          </el-form-item>
          <el-form-item :label="$t('tool.server.desc')" prop="desc">
            <el-input v-model="ruleForm.desc"></el-input>
          </el-form-item>
        </el-form>
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button @click="handleClose" size="mini">{{ $t('common.button.cancel') }}</el-button>
        <el-button
          type="primary"
          size="mini"
          @click="submitForm('ruleForm')"
          :loading="publishLoading"
        >
          {{ $t('common.button.confirm') }}
        </el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>
import {addServer, editServer} from "@/api/mcp";
import {uploadAvatar} from "@/api/user";

export default {
  data() {
    return {
      dialogVisible: false,
      mcpServerId: "",
      title: '',
      basePath: this.$basePath,
      defaultAvatar: require("@/assets/imgs/mcp_active.svg"),
      ruleForm: {
        MCPServerId: "",
        avatar: {
          key: "",
          path: ""
        },
        name: "",
        desc: ""
      },
      rules: {
        name: [{required: true, message: "请输入服务名称", trigger: "blur"}],
        desc: [{required: true, message: "请输入服务描述", trigger: "blur"}]
      },
      publishLoading: false
    };
  },
  methods: {
    showDialog(item) {
      this.dialogVisible = true
      if (item) {
        this.mcpServerId = item.mcpServerId
        this.ruleForm = item
        this.title = '修改MCP服务'
      } else this.title = '创建MCP服务'
    },
    handleUploadAvatar(data) {
      if (data.file) {
        const formData = new FormData()
        const config = {headers: {"Content-Type": "multipart/form-data"}}
        formData.append('avatar', data.file)
        uploadAvatar(formData, config).then(res => {
          if (res.code === 0) {
            this.ruleForm.avatar = res.data
          }
        })
      }
    },
    handleUploadError() {
      this.$message.error(this.$t('common.message.uploadError'))
    },
    handleClose() {
      this.dialogVisible = false
      this.$emit("handleClose", false)
      this.$refs.ruleForm.resetFields()
      this.$refs.ruleForm.clearValidate()
      this.mcpServerId = ''
      this.ruleForm = {
        MCPServerId: "",
        avatar: {
          key: "",
          path: ""
        },
        name: "",
        desc: ""
      }
    },
    submitForm(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.publishLoading = true
          const params = {
            ...this.ruleForm
          }
          if (this.mcpServerId) editServer(params).then((res) => {
            if (res.code === 0) {
              this.$message.success("发布成功")
              this.$emit("handleFetch", false)
              this.handleClose()
            }
          }).finally(() => this.publishLoading = false)
          else addServer(params).then((res) => {
            if (res.code === 0) {
              this.$message.success("发布成功")
              this.$emit("handleFetch", false)
              this.handleClose()
              this.$router.push({path: `/tool/detail/server?mcpServerId=${res.data.mcpServerId}`})
            }
          }).finally(() => this.publishLoading = false)
        }
      });
    }
  },
};
</script>
<style lang="scss" scoped>
.required-label::after {
  content: '*';
  position: absolute;
  color: #eb0a0b;
  font-size: 20px;
  margin-left: 4px;
}

.add-dialog {
  .el-button.is-disabled {
    &:active {
      background: transparent !important;
      border-color: #ebeef5 !important;
    }
  }

  .avatar {
    display: flex;
    align-items: center;
    margin-top: 4px;

    .avatar-uploader {
      width: 32px;
      height: 32px;
      flex-shrink: 0;

      /deep/ {
        .el-upload {
          width: 100%;
          height: 100%;
          border-radius: 6px;
          border: 1px solid #DCDFE6;
          overflow: hidden;
        }

        .echo-img {
          width: 100%;
          height: 100%;
          position: relative;

          img {
            object-fit: cover;
            height: 100%;
          }
        }
      }
    }

    .row {
      flex: 1;
    }
  }
}
</style>