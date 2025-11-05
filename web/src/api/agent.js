import service from "@/utils/request"
const BASE_URL = '/user/api/v1'
export const createAgent = (data)=>{
    return service({
        url: `${BASE_URL}/assistant`,
        method: 'post',
        data
    })
}

export const updateAgent = (data)=>{
    return service({
        url: `${BASE_URL}/assistant`,
        method: 'put',
        data
    })
}
export const delAgent = (data)=>{
    return service({
        url: `${BASE_URL}/assistant`,
        method: 'delete',
        data
    })
}
export const getAgentInfo = (params)=>{
    return service({
        url: `${BASE_URL}/assistant`,
        method: 'get',
        params
    })
}
export const getAgentDetail = (params)=>{
    return service({
        url: `${BASE_URL}/assistant/draft`,
        method: 'get',
        params
    })
}
export const putAgentInfo = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/config`,
        method: 'put',
        data
    })
}
export const createConversation = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/conversation`,
        method: 'post',
        data
    })
}
export const delConversation = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/conversation`,
        method: 'delete',
        data
    })
}
export const getConversationHistory = (params)=>{
    return service({
        url: `${BASE_URL}/assistant/conversation/detail`,
        method: 'get',
        params
    })
}
export const getConversationlist = (params)=>{
    return service({
        url: `${BASE_URL}/assistant/conversation/list`,
        method: 'get',
        params
    })
}
export const getActionInfo = (params)=>{
    return service({
        url: `${BASE_URL}/assistant/action`,
        method: 'get',
        params
    })
}
export const editActionInfo = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/action`,
        method: 'put',
        data
    })
}
export const addActionInfo = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/action`,
        method: 'post',
        data
    })
}
export const delActionInfo = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/action`,
        method: 'delete',
        data
    })
}
export const enableAction = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/action/enable`,
        method: 'put',
        data
    })
}
export const addWorkFlowInfo = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/tool/workflow`,
        method: 'post',
        data
    })
}
export const delWorkFlowInfo = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/tool/workflow`,
        method: 'delete',
        data
    })
}
export const enableWorkFlow = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/tool/workflow/switch`,
        method: 'put',
        data
    })
}
export const agentStream = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/stream`,
        method: 'post',
        data
    })
}
export const agentTestStream = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/test/stream`,
        method: 'post',
        data
    })
}
export const getAgentList = (params)=>{
    return service({
        url: `${BASE_URL}/assistant/list`,
        method: 'get',
        params
    })
}

//删除mcp工具
export const deleteMcp = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/tool/mcp`,
        method: 'delete',
        data
    })
}
//添加mcp工具
export const addMcp = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/tool/mcp`,
        method: 'post',
        data
    })
}
//启停mcp工具
export const enableMcp = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/tool/mcp/switch`,
        method: 'put',
        data
    })
}

// 删除自定义、内置工具
export const delCustomBuiltIn = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/tool`,
        method: 'delete',
        data
    })
}
// 添加自定义、内置工具
export const addCustomBuiltIn = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/tool`,
        method: 'post',
        data
    })
}
// 启停自定义、内置工具
export const switchCustomBuiltIn = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/tool/switch`,
        method: 'put',
        data
    })
}
//工具列表
export const toolList = (data)=>{
    return service({
        url: `${BASE_URL}/tool/select`,
        method: 'get',
        params:data
    })
}
//工具下面的action列表
export const toolActionList = (data)=>{
    return service({
        url: `${BASE_URL}/tool/action/list`,
        method: 'get',
        params:data
    })
}
//内置工具下面的action详情
export const toolActionDetail = (data)=>{
    return service({
        url: `${BASE_URL}/tool/action/detail`,
        method: 'get',
        params:data
    })
}
//mcp工具列表
export const mcptoolList = (data)=>{
    return service({
        url: `${BASE_URL}/mcp/select`,
        method: 'get',
        params:data
    })
}
//mcp工具下面的action列表
export const mcpActionList = (data)=>{
    return service({
        url: `${BASE_URL}/mcp/action/list`,
        method: 'get',
        params:data
    })
}
    
//编辑url
export const editOpenurl = (data)=>{
    return service({
        url: `${BASE_URL}/appspace/app/openurl`,
        method: 'put',
        data
    })
}
//创建url
export const createOpenurl = (data)=>{
    return service({
        url: `${BASE_URL}/appspace/app/openurl`,
        method: 'post',
        data
    })
}
//删除应用url
export const delOpenurl = (data)=>{
    return service({
        url: `${BASE_URL}/appspace/app/openurl`,
        method: 'delete',
        data
    })
}
//获取应用url列表
export const getOpenurl = (data)=>{
    return service({
        url: `${BASE_URL}/appspace/app/openurl/list`,
        method: 'get',
        params:data
    })
}
//启停应用url状态
export const switchOpenurl = (data)=>{
    return service({
        url: `${BASE_URL}/appspace/app/openurl/status`,
        method: 'put',
        data
    })
}


//获取智能体openurl信息
export const getOpenurlInfo = (suffix,config={})=>{
    return service({
        url: `/service/url/openurl/v1/agent/${suffix}`,
        method: 'get',
        ...config,
        isOpenUrl:true
    })
}
//智能体openurl创建智能体对话
export const openurlConversation = (data,suffix,config={})=>{
    return service({
        url: `/service/url/openurl/v1/agent/${suffix}/conversation`,
        method: 'post',
        data,
        ...config,
        isOpenUrl:true
    })
}
//删除智能体openurl创建智能体对话
export const delOpenurlConversation = (data,suffix,config={})=>{
    return service({
        url: `/service/url/openurl/v1/agent/${suffix}/conversation`,
        method: 'delete',
        data,
        ...config,
        isOpenUrl:true
    })
}
//智能体openurl详情历史列表
export const OpenurlConverHistory = (data,suffix,config={})=>{
    return service({
        url: `/service/url/openurl/v1/agent/${suffix}/conversation/detail`,
        method: 'get',
        params:data,
        ...config,
        isOpenUrl:true
    })
}
//智能体openurl对话列表
export const OpenurlConverList = (suffix,config={})=>{
    return service({
        url: `/service/url/openurl/v1/agent/${suffix}/conversation/list`,
        method: 'get',
        ...config,
        isOpenUrl:true
    })
}
//智能体openurl流式对话
export const OpenurlStream = (data,suffix,config={})=>{
    return service({
        url: `/service/url/openurl/v1/agent/${suffix}/stream`,
        method: 'post',
        data,
        ...config,
        isOpenUrl:true
    })
}
//更新博查rerank模型
export const updateRerank = (data)=>{
    return service({
        url: `${BASE_URL}/assistant/tool/config`,
        method: 'put',
        data
    })
}