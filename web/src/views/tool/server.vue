<template>
  <div class="mcp-content-box customize">
    <div class="mcp-content">
      <div class="card-search card-search-cust">
        <div>
          <p class="card-search-des" style="display: flex; align-items: center">
            <span>{{$t('tool.server.slogan')}}</span>
          </p>
        </div>
        <div>
          <search-input :placeholder="$t('tool.server.search')" ref="searchInput" @handleSearch="fetchList" />
        </div>
      </div>

      <div class="card-box">
        <div class="card card-item-create">
          <div class="app-card-create" @click="handleAddServer('')">
            <div class="create-img-wrap">
              <img class="create-type" src="@/assets/imgs/create_mcp.svg" alt="" />
              <img class="create-img" src="@/assets/imgs/create_icon.png" alt="" />
              <div class="create-filter"></div>
            </div>
            <span>{{$t('tool.server.create')}}</span>
          </div>
        </div>
        <div
            v-if="list && list.length"
            class="card"
            v-for="(item, index) in list"
            :key="index"
            @click.stop="handleClick(item.mcpServerId)"
        >
          <div class="card-title">
            <img class="card-logo" :src="item.avatar.path ? basePath + '/user/api/' + item.avatar.path : defaultAvatar" />
            <div class="mcp_detailBox">
              <span class="mcp_name">{{ item.name }}</span>
              <span class="mcp_from">
                <label>
                  已绑定{{ item.toolNum }}个工具
                </label>
              </span>
            </div>
            <i class="el-icon-edit-outline action-icon"
               style="margin-right: 20px"
               @click.stop="handleAddServer(item)"/>
            <i class="el-icon-delete-solid action-icon"
               @click.stop="handleDelete(item)"/>
          </div>
          <div class="card-des">{{ item.desc }}</div>
        </div>
      </div>

      <el-empty class="noData" v-if="!(list && list.length)" :description="$t('common.noData')"/>
    </div>
    <serverDialog ref="serverDialog" @handleFetch="fetchList()"/>
  </div>
</template>
<script>
import serverDialog from "./serverDialog.vue";
import SearchInput from "@/components/searchInput.vue"
import { getServerList, deleteServer } from "@/api/mcp";
export default {
  components: { SearchInput, serverDialog},
  data() {
    return {
      basePath: this.$basePath,
      defaultAvatar: require("@/assets/imgs/mcp_active.svg"),
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
      getServerList(params)
          .then((res) => {
            this.list = res.data.list || []
          })
    },
    handleClick(mcpServerId) {
      this.$router.push({path: `/tool/detail/server?mcpServerId=${mcpServerId}`})
    },
    handleAddServer(item) {
      this.$refs.serverDialog.showDialog(item)
    },
    handleDelete(item) {
      this.$confirm(
          "确定要删除 <span style='font-weight: bold;'>" + item.name + "</span> 该服务吗？",
          "提示",
          {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            dangerouslyUseHTMLString: true,
            type: "warning",
            center: true,
          }
      ).then(async () => {
        deleteServer({
          mcpServerId: item.mcpServerId,
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