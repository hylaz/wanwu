<template>
  <div class="api_box">
    <!--知识库节点添加模型选择-->
    <div class="model-box" style="margin-top: 20px">
      <p class="params-type-span">模型</p>
      <el-form ref="form" label-position="left" label-width="90px">
        <el-form-item label="Rerank模型">
          <el-select
              v-model="settings.model"
              placeholder="选择模型"
          >
            <el-option
                v-for="item in models"
                :label="item.modelName"
                :value="item.modelId"
                :key="item.modelId"
            ></el-option>
          </el-select>
        </el-form-item>
      </el-form>
    </div>

    <!--输入-->
    <div class="node-params">
      <div class="params-type">
        <i class="el-icon-caret-bottom"></i>
        <span class="params-type-span">输入</span>
      </div>
      <!--form-->
      <div v-for="(n, i) in nodeData.data.inputs" :key="`${i}ipt`">
        <!--query-->
        <div class="params-form" v-if="n.name ==='query'">
          <div class="form-item">
            <div class="item">参数名</div>
            <div class="item">类型</div>
            <div class="item last-item">值</div>
          </div>
          <div class="form-item" >
            <el-input class="item" size="mini" v-text="n.name"></el-input>
            <el-select
                    class="item"
                    size="mini"
                    v-model="n.value.type"
                    @change="(val) => selectTypeChange(val, i,n)"
            >
              <el-option value="generated" label="String"></el-option>
              <el-option value="ref" label="引用"></el-option>
            </el-select>
            <!--非引用-->
            <el-input
                    class="item last-item"
                    size="mini"
                    v-if="n.value.type === 'generated'"
                    v-model="n.value.content"
            ></el-input>
            <!--引用-->
            <el-popover placement="bottom" width="260" trigger="click">
              <div>
                <div
                        class="popover-select-item"
                        v-for="(m, j) in preAllNode"
                        :key="`${j}pn`"
                >
                  <p class="node-name">
                    <i class="el-icon-caret-bottom"></i>{{ m.name }}
                  </p>
                  <div
                          class="node-content"
                          v-for="(p, l) in m.data.data.outputs"
                          :key="`${l}opt`"
                          @click="refValueClick(n, i, m, p, l)"
                  >
                    <span class="name">{{ p.name }}</span>
                    <span class="type">{{ p.type }}</span>
                  </div>
                </div>
              </div>
              <div
                      slot="reference"
                      v-show="n.value.type === 'ref'"
                      class="item last-item popover-select"
              >
                {{ n.newRefContent }}
              </div>
            </el-popover>
          </div>
        </div>
      </div>
      <p v-if="nodeData.validate && JSON.parse(nodeData.validate).inputValidate === false" class="workflow-errormsg">{{JSON.parse(nodeData.validate).message}}</p>
    </div>

    <!--知识库参数-->
    <div class="node-params">
      <div class="params-type">
        <i class="el-icon-caret-bottom"></i>
        <span class="params-type-span">知识库参数</span>
      </div>
      <div class="rl" v-loading="knowledgeLoading">
        <treeselect
                class="treeselect"
                :clearable="true"
                :options="knowledgeData"
                no-options-text="暂无数据"
                :multiple="true"
                :normalizer="normalizeOptions"
                v-model="settings.knowledgeBase"
                placeholder="请选择知识库"
                ref="optionRef"
                @select="knowledgeChange"
                @deselect="knowledgeChange"
                flat
                :disableBranchNodes="true"
        ></treeselect>
      </div>
      <!--<el-select class="knowledgeBase-select" placeholder="请选择知识库" size="mini" multiple clearable></el-select>-->
      <el-form label-width="80px">
        <div v-for="(n, i) in nodeData.data.inputs" :key="`${i}ipt`">
          <!--过滤阈值 threshold-->
          <div class="params-form" v-if="n.name ==='threshold'">
            <el-form-item label="过滤阈值">
              <el-row class="slider-box">
                <el-col :span="11" class="slider">
                  <el-slider v-model="n.value.content" :min="0" :max="1" :step="0.01"></el-slider>
                </el-col>
                <el-col :span="8" class="input">
                  <el-input-number size="mini" controls-position="right" v-model="n.value.content" :min="0" :max="1" :step="0.01"></el-input-number>
                </el-col>
              </el-row>
            </el-form-item>
          </div>

        <!--知识条数 top_k-->
        <div class="params-form" v-if="n.name ==='top_k'">
          <el-form-item label="知识条数">
            <el-row class="slider-box">
              <el-col :span="11" class="slider">
                <el-slider v-model="n.value.content" :min="0" :max="20" :step="1"></el-slider>
              </el-col>
              <el-col :span="8" class="input">
                <el-input-number size="mini" controls-position="right" v-model="n.value.content" :min="0" :max="20" :step="1"></el-input-number>
              </el-col>
            </el-row>
          </el-form-item>
        </div>
      </div>
      </el-form>
    </div>

    <!--输出-->
    <div class="node-params">
      <div class="params-type">
        <i class="el-icon-caret-bottom"></i>
        <span class="params-type-span">输出</span>
      </div>
      <div class="params-form">
        <div class="form-item">
          <div class="item">参数名</div>
          <div class="item">类型</div>
        </div>
        <div class="form-item" v-for="(n, i) in nodeData.data.outputs" :key="`${i}opts`">
          <el-input class="item" size="mini" v-model="n.name" disabled></el-input>
          <el-select class="item" size="mini" v-model="n.type" disabled>
            <el-option value="string" label="String"></el-option>
            <el-option value="object" label="引用"></el-option>
          </el-select>
        </div>
      </div>
      <p v-if="nodeData.validate && JSON.parse(nodeData.validate).outputValidate === false" class="workflow-errormsg">{{JSON.parse(nodeData.validate).message}}</p>
    </div>

  </div>
