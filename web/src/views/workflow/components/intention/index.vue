<template>
  <div class="model_box">
    <!--无默认模型，隐藏模式选择，默认展示完整模式-->
    <!--<div class="node-params">
      <div class="params-type">
        <i class="el-icon-caret-bottom"></i>
        <span class="params-type-span">模式</span>
      </div>
      <el-radio-group  v-model="nodeData.data.settings.lite_mode" @input="modelChange">
        <el-radio :label="true">极速模式
          <el-popover
              title=""
              width="200"
              trigger="hover"
              content="极速模式：适用于意图明确且简单，对响应速度要求高的场景。">
              <i class="el-icon-question" slot="reference"></i>
            </el-popover>
        </el-radio>
        <el-radio :label="false">完整模式
          <el-popover
              title=""
              width="200"
              trigger="hover"
              content="完整模式：支持配置意图例句及大模型参数，耗时相对较长。">
              <i class="el-icon-question" slot="reference"></i>
            </el-popover>
        </el-radio>
        </el-radio-group>
    </div>-->
    <!--输入-->
    <div class="node-params">
      <div class="params-type">
        <i class="el-icon-caret-bottom"></i>
        <span class="params-type-span">输入</span>
      </div>
      <!--form-->
      <div class="params-form">
        <div class="form-item">
          <div class="item">参数名</div>
          <div class="item">类型</div>
          <div class="item last-item">值</div>
        </div>
        <div
          class="form-item"
          v-for="(n, i) in nodeData.data.inputs"
          :key="`${i}ipt`"
        >
          <el-input class="item" size="mini" v-text="n.name"></el-input>
          <el-select
            class="item"
            size="mini"
            popper-class="workflow-select"
            v-model="n.value.type"
            @change="(val) => selectTypeChange(val, i, n)"
          >
            <el-option value="generated" label="String"></el-option>
            <el-option value="ref" label="引用"></el-option>
          </el-select>
          <!--非引用-->
          <el-input
            class="item last-item"
            size="mini"
            style="width: 190px"
            :placeholder="n.desc"
            v-if="n.value.type == 'generated'"
            v-model="n.value.content"
          ></el-input>

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
      <p v-if="nodeData.validate && JSON.parse(nodeData.validate).inputValidate === false" class="workflow-errormsg">{{JSON.parse(nodeData.validate).message}}</p>
    </div>

    <!--意图配置-->
    <div class="node-params">
      <div class="params-type">
        <i class="el-icon-caret-bottom"></i>
        <span class="params-type-span">意图配置</span>
      </div>
      <div class="params-form">
        <div v-for="(item,index) in nodeData.data.settings.intentions" class="iten-item" :key="index">
          <div>
              <div class="form-item">
                  <el-input class="item-label" size="mini" v-text="'意图名称'"></el-input>
                  <el-input
                          class="item last-item"
                          size="mini"
                          placeholder="请填写意图名称"
                          :disabled="index === nodeData.data.settings.intentions.length-1"
                          v-model="item.name"
                      ></el-input>
              </div>
              <div v-if="index !== nodeData.data.settings.intentions.length-1" class="form-item">
                  <el-input class="item-label" size="mini" v-text="'意图描述'"></el-input>
                  <el-input
                          class="item last-item"
                          size="mini"
                          type="textarea"
                          placeholder="请描述意图的含义、使用场景，或提供例句，便于大模型更好识别该意图"
                          v-model="item.desc"
                      ></el-input>
              </div>
              <div v-if="nodeData.data.settings.lite_mode===false && index !== nodeData.data.settings.intentions.length-1 " class="form-item item-vertical">
                <div class="iten-ext">
                    <el-input class="item-label" size="mini" v-text="'意图例句'"></el-input>
                      <div>
                        <span>{{item.sentences.length}}/10</span>
                        <i
                          v-if="item.sentences.length<10"
                          class="el-icon-circle-plus-outline"
                          @click.stop="addSentence(index)"
                          ></i>
                      </div>
                  </div>
                  <div class="iten-sentence" v-for="(iitem, iindex) in item.sentences" :key="`${iindex}sentence`">
                    <el-input
                        class="item last-item"
                        size="mini"
                        placeholder="请输入意图例句"
                        v-model="item.sentences[iindex]"
                    ></el-input>
                    <i
                      class="el-icon-remove-outline"
                      @click.stop="deleteSentence(index,iindex)"
                      ></i>
                  </div>
              </div>
            </div>
            <i
              v-if="index>0 && index!==nodeData.data.settings.intentions.length-1"
              class="el-icon-remove-outline condition-item-del-icon"
              @click.stop="preDelCondition(index)"
              ></i>
        </div>
        <p v-if="showError" class="workflow-errormsg">意图名称存在重复，请修改名称</p>
          <el-button
                    v-if=" nodeData.data.settings.intentions.length <10"
                    class="add-switch-bt"
                    size="mini"
                    icon="el-icon-plus"
                    @click="preAdd"
            >添加意图</el-button>
      </div>
    </div>

    <div class="node-params" v-if="!nodeData.data.settings.lite_mode">
      <div class="params-type">
        <i class="el-icon-caret-bottom"></i>
        <span class="params-type-span">高级配置</span>
        <!-- <i class="el-icon-plus add-icon" @click="preAddOutputParams"></i>-->
      </div>
      <el-form ref="form" label-position="left" label-width="80px">
        <el-form-item label="选择模型">
          <el-select
           :disabled="nodeData.data.settings.lite_mode"
            v-model="nodeData.data.settings.model.model_name"
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
        <el-form-item label="多样性">
          <el-row class="slider-box">
            <el-col :span="10" class="slider">
              <el-slider
                v-model="nodeData.data.settings.model.top_p"
                :min="0"
                :max="1"
                :step="0.01"
              ></el-slider>
            </el-col>
            <el-col :span="10" class="input">
              <el-input-number
                size="mini"
                controls-position="right"
                v-model="nodeData.data.settings.model.top_p"
                :min="0.01"
                :max="1"
                :step="0.01"
              ></el-input-number>
            </el-col>
          </el-row>
        </el-form-item>
       <el-form-item label="附加提示词">
           <el-input
            class="item"
            type="textarea"
            placeholder="请输入附加提示词，提高意图识别准确率"
            v-model="nodeData.data.settings.model.system"
          ></el-input>
        </el-form-item>
      </el-form>
    </div>

    <!--输出-->
    <div class="node-params">
      <div class="params-type">
        <i class="el-icon-caret-bottom"></i>
        <span class="params-type-span">输出</span>
        <!-- <i class="el-icon-plus add-icon" @click="preAddOutputParams"></i>-->
      </div>
      <div class="params-form">
        <div class="form-item">
          <div class="item">参数名</div>
          <div class="item">类型</div>
        </div>
        <div
          class="form-item"
          v-for="(n, i) in nodeData.data.outputs"
          :key="`${i}opts`"
        >
          <el-input
            class="item"
            size="mini"
            v-model="n.name"
            disabled
          ></el-input>
          <el-select class="item" size="mini" v-model="n.type" disabled>
            <el-option value="string" label="String"></el-option>
          </el-select>
        </div>
      </div>
      <p v-if="nodeData.validate && JSON.parse(nodeData.validate).outputValidate === false" class="workflow-errormsg">{{JSON.parse(nodeData.validate).message}}</p>
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";
import nodeMethod from "@/views/workflow/mixins/nodeMethod";
import {getModels} from "@/api/workflow";

