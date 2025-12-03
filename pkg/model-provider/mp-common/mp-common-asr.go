package mp_common

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

// --- openapi request ---

type AsrReq struct {
	File   *multipart.FileHeader `form:"file" json:"file" validate:"required"`
	Config AsrConfigOut          `form:"config" json:"config" validate:"required"`
	ApiKey string                `json:"api_key"`
}

type AsrConfigOut struct {
	Config AsrConfig `form:"config" json:"config" validate:"required"`
}

type AsrConfig struct {
	SessionId           string  `json:"session_id" validate:"required"`
	AddPunc             int     `json:"add_punc,omitempty"`
	ItnSwitch           int     `json:"itn_switch,omitempty"`
	VadSwitch           int     `json:"vad_switch,omitempty"`
	Diarization         int     `json:"diarization,omitempty"`
	SpkNum              int     `json:"spk_num,omitempty"`
	Translate           int     `json:"translate,omitempty"`
	Sensitive           int     `json:"sensitive,omitempty"`
	Language            int     `json:"language,omitempty"`
	AudioClassification int     `json:"audio_classification,omitempty"`
	DiarizationMode     int     `json:"diarization_mode,omitempty"`
	MaxEndSil           int     `json:"max_end_sil,omitempty"`
	MaxSingleSeg        int     `json:"max_single_seg,omitempty"`
	SpeechNoiseThres    float64 `json:"speech_noise_thres,omitempty"`
}

func (req *AsrReq) Check() error {
	return nil
}

func (req *AsrReq) Data() (map[string]interface{}, error) {
	m := make(map[string]interface{})
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// --- openapi response ---

type AsrResp struct {
	Status  string    `json:"status"`
	Code    int       `json:"code"`
	Message string    `json:"msg"`
	Uuid    string    `json:"uuid"`
	Result  AsrResult `json:"result"`
}
type AsrResult struct {
	Diarization []DiarizationObj `json:"diarization"`
}
type DiarizationObj struct {
	Start   float32 `json:"start"`
	End     float32 `json:"end"`
	Speaker int     `json:"speaker"`
	Text    string  `json:"text"`
	Trans   string  `json:"trans"`
}

// --- request ---

type IAsrReq interface {
	Data() *AsrReq
}

// asrReq implementation of IAsrReq
type asrReq struct {
	data *AsrReq
}

func NewAsrReq(data *AsrReq) IAsrReq {
	return &asrReq{data: data}
}

func (req *asrReq) Data() *AsrReq {
	return req.data
}

// --- response ---

type IAsrResp interface {
	String() string
	Data() (interface{}, bool)
	ConvertResp() (*AsrResp, bool)
}

// asrResp implementation of IAsrResp
type asrResp struct {
	raw string
}

func NewAsrResp(raw string) IAsrResp {
	return &asrResp{raw: raw}
}

func (resp *asrResp) String() string {
	return resp.raw
}

func (resp *asrResp) Data() (interface{}, bool) {
	ret := make(map[string]interface{})
	if err := json.Unmarshal([]byte(resp.raw), &ret); err != nil {
		log.Errorf("asr resp (%v) convert to data err: %v", resp.raw, err)
		return nil, false
	}
	return ret, true
}

func (resp *asrResp) ConvertResp() (*AsrResp, bool) {
	var ret *AsrResp
	if err := json.Unmarshal([]byte(resp.raw), &ret); err != nil {
		log.Errorf("asr resp (%v) convert to data err: %v", resp.raw, err)
		return nil, false
	}

	log.Infof("asr resp: %v", resp.raw)
	if err := util.Validate(ret); err != nil {
		log.Errorf("asr resp validate err: %v", err)
		return nil, false
	}
	return ret, true
}

// --- asr ---

func Asr(ctx *gin.Context, provider, apiKey, url string, req *AsrReq, headers ...Header) ([]byte, error) {
	if apiKey != "" {
		headers = append(headers, Header{
			Key:   "Authorization",
			Value: "Bearer " + apiKey,
		})
	}
	file, err := req.File.Open()
	if err != nil {
		return nil, fmt.Errorf("request %v %v asr err: %v", url, provider, err)
	}
	configJSON, err := json.Marshal(req.Config)
	if err != nil {
		return nil, fmt.Errorf("marshal config failed: %v", err)
	}
	defer file.Close()
	request := resty.New().
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}). // 关闭证书校验
		SetTimeout(0).                                             // 关闭请求超时
		R().
		SetContext(ctx).
		SetHeader("Content-Type", "multipart/form-data").
		SetHeader("Accept", "application/json").
		SetFileReader("file", req.File.Filename, file).
		SetMultipartField("config", "", "application/json", strings.NewReader(string(configJSON))).
		SetDoNotParseResponse(true)
	for _, header := range headers {
		request.SetHeader(header.Key, header.Value)
	}
	resp, err := request.Post(url)
	if err != nil {
		return nil, fmt.Errorf("request %v %v asr err: %v", url, provider, err)
	}
	b, err := io.ReadAll(resp.RawResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("request %v %v asr read response body err: %v", url, provider, err)
	}
	if resp.StatusCode() >= 300 {
		return nil, fmt.Errorf("request %v %v asr http status %v msg: %v", url, provider, resp.StatusCode(), string(b))
	}

	return b, nil
}
