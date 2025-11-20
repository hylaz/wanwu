<template>
  <div class="mcp-content-box customize">
    <div class="mcp-content">
      <div class="card-search card-search-cust">
        <div>
          <p class="card-search-des">
            {{ $t('tool.custom.slogan') }}
          </p>
        </div>
        <div>
          <search-input :placeholder="$t('tool.custom.search')" ref="searchInput" @handleSearch="fetchList" />
        </div>
      </div>

      <div class="card-box">
        <div class="card card-item-create">
          <div class="app-card-create" @click="handleAddMCP('')">
            <div class="create-img-wrap">
              <img class="create-type" src="@/assets/imgs/create_tools.svg" alt="" />
              <img class="create-img" src="@/assets/imgs/create_icon.png" alt="" />
              <div class="create-filter"></div>
            </div>
            <span>{{ $t('tool.custom.addTitle') }}</span>
          </div>
        </div>
        <div
          v-if="list && list.length"
          class="card"
          v-for="(item, index) in list"
          :key="index"
          @click.stop="handleClick(item.customToolId)"
        >
          <div class="card-title">
            <img class="card-logo" :src="(item.avatar && item.avatar.path) ? avatarSrc(item.avatar.path) : defaultAvatar" alt=""/>
            <div class="mcp_detailBox">
              <span class="mcp_name">{{ item.name }}</span>
            </div>
            <el-dropdown
                placement="bottom">
              <span class="el-dropdown-link">
                <i class="el-icon-more"
                    @click.stop/>
              </span>
              <el-dropdown-menu slot="dropdown"  style="margin-top: -10px">
                <el-dropdown-item
                    @click.native="handleAddMCP(item.customToolId)">
                  {{$t('common.button.edit')}}
                </el-dropdown-item>
                <el-dropdown-item
                    @click.native="handleDelete(item)">
                  {{$t('common.button.delete')}}
                </el-dropdown-item>
              </el-dropdown-menu>
            </el-dropdown>
          </div>
          <div class="card-des">{{ item.description }}</div>
        </div>
      </div>
      <el-empty class="noData" v-if="!(list && list.length)" :description="$t('common.noData')"></el-empty>
    </div>
    <addDialog  ref="addDialog" @handleFetch="fetchList"></addDialog>
  </div>
</template>
<script>
import addDialog from "./addDialog.vue";
import SearchInput from "@/components/searchInput.vue"
import { getCustomList, deleteCustom } from "@/api/mcp";
import {avatarSrc} from "@/utils/util";
export default {
  components: { SearchInput, addDialog },
  data() {
    return {
      defaultAvatar: require("@/assets/imgs/toolImg.png"),
      list: [],
    };
  },
  mounted() {
    this.fetchList()
  },
  methods: {
    avatarSrc,
    fetchList() {
      const searchInput = this.$refs.searchInput
      const params = {
        name: searchInput.value,
      }
      getCustomList(params)
        .then((res) => {
          this.list = res.data.list || []
        })
    },
    handleClick(customToolId) {
      this.$refs.addDialog.showDialog(customToolId, true);
    },
    handleAddMCP(customToolId) {
      this.$refs.addDialog.showDialog(customToolId, false);
    },
    handleDelete(item) {
      this.$confirm(
        this.$t('tool.custom.deleteHint'),
        this.$t('common.confirm.title'),
        {
          confirmButtonText: this.$t('common.confirm.confirm'),
          cancelButtonText: this.$t('common.confirm.cancel'),
          dangerouslyUseHTMLString: true,
          type: "warning",
          center: true,
        }
      ).then(async () => {
        deleteCustom({
          customToolId: item.customToolId,
        }).then((res) => {
          if (res.code === 0) {
            this.$message.success(this.$t('common.info.delete'))
            this.fetchList()
          } else {
            this.$message.error( res.msg || this.$t('common.info.deleteErr'))
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
  width: 40px;
  height: 40px;
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