<template>
  <el-dialog
    :title="title"
    :visible.sync="dialogBasicVisible"
    width="50%"
    :before-close="cancel">
    <div class="action">
      <el-form
        :model="form"
        :rules="rules"
        ref="form">
        <el-col :span="24" class="left-col">
          <div class="action-form">
            <div class="block prompt-box" v-show="!dialogDetailVisible">
              <p class="block-title required-label rl">{{ dialogToolVisible ? $t('tool.custom.app') : $t('tool.custom.tool') }}</p>
              <el-form-item prop="name">
                <el-input class="name-input" v-model="form.name" :placeholder="$t('common.input.placeholder') + $t('tool.custom.tool')"></el-input>
              </el-form-item>
            </div>
            <div class="block prompt-box" v-show="!dialogToolVisible">
              <p class="block-title required-label rl">{{ $t('tool.custom.desc') }}</p>
              <div v-show="dialogDetailVisible">{{ form.description }}</div>
              <el-form-item prop="description" v-show="!dialogDetailVisible">
                <el-input class="name-input" v-model="form.description" :placeholder="$t('common.input.placeholder') + $t('tool.custom.desc')"></el-input>
              </el-form-item>
            </div>
            <div class="block prompt-box" v-show="!dialogDetailVisible">
              <p class="block-title required-label rl">{{ $t('tool.custom.apiAuth') }}</p>
              <div class="rl" @click="preAuthorize">
                <el-form-item prop="apiAuth">
                  <div class="api-key">{{ form.apiAuth.type }}</div>
                </el-form-item>
                <img class="auth-icon" src="@/assets/imgs/auth.png"/>
              </div>
            </div>

            <div class="block prompt-box" v-show="!dialogDetailVisible">
              <p class="block-title required-label rl">Schema</p>
              <div class="rl">
                <div class="flex" style="margin-bottom: 10px">
                  <el-select v-model="example" :placeholder="$t('tool.custom.schema')" style="width:100%;"
                             @change="exampleChange">
                    <!--<el-option label="模板样例导入" value="json"></el-option>-->
                    <el-option :label="'JSON' + $t('tool.custom.example')" value="json"></el-option>
                    <el-option :label="'YAML' + $t('tool.custom.example')" value="yaml"></el-option>
                  </el-select>
                </div>
                <el-form-item prop="schema">
                  <el-input class="schema-textarea" v-model="form.schema" @blur="listenerSchema"
                            :placeholder="$t('tool.custom.schemaHint')"
                            type="textarea"></el-input>
                </el-form-item>
              </div>
            </div>

            <div class="block prompt-box">
              <p class="block-title required-label rl">{{ $t('tool.custom.api') }}</p>
              <div class="api-list">
                <el-form-item prop="apiTable">
                  <el-table
                    ref="apiTable"
                    :data="apiList"
                    border
                    size="mini"
                    class="api-table"
                    :header-cell-style="{ textAlign: 'center' }"
                  >
                    <el-table-column
                      v-if="dialogToolVisible"
                      type="selection"
                      width="55"
                      align="center">
                    </el-table-column>
                    <el-table-column
                      prop="name"
                      label="Name">
                    </el-table-column>
                    <el-table-column
                      prop="method"
                      label="Method">
                    </el-table-column>
                    <el-table-column
                      prop="path"
                      label="Path">
                    </el-table-column>
                  </el-table>
                </el-form-item>
              </div>
            </div>
            <div class="block prompt-box" v-show="!dialogDetailVisible">
              <p class="block-title rl">{{ $t('tool.custom.privacy') }}</p>
              <el-form-item prop="privacyPolicy">
                <el-input class="name-input" v-model="form.privacyPolicy"
                          :placeholder="$t('tool.custom.privacyHint')"></el-input>
              </el-form-item>
            </div>
          </div>
        </el-col>
      </el-form>

      <!--认证弹窗-->
      <el-dialog
        :title="$t('tool.custom.auth.title')"
        :visible.sync="dialogAuthVisible"
        width="600px"
        append-to-body
        :close-on-click-modal="false"
        @close="beforeApiAuthClose"
      >
        <div class="action-form">
          <el-form :rules="apiAuthRules" ref="apiAuthForm" :inline="false" :model="form.apiAuth">
            <el-form-item :label="$t('tool.custom.auth.type')">
              <el-radio-group v-model="form.apiAuth.type">
                <el-radio label="None">None</el-radio>
                <el-radio label="API Key">API Key</el-radio>
                <!--<el-radio label="3">OAuth</el-radio>-->
              </el-radio-group>
            </el-form-item>
            <!--API Key-->
            <div v-show="form.apiAuth.type === 'API Key'">
              <el-form-item label="API key" prop="apiKey">
                <el-input class="desc-input " v-model="form.apiAuth.apiKey" placeholder="API key" clearable></el-input>
              </el-form-item>
              <el-form-item :label="$t('tool.custom.auth.authType')">
                <el-radio-group v-model="form.apiAuth.authType">
                  <!--<el-radio label="1">Basic</el-radio>
                  <el-radio label="2">Bearer</el-radio>-->
                  <el-radio label="Custom">Custom</el-radio>
                </el-radio-group>
              </el-form-item>
              <el-form-item label="Custom Header Name" prop="customHeaderName">
                <el-input class="desc-input " v-model="form.apiAuth.customHeaderName" placeholder="Custom Header Name"
                          clearable></el-input>
              </el-form-item>
            </div>
          </el-form>
        </div>
        <span slot="footer" class="dialog-footer">
                <el-button @click="beforeApiAuthClose">{{ $t('common.button.cancel') }}</el-button>
                <el-button type="primary" @click="listenerApiKey">{{ $t('common.button.confirm') }}</el-button>
            </span>
      </el-dialog>
    </div>
    <span slot="footer" class="dialog-footer" v-show="!dialogDetailVisible">
        <el-button @click="cancel">{{ $t('common.button.cancel') }}</el-button>
        <el-button
          type="primary"
          @click="submit"
          :loading="loading">{{ $t('common.button.confirm') }}</el-button>
    </span>
    <span slot="footer" class="dialog-footer" v-show="dialogDetailVisible">
        <el-button
          type="primary"
          @click="dialogDetailVisible = false; title = $t('tool.custom.editTitle')">{{$t('common.button.edit')}}</el-button>
    </span>
  </el-dialog>
