package mp_yuanjing

import (
	"net/url"

	mp_common "github.com/UnicomAI/wanwu/pkg/model-provider/mp-common"
	"github.com/gin-gonic/gin"
)

type Text2Image struct {
	ApiKey      string `json:"apiKey"`      // ApiKey
	EndpointUrl string `json:"endpointUrl"` // 推理url
}

func (cfg *Text2Image) Tags() []mp_common.Tag {
	tags := []mp_common.Tag{
		{
			Text: mp_common.TagText2Image,
		},
	}
	return tags
}
func (cfg *Text2Image) NewReq(req *mp_common.Text2ImageReq) (mp_common.IText2ImageReq, error) {
	m, err := req.Data()
	if err != nil {
		return nil, err
	}
	return mp_common.NewText2ImageReq(m), nil
}

func (cfg *Text2Image) Text2Image(ctx *gin.Context, req mp_common.IText2ImageReq, headers ...mp_common.Header) (mp_common.IText2ImageResp, error) {
	b, err := mp_common.Text2Image(ctx, "yuanjing", cfg.ApiKey, cfg.text2ImageUrl(), req.Data(), headers...)
	if err != nil {
		return nil, err
	}
	return mp_common.NewText2ImageResp(string(b)), nil
}

func (cfg *Text2Image) text2ImageUrl() string {
	ret, _ := url.JoinPath(cfg.EndpointUrl, "")
	return ret
}
