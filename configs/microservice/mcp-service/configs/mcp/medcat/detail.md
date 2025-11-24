# MedCAT MCP

这是一个基于 `FastMCP` 框架构建的、功能强大的临床文本处理微服务。它深度集成了业界领先的医学文本挖掘工具包 `MedCAT` 和先进的大型语言模型 (LLM) ，提供了一套全面的API，用于从非结构化临床文本中智能地提取、分析和保护有价值的信息。

## 核心功能

1.  **高级实体识别与链接 (Advanced Entity Recognition & Linking)**
    *   **功能**: 自动识别临床文本中的多种医学实体，如疾病、症状、药物、解剖部位、检查和手术等。
    *   **技术**: 主要依赖 `MedCAT` 的概念数据库 (CDB) 和词汇表 (Vocab)，能够将识别出的实体链接到标准医学概念唯一标识符 (CUI)，实现文本的标准化和结构化。
    *   **备用方案**: 在 `MedCAT` 模型不可用时，会切换到基于正则表达式的演示模式，保证基本功能可用。

2.  **智能上下文情景分析 (Intelligent Contextual Analysis)**
    *   **功能**: 精准判断每个医学实体的上下文语义，区分实体是**肯定的** (Affirmed)、**否定的** (Negated)、**历史性的** (Historical)、**假设性的** (Hypothetical) 还是**家族史** (FamilyHistory)。
    *   **技术**:
        *   优先使用 `MetaCAT` 模型进行高效、准确的上下文分析。
        *   若 `MetaCAT` 不可用，则自动调用 **LLM 增强功能**，通过专业的医学Prompt模板进行深度语义理解，实现高精度的上下文判断。
        *   在演示模式下，使用规则匹配作为基础实现。

3.  **深层关系抽取 (Deep Relation Extraction)**
    *   **功能**: 抽取医学实体之间的语义关系三元组，例如 “药物A `治疗` 疾病B”、“症状C `由` 疾病D `引起`” 等。
    *   **技术**:
        *   优先使用 `RelCAT` 模型来识别和提取预定义的关系类型。
        *   当 `RelCAT` 不可用或需要补充更复杂的关系时，可调用 **LLM 增强功能**，基于文本证据提取关系，并提供置信度和推理依据。
        *   在演示模式下，使用基于关键词和实体类型的规则进行基础关系推断。

4.  **临床文本去标识化 (Clinical Text De-identification)**
    *   **功能**: 自动检测并移除或遮蔽文本中的个人敏感信息 (PHI)，如姓名、电话、地址、病历号等，以保护患者隐私，满足合规要求。
    *   **技术**: 结合了 `MedCAT` 的 `DeIdModel` 模块和精细的正则表达式规则，提供 `mask` (遮蔽)、`replace` (替换为标签)、`remove` (移除) 等多种去标识化模式。

5.  **高效批量处理 (Efficient Batch Processing)**
    *   **功能**: 支持对大量文档进行并发处理，能够高效地完成大规模临床记录的实体识别任务。
    *   **技术**: 利用 `concurrent.futures.ThreadPoolExecutor` 实现多线程并发处理，显著提升数据处理吞吐量。

## 技术亮点

*   **混合模型架构 (Hybrid Architecture)**: 巧妙地融合了 `MedCAT` 这一专业领域的 NLP 模型和通用大语言模型 (LLM)，兼顾了医学领域的准确性和深度语义理解的灵活性。
*   **智能降级与增强 (Smart Fallback & Enhancement)**: 服务具备强大的鲁棒性。当 `MedCAT` 的高级组件（如 `MetaCAT`, `RelCAT`）不可用时，系统会自动切换到 LLM 进行功能增强；当所有高级模型都不可用时，则降级为基于规则的演示模式，确保核心功能的持续可用。
*   **高度可配置 (Highly Configurable)**: 通过环境变量可以轻松配置所有关键路径和参数，包括 `MedCAT` 模型路径、LLM 的 API Key、Base URL 和模型名称，极大地简化了部署和维护。
*   **标准化接口 (Standardized Interface)**: 基于 `FastMCP` 和 `Pydantic` 构建，所有工具都具有清晰、类型安全的输入输出定义，使得服务易于被其他系统集成和调用。
*   **全面的日志系统**: 提供详细的启动和运行时日志，清晰展示模型加载状态、LLM 功能是否启用以及当前运行模式，便于监控和调试。

## 配置环境变量（可选，用于LLM增强）

```bash
# DeepSeek 或其它大模型 API 配置
export OPENAI_API_KEY="sk-xxxxxxx"
# medcat 模型，可选
export MEDCAT_MODEL_PATH=xxxx # MedCAT模型包路径
export MEDCAT_CDB_PATH=xxxx # 概念数据库路径
export MEDCAT_VOCAB_PATH=xxxx # 词汇表路径
export METACAT_MODEL_PATH=xxxx # MetaCAT模型路径
export RELCAT_MODEL_PATH=xxxx # RelCAT模型路径
```

### 启动MCP服务器

```bash
python medcat_mcp_server.py
```

服务器将在 `http://0.0.0.0:8007/mcp` 上启动。