</template>
<script>
import {getCustom, addCustom, editCustom, getSchema, addOpenapi} from "@/api/mcp";
import {schemaConfig} from '@/utils/schema.conf';

export default {
  data() {
    const validateApiAuthFields = (rule, value, callback) => {
      if (this.form.apiAuth.type === 'API Key' &&
        (!this.form.apiAuth.apiKey || !this.form.apiAuth.customHeaderName)) {
        callback(new Error(rule.message));
      } else {
        callback();
      }
    }
    const validateApiTableFields = (rule, value, callback) => {
      if (this.dialogToolVisible && this.$refs.apiTable.selection.length === 0) {
        callback(new Error(rule.message));
      } else {
        callback();
      }
    }
    return {
      dialogBasicVisible: false,
      dialogDetailVisible: false,
      dialogToolVisible: false,
      title: '',
      apiList: [],
      example: '',
      form: {
        description: '',
        customToolId: '',
        mcpServerId: '',
        name: '',
        schema: '',
        privacyPolicy: '',
        apiAuth: {
          type: 'None',
          authType: 'Custom',
          apiKey: '',
          customHeaderName: '',
        }
      },
      //认证表单
      dialogAuthVisible: false,
      rules: {
        description: [{required: true, message: this.$t('common.input.placeholder'), trigger: 'blur'}],
        name: [{required: true, message: this.$t('common.input.placeholder'), trigger: 'blur'}],
        schema: [{required: true, message: this.$t('common.input.placeholder'), trigger: 'blur'}],
        apiAuth: [{validator: validateApiAuthFields, message: this.$t('tool.custom.apiAuthHint'), trigger: 'blur'}],
        apiTable: [{validator: validateApiTableFields, message: this.$t('tool.custom.apiHint'), trigger: 'blur'}],
      },
      apiAuthRules: {
        apiKey: [{required: true, message: this.$t('common.input.placeholder'), trigger: 'blur'}],
        customHeaderName: [{required: true, message: this.$t('common.input.placeholder'), trigger: 'blur'}],
      },
      schemaConfig: schemaConfig,
      loading: false
    }
  },
  methods: {
    showDialog(customToolId, dialogDetailVisible) {
      this.dialogDetailVisible = dialogDetailVisible
      this.dialogBasicVisible = true
      if (customToolId) {
        if (!dialogDetailVisible) this.title = this.$t('tool.custom.editTitle')
        const params = {
          customToolId: customToolId
        }
        getCustom(params)
          .then((res) => {
            const {list, ...form} = res.data
            this.form = form
            if (dialogDetailVisible) {
              this.title = form.name
            }
            this.listenerSchema()
          })
      } else this.title = this.$t('tool.custom.addTitle')
    },
    showToolDialog(mcpServerId) {
      this.form.mcpServerId = mcpServerId
      this.dialogBasicVisible = true
      this.dialogToolVisible = true
      this.form.description = ' '
      this.title = this.$t('tool.custom.toolTitle')
    },
    exampleChange(value) {
      this.form.schema = this.schemaConfig[value]
      this.listenerSchema()
    },
    beforeApiAuthClose() {
      this.$refs.apiAuthForm.clearValidate()
      this.dialogAuthVisible = false
    },
    listenerApiKey() {
      this.$refs.apiAuthForm.validate(async (valid) => {
        if (!valid) return;
        this.dialogAuthVisible = false
      })
    },
    listenerSchema() {
      const params = JSON.stringify({
        schema: this.form.schema
      })
      getSchema(params)
        .then((res) => {
          this.apiList = res.data.list || []
        })
    },
    preAuthorize() {
      this.dialogAuthVisible = true
    },
    submit() {
      this.$refs.form.validate(async (valid) => {
        if (!valid) return;
        this.loading = true
        if (this.form.apiAuth.type === 'None') {
          this.form.apiAuth.authType = '';
          this.form.apiAuth.apiKey = '';
          this.form.apiAuth.customHeaderName = ''
        }
        const params = {
          ...this.form
        }
        if (this.form.customToolId) {
          editCustom(params)
            .then((res) => {
              if (res.code === 0) {
                this.$message.success(this.$t('common.info.edit'))
                this.$emit("handleFetch", false)
                this.cancel()
              }
            }).finally(() => this.loading = false)
        } else {
          delete params.customToolId
          if (this.dialogToolVisible) {
            params.methodNames = this.$refs.apiTable.selection.map(item => item.name)
            addOpenapi(params).then((res) => {
              if (res.code === 0) {
                this.$message.success(this.$t('common.info.create'))
                this.$emit("handleFetch", false)
                this.cancel()
              }
            }).finally(() => this.loading = false)
          } else {
            addCustom(params).then((res) => {
              if (res.code === 0) {
                this.$message.success(this.$t('common.info.create'))
                this.$emit("handleFetch", false)
                this.cancel()
              }
            }).finally(() => this.loading = false)
          }
        }
      })
    },
    cancel() {
      this.$emit("handleClose", false)
      this.loading = false
      this.dialogBasicVisible = false
      this.dialogDetailVisible = false
      this.dialogToolVisible = false
      this.apiList = []
      this.example = ''
      this.title = ''
      this.$refs.form.clearValidate()
      this.form = {
        description: '',
        customToolId: '',
        mcpServerId: '',
        name: '',
        schema: '',
        privacyPolicy: '',
        apiAuth: {
          type: 'None',
          authType: 'Custom',
          apiKey: '',
          customHeaderName: ''
        }
      }
    }
  },

}
</script>

