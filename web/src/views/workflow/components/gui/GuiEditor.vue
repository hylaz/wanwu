<template>
  <div class="edit-api-box">
    <div class="edit-modal">
      <div class="edit-header">
        <span>编辑API</span>
        <i @click="closeEdit" class="el-icon-d-arrow-right arrow-icon"></i>
      </div>
      <div class="req-box">
        <!--url-->
        <div class="req-url">
          <el-select
            v-model="settings.http_method"
            popper-class="workflow-select"
          >
            <el-option label="GET" value="GET"></el-option>
            <el-option label="POST" value="POST"></el-option>
          </el-select>
          <el-input v-model="settings.url"></el-input>
          <el-button class="send-bt" type="primary" size="mini" @click="preSend"
            >发 送</el-button
          >
        </div>
        <!--params-->
        <div class="req-params">
          <div class="req-params-tabs">
            <el-radio-group size="small" v-model="reqType" class="http-type">
              <el-radio-button label="req">请求</el-radio-button>
              <el-radio-button label="res">响应</el-radio-button>
            </el-radio-group>

            <el-tabs
              v-if="reqType === 'req'"
              v-model="reqItem"
              @tab-click="reqTabClick"
            >
              <el-tab-pane label="Params" name="Params" class="rl">
                <p class="tab-content-title">Query Params</p>
                <br />
                <el-table
                  class="noborder-table"
                  :data="paramsTable"
                  style="width: 100%"
                  :header-cell-style="{
                    background: '#F9F9F9',
                    color: '#999999',
                  }"
                >
                  <el-table-column prop="name" label="参数名">
                    <template slot-scope="scope">
                      <el-input size="mini" v-model="scope.row.name"></el-input>
                    </template>
                  </el-table-column>
                  <el-table-column prop="value" label="MOCK值">
                    <template slot-scope="scope">
                      <el-input
                        size="mini"
                        v-model="scope.row.newValue"
                      ></el-input>
                    </template>
                  </el-table-column>
                  <el-table-column fixed="right" label="操作" width="100">
                    <template slot-scope="scope">
                      <i
                        class="el-icon-remove-outline"
                        v-if="scope.$index !== paramsTable.length - 1"
                        @click="preDelParams(scope.$index)"
                      ></i>
                    </template>
                  </el-table-column>
                </el-table>
              </el-tab-pane>
              <el-tab-pane label="Headers" name="Headers">
                <p class="tab-content-title">Headers</p>
                <br />
                <el-table
                  class="noborder-table"
                  :data="headersTable"
                  style="width: 100%"
                  :header-cell-style="{
                    background: '#F9F9F9',
                    color: '#999999',
                  }"
                >
                  <el-table-column prop="name" label="参数名">
                    <template slot-scope="scope">
                      <el-input size="mini" v-model="scope.row.key"></el-input>
                    </template>
                  </el-table-column>
                  <el-table-column prop="value" label="参数值">
                    <template slot-scope="scope">
                      <el-input
                        size="mini"
                        v-model="scope.row.value"
                      ></el-input>
                    </template>
                  </el-table-column>
                  <el-table-column fixed="right" label="操作" width="100">
                    <template slot-scope="scope">
                      <i
                        class="el-icon-remove-outline"
                        v-if="scope.$index !== headersTable.length - 1"
                        @click="preDelHeaders(scope.$index)"
                      ></i>
                    </template>
                  </el-table-column>
                </el-table>
              </el-tab-pane>
              <!-- <el-tab-pane label="Authorization" name="Authorization">
                <el-form
                  ref="form"
                  class="authorization-form"
                  size="mini"
                  :model="AuthorizationForm"
                  label-width="100px"
                >
                  <el-form-item label="鉴权方式">
                    <el-radio-group v-model="AuthorizationForm.type">
                      <el-radio label="apikey">API Key</el-radio>
                    </el-radio-group>
                  </el-form-item>
                  <el-form-item label="密钥位置">
                    <el-radio-group v-model="AuthorizationForm.location">
                      <el-radio label="header">Header</el-radio>
                    </el-radio-group>
                  </el-form-item>
                  <el-form-item label="密钥参数名">
                    <el-input v-model="AuthorizationForm.secretKey"></el-input>
                  </el-form-item>
                  <el-form-item label="密钥值">
                    <el-input
                      v-model="AuthorizationForm.secretValue"
                    ></el-input>
                  </el-form-item>
                </el-form>
              </el-tab-pane> -->
              <el-tab-pane
                label="Body"
                name="Body"
                :disabled="settings.http_method != 'POST'"
              >
                <p class="tab-content-title">Body（json）</p>
                <br />
                <el-table
                  class="noborder-table"
                  :data="bodyTable"
                  style="width: 100%"
                  :header-cell-style="{
                    background: '#F9F9F9',
                    color: '#999999',
                  }"
                >
                  <el-table-column prop="name" label="参数名">
                    <template slot-scope="scope">
                      <el-input size="mini" v-model="scope.row.name"></el-input>
                    </template>
                  </el-table-column>
                  <el-table-column prop="value" label="MOCK值">
                    <template slot-scope="scope">
                      <el-input
                        size="mini"
                        v-model="scope.row.newValue"
                      ></el-input>
                    </template>
                  </el-table-column>
                  <el-table-column fixed="right" label="操作" width="100">
                    <template slot-scope="scope">
                      <i
                        class="el-icon-remove-outline"
                        v-if="scope.$index !== bodyTable.length - 1"
                        @click="preDelParams(scope.$index)"
                      ></i>
                    </template>
                  </el-table-column>
                </el-table>
              </el-tab-pane>
            </el-tabs>

            <div
              v-if="reqType === 'res'"
              style="padding-top: 50px; height: 97%"
            >
              <el-table
                class="res-table"
                height="100%"
                :data="nodeData.data.outputs"
                :header-cell-style="{
                  background: '#F9F9F9',
                  color: '#999999',
                }"
                style="width: 100%"
              >
                <el-table-column prop="name" label="参数名">
                  <template slot-scope="scope">
                    <el-input size="mini" v-text="scope.row.name"></el-input>
                  </template>
                </el-table-column>
                <el-table-column prop="value" label="参数类型">
                  <template slot-scope="scope">
                    <el-input size="mini" v-text="scope.row.type"></el-input>
                  </template>
                </el-table-column>
                <el-table-column fixed="right" label="操作" width="100">
                  <template slot-scope="scope">
                    <i
                      class="el-icon-remove-outline"
                      @click="preDelResParams(scope.$index)"
                    ></i>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </div>
        </div>
      </div>

      <!--响应结果-->
      <div class="res-box">
        <div class="res-header">
          <span>响应结果</span>
          <div class="right-status" v-if="mockResData.code === 0">
            <div class="status-item">
              <span>请求状态码: </span><span>200</span>
            </div>
            <div class="status-item">
              <span>content type: </span
              ><span>application/json;charset=UTF-8</span>
            </div>
          </div>
        </div>
        <div class="res-content" v-loading="resLoading">
          <div class="left">
            <el-tabs v-model="resItem" @tab-click="resTabClick">
              <el-tab-pane label="Response" name="Response">
                <div class="editable--input" v-show="mockResDataStr">
                  <!--<codeEditor
                                            style="height: 400px;overflow: auto"
                                            :value="mockResDataStr"
                                            :language="'json'"
                                            :readOnly="false"
                                            :theme="'vs'"
                                    ></codeEditor>-->
                  {{ mockResData }}
                </div>
              </el-tab-pane>
            </el-tabs>
            <el-button
              class="parse-res"
              type="primary"
              size="mini"
              @click="parseResData"
            >
              <i class="el-icon-refresh"></i>
              解析到输出参数</el-button
            >
          </div>
          <div class="right">
            <el-tabs>
              <el-tab-pane label="节点输出参数" name="节点输出参数">
                <div
                  class="editable--input"
                  contenteditable="true"
                  v-show="mockResDataStr"
                >
                  {{ mockResData.result }}
                </div>
              </el-tab-pane>
            </el-tabs>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { runNode } from "@/api/workflow";
