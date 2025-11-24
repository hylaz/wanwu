package iam

import (
	"context"
	"fmt"
	"strconv"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	iam_service "github.com/UnicomAI/wanwu/api/proto/iam-service"
	"github.com/UnicomAI/wanwu/internal/iam-service/pkg/util"
	utils "github.com/UnicomAI/wanwu/pkg/util"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) GetCaptcha(ctx context.Context, req *iam_service.GetCaptchaReq) (*iam_service.GetCaptchaResp, error) {
	code, b64s, err := util.GenerateCaptcha()
	if err != nil {
		return nil, errStatus(errs.Code_IAMCaptcha, toErrStatus("iam_login_captcha", err.Error()))
	}
	if err := s.cli.RefreshCaptcha(ctx, req.Key, code); err != nil {
		return nil, errStatus(errs.Code_IAMCaptcha, err)
	}
	return &iam_service.GetCaptchaResp{
		Key: req.Key,
		B64: b64s,
	}, nil
}

func (s *Service) Login(ctx context.Context, req *iam_service.LoginReq) (*iam_service.LoginResp, error) {
	// captcha
	if err := s.cli.CheckCaptcha(ctx, req.Key, req.Code); err != nil {
		return nil, errStatus(errs.Code_IAMCaptcha, err)
	}
	// login
	user, permission, err := s.cli.Login(ctx, req.UserName, req.Password, req.Language)
	if err != nil {
		return nil, errStatus(errs.Code_IAMLogin, err)
	}
	return &iam_service.LoginResp{
		User:       toUserInfo(user),
		Permission: toPermission(permission),
	}, nil
}

func (s *Service) LoginByEmail(ctx context.Context, req *iam_service.LoginByEmailReq) (*iam_service.LoginByEmailResp, error) {
	// captcha
	if err := s.cli.CheckCaptcha(ctx, req.Key, req.Code); err != nil {
		return nil, errStatus(errs.Code_IAMCaptcha, err)
	}
	// login
	emailLoginInfo, err := s.cli.LoginByEmail(ctx, req.UserName, req.Password)
	if err != nil {
		return nil, errStatus(errs.Code_IAMLogin, err)
	}
	return &iam_service.LoginByEmailResp{
		UserId:               strconv.Itoa(int(emailLoginInfo.ID)),
		IsEmailCheck:         emailLoginInfo.IsEmailCheck,
		LastUpdatePasswordAt: emailLoginInfo.LastUpdatePasswordAt,
	}, nil
}

func (s *Service) LoginSendEmailCode(ctx context.Context, req *iam_service.LoginSendEmailCodeReq) (*emptypb.Empty, error) {
	if err := s.cli.LoginSendEmailCode(ctx, req.Email); err != nil {
		return nil, errStatus(errs.Code_IAMUser, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) LoginEmailCheck(ctx context.Context, req *iam_service.LoginEmailCheckReq) (*iam_service.LoginResp, error) {
	// login
	fmt.Println("用户ID", utils.MustU32(req.UserId))
	user, permission, err := s.cli.LoginEmailCheck(ctx, utils.MustU32(req.UserId), req.Email, req.Code, req.Language)
	if err != nil {
		return nil, errStatus(errs.Code_IAMLogin, err)
	}
	return &iam_service.LoginResp{
		User:       toUserInfo(user),
		Permission: toPermission(permission),
	}, nil
}

func (s *Service) ChangeUserPasswordByEmail(ctx context.Context, req *iam_service.ChangeUserPasswordByEmailReq) (*iam_service.LoginResp, error) {
	// login
	user, permission, err := s.cli.ChangeUserPasswordByEmail(ctx,
		utils.MustU32(req.UserId),
		req.OldPassword,
		req.NewPassword,
		req.Email,
		req.Code,
		req.Language)
	if err != nil {
		return nil, errStatus(errs.Code_IAMLogin, err)
	}
	return &iam_service.LoginResp{
		User:       toUserInfo(user),
		Permission: toPermission(permission),
	}, nil
}
