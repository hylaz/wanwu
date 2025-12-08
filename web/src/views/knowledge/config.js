import { i18n } from '@/lang';
import {
  STATUS_PENDING,
  STATUS_PROCESSING,
  STATUS_FINISHED,
  STATUS_FAILED,
  POWER_TYPE_READ,
  POWER_TYPE_EDIT,
  POWER_TYPE_ADMIN,
  POWER_TYPE_SYSTEM_ADMIN,
  SEGMENT_TYPE_AUTO,
  SEGMENT_TYPE_CUSTOM,
  SEGMENT_TYPE_COMMON,
  SEGMENT_TYPE_PARENTSON,
  ANALYZER_TYPE_TEXT,
  ANALYZER_TYPE_OCR,
  ANALYZER_TYPE_MODEL,
  REPORT_STATUS_PENDING,
  REPORT_STATUS_GENERATING,
  REPORT_STATUS_GENERATED,
  REPORT_STATUS_GENERATION_FAILED,
  KNOWLEDGE_STATUS_UPLOADED,
  KNOWLEDGE_STATUS_ALL,
  KNOWLEDGE_STATUS_PENDING_PROCESSING,
  KNOWLEDGE_STATUS_FINISH,
  KNOWLEDGE_STATUS_CHECKING,
  KNOWLEDGE_STATUS_ANALYSING,
  KNOWLEDGE_STATUS_CHECK_FAIL,
  KNOWLEDGE_STATUS_FAIL,
  QA_STATUS_ALL,
  QA_STATUS_PENDING,
  QA_STATUS_PROCESSING,
  QA_STATUS_FINISHED,
  QA_STATUS_FAILED,
} from '@/views/knowledge/constants';

export const FAT_SON_BLOCK = [
  {
    title: i18n.t('knowledgeManage.config.parentBlock'),
    level: 'parent',
    key: 'splitter',
    splitter: 'splitter',
    maxSplitter: 'maxSplitter',
    splitterProp: 'docSegment.splitter',
    maxSplitterProp: 'docSegment.maxSplitter',
    maxSplitterNum: 4000,
  },
  {
    title: i18n.t('knowledgeManage.config.sonBlock'),
    level: 'son',
    key: 'subSplitter',
    splitter: 'subSplitter',
    maxSplitter: 'subMaxSplitter',
    splitterProp: 'docSegment.subSplitter',
    maxSplitterProp: 'docSegment.subMaxSplitter',
    maxSplitterNum: 4000,
  },
];
export const SEGMENT_COMMON_LIST = [
  {
    label: SEGMENT_TYPE_AUTO,
    text: i18n.t('knowledgeManage.config.autoChunk'),
    desc: i18n.t('knowledgeManage.config.autoChunkDesc'),
  },
  {
    label: SEGMENT_TYPE_CUSTOM,
    text: i18n.t('knowledgeManage.config.customChunk'),
    desc: i18n.t('knowledgeManage.config.customChunkDesc'),
  },
];
export const SEGMENT_LIST = [
  {
    label: SEGMENT_TYPE_COMMON,
    img: 'setting-gear.png',
    text: i18n.t('knowledgeManage.config.commonSegment'),
    desc: i18n.t('knowledgeManage.config.commonSegmentDesc'),
  },
  {
    label: SEGMENT_TYPE_PARENTSON,
    img: 'setting-effect.png',
    text: i18n.t('knowledgeManage.config.parentSonSegment'),
    desc: i18n.t('knowledgeManage.config.parentSonSegmentDesc'),
  },
];
export const DOC_ANALYZER_LIST = [
  {
    label: ANALYZER_TYPE_TEXT,
    text: i18n.t('knowledgeManage.config.textExtraction'),
    desc: i18n.t('knowledgeManage.config.textExtractionDesc'),
  },
  {
    label: ANALYZER_TYPE_OCR,
    text: i18n.t('knowledgeManage.OCRAnalysis'),
    desc: i18n.t('knowledgeManage.config.OCRAnalysisDesc'),
  },
  {
    label: ANALYZER_TYPE_MODEL,
    text: i18n.t('knowledgeManage.config.modelAnalysis'),
    desc: i18n.t('knowledgeManage.config.modelAnalysisDesc'),
  },
];
export const MODEL_TYPE_TIP = {
  [ANALYZER_TYPE_OCR]: {
    label: i18n.t('knowledgeManage.config.OCRModel'),
    desc: i18n.t('knowledgeManage.config.OCRModelDesc'),
  },
  [ANALYZER_TYPE_MODEL]: {
    label: i18n.t('knowledgeManage.config.documentAnalysis'),
    desc: i18n.t('knowledgeManage.config.documentAnalysisDesc'),
  },
};
export const POWER_TYPE = {
  [POWER_TYPE_READ]: i18n.t('knowledgeManage.config.read'),
  [POWER_TYPE_EDIT]: i18n.t('knowledgeManage.config.edit'),
  [POWER_TYPE_ADMIN]: i18n.t('knowledgeManage.config.admin'),
  [POWER_TYPE_SYSTEM_ADMIN]: i18n.t('knowledgeManage.config.systemAdmin'),
};
export const KNOWLEDGE_GRAPH_TIPS = [
  {
    title: i18n.t('knowledgeManage.config.functionDescription'),
    content: i18n.t('knowledgeManage.config.functionDescriptionContent'),
  },
  {
    title: i18n.t('knowledgeManage.config.sceneDescription'),
    content: i18n.t('knowledgeManage.config.sceneDescriptionContent'),
  },
  {
    title: i18n.t('knowledgeManage.config.attentionDescription'),
    content: i18n.t('knowledgeManage.config.attentionDescriptionContent'),
  },
];
export const COMMUNITY_REPORT_STATUS = {
  [REPORT_STATUS_PENDING]: '-',
  [REPORT_STATUS_GENERATING]: i18n.t('knowledgeManage.config.generating'),
  [REPORT_STATUS_GENERATED]: i18n.t('knowledgeManage.config.generated'),
  [REPORT_STATUS_GENERATION_FAILED]: i18n.t(
    'knowledgeManage.config.generationFailed',
  ),
};
export const KNOWLEDGE_GRAPH_STATUS = {
  [STATUS_PENDING]: i18n.t('knowledgeManage.config.pending'),
  [STATUS_PROCESSING]: i18n.t('knowledgeManage.config.processing'),
  [STATUS_FINISHED]: i18n.t('knowledgeManage.config.finished'),
  [STATUS_FAILED]: i18n.t('knowledgeManage.config.failed'),
};
export const COMMUNITY_IMPORT_STATUS = {
  [STATUS_PENDING]: i18n.t('knowledgeManage.communityReport.taskPending'),
  [STATUS_PROCESSING]: i18n.t('knowledgeManage.communityReport.taskProcessing'),
  [STATUS_FINISHED]: i18n.t('knowledgeManage.communityReport.taskFinished'),
  [STATUS_FAILED]: i18n.t('knowledgeManage.communityReport.taskFailed'),
};

