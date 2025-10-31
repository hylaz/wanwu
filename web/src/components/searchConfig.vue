<template>
  <el-form
    :model="formInline"
    ref="formInline"
    :inline="false"
    class="searchConfig"
  >
    <el-form-item
      class="vertical-form-item"
    >
    <template #label>
        <span v-if="!setType" class="vertical-form-title">检索方式配置</span>
    </template>
      <div
        v-for="item in searchTypeData"
        :class="['searchType-list',{ 'active': item.showContent }]"
      >
        <div
          class="searchType-title"
          @click="clickSearch(item)"
        >
          <span :class="[item.icon,'img']"></span>
          <div class="title-content">
            <div class="title-box">
              <h3 class="title-name">{{item.name}}</h3>
              <p class="title-desc">{{item.desc}}</p>
            </div>
            <span :class="item.showContent?'el-icon-arrow-up':'el-icon-arrow-down'"></span>
          </div>
        </div>
        <div
          class="searchType-content"
          v-if="item.showContent"
        >
          <div
            v-if="item.isWeight"
            class="weightType-box"
          >
            <div
              v-for="mixItem in item.mixType"
              :class="['weightType',{ 'active': mixItem.value === item.mixTypeValue }]"
              @click.stop="mixTypeClick(item,mixItem)"
            >
              <p class="weightType-name">{{mixItem.name}}</p>
              <p class="weightType-desc">{{mixItem.desc}}</p>
            </div>
          </div>
          <el-row
            v-if="item.isWeight && item.mixTypeValue === 'weight'"
            @click.stop
          >
            <el-col class="mixTypeRange-title">
              <span>语义[{{item.mixTypeRange}}]</span>
              <span>关键词[{{(1 - (item.mixTypeRange || 0)).toFixed(1)}}]</span>
            </el-col>
            <el-col>
              <el-slider
                v-model="item.mixTypeRange"
                show-stops
                :step="0.1"
                :max="1"
                @change="rangeChage($event)"
              >
              </el-slider>
            </el-col>
          </el-row>
          <el-row v-if="showRerank(item)">
            <el-col>
              <span class="content-name">Rerank模型</span>
              <el-tooltip
                class="item"
                effect="dark"
                content="重排序模型会根据候选文档与用户问题的语义匹配度，对初步检索结果进行重新排序从而进一步提升最终返回结果的相关性和准确性。"
                placement="right"
              >
                <span class="el-icon-question tips"></span>
              </el-tooltip>
            </el-col>
            <el-col>
              <el-select
                clearable
                filterable
                style="width:100%;"
                loading-text="模型加载中..."
                v-model="formInline.knowledgeMatchParams.rerankModelId"
                @visible-change="visibleChange($event)"
                @change="handleRerankChange"
                placeholder="请选择"
                :loading="rerankLoading"
              >
                <el-option
                  v-for="item in rerankOptions"
                  :key="item.modelId"
                  :label="item.displayName"
                  :value="item.modelId"
                >
                </el-option>
              </el-select>
            </el-col>
          </el-row>
          <el-row>
            <el-col>
              <span class="content-name">TopK</span>
              <el-tooltip
                class="item"
                effect="dark"
                content="用于控制检索阶段返回的最相关的文档片段的数量。这些文档片段将被送入生成模型中，用于 生成最终的回答。"
                placement="right"
              >
                <span class="el-icon-question tips"></span>
              </el-tooltip>
            </el-col>
            <el-col>
              <el-slider
                :min="1"
                :max="10"
                :step="1"
                v-model="formInline.knowledgeMatchParams.topK"
                show-input
              >
              </el-slider>
            </el-col>
          </el-row>
          <el-row v-if=showHistory(item)>
            <el-col>
              <span class="content-name">最长上下文</span>
              <el-tooltip
                class="item"
                effect="dark"
                content="保存的最长的上下文对话轮数。"
                placement="right"
              >
                <span class="el-icon-question tips"></span>
              </el-tooltip>
            </el-col>
            <el-col>
              <el-slider
                :min="0"
                :max="100"
                :step="1"
                v-model="formInline.knowledgeMatchParams.maxHistory"
                show-input
              >
              </el-slider>
            </el-col>
          </el-row>
          <el-row>
            <el-col>
              <span class="content-name">Score阈值</span>
              <el-tooltip
                class="item"
                effect="dark"
                content="检索结果的相似度阈值，低于该值的结果将被过滤。"
                placement="right"
              >
                <span class="el-icon-question tips"></span>
              </el-tooltip>
            </el-col>
            <el-col>
              <el-slider
                :min="0"
                :max="1"
                :step="0.1"
                v-model="formInline.knowledgeMatchParams.threshold"
                show-input
              >
              </el-slider>
            </el-col>
          </el-row>
        </div>
      </div>
    </el-form-item>
  </el-form>
