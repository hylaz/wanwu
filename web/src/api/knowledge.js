import service from "@/utils/request";
import {KNOWLEDGE_API} from "@/utils/requestConstants";
//新增文档知识分类
export const createDoc = (data) => {
    return service({
        url: `${KNOWLEDGE_API}/ux/doccategory`,
        method: "post",
        data: data,
    });
};

//修改文档知识分类
export const editDoc = (data) => {
    return service({
        url: `${KNOWLEDGE_API}/ux/doccategory`,
        method: "put",
        data: data,
    });
};

//删除文档知识分类
export const removeDoc = (data) => {
    return service({
        url: `${KNOWLEDGE_API}/ux/doccategory`,
        method: "delete",
        data: data,
    });
};

//获取文档列表
export const getList = (data) => {
    return service({
        url: `${KNOWLEDGE_API}/ux/doc/list`,
        method: "post",
        data: data,
    });
};
//编辑文档
export const modifyDoc = (data) => {
    return service({
        url: `${KNOWLEDGE_API}/ux/doc`,
        method: "put",
        data: data,
    });
};
//删除文档
export const deleteDoc = (data) => {
    return service({
        url: `${KNOWLEDGE_API}/ux/doc`,
        method: "delete",
        data: data,
    });
};
//上传文档
export const importDoc = (data,source) => {
    return service({
        url: `${KNOWLEDGE_API}/ux/doc/import`,
        method: "post",
        cancleToken:source,
        data: data,
        headers: {"Content-Type": "multipart/form-data"}
    });
};
//保存上传文档
export const saveImportDoc = (data) => {
    return service({
        url: `${KNOWLEDGE_API}/ux/doc/save`,
        method: "put",
        data: data,
    });
};
//获取文档下载链接
export const getDocLink = (id) => {
    return service({
        url: `${KNOWLEDGE_API}/ux/doc/download_url?id=${id}`,
        method: "get"
    });
};
//命中测试
export const test = (data) => {
    return service({
        url: `${KNOWLEDGE_API}/ux/chunk/evaluate`,
        method: "post",
        data: data
    });
}
//上传文件删除无效数据
export const deleteInvalid = (data) => {
    return service({
        url: `${KNOWLEDGE_API}/ux/doc/invalid`,
        method: "delete",
        data: data
    });
}
//从url上传
export const setUploadURL = (data)=>{
    return service({
        url: `${KNOWLEDGE_API}/ux/doc/importUrl`,
        method: 'post',
        data
    })
};
// 解析url
export const analyzeURL = (data)=>{
    return service({
        url: `${KNOWLEDGE_API}/ux/doc/analysisUrl`,
        method: 'post',
        data
    })
};

// 查看分段结果列表
export const getContentList = (data,config)=>{
    return service({
        url: `${KNOWLEDGE_API}/ux/doc/fileSplit`,
        method: 'post',
        data: data,
        // config
    })
};

// 查看分段结果列表
export const setStatus = (data,config)=>{
    return service({
        url: `${KNOWLEDGE_API}/ux/doc/updateFileStatus`,
        method: 'post',
        data: data,
        // config
    })
};

// url上传批量
export const batchurl = (data,source)=>{
    return service({
        url: `${KNOWLEDGE_API}/ux/doc/analysisBatchUrl`,
        method: 'post',
        data: data,
        cancelToken: source,
        headers: {"Content-Type": "multipart/form-data"}
    })
};
export const batchUrlTaskStatus = (data)=>{
    return service({
        url: `${KNOWLEDGE_API}/ux/doc/batchUrlTaskStatus`,
        method: 'get',
        params: data
    })
};
export const importBatchUrl = (data)=>{
    return service({
        url: `${KNOWLEDGE_API}/ux/doc/importBatchUrl`,
        method: 'get',
        params: data
    })
};
export const BatchUrlDemo = ()=>{
    return service({
        url: `${KNOWLEDGE_API}/ux/doc/downloadDemoFile`,
        method: 'get'
    })
};


