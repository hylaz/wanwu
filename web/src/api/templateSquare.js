import request from "@/utils/request"
const BASE_URL = '/user/api/v1'

/*---工作流模板---*/
export const getWorkflowTempList = (data)=>{
    return request({
        url: `${BASE_URL}/workflow/template/list`,
        method: 'get',
        params: data
    })
};
export const getWorkflowTempInfo = (data)=>{
    return request({
        url: `${BASE_URL}/workflow/template/detail`,
        method: 'get',
        params: data
    })
};
export const getWorkflowRecommendsList = (data)=>{
    return request({
        url: `${BASE_URL}/workflow/template/recommend`,
        method: 'get',
        params: data
    })
};
export const downloadWorkflow = (params) => {
    return request({
        url: `${BASE_URL}/workflow/template/download`,
        method: "get",
        params,
        responseType: 'blob'
    });
};
export const copyWorkflowTemplate = (data)=>{
    return request({
        url: `${BASE_URL}/workflow/template`,
        method: 'post',
        data
    })
};