package operate

import (
	"context"

	"github.com/UnicomAI/wanwu/api/proto/common"
	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	operate_service "github.com/UnicomAI/wanwu/api/proto/operate-service"
	"github.com/UnicomAI/wanwu/internal/operate-service/client/orm"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) AddClientRecord(ctx context.Context, req *operate_service.AddClientRecordReq) (*emptypb.Empty, error) {
	if err := s.cli.AddClientRecord(ctx, req.ClientId); err != nil {
		return nil, errStatus(errs.Code_OperateRecord, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) GetClientStatistic(ctx context.Context, req *operate_service.GetClientStatisticReq) (*operate_service.ClientStatistic, error) {
	stats, err := s.cli.GetClientStatistic(ctx, req.StartDate, req.EndDate)
	if err != nil {
		return nil, errStatus(errs.Code_OperateRecord, err)
	}
	return &operate_service.ClientStatistic{
		Overview: &operate_service.ClientOverViewInfo{
			Cumulative: convertClientOverview(stats.Overview.Cumulative),
			New:        convertClientOverview(stats.Overview.New),
			Active:     convertClientOverview(stats.Overview.Active),
		},
		Trend: &operate_service.ClientTrendInfo{
			Client: convertStatisticChart(stats.Trend.Client),
		},
	}, nil
}

// --- internal ---

func convertClientOverview(stats orm.ClientOverviewItem) *operate_service.ClientOverviewItem {
	return &operate_service.ClientOverviewItem{
		Value:            stats.Value,
		PeriodOverperiod: stats.PeriodOverPeriod,
	}
}

func convertStatisticChart(chart orm.StatisticChart) *common.StatisticChart {
	pbChart := &common.StatisticChart{
		TableName:  chart.Name,
		ChartLines: make([]*common.StatisticChartLine, 0, len(chart.Lines)),
	}
	for _, line := range chart.Lines {
		pbLine := &common.StatisticChartLine{
			LineName: line.Name,
			Items:    make([]*common.StatisticChartLineItem, 0, len(line.Items)),
		}
		for _, respItem := range line.Items {
			pbLine.Items = append(pbLine.Items, &common.StatisticChartLineItem{
				Key:   respItem.Key,
				Value: respItem.Value,
			})
		}
		pbChart.ChartLines = append(pbChart.ChartLines, pbLine)
	}
	return pbChart
}
