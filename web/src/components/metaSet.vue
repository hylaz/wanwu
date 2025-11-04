<template>
    <div class="metaSet">
        <div class="tool-typ">
            <el-button icon="el-icon-plus" type="primary" @click="addMataItem" size="small">新增条件</el-button>
            <el-switch v-model="metaDataFilterParams.filterEnable" active-color="var(--color)" @change="switchChange"></el-switch>
        </div>
        <div class="docMetaData">
            <div :class="['docMetaBox',metaDataFilterParams.metaFilterParams.length > 1 ? 'docMetaContainer':'']">
                <div
                    v-for="(item,index) in metaDataFilterParams.metaFilterParams"
                    class="docItem"
                >
                    <div class="docItem_data">
                        <el-select
                            v-model="item.key"
                            placeholder="请选择key"
                            @change="keyChange($event,item,index)"
                        >
                           <template #prefix>
                              <img class="prefix" src="@/assets/imgs/key.png" />
                           </template>
                            <el-option
                            v-for="meta in keyOptions"
                            :key="meta.metaKey"
                            :label="meta.metaKey + ' | ' + '[ '+meta.metaValueType+' ]'"
                            :value="meta.metaKey"
                            >
                            </el-option>
                        </el-select>
                    </div>
                    <el-divider direction="vertical"></el-divider>
                    <div class="docItem_data">
                        <el-select
                            v-model="item.condition"
                            placeholder="请选择条件"
                            @change="conditionChange($event,item)"
                        >
                            <template #prefix>
                              <img class="prefix" src="@/assets/imgs/condition.png" style="width:18px;"/>
                           </template>
                            <el-option
                            v-for="item in conditionOptions[item.type]"
                            :key="item.value"
                            :label="item.label"
                            :value="item.value"
                            >
                            </el-option>
                        </el-select>
                    </div>
                    <el-divider direction="vertical"></el-divider>
                    <div class="docItem_data">
                        <el-input
                            v-model="item.value"
                            v-if="item.type === 'string' || item.type === ''"
                            @blur="metaValueBlur(item)"
                            placeholder="请输入value"
                            :disabled="item.condition === 'empty' || item.condition === 'not empty'"
                        >
                        <template #prefix>
                            <img class="prefix" src="@/assets/imgs/value.png" style="width:16px;"/>
                        </template>
                        </el-input>
                        <el-input
                            v-model="item.value"
                            v-if="item.type === 'number'"
                            @blur="metaValueBlur(item)"
                            type="number"
                            placeholder="number"
                            :disabled="item.condition === 'empty'|| item.condition === 'not empty'"
                        >
                        <template #prefix>
                            <img class="prefix" src="@/assets/imgs/number.png" style="width:16px;"/>
                        </template>
                        </el-input>
                        <el-date-picker
                            v-if="item.type === 'time'"
                            v-model="item.value"
                            align="right"
                            format="yyyy-MM-dd HH:mm:ss"
                            value-format="timestamp"
                            type="datetime"
                            placeholder="选择日期时间"
                            :disabled="item.condition === 'empty' || item.condition === 'not empty'"
                        >
                        </el-date-picker>
                    </div>
                    <el-divider direction="vertical"></el-divider>
                    <div class="docItem_data docItem_data_btn">
                        <span
                            class="el-icon-delete setBtn"
                            @click="delMataItem(index)"
                        ></span>
                    </div>
                </div>
                <el-select
                    v-if="metaDataFilterParams.metaFilterParams.length > 1"
                    v-model="metaDataFilterParams.filterLogicType"
                    class="orAnd"
                    placeholder="条件"
                >
                    <el-option
                    v-for="item in conditions"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                    >
                    </el-option>
                </el-select>
            </div>
        </div>
    </div>
