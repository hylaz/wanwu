package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/UnicomAI/wanwu/internal/agent-service/model/request"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/eino-contrib/jsonschema"
	"github.com/getkin/kin-openapi/openapi3"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

// convertOpenAPISchemaToJSONSchema 将 OpenAPI 3.0 Schema 转换为 eino-contrib/jsonschema Schema
func convertOpenAPISchemaToJSONSchema(openAPISchema *openapi3.Schema) *jsonschema.Schema {
	if openAPISchema == nil {
		return nil
	}

	js := &jsonschema.Schema{
		Type:        openAPISchema.Type,
		Description: openAPISchema.Description,
		Required:    openAPISchema.Required,
		Enum:        openAPISchema.Enum,
		Default:     openAPISchema.Default,
		Title:       openAPISchema.Title,
		Pattern:     openAPISchema.Pattern,
		Format:      openAPISchema.Format,
		ReadOnly:    openAPISchema.ReadOnly,
		WriteOnly:   openAPISchema.WriteOnly,
		Deprecated:  openAPISchema.Deprecated,
	}

	if openAPISchema.MinLength > 0 {
		minLen := uint64(openAPISchema.MinLength)
		js.MinLength = &minLen
	}
	if openAPISchema.MaxLength != nil {
		maxLen := uint64(*openAPISchema.MaxLength)
		js.MaxLength = &maxLen
	}
	if openAPISchema.Min != nil {
		js.Minimum = json.Number(fmt.Sprintf("%v", *openAPISchema.Min))
	}
	if openAPISchema.Max != nil {
		js.Maximum = json.Number(fmt.Sprintf("%v", *openAPISchema.Max))
	}
	if openAPISchema.MinItems > 0 {
		minItems := uint64(openAPISchema.MinItems)
		js.MinItems = &minItems
	}
	if openAPISchema.MaxItems != nil {
		maxItems := uint64(*openAPISchema.MaxItems)
		js.MaxItems = &maxItems
	}
	if openAPISchema.MinProps > 0 {
		minProps := uint64(openAPISchema.MinProps)
		js.MinProperties = &minProps
	}
	if openAPISchema.MaxProps != nil {
		maxProps := uint64(*openAPISchema.MaxProps)
		js.MaxProperties = &maxProps
	}

	if openAPISchema.UniqueItems {
		js.UniqueItems = true
	}

	if len(openAPISchema.Properties) > 0 {
		js.Properties = orderedmap.New[string, *jsonschema.Schema]()
		for name, schemaRef := range openAPISchema.Properties {
			if schemaRef != nil && schemaRef.Value != nil {
				js.Properties.Set(name, convertOpenAPISchemaToJSONSchema(schemaRef.Value))
			}
		}
	}

	if openAPISchema.Items != nil && openAPISchema.Items.Value != nil {
		js.Items = convertOpenAPISchemaToJSONSchema(openAPISchema.Items.Value)
	}

	if openAPISchema.AdditionalProperties.Has != nil && *openAPISchema.AdditionalProperties.Has {
		if openAPISchema.AdditionalProperties.Schema != nil && openAPISchema.AdditionalProperties.Schema.Value != nil {
			js.AdditionalProperties = convertOpenAPISchemaToJSONSchema(openAPISchema.AdditionalProperties.Schema.Value)
		}
	}

	if len(openAPISchema.AllOf) > 0 {
		js.AllOf = make([]*jsonschema.Schema, len(openAPISchema.AllOf))
		for i, schemaRef := range openAPISchema.AllOf {
			if schemaRef != nil && schemaRef.Value != nil {
				js.AllOf[i] = convertOpenAPISchemaToJSONSchema(schemaRef.Value)
			}
		}
	}

	if len(openAPISchema.AnyOf) > 0 {
		js.AnyOf = make([]*jsonschema.Schema, len(openAPISchema.AnyOf))
		for i, schemaRef := range openAPISchema.AnyOf {
			if schemaRef != nil && schemaRef.Value != nil {
				js.AnyOf[i] = convertOpenAPISchemaToJSONSchema(schemaRef.Value)
			}
		}
	}

	if len(openAPISchema.OneOf) > 0 {
		js.OneOf = make([]*jsonschema.Schema, len(openAPISchema.OneOf))
		for i, schemaRef := range openAPISchema.OneOf {
			if schemaRef != nil && schemaRef.Value != nil {
				js.OneOf[i] = convertOpenAPISchemaToJSONSchema(schemaRef.Value)
			}
		}
	}

	if openAPISchema.Not != nil && openAPISchema.Not.Value != nil {
		js.Not = convertOpenAPISchemaToJSONSchema(openAPISchema.Not.Value)
	}

	return js
}