export const DROPDOWN_GROUPS = [
  {
    label: i18n.t('common.button.export'),
    icon: 'el-icon-arrow-down',
    items: [
      {
        command: 'exportData',
        label: i18n.t('knowledgeManage.qaDatabase.exportData'),
      },
      {
        command: 'exportRecord',
        label: i18n.t('knowledgeManage.qaDatabase.exportRecord'),
      },
    ],
  },
  {
    label: i18n.t('common.button.add'),
    icon: 'el-icon-arrow-down',
    items: [
      {
        command: 'createQaPair',
        label: i18n.t('knowledgeManage.qaDatabase.createQaPair'),
      },
      { command: 'fileUpload', label: i18n.t('knowledgeManage.fileUpload') },
    ],
  },
  {
    label: i18n.t('knowledgeManage.hitTest.graph'),
    icon: 'el-icon-arrow-down',
    items: [
      {
        command: 'goKnowledgeGraph',
        label: i18n.t('knowledgeManage.hitTest.graph'),
      },
      {
        command: 'goCommunityReport',
        label: i18n.t('knowledgeManage.hitTest.communityReport'),
      },
    ],
  },
];

export const KNOWLEDGE_STATUS_OPTIONS = [
  { label: i18n.t('knowledgeManage.all'), value: KNOWLEDGE_STATUS_ALL },
  { label: i18n.t('knowledgeManage.finish'), value: KNOWLEDGE_STATUS_FINISH },
  { label: i18n.t('knowledgeManage.fail'), value: KNOWLEDGE_STATUS_FAIL },
  {
    label: i18n.t('knowledgeManage.analysing'),
    value: KNOWLEDGE_STATUS_ANALYSING,
  },
  // {label: i18n.t("knowledgeManage.checkFail"), value: KNOWLEDGE_STATUS_CHECK_FAIL},
  {
    label: i18n.t('knowledgeManage.pendingProcessing'),
    value: KNOWLEDGE_STATUS_PENDING_PROCESSING,
  },
  // { label: i18n.t("knowledgeManage.checking"), value: KNOWLEDGE_STATUS_CHECKING }
];

export const QA_STATUS_OPTIONS = [
  { label: i18n.t('knowledgeManage.all'), value: QA_STATUS_ALL },
  {
    label: i18n.t('knowledgeManage.communityReport.taskPending'),
    value: QA_STATUS_PENDING,
  },
  {
    label: i18n.t('knowledgeManage.communityReport.taskProcessing'),
    value: QA_STATUS_PROCESSING,
  },
  {
    label: i18n.t('knowledgeManage.communityReport.taskFinished'),
    value: QA_STATUS_FINISHED,
  },
  {
    label: i18n.t('knowledgeManage.communityReport.taskFailed'),
    value: QA_STATUS_FAILED,
  },
];
