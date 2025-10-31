<template>
  <div>
    <el-dialog
      title="绑定应用"
      :visible.sync="dialogVisible"
      width="50%"
      :before-close="handleClose">
      <div class="tool-typ">
        <div class="toolbtn">
          <div
            v-for="(item,index) in toolList" :key="index" @click="clickTool(item)"
            :class="[{'active':activeValue === item.value}]">
            {{ item.name }}
          </div>
        </div>
        <el-input
          v-model="toolName" placeholder="搜索应用" class="tool-input" suffix-icon="el-icon-search"
          @keyup.enter.native="searchTool" clearable/>
      </div>
      <div class="toolContent">
        <div @click="goCreate" class="createTool"><span class="el-icon-plus add"/>{{ createText() }}</div>
        <template v-for="(items, type) in contentMap">
          <div
            v-if="activeValue === type"
            v-for="(item,i) in items"
            :key="i">
            <div class="toolContent_item" @click="handleClick(item)">
              <div>
                <div v-if="item.type === 'custom'">{{ item.name }}</div>
                <div v-else>{{ item.name }}</div>
                <div class="tool-description">{{ item.description }}</div>
              </div>
              <el-button
                v-show="item.type !== 'custom'"
                v-if="!item.checked"
                type="text"
                @click="addTool(item)">
                添加
              </el-button>
              <el-button
                v-show="item.type !== 'custom'"
                v-else
                type="text"
                style="color:#ccc;">
                已添加
              </el-button>
              <i
                v-show="item.type === 'custom'"
                v-if="item.showTools"
                class="el-icon-caret-bottom"/>
              <i
                v-show="item.type === 'custom'"
                v-else
                class="el-icon-caret-left"/>
            </div>

            <!-- 展开的方法列表 -->
            <div
              class="toolContent_item"
              v-if="item.type === 'custom'"
              v-show="item.showTools"
              v-for="(method, index) in item.methods"
              :key="index">
              <div>
                <div>{{ method.methodName }}</div>
                <div class="tool-description">{{ item.description }}</div>
              </div>
              <el-button
                v-if="!method.checked"
                type="text"
                @click="addTool(item, method)">
                添加
              </el-button>
              <el-button
                v-else
                type="text"
                style="color:#ccc;">
                已添加
              </el-button>
            </div>
          </div>
        </template>
      </div>
    </el-dialog>
  </div>
</template>
<script>
import {getServerToolListCustom, addServerTool} from "@/api/mcp";

export default {
  props: ['assistantId'],
  data() {
    return {
      mcpServerId: '',
      toolName: '',
      dialogVisible: false,
      activeValue: 'custom',
      customInfos: [],
      customList: [],
      toolList: [
        {
          value: 'custom',
          name: '自定义工具'
        },
        // {
        //   value: 'agent',
        //   name: '智能体'
        // },
        // {
        //   value: 'workflow',
        //   name: '工作流'
        // },
        // {
        //   value: 'rag',
        //   name: '文本问答'
        // },
        // {
        //   value: 'builtIn',
        //   name: '内置工具'
        // },
      ]
    }
  },
  computed: {
    contentMap() {
      return {
        custom: this.customInfos,
      }
    }
  },
  created() {
    this.getCustomList('')
  },
  methods: {
    showDialog(detail) {
      this.mcpServerId = detail.mcpServerId
      this.customList = detail.tools.filter(tool => tool.type === 'custom');
      this.customInfos.forEach(item => {
        if (item.methods) {
          item.methods.forEach(method => {
            method.checked = this.customList.some(custom => custom.methodName === method.methodName)
          })
        }
        item.checked = false
        item.showTools = false
      })

      this.dialogVisible = true
    },
    getCustomList(name) {
      getServerToolListCustom({name}).then(res => {
        if (res.code === 0) {
          this.customInfos = (res.data.list || []).map(item => ({
            ...item,
            type: 'custom',
            showTools: false
          }))
        }
      })
    },
    goCreate() {
      if (this.activeValue === 'custom') {
        this.$router.push({path: '/tool?type=tool&mcp=custom'})
      }
    },
    createText() {
      if (this.activeValue === 'custom') {
        return '创建自定义工具'
      }
    },
    addTool(item, method) {
      if (method) {
        item = {...item, ...method}
      }
      const params = {
        id: item.customToolId,
        type: item.type,
        methodName: item.methodName,
        mcpServerId: this.mcpServerId
      }
      addServerTool(params).then(res => {
        if (res.code === 0) {
          method.checked = true;
          this.$message.success('工具添加成功');
          this.$emit('handleFetch');
          this.$nextTick(() => {
            this.$forceUpdate();
          });
        }
      })
    },
    searchTool() {
      if (this.activeValue === 'custom') {
        this.getCustomList(this.toolName)
      }
    },
    handleClick(item) {
      if (item.type === 'custom') {
        item.showTools = !item.showTools;
      }
    },
    handleClose() {
      this.activeValue = 'custom'
      this.dialogVisible = false;
    },
    clickTool(item) {
      this.activeValue = item.value;
    },
  }
}
</script>
<style lang="scss" scoped>
/deep/ {
  .el-dialog__body {
    padding: 10px 20px;
  }
}

.createTool {
  padding: 10px;
  cursor: pointer;

  .add {
    padding-right: 5px;
  }
}

.createTool:hover {
  color: $color;
}

.tool-typ {
  display: flex;
  justify-content: space-between;
  padding: 10px 0;
  border-bottom: 1px solid #dbdbdb;

  .toolbtn {
    display: flex;
    justify-content: flex-start;
    gap: 20px;

    div {
      text-align: center;
      padding: 5px 20px;
      border-radius: 6px;
      border: 1px solid #ddd;
      cursor: pointer;
    }
  }

  .tool-input {
    width: 200px;
  }
}

.toolContent {
  padding: 10px 0;
  max-height: 300px;
  overflow-y: auto;

  .toolContent_item {
    padding: 5px 20px;
    border: 1px solid #dbdbdb;
    border-radius: 6px;
    margin-bottom: 10px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: space-between;

    .tool-name {
      font-weight: bold;
      margin-bottom: 4px;
    }

    .tool-description {
      font-size: 10px;
      color: #999;
    }
  }

  .toolContent_item:hover {
    background: $color_opacity;
  }
}

.active {
  border: 1px solid $color !important;
  color: #fff;
  background: $color;
}
</style>