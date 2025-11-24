# 且慢MCP

盈米且慢推出国内首个兼容MCP协议的财富管理服务平台—Qieman MCP Server（简称且慢MCP）。且慢MCP提供了完整的金融数据与专业分析工具支持，让您的AI大模型能实时提供高质量、可用且数据真实准确的金融服务。助力财富管理的行业伙伴及投资人，能一同探索体验，解决自身以及客户的真实金融问题。

## 主要特点

1. **精确的金融数据**：调用且慢MCP可获取实时准确的金融数据，可避免常见的大模型幻觉；
2. **全面的投资研究系统**：整合了金融专家的最新分析观点和汇总信息，以提高金融信息服务的质量。
3. **专业的投资顾问系统**：调用专业投顾测算工具与各类资产分析能力，构建金融与投顾服务核心优势；
4. **便捷的用户体验**：初始配置后，服务将自动更新和升级，确保能持续使用最新的金融工具服务包。

## 工具列表

| 功能/API名称                | 描述                                                         |
| --------------------------- | ------------------------------------------------------------ |
| BatchGetFundsDetail         | 批量查询多个基金的详细信息，包括名称、类型、规模、风险级别、经理、成立日期、投资范围等。 |
| BatchGetFundNavHistory      | 批量获取多个基金的历史净值数据，支持按不同时间维度查询。     |
| GetBatchFundPerformance     | 批量获取基金的业绩数据，包括业绩分析指标和阶段回报。         |
| BatchGetFundsHolding        | 批量获取多个基金的持仓信息，包括十大重仓股、债券持仓等。     |
| BatchGetFundsHolderInfo     | 批量获取基金的资产规模和持有人结构数据，包括单份数据、总份数据和持有人结构信息。 |
| BatchGetFundsFeeRule        | 批量获取多个基金的费用规则，包括认购费、申购费、赎回费和操作费。 |
| BatchGetFundTradeRules      | 批量获取基金的交易规则，包括最低/最高购买金额、预期确认日期、预期到账日期、费用规则等。 |
| BatchGetFundTradeLimit      | 批量获取多个基金的交易限制信息，包括基金是否可交易、最低购买金额、最低持有份额、定投起点等。 |
| BatchGetFundsSplitHistory   | 批量获取基金的拆分记录信息，包括拆分日期和比例。             |
| BatchGetFundsDividendRecord | 批量获取基金的分红记录，包括权益登记日、红利发放日和每份红利金额。 |
| GetFundAnnouncements        | 查询基金公告，包括基金代码、名称、全名、公告ID、日期、来源、标题、链接和类型。支持按时间范围、公告类型和标题关键词查询。 |
| SearchFunds                 | 基于名称、代码或其他条件搜索基金，支持按回报、规模、费用等排序。 |
| GetPopularFund              | 获取最近受欢迎基金的列表，了解市场焦点和投资趋势。           |
| GuessFundCode               | 根据基金名称匹配最接近的基金代码。                           |
| GetAssetAllocation          | 分析基金组合的资产配置，提供资产类别分布、雷达图评分和诊断结果。 |
| GetFundsCorrelation         | 分析多个基金之间的相关系数，了解组合中每个基金趋势的相互影响。 |
| GetFundsBackTest            | 对给定的基金组合进行历史回测分析，计算关键指标并提供全面的诊断结果和评分。 |
| MonteCarloSimulate          | 基于给定的资产配置权重执行蒙特卡洛模拟计算，生成预期收益分布、波动率情景和各种百分位数的收益率数据。 |
| AnalyzeFundRisk             | 获取多个基金的风险评分和详细风险解释，计算风险评分、R平方、残差方差、标准误差等指标。 |
| AnalyzePortfolioRisk        | 对给定的基金组合进行风险评估分析，计算多维风险指标。         |
| SearchFinancialNews         | 搜索工具支持按关键词和时间范围进行筛选和分页，返回详细的新闻内容。 |
| SearchManagerViewpoint      | 获取基金经理对各行业的见解和市场分析，支持按行业主题、时间范围和关键词进行筛选。 |
| GetFundDiagnosis            | 提供全面的基金诊断分析，包括最新信息、风险和机会评估、行业持仓、资产配置、业绩等。 |
| DiagnoseFundPortfolio       | 提供全面的基金组合诊断分析，包括资产配置分析、基金相关性分析和回测诊断。 |
| GetAssetAllocationPlan      | 基于投资参数获取资产配置计划。                               |
| GetCompositeModel           | 通过资产配置计划ID获取相应的复合模型。                       |
| RenderEchart                | 根据提供的ECharts配置渲染图表并转换为图像，支持返回base64编码的图像内容或OSS访问URL。 |
| RenderHtmlToPdf             | 将HTML内容转换为PDF文档，支持自定义PDF格式和边距设置。       |
| GetCurrentTime              | 获取服务器的当前时间，返回格式化的日期和时间字符串。         |
| GetTxnDayRange              | 基于中心时间获取指定期间内的交易日列表，支持向前和向后推算指定天数。 |

## 官方网站

https://qieman.com/mcp/landing

## 开始使用

1. **申请开通Qieman MCP服务**： 访问且慢MCP官方网站，登录或注册且慢账户，填写表格获取您专属的API密钥。
2. **使用链接开始获取服务**： https://stargate.yingmi.com/mcp/sse?apiKey=[您的API密钥]

## 用户指南

https://qieman.com/mcp/how-to-use

## 联系我们

oap@yingmi.cn