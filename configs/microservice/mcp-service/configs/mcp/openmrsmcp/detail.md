
# MCP 服务模块文档

## 核心功能模块 (metrics_processor.py)

### 功能说明
提供基于OpenMRS数据指标计算、系统信息获取和大数据处理的核心功能

#### 数据指标计算和报表生成主要工具函数
1. **calculate_metrics** - 多维度指标计算引擎
   - 输入: raw_data(数据集), metric_config(指标配置)
   - 输出: 包含计算状态、指标结果、数据问题的字典
   - 特性: 支持50+技术指标，含数据质量检查与超时控制

2. **process_large_data** - 大数据集分块处理
   - 特性: 分块处理避免内存溢出，支持并行计算

3. **generate_report** - 生成可视化报表
   - 支持格式: markdown/html
   - 输入: kpi_data(指标数据字典)


## 患者管理模块 (patient_management.py)

### 功能说明
提供基于OpenMRS的患者CRUD操作的完整生命周期管理，提供患者兵力信息管理、交互式多条件检索患者工具和报表与患者匹配去重工具。

#### 主要工具函数
1. **患者查询**:
   - search_patient_by_identifier: 通过标识符查询
   - search_patient_by_name: 通过姓名查询
   - search_patient_by_uuid: 通过UUID查询

2. **患者管理**:
   - create_patient: 创建新患者记录
     - 必填字段: 标识符、姓名、性别、出生日期等
   - update_patient: 更新患者信息
   - manage_patient: 删除/恢复患者(软删除)

3. **错误测试**:
   - test_search_errors: 测试错误场景

### 运行方式
```bash
python patient_management.py
```

## 国际医疗数据标准交换工具(便于模块 (fhir_adapter.py)

### 功能说明
实现基于国际医疗数据标注，特别是FHIR标准与OpenMRS系统的数据转换

#### 主要工具函数
1. **import_fhir_patient**:
   - 功能: 将FHIR患者数据转换为OpenMRS格式
   - 输入: fhir_data(FHIR格式数据)
   - 输出: OpenMRS兼容的患者数据

### 运行方式
```bash
python fhir_adapter.py
```

## 药品管理模块 (drug_management.py)

### 功能说明
提供药品创建和概念查询、校验功能

#### 主要功能
1. **create_drug**:
   - 创建药品记录
   - 必需参数: 药品名称、描述、概念UUID等

2. **get_concept_uuid**:
   - 根据概念名称查询UUID

3. **系统验证**:
   - validate_openmrs_connection: 验证服务器连接

# 门诊预约排班系统

## 模块说明

### 1. 药品管理模块 (drug_management.py)

#### 功能说明
提供药品创建和概念查询、校验功能

#### 主要功能
- **create_drug**:
  - 创建药品记录
  - 必需参数: 药品名称、描述、概念UUID等
- **get_concept_uuid**:
  - 根据概念名称查询UUID
- **系统验证**:
  - validate_openmrs_connection: 验证服务器连接

### 4. 预约排班模块 (appointment_scheduling.py)

#### 功能说明
提供门诊医生排班管理和患者预约功能

#### 主要功能
- **create_schedule**:
  - 创建医生排班记录
  - 必需参数: 医生ID、科室、排班日期、时间段、最大预约数等
- **query_available_slots**:
  - 查询指定日期/科室的可用预约时段
  - 可选参数: 科室、日期范围、医生ID等
- **book_appointment**:
  - 患者预约功能
  - 必需参数: 患者ID、排班ID、预约时间等
- **cancel_appointment**:
  - 取消预约功能
  - 必需参数: 预约ID
- **系统验证**:
  - validate_system_connection: 验证预约系统连接状态

### 通用设置
- 所有服务均支持SSE传输协议
- 默认端口配置:
  - 核心功能: 9062
  - 患者管理: 8005
  - FHIR适配器: 9062
  - 预约系统: 9070

### 排班系统专用配置
- 排班时间粒度: 30分钟
- 最大预约提前天数: 30天

## 通用配置
- 所有服务均支持SSE传输协议
- 默认端口:
  - 核心功能: 9062
  - 患者管理: 8005
  - FHIR适配器: 9062

## 错误处理
所有模块均采用统一错误响应格式:
```json
{
    "success": bool,
    "error": {
        "code": int,
        "msg": str,
        "solution": str
    }
}
```

## 测试建议
各模块均包含测试用例，可通过取消注释`asyncio.run(test_*)`执行本地测试
```
