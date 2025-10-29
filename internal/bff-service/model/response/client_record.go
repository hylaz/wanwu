package response

type ClientStatistic struct {
	Overview ClientOverView `json:"overview"` // 客户端统计面板
	Trend    ClientTrends   `json:"trend"`    // 客户端统计趋势
}

type ClientOverView struct {
	ActiveClient   ClientOverviewItem `json:"activeClient"`   // 活跃客户端
	TotalClient    ClientOverviewItem `json:"totalClient"`    // 累计客户端
	AdditionClient ClientOverviewItem `json:"additionClient"` // 新增客户端
}

type ClientOverviewItem struct {
	Value            float32 `json:"value"`            // 数量
	PeriodOverPeriod float32 `json:"periodOverPeriod"` // 环比上周期百分比
}