</template>

<script>
import { getRerankModels } from "@/api/workflow";
import { mapState, mapActions } from "vuex";

import nodeMethod from "@/views/workflow/mixins/nodeMethod";
import { getKnowledgeList } from "@/api/knowledge";
import Treeselect from "@riophae/vue-treeselect";
import "@riophae/vue-treeselect/dist/vue-treeselect.css";

export default {
  props: ["graph", "node"],
  data() {
    return {
      editVisible: false,
      nodeData: {},
      preNodeOutputs: [],
      settings: {
          model:'',
          "headers": {},
          "knowledgeBase": [],
          "userId": JSON.parse(
              localStorage.getItem("access_cert")
          ).user.userInfo.uid,
          "content_type": "application/json"
      },
      preAllNode: [],
      models:[],
      //知识库
      //knowledgeData: [],
      knowledgeNameObj: {},
      knowledgeLoading:false
    };
  },
  mixins: [nodeMethod],
  computed: {
    ...mapState({
      nodeIdMap: (state) => state.workflow.nodeIdMap,
      knowledgeData:(state) => state.workflow.knowledgeData,
    }),
  },
  components:{ Treeselect },
  created() {
    this.nodeData = this.node.data;
    this.settings = this.nodeData.data.settings
    this.getModels()

    console.log('获取knowledgeData。。。',this.knowledgeData)
    if(!this.knowledgeData.length){
        this.getClassfyDoc();  //todo 每次展开节点都会触发右侧DagNode组件刷新,导致知识库闪动
    }

    //this.parseNodeData(this.nodeData);

  },
  watch: {
      "node.data.id": {
          handler: function (newVal, oldVal) {
              console.log("watch:", newVal);
              this.$forceUpdate()
          },
          deep: true,
        },
  },

  methods: {
    ...mapActions("workflow", [
        "setKnowledgeList"
    ]),
    getModels() {
      getRerankModels().then((res) => {
        const {list} = res.data || {}
        const models = list ? list.map((item) => ({...item, modelName: item.displayName || item.model})) : []
        this.models = models
      });
    },
    preEdit() {
      this.editVisible = true;
    },
    closeEdit() {
      this.editVisible = false;
    },
    inputsChange(newInputs) {
      this.nodeData.inputs = newInputs.filter((n) => {
        return n.name;
      });
    },
      selectTypeChange(val, i) {
          switch (val) {
              case "ref":
                  break;
              case "generated":
                  let newItem = {
                      desc: "",
                      extra: {
                          location: "body",
                      },
                      list_schema: null,
                      name: this.nodeData.data.inputs[i].name,
                      object_schema: null,
                      required: false,
                      type: "string",
                      value: {
                          content: "",
                          type: "generated",
                      },
                      newRefContent: "",
                  };
                  this.$set(this.nodeData.data.inputs, i, newItem);
                  break;
          }
      },
      //知识库
      normalizeOptions(node) {
          if (node.children == null || node.children == "null") {
              delete node.children;
          }
          return {
              id: node.name,
              label: node.name,
              children: node.children,
          };
      },
      async getClassfyDoc() {
          //获取文档知识分类
          this.knowledgeLoading = true
          const res = await getKnowledgeList({});
          if (res.code === 0) {
              this.knowledgeLoading = false
              //this.knowledgeData = res.data;
              const knowledgeList = res.data ? (res.data.knowledgeList || []) : []
              this.setKnowledgeList(knowledgeList)
              //保存id和name映射关系
              this.saveKnowledgeName(knowledgeList);
          } else {
              this.$message.error(res.message);
          }
      },
      saveKnowledgeName(data) {
          data.forEach((n) => {
              if (n.children) {
                  this.saveKnowledgeName(n.children);
              } else {
                  this.knowledgeNameObj[n.knowledgeId] = n.name;
              }
          });
      },
      knowledgeChange() {

      },
  },
};
</script>

