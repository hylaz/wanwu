<template>
  <div>
    <el-dialog
      title="选择应用"
      :visible.sync="dialogVisible"
      width="500px"
      v-loading="loading"
      append-to-body
      :close-on-click-modal="false"
      :before-close="handleClose"
    >
      <div v-if="options.length > 0">
        <label>选择和工作流关联的应用：</label>
        <el-select
          v-model="value"
          size="mini"
          placeholder="请选择"
          @change="handleChange"
        >
          <el-option
            v-for="item in options"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          >
          </el-option>
        </el-select>
      </div>
      <div v-else class="no-app">
        您暂无应用，请前往API Key管理进行
        <span @click="handleOpen">创建</span>
      </div>
      <div class="tips">
        <p>
          为了确保工作流的正常使用，请务必绑定应用并生成Token。Token是工作流运行的关键凭证，若未绑定应用或删除应用，将导致工作流无法正常工作。请按照以下步骤操作：
        </p>

                <li>在调试工作流中找到“选择和工作流关联的应用”选项。</li>
                <li>如果该选项中不存在可用应用，请前往<span @click="handleOpen">API Key管理</span>界面创建新的应用。</li>
                <li>按照提示完成应用绑定并生成Token。</li>

        <p class="warning">
          注意：删除应用将导致Token失效，工作流将无法继续使用。
        </p>
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button size="mini" @click="handleClose">取 消</el-button>
        <el-button size="mini" type="primary" @click="doSubmit"
          >确 定</el-button
        >
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { getAppList, getStaticToken } from "@/api/workflow";
export default {
  data() {
    return {
      options: [],
      dialogVisible: false,
      value: "",
      loading: false,
    };
  },
  created() {
    this.getAppList();
  },
  methods: {
    openDialog(row) {
      this.dialogVisible = true;
    },
    // 获取应用列表
    getAppList() {
      getAppList({ pageNo: 1, pageSize: 100 }).then((res) => {
          this.options = [];
          const data = res.data.list || [];
          for (let i = 0; i < data.length; i++) {
              if (data[i].status === true) {
                  this.options.push({
                  label: data[i].appName,
                  value:
                      data[i].apiKey + "|" + data[i].appId,
                  });
              }
          }
          if(this.options[0]){
              this.value = this.options[0].value
          }
      });
    },
    getToken(){},
    handleClose(){
        this.dialogVisible = false
    },
    handleChange(val){
        this.value = val
    },
    doSubmit(){
        if(!this.value){
            this.dialogVisible = false
            return
        }
        let params = this.value.split("|")
        this.loading = true
        getStaticToken({
            appid: params[1],
            apiKey: params[0],
        }).then((res)=>{
            this.loading = false
            if(res.code == 0){
                this.$emit("getToken",res.data.static_token)
                this.dialogVisible = false
            }
        }).catch((e)=>{
            this.loading = false
        })
    },
    handleOpen(){
          window.open("/aibase/userCenter/appV2");
    },
  }
}
</script>

<style lang="scss" scoped>
.no-app {
  text-align: center;
  span {
    margin: 0 3px;
    color: #d33a3a;
    font-weight: 700;
    cursor: pointer;
    text-decoration: underline;
  }
}

.tips{
    margin-top: 20px;
    p{
        font-size: 12px;
        margin: 10px 0px;
        color: #999;
    }
    .warning{
        color: #d33a3a;
        font-weight: 600;
    }
    li{
        font-size: 12px;
        margin-left:5px;
        color: #999;
        span{
            cursor:pointer;
            font-size: 12px;
            color: #d33a3a;
            padding:0 2px;
            opacity: 0.8;
            font-weight: 700;
        }
    }
}
</style>
