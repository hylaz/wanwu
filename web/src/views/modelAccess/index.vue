<template>
  <div class="page-wrapper modelAccess">
    <div class="table-wrap list-common wrap-fullheight">
      <div class="page-title">
        <img class="page-title-img" src="@/assets/imgs/model.png" alt="" />
        <span class="page-title-name">{{$t('modelAccess.title')}}</span>
      </div>
      <div class="table-box">
        <div class="table-form">
          <el-select
            v-model="params.provider"
            :placeholder="$t('modelAccess.table.publisher')"
            class="modelAccess-select no-border-select"
            clearable
            @change="searchData()"
          >
            <el-option v-for="item in providerList" :key="item.key" :label="item.name" :value="item.key" />
          </el-select>

          <el-select
            v-model="params.modelType"
            :placeholder="$t('modelAccess.table.modelType')"
            class="modelAccess-select no-border-select"
            style="margin-left: 15px"
            clearable
            @change="searchData()"
          >
            <el-option v-for="item in modelTypeList" :key="item.key" :label="item.name" :value="item.key" />
          </el-select>
          <div style="width: 100%; display: inline-block; float: right; text-align: right; margin-top: -30px">
            <el-input
              v-model="params.displayName"
              prefix-icon="el-icon-search"
              class="no-border-input"
              style="width: 240px; margin-right: 10px"
              :placeholder="$t('modelAccess.table.modelName')"
              @keyup.enter.native="searchName"
              @clear="searchData()"
              clearable
            />
            <el-button class="add-bt" size="mini" type="primary" @click="preInsert">
              <img
                style="width: 14px; margin-right: 5px; display: inline-block; vertical-align: middle"
                src="@/assets/imgs/modelImport.png"
                alt=""
              />
              <span style="display: inline-block; vertical-align: middle">{{$t('modelAccess.import')}}</span>
            </el-button>
          </div>
        </div>
        <div class="card-wrapper">
          <div class="card-item card-item-create">
            <div class="app-card-create" @click="preInsert">
              <div class="create-img-wrap">
                <img class="create-type" src="@/assets/imgs/create_model.svg" alt="" />
                <img class="create-img" src="@/assets/imgs/create_icon.png" alt="" />
                <div class="create-filter"></div>
              </div>
              <span>{{$t('modelAccess.import')}}</span>
            </div>
          </div>
          <div
            v-if="tableData && tableData.length"
            class="card-item"
            v-for="(item, index) in tableData"
            :key="item.model + index"
            @click="preUpdate(item)"
          >
            <div class="card-top">
              <img class="card-img" v-if="item.avatar && item.avatar.path" :src="basePath + '/user/api/' + item.avatar.path">
              <div class="card-title">
                <div class="card-name" :title="item.displayName || item.model">
                  {{item.displayName || item.model}}
                </div>
              </div>
              <div class="card-top-right" @click.stop="">
                <el-switch
                  @change="(val) => {changeStatus(item, val)}"
                  style="width: 32px;"
                  v-model="item.isActive"
                  active-text=""
                  inactive-text=""
                />
                <el-dropdown @command="handleCommand" placement="top">
                  <span class="el-dropdown-link">
                    <i class="el-icon-more more"></i>
                  </span>
                  <el-dropdown-menu slot="dropdown">
                    <el-dropdown-item :command="{type: 'edit', item}">
                      <i class="el-icon-edit-outline card-opera-icon"></i>
                      {{$t('common.button.edit')}}
                    </el-dropdown-item>
                    <el-dropdown-item class="card-delete" :command="{type: 'delete', item}">
                      <i class="el-icon-delete card-opera-icon" />
                      {{$t('common.button.delete')}}
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </el-dropdown>
              </div>
            </div>
            <div class="card-middle">
              <div v-if="item.tags" class="card-type" v-for="it in item.tags">{{it.text}}</div>
            </div>
            <div class="card-bottom">
              <div
                :class="['card-bottom-provider', {'no-publishData': !item.updatedAt}]"
                :title="providerObj[item.provider] || '--'"
              >
                {{$t('modelAccess.table.publisher')}}: {{providerObj[item.provider] || '--'}}
              </div>
              <div style="float: right">
                {{item.updatedAt ? item.updatedAt.split(' ')[0] : '--'}} {{$t('modelAccess.table.update')}}
              </div>
            </div>
          </div>
        </div>
        <el-empty class="noData" v-if="!(tableData && tableData.length)" :description="$t('common.noData')"></el-empty>
      </div>
      <CreateSelectDialog ref="createSelectDialog" @showCreate="showCreate" />
      <CreateDialog ref="createDialog" @reloadData="searchData" />
    </div>
  </div>
