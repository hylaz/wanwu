package oauth2_util

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/UnicomAI/wanwu/pkg/redis"
	"github.com/UnicomAI/wanwu/pkg/util"
)

var (
	yamlFile string
)

type YamlConfig struct {
	Redis redis.Config `json:"redis" mapstructure:"redis"`
}

func TestMain(m *testing.M) {
	flag.StringVar(&yamlFile, "config", "../../../../configs/microservice/bff-service/configs/config.yaml", "conf yaml file")

	flag.Parse()
	os.Exit(m.Run())
}

func TestOAuth2RedisUtil(t *testing.T) {
	ctx := context.TODO()

	cfg := YamlConfig{}
	if err := util.LoadConfig(yamlFile, &cfg); err != nil {
		t.Fatal(err)
	}
	if err := redis.InitOP(ctx, cfg.Redis); err != nil {
		t.Fatal(err)
	}
	if err := Init(redis.OP().Cli()); err != nil {
		t.Fatal(err)
	}
	defer redis.OP().Stop()

	code := "code_123"
	if err := SaveCode(ctx, code, CodePayload{
		ClientID: "client_1",
		UserID:   "user_1",
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := ValidateCode(ctx, code, "client_1"); err != nil {
		t.Fatal(err)
	}

	if err := SaveRefreshToken(ctx, "refresh_123", time.Minute); err != nil {
		t.Fatal(err)
	}

	if err := ValidateRefreshToken(ctx, "refresh_123"); err != nil {
		t.Fatal(err)
	}
}
