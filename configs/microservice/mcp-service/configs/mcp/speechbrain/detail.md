# 医疗声纹识别 - 声纹认证和身份识别服务

## 📋 服务信息

**服务名称**: 医疗声纹识别 - 声纹认证和身份识别服务

**基于项目**: SpeechBrain

**简短描述**: 基于 SpeechBrain ECAPA-TDNN 模型，专门针对医疗场景优化的声纹识别服务。支持医生、患者、护士等医疗角色的声纹注册、验证和识别，提供医疗场景说话人分析、多说话人聚类、声纹特征提取等功能。适用于医疗身份认证、病历录音分析、医疗会议记录等应用场景。

**部署URL**: http://localhost:9098/sse

---

## 🚀 快速启动

### 安装依赖
```bash
pip install speechbrain torch scikit-learn librosa soundfile fastmcp
```

### 启动服务
```bash
python speechbrain_medical_voiceprint_mcp.py
```

### 运行测试
```bash
python test_speechbrain_medical_voiceprint.py
```

---

## 🔧 Cursor MCP 配置

在 `.cursor/mcp_settings.json` 中添加：

```json
{
  "mcpServers": {
    "speechbrain_medical_voiceprint": {
      "type": "sse",
      "url": "http://localhost:9098/sse",
      "timeout": 300
    }
  }
}
```

---

## 🎯 核心功能

- **声纹注册**: 医疗人员身份登记和声纹采集
- **声纹验证**: 1:1 身份验证（验证是否为指定人员）
- **声纹识别**: 1:N 身份识别（从数据库中识别说话人）
- **医疗场景分析**: 门诊咨询、病房查房、医疗会议等场景分析
- **说话人聚类**: 多说话人自动聚类和分组
- **特征提取**: 声纹特征提取和质量评估
- **数据库管理**: 完整的说话人数据库CRUD操作

---

## 💡 使用示例

```
@speechbrain_medical_voiceprint 注册张医生的声纹，角色为医生，科室为心内科
@speechbrain_medical_voiceprint 验证这个声音是否是张医生
@speechbrain_medical_voiceprint 识别这个病历录音中的说话人
@speechbrain_medical_voiceprint 分析这个医患对话中的角色分布
@speechbrain_medical_voiceprint 对多个会议录音进行说话人聚类
```

---

## 🏥 支持的医疗角色

| 角色 | 描述 | 优先级 |
|-----|------|--------|
| doctor | 主治医生、专家医生等 | 1 |
| patient | 病人、患者家属等 | 2 |
| nurse | 护士、护理人员等 | 3 |
| therapist | 康复师、心理治疗师等 | 4 |
| technician | 检验师、影像技师等 | 5 |
| administrator | 医院管理人员等 | 6 |

---

## 📊 技术特性

- **高精度**: 基于深度学习的 ECAPA-TDNN 模型
- **实时处理**: 支持实时声纹特征提取和匹配
- **多格式支持**: WAV、MP3、FLAC、M4A、Base64编码音频
- **医疗专业化**: 针对医疗场景角色和权限优化
- **数据安全**: 本地化数据存储，隐私保护
- **高兼容性**: 支持CPU/GPU运算，自适应设备选择

---

## 📁 数据存储

服务会在当前目录创建 `speechbrain_medical_data` 文件夹存储：
- `speaker_database.json`: 说话人档案数据库
- `voice_features_cache.pkl`: 声纹特征缓存

---

## ⚠️ 注意事项

1. 首次启动时会自动下载 SpeechBrain 模型（需要网络连接）
2. 建议音频时长在 1-30 秒之间以获得最佳识别效果
3. 支持 16kHz 采样率音频，其他采样率会自动转换
4. GPU 加速需要正确安装 CUDA 和 PyTorch GPU 版本