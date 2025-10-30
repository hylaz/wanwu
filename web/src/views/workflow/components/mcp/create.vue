<template>
  <div v-if="dialogVisible">
    <el-dialog
      title="添加MCP服务"
      :visible="dialogVisible"
      width="680px"
      append-to-body
      :close-on-click-modal="false"
      :before-close="handleClose"
    >
      <div v-loading="loading">
        <el-button-group class="mcpBtn">
          <el-button
            type="primary"
            @click="handleClick('自定义MCP')"
            :class="{ active: active === '自定义MCP' }"
            >自定义MCP</el-button
          >
        </el-button-group>
        <el-collapse
          accordion
          v-if="serverOptions.length > 0"
          class="mcpCollapse"
          v-model="activeCollapseName"
          @change="changeItem"
        >
          <el-collapse-item
            v-for="(item, index) in serverOptions"
            :key="index"
            :name="item.name"
          >
            <template slot="title">
              <div style="padding-left: 10px" @click.stop class="title-box">
                <span
                  class="title-text"
                  @click="handleTitleClick(item.name)"
                  style="padding-left: 10px"
                  >{{ item.name }}</span
                >
              </div>
            </template>

            <div>
              <ul class="detail-mcpList" v-if="item.tools">
                <li v-for="(i, n) in item.tools" :key="n">
                  <el-checkbox
                    :disabled="i.disabled"
                    @change="handleRadioClick(index, n, i.name)"
                    :value="i.checked"
                  >
                  </el-checkbox>
                  <span class="title-text" style="padding-left: 10px">
                    {{i.name }}
                  </span>
                </li>
              </ul>
            </div>
          </el-collapse-item>
        </el-collapse>
        <div class="no-list" v-if="serverOptions.length === 0">
          <div>
            您还没有可选用的MCP<br />
            <span>请前往“<strong @click="handleOpenMcp">MCP管理</strong>”进行添加</span>
          </div>
        </div>
      </div>

      <span slot="footer" class="dialog-footer">
        <el-button size="small" @click="reset">取消</el-button>
        <el-button
          size="small"
          :disabled="!mcp_node"
          type="primary"
          @click="doSubmit"
        >
          确定
        </el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { getMcpToolList, getList } from "@/api/workflow";