//new 获取知识库列表
import {USER_API} from "@/utils/requestConstants"
export const getKnowledgeList = (data)=>{
    return service({
        url: `${USER_API}/knowledge/select`,
        method: 'post',
        data
    })
};
// export const getKnowledgeItem = (params)=>{
//     return service({
//         url: `${USER_API}/knowledge`,
//         method: 'get',
//         params
//     })
// };
export const delKnowledgeItem = (data)=>{
    return service({
        url: `${USER_API}/knowledge`,
        method: 'delete',
        data
    })
};
export const createKnowledgeItem = (data)=>{
    return service({
        url: `${USER_API}/knowledge`,
        method: 'post',
        data
    })
};
export const editKnowledgeItem = (data)=>{
    return service({
        url: `${USER_API}/knowledge`,
        method: 'put',
        data
    })
};
export const getDocList = (params)=>{
    return service({
        url: `${USER_API}/knowledge/doc/list`,
        method: 'get',
        params
    })
};
export const delDocItem = (data)=>{
    return service({
        url: `${USER_API}/knowledge/doc`,
        method: 'delete',
        data
    })
};
// 上传文件提示接口
export const uploadFileTips = (params)=>{
    return service({
        url: `${USER_API}/knowledge/doc/import/tip`,
        method: 'get',
        params
    })
};
export const getSectionList = (params)=>{
    return service({
        url: `${USER_API}/knowledge/doc/segment/list`,
        method: 'get',
        params
    })
};
//更新文档切片标签
export const sectionLabels = (data)=>{
    return service({
        url: `${USER_API}/knowledge/doc/segment/labels`,
        method: 'post',
        data
    })
};
export const setSectionStatus = (data)=>{
    return service({
        url: `${USER_API}/knowledge/doc/segment/status/update`,
        method: 'post',
        data
    })
};

export const setAnalysis = (data)=>{
    return service({
        url: `${USER_API}/knowledge/doc/url/analysis`,
        method: 'post',
        data
    })
};

export const docImport = (data)=>{
    return service({
        url: `${USER_API}/knowledge/doc/import`,
        method: 'post',
        data
    })
};

//删除知识库标签
export const delTag = (data)=>{
    return service({
        url: `${USER_API}/knowledge/tag`,
        method: 'delete',
        data
    })
};
//查询知识库标签列表
export const tagList = (params)=>{
    return service({
        url: `${USER_API}/knowledge/tag`,
        method: 'get',
        params
    })
};
//创建知识库标签
export const createTag = (data)=>{
    return service({
        url: `${USER_API}/knowledge/tag`,
        method: 'post',
        data
    })
};
//修改知识库标签
export const editTag = (data)=>{
    return service({
        url: `${USER_API}/knowledge/tag`,
        method: 'put',
        data
    })
};
//绑定修改知识库标签
export const bindTag = (data)=>{
    return service({
        url: `${USER_API}/knowledge/tag/bind`,
        method: 'post',
        data
    })
};

//查询标签绑定知识库数量
export const bindTagCount = (params)=>{
    return service({
        url: `${USER_API}/knowledge/tag/bind/count`,
        method: 'get',
        params
    })
};

//命中测试接口
export const hitTest = (data)=>{
    return service({
        url: `${USER_API}/knowledge/hit`,
        method: 'post',
        data
    })
};
export const ocrSelectList = ()=>{
    return service({
        url: `${USER_API}/model/select/ocr`,
        method: 'get',
    })
};
export const updateDocMeta = (data)=>{
    return service({
        url: `${USER_API}/knowledge/doc/meta`,
        method: 'post',
        data
    })
};
export const delSplitter = (data)=>{
    return service({
        url: `${USER_API}/knowledge/splitter`,
        method: 'delete',
        data
    })
};
export const getSplitter = (data)=>{
    return service({
        url: `${USER_API}/knowledge/splitter`,
        method: 'get',
        params:data
    })
};
export const createSplitter = (data)=>{
    return service({
        url: `${USER_API}/knowledge/splitter`,
        method: 'post',
        data
    })
};
export const editSplitter = (data)=>{
    return service({
        url: `${USER_API}/knowledge/splitter`,
        method: 'put',
        data
    })
};
export const createSegment = (data)=>{
    return service({
        url: `${USER_API}/knowledge/doc/segment/create`,
        method: 'post',
        data
    })
};
export const createBatchSegment = (data)=>{
    return service({
        url: `${USER_API}/knowledge/doc/segment/batch/create`,
        method: 'post',
        data
    })
};
export const delSegment = (data)=>{
    return service({
        url: `${USER_API}/knowledge/doc/segment/delete`,
        method: 'delete',
        data
    })
};
export const editSegment = (data)=>{
    return service({
        url: `${USER_API}/knowledge/doc/segment/update`,
        method: 'post',
        data
    })
};
export const metaSelect = (params)=>{
    return service({
        url: `${USER_API}/knowledge/meta/select`,
        method: 'get',
        params
    })
};
export const parserSelect = ()=>{
    return service({
        url: `${USER_API}/model/select/pdf-parser`,
        method: 'get'
    })
};
export const getSegmentChild = (params)=>{
    return service({
        url: `${USER_API}/knowledge/doc/segment/child/list`,
        method: 'get',
        params
    })
};

