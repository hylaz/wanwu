package operate

import (
	"context"

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

func (s *Service) GetClientOverview(ctx context.Context, req *operate_service.GetClientOverviewReq) (*operate_service.ClientOverViewInfo, error) {
	stats, err := s.cli.GetClientOverview(ctx, req.StartDate, req.EndDate)
	if err != nil {
		return nil, errStatus(errs.Code_OperateRecord, err)
	}
	return toClientOverviewInfo(stats), nil
}

func (s *Service) GetClientTrend(ctx context.Context, req *operate_service.GetClientTrendReq) (*operate_service.ClientTrendInfo, error) {
	trend, err := s.cli.GetClientTrend(ctx, req.StartDate, req.EndDate)
	if err != nil {
		return nil, errStatus(errs.Code_OperateRecord, err)
	}
	return &operate_service.ClientTrendInfo{
		Client: convertStatisticChart(trend),
	}, nil
}

func (s *Service) GetCumulativeClientOverview(ctx context.Context, req *operate_service.GetCumulativeClientOverviewReq) (*operate_service.ClientOverViewInfo, error) {
	stats, err := s.cli.GetCumulativeClientOverview(ctx, req.EndAt)
	if err != nil {
		return nil, errStatus(errs.Code_OperateRecord, err)
	}
	return toCumulativeClientOverviewInfo(stats), nil
}

func convertStatisticChart(trend *orm.ClientTrends) *operate_service.StatisticChart {
	pbChart := &operate_service.StatisticChart{
		TableName:  trend.Client.TableName,
		ChartLines: make([]*operate_service.StatisticChartLine, 0, len(trend.Client.Lines)),
	}
	for _, respLine := range trend.Client.Lines {
		pbLine := &operate_service.StatisticChartLine{
			LineName: respLine.LineName,
			Items:    make([]*operate_service.StatisticChartLineItem, 0, len(respLine.Items)),
		}
		for _, respItem := range respLine.Items {
			pbLine.Items = append(pbLine.Items, &operate_service.StatisticChartLineItem{
				Key:   respItem.Key,
				Value: respItem.Value,
			})
		}
		pbChart.ChartLines = append(pbChart.ChartLines, pbLine)
	}
	return pbChart
}

// --- internal ---
func toClientOverviewInfo(stats *orm.ClientOverView) *operate_service.ClientOverViewInfo {
	ret := &operate_service.ClientOverViewInfo{
		ActiveClient: &operate_service.ClientOverviewItem{
			Value:            stats.ActiveClient.Value,
			PeriodOverperiod: stats.ActiveClient.PeriodOverPeriod,
		},
		AdditionClient: &operate_service.ClientOverviewItem{
			Value:            stats.AdditionClient.Value,
			PeriodOverperiod: stats.AdditionClient.PeriodOverPeriod,
		},
	}
	return ret
}

func toCumulativeClientOverviewInfo(stats *orm.ClientOverView) *operate_service.ClientOverViewInfo {
	ret := &operate_service.ClientOverViewInfo{
		TotalClient: &operate_service.ClientOverviewItem{
			Value: stats.TotalClient.Value,
		},
	}
	return ret
}
