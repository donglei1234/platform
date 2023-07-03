package public

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"

	"github.com/donglei1234/platform/services/auth/tools"
	"github.com/donglei1234/platform/services/common/utils"
	pb "github.com/donglei1234/platform/services/iap/generated/grpc/go/iap/api"
)

func (s *Service) CheckIAPToken(ctx context.Context, request *pb.IAPRequest) (response *pb.IAPResponse, err error) {
	// 检验token并获取user id
	userToken, err := tools.ParseJwtToken(ctx)
	if err != nil {
		s.logger.Error("get user token error", zap.Any("error", err))
		return nil, ErrGetUserId
	}
	userId, err := utils.DecodeToken(userToken)
	if err != nil {
		s.logger.Error("decode user token error", zap.Any("error", err))
		return nil, ErrGetUserId
	}

	sys := request.Sys
	// 根据sys的类型判断请求url是什么平台
	switch sys {
	case pb.SYS_ANDROID:
		ctx := context.Background()

		//通过google颁发的证书进行认证
		service, err := androidpublisher.NewService(ctx, option.WithCredentialsJSON([]byte(s.iapCredential)))
		if err != nil {
			s.logger.Error("error service", zap.Any("err", err.Error()))
			return nil, ErrApproveRequest
		}

		sub := service.Purchases.Products
		call := sub.Get(request.AppStoreId, request.ProductId, request.ProductToken)
		msg, err := call.Do()
		if err != nil {
			s.logger.Error("call google play error", zap.Any("error", err.Error()))
			return nil, ErrApproveRequest
		}

		orderInfo, err := json.Marshal(msg)
		if err != nil {
			s.logger.Error("json marshal error", zap.Any("err", err.Error()))
			return response, nil
		}

		response = &pb.IAPResponse{
			Code: OKCODE,
			Msg:  OKMSG,
			Sys:  sys,
			Data: string(orderInfo),
		}

		s.addTradingRecord(userId, request, msg)

		return response, nil
	case pb.SYS_IOS:
		return nil, nil
	default:
		return nil, ErrRequestSysType
	}
}

// game server发放奖励之后请求接口
func (s *Service) ConsumeCallBack(ctx context.Context, request *pb.IAPRequest) (*pb.IAPResponse, error) {
	appName := request.AppStoreId
	token := request.ProductToken
	status, err := s.db.CheckToken(appName, token)
	if err != nil {
		s.logger.Error("check token error", zap.Any("err", err.Error()))
		return nil, ErrRequestToken
	}
	if !status {
		s.logger.Error("check token not exist", zap.Any("token", token))
		return nil, ErrRequestToken
	}
	msg, err := s.db.GetIAPTradingRecord(appName, token)
	if err != nil {
		s.logger.Error("get token error", zap.Any("token", token))
		return nil, ErrRequestToken
	}

	record := &IAPTradingRecord{}
	json.Unmarshal([]byte(msg), record)

	record.ConsumptionState = 1

	jsonLog, err := json.Marshal(record)

	s.db.UpdateIAPTradingRecord(appName, token, string(jsonLog))
	response := &pb.IAPResponse{
		Code: OKCODE,
		Msg:  OKMSG,
	}
	return response, nil
}

func (s *Service) addTradingRecord(userId string, req *pb.IAPRequest, gresponse *androidpublisher.ProductPurchase) {
	record := &IAPTradingRecord{
		UserId:           userId,
		Token:            req.ProductToken,
		PurchaseState:    gresponse.PurchaseState,
		ConsumptionState: gresponse.ConsumptionState,
		Sys:              req.Sys,
	}
	msg, err := json.Marshal(record)
	s.db.InsertIAPTradingRecord(req.AppStoreId, req.ProductToken, string(msg))
	if err != nil {
		s.logger.Error("add record error", zap.Any("error", err.Error()))
	}
}