const models =  [
            {
                "modelId": "yuanjing-70b-chat",
                "modelName": "Yuanjing-70B-Chat",
            },
            {
                "modelId": "deepseek-r1",
                "modelName": "DeepSeek-R1-671B",
            },
            {
                "modelId": "deepseek-v3",
                "modelName": "DeepSeek-V3-0324",
            },
            {
                "modelId": "deepseek-r1-distill-qwen-1.5b",
                "modelName": "DeepSeek-R1-Distill-Qwen-1.5B",
            },
            {
                "modelId": "deepseek-r1-distill-qwen-7b",
                "modelName": "DeepSeek-R1-Distill-Qwen-7B",
            },
            {
                "modelId": "deepseek-r1-distill-qwen-14b",
                "modelName": "DeepSeek-R1-Distill-Qwen-14B",
            },
            {
                "modelId": "deepseek-r1-distill-qwen-32b",
                "modelName": "DeepSeek-R1-Distill-Qwen-32B",
            },
            {
                "modelId": "deepseek-r1-distill-llama-8b",
                "modelName": "DeepSeek-R1-Distill-Llama-8B",
            },
            {
                "modelId": "deepseek-r1-distill-llama-70b",
                "modelName": "DeepSeek-R1-Distill-Llama-70B",
            }
        ]

