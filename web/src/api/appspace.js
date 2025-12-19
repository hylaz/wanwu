import request from '@/utils/request';
import { USER_API } from '@/utils/requestConstants';

// 生成apikey
export const createApiKey = data => {
  return request({
    url: `${USER_API}/appspace/app/key`,
    method: 'post',
    data,
  });
};
// 删除apikey
export const delApiKey = data => {
  return request({
    url: `${USER_API}/appspace/app/key`,
    method: 'delete',
    data,
  });
};
// 获取apikey列表
export const getApiKeyList = params => {
  return request({
    url: `${USER_API}/appspace/app/key/list`,
    method: 'get',
    params,
  });
};
// 获取apikey根地址
export const getApiKeyRoot = params => {
  return request({
    url: `${USER_API}/appspace/app/url`,
    method: 'get',
    params,
  });
};

// 获取智能体/文本问答/工作流列表
export const getAppSpaceList = params => {
  return request({
    url: `${USER_API}/appspace/app/list`,
    method: 'get',
    params,
  });
};

//发布app
export const appPublish = data => {
  return request({
    url: `${USER_API}/appspace/app/publish`,
    method: 'post',
    data,
  });
};

// 取消发布app
export const appCancelPublish = data => {
  return request({
    url: `${USER_API}/appspace/app/publish`,
    method: 'delete',
    data,
  });
};

//统一删除工作室应用接口
export const deleteApp = data => {
  return request({
    url: `${USER_API}/appspace/app`,
    method: 'delete',
    data,
  });
};

//获取应用最新版本信息
export const getAppLatestVersion = params => {
  return request({
    url: `${USER_API}/appspace/app/version`,
    method: 'get',
    params,
  });
};

//更新应用版本信息
export const updateAppVersion = data => {
  return request({
    url: `${USER_API}/appspace/app/version`,
    method: 'put',
    data,
  });
};

//获取应用版本列表
export const getAppVersionList = params => {
  return request({
    url: `${USER_API}/appspace/app/version/list`,
    method: 'get',
    params,
  });
};

//回滚应用到指定版本
export const rollbackAppVersion = data => {
  return request({
    url: `${USER_API}/appspace/app/version/rollback`,
    method: 'post',
    data,
  });
};

//智能体模版
export const agentTemplateList = params => {
  return request({
    url: `${USER_API}/assistant/template/list`,
    method: 'get',
    params,
  });
};
//复制智能体
export const copyAgentTemplate = data => {
  return request({
    url: `${USER_API}/assistant/template`,
    method: 'post',
    data,
  });
};
//智能体模版详情
export const agentTemplateDetail = params => {
  return request({
    url: `${USER_API}/assistant/template`,
    method: 'get',
    params,
  });
};
//复制文本问答应用
export const copyTextQues = data => {
  return request({
    url: `${USER_API}/appspace/rag/copy`,
    method: 'post',
    data,
  });
};
//复制智能体应用
export const copyAgentApp = data => {
  return request({
    url: `${USER_API}/assistant/copy`,
    method: 'post',
    data,
  });
};