</template>
<script>
import { getRerankList } from "@/api/modelAccess";
export default {
  props:['setType','config'],
  data() {
    return {
      debounceTimer:null,
      rerankOptions: [],
      rerankLoading: false,
      isSettingFromConfig: false, // 添加标志位，用于区分是否是从config设置的值
      formInline: {
        knowledgeMatchParams: {
          keywordPriority: 0.8, //关键词权重
          matchType: "", //vector（向量检索）、text（文本检索）、mix（混合检索：向量+文本）
          priorityMatch: 1, //权重匹配，只有在混合检索模式下，选择权重设置后，这个才设置为1
          rerankModelId: "", //rerank模型id
          threshold: 0.4, //过滤分数阈值
          semanticsPriority: 0.2, //语义权重
          topK:5, //topK 获取最高的几行
          maxHistory:0//最长上下文
        },
      },
      initialEditForm:null,
      searchTypeData: [
        {
          name: "向量检索",
          value: "vector",
          desc: "通过向量相似度找到语义相近、表达多样的文本片段，适用于理解和召回语义相关信息。",
          icon: "el-icon-menu",
          isWeight: false,
          showContent: false,
        },
        {
          name: "全文检索",
          value: "text",
          desc: "基于关键词匹配，能够高效查询包含指定词汇的文本片段，适用于精确查找",
          icon: "el-icon-document",
          isWeight: false,
          showContent: false,
        },
        {
          name: "混合检索",
          value: "mix",
          desc: "结合向量和关键词检索，融合语义理解与关键词匹配，兼顾相关性和准确性，提升检索效果。",
          icon: "el-icon-s-grid",
          isWeight: true,
          Weight: "",
          mixTypeValue: "weight",
          showContent: false,
          mixTypeRange: 0.2,
          mixType: [
            {
              name: "权重设置",
              value: "weight",
              desc: "权重设置功能用于调整不同检索方式的影响力。通过设置权重，可以控制语义相似度和关键词匹配在最终排序中的占比。",
            },
            {
              name: "Rerank模型",
              value: "rerank",
              desc: "重排序模型会根据候选文档与用户问题的语义匹配度，对初步检索结果进行重新排序从而进一步提升最终返回结果的相关性和准确性。",
            },
          ],
        },
      ],
    };
  },
  watch: {
    formInline: {
      handler(newVal) {
        // 如果是从config设置的值，不触发sendConfigInfo
        if (this.isSettingFromConfig) {
          return;
        }
        
        if (this.debounceTimer) {
          clearTimeout(this.debounceTimer);
        }
        this.debounceTimer = setTimeout(() => {
          const props = ['knowledgeMatchParams'];
          const changed = props.some(prop => {
            return JSON.stringify(newVal[prop]) !== JSON.stringify(
                (this.initialEditForm || {})[prop]
              );
            });
          if (changed) {
            if(!this.setType){
              delete this.formInline.knowledgeMatchParams.maxHistory;
            }
            this.$emit('sendConfigInfo', this.formInline);
          }
        }, 200);
      },
      deep: true,
      immediate: false
    },
    config:{
      handler(newVal) {
        if(newVal && Object.keys(newVal).length > 0){
          this.isSettingFromConfig = true; // 设置标志位
          const formData = JSON.parse(JSON.stringify(newVal))
          this.formInline.knowledgeMatchParams = formData;
          const { matchType,priorityMatch } = this.formInline.knowledgeMatchParams;
          if(matchType !== ''){
              this.searchTypeData = this.searchTypeData.map((item) => ({
              ...item,
              showContent: item.value === matchType ? true : false,
            }));
            if(matchType === 'mix'){
              this.searchTypeData[2]['mixTypeValue'] = priorityMatch === 1 ? 'weight' : 'rerank';
            }
          }

          // 使用nextTick确保DOM更新完成后再重置标志位
          this.$nextTick(() => {
            this.isSettingFromConfig = false;
          });
        }
      },
      deep: true,
      immediate: true
    }
  },
  mounted() {
    this.$nextTick(() => {
      this.initialEditForm = JSON.parse(JSON.stringify(this.formInline));
    });
  },
  created() {
    // 预加载数据，避免首次打开下拉框时的延迟
    this.getRerankData();
  },
  methods: {
    rangeChage(val){
      this.formInline.knowledgeMatchParams.keywordPriority = Number((1 - (val || 0)).toFixed(1));
      this.formInline.knowledgeMatchParams.semanticsPriority = val;
    },
    mixTypeClick(item, n) {
      item.mixTypeValue = n.value;
      const { knowledgeMatchParams } = this.formInline;
      knowledgeMatchParams.priorityMatch = n.value === "weight" ? 1 : 0;
      // if(n.value === 'weight'){
      //   knowledgeMatchParams.rerankModelId = '';
      // }
    },
    showRerank(n) {
      return (
        n.value === "vector" ||
        n.value === "text" ||
        (n.value === "mix" && n.mixTypeValue === "rerank")
      );
    },
    showHistory(n){
      return (
       (this.setType === 'rag'||this.setType === 'agent') &&
        (n.value === "vector" ||
         n.value === "text" ||
         (n.value === "mix") //&& n.mixTypeValue === "rerank"
        )
      )
    },
    clickSearch(n) {
      this.formInline.knowledgeMatchParams.matchType = n.value;
      this.searchTypeData = this.searchTypeData.map((item) => ({
        ...item,
        showContent: item.value === n.value ? !item.showContent : false,
      }));
      this.formInline.knowledgeMatchParams.priorityMatch = n.value !== 'mix' ? 0 : 1;
      this.clear();
    },
    clear() {
      this.formInline.knowledgeMatchParams.rerankModelId = "";
      this.formInline.knowledgeMatchParams.keywordPriority = 0.8;
      this.formInline.knowledgeMatchParams.semanticsPriority = 0.2;
      this.formInline.knowledgeMatchParams.threshold = 0.4;
      this.formInline.knowledgeMatchParams.topK = 5;
    },
    getRerankData() {
      this.rerankLoading = true;
      getRerankList().then((res) => {
        if (res.code === 0) {
          this.rerankOptions = res.data.list || [];
        }
      }).finally(() => {
        this.rerankLoading = false;
      });
    },
    visibleChange(val) {
      if (val && this.rerankOptions.length === 0) {
        this.getRerankData();
      }
    },
    handleRerankChange(value) {
      // 直接触发事件，避免防抖延迟
      if(!this.setType){
        const formData = JSON.parse(JSON.stringify(this.formInline));
        delete formData.knowledgeMatchParams.maxHistory;
        this.$emit('sendConfigInfo', formData);
      } else {
        this.$emit('sendConfigInfo', this.formInline);
      }
    },
  },
};
</script>
<style lang="scss" scoped>
/deep/ {
  .vertical-form-item {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    .vertical-form-title{
      color:#000;
      font-size:14px;
    }
  }
  .vertical-form-item .el-form-item__label {
    line-height: unset;
    font-size: 14px;
    font-weight: bold;
  }
  .el-form-item__content {
    width: 100%;
  }
  .el-input-number--small{
    line-height: 28px!important;
  }
}
.active {
  border: 1px solid $color !important;
}
.searchConfig {
  .searchType-list:hover {
    border: 1px solid $color;
  }
  .searchType-list {
    border: 1px solid #c0c4cc;
    border-radius: 4px;
    margin: 20px 0;
    padding: 0 10px;
    cursor: pointer;
    .searchType-title {
      display: flex;
      align-items: center;
      .img {
        font-size: 30px;
        text-align: center;
        line-height: 50px;
        color: $color;
        background-color: #fff;
        width: 50px;
        height: 50px;
        border-radius: 8px;
        border: 1px solid #e9e9eb;
        box-shadow: 4px 2px 4px #f1f1f1;
      }
      .title-content {
        flex: 1;
        display: flex;
        margin-left: 10px;
        justify-content: space-between;
        align-items: center;
        .title-name {
          font-size: 16px;
          font-weight: bold;
          line-height: 1;
          padding-top: 10px;
        }
        .title-desc {
          color: #888;
        }
      }
    }
    .searchType-content {
      padding: 20px;
      .tips {
        color: #888;
        margin-left: 5px;
      }
      .content-name {
        font-weight: bold;
      }
      .weightType-box {
        display: flex;
        gap: 20px;
        .weightType {
          border: 1px solid #c0c4cc;
          border-radius: 4px;
          .weightType-name {
            text-align: center;
            font-weight: bold;
            line-height: 2;
            font-size: 16px;
            padding-top: 5px;
          }
          .weightType-desc {
            text-align: center;
            line-height: 1.5;
            padding: 10px;
            color: #888;
          }
        }
      }
      .mixTypeRange-title {
        display: flex;
        align-items: center;
        justify-content: space-between;
        font-weight: bold;
        margin-top: 20px;
        line-height: 1;
      }
    }
  }
}
</style>