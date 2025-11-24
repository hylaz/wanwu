# 医药语音识别 - 专业医药语音转写服务

## 📋 服务信息

**服务名称**: 医药语音识别 - 专业医药语音转写服务

**基于项目**: Paraformer (中科院语音所)

**简短描述**: 基于中科院语音所 Paraformer-large 模型的专业医药语音识别服务。支持 16kHz 中文语音识别，医药专业术语增强，流式音频处理和文件识别。内置丰富的医药专业词汇库，包括常用药品名、医学术语、药理术语、解剖术语等，可自动纠正识别错误并优化医药文本准确性。

**部署URL**: http://localhost:9100/sse

---

## 🚀 快速启动

### 安装依赖
```bash
pip install funasr librosa soundfile fastmcp
```

### 启动服务
```bash
python paraformer_medical_asr_mcp.py
```

### 运行测试
```bash
python test_paraformer_medical_asr.py
```

---

## 🔧 Cursor MCP 配置

在 `.cursor/mcp_settings.json` 中添加：

```json
{
  "mcpServers": {
    "paraformer_medical_asr": {
      "type": "sse",
      "url": "http://localhost:9100/sse",
      "timeout": 600
    }
  }
}
```

---

## 🎯 核心功能

- **医药语音流式识别**: 实时处理Base64编码音频数据，快速转写
- **医药语音文件识别**: 支持多种音频文件格式的批量转写
- **医药文本后处理增强**: 自动纠正医药术语识别错误，提升准确性
- **医药专业词汇管理**: 内置4大类专业词汇，支持词汇查询和分析
- **音频格式检查**: 智能检测音频兼容性，提供优化建议
- **模型信息监控**: 实时获取模型状态和性能信息

---

## 💡 使用示例

```
@paraformer_medical_asr 识别这个医药语音文件中的内容
@paraformer_medical_asr 转写这段医生处方录音为文本
@paraformer_medical_asr 增强这段医药文本的准确性
@paraformer_medical_asr 获取支持的医药专业词汇
@paraformer_medical_asr 检查这个音频文件的格式兼容性
```

---

## 🏥 医药专业词汇库

| 类别 | 描述 | 示例 |
|-----|------|------|
| 药品名 | 常用药品和药物名称 | 阿司匹林、布洛芬、胰岛素 |
| 医学术语 | 疾病和症状术语 | 高血压、糖尿病、心房颤动 |
| 药理术语 | 用药方法和剂量 | 口服、静脉注射、每日三次 |
| 解剖术语 | 人体器官和组织 | 心脏、肝脏、血管、神经 |

---

## 📊 技术特性

- **高精度识别**: 基于 Paraformer-large 大模型，8404词汇表
- **医药专业化**: 专门针对医药场景优化的词汇纠错
- **实时处理**: 支持流式音频和文件两种处理模式
- **多格式支持**: WAV、MP3、FLAC、M4A、OGG等格式
- **智能增强**: 自动识别和纠正常见医药术语错误
- **16kHz优化**: 专为16kHz采样率优化，自动重采样

---

## 🔧 工具函数

| 函数名 | 功能描述 |
|--------|----------|
| `medical_asr_stream` | 医药语音流式识别 |
| `medical_asr_file` | 医药语音文件识别 |
| `enhance_medical_text` | 医药文本后处理增强 |
| `get_medical_vocab` | 获取医药专业词汇 |
| `check_audio_format` | 检查音频格式兼容性 |
| `get_model_info` | 获取模型信息 |

---

## 📁 模型配置

```python
MODEL_CONFIG = {
    "model_path": "/Geo_mcp/iic/speech_paraformer-large_asr_nat-zh-cn-16k-common-vocab8404-pytorch",
    "sample_rate": 16000,
    "device": "cuda",  # 或 "cpu"
    "batch_size": 1
}
```

---

## ⚠️ 注意事项

1. 首次启动需要确保 Paraformer 模型已正确下载到指定路径
2. 建议音频时长在 0.5-60 秒之间以获得最佳识别效果
3. 支持 16kHz 采样率，其他采样率会自动转换
4. GPU 加速需要 CUDA 环境支持
5. 多声道音频会自动转换为单声道处理

---

## 🎯 应用场景

- **医院门诊**: 医生处方录音自动转写
- **药房管理**: 药品名称语音识别和核对
- **医学教育**: 医学术语语音学习和测试
- **病历录入**: 医生查房记录语音转文字
- **药物咨询**: 患者用药咨询录音整理
- **医疗会议**: 学术讨论和病例分析记录