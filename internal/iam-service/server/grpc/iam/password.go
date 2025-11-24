package iam

import (
	"context"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	iam_service "github.com/UnicomAI/wanwu/api/proto/iam-service"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) ResetPasswordByEmail(ctx context.Context, req *iam_service.ResetPasswordByEmailReq) (*emptypb.Empty, error) {
	if err := s.cli.ResetPasswordByEmail(ctx, req.Email, req.Password, req.Code); err != nil {
		return nil, errStatus(errs.Code_IAMUser, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) ResetPasswordSendEmailCode(ctx context.Context, req *iam_service.ResetPasswordSendEmailCodeReq) (*emptypb.Empty, error) {
	if err := s.cli.ResetPasswordSendEmailCode(ctx, req.Email); err != nil {
		return nil, errStatus(errs.Code_IAMUser, err)
	}
	return &emptypb.Empty{}, nil
}