// openAPITool 实现了 tool.InvokableTool 接口
type openAPITool struct {
	info    *schema.ToolInfo
	handler func(ctx context.Context, arguments string) (string, error)
}

// Info 返回工具的元信息
func (t *openAPITool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return t.info, nil
}

// InvokableRun 执行工具
func (t *openAPITool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	if t.handler == nil {
		return "", fmt.Errorf("tool handler is not set")
	}
	return t.handler(ctx, argumentsInJSON)
}

func GetToolsFromOpenAPISchema(ctx context.Context, pluginToolList []*request.PluginToolInfo) ([]tool.BaseTool, error) {
	if len(pluginToolList) == 0 {
		return nil, nil
	}
	var allTools []tool.BaseTool

	for _, wrapper := range pluginToolList {
		if wrapper.APISchema == nil {
			continue
		}

		loader := openapi3.NewLoader()
		loader.IsExternalRefsAllowed = true

		if err := wrapper.APISchema.Validate(ctx, openapi3.EnableExamplesValidation()); err != nil {
			log.Printf("Warning: OpenAPI schema validation failed: %v", err)
		}

		for path, pathItem := range wrapper.APISchema.Paths {
			operations := map[string]*openapi3.Operation{
				"get":    pathItem.Get,
				"post":   pathItem.Post,
				"put":    pathItem.Put,
				"delete": pathItem.Delete,
				"patch":  pathItem.Patch,
			}

			for method, operation := range operations {
				if operation == nil {
					continue
				}

				toolName := operation.OperationID
				if toolName == "" {
					toolName = fmt.Sprintf("%s_%s", method, path)
				}

				toolDesc := operation.Summary
				if toolDesc == "" {
					toolDesc = operation.Description
				}

				params := buildParametersSchema(operation)

				toolInfo := &schema.ToolInfo{
					Name: toolName,
					Desc: toolDesc,
				}

				if params != nil {
					toolInfo.ParamsOneOf = schema.NewParamsOneOfByJSONSchema(convertOpenAPISchemaToJSONSchema(params))
				}

				serverURL := ""
				if len(wrapper.APISchema.Servers) > 0 {
					serverURL = wrapper.APISchema.Servers[0].URL
				}

				contentType := getRequestContentType(operation)
				handler := createHTTPHandler(serverURL, path, method, wrapper.APIAuth, contentType)

				tools := &openAPITool{
					info:    toolInfo,
					handler: handler,
				}

				// 打印工具详细信息
				paramsInfo := "no parameters"
				if toolInfo.ParamsOneOf != nil {
					jsonSchema, err := toolInfo.ParamsOneOf.ToJSONSchema()
					if err == nil && jsonSchema != nil {
						paramsJSON, _ := json.MarshalIndent(jsonSchema, "", "  ")
						paramsInfo = string(paramsJSON)
					}
				}
				log.Printf("Loaded OpenAPI tool: %s\n  Description: %s\n  Method: %s %s\n  Parameters Schema:\n%s",
					toolName, toolDesc, method, path, paramsInfo)

				allTools = append(allTools, tools)
			}
		}
	}

	return allTools, nil
}

func getRequestContentType(operation *openapi3.Operation) string {
	if operation.RequestBody != nil && operation.RequestBody.Value != nil {
		for contentType := range operation.RequestBody.Value.Content {
			return contentType
		}
	}
	return "application/json"
}

