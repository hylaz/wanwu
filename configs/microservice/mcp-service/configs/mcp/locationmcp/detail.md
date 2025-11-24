# 位置信息查询 MCP 服务

这是一个基于 MCP (Model Context Protocol) 的位置信息查询服务，提供多种位置相关信息的查询功能。

## 功能特点

该服务提供以下四种查询功能：

1. **手机号归属地查询**：查询手机号的运营商和归属地信息
2. **IP地址归属地查询**：查询IP地址的所属地区和ISP信息
3. **银行卡信息查询**：查询银行卡的发卡行和卡类型
4. **身份证归属地查询**：查询身份证号对应的行政区划信息

## 技术栈

- TypeScript
- Node.js
- @modelcontextprotocol/sdk：MCP协议SDK
- DuckDB：用于高效查询手机号和IP地址数据
- CSV解析：用于银行卡和身份证号数据的查询

## 项目结构

```
.
├── data/                  # 数据文件目录
│   ├── mobile_location.duckdb    # 手机号归属地数据库
│   ├── ip_location.duckdb        # IP地址归属地数据库
│   ├── bank_card_bin_data.csv    # 银行卡BIN数据
│   └── id_location.csv           # 身份证归属地数据
├── src/                   # 源代码目录
│   ├── index.ts                  # 主入口文件
│   ├── mobile_location.ts        # 手机号归属地查询服务
│   ├── ip_location.ts            # IP地址归属地查询服务
│   ├── bank_card.ts              # 银行卡信息查询服务
│   └── id_location.ts            # 身份证归属地查询服务
├── dist/                  # 编译后的JavaScript文件
├── package.json           # 项目依赖配置
└── tsconfig.json          # TypeScript配置
```

## 安装和使用

### 安装依赖

```
npm install
```

### 编译TypeScript

```
npm run build
```

### 运行服务

```
npm start
```

或者在开发模式下运行：

```
npm run dev
```

## MCP配置示例

在使用大语言模型时，可以将此服务配置为MCP服务器。以下是配置示例：

```
{
  "mcpServers": {
    "location": {
      "command": "npx",
      "args": ["@data_wise/location-mcp"],
      "env": {}
    }
  }
}
```

您可以将此配置添加到您的MCP配置文件中，然后大语言模型就可以调用这些位置查询工具。

## 数据文件

- **mobile_location.duckdb**：包含手机号段与归属地的对应关系
- **ip_location.duckdb**：包含IP地址段与归属地的对应关系
- **bank_card_bin_data.csv**：包含银行卡BIN码、发卡行和卡类型信息
- **id_location.csv**：包含身份证号行政区划代码与对应地区名称

## 查询接口

### 手机号归属地查询

```
queryMobileLocation: 查询手机号归属地
参数: {
  mobileNumbers: string[]  // 手机号数组
}
返回: {
  phoneNumber: string      // 手机号
  province: string         // 省份
  city: string             // 城市
  operator: string         // 运营商
}
```

### IP归属地查询

```
queryIpLocation: 查询IP地址归属地
参数: {
  ipAddresses: string[]    // IP地址数组
}
返回: {
  ip: string               // IP地址
  isp: string              // 互联网服务提供商
  region: string           // 地区
}
```

### 银行卡信息查询

```
queryBankCard: 查询银行卡信息
参数: {
  cardNumbers: string[]    // 银行卡号数组
}
返回: {
  cardNumber: string       // 银行卡号
  bankName: string         // 银行名称
  cardType: string         // 卡类型（借记卡/信用卡）
}
```

### 身份证归属地查询

```
queryIdLocation: 查询身份证归属地
参数: {
  idNumbers: string[]      // 身份证号数组
}
返回: {
  idNumber: string         // 身份证号
  location: string         // 归属地区
}
```

## 使用示例

您可以通过标准输入输出与服务进行通信，符合MCP协议规范。也可以使用提供的测试脚本进行测试：

```
# 测试手机号归属地查询
npx ts-node src/test.ts

# 测试IP归属地查询
npx ts-node src/ip_test.ts

# 测试银行卡信息查询
npx ts-node src/bank_card_test.ts

# 测试身份证归属地查询
npx ts-node src/id_test.ts
```

## 开发说明

### 添加新的查询服务

1. 在`src`目录下创建新的服务类文件
2. 实现相应的查询逻辑
3. 在`index.ts`中引入并注册新的工具

### 更新数据

各服务的数据文件位于`data`目录下，可以按需更新：

- 手机号和IP地址数据使用DuckDB数据库
- 银行卡和身份证号数据使用CSV文件

## 许可证

本项目采用MIT许可证