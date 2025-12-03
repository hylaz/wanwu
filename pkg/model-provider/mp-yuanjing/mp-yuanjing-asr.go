package mp_yuanjing

import (
	"net/url"

	mp_common "github.com/UnicomAI/wanwu/pkg/model-provider/mp-common"
	"github.com/gin-gonic/gin"
)

type Asr struct {
	ApiKey      string `json:"apiKey"`      // ApiKey
	EndpointUrl string `json:"endpointUrl"` // 推理url
}

func (cfg *Asr) Tags() []mp_common.Tag {
	tags := []mp_common.Tag{
		{
			Text: mp_common.TagAsr,
		},
	}
	return tags
}
func (cfg *Asr) NewReq(req *mp_common.AsrReq) (mp_common.IAsrReq, error) {
	return mp_common.NewAsrReq(req), nil
}

func (cfg *Asr) Asr(ctx *gin.Context, req mp_common.IAsrReq, headers ...mp_common.Header) (mp_common.IAsrResp, error) {
	b, err := mp_common.Asr(ctx, "yuanjing", cfg.ApiKey, cfg.asrUrl(), req.Data(), headers...)
	if err != nil {
		return nil, err
	}
	return mp_common.NewAsrResp(string(b)), nil
}

func (cfg *Asr) asrUrl() string {
	ret, _ := url.JoinPath(cfg.EndpointUrl, "")
	return ret
}
