import { getHistoryList } from "@/api/explore";
import { getEmbeddingList } from "@/api/modelAccess";
export const app = {
    namespaced: true,
    state: {
        sessionStatus:-1,
        basicForm:{
            assistantId:'',
            avatar:'',
            instructions:'',
            name:'',
            description:''
        },
        expandForm:{
            fileList:[],
            starterPrompts:[{value:''}],
            models:[]
        },
        maxPicNum:6,
        //starterPrompts:[{value:''}],
        cacheData:{},
        historyAppList:[],
        embeddingList:[],
        fromList:'',
        permissionType:''
    },

    mutations: {
        SET_PERMISSION_TYPE(state,data){
            state.permissionType = data
        },
        CLEAR_PERMISSION_TYPE(state){
            state.permissionType = ''
        },
        SET_MAX_PICNUM(state,data){
            state.maxPicNum = data
        },
        CLEAR_MAX_PICNUM(state){
            state.maxPicNum = 6
        },
        SET_FROM_LIST(state,data){
            state.fromList = data
        },
        SET_SESSION_STATUS(state,data){
            state.sessionStatus = data
        },
        SET_BASIC_FORM(state,data){
            state.basicForm = data
        },
        SET_EXPAND_FORM(state,data){
            state.expandForm = data
        },
        SET_STARTER_PROMPTS(state,data){
            state.expandForm.starterPrompts = data
        },
        CACHE_APP_DETAIL(state,data){
            state.cacheData = data
        },
        INIT_STATE(state) {
            state.sessionStatus = -1
            state.basicForm = {
                assistantId:'',
                avatar:'',
                instructions:'',
                name:'',
                description:''
            }
            state.expandForm = {
                fileList:[],
                starterPrompts:[{value:''}],
                models:[]
            }
            state.starterPrompts = [{value:''}]
            state.cacheData = {}
        },
        REFRESH_RECENT_APP(state){

        },
        SET_HISTORYAPP_LIST(state,data){
            state.historyAppList = data
        },
        SET_EMBEDDING_LIST(state,data){
            state.embeddingList = data
        }
    },
    actions: {
        setPermissionType({ commit },data){
            commit('SET_PERMISSION_TYPE',data)
        },
        clearPermissionType({ commit }){
            commit('CLEAR_PERMISSION_TYPE')
        },
        setMaxPicNum({ commit },data){
            commit('SET_MAX_PICNUM',data)
        },
        clearMaxPicNum({ commit }){
            commit('CLEAR_MAX_PICNUM')
        },
        setFromList({ commit },data){
            commit('SET_FROM_LIST',data)
        },
        setStoreSessionStatus({ commit },data){
            // console.trace(data,'----status');
            commit('SET_SESSION_STATUS',data)
        },
        setBasicForm({ commit },data) {
            commit('SET_BASIC_FORM',data)
        },
        setExpandForm({ commit },data) {
            commit('SET_EXPAND_FORM',data)
        },
        setStarterPrompts({ commit },data) {
            commit('SET_STARTER_PROMPTS',data)
        },
        cacheAppDetail({ commit },data) {
            commit('CACHE_APP_DETAIL',data)
        },
        initState({ commit }){
            commit('INIT_STATE')
        },
        refreshRecentApp({ commit }){
            commit('REFRESH_RECENT_APP')
        },
        getHistoryList({commit},params=''){
            getHistoryList({history:params}).then(res=>{
                if(res.code === 0){
                    if(res.data.list && res.data.list.length >0){
                        const list = res.data.list.map(n =>{
                            return {...n,path: n.appType === 'agent'?`/explore/agent?id=${n.appId}`:`/explore/rag?id=${n.appId}`,active:'',hover:false,edit:false}
                        })
                        commit('SET_HISTORYAPP_LIST',list)
                    }
                }
            })
        },
        getEmbeddingList({commit}){
            getEmbeddingList({provider:'',model:''}).then(res =>{
                if(res.code === 0){
                    const list = res.data.list || []
                    commit('SET_EMBEDDING_LIST',list)
                }
            })
        }
    },
    getters: {
        basicForm:(state)=> state.basicForm,
        expandForm:(state)=> state.expandForm,
       // starterPrompts:(state)=> state.starterPrompts,
        cacheData:(state)=> state.cacheData,
        sessionStatus:(state)=> state.sessionStatus,
        historyAppList:(state) => state.historyAppList,
        embeddingList:(state) => state.embeddingList,
        fromList:(state) => state.fromList,
        maxPicNum:(state) => state.maxPicNum,
        permissionType:(state) => state.permissionType
    }
}
