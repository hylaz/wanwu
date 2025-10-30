<template>
  <div
    class="dag-node"
    :id="node.id"
    tabindex="0"
    @keydown="(e) => handleKeydown(e, nodeData)"
  >
    <div class="node-header">
      <div>
        <i class="el-icon-caret-bottom"></i>
        <img class="node-header-icon" :src="node.store.data.icon" />
        <span class="node-header-name">
          {{ nodeData.name }}
          <!--<span v-if="['StartNode','EndNode'].includes(nodeData.type)">{{nodeData.name}}</span>
          <el-input v-else size="mini" v-model="nodeData.name"></el-input>-->
        </span>
      </div>
      <el-dropdown
        v-if="!['StartNode', 'EndNode', 'EndStreamingNode'].includes(nodeData.type)"
        trigger="click"
        class="controls"
        @command="handleDeleteNode"
      >
        <span class="el-dropdown-link">
          操作<i class="el-icon-arrow-down el-icon--right"></i>
        </span>
        <el-dropdown-menu slot="dropdown">
          <el-dropdown-item>删除</el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </div>
    <!--输入-->
    <!--开始节点-->
    <div class="node-params" v-if="nodeData.type === 'StartNode'">
      <el-collapse v-model="activeName">
        <el-collapse-item name="1">
          <template slot="title">
            <el-row>
              <el-col :span="24">
                <div class="params-type">
                  <span class="params-type-span">输入</span>
                </div>
              </el-col>
            </el-row>
          </template>
          <el-row>
            <el-col :span="12">
              <div class="params-content">
                <div
                  class="params-content-item"
                  v-for="(n, i) in nodeData.data.outputs"
                  :key="`${i}opts`"
                >
                  <span>{{ n.name || "未命名" }}</span>
                  <span>{{ n.type }}</span>
                </div>
              </div>
            </el-col>
          </el-row>
        </el-collapse-item>
      </el-collapse>
      <div class="params-content" v-if="nodeData.data.settings.staticAuthToken">
        <div class="params-content-item">
          <span>token</span>
          <span>{{
            nodeData.data.settings.staticAuthToken.slice(0, 6) + "******"
          }}</span>
        </div>
      </div>
    </div>
    <!--其他节点-->
    <div
      class="node-params"
      v-if="
        [
          'ApiNode',
          'PythonNode',
          'TemplateTransformNode',
          'LLMNode',
          'LLMStreamingNode',
          'RAGNode',
          'GUIAgentNode',
          'FileGenerateNode',
          'FileParseNode',
          'MCPClientNode',
          'IntentionNode',
        ].includes(nodeData.type)
      "
    >
      <el-collapse v-model="activeName">
        <el-collapse-item name="1">
          <template slot="title">
            <el-row>
              <el-col :span="14">
                <div class="params-type">
                  <span class="params-type-span">输入</span>
                </div>
              </el-col>
              <el-col :span="10" v-if="nodeData.type !== 'StartNode'">
                <div class="params-type">
                  <span class="params-type-span">值</span>
                </div>
              </el-col>
            </el-row>
          </template>
          <div v-for="(n, i) in nodeData.data.inputs" :key="`${i}ipts`">
            <el-row
              v-if="
                nodeData.type === 'ApiNode'
                  ? n.name
                    ? true
                    : false
                  : true && !['threshold', 'top_k'].includes(n.name)
              "
            >
              <el-col :span="14">
                <div class="params-content">
                  <div class="params-content-item">
                    <span>{{ n.name || "未命名" }}</span>
                    <span>{{ n.value.type === "ref" ? "引用" : n.type }}</span>
                  </div>
                </div>
              </el-col>
              <el-col :span="10" v-if="nodeData.type !== 'StartNode'">
                <div class="params-content">
                  <div class="params-content-item">
                    <span
                      class="params-content-value"
                      v-if="n.value.type === 'generated'"
                      >{{ n.value.content }}</span
                    >
                    <span
                      class="params-content-value"
                      v-if="n.value.type === 'ref' && n.value.content"
                    >
                      {{
                        n.value.content.ref_node_id
                          ? `${nodeIdMap[n.value.content.ref_node_id]}/${
                              n.value.content.ref_var_name
                            }`
                          : "未选择"
                      }}
                    </span>
                  </div>
                </div>
              </el-col>
            </el-row>
          </div>
        </el-collapse-item>
      </el-collapse>
    </div>

    <!--大模型 -->
    <div
      class="node-params"
      v-if="['LLMNode', 'LLMStreamingNode'].includes(nodeData.type)"
    >
      <div class="params-content">
        <div class="params-content-item">
          <span>模型</span>
          <span>{{ modelName }}</span>
        </div>
        <div class="params-content-item">
          <span>提示词</span>
          <span>{{ nodeData.data.modelForm.input }}</span>
        </div>
      </div>
    </div>
    <!--知识库-->
    <div class="node-params" v-if="['RAGNode'].includes(nodeData.type)">
      <el-collapse v-model="activeName">
        <el-collapse-item name="1">
          <template slot="title">
            <el-row>
              <el-col :span="14">
                <div class="params-type">
                  <span class="params-type-span">知识库</span>
                </div>
              </el-col>
            </el-row>
          </template>
          <div
            v-for="(n, i) in nodeData.data.settings.knowledgeBase"
            :key="`${i}ipts`"
            class="knowledge-span-item"
          >
            <i class="el-icon-document"></i>
            <span class="knowledge-span">
              {{ n }}
            </span>
          </div>
        </el-collapse-item>
      </el-collapse>
    </div>

    <!--分支器-->
    <div v-if="['SwitchNode'].includes(nodeData.type)">
      <div
        v-for="(n, i) in nodeData.data.inputs"
        :key="`${i}ips`"
        class="node-params switch-item"
      >
        <div class="params-type">
          <i class="el-icon-caret-bottom"></i>
          <span class="params-type-span">
            {{
              i === 0
                ? "如果"
                : i === nodeData.data.inputs.length - 1
                ? "否则"
                : "否则如果"
            }}
          </span>
        </div>
        <div class="conditions-box">
          <div
            v-if="n.conditions"
            v-for="(m, j) in n.conditions"
            :key="`${j}cds`"
            class="conditions"
          >
            <el-row class="condition-item" :gutter="8">
              <el-col :span="8" class="left">
                <span
                  class="params-content-value"
                  v-if="m.left.value.type === 'generated'"
                  >{{ m.left.value.content }}</span
                >
                <span
                  class="params-content-value"
                  v-if="m.left.value.type === 'ref' && m.left.value.content"
                >
                  {{
                    m.left.value.content.ref_node_id
                      ? `${nodeIdMap[m.left.value.content.ref_node_id]}/${
                          m.left.value.content.ref_var_name
                        }`
                      : "未选择"
                  }}
                </span>
              </el-col>
              <el-col :span="8" class="center">
                {{ switchOperatorConfig[m.operator] || "请选择" }}
              </el-col>
              <el-col
                :span="8"
                class="right"
                v-show="!disabledArr.includes(m.operator)"
              >
                <span
                  class="params-content-value"
                  v-if="m.right.value.type === 'generated'"
                  >{{ m.right.value.content }}</span
                >
                <span
                  class="params-content-value"
                  v-if="m.right.value.type === 'ref' && m.right.value.content"
                >
                  {{
                    m.right.value.content.ref_node_id
                      ? `${nodeIdMap[m.right.value.content.ref_node_id]}/${
                          m.right.value.content.ref_var_name
                        }`
                      : "未选择"
                  }}
                </span>
              </el-col>
            </el-row>
          </div>
          <!--连接线-->
          <div class="line" v-if="n.conditions && n.conditions.length > 1">
            <span :class="['line-span', n.logic]">
              {{ switchLogicConfig[n.logic] }}
            </span>
          </div>
        </div>
      </div>
    </div>

      <!--意图识别-->
    <div class="node-params" v-if="['IntentionNode'].includes(nodeData.type)">
      <h2>意图</h2>
      <div
        v-for="(n, i) in nodeData.data.settings.intentions"
        :key="`${i}ips`"
        class="switch-item iten-item"
      >
        {{n.name}}
      </div>
    </div>

    <!--输出-->
    <!--结束节点-->
    <div
      class="node-params"
      v-if="nodeData.type === 'EndNode' || nodeData.type === 'EndStreamingNode'"
    >
      <el-collapse v-model="activeName">
        <el-collapse-item name="1">
          <template slot="title">
            <el-row>
              <el-col :span="14">
                <div class="params-type">
                  <span class="params-type-span">输出</span>
                </div>
              </el-col>
              <el-col :span="10">
                <div class="params-type">
                  <span class="params-type-span">值</span>
                </div>
              </el-col>
            </el-row>
          </template>
          <div v-for="(n, i) in nodeData.data.inputs" :key="`${i}nipts`">
            <el-row>
              <el-col :span="14">
                <div class="params-content">
                  <div class="params-content-item">
                    <span>{{ n.name || "未命名" }}</span>
                    <span>{{ n.value.type === 'ref' ? '引用' : n.type }}</span>
                  </div>
                </div>
              </el-col>
              <el-col :span="10" v-if="nodeData.type !== 'StartNode'">
                <div class="params-content">
                  <div class="params-content-item">
                    <span
                      class="params-content-value"
                      v-if="n.value.type === 'generated'"
                      >{{ n.value.content }}</span
                    >
                    <span
                      class="params-content-value"
                      v-if="n.value.type === 'ref' && n.value.content"
                    >
                      {{
                        n.value.content.ref_node_id
                          ? `${nodeIdMap[n.value.content.ref_node_id]}/${
                              n.value.content.ref_var_name
                            }`
                          : "未选择"
                      }}
                    </span>
                  </div>
                </div>
              </el-col>
            </el-row>
          </div>
        </el-collapse-item>
      </el-collapse>

      <!-- <div class="params-content">
          <div class="params-content-item" v-for="(n,i) in nodeData.data.inputs">
              <span>{{n.name}}</span>
              <span>{{n.type}}</span>
          </div>
      </div>-->
    </div>
    <!--其他节点-->
    <div
      class="node-params"
      v-if="
        [
          'ApiNode',
          'PythonNode',
          'TemplateTransformNode',
          'LLMNode',
          'LLMStreamingNode',
          'RAGNode',
          'GUIAgentNode',
          'FileGenerateNode',
          'FileParseNode',
          'MCPClientNode',
          'IntentionNode',
        ].includes(nodeData.type)
      "
    >
      <el-collapse v-model="activeName">
        <el-collapse-item name="2">
          <template slot="title">
            <el-row>
              <el-col :span="24">
                <div class="params-type">
                  <span class="params-type-span">输出</span>
                </div>
              </el-col>
            </el-row>
          </template>
          <div class="params-content">
            <div
              class="params-content-item"
              v-for="(n, i) in nodeData.data.outputs"
              :key="`${i}oupts`"
            >
              <span>{{ n.name || "未命名" }}</span>
              <span>{{ n.type }}</span>
            </div>
          </div>
        </el-collapse-item>
      </el-collapse>
    </div>

    <!--调试结果-->
    <!--分支器节点不需要显示状态-->
    <div class="node-status" v-if="nodeData.node_status ">
      <div class="header">
        <i
          v-if="nodeData.node_status === 'loading'"
          class="el-icon-warning loading status-icon"
        ></i>
        <i
          v-if="nodeData.node_status === 'success'"
          class="el-icon-success success status-icon"
        ></i>
        <i
          v-if="nodeData.node_status === 'failed'"
          class="el-icon-error failed status-icon"
        ></i>
        <i
          v-if="nodeData.node_status === 'init'"
          class="el-icon-warning init status-icon"
        ></i>
        {{ statusObj[nodeData.node_status] }}
      </div>

      <div>
        <div
          v-if="nodeData.node_status === 'loading'"
          style="line-height: 140px; text-align: center"
        >
          <i
            class="el-icon-loading"
            style="color: cornflowerblue; margin: auto; font-size: 26px"
          ></i>
        </div>
        <div class="node-status-content" v-else style="height: 220px; overflow: auto">
          <!--错误提示-->
          <div
            class="params node-message"
            v-if="nodeData.node_status === 'failed'"
          >
            错误信息: {{ nodeData.node_message }}
          </div>
          <!--未运行的节点不显示结果,只显示node_message-->
          <div v-if="nodeData.type !== 'SwitchNode'">
            <div class="params" v-if="nodeData.type === 'StartNode'">
              <p><i class="el-icon-caret-bottom"></i> 输入： <i
                  class="el-icon-document-copy copy-icon"
                  @click="preCopy(nodeData.res_outputs)"
                ></i></p>
              <div><pre v-html="handleFilter(nodeData.res_outputs)"></pre></div>
            </div>
            <div
              class="params"
              v-if="
                !['StartNode', 'EndNode', 'EndStreamingNode'].includes(
                  nodeData.type
                )
              "
            >
              <p><i class="el-icon-caret-bottom"></i> 输入：<i
                  class="el-icon-document-copy copy-icon"
                  @click="preCopy(nodeData.res_inputs)"
                ></i></p>
              <div>
                <pre v-html="handleFilter(nodeData.res_inputs)"></pre>
              </div>
            </div>
            <div
              class="params"
              v-if="
                !['StartNode', 'EndNode', 'EndStreamingNode'].includes(
                  nodeData.type
                )
              "
            >
              <p><i class="el-icon-caret-bottom"></i> 输出：<i
                  class="el-icon-document-copy copy-icon"
                  @click="preCopy(nodeData.res_outputs)"
                ></i></p>
              <pre v-html="handleFilter(nodeData.res_outputs)"></pre>
            </div>
            <div
              class="params"
              v-if="
                nodeData.type === 'EndNode' ||
                nodeData.type === 'EndStreamingNode'
              "
            >
              <p><i class="el-icon-caret-bottom"></i> 输出：<i
                  class="el-icon-document-copy copy-icon"
                  @click="preCopy(nodeData.res_inputs)"
                ></i></p>
              <div>
                <pre v-html="handleFilter(nodeData.res_inputs)"></pre>
              </div>
            </div>
          </div>
          <div class="params" v-else>
            <p><i class="el-icon-caret-bottom"></i> 输出：</p>
            <div>{{ nodeData.node_message }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
//import { mapState } from "vuex";
import { switchLogicConfig, switchOperatorConfig } from "../../mock/nodeConfig";
import { getModels } from "@/api/workflow";
import { getQueryString } from "@/utils/util.js";
export default {
  inject: ["getNode"],
  // props: ["nodeIdMap"],
  data() {
    return {
      isStream: getQueryString("isStream"),
      activeName: [],
      old: [],
      node: {},
      nodeData: {},
      models: [],
      statusObj: {
        success: "成功",
        failed: "失败",
        init: "等待",
        loading: "运行中",
        running_skip: "未运行",
      },
      switchLogicConfig: switchLogicConfig,
      switchOperatorConfig: switchOperatorConfig,
      disabledArr: ["empty", "not_empty"],
      //nodeIdMap:window.localStorage.getItem('nodeIdMap')
    };
  },
  watch: {
    "nodeData.data.inputs": {
      handler: function (newVal, oldVal) {
        //console.log('监听到输入值：',newVal)
      },
      deep: true,
    },
    "nodeData.data.outputs": {
      handler: function (newVal, oldVal) {
        //console.log('监听到输出值：',newVal)
      },
      deep: true,
    },
  },
  computed: {
    nodeIdMap() {
      return JSON.parse(window.localStorage.getItem("nodeIdMap"));
    },
    modelName() {
      let result = "";

      if (this.models && this.models.length) {
        this.models.forEach((item) => {
          if (item.modelId == this.nodeData.data.modelForm.model) {
            result = item.modelName;
          }
        });
      }
      return result;
    },
    /*...mapState({
          nodeIdMap: (state) => state.workflow.nodeIdMap,
      }),*/
  },
  created() {
    this.activeName = sessionStorage.getItem("activeNode").split(",");
    this.old = sessionStorage.getItem("activeNode").split(",");
    const node = this.getNode();
    this.getModels();
    this.node = node;
    this.nodeData = node.data;

    setInterval(() => {
      if (sessionStorage.getItem("activeNode").split(",")[0] !== this.old[0]) {
        this.activeName = sessionStorage.getItem("activeNode").split(",");
        this.old = sessionStorage.getItem("activeNode").split(",");
      }
    }, 100);
  },
  methods: {
    handleKeydown(e, nodeData) {
      if (!['StartNode', 'EndNode', 'EndStreamingNode'].includes(nodeData.type)) {
        if (e.key === 'Backspace' || e.keyCode === 8 || e.key === 'Delete' || e.keyCode === 46) {
          this.handleDeleteNode()
        }
      }
    },
    handleDeleteNode() {
      window.graph.removeNode(this.node.id);
    },
    preCopy(val) {
      this.$copy(JSON.stringify(val));
      this.$message.success("内容已复制到粘贴板");
    },
    getModels() {
      getModels().then((res) => {
        const {list} = res.data || {}
        const models = list ? list.map((item) => ({...item, modelName: item.displayName || item.model})) : []
        this.models = models
      });
    },
    handleFilter(data) {
      return data;
    },
  },
};
</script>

<style lang="scss" scoped>
.node-status {
  position: absolute;
  width: 380px;
  border-radius: 10px;
  border: 1px solid #ddd;
  background-color: #333;
  color: #fff;
  bottom: -255px;
  left: 0;
  right: 0;
  height: 250px;
  .header {
    position: relative;
    line-height: 28px;
    background: #666;
    padding: 0 10px;
    border-radius: 6px 6px 0 0;
    
  }
  .status-icon {
    font-size: 18px;
  }
  .success {
    color: #5a9600;
  }
  .failed {
    color: red;
  }
  .init {
    color: orange;
  }
  .params {
    padding: 10px;
    .copy-icon {
      cursor: pointer;
      &:hover {
        color: cornflowerblue;
      }
    }
  }
  .loading {
    color: cornflowerblue;
  }
}
.dag-node {
  position: relative;
  width: 380px;
  min-height: 150px;
  border-radius: 10px;
  padding: 20px;
  border: 1px solid #ddd;
  background-color: #fff;
}
.node-header {
  display: flex;
  justify-content: space-between;
  .node-header-icon {
    width: 30px;
    height: 20px;
    object-fit: contain;
    border-radius: 10px;
    vertical-align: middle;
    margin: 0 6px 4px 0;
  }
  .node-header-name {
  }
  .controls {
    cursor: pointer;
  }
}
.node-params {
  position: relative;
  padding: 10px;
  background-color: #f9f9fb;
  color: #151b26;
  border-radius: 10px;
  margin-top: 10px;
  h2{
    font-weight:400;
  }
  .iten-item{
    line-height:30px;
    background:#ffffff;
    margin-top:10px;
    padding-left:15px;
    border-radius:3px;
    height: 30px;
  }
  .params-type {
    .params-type-span {
      margin-left: 10px;
    }
  }
  .params-content {
    margin-top: 4px;
    .params-content-item {
      padding: 4px 0;
      display: flex;
      span {
        font-size: 12px;
        display: inline-block;
        vertical-align: middle;
      }
      span:nth-child(1) {
        color: #876300;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        max-width: 106px;
      }
      span:nth-child(2) {
        margin-left: 6px;
        padding: 0 5px;
        white-space: nowrap;
        border-radius: 4px;
        background-color: #e8e9eb;
        color: #5c5f66;
        max-width: 75%;
        overflow: hidden;
        text-overflow: ellipsis;
      }
    }
    .params-content-value {
      display: block;
      box-sizing: border-box;
      width: fit-content;
      max-width: 100%;
      padding: 0 4px;
      /*border: 1px solid #e8e9eb;*/
      border-radius: 4px;
      background-color: #fff;
      color: #333;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
    }
  }
  .value-box {
    border: 1px solid red;
    min-height: 100px;
  }
}

.switch-item {
  position: relative;
  .conditions-box {
    position: relative;
    padding-left: 22px;
    .conditions {
      margin-top: 10px;
      .condition-item {
        .left,
        .right {
          background-color: #fff;
          padding: 2px 6px;
          border-radius: 4px;
          display: inline-block;
          height: 24px;
          font-size: 12px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
        .center {
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          color: #666;
        }
        .center,
        .del-icon {
          text-align: center;
          line-height: 24px;
          font-size: 12px;
        }
      }
    }
    .switch-item-bt-box {
      margin-top: 10px;
      .el-col {
        margin-right: 5px;
      }
      .switch-item-bt {
        color: #2c7eea;
        margin-right: 5px;
      }
      .sort-icon {
        transform: rotate(90deg);
      }
    }
    .line {
      position: absolute;
      top: 14px;
      bottom: 12px;
      left: 7px;
      width: 10px;
      border: 1px solid #d4d6d9;
      border-right: rgba(0, 0, 0, 0);
      border-radius: 4px 0 0 4px;
      .line-span {
        position: absolute;
        top: 50%;
        left: -10px;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 18px;
        height: 18px;
        transform: translateY(-50%);
        border-radius: 4px;
        font-size: 12px;
      }
      .or {
        color: #ff9326;
        background: #fff4e6;
      }
      .and {
        color: #2468f2;
        background: #e6f0ff;
      }
    }
  }
}
/deep/.collapse-transition {
  -webkit-transition: 0s height, 0s padding-top, 0s padding-bottom;
  transition: 0s height, 0s padding-top, 0s padding-bottom;
}

/deep/.horizontal-collapse-transition {
  -webkit-transition: 0s width, 0s padding-left, 0s padding-right;
  transition: 0s width, 0s padding-left, 0s padding-right;
}

/deep/.horizontal-collapse-transition
  .el-submenu__title
  .el-submenu__icon-arrow {
  -webkit-transition: 0s;
  transition: 0s;
  opacity: 0;
}
.knowledge-span-item {
  i {
    color: #03b1ee;
  }
  .knowledge-span {
    margin-left: 4px;
    color: #555;
  }
}

pre {
  width: 100%;
  white-space: pre-wrap;
}
.node-status-content::-webkit-scrollbar {
  display: block;
  width: 10px;
}
</style>
