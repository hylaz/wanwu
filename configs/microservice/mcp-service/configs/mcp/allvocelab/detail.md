官方AllVoicelab模型上下文协议（MCP）服务器，支持与强大的文本到语音和视频翻译API的交互。启用MCP客户，例如Claude Desktop，Cursor，Windsurf，OpenAI代理，可以产生语音，翻译视频并执行智能的语音转换。提供场景，例如全球市场的短剧本地化，AI生成的有声读物，AI驱动的电影/电视叙事制作。
链接：https://www.modelscope.cn/mcp/servers/@allvoicelab/AllVoiceLab

为什么选择AllVoicelab MCP服务器？
多引擎技术解锁了语音的无限可能性：使用简单的文本输入，您可以访问视频生成，语音合成，语音克隆等。
AI语音生成器（TTS）：具有超高现实主义的30多种语言的自然语音生成
语音改变者：实时语音转换，非常适合游戏，直播和隐私保护
人声分离：超快速5ms人声和背景音乐的分离，具有行业领先的精度
多语言配音：一单击的翻译和配音，用于简短的视频/电影，保持情感语气和节奏
语音到文本（STT）：AI驱动的多语言字幕生成，精度超过98％
字幕删除：即使在复杂的背景下，无缝的硬字幕擦除
语音克隆：三秒钟的超快速克隆与类似人类的声音综合
文档
中文文档

Quickstart
从中获取API键AllVoiceLab.
安装uv（Python软件包管理器），安装curl -LsSf https://astral.sh/uv/install.sh | sh
重要：不同区域中API的服务器地址需要匹配相应区域的键，否则将出现该工具不可用的错误。
区域	全局	大陆
allvoicelab_api_key	去AllVoiceLab	去AllVoiceLab
allvoicelab_api_domain	https://api.allvoicelab.com	https://api.allvoicelab.cn
克劳德桌面
转到Claude>设置>开发人员>编辑配置> Claude_desktop_config.json以包含以下内容：

{
  "mcpServers": {
    "AllVoceLab": {
      "command": "uvx",
      "args": ["allvoicelab-mcp"],
      "env": {
        "ALLVOICELAB_API_KEY": "<insert-your-api-key-here>",
        "ALLVOICELAB_API_DOMAIN": "<insert-api-domain-here>",
        "ALLVOICELAB_BASE_PATH":"optional, default is user home directory.This is uesd to store the output files."
      }
    }
  }
}
如果您使用的是Windows，则必须在Claude Desktop中启用“开发人员模式”才能使用MCP服务器。在左上方的汉堡菜单中单击“帮助”，然后选择“启用开发人员模式”。

光标
转到光标 - >首选项 - >光标设置 - > MCP->添加新的全局MCP服务器以添加上面的配置。

就是这样。您的MCP客户端现在可以与Allvoicelab进行交互。

可用方法
方法	简短说明
text_to_speech	将文字转换为语音
speep_to_speech	在保留语音内容的同时，将音频转换为另一个声音
isaly_human_voice	通过删除背景噪音和非语音声音来提取清洁人的声音
clone_voice	通过从音频示例克隆来创建自定义的语音配置文件
remove_subtitle	使用OCR从视频中删除字幕
video_translation_dubbing	将视频演讲转换为不同的语言
text_translation	将文本文件翻译成另一种语言
subtitle_extraction	使用OCR从视频中提取字幕
示例用法
⚠️警告：使用这些工具需要Allvoicelab信用。

1。语音文字
尝试询问：转换“在所有语音实验室，我们正在使用AI驱动的解决方案重塑音频工作流的未来，从而使各地的创造者都可以访问真实的声音。”变成声音。

image

2。语音转换
从上一个示例生成音频后，选择音频文件并询问：将其转换为男性声音。

image

3。删除背景噪音
选择一个带有丰富声音的音频文件（同时包含BGM和人的声音），然后问：删除背景噪声。

image

4。语音克隆
选择一个带有单个语音的音频文件，然后询问：克隆此语音。

image

5。视频翻译
选择一个视频文件（英语）并询问：将此视频转换为日语。

image

原始视频：

image

翻译后：

image

6。删除字幕
选择带有字幕的视频，然后询问：从此视频中删除字幕。

image

原始视频：

image

任务完成后：

image

7。文字翻译
选择一个长文字（例如，“愚蠢的老人去除山脉”），然后问：将此文字翻译成日语。 如果未指定语言，默认情况下它将翻译为英语。

image

8。字幕提取
选择带有字幕的视频，然后询问：从此视频中提取字幕。

image

任务完成后，您将获得一个SRT文件，如下所示：

image

故障排除
可以在：

Windows：c：\ users \ <用户名> \。MCP\ allvoicelab_mcp.log
macOS：〜/.mcp/allvoicelab_mcp.log
请通过电子邮件（tech@allvoicelab.com）与我们联系，并与日志文件联系