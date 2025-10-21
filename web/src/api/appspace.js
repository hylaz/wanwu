import request from "@/utils/request";
const BASE_URL = '/user/api/v1'

// 生成apikey
export const createApiKey = (data)=>{
    return request({
        url: `${BASE_URL}/appspace/app/key`,
        method: 'post',
        data
    })
};
// 删除apikey
export const delApiKey = (data)=>{
    return request({
        url: `${BASE_URL}/appspace/app/key`,
        method: 'delete',
        data
    })
};
// 获取apikey列表
export const getApiKeyList = (params)=>{
    return request({
        url: `${BASE_URL}/appspace/app/key/list`,
        method: 'get',
        params
    })
};
// 获取apikey根地址
export const getApiKeyRoot = (params)=>{
    return request({
        url: `${BASE_URL}/appspace/app/url`,
        method: 'get',
        params
    })
};

// 获取智能体/文本问答/工作流列表
export const getAppSpaceList = (params)=>{
    return request({
        url: '/user/api/v1/appspace/app/list',
        method: 'get',
        params
    })
}

//发布app
export const appPublish = (data)=>{
    return request({
        url: `${BASE_URL}/appspace/app/publish`,
        method: 'post',
        data
    })
};

// 取消发布app
export const appCancelPublish = (data)=>{
    return request({
        url: `${BASE_URL}/appspace/app/publish`,
        method: 'delete',
        data
    })
};

//统一删除工作室应用接口
export const deleteApp = (data)=>{
    return request({
        url: `${BASE_URL}/appspace/app`,
        method: 'delete',
        data
    })
};

//智能体模版
export const agnetTemplateList = (params)=>{
    return request({
        url: `${BASE_URL}/assistant/template/list`,
        method: 'get',
        params
    })
};
//复制智能体
export const copyAgnetTemplate = (data)=>{
    return request({
        url: `${BASE_URL}/assistant/template`,
        method: 'post',
        data
    })
};
//智能体模版详情
export const agnetTemplateDetail = (params)=>{
    return request({
        url: `${BASE_URL}/assistant/template`,
        method: 'get',
        params
    })
};
//复制文本问答应用
export const copyTextQues = (data)=>{
    return request({
        url: `${BASE_URL}/appspace/rag/copy`,
        method: 'post',
        data
    })
};

