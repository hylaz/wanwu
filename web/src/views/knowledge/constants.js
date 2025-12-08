// 初始化标志常量
export const INITIAL = -1;

// 通用状态常量
export const STATUS_PENDING = 0;
export const STATUS_PROCESSING = 1;
export const STATUS_FINISHED = 2;
export const STATUS_FAILED = 3;

// 权限类型常量
export const POWER_TYPE_READ = 0;
export const POWER_TYPE_EDIT = 10;
export const POWER_TYPE_ADMIN = 20;
export const POWER_TYPE_SYSTEM_ADMIN = 30;

// 分段类型常量
export const SEGMENT_TYPE_AUTO = '0';
export const SEGMENT_TYPE_CUSTOM = '1';
export const SEGMENT_TYPE_COMMON = '0';
export const SEGMENT_TYPE_PARENTSON = '1';

// 分析器类型常量
export const ANALYZER_TYPE_TEXT = 'text';
export const ANALYZER_TYPE_OCR = 'ocr';
export const ANALYZER_TYPE_MODEL = 'model';

// 报告状态常量
export const REPORT_STATUS_PENDING = 0;
export const REPORT_STATUS_GENERATING = 1;
export const REPORT_STATUS_GENERATED = 2;
export const REPORT_STATUS_GENERATION_FAILED = 3;

// 知识库文档解析状态常量
export const KNOWLEDGE_STATUS_UPLOADED = -2;
export const KNOWLEDGE_STATUS_ALL = -1;
export const KNOWLEDGE_STATUS_PENDING_PROCESSING = 0;
export const KNOWLEDGE_STATUS_FINISH = 1;
export const KNOWLEDGE_STATUS_CHECKING = 2;
export const KNOWLEDGE_STATUS_ANALYSING = 3;
export const KNOWLEDGE_STATUS_CHECK_FAIL = 4;
export const KNOWLEDGE_STATUS_FAIL = 5;

// 问答库文档解析状态常量
export const QA_STATUS_ALL = -1;
export const QA_STATUS_PENDING = 0;
export const QA_STATUS_PROCESSING = 1;
export const QA_STATUS_FINISHED = 2;
export const QA_STATUS_FAILED = 3;