import codeEditor from "@/views/codeEditor/components/codeEditor.vue";
import { md } from "@/mixins/marksown-it";

const paramsModel = {
  desc: "",
  extra: {
    location: "query",
  },
  list_schema: null,
  name: "",
  object_schema: null,
  required: false,
  type: "string",
  value: {
    content: {
      ref_node_id: "",
      ref_var_name: "",
    },
    type: "ref",
  },
};
const bodyModel = {
  desc: "",
  extra: {
    location: "body",
  },
  list_schema: null,
  name: "",
  object_schema: null,
  required: false,
  type: "string",
  value: {
    content: {
      ref_node_id: "",
      ref_var_name: "",
    },
    type: "ref",
  },
};

export default {
  props: ["node"],
  components: { codeEditor },
  data() {
    return {
      editVisible: false,
      nodeData: {},
      settings: {},
      //请求tabs
      reqType: "req",
      reqItem: "Params",
      resItem: "Response",
      resLoading: false,
      //参数表格
      paramsTable: [JSON.parse(JSON.stringify(paramsModel))],
      //headers
      headersTable: [{ key: "", value: "" }],
      bodyTable: [JSON.parse(JSON.stringify(bodyModel))],
      //Authorization
      AuthorizationForm: {
        type: "apikey",
        location: "header",
        secretKey: "key",
        secretValue: "77b5f0d102c848d443b791fd469b732d",
      },
      //响应
      resParamsTable: [],
      //单个节点调试结果
      mockResData: {},
      mockResDataStr: "",
    };
  },
  watch: {
    // 'nodeData.data.inputs': {
    //     handler:function (newVal, oldVal) {
    //         if(this.nodeData.data.inputs[this.nodeData.data.inputs.length-1].name){
    //             this.preAddTable('Params')
    //         }
    //     },
    //     deep:true
    // },
    headersTable: {
      handler: function (newVal, oldVal) {
        if (this.headersTable[this.headersTable.length - 1].key) {
          this.preAddTable("Headers");
        }
      },
      deep: true,
    },
    bodyTable: {
      handler: function (newVal, oldVal) {
        if (this.bodyTable[this.bodyTable.length - 1].name) {
          this.preAddTable("Body");
        }
        this.setInputs(true);
      },
      deep: true,
    },
    paramsTable: {
      handler: function (newVal, oldVal) {
        if (this.paramsTable[this.paramsTable.length - 1].name) {
          this.preAddTable("Params");
        }
        this.setInputs(this.settings.http_method == "POST");
      },
      deep: true,
    },
    "settings.http_method": {
      handler: function (newVal, oldVal) {
        this.reqItem = "Params";
        this.setInputs(this.settings.http_method == "POST");
      },
    },
  },
  created() {
    //this.mockResData = {}
    //this.mockResDataStr = ''
    this.nodeData = this.node.data;
    this.settings = this.nodeData.data.settings;
    this.nodeData.data.inputs = this.nodeData.data.inputs.map((n) => {
      return {
        ...n,
        newValue: n.value.type === "generated" ? n.value.content : "",
      };
    });
    this.setParasm();
  },
  methods: {
    setParasm() {
      let params = [];
      let body = [];
      this.nodeData.data.inputs.map((n) => {
        if (n.extra.location == "body") {
          body.push(n);
        } else {
          params.push(n);
        }
      });
      this.paramsTable =
        params.length == 0 ? [JSON.parse(JSON.stringify(paramsModel))] : params;
      this.bodyTable =
        body.length == 0 ? [JSON.parse(JSON.stringify(bodyModel))] : body;
    },
    setInputs(hasBody) {
      let result = [];
      result = result.concat(
        this.paramsTable.map((n) => {
          return {
            ...n,
          };
        })
      );
      if (hasBody) {
        result = result.concat(
          this.bodyTable.map((n) => {
            return {
              ...n,
            };
          })
        );
      }
      this.nodeData.data.inputs = result;
    },
    parseResData() {
      this.reqType = "res";
      this.resParamsTable = [];
      let obj = this.mockResData;
      for (let key in obj) {
        this.resParamsTable.push({
          name: key,
          value: {
            content: "",
            type: "generated",
          },
          type: typeof obj[key],
        });
      }
      this.nodeData.data.outputs = this.resParamsTable;
    },
    preDelResParams(index) {
      // let tempList = JSON.parse(JSON.stringify(this.resParamsTable))
      // tempList.splice(index,1)
      // //this.resParamsTable = JSON.parse(JSON.stringify(tempList))
      // this.nodeData.data.outputs = JSON.parse(JSON.stringify(tempList))
      this.nodeData.data.outputs.splice(index, 1);
    },
    preDelParams(index) {
      // this.nodeData.data.inputs.splice(index,1)
      this.paramsTable.splice(index, 1);
      //this.$emit('inputsChange',this.paramsTable)
    },
    preDelHeaders(index) {
      this.headersTable.splice(index, 1);
    },
    preDelBody(index) {
      this.bodyTable.splice(index, 1);
    },
    async preSend() {
      let headers = {};
      this.headersTable.forEach((n) => {
        if (n.key && n.value) {
          headers[n.key] = n.value;
        }
      });

      let nodeSchema = {
        nodeSchema: {
          nodes: [
            {
              id: this.node.data.id,
              type: this.node.data.type,
              name: this.node.data.name,
              data: {
                settings: {
                  ...this.nodeData.data.settings,
                  headers: headers,
                },
                inputs: this.nodeData.data.inputs
                  .filter((m) => {
                    return m.name && m.newValue;
                  })
                  .map((n, i) => {
                    return {
                      name: n.name,
                      type: n.type,
                      desc: "",
                      value: {
                        type: "generated",
                        content: n.newValue,
                      },
                      // 注意这里写死null，因为api节点的input传参定义，只支持第一层级的解析定义结构
                      object_schema: null,
                      // 注意这里写死null，因为api节点的input传参定义，只支持第一层级的解析定义结构
                      list_schema: null,
                      extra: n.extra,
                    };
                  }),
                // 在单个节点调试运行时，其output出参要求为空数组，与整个workflow画布调试运行时的场景，作区分
                outputs: [], //this.nodeData.data
              },
            },
          ],
        },
      };
      console.log(nodeSchema);
      this.resLoading = true;
      let res = await runNode(nodeSchema);
      if (res.code === 0) {
        this.resLoading = false;
        this.mockResData = res.data;
        this.ParseResData(res.data);
      }
    },
    ParseResData(data) {
      this.mockResDataStr = JSON.stringify(data);
    },
    reqTabClick(val) {},
    resTabClick(val) {},
    preAddTable(type) {
      switch (type) {
        case "Params":
          let ParamsItem = JSON.parse(JSON.stringify(paramsModel));
          // this.$set(this.nodeData.data.inputs,this.nodeData.data.inputs.length,ParamsItem)
          //通知父级更新数据
          //this.$emit('inputsChange',this.paramsTable)
          this.$set(this.paramsTable, this.paramsTable.length, ParamsItem);
          break;
        case "Headers":
          let HeadersItem = { key: "", value: "" };
          this.$set(this.headersTable, this.headersTable.length, HeadersItem);
          //同步到settings里
          let _headers = {};
          this.headersTable.forEach((n) => {
            if (n.key) {
              _headers[n.key] = n.value;
            }
          });
          this.nodeData.data.settings.headers = _headers;
          break;
        case "Body":
          let BodyItem = JSON.parse(JSON.stringify(bodyModel));
          this.$set(this.bodyTable, this.bodyTable.length, BodyItem);
          break;
      }
    },
    closeEdit() {
      this.$emit("close");
    },
  },
};
</script>


