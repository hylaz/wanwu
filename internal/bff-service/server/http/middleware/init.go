package middleware

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/config"
	mid "github.com/UnicomAI/wanwu/pkg/gin-util/mid-wrap"
	"github.com/UnicomAI/wanwu/pkg/gin-util/route"
	jwt_util "github.com/UnicomAI/wanwu/pkg/jwt-util"
)

func Init() {
	jwt_util.InitUserJWT(config.Cfg().JWT.SigningKey)

	mid.InitWrapper(Record)

	// --- openapi ---
	mid.NewSub("openapi", "对外提供原子能力", route.PermNone, false, false)

	// --- callback ---
	mid.NewSub("callback", "系统内部调用", route.PermNone, false, false)

	// --- openurl ---
	mid.NewSub("openurl", "智能体Url", route.PermNone, false, false)

	// --- guest ---
	mid.NewSub("guest", "", route.PermNone, false, false)

	// --- common ---
	mid.NewSub("common", "", route.PermNeedEnable, false, false, JWTUser, CheckUserEnable)

	// --- model ---
	mid.NewSub("model", "模型管理", route.PermNeedCheck, true, true, JWTUser, CheckUserPerm)

	// --- knowledge ---
	mid.NewSub("knowledge", "知识库", route.PermNeedCheck, true, true, JWTUser, CheckUserPerm)

	// --- mcp ---
	mid.NewSub("mcp", "MCP广场", route.PermNeedCheck, true, true, JWTUser, CheckUserPerm)

	// --- tool ---
	mid.NewSub("tool", "资源库", route.PermNeedCheck, true, true, JWTUser, CheckUserPerm)

	// --- safety ---
	mid.NewSub("safety", "安全护栏", route.PermNeedCheck, true, true, JWTUser, CheckUserPerm)

	// --- rag ---
	mid.NewSub("rag", "文本问答", route.PermNeedCheck, true, true, JWTUser, CheckUserPerm)

	// --- workflow ---
	mid.NewSub("workflow", "工作流", route.PermNeedCheck, true, true, JWTUser, CheckUserPerm)

	// --- agent ---
	mid.NewSub("agent", "智能体", route.PermNeedCheck, true, true, JWTUser, CheckUserPerm)

	// --- exploration ---
	mid.NewSub("exploration", "应用广场", route.PermNeedCheck, true, true, JWTUser, CheckUserPerm)

	// --- permission ---
	mid.NewSub("permission", "组织管理", route.PermNeedCheck, true, true, JWTUser, CheckUserPerm)

	// permission.user
	mid.Sub("permission").NewSub("user", "用户", route.PermNeedCheck, true, true)

	// permission.org
	mid.Sub("permission").NewSub("org", "组织", route.PermNeedCheck, true, true)

	// permission.role
	mid.Sub("permission").NewSub("role", "角色", route.PermNeedCheck, true, true)

	// --- setting ---
	mid.NewSub("setting", "平台配置", route.PermNeedCheck, true, true, JWTUser, CheckUserPerm)

	// --- statistic_client ---
	mid.NewSub("statistic_client", "统计分析", route.PermNeedCheck, true, true, JWTUser, CheckUserPerm)

	// --- oauth ---
	mid.NewSub("oauth", "OAuth密钥管理", route.PermNeedCheck, true, true, JWTUser, CheckUserPerm)
}
