package service

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	net_url "net/url"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/bff-service/config"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func FileUrlConvertBase64(ctx *gin.Context, url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_file_http_get", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_file_http_get", fmt.Sprintf("StatusCode: %d", resp.StatusCode))
	}
	fileData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_file_read", err.Error())
	}
	base64Data := base64.StdEncoding.EncodeToString(fileData)

	return base64Data, nil
}

func UploadFileToWorkflow(ctx *gin.Context, req *request.WorkflowUploadFileReq) (*response.UploadFileByWorkflowResp, error) {
	file, err := req.File.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	base64Str := base64.StdEncoding.EncodeToString(fileBytes)
	return UploadFileBase64ToWorkflow(ctx, req.FileName, base64Str)
}

func UploadFileBase64ToWorkflow(ctx *gin.Context, fileName, file string) (*response.UploadFileByWorkflowResp, error) {
	url, _ := net_url.JoinPath(config.Cfg().Workflow.Endpoint, config.Cfg().Workflow.UploadFileUri)
	ret := &response.UploadFileByWorkflowResp{}
	requestBody := map[string]string{
		"name": fileName,
		"data": file,
	}
	if resp, err := resty.New().
		R().
		SetContext(ctx.Request.Context()).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeaders(workflowHttpReqHeader(ctx)).
		SetBody(requestBody).
		SetResult(ret).
		Post(url); err != nil {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_upload_file", err.Error())
	} else if resp.StatusCode() >= 300 {
		b, err := io.ReadAll(resp.RawResponse.Body)
		if err != nil {
			return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_upload_file", fmt.Sprintf("[%v] %v", resp.StatusCode(), err))
		}
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_upload_file", fmt.Sprintf("[%v] %v", resp.StatusCode(), string(b)))
	}
	return ret, nil
}