<style lang="scss">
@import "@/style/workflow.scss";

.edit-api-box {
  position: fixed;
  height: calc(100vh - 40px);
  overflow-y: auto;
  right: 400px;
  left: 0;
  top: 20px;
  bottom: 20px;
  background-color: #070c1480;
  z-index: 10;
  overflow: auto;
  .info-title {
    font-size: 16px;
    margin: 10px 0;
  }
  .info-item {
    display: flex;
    margin: 15px 0;
    label {
      display: block;
      flex: 1;
      color: #5c5f66;
    }
    span {
      display: block;
      flex: 3;
    }
  }
  .edit-api {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    color: #151b26;
    border: 1px solid #e8e9eb;
    padding: 8px 0;
    border-radius: 6px;
    cursor: pointer;
    i {
      margin-right: 6px;
    }
  }
}
.edit-modal {
  float: right;
  width: 75%;
  height: 100%;
  background: #fff;
  overflow-y: auto;
  font-size: 12px;
  .edit-header {
    position: relative;
    padding: 14px 20px;
    border-bottom: 1px solid #eee;
    .arrow-icon {
      position: absolute;
      right: 20px;
    }
  }

  .req-box {
    flex: 1;
    padding: 14px 20px;
    height: 450px;
    box-sizing: border-box;
    .req-url {
      display: flex;
      .el-select/deep/.el-input__inner {
        background-color: #f7f7f9;
        font-size: 12px;
      }
      .el-input {
        margin-left: -1px;
      }
      .send-bt {
        margin: 0 5px;
      }
    }
    .req-params {
      position: relative;
      height: calc(100% - 27px);
      padding-top: 10px;
      .req-params-tabs {
        height: 100%;
        .http-type {
          position: absolute;
          top: 15px;
          z-index: 10;
        }
      }
      /deep/.el-radio-button {
        &:hover {
          .el-radio-button__inner {
            color: #606266;
          }
        }
      }
      /deep/.is-active {
        &:hover {
          .el-radio-button__inner {
            color: #fff;
          }
        }
        .el-radio-button__inner {
          background: #e60001;
          border-color: #e60001;
          box-shadow: 1px 0 0 0 #e60001;
        }
      }
      /deep/.res-table {
        // height: 100% !important;
        .cell {
          font-size: 12px;
        }
      }
      /deep/.el-tabs {
        height: 100%;
        margin-top: 0;
        .is-active {
          color: #e60001;
        }
        .el-tabs__active-bar {
          border-color: #e60001;
          background-color: #e60001;
        }
      }
      /deep/.el-tabs__nav {
        margin-left: 150px;
      }
      /deep/.el-tabs__content {
        height: calc(100% - 60px);
        margin-top: 10px;
        padding: 5px 10px;
        border-radius: 5px;
        overflow: auto;
      }
      /deep/.el-tabs__nav-wrap::after {
        height: 0;
      }
    }
  }
  .res-box {
    height: calc(100vh - 536px);
    border-top: 1px solid #eee;
    font-size: 12px;
    .res-header {
      padding: 10px 20px;
      border-bottom: 1px solid #eee;
      .right-status {
        float: right;
        display: flex;
        .status-item {
          margin-left: 10px;
          vertical-align: middle;
          span:nth-child(1) {
            color: #84868c;
          }
          span:nth-child(2) {
            color: #333;
          }
          span {
            font-size: 12px;
          }
        }
      }
    }
    .res-content {
      display: flex;
      height: calc(100% - 42px);

      .el-tabs {
        margin-top: 0 !important;
      }
      .left,
      .right {
        width: 50%;
        flex: 1;
        /deep/.el-tabs__nav-wrap {
          padding: 0 20px;
        }
        /deep/.el-tabs {
          height: 100%;
        }
      }
      /deep/.el-tabs__content {
        height: calc(100% - 39px);
        overflow: auto;
      }
      .left {
        position: relative;
        border-right: 1px solid #e4e7ed;

        /deep/.parse-res {
          position: absolute;
          right: 15px;
          top: 7px;
          font-weight: bold;
          background: transparent !important;
          border: 0 !important;
          border-radius: 0 !important;
          padding: 5px !important;
          margin-left: 0 !important;
          color: #0079bf;
          span {
            font-size: 12px;
          }
          &:hover {
            border-radius: 4px !important;
            background: #0000000d !important;
          }
        }
      }
      .right {
        /deep/.is-active {
          color: #303133;
        }
        /deep/.el-tabs__active-bar {
          border-color: transparent;
          background-color: transparent;
        }
        /deep/.el-tabs__item {
          color: #303133;
        }
      }
    }
    .editable--input {
      padding: 20px;
      min-height: 100%;
      word-break: break-all;
    }
  }
}

.authorization-form {
  margin-top: 20px;

  /deep/.el-radio__label {
    font-size: 12px;
  }
}
.authorization-form /deep/.el-form-item__label {
  text-align: left !important;
  font-size: 12px;
}

.el-tabs {
  margin-top: 10px;
  /deep/.el-tabs__header {
    margin-bottom: 0;
  }
}
.tab-content-title {
  margin-top: 5px;
  color: #555;
}

.noborder-table /deep/ .el-input__inner {
  border: none !important;
}
</style>