export const createSegmentChild = (data)=>{
    return service({
        url: `${USER_API}/knowledge/doc/segment/child/create`,
        method: 'post',
        data
    })
};
export const delSegmentChild = (data)=>{
    return service({
        url: `${USER_API}/knowledge/doc/segment/child/delete`,
        method: 'delete',
        data
    })
};
export const updateSegmentChild = (data)=>{
    return service({
        url: `${USER_API}/knowledge/doc/segment/child/update`,
        method: 'post',
        data
    })
};
// 获取知识库组织列表
export const getOrgList = (data)=>{
    return service({
        url: `${USER_API}/knowledge/org`,
        method: 'get',
        params:data
    })
};
// 获取知识库组织列表
export const getOrgUser = (data)=>{
    return service({
        url: `${USER_API}/knowledge/user/no/permit`,
        method: 'get',
        params:data
    })
};
// 获取知识库用户权限列表
export const getUserPower = (data)=>{
    return service({
        url: `${USER_API}/knowledge/user`,
        method: 'get',
        params:data
    })
};
// 新增知识库用户权限
export const addUserPower = (data)=>{
    return service({
        url: `${USER_API}/knowledge/user/add`,
        method: 'post',
        data
    })
};
// 转让知识库管理权限
export const transferUserPower = (data)=>{
    return service({
        url: `${USER_API}/knowledge/user/admin/transfer`,
        method: 'post',
        data
    })
};
// 修改知识库用户权限
export const editUserPower = (data)=>{
    return service({
        url: `${USER_API}/knowledge/user/edit`,
        method: 'post',
        data
    })
};
// 删除知识库用户权限
export const delUserPower = (data)=>{
    return service({
        url: `${USER_API}/knowledge/user/delete`,
        method: 'delete',
        data
    })
};
//更新文档元数据
export const updateMetaData = (data)=>{
    return service({
        url: `${USER_API}/knowledge/meta/value/update`,
        method: 'post',
        data
    })
};

//获取文档元数据列表
export const getDocMetaList = (data)=>{
    return service({
        url: `${USER_API}/knowledge/meta/value/list`,
        method: 'post',
        data
    })
};

//获取知识图谱详情
export const getGraphDetail = (data)=>{
    return service({
        url: `${BASE_URL}/knowledge/graph`,
        method: 'get',
        params:data
    })
};

//单条新增社区报告
export const createCommunityReport = (data)=>{
    return service({
        url: `${USER_API}/knowledge/report/add`,
        method: 'post',
        data
    })
};
//批量新增社区报告
export const createBatchCommunityReport = (data)=>{
    return service({
        url: `${USER_API}/knowledge/report/batch/add`,
        method: 'post',
        data
    })
};
//删除社区报告
export const delCommunityReport = (data)=>{
    return service({
        url: `${USER_API}/knowledge/report/delete`,
        method: 'delete',
        data
    })
};
//生成社区报告
export const generateCommunityReport = (data)=>{
    return service({
        url: `${USER_API}/knowledge/report/generate`,
        method: 'post',
        data
    })
};
//获取社区报告
export const getCommunityReportList = (data)=>{
    return service({
        url: `${USER_API}/knowledge/report/list`,
        method: 'get',
        params:data
    })
};
//编辑社区报告
export const editCommunityReportList = (data)=>{
    return service({
        url: `${USER_API}/knowledge/report/update`,
        method: 'post',
        data
    })
};