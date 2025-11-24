Hive MCP 服务器
smithery 徽章

一个通过模型上下文协议使 AI 助手能够与 Hive 区块链交互的 MCP 服务器。

概述
该服务器为 AI 助手（如 Claude）和 Hive 区块链之间提供了桥梁，允许 AI 模型执行以下操作：

获取账户信息和历史记录
检索博客文章和讨论
按标签或用户获取文章
对内容进行投票并创建文章（在适当认证后）
向其他账户发送 HIVE 或 HBD 代币
使用 Hive 密钥签名和验证消息
发送和接收加密消息
功能
提示词
create-post - 创建结构化的提示词，引导 AI 以正确的格式和标签创建新的 Hive 文章
analyze-account - 生成提示词，分析 Hive 账户的统计数据、发布历史和活动模式
工具
读取数据
get_account_info - 获取有关 Hive 区块链账户的详细信息
get_post_content - 根据作者和永久链接检索特定的文章
get_posts_by_tag - 按标签和类别（热门、趋势等）检索文章
get_posts_by_user - 从特定用户或其订阅中获取文章
get_account_history - 获取账户的交易历史记录，并可选择过滤操作
get_chain_properties - 获取当前 Hive 区块链的属性和统计信息
get_vesting_delegations - 获取由特定账户发起的质押委托列表
区块链交互（需要认证）
vote_on_post - 对 Hive 内容进行投票（需要发布密钥）
create_post - 在 Hive 区块链上创建新博客文章（需要发布密钥）
create_comment - 对现有文章发表评论或回复评论（需要发布密钥）
send_token - 将 HIVE 或 HBD 加密货币发送到其他账户（需要活跃密钥）
加密技术
sign_message - 使用 Hive 私钥对消息进行签名
verify_signature - 验证消息签名是否与 Hive 公钥匹配
加密消息
encrypt_message - 为特定的 Hive 账户加密消息
decrypt_message - 解密来自特定 Hive 账户的加密消息
send_encrypted_message - 使用代币转账发送加密消息
get_encrypted_messages - 从账户历史记录中检索并可选地解密消息
使用 MCP Inspector 进行调试
MCP Inspector 提供了一个交互式界面用于测试和调试服务器：

npx @modelcontextprotocol/inspector npx @gluneau/hive-mcp-server
认证配置
要启用认证操作（投票、发帖、发送代币），您需要设置环境变量：

export HIVE_USERNAME=your-hive-username
export HIVE_POSTING_KEY=your-hive-posting-private-key  # For content operations
export HIVE_ACTIVE_KEY=your-hive-active-private-key    # For token transfers
export HIVE_MEMO_KEY=your-hive-memo-private-key        # For encrypted messaging
安全注意事项：切勿共享您的私钥或将它们提交到版本控制系统。请使用环境变量或安全的配置方法。

与 AI 助手集成
Claude 桌面版
要将此服务器与 Claude Desktop 一起使用：

确保已安装 Claude Desktop

打开或创建 Claude 配置文件：

macOS: ~/Library/Application Support/Claude/claude_desktop_config.json
Windows: %APPDATA%\Claude\claude_desktop_config.json
Linux: ~/.config/Claude/claude_desktop_config.json
将此服务器添加到您的配置中：

{
  "mcpServers": {
    "hive": {
      "command": "npx",
      "args": ["-y", "@gluneau/hive-mcp-server"],
      "env": {
        "HIVE_USERNAME": "your-hive-username",
        "HIVE_POSTING_KEY": "your-hive-posting-private-key",
        "HIVE_ACTIVE_KEY": "your-hive-active-private-key",
        "HIVE_MEMO_KEY": "your-hive-memo-private-key"
      }
    }
  }
}
Windsurf 和 Cursor
相同的 JSON 配置适用于 Windsurf（在 windsurf_config.json 中）和 Cursor（对于版本 >= 0.47，在 ~/.cursor/mcp.json 中）。

在早期版本中，您需要在设置的 MCP 部分使用单行命令格式： env HIVE_USERNAME=your-hive-username env HIVE_POSTING_KEY=your-hive-posting-private-key env HIVE_ACTIVE_KEY=your-hive-active-private-key env HIVE_MEMO_KEY=your-hive-memo-private-key npx -y @gluneau/hive-mcp-server

示例
连接到 MCP 客户端后，您可以提出如下问题：