<style lang="scss" scoped>
/deep/ .el-radio__input.is-checked .el-radio__inner {
  border-color: #D33A3A !important;
  background: transparent !important;
}

/deep/ .el-radio__input.is-checked .el-radio__inner::after {
  background: #eb0a0b !important;
  width: 7px !important;
  height: 7px !important;
}

::selection {
  color: #1a2029;
  background: #c8deff;
}

.left-col {
  // background-color: #fafafa;
  overflow: auto;
  height: 100%;

  .left-col-header {
    width: 100%;
    padding: 30px 40px;
    text-align: center;

    .back-icon {
      position: absolute;
      left: 35px;
      font-size: 14px;
      cursor: pointer;
      border-radius: 20px;
      border: 1px solid #e1e1e1;
      padding: 6px;
      color: #444;

      &:hover {
        font-weight: bold;
      }
    }

    .header-title {
      font-size: 18px;
      font-weight: bold;
      color: #303133;
    }

    .bt-box {
      position: absolute;
      width: 160px;
      height: 30px;
      right: 20px;
      top: 0;
      bottom: 0;
      margin: auto;

      .del-bt {
        margin-left: 10px;
      }
    }
  }

  .action-form {
    padding: 0 40px;

    /deep/ .schema-textarea {
      .el-textarea__inner {
        height: 200px !important;
      }
    }

    .api-key {
      background-color: transparent !important;
      border: 1px solid #d3d7dd !important;
      padding: 0 15px;
      -webkit-appearance: none;
      background-image: none;
      border-radius: 4px;
      box-sizing: border-box;
      color: #606266;
      display: inline-block;
      height: 40px;
      line-height: 40px;
      outline: 0;
      transition: border-color .2s cubic-bezier(.645, .045, .355, 1);
      width: 100%;
    }

    .auth-icon {
      position: absolute;
      right: 0;
      height: 39px;
      top: 0;
      cursor: pointer;
      border-left: 1px solid #d3d7dd;
      padding: 7px 9px;
    }
  }

}

