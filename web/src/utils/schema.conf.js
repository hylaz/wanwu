export const schemaConfig = {
  json: `{
    "openapi": "3.0.0",
    "info":
      {
        "title": "心知天气API",
        "version": "1.0.0",
        "description": "提供当前天气信息的API，包括温度、天气状况等。"
      },
    "servers":
      [
        {
          "url": "https://api.seniverse.com/v3"
        }
      ],
    "paths":{
      "/weather/now.json": {
        "get": {
          "summary": "天气查询工具",
          "operationId": "getWeatherNow",
          "description": "根据地点获取当前的天气情况，包括温度和天气状况描述。",
          "parameters": [{
            "name": "location",
            "description": "查询的地点，可以是城市名、邮编等。",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string"
            }
          }],
          "responses": {
            "200": {
              "description": "成功获取天气信息",
              "content": {
                "application/json": {
                  "schema": {
                    "type": "object",
                    "properties": {
                      "location": { "type": "string"},
                      "text": { "type": "string" },
                      "code": { "type": "string"},
                      "temperature": { "type": "string" }
                    }
                  }
                }
              }
            },
            "default": {
              "description": "请求失败时的错误信息",
              "content": {
                "application/json": {
                  "schema": {
                    "type": "object",
                    "properties": {
                      "error": { "type": "string" }
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }`,
  yaml: `
    openapi: 3.0.0
    info:
      title: 心知天气API
      version: 1.0.0
      description: 提供当前天气信息的API，包括温度、天气状况等。
    servers:
      - url: https://api.seniverse.com/v3
    paths:
      /weather/now.json:
        get:
          summary: 天气查询工具
          operationId: getWeatherNow
          description: 根据地点获取当前的天气情况，包括温度和天气状况描述。
          parameters:
            - name: location
              description: 查询的地点，可以是城市名、邮编等。
              in: query
              required: true
              schema:
                type: string
          responses:
            '200':
              description: 成功获取天气信息
              content:
                application/json:
                  schema:
                    type: object
                    properties:
                      location:
                        type: string
                      text:
                        type: string
                      code:
                        type: string
                      temperature:
                        type: string
            default:
              description: 请求失败时的错误信息
              content:
                application/json:
                  schema:
                    type: object
                    properties:
                      error:
                        type: string
  `,
};
