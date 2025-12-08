<template>
  <div>
    <el-dialog
      title="召回参数配置"
      :visible.sync="dialogVisible"
      width="50%"
      :before-close="handleClose"
    >
      <span>
        <el-form
          :model="ruleForm"
          ref="ruleForm"
          label-width="100px"
          class="demo-ruleForm"
        >
          <el-form-item
            :label="item.label"
            :prop="item.props"
            v-for="(item, index) in konwledgeSet"
            :key="index"
          >
            <el-row>
              <el-col :span="1">
                <el-tooltip
                  class="item"
                  effect="light"
                  :content="item.desc"
                  placement="bottom"
                >
                  <span class="el-icon-question question"></span>
                </el-tooltip>
              </el-col>
              <el-col :span="2">
                <el-switch v-model="ruleForm[item.btnProps]"></el-switch>
              </el-col>
              <el-col :span="20">
                <el-slider
                  v-model="ruleForm[item.props]"
                  show-input
                  :min="item.min"
                  :max="item.max"
                  :step="item.step"
                ></el-slider>
              </el-col>
            </el-row>
          </el-form-item>
        </el-form>
      </span>
      <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="submit">确 定</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>
export default {
  props: {
    knowledgeConfig: {
      type: Object,
      default: null,
    },
  },
  data() {
    return {
      dialogVisible: false,
      ruleForm: {
        maxHistory: 0,
        threshold: 0.4,
        topK: 5,
        maxHistoryEnable: true,
        thresholdEnable: true,
        topKEnable: true,
      },
      konwledgeSet: [
        {
          label: '最长上下文',
          desc: '保存的最长的上下文对话轮数',
          props: 'maxHistory',
          btnProps: 'maxHistoryEnable',
          min: 0,
          max: 100,
          step: 1,
        },
        {
          label: '过滤阀值',
          desc: '检索结果匹配度的最小值，小于阈值的结果会被过滤掉',
          props: 'threshold',
          btnProps: 'thresholdEnable',
          precision: 1,
          min: 0,
          max: 1,
          step: 0.1,
        },
        {
          label: '知识条数',
          desc: '检索召回的知识片段数量的最大值，当检索到的知识数量大于知识条数，也只返回最大知识条数',
          props: 'topK',
          btnProps: 'topKEnable',
          min: 1,
          max: 20,
          step: 1,
        },
      ],
    };
  },
  methods: {
    showDialog() {
      this.dialogVisible = true;
      if (this.knowledgeConfig !== null) {
        const data = JSON.parse(JSON.stringify(this.knowledgeConfig));
        this.ruleForm = data;
      }
    },
    handleClose() {
      this.dialogVisible = false;
    },
    submit() {
      this.dialogVisible = false;
      this.$emit('setKnowledgeSet', this.ruleForm);
    },
  },
};
</script>
<style lang="scss" scoped>
/deep/ {
  .el-input-number--small {
    line-height: 28px !important;
  }
}
.question {
  cursor: pointer;
  color: #ccc;
}
</style>
