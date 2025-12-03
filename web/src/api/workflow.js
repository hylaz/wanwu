import request from "@/utils/request";
import { SERVICE_API, USER_API, WORKFLOW_API } from "@/utils/requestConstants"

export const getWorkFlowParams = (params) => {
    return request({
        url: `${WORKFLOW_API}/workflow/parameter`,
        method: "get",
        params,
    });
};
export const useWorkFlow = (data)=>{
    return request({
        url: `${WORKFLOW_API}/api/workflow/use`,
        method: 'post',
        data
    })
};
//应用广场工作流列表
export const getExplorationFlowList = (params)=>{
    return request({
        url: `${USER_API}/exploration/app/list`,
        method: 'get',
        params
    })
};
export const readWorkFlow = (data)=>{
    return request({
        url: `${WORKFLOW_API}/workflow/openapi_schema`,
        method: 'get',
        params:data
    })
};
export const externalUpload = (data, config) => {
    return request({
        url: `${SERVICE_API}/proxy/file/upload`,
        method: "post",
        data,
        config,
        isHandleRes: false
    });
};

// 图片上传
export const uploadFile = (data) => {
    return request({
        url: `/api/bot/upload_file`,
        method: "post",
        data
    })
}

// 创建
export const createWorkFlow = (data, appType)=>{
    return request({
        url: `${USER_API}/appspace/${appType || 'workflow'}`,
        method: 'post',
        data
    })
};

// 复制
export const copyWorkFlow = (data, appType)=>{
    return request({
        url: `${USER_API}/appspace/${appType || 'workflow'}/copy`,
        method: 'post',
        data
    })
};

// 导入
export const importWorkflow = (data, config, appType) => {
    return request({
        url: `${USER_API}/appspace/${appType || 'workflow'}/import`,
        method: 'post',
        data,
        config
    });
};

// 导出
export const exportWorkflow = (params, appType) => {
    return request({
        url: `${USER_API}/appspace/${appType || 'workflow'}/export`,
        method: "get",
        params,
        responseType: 'blob'
    });
};

// 工作流/对话流互转
export const transformWorkflow = (data, appType) => {
    return request({
        url: `${USER_API}/appspace/${appType || 'workflow'}/convert`,
        method: "post",
        data,
    });
};