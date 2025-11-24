# Google Maps MCP Server

用于 Google Maps API 的 MCP Server。

## 工具

1. `maps_geocode`
   - 将地址转换为坐标
   - 输入: `address` (字符串)
   - 返回: location, formatted_address, place_id

2. `maps_reverse_geocode`
   - 将坐标转换为地址
   - 输入:
     - `latitude` (数字)
     - `longitude` (数字)
   - 返回: formatted_address, place_id, address_components

3. `maps_search_places`
   - 使用文本查询搜索地点
   - 输入:
     - `query` (字符串)
     - `location` (可选): { latitude: 数字, longitude: 数字 }
     - `radius` (可选): 数字（米，最大 50000）
   - 返回: 包含名称、地址、位置的地点数组

4. `maps_place_details`
   - 获取地点的详细信息
   - 输入: `place_id` (字符串)
   - 返回: name, address, contact info, ratings, reviews, opening hours

5. `maps_distance_matrix`
   - 计算点之间的距离和时间
   - 输入:
     - `origins` (字符串数组)
     - `destinations` (字符串数组)
     - `mode` (可选): "driving" | "walking" | "bicycling" | "transit"
   - 返回: 距离和持续时间矩阵

6. `maps_elevation`
   - 获取地点的海拔数据
   - 输入: `locations` (包含 {latitude, longitude} 的数组)
   - 返回: 每个点的海拔数据

7. `maps_directions`
   - 获取点之间的路线
   - 输入:
     - `origin` (字符串)
     - `destination` (字符串)
     - `mode` (可选): "driving" | "walking" | "bicycling" | "transit"
   - 返回: 包含步骤、距离、持续时间的路线详情

## 设置

### API 密钥
按照 这里 的说明获取 Google Maps API 密钥。

### 与 Claude Desktop 的使用

将以下内容添加到您的 `claude_desktop_config.json` 中：

#### Docker

```json
{
  "mcpServers": {
    "google-maps": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "-e",
        "GOOGLE_MAPS_API_KEY",
        "mcp/google-maps"
      ],
      "env": {
        "GOOGLE_MAPS_API_KEY": "<YOUR_API_KEY>"
      }
    }
  }
}
```

### NPX

```json
{
  "mcpServers": {
    "google-maps": {
      "command": "npx",
      "args": [
        "-y",
        "@modelcontextprotocol/server-google-maps"
      ],
      "env": {
        "GOOGLE_MAPS_API_KEY": "<YOUR_API_KEY>"
      }
    }
  }
}
```

## 构建

Docker 构建：

```bash
docker build -t mcp/google-maps -f src/google-maps/Dockerfile .
```

## 许可证

此 MCP server 采用 MIT 许可证。这意味着您可以自由使用、修改和分发该软件，但需遵循 MIT 许可证的条款和条件。有关更多详细信息，请参见项目存储库中的 LICENSE 文件。
