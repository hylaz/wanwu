import service from '@/utils/request';
import { USER_API } from '@/utils/requestConstants';

// 获取列表
export const fetchModelList = params => {
  return service({
    url: `${USER_API}/model/list`,
    method: 'get',
    params,
  });
};

// 获取单个模型
export const getModelDetail = params => {
  return service({
    url: `${USER_API}/model`,
    method: 'get',
    params,
  });
};

// 创建
export const addModel = data => {
  return service({
    url: `${USER_API}/model`,
    method: 'post',
    data,
  });
};
// 编辑
export const editModel = data => {
  return service({
    url: `${USER_API}/model`,
    method: 'put',
    data,
  });
};
// 删除
export const deleteModel = data => {
  return service({
    url: `${USER_API}/model`,
    method: 'delete',
    data,
  });
};
// 修改状态
export const changeModelStatus = data => {
  return service({
    url: `${USER_API}/model/status`,
    method: 'put',
    data,
  });
};

//获取embedding列表
export const getEmbeddingList = params => {
  return service({
    url: `${USER_API}/model/select/embedding`,
    method: 'get',
    params,
  });
};

//获取rerank模型列表
export const getRerankList = () => {
  return service({
    url: `${USER_API}/model/select/rerank`,
    method: 'get',
  });
};

//获取下来选择模型列表
export const selectModelList = () => {
  return service({
    url: `${USER_API}/model/select/llm`,
    method: 'get',
  });
};
