<template>
  <div class="tempSquare-management">
    <div class="tempSquare-content-box tempSquare-third">
      <div class="tempSquare-main">
        <div class="tempSquare-content">
          <div class="tempSquare-card-box">
            <div class="card-search card-search-cust" v-if="!templateUrl">
              <div>
                <span
                  v-for="item in typeList"
                  :key="item.key"
                  :class="['tab-span', {'is-active': typeRadio === item.key}]"
                  @click="changeTab(item.key)"
                >
                  {{ item.name }}
                </span>
              </div>
              <search-input
                style="margin-right: 2px"
                placeholder="输入名称搜索"
                ref="searchInput"
                @handleSearch="doGetWorkflowTempList"
              />
            </div>

            <div class="card-loading-box" v-if="list.length && !templateUrl">
              <div class="card-box" v-loading="loading">
                <div
                  class="card"
                  v-for="(item, index) in list"
                  :key="index"
                  @click.stop="handleClick(item)"
                >
                  <div class="card-title">
                    <img
                      class="card-logo"
                      v-if="item.avatar && item.avatar.path"
                      :src="item.avatar.path"
                    />
                    <div class="mcp_detailBox">
                      <span class="mcp_name">{{ item.name }}</span>
                      <span class="mcp_from">
                        <label>
                          作者：{{ item.author }}
                        </label>
                      </span>
                    </div>
                  </div>
                  <div class="card-des">{{ item.desc }}</div>
                  <div class="card-bottom">
                    <div class="card-bottom-left">下载量：{{item.downloadCount || 0}}</div>
                    <div class="card-bottom-right">
                      <i v-if="!isPublic" class="el-icon-copy-document" title="复制" @click.stop="copyTemplate(item)"></i>
                      <i class="el-icon-download" title="下载" @click.stop="downloadTemplate(item)"></i>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div v-else class="empty">
              <el-empty :description="$t('common.noData')"></el-empty>
            </div>
          </div>
        </div>
      </div>
    </div>
    <HintDialog :templateUrl="templateUrl" ref="hintDialog" />
    <CreateForm type="clone" ref="cloneWorkflowDialog" />
  </div>
</template>
<script>
import { getWorkflowTempList, downloadWorkflow } from "@/api/templateSquare"
import SearchInput from "@/components/searchInput.vue"
import HintDialog from "./components/hintDialog.vue"
import CreateForm from "@/views/workflowList/components/createForm.vue"
export default {
  components: { SearchInput, HintDialog, CreateForm },
  props: {
    isPublic: true,
    type: ''
  },
  data() {
    return {
      basePath: this.$basePath,
      category: '全部',
      list: [],
      templateUrl: '',
      loading: false,
      typeRadio: 'all',
      typeList: [
        {name: '全部', key: 'all'},
        {name: '政务', key: 'gov'},
        {name: '工业', key: 'industry'},
        {name: '文教', key: 'edu'},
        {name: '医疗', key: 'medical'},
        {name: '数据', key: 'data'},
        {name: '创作', key: 'create'},
        {name: '搜索', key: 'search'},
      ]
    };
  },
  mounted() {
    this.doGetWorkflowTempList()
  },
  methods: {
    changeTab(key) {
      this.typeRadio = key
      this.$refs.searchInput.value = ''
      this.doGetWorkflowTempList()
    },
    showHintDialog() {
      this.$refs.hintDialog.openDialog()
    },
    doGetWorkflowTempList() {
      const searchInput = this.$refs.searchInput
      let params = {
        name: searchInput.value,
        category: this.typeRadio,
      }

      getWorkflowTempList(params)
        .then((res) => {
          const {downloadLink = {}, list} = res.data || {}
          this.templateUrl = downloadLink.url
          if (downloadLink.templateUrl) this.showHintDialog()

          this.list = list || []
          this.loading = false
        })
        .catch(() => this.loading = false)
    },
    copyTemplate(item) {
      this.$refs.cloneWorkflowDialog.openDialog(item)
    },
    downloadTemplate(item) {
      downloadWorkflow({ templateId : item.templateId }).then(response => {
        const blob = new Blob([response], { type: response.type })
        const url = URL.createObjectURL(blob);
        const link = document.createElement("a")
        link.href = url
        link.download = row.name + '.json'
        link.click()
        window.URL.revokeObjectURL(link.href)
      })
    },
    handleClick(val) {
      const path = `${this.isPublic ? '/public' : ''}/templateSquare/detail`
      this.$router.push({path, query: { templateSquareId: val.templateId, type: this.type }})
    },
  },
};
</script>

