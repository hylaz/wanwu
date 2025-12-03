import { i18n } from "@/lang";

export const FAT_SON_BLOCK = [
  {
    title: i18n.t("knowledgeManage.config.parentBlock"),
    level: "parent",
    key: "splitter",
    splitter: "splitter",
    maxSplitter: "maxSplitter",
    splitterProp: "docSegment.splitter",
    maxSplitterProp: "docSegment.maxSplitter",
    maxSplitterNum: 4000,
  },
  {
    title: i18n.t("knowledgeManage.config.sonBlock"),
    level: "son",
    key: "subSplitter",
    splitter: "subSplitter",
    maxSplitter: "subMaxSplitter",
    splitterProp: "docSegment.subSplitter",
    maxSplitterProp: "docSegment.subMaxSplitter",
    maxSplitterNum: 4000,
  },
];
export const SEGMENT_COMMON_LIST = [
  {
    label: "0",
    text: i18n.t("knowledgeManage.config.autoChunk"),
    desc: i18n.t("knowledgeManage.config.autoChunkDesc"),
  },
  {
    label: "1",
    text: i18n.t("knowledgeManage.config.customChunk"),
    desc: i18n.t("knowledgeManage.config.customChunkDesc"),
  },
];
export const SEGMENT_LIST = [
  {
    label: "0",
    img: "setting-gear.png",
    text: i18n.t("knowledgeManage.config.commonSegment"),
    desc: i18n.t("knowledgeManage.config.commonSegmentDesc"),
  },
  {
    label: "1",
    img: "setting-effect.png",
    text: i18n.t("knowledgeManage.config.parentSonSegment"),
    desc: i18n.t("knowledgeManage.config.parentSonSegmentDesc"),
  },
];
export const DOC_ANALYZER_LIST = [
  {
    label: "text",
    text: i18n.t("knowledgeManage.config.textExtraction"),
    desc: i18n.t("knowledgeManage.config.textExtractionDesc"),
  },
  {
    label: "ocr",
    text: i18n.t("knowledgeManage.OCRAnalysis"),
    desc: i18n.t("knowledgeManage.config.OCRAnalysisDesc"),
  },
  {
    label: "model",
    text: i18n.t("knowledgeManage.config.modelAnalysis"),
    desc: i18n.t("knowledgeManage.config.modelAnalysisDesc"),
  },
];
export const MODEL_TYPE_TIP = {
  ocr: {
    label: i18n.t("knowledgeManage.config.OCRModel"),
    desc: i18n.t("knowledgeManage.config.OCRModelDesc"),
  },
  model: {
    label: i18n.t("knowledgeManage.config.documentAnalysis"),
    desc: i18n.t("knowledgeManage.config.documentAnalysisDesc"),
  },
};
export const POWER_TYPE = {
  0: i18n.t("knowledgeManage.config.read"),
  10: i18n.t("knowledgeManage.config.edit"),
  20: i18n.t("knowledgeManage.config.admin"),
  30: i18n.t("knowledgeManage.config.systemAdmin"),
};
export const KNOWLEDGE_GRAPH_TIPS = [
  {
    title: i18n.t("knowledgeManage.config.functionDescription"),
    content: i18n.t("knowledgeManage.config.functionDescriptionContent"),
  },
  {
    title: i18n.t("knowledgeManage.config.sceneDescription"),
    content: i18n.t("knowledgeManage.config.sceneDescriptionContent"),
  },
  {
    title: i18n.t("knowledgeManage.config.attentionDescription"),
    content: i18n.t("knowledgeManage.config.attentionDescriptionContent"),
  },
];
export const COMMUNITY_REPORT_STATUS = {
  0: "-",
  1: i18n.t("knowledgeManage.config.generating"),
  2: i18n.t("knowledgeManage.config.generated"),
  3: i18n.t("knowledgeManage.config.generationFailed"),
};
export const KNOWLEDGE_GRAPH_STATUS = {
  0: i18n.t("knowledgeManage.config.pending"),
  1: i18n.t("knowledgeManage.config.processing"),
  2: i18n.t("knowledgeManage.config.finished"),
  3: i18n.t("knowledgeManage.config.failed"),
};
export const COMMUNITY_IMPORT_STATUS = {
  0: i18n.t("knowledgeManage.communityReport.taskPending"),
  1: i18n.t("knowledgeManage.communityReport.taskProcessing"),
  2: i18n.t("knowledgeManage.communityReport.taskFinished"),
  3: i18n.t("knowledgeManage.communityReport.taskFailed"),
};

export const DROPDOWN_GROUPS = [
  {
    label: i18n.t("common.button.export"),
    icon: "el-icon-arrow-down",
    items: [
      { command: "exportData", label: i18n.t("knowledgeManage.qaDatabase.exportData") },
      { command: "exportRecord", label: i18n.t("knowledgeManage.qaDatabase.exportRecord") },
    ],
  },
  {
    label: i18n.t("common.button.add"),
    icon: "el-icon-arrow-down",
    items: [
      { command: "createQaPair", label: i18n.t("knowledgeManage.qaDatabase.createQaPair") },
      { command: "fileUpload", label: i18n.t("knowledgeManage.fileUpload") },
    ]
  }
]