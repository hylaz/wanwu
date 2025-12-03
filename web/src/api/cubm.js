import request from "@/utils/request";
import {MODEL_API, DATACENTER_API} from "@/utils/requestConstants";

/*----元景------*/
//对话列表
export const getConversationList = (data) => {
  return request({
    url: `${MODEL_API}/chatllm/conversation/list`,
    method: 'get',
    params: data
  })
};
//创建对话
export const createConversation = (data) => {
  return request({
    url: `${MODEL_API}/chatllm/conversation/create`,
    method: 'post',
    data
  })
};
//删除对话
export const deleteConversation = (data) => {
  return request({
    url: `${MODEL_API}/chatllm/conversation/delete`,
    method: 'delete',
    data
  })
};
//对话详情
export const getConversationDetail = (data) => {
  return request({
    url: `${MODEL_API}/chatllm/conversation/detail`,
    method: 'get',
    params: data
  })
};
export const addAction = (data) => {
  return request({
    url: `${MODEL_API}/assistant/action/create`,
    method: 'post',
    data
  })
};
export const updateAction = (data) => {
  return request({
    url: `${MODEL_API}/assistant/action/update`,
    method: 'put',
    data
  })
};
export const deleteAction = (data) => {
  return request({
    url: `${MODEL_API}/assistant/action/delete`,
    method: 'delete',
    data
  })
};
export const getActionDetail = (data) => {
  return request({
    url: `${MODEL_API}/assistant/action/info`,
    method: 'get',
    params: data
  })
};
export const deleteConversationHistory = (data) => {
  return request({
    url: `${MODEL_API}/assistant/conversation/detail/delete`,
    method: 'delete',
    data
  })
};
//获取模型列表
export const getModelList = (data) => {
  return request({
    url: `${DATACENTER_API}/infer/publish/model/select`,
    method: "get",
    params: data
  });
}

//AI自动生成原生应用
export const autoCreate = (data) => {
  return request({
    url: `${MODEL_API}/assistant/auto/create`,
    method: "post",
    data
  });
}