export default {
  components: {},
  props: ["graph", "node"],
  data() {
    return {
      placeholder:"",
      settings: {},
      nodeData: {},
      preNode: {},
      preNodeOutputs: [],
      preAllNode: [],
      showError: false,
      models: [],
    };
  },
  computed: {
    ...mapState({
      nodeIdMap: (state) => state.workflow.nodeIdMap,
    }),
  },
  watch:{
     "nodeData.data.settings.intentions": {
        handler: function (newVal, oldVal) {
          let intentionsMap = []
          let showError = false
          newVal.forEach((item)=>{
            if(intentionsMap.indexOf(item.name) === -1){
              intentionsMap.push(item.name)
            }else{
              showError = true
            }
          })
          this.showError = showError
        },
        deep: true,
      },
  },
  mixins: [nodeMethod],
  created() {
    //this.nodeData = this.node.data;
    this.getModels();
  },
  methods: {
    getModels() {
      getModels().then((res) => {
        const {list} = res.data || {}
        const models = list ? list.map((item) => ({...item, modelName: item.displayName || item.model})) : []
        this.models = models
      });
    },
    selectTypeChange(val, i, n) {
      switch (val) {
        case "ref":
          this.preAllNode = [];
          this.getPreNode(this.nodeData.id);
          break;
        case "generated":
          let newItem = {
            desc: n.desc||"",
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
          this.$nextTick(() => {
            this.$set(this.nodeData.data.inputs, i, newItem);
          });

          break;
      }
    },
    preAdd(){
      let lengths = this.node.data.data.settings.intentions.length
        this.node.data.data.settings.intentions.splice(lengths-1,0,{
             "name": `意图${lengths}`,
            "desc": "",
            "global_selected": true, // 预留，前端不使用
            "target_node_id": "",
            "sentences": [],
            "params": []
        })

        setTimeout(() => {
                    let port = this.node.port.ports;
                    if (port[0].group === "left") {
                        port.shift();
                    }
                    let lastPort = port[port.length - 1];

                    const edges = [];
                    this.graph.getEdges().forEach((edge) => {
                        const sourcePortId = edge.getSourcePortId();
                        const targetPortId = edge.getTargetPortId();
                        if (sourcePortId === lastPort.id || targetPortId === lastPort.id) {
                            edges.push(edge);
                        }
                    });

                    this.node.addPort({
                        group: "absolute",
                        id: this.node.id + port.length + "-right",
                        args: {
                            x: "100%",
                            y: lastPort.args.y + 40,
                        },
                    });
                    if (edges.length > 0) {
                        let newPort = this.node.port.ports;
                        edges[0].setSource({
                            cell: this.node.id,
                            port: newPort[newPort.length - 1].id,
                        });
                    }
                    let ee = document.querySelector(
                        `.x6-node[data-cell-id="${this.node.id}"] foreignObject body .dag-node`
                    );
                    this.node.resize(ee.offsetWidth, ee.offsetHeight);
                    const portsdom = document.querySelectorAll(".x6-port-body");
                    for (let i = 0, len = portsdom.length; i < len; i = i + 1) {
                        portsdom[i].style.visibility = "visible";
                    }
                }, 10);
    },
    preDelCondition(index){
         this.node.data.data.settings.intentions.splice(index,1)
           let port = JSON.parse(JSON.stringify(this.node.port.ports));

                if (port[0].group === "left") {
                    port.shift();
                }
                let lastPort = port[port.length - 1];
                // port.shift();
                // port.pop();
                let edgesMap = {};
                let edgesArr = [];
                const edges = this.graph.getEdges().filter((edge) => {
                    const target = edge.getTarget();
                    const source = edge.getSource();
                    const sourcePortId = edge.getSourcePortId();
                    const targetPortId = edge.getTargetPortId();
                    edgesMap[source.port] = edge;
                    if (sourcePortId === lastPort.id || targetPortId === lastPort.id) {
                        edgesArr.push(edge);
                    }
                    return target.port === port[index].id || source.port === port[index].id;
                });
                setTimeout(() => {
                    setTimeout(() => {
                        setTimeout(() => {
                            port.forEach((iport, poindex) => {
                                if (
                                    poindex > index &&
                                    edgesMap[iport.id] &&
                                    poindex != port.length - 1
                                ) {
                                    edgesMap[iport.id].setSource({
                                        cell: this.node.id,
                                        port: port[poindex - 1].id,
                                    });
                                }
                            });
                            edges[0].delSource = 'condition'
                            this.graph.removeEdge(edges[0]);
                        }, 10);
                        if (edgesMap[port[port.length - 1].id]) {
                            edgesMap[port[port.length - 1].id].setSource({
                                cell: this.node.id,
                                port: port[port.length - 2].id,
                            });
                        }
                        this.node.removePort(port[port.length - 1].id);
                        // this.node.scale(1,1)
                        // refreshNode(cnode)
                    }, 10);
                }, 50);
    },
    modelChange(val){
      if(val){
        this.node.data.data.settings.model.model_name = "yuanjing-70b-chat"
         this.node.data.data.outputs = [ {
              "name": "classification",
              "type": "string",
              "desc": "",
              "object_schema": null,
              "list_schema": null,
              "required": false,
              "value": {
                "type": "generated",
                "content": null
              }
            },
            {
              "name": "classificationID",
              "type": "integer",
              "desc": "",
              "object_schema": null,
              "list_schema": null,
              "required": false,
              "value": {
                "type": "generated",
                "content": null
              },
            }]
        this.node.data.data.settings.intentions.forEach((item)=>{
          item.sentences = []
        })
      }else{
         this.node.data.data.outputs = [{
            "name": "thought",
            "type": "string",
            "desc": "",
            "object_schema": null,
            "list_schema": null,
            "required": false,
            "value": {
              "type": "generated",
              "content": null
            }
          },
          {
            "name": "classification",
            "type": "string",
            "desc": "",
            "object_schema": null,
            "list_schema": null,
            "required": false,
            "value": {
              "type": "generated",
              "content": null
            }
          },
          {
            "name": "classificationID",
            "type": "integer",
            "desc": "",
            "object_schema": null,
            "list_schema": null,
            "required": false,
            "value": {
              "type": "generated",
              "content": null
            }
          }]
      }
    },
    addSentence(index){
      this.node.data.data.settings.intentions[index].sentences.push("")
    },
    deleteSentence(index, iindex){
      this.node.data.data.settings.intentions[index].sentences.splice(iindex,1)
    },
  },
};
</script>

<style lang="scss" scoped>
   .add-switch-bt {
    width: 100%;
  }

  .iten-item{
    position:relative;
    margin-bottom:15px;
    .item-label{
        flex: 1;
        max-width: 80px;
        line-height: 30px;
    }
    .form-item{
        padding-right:35px;
    }
    .condition-item-del-icon{
        position:absolute;
        right:10px;
        top: 10px;
    }
  }
  .iten-ext{
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    i{
      cursor:pointer;
    }
    span{
      color: $disabled_color;
      font-size: 13px;
      padding-right:5px;
    }
  }
  .item-vertical{
    flex-direction: column;
  }
  .iten-sentence{
    display:flex;
    align-items: center;
    i{
      margin-left:5px;
      cursor:pointer;
    }
  }
</style>