<style lang="scss" scoped>
  .node-params {
    padding: 10px;
    background-color: #f9f9fb;
    color: #151b26;
    border-radius: 10px;
    margin-top: 20px;
    .params-type {
      position: relative;
      .add-icon {
        position: absolute;
        right: 20px;
        top: 10px;
        font-size: 16px;
      }
      .params-type-span {
        margin-left: 10px;
      }
    }
    .params-content {
      margin-top: 8px;
      min-height: 40px;
      .params-content-item {
        span {
          font-size: 12px;
        }
        span:nth-child(1) {
          color: #876300;
        }
        span:nth-child(2) {
          margin-left: 6px;
          height: 20px;
          padding: 0 5px;
          white-space: nowrap;
          border-radius: 4px;
          background-color: #e8e9eb;
          color: #5c5f66;
        }
        span:nth-child(3) {
          margin-left: 6px;
          height: 20px;
          padding: 0 5px;
          white-space: nowrap;
        }
      }
    }
    .params-form {
      margin-top: 10px;
      /deep/ .el-form-item{
        margin-bottom: 0!important;
      }
      /deep/.el-form-item__label{
        text-align: left!important;
      }
      .form-item {
        display: flex;
        gap: 5px;
        margin: 5px 0;
        .item {
          flex: 1;
          font-size: 12px;
          display: flex;
          align-items: center;
        }
        .last-item {
          flex: 2;
        }
        .del-icon {
          line-height: 30px;
        }
      }
    }
  }
  .node-params{
    margin-top: 20px;
    .params-type-span{
      margin-bottom: 10px;
      .required{
        color: #e60001;
        margin-left: 5px;
      }
    }
    .slider-box{
      .slider{

      }
      .input{
        margin-left: 10px;
      }
    }
  }
.treeselect{
  width: 100%;
  margin-top: 10px;
  z-index: 1100;
}
  .treeselect/deep/.vue-treeselect__control {
    background-color: transparent !important;
  }
</style>
