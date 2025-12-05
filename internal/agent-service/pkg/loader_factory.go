package pkg

import (
	"fmt"
	"sort"
)

var loadServiceList []*LoaderService

func AddContainer(service LoaderService) {
	loadServiceList = append(loadServiceList, &service)
}

func InitAllService() error {
	if len(loadServiceList) >= 0 {
		for _, service := range loadServiceList {
			fmt.Printf("load service %s start \n", (*service).LoadType())
			err := (*service).Load()
			if err != nil {
				fmt.Printf("load service %s error: %s \n", (*service).LoadType(), err.Error())
				return err
			} else {
				fmt.Printf("load service %s start success\n", (*service).LoadType())
			}
		}
	}
	return nil
}

func StopAllService() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in StopAllService", r)
		}
	}()
	if len(loadServiceList) >= 0 {
		sort.Stable(LoaderServiceSlice(loadServiceList))
		for _, service := range loadServiceList {
			fmt.Printf("load service %s stop ", (*service).LoadType())
			err := (*service).Stop()
			if err != nil {
				fmt.Printf("StopAllService error:%s", err.Error())
			} else {
				fmt.Printf("load service %s stop success", (*service).LoadType())
			}
		}
	}
}

type LoaderServiceSlice []*LoaderService

func (p LoaderServiceSlice) Len() int { return len(p) }
func (p LoaderServiceSlice) Less(i, j int) bool {
	return (*p[i]).StopPriority() < (*p[j]).StopPriority()
}
func (p LoaderServiceSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
