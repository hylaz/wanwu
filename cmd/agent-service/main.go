package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/UnicomAI/wanwu/internal/agent-service/pkg"
	_ "github.com/UnicomAI/wanwu/internal/agent-service/pkg/minio"
	_ "github.com/UnicomAI/wanwu/internal/agent-service/server/http/router"
)

var (
	isVersion    bool
	buildTime    string //编译时间
	buildVersion string //编译版本
	gitCommitID  string //git的commit id
	gitBranch    string //git branch
	builder      string //构建者
)

func main() {
	flag.BoolVar(&isVersion, "v", false, "编译信息")
	flag.Parse()
	if len(os.Args) > 1 && !isVersion {
		flag.Usage()
		return
	}

	if isVersion {
		versionPrint()
		return
	}
	err := pkg.InitAllService()
	if err != nil {
		panic(err)
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGTERM)
	<-sc
	pkg.StopAllService()
}

func versionPrint() {
	fmt.Printf("build_time: %s\n", buildTime)
	fmt.Printf("build_version: %s\n", buildVersion)
	fmt.Printf("git_commit_id: %s\n", gitCommitID)
	fmt.Printf("git branch: %s\n", gitBranch)
	fmt.Printf("runtime version: %s\n", runtime.Version())
	fmt.Printf("builder: %s\n", builder)
}
