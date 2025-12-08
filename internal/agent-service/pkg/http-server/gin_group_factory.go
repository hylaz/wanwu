package http_server

import (
	"github.com/gin-gonic/gin"
)

var serviceList []*GinGroup

type GinGroupService interface {
	RegisterPath() string
}

type GinGroup struct {
	GroupPath     string
	InterfaceList []*GinInterface
}

type GinInterface struct {
	RelPath     string
	method      string
	handler     gin.HandlerFunc
	desc        string
	middlewares []gin.HandlerFunc
}

func Group(path string) *GinGroup {
	ginGroup := &GinGroup{
		GroupPath: path,
	}
	serviceList = append(serviceList, ginGroup)
	return ginGroup
}

func (g *GinGroup) Register(relPath, method string, handler gin.HandlerFunc, desc string, middlewares ...gin.HandlerFunc) {
	g.InterfaceList = append(g.InterfaceList, &GinInterface{
		RelPath:     relPath,
		method:      method,
		handler:     handler,
		desc:        desc,
		middlewares: middlewares,
	})
}

func InitGinGroup(ginEngine *gin.Engine, middlewares ...gin.HandlerFunc) error {
	for _, service := range serviceList {
		group := ginEngine.Group(service.GroupPath)
		for _, ginInterface := range service.InterfaceList {
			handlerFuncs := append(middlewares, ginInterface.middlewares...)
			reg(group, ginInterface.RelPath, ginInterface.method, ginInterface.handler, handlerFuncs...)
		}
	}
	return nil
}

func reg(rg *gin.RouterGroup, relPath, method string, handler gin.HandlerFunc, middlewares ...gin.HandlerFunc) {
	// handler
	var handlers []gin.HandlerFunc
	handlers = append(handlers, middlewares...)
	handlers = append(handlers, handler)
	rg.Handle(method, relPath, handlers...)
}
