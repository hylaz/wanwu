基本能力
产品定位
Image Generation MCP Server 是一个图像生成服务，基于Replicate Flux模型，提供通过文本描述生成图像的功能。

核心功能
使用文本提示生成图像
支持多种图像格式（webp、jpg、png）
可配置图像宽高比
支持生成多张图像（1-4张）
支持随机种子设置以实现可重复生成
适用场景
需要快速生成图像的创意设计
基于文本描述生成视觉内容的应用
需要批量生成图像的自动化流程
工具列表
generate_image: 根据文本提示生成图像
参数:
prompt (必填): 图像生成的文本描述
seed (可选): 随机种子
aspect_ratio (可选): 图像宽高比（默认"1:1"）
output_format (可选): 输出格式（"webp", "jpg", 或 "png"，默认"webp"）
num_outputs (可选): 生成图像数量（1-4，默认1）
使用教程
使用依赖
需要Node.js环境和npm包管理器。

安装教程
通过Smithery安装
bash
npx -y @smithery/cli install @GongRzhe/Image-Generation-MCP-Server --client claude
方法1: NPX方式（无需本地安装）
bash
# 无需安装 - npx会自动处理
方法2: 本地安装
bash
# 全局安装
npm install -g @gongrzhe/image-gen-server

# 或本地安装
npm install @gongrzhe/image-gen-server
调试方式
获取Replicate API Token:

注册/登录 https://replicate.com
访问 https://replicate.com/account/api-tokens
创建新API token
复制token并替换配置中的your-replicate-api-token
配置Claude Desktop:

MacOS: ~/Library/Application Support/Claude/claude_desktop_config.json
Windows: %APPDATA%/Claude/claude_desktop_config.json
NPX配置（推荐）
json
{
  "mcpServers": {
    "image-gen": {
      "command": "npx",
      "args": ["@gongrzhe/image-gen-server"],
      "env": {
        "REPLICATE_API_TOKEN": "your-replicate-api-token",
        "MODEL": "alternative-model-name"
      },
      "disabled": false,
      "autoApprove": []
    }
  }
}
本地安装配置
json
{
  "mcpServers": {
    "image-gen": {
      "command": "node",
      "args": ["/path/to/image-gen-server/build/index.js"],
      "env": {
        "REPLICATE_API_TOKEN": "your-replicate-api-token",
        "MODEL": "alternative-model-name"
      },
      "disabled": false,
      "autoApprove": []
    }
  }
}