.right-col {
  height: 100%;
  // background-color: #f6f7f9;
  .right-title {
    line-height: 84px;
    font-size: 18px;
    font-weight: bold;
    text-align: center;
    color: #303133;
  }

  .smart-center {
    min-width: 0;
    height: calc(100% - 84px);
    flex: 1;
    background-size: 100% 100%;
    position: relative;
  }
}

/*通用*/
.action {
  position: relative;
  height: 100%;

  /deep/ .el-input__inner, /deep/ .el-textarea__inner {
    background-color: transparent !important;
    border: 1px solid #d3d7dd !important;
    font-family: 'Microsoft YaHei', Arial, sans-serif;
    padding: 15px;
  }

  .flex {
    width: 100%;
    display: flex;
    justify-content: space-between;
  }

  .block {
    margin-bottom: 20px;

    .block-title {
      line-height: 30px;
      font-size: 15px;
      font-weight: bold;
    }

    .required-label::after {
      content: '*';
      position: absolute;
      color: #eb0a0b;
      font-size: 20px;
      margin-left: 4px;
    }

    .block-tip {
      color: #919eac;
    }
  }

  .el-input__count {
    color: #909399;
    // background: #fafafa;
    position: absolute;
    font-size: 12px;
    bottom: 5px;
    right: 10px;
  }
}

.action-form /deep/ .el-form-item__label {
  display: block;
  width: 100%;
  text-align: left;
}

.api-list {
  .api-table /deep/ .el-table__body tr td,
  .api-table /deep/ .el-table__header tr th,
  .api-table /deep/ .el-table__body tr:hover > td {
    background-color: transparent !important;
  }
}

</style>