func buildParametersSchema(operation *openapi3.Operation) *openapi3.Schema {
	properties := make(map[string]*openapi3.SchemaRef)
	var required []string

	for _, param := range operation.Parameters {
		paramVal := param.Value
		if paramVal == nil {
			continue
		}

		if paramVal.In == "header" || paramVal.In == "path" {
			continue
		}

		if paramVal.Schema != nil {
			properties[paramVal.Name] = paramVal.Schema
			if paramVal.Required {
				required = append(required, paramVal.Name)
			}
		}
	}

	if operation.RequestBody != nil && operation.RequestBody.Value != nil {
		for _, content := range operation.RequestBody.Value.Content {
			if content.Schema != nil && content.Schema.Value != nil {
				for propName, propSchema := range content.Schema.Value.Properties {
					properties[propName] = propSchema
				}
				if len(content.Schema.Value.Required) > 0 {
					required = append(required, content.Schema.Value.Required...)
				}
			}
			break
		}
	}

	if len(properties) == 0 {
		return nil
	}

	objectType := "object"
	return &openapi3.Schema{
		Type:       objectType,
		Properties: properties,
		Required:   required,
	}
}

func createHTTPHandler(serverURL, path, method string, auth *request.APIAuth, contentType string) func(ctx context.Context, arguments string) (string, error) {
	return func(ctx context.Context, arguments string) (string, error) {
		requestURL := serverURL + path

		var body io.Reader
		var actualContentType string

		// 解析 URL 以便添加查询参数(包括认证参数)
		parsedURL, err := url.Parse(requestURL)
		if err != nil {
			return "", fmt.Errorf("failed to parse URL: %w", err)
		}

		// 获取现有的查询参数
		queryValues := parsedURL.Query()

		// 如果认证方式是 query 参数,先添加认证参数
		if auth != nil && auth.Type == "apiKey" && auth.In == "query" {
			queryValues.Set(auth.Name, auth.Value)
		}

		// 处理 GET 请求的查询参数
		if method == "get" && arguments != "" {
			var params map[string]interface{}
			if err := json.Unmarshal([]byte(arguments), &params); err != nil {
				return "", fmt.Errorf("failed to parse arguments: %w", err)
			}

			// 添加业务查询参数
			for key, value := range params {
				queryValues.Set(key, fmt.Sprintf("%v", value))
			}
		}

		// 更新 URL 的查询参数
		parsedURL.RawQuery = queryValues.Encode()
		requestURL = parsedURL.String()

		if method == "post" || method == "put" || method == "patch" {
			if contentType == "multipart/form-data" {
				var params map[string]interface{}
				if err := json.Unmarshal([]byte(arguments), &params); err != nil {
					return "", fmt.Errorf("failed to parse arguments: %w", err)
				}

				bodyBuf := &bytes.Buffer{}
				writer := multipart.NewWriter(bodyBuf)

				for key, value := range params {
					var valueStr string
					switch v := value.(type) {
					case string:
						valueStr = v
					default:
						valueBytes, _ := json.Marshal(v)
						valueStr = string(valueBytes)
					}

					if err := writer.WriteField(key, valueStr); err != nil {
						return "", fmt.Errorf("failed to write field %s: %w", key, err)
					}
				}

				if err := writer.Close(); err != nil {
					return "", fmt.Errorf("failed to close multipart writer: %w", err)
				}

				body = bodyBuf
				actualContentType = writer.FormDataContentType()
			} else {
				body = bytes.NewBufferString(arguments)
				actualContentType = "application/json"
			}
		}

		methodUpper := http.MethodPost
		switch method {
		case "get":
			methodUpper = http.MethodGet
		case "post":
			methodUpper = http.MethodPost
		case "put":
			methodUpper = http.MethodPut
		case "delete":
			methodUpper = http.MethodDelete
		case "patch":
			methodUpper = http.MethodPatch
		}

		req, err := http.NewRequestWithContext(ctx, methodUpper, requestURL, body)
		if err != nil {
			return "", fmt.Errorf("failed to create request: %w", err)
		}

		if auth != nil && auth.Type == "apiKey" && auth.In == "header" {
			req.Header.Set(auth.Name, auth.Value)
		}

		if body != nil {
			req.Header.Set("Content-Type", actualContentType)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return "", fmt.Errorf("failed to execute request: %w", err)
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode >= 400 {
			return "", fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
		}

		return string(respBody), nil
	}
}