<style lang="scss" scoped>
.tempSquare-management {
  height: calc(100% - 50px);
  .tempSquare-content-box {
    height: calc(100% - 145px);
  }
  .tempSquare-content {
    padding: 0 20px;
    width: 100%;
    height: 100%;
  }

  .tempSquare-third{
    min-height: 600px;
    .tab-span {
      display: inline-block;
      vertical-align: middle;
      padding: 6px 12px;
      border-radius: 6px;
      color: $color_title;
      cursor: pointer;
    }
    .tab-span.is-active {
      color: $color;
      background: #fff;
      font-weight: bold;
    }
    .tempSquare-main{
      display: flex;
      padding: 0 20px;
      height: 100%;
      .tempSquare-content{
        display: flex;
        width:100%;
        padding: 0;
        height: 100%;
        .tempSquare-menu{
          margin-top: 10px;
          margin-right: 20px;
          width: 90px;
          height: 450px;
          border: 1px solid $border_color;
          text-align: center;
          border-radius: 6px;
          color: #333;
          p{
            line-height: 28px;
            margin:10px 0;
          }
          .active{
            background: rgba(253, 231, 231, 1);
          }
        }
        .tempSquare-card-box{
          width: 100%;
          height: 100%;
          .input-with-select {
            width: 300px;
          }
          .card-loading-box{
            .card-box {
              display: flex;
              flex-wrap: wrap;
              margin: 6px -10px 0;
              align-content: flex-start;
              padding-bottom: 20px;
              /*overflow: auto;*/
              .card {
                position: relative;
                padding: 20px 16px;
                border-radius: 12px;
                height: fit-content;
                background: #fff;
                display: flex;
                flex-direction: column;
                align-items: center;
                width: calc((100% / 4) - 20px);
                margin: 0 10px 20px;
                box-shadow: 0 1px 4px 0 rgba(0, 0, 0, 0.15);
                border: 1px solid rgba(0, 0, 0, 0);
                &:hover {
                  cursor: pointer;
                  box-shadow: 0 2px 8px #171a220d, 0 4px 16px #0000000f;
                  border: 1px solid $border_color;

                  .action-icon {
                    display: block;
                  }
                }
                .card-title {
                  display: flex;
                  width: 100%;
                  border-bottom: 1px solid #ddd;
                  padding-bottom: 7px;
                  .svg-icon {
                    width: 50px;
                    height: 50px;
                  }
                  .mcp_detailBox {
                    width: calc(100% - 70px);
                    margin-left: 10px;
                    display: flex;
                    flex-direction: column;
                    justify-content: space-between;
                    padding: 3px 0;
                    .mcp_name {
                      display: block;
                      font-size: 15px;
                      font-weight: 700;
                      overflow: hidden;
                      white-space: nowrap;
                      text-overflow: ellipsis;
                      color: #5d5d5d;
                    }
                    .mcp_from {
                      label {
                        padding: 3px 7px;
                        font-size: 12px;
                        color: #84868c;
                        background: #f2f5f9;
                        border-radius: 3px;
                        display: block;
                        height: 22px;
                        width: 100%;
                        overflow: hidden;
                        text-overflow: ellipsis;
                        white-space: nowrap;
                      }
                    }
                  }

                  margin-bottom: 13px;
                }
                .card-des {
                  width: 100%;
                  display: -webkit-box;
                  text-overflow: ellipsis;
                  color: #5d5d5d;
                  font-weight: 400;
                  overflow: hidden;
                  -webkit-line-clamp: 2;
                  line-clamp: 2;
                  -webkit-box-orient: vertical;
                  font-size: 13px;
                  height: 36px;
                  word-wrap: break-word;
                }
                .card-bottom {
                  width: 100%;
                  display: flex;
                  align-items: center;
                  justify-content: space-between;
                  margin-top: 14px;
                  margin-bottom: -6px;
                  .card-bottom-left {
                    color: #888;
                  }
                  .card-bottom-right {
                    i {
                      margin-left: 5px;
                      cursor: pointer;
                    }
                  }
                }
              }

              .loading-tips{
                height: 20px;
                color: #999;
                text-align: center;
                display: block;
                width: 100%;
                i{
                  font-size: 18px;
                }
              }
            }
          }
        }
      }
    }
    .card-logo{
      width: 50px;
      height: 50px;
      object-fit: cover;
    }
  }
  .card-search {
    text-align: right;
    padding: 10px 0;
  }
  .card-search-cust {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .empty{
    width: 200px;
    height: 100px;
    margin: 50px auto;
  }
}
</style>