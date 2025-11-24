MCP Tool Server 用于 Logo 生成
链接：https://www.modelscope.cn/mcp/servers/@sshtunnelvision/MCP-LOGO-GEN
此服务器使用 FAL AI 提供 logo 生成功能，包括图像生成、背景移除和自动缩放等工具。

演示
MCP Tool Server 演示

安装
安装 uv（通用虚拟环境）：
curl -LsSf https://astral.sh/uv/install.sh | sh
创建并激活一个虚拟环境：
uv venv
source .venv/bin/activate  # On Unix/macOS
# or
.venv\Scripts\activate     # On Windows
安装依赖项：
uv pip install -r requirements.txt
设置你的环境变量：
在根目录下创建一个 .env 文件
添加你的 FAL AI API 密钥：
FAL_KEY=your_fal_ai_key_here
运行服务器
使用以下命令启动服务器：

python run_server.py
服务器将在 http://127.0.0.1:7777 上可用。

故障排除
如果你在 Windows 上运行服务器时遇到 FileNotFoundError，请确保你是在项目的根目录下运行该命令。如果问题仍然存在，请尝试更新到包含 Windows 兼容性修复的最新版本仓库。

对于 Windows 用户特别提示：

确保你已经通过 .venv\Scripts\activate 激活了虚拟环境
从项目的根目录用 python run_server.py 命令运行服务器
如果看到任何路径相关的错误，请在仓库的问题部分报告它们
Cursor IDE 配置
打开 Cursor 设置
导航到 MCP 部分
添加以下配置：
URL: http://127.0.0.1:7777/sse
连接类型: SSE
启用连接
注意事项
在 Cursor Composer 中始终引用 @logo-creation.mdc 以获得一致的结果
步骤定义在 @logo-creation.mdc 中，但工具可以独立使用
所有生成的 logo 将保存在 downloads 目录中
每个 logo 自动以三种尺寸生成：
原始尺寸
32x32 像素
128x128 像素
所有 logo 在最终 PNG 格式中保持透明度
由代理创建的提示信息基于在 server.py 中看到的例子和提示结构。你可以通过编辑 server.py 文件来自定义提示结构。
你可以使用 generate_image 工具来生成任何你想要的图像，而不仅仅是 logo
要求
Python 3.8+
FAL AI API 密钥（用于图像生成）
活跃的互联网连接