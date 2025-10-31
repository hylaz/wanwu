<template>
  <div class="mcp-content-box customize">
    <div class="mcp-content">
      <div class="card-search card-search-cust">
        <div>
          <p class="card-search-des" style="display: flex; align-items: center">
            <span>集合各MCP server，即选即用。可支持在插件编排工具中使用。如使用自定义MCP，请在自定义TAB中进行添加和查看</span>
            <LinkIcon type="mcp" />
          </p>
          <!--<div class="radio-box" >
            <el-radio-group v-model="source" @input="radioChange" size="mini">
              <el-radio label="">全部</el-radio>
              <el-radio label="1">来自我添加</el-radio>
              <el-radio label="2">来自广场发送</el-radio>
            </el-radio-group>
          </div>-->
        </div>
        <div>
          <search-input placeholder="请输入MCP名称进行搜索" ref="searchInput" @handleSearch="fetchList" />
          <el-button size="mini" type="primary" @click="handleAddMCP">导入</el-button>
        </div>
      </div>

      <div class="card-box">
        <div class="card card-item-create">
          <div class="app-card-create" @click="handleAddMCP">
            <div class="create-img-wrap">
              <img class="create-type" src="@/assets/imgs/create_mcp.svg" alt="" />
              <img class="create-img" src="@/assets/imgs/create_icon.png" alt="" />
              <div class="create-filter"></div>
            </div>
            <span>导入MCP</span>
          </div>
        </div>
        <div
          v-if="list && list.length"
          class="card"
          v-for="(item, index) in list"
          :key="index"
          @click.stop="handleClick(item)"
        >
          <div class="card-title">
            <img class="card-logo" v-if="item.avatar && item.avatar.path" :src="basePath + '/user/api/' + item.avatar.path" />
            <div class="mcp_detailBox">
              <span class="mcp_name">{{ item.name }}</span>
              <span class="mcp_from">
                <label>
                  {{ item.from }}
                </label>
              </span>
            </div>
            <i class="el-icon-edit-outline action-icon"
               style="margin-right: 20px"
              @click.stop="handleEdit(item)"/>
            <i
              class="el-icon-delete-solid action-icon"
              @click.stop="handleDelete(item)"
            ></i>
          </div>
          <div class="card-des">{{ item.desc }}</div>
        </div>
      </div>

      <!--<div class="no-list" v-if="list.length === 0 && is">
        <div>
          <i class="el-icon-circle-plus-outline" @click="handleAddMCP"></i>
          <span>添加你的第一个MCP Server</span>
        </div>
      </div>-->
      <el-empty class="noData" v-if="!(list && list.length)" :description="$t('common.noData')"></el-empty>
    </div>
    <addDialog :dialogVisible="addOpen" :initialData="dialogParams" @handleFetch="fetchList()" @handleClose="handleClose"></addDialog>
  </div>
</template>
<script>
import addDialog from "./addDialog.vue";
import SearchInput from "@/components/searchInput.vue"
import { getList, setDelete } from "@/api/mcp";
import LinkIcon from "@/components/linkIcon.vue"
export default {
  components: { LinkIcon, SearchInput, addDialog },
  data() {
    return {
      basePath: this.$basePath,
      addOpen: false, // 自定义添加mcp开关
      dialogParams: {
        name: "",
        from: "",
        sseUrl: "",
        desc: "",
      }, // 添加自定义mcp参数
      list: [],
    };
  },
  mounted() {
    this.fetchList()
  },
  methods: {
    fetchList() {
      const searchInput = this.$refs.searchInput
      const params = {
        name: searchInput.value,
      }
      getList(params)
        .then((res) => {
          this.list = res.data.list || []
        })
    },
    handleClick(val) {
      // mcpSquareId 有值 MCP广场, 否则自定义
      this.$router.push({path: `/tool/detail/custom?mcpId=${val.mcpId}&mcpSquareId=${val.mcpSquareId}`})
    },
    handleAddMCP() {
      this.addOpen = true;
    },
    handleEdit(item) {
      this.addOpen = true;
      this.dialogParams = {
        ...item
      };
    },
    handleClose() {
      this.addOpen = false
      this.dialogParams = {
        name: "",
        from: "",
        sseUrl: "",
        desc: "",
      };
    },
    handleDelete(item) {
      this.$confirm(
        "删除后，历史引用了本MCP服务的智能体将自动取消引用，且此操作不可撤回,确定要删除吗？",
        "提示",
        {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          dangerouslyUseHTMLString: true,
          type: "warning",
          center: true,
        }
      ).then(async () => {
        setDelete({
          mcpId: item.mcpId,
        }).then((res) => {
          if (res.code === 0) {
            this.$message.success("删除成功")
            this.fetchList()
          } else {
            this.$message.error( res.msg || '删除失败')
          }
        })
      })
    },
  },
};
</script>
<style lang="scss">
.card-search-cust {
  text-align: left !important;

  .radio-box {
    margin: 20px 0 0 0 !important;
  }
}
.card-logo{
  width: 50px;
  height: 50px;
  object-fit: cover;
}
.mcp-content-box .noData {
  width: 100%;
  text-align: center;
  margin-top: -60px;
  /deep/ .el-empty__description p {
    color: #B3B1BC;
  }
}
</style>