</template>
<script>
import {metaSelect} from "@/api/knowledge"
export default {
    props:{
        knowledgeId:{
            type:String,
            required: true,
            default:''
        },
        currentMetaData:{
            type:Object,
            default: () => ({})
        }
    },
    data(){
        return {
            metaDataFilterParams:{
                filterEnable:false,
                filterLogicType:'and',
                metaFilterParams:[]
            },
            keyOptions:[],
            conditions:[
                {
                    value:'and',
                    label:'且'
                },
                {
                    value:'or',
                    label:'或'
                }
            ],
            conditionOptions:{
                time:[
                        {
                            value:'is',
                            label:'是'
                        },
                        {
                            value:'before',
                            label:'早于'
                        },
                        {
                            value:'after',
                            label:'晚于'
                        },
                        {
                            value:'empty',
                            label:'为空'
                        },
                        {
                            value:'not empty',
                            label:'不为空'
                        }
                    ],
                string:[
                    {
                        value:'is',
                        label:'是'
                    },
                    {
                        value:'is not',
                        label:'不是'
                    },
                    {
                        value:'contains',
                        label:'包含'
                    },
                    {
                        value:'not contains',
                        label:'不包含'
                    },
                    {
                        value:'start with',
                        label:'开始是'
                    },
                    {
                        value:'end with',
                        label:'结束是'
                    },
                    {
                        value:'empty',
                        label:'为空'
                    },
                    {
                        value:'not empty',
                        label:'不为空'
                    }
                ],
                number:[
                    {
                        value:'=',
                        label:'等于'
                    },
                    {
                        value:'≠',
                        label:'不等于'
                    },
                    {
                        value:'>',
                        label:'大于'
                    },
                    {
                        value:'≥',
                        label:'大于等于'
                    },
                    {
                        value:'<',
                        label:'小于'
                    },
                    {
                        value:'≤',
                        label:'小于等于'
                    },
                    {
                        value:'empty',
                        label:'为空'
                    },
                    {
                        value:'not empty',
                        label:'不为空'
                    }
                ]
            }
        }
    },
    watch:{
       currentMetaData:{
          handler: function (val,old){
            if(val === null ){
                this.metaDataFilterParams = {
                     filterEnable:false,
                    filterLogicType:'and',
                    metaFilterParams:[]
                }
                return;
            }

            if(Object.keys(val).length > 0){
                this.metaDataFilterParams = JSON.parse(JSON.stringify(val))
            }
          },
          immediate:true,
          deep:true
       },
       knowledgeId:{
        handler(val,old){
            if(val){
                if(val !== old){
                    this.getList()
                }
            }
        },
        immediate:true,
       }
    },
    created(){
        // this.getList()
    },
    methods:{
        switchChange(val){
            if (!val) {
                this.metaDataFilterParams.metaFilterParams = [];
            }
        },
        conditionChange(e,item){
            item.value = '';
        },
        getMetaData(){
            this.metaDataFilterParams.metaFilterParams = this.metaDataFilterParams.metaFilterParams.map(item => {
                if (item.type === 'time') {
                    return {
                        ...item,
                        value: String(item.value)
                    };
                }
                return item;
            });
            return  { 'metaDataFilterParams':this.metaDataFilterParams}
        },
        getList(){
            metaSelect({knowledgeId:this.knowledgeId}).then(res =>{
                if(res.code === 0){
                    this.keyOptions = res.data.knowledgeMetaList || []
                }
            }).catch(() =>{})
        },
        metaValueBlur(item){
            if(item.condition === 'empty' || item.condition ==='not empty'){
                 return true;
            }else{
                if (!item.value) {
                    this.$message.warning("请输入value值");
                    return;
                }
            }
            
        },
        isEmpty(value){
            if (value === null ||value === 'null' || value === undefined || value === '') return true;
            return false;
        },
        validateRequiredFields(data){
            return data.some(field => {
               if (field && typeof field === 'object' && (field.condition === 'empty' || field.condition === 'not empty') ) {
                    return false;
                }
                if (field && typeof field === 'object' && 'value' in field) {
                    return this.isEmpty(field.value);
                }
         });
        },
        addMataItem(){
            if(this.metaDataFilterParams.filterEnable === false){
                this.$message.warning('请开启元数据配置后再进行添加')
                return;
            }
            if(this.metaDataFilterParams.metaFilterParams.length > 0){
                 if(this.validateRequiredFields(this.metaDataFilterParams.metaFilterParams)){
                    this.$message.warning('存在未填信息,请补充')
                    return
                 }
            }
            this.metaDataFilterParams.metaFilterParams.push({
                key:'',
                type:'',
                condition:'',
                value:''
            })
        },
        clearData(){
            this.metaDataFilterParams.metaFilterParams = [];
            this.metaDataFilterParams.filterLogicType = 'and';
            this.metaDataFilterParams.filterEnable = false
        },
        keyChange(val,item){
           item.value = '';
           item.condition = '';
           item.type = this.keyOptions.filter(i => i.metaKey === val).map(e => e.metaValueType)[0];
        },
        delMataItem(index){
            this.metaDataFilterParams.metaFilterParams.splice(index,1)
            if(this.metaDataFilterParams.metaFilterParams.length === 0){
                this.metaDataFilterParams.filterLogicType = 'and';
            }
        }
    }
}
</script>
<style lang="scss" scoped>
/deep/{
    .el-dialog__body{
        padding:10px 20px;
    }
    .el-divider--vertical{
        margin:0;
    }
    .el-input__prefix{
        display: flex;
        align-items:center;
        justify-content: center;
    }
    .el-icon-time{
        color: #727ff9;
    }
}
.metaSet{
    width:100%;
    .tool-typ{
        display: flex;
        justify-content:space-between;
        align-items:center;
    }
    .docMetaData {
        display:flex;
        justify-content:space-between;
        align-items:center;
        .docMetaBox{
            width:100%;
        }
        .docMetaContainer{
            position: relative;
            margin-left:65px;
            margin-top:15px;
        }
        .docMetaContainer::after{
            content: "";
            display: block;
            position: absolute;
            left: -20px;
            top: 50%;
            bottom: 0;
            width: 15px;
            height: 90%;
            transform: translateY(-50%);
            border-top-left-radius: 8px;
            border-bottom-left-radius: 8px;
            border: 1px solid rgba(16, 24, 40, .1411764706);
            border-right-width: 0;
        }
        .docItem {
            display: flex;
            flex:1;
            align-items: center;
            border-radius: 8px;
            background: #f7f8fa;
            margin-top: 10px;
            width:100%;
            .docItem_data {
            width:32%;
            display: flex;
            align-items: center;
            flex-wrap: wrap;
            padding: 5px 10px;
            .el-input,
            .el-select,
            .el-date-picker {
                width:100%;
            }
            .prefix{
                width:14px;
                margin-left: 5px;
            }
            .docItem_data_label {
                margin-right: 8px;
                display: flex;
                align-items: center;
                .question {
                color: #aaadcc;
                margin-left: 2px;
                cursor: pointer;
                }
            }
            .setBtn {
                font-size: 16px;
                cursor: pointer;
                color: $btn_bg;
            }
            }
            .docItem_data_type{
                width:80px;
            }
            .docItem_data_btn {
            display: flex;
            justify-content: center;
            width:30px;
            .el-icon-delete {
                margin-left: 5px;
            }
            }
        }
        .orAnd{
            width:65px;
            position: absolute;
            left: 0;
            top: 50%;
            transform: translate(-110%, -50%);
            padding: 6px 2px;
            z-index: 1;
        }
    }
}

</style>