"Hive 上 #photography 标签下的热门帖子有哪些？"
"显示用户名 'alice' 的最新帖子"
"查询 'bob' 的账户余额和详细信息"
"获取 'charlie' 的交易历史记录"
"能否给 'dave' 发布的带有 permlink 'my-awesome-post' 的帖子点赞？"
"在 Hive 上创建一篇关于 AI 技术的新帖子"
"向用户 'frank' 发送 1 HIVE 并附上备注 'Thanks for your help!'"
"用我的 Hive 发帖密钥签名这条消息：'Verifying my identity'"
"当前 Hive 区块链的属性是什么？"
"显示用户 'grace' 做出的质押委托"
"为用户 'alice' 加密这条消息：'This is a secret message'"
"解密来自 'bob' 的消息：'#4f3a5b...'"
"向 'charlie' 发送一条加密消息，内容是 'Let's meet tomorrow'"
"显示并解密我的加密消息"
"获取我与 'dave' 交换的最后 10 条加密消息"
工具文档
get_account_info
获取 Hive 区块链账户的详细信息，包括余额、权限、投票权等其他指标。

参数：
username: 要查询信息的 Hive 用户名
get_post_content
通过作者和永久链接检索特定的 Hive 博客文章。

参数：
author: 文章作者
permlink: 文章的永久链接
get_posts_by_tag
根据特定标签筛选并按类别排序来检索 Hive 文章。

参数：
category: 排序类别（如热门、最新、创建时间等）
tag: 用于筛选文章的标签
limit: 返回的文章数量（1-20）
get_posts_by_user
检索特定 Hive 用户发布的或在其动态中的文章。

参数：
category: 要获取的用户文章类型（博客或动态）
username: 要为其获取文章的 Hive 用户名
limit: 返回的文章数量（1-20）
get_account_history
检索 Hive 账户的交易历史记录，并可选择性地按操作类型过滤。

参数：

account: 要查询历史记录的 Hive 账户名
start: 开始的交易索引
end: 结束的交易索引
operation_types: 可选参数，指定要过滤的操作类型列表
参数：

username: Hive 用户名
limit: 返回的操作数量
operation_filter: 可选的操作类型列表，用于过滤
get_chain_properties
获取当前的 Hive 区块链属性和统计数据。

参数：无
get_vesting_delegations
获取特定 Hive 账户的委托列表。

参数：
username: 获取委托的 Hive 账户
limit: 要检索的最大委托数量
from: 可选的分页起始账户
vote_on_post
使用配置的 Hive 账户对 Hive 帖子进行投票（点赞或点踩）。

参数：
author: 要投票的帖子的作者
permlink: 要投票的帖子的永久链接
weight: 投票权重，从 -10000（100% 点踩）到 10000（100% 点赞）
create_post
使用配置的账户在 Hive 区块链上创建新的博客文章。

参数：
title: 博客文章的标题
body: 博客文章的内容（支持 Markdown）
tags: 文章的标签
各种可选参数，如奖励、受益人等
create_comment
在现有的 Hive 帖子上发表评论或回复其他评论。

参数：
parent_author: 您要回复的帖子作者或评论者的用户名
parent_permlink: 您要回复的帖子或评论的永久链接
body: 评论的内容（支持 Markdown）
各种可选参数，如奖励、受益人等
send_token
使用配置的账户向另一个 Hive 账户发送 HIVE 或 HBD 代币。

参数：
to: 收件人的 Hive 用户名
amount: 要发送的代币数量
currency: 要发送的货币（HIVE 或 HBD）
memo: 可选的交易备注
sign_message
使用环境变量中的 Hive 私钥签署消息。

参数：
message: 要签名的消息
key_type: 使用的密钥类型（posting、active 或 memo）
verify_signature
验证数字签名是否与 Hive 公钥匹配。

参数：
message_hash: 消息的 SHA-256 哈希值（十六进制格式）
signature: 要验证的签名字符串
public_key: 用于验证的公钥
encrypt_message
使用备忘录加密为特定 Hive 账户加密消息。

参数：
message: 要加密的消息
recipient: 收件人的 Hive 用户名
decrypt_message
解密从特定 Hive 账户收到的加密消息。

参数：
encrypted_message: 加密的消息（以 # 开头）
sender: 发送者的 Hive 用户名
send_encrypted_message
使用小额代币转账发送加密消息给 Hive 账户。

参数：
message: 要加密并发送的消息
recipient: 收件人的 Hive 用户名
amount: 要发送的 HIVE 数量（最小 0.001，默认：0.001）
get_encrypted_messages
从账户历史中检索加密消息，并可选择解密。

参数：
username: 要获取加密消息的Hive用户名
limit: 要检索的最大消息数量（默认：20）
decrypt: 是否尝试对消息进行解密（默认：false）
开发
项目结构
src/index.ts - 主服务器实现
src/tools/ - 所有工具的实现
src/schemas/ - 工具参数的Zod模式
src/utils/ - 与Hive区块链交互的实用函数
src/config/ - 客户端配置和日志级别处理
依赖项
@hiveio/dhive - Hive区块链客户端
@modelcontextprotocol/sdk - MCP SDK
zod - 模式验证
许可证
ISC

贡献
欢迎贡献！请随时提交Pull Request。

有关更详细的贡献指南，请参阅CONTRIBUTING.md文件。