</template>

<script>
  import Pagination from "@/components/pagination.vue"
  import { fetchModelList, deleteModel, changeModelStatus, getModelDetail } from "@/api/modelAccess"
  import CreateDialog from './components/createDialog.vue'
  import CreateSelectDialog from "./components/createSelectDialog.vue"
  import { MODEL_TYPE_OBJ, PROVIDER_OBJ, PROVIDER_TYPE, MODEL_TYPE } from "./constants"

  export default {
    components: { Pagination, CreateSelectDialog, CreateDialog },
    data() {
      return {
        listApi: fetchModelList,
        providerList: PROVIDER_TYPE,
        modelTypeList: MODEL_TYPE,
        basePath: this.$basePath,
        modelTypeObj: MODEL_TYPE_OBJ,
        providerObj: PROVIDER_OBJ,
        tableData: [],
        params: {
          provider: '',
          modelType: '',
          displayName: ''
        },
        loading: false,
      }
    },
    mounted() {
      this.getTableData()
    },
    methods: {
      async getTableData(params) {
        this.loading = true
        try {
          const res = await fetchModelList({...params})
          const tableData = res.data ? res.data.list || [] : []
          this.tableData = [...tableData]
        } finally {
          this.loading = false
        }
      },
      searchData(isCreate) {
        if (isCreate) {
          this.params.provider = ''
          this.params.modelType = ''
          this.params.displayName = ''
        }
        this.getTableData({...this.params})
      },
      searchName(e) {
        if (e.keyCode === 13) {
          this.searchData()
        }
      },
      handleCommand(value) {
        const {type, item} = value || {}
        switch (type) {
          case 'edit':
            this.preUpdate(item)
            break
          case 'delete':
            this.preDel(item)
            break
        }
      },
      preInsert() {
        this.$refs.createSelectDialog.openDialog()
      },
      showCreate(item) {
        this.$refs.createDialog.openDialog(item.key)
      },
      preUpdate(row) {
        const {modelId, provider} = row || {}

        getModelDetail({modelId}).then(res => {
          const rowObj = res.data || {}
          const newRow = {...rowObj, ...rowObj.config}
          this.$refs.createDialog.openDialog(provider, newRow)
        })
      },
      preDel(row) {
        this.$confirm(this.$t('modelAccess.confirm.delete'), this.$t('common.confirm.title'), {
          confirmButtonText: this.$t('common.confirm.confirm'),
          cancelButtonText: this.$t('common.confirm.cancel'),
          type: 'warning'
        }).then(async () => {
          const {modelId} = row || {}
          let res = await deleteModel({modelId})
          if (res.code === 0) {
            this.$message.success(this.$t('common.message.success'))
            await this.getTableData()
          }
        })
      },
      changeStatus(row, val) {
        this.$confirm(val ? this.$t('modelAccess.confirm.start') : this.$t('modelAccess.confirm.stop'), this.$t('common.confirm.title'), {
          confirmButtonText: this.$t('common.confirm.confirm'),
          cancelButtonText: this.$t('common.confirm.cancel'),
          type: 'warning'
        }).then(async() => {
          const {modelId} = row || {}
          let res = await changeModelStatus({modelId, isActive: val})
          if (res.code === 0) {
            this.$message.success(this.$t('common.message.success'))
            await this.getTableData()
          }
        }).catch(() => {
          this.getTableData()
        })
      },
    }
  }
</script>

<style lang="scss" scoped>
.routerview-container {
  top:0;
}
.table-box {
  padding: 20px 20px 0;
  .table-form {
    width: 100%;
    padding-bottom: 20px;
    clear: both;
  }
  .table-header {
    font-size: 16px;
    font-weight: bold;
    color: #555;
  }
  .add-bt {
    margin: 0 2px 20px 5px;
  }
}
.modelAccess-select {
  width: 200px;
}
.mark-textArea /deep/ {
  .el-textarea__inner {
    font-family: inherit;
    font-size: inherit;
  }
}
.card-wrapper {
  margin: 0 -10px;
}
.card-item {
  display: inline-block;
  width: calc((100% / 4) - 20px);
  height: 180px;
  vertical-align: middle;
  margin: 0 10px 20px;
  background: #FFFFFF;
  box-shadow: 0 1px 4px 0 rgba(0,0,0,0.15);
  border-radius: 8px;
  padding: 18px 10px 16px 14px;
  position: relative;
  cursor: pointer;
  .card-top {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  .card-img {
    width: 46px;
    height: 46px;
    object-fit: cover;
    /*padding: 10px 5px;*/
    background: #FFFFFF;
    box-shadow: 0 1px 4px 0 rgba(0,0,0,0.15);
    border-radius: 8px;
    border: 0 solid #D9D9D9;
    margin-right: 10px;
  }
  .card-title {
    width: calc(100% - 90px);
  }
  .card-name {
    font-size: 18px;
    color: #434343;
    font-weight: bold;
    width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    line-clamp: 2;
    -webkit-line-clamp: 2;
  }
  .card-middle {
    padding-top: 14px;
  }
  .card-type {
    display: inline-block;
    padding: 2px 12px;
    border-radius: 2px;
    color: #6977F9;
    background: #ECEEFE;
    margin-top: 2px;
    margin-right: 8px;
  }
  .card-top-right {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    margin-left: 5px;
  }
  .more {
    margin-left: 5px;
    cursor: pointer;
    transform: rotate(90deg);
    font-size: 16px;
    color: #8c8c8f;
  }
  .more:hover {
    color: $color;
  }
  .card-bottom {
    position: absolute;
    color: #88888B;
    bottom: 14px;
    left: 15px;
    right: 12px;
    div {
      display: inline-block;
    }
    .card-bottom-provider {
      width: calc(100% - 100px);
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    .no-publishData.no-publishData {
      width: calc(100% - 45px);
    }
  }
}
.card-item:hover {
  /*background-image: url('../../assets/imgs/cardBg.png');
  background-size: cover;*/
  border: 1px solid $color;
}
.card-item-create {
  background: #F4F5FF;
  box-shadow: 0 1px 4px 0 rgba(0,0,0,0.15);
  border: 1px solid rgba(56,75,247,0.47);
  .app-card-create {
    width: 100%;
    height: 100%;
    text-align: center;
    display: flex;
    align-items: center;
    justify-content: center;
    .create-img-wrap {
      display: inline-block;
      vertical-align: middle;
      margin-right: 10px;
      position: relative;
      .create-img {
        width: 40px;
        height: 40px;
        border-radius: 8px;
        background: $color;
        padding: 10px;
      }
      .create-filter {
        width: 40px;
        height: 8px;
        background: rgba(2, 81, 252, 0.3);
        filter: blur(5px);
        position: absolute;
        bottom: -6px;
      }
      .create-type {
        width: 30px;
        position: absolute;
        background: rgba(171,198,255,0.5);
        backdrop-filter: blur(6.55px);
        border-radius: 5px;
        padding: 5px;
        top: -10px;
        left: -12px;
      }
    }
    span {
      display: inline-block;
      vertical-align: middle;
      font-size: 16px;
      color: #434C6C;
      font-weight: bold;
    }
  }
}
/deep/ .el-dropdown-menu__item.card-delete:hover {
  color: #FF4D4F !important;
  background: #FBEAE8 !important;
}
.card-opera-icon {
  font-size: 15px;
}
.modelAccess .noData {
  width: 100%;
  text-align: center;
  margin-top: -60px;
  /deep/ .el-empty__description p {
    color: #B3B1BC;
  }
}
</style>