export default {
  data() {
    return {
      activeCollapseName: "",
      serverOptions: [],
      active: "自定义MCP",
      loading: false,
      form: {
        mcpName: "",
        mcpDesc: "",
        mcpServerUrl: "",
      },
      rules: {
        mcpName: [
          { required: true, message: "请选择服务名称", trigger: "change" },
        ],
        mcpDesc: [
          { required: true, message: "请输入功能描述", trigger: "blur" },
        ],
        mcpServerUrl: [
          { required: true, message: "请输入ServerURL", trigger: "blur" },
        ],
      },
      dialogVisible: false,
      mcp_node: "",
      toolActive: {},
      toollist: [],
    };
  },
  mounted() {
    // 打开弹窗时再请求，避免重复请求问题
    // this.init();
  },
  methods: {
    async changeItem(value) {
      if (!value) return

      const item = this.serverOptions.find(item => item.name === value) || {}
      const itemIndex = this.serverOptions.findIndex(item => item.name === value)
      this.serverOptions = this.serverOptions.map((item) => ({...item, tools: []}))

      try {
        this.loading = true
        const res = await getMcpToolList({
          serverUrl: item.serverUrl,
        });
        if (res.code === 0) {
          this.serverOptions[itemIndex].tools = res.data && res.data.tools
              ? res.data.tools.map(item => ({...item, checked: false, disabled: false}))
              : []
        }
      } finally {
        this.loading = false
      }
    },
    handleOpenMcp() {
      /*window.parent.postMessage({ action: "navigateTo", path: "/mcp" }, "*");*/
      this.$router.push({path: "/mcp"})
    },
    handleTitleClick(name) {
      this.activeCollapseName = this.activeCollapseName === name ? "" : name;
    },
    handleRadioClick(index, toolIndex, val) {
      let obj = JSON.parse(JSON.stringify(this.serverOptions[index]));
      obj.tools[toolIndex].checked = !obj.tools[toolIndex].checked;
      this.$set(this.serverOptions, index, obj);
      if (obj.tools[toolIndex].checked) {
        this.mcp_node = val;
        this.toolActive = obj.tools[toolIndex];
        this.form.mcpName = obj.name;
        this.form.mcpDesc = obj.description;
        this.form.mcpServerUrl = obj.serverUrl;

        this.serverOptions.forEach((item, n) => {
          item.tools.forEach((i) => {
            if (val === i.name && n === index) {
              i.disabled = false;
            } else {
              i.disabled = true;
            }
          });
          this.$set(this.serverOptions, n, item);
        });
      } else {
        this.mcp_node = "";
        this.toolActive = {};
        this.form.mcpName = "";
        this.form.mcpDesc = "";
        this.form.mcpServerUrl = "";
        this.serverOptions.forEach((item, n) => {
          item.tools.forEach((i) => {
            i.disabled = false;
          });
          this.$set(this.serverOptions, n, item);
        });
      }
    },
    init() {
      this.loading = true;
      getList()
        .then((res) => {
          this.getToolsList(res.data.list);
        })
        .catch((err) => {
          this.loading = false;
        });
    },
    async getToolsList(list) {
      let mcpList = list || [];
      /*for (const item of mcpList) {
        try {
          const res = await getMcpToolList({
            serverUrl: item.serverUrl,
          });
          if (res.code === 0) {
            item.tools = res.data.tools || [];
            await item.tools.forEach((i) => {
              i.checked = false;
              i.disabled = false;
            });
          } else {
            item.tools = [];
          }
        } catch (error) {
          this.loading = false
        }
      }*/
      this.serverOptions = mcpList;
      this.loading = false;
    },
    handleClick(active) {
      this.active = active;
    },
    openDialog() {
      this.dialogVisible = true;
      this.init();
    },
    handleClose() {
      this.mcp_node = "";
      this.serverOptions = [];
      this.activeCollapseName = "";
      this.toolActive = {};
      this.form.mcpName = "";
      this.form.mcpDesc = "";
      this.form.mcpServerUrl = "";
      this.dialogVisible = false;
    },
    reset() {
      this.handleClose();
    },
    preInsert() {
      this.operation = 1;
      this.dialogVisible = true;
    },
    doSubmit() {
      let inputSchema = this.toolActive.inputSchema;
      let inputs = [];
      for (let o in inputSchema.properties) {
        inputs.push({
          name: o,
          desc: inputSchema.properties[o].description,
          type: "string",
          list_schema: null,
          object_schema: null,
          value: {
            content: "",
            type: "generated",
          },
          extra: {
            location: "body",
          },
        });
      }
      const form = JSON.parse(JSON.stringify(this.form));
      this.$emit("createMcp", inputs, form, this.mcp_node);
      this.handleClose();
    },
  },
};
</script>

<style lang="scss" scoped>
.form-content {
  /deep/.el-form-item {
    margin-bottom: 7px;
  }
  .el-select {
    width: 100%;
  }
}
.mcpCollapse {
  margin-top: 10px;
  /deep/.el-collapse-item__header {
    position: relative;
    display: flex !important;
    height: 48px !important;
    line-height: 48px !important;
    background-color: #fff !important;
    border-top: 1px solid #ebeef5 !important;
  }
  /deep/.el-icon-arrow-right {
    position: relative;
    top: auto;
  }
  .detail-mcpList {
    max-height: 250px;
    overflow-y: auto;
    margin: 10px 0 20px;
    li {
      padding: 10px;
      border: 1px solid #ddd;
      border-radius: 5px;
      margin-bottom: 10px;
    }
  }
  ul {
    list-style: none;
  }
  .el-collapse-item {
    &/deep/.is-active.el-collapse-item__header {
      //   & > div {
      background-color: $color !important;

      span,
      i {
        color: #fff;
      }
      //   }
    }
  }
}
.no-list {
  padding: 70px 0;
  text-align: center;

  strong {
    color: $color;
    cursor: pointer;
    text-decoration: underline;
  }
}
.mcpbtn {
  margin: 10px 0;
  float: right;
}
.mcpBtn {
  width: 100%;
  text-align: center;
  .el-button {
    float: none;
    background: transparent !important;
    color: $txt_color;
    border-color: transparent !important;
    border-radius: 0 !important;

    &.active {
      border-bottom-color: $color !important;
    }
  }
  /deep/ .el-button--primary span {
    color: $txt_color !important;
  }
}
.mcp-box {
  width: 100%;
  float: left;
  /deep/.el-radio-group {
    margin-top: 10px;
    display: flex;
    flex-direction: column;
    .el-radio {
      margin: 0px 30px 5px 10px;
    }
  }
}

/deep/.el-dialog__body > div {
  overflow: hidden;
}
</style>
