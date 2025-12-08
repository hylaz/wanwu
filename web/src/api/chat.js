import request from '@/utils/request';
import { MODEL_API, SERVICE_API } from '@/utils/requestConstants';

export const createApp = data => {
  return request({
    url: `${MODEL_API}/assistant/create`,
    method: 'post',
    data,
  });
};
export const updateApp = data => {
  return request({
    url: `${MODEL_API}/assistant/update`,
    method: 'put',
    data,
  });
};
export const getAppDetail = data => {
  return request({
    url: `${MODEL_API}/assistant/info`,
    method: 'get',
    params: data,
  });
};
export const deleteApp = data => {
  return request({
    url: `${MODEL_API}/assistant/delete`,
    method: 'delete',
    data,
  });
};
export const publishApp = data => {
  return request({
    url: `${MODEL_API}/assistant/publish`,
    method: 'post',
    data,
  });
};
export const getAppDraftList = data => {
  return request({
    url: `${MODEL_API}/assistant/draft_list`,
    method: 'get',
    params: data,
  });
};
export const getAppMoreList = data => {
  return request({
    url: `${MODEL_API}/assistant/more_list`,
    method: 'get',
    params: data,
  });
};
export const getMyAppList = data => {
  return request({
    url: `${MODEL_API}/assistant/list`,
    method: 'get',
    params: data,
  });
};
//头像上传
export const fileUpload = (data, config) => {
  return request({
    url: `${SERVICE_API}/model/expansion/file/batch/upload`,
    method: 'post',
    data,
    config,
  });
};
//知识增强文件上传
export const knowledgeFileUpload = (data, config) => {
  return request({
    url: `${MODEL_API}/assistant/knowledge/file/upload`,
    method: 'post',
    data,
    config,
  });
};
//查询已上传文件列表
export const getKnowledgeFileList = data => {
  return request({
    url: `${MODEL_API}/assistant/knowledge/file/list`,
    method: 'get',
    params: data,
  });
};
export const deleteKnowledgeFile = data => {
  return request({
    url: `${MODEL_API}/assistant/knowledge/file/delete`,
    method: 'delete',
    data,
  });
};
//常用应用
export const getRecentApp = data => {
  return request({
    url: `${MODEL_API}/assistant/common/list`,
    method: 'get',
    params: data,
  });
};
//删除常用应用
export const deleteRecentApp = data => {
  return request({
    url: `${MODEL_API}/assistant/common/delete`,
    method: 'delete',
    data,
  });
};

//对话列表
export const getConversationList = data => {
  return request({
    url: `${MODEL_API}/assistant/conversation/list`,
    method: 'get',
    params: data,
  });
};
//创建对话
export const createConversation = data => {
  return request({
    url: `${MODEL_API}/assistant/conversation/create`,
    method: 'post',
    data,
  });
};
//删除对话
export const deleteConversation = data => {
  return request({
    url: `${MODEL_API}/assistant/conversation/delete`,
    method: 'delete',
    data,
  });
};
//对话详情
export const getConversationDetail = data => {
  return request({
    url: `${MODEL_API}/assistant/conversation/detail`,
    method: 'get',
    params: data,
  });
};
/*----元景------*/
//对话列表
export const getConversationListCUBM = data => {
  return request({
    url: `${MODEL_API}/chatllm/conversation/list`,
    method: 'get',
    params: data,
  });
};
//创建对话
export const createConversationCUBM = data => {
  return request({
    url: `${MODEL_API}/chatllm/conversation/create`,
    method: 'post',
    data,
  });
};
//删除对话
export const deleteConversationCUBM = data => {
  return request({
    url: `${MODEL_API}/chatllm/conversation/delete`,
    method: 'delete',
    data,
  });
};
//对话详情
export const getConversationDetailCUBM = data => {
  return request({
    url: `${MODEL_API}/chatllm/conversation/detail`,
    method: 'get',
    params: data,
  });
};
//批量文件上传
export const batchUpload = (data, config) => {
  return request({
    url: `${MODEL_API}/file/batch/upload`,
    method: 'post',
    data,
    config,
  });
};
// app接入
export const linkAPP = data => {
  return request({
    url: `${MODEL_API}/assistant/app/publish`,
    method: 'post',
    data,
  });
};

//推荐智能体列表
export const recommendList = data => {
  return request({
    url: `${MODEL_API}/assistant/recommend/list`,
    method: 'get',
    params: data,
  });
};
//标记推荐智能体
export const recommendMark = data => {
  return request({
    url: `${MODEL_API}/assistant/recommend/update`,
    method: 'put',
    data,
  });
};

//上传文件确认路径
export const confirmPath = data => {
  return request({
    url: `${MODEL_API}/file/confirmPath`,
    method: 'post',
    data,
  });
};
