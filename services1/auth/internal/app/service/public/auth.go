package public

import (
	"context"
	pb "github.com/donglei1234/platform/services/proto/gen/auth/api"

	"github.com/donglei1234/platform/services/auth/tools"
	fb "github.com/huandu/facebook"
	"go.uber.org/zap"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

const (
	PlatformTypeGoogle   = "google"
	PlatformTypeFacebook = "facebook"
	DebugToken           = "/debug_token"
	InputToken           = "input_token"
	MethodValidateToken  = "/auth.pb.AuthPublic/ValidateToken"
	MethodAuthenticate   = "/auth.pb.AuthPublic/Authenticate"
)

type FacebookUserInfo struct {
	Data struct {
		AppID               string `json:"app_id"`
		Type                string `json:"type"`
		Application         string `json:"application"`
		DataAccessExpiresAt int    `json:"data_access_expires_at"`
		ExpiresAt           int    `json:"expires_at"`
		IsValid             bool   `json:"is_valid"`
		Metadata            struct {
			AuthType string `json:"auth_type"`
		} `json:"metadata"`
		Scopes []string `json:"scopes"`
		UserID string   `json:"user_id"`
	} `json:"data"`
}

func (s *Service) Authenticate(
	ctx context.Context,
	req *pb.AuthenticateRequest,
) (response *pb.AuthenticateResponse, err error) {
	resp := &pb.AuthenticateResponse{}
	if user, err := s.doc.LoadOrCreateUserAuth(req.AppId, req.Username); err != nil {
		s.logger.Error("get id error", zap.Error(err))
		return nil, ErrGeneralFailure
	} else if jwtToken, err := tools.CreatProfileJwt(user.ProfileId, s.jwtSecret); err != nil {
		s.logger.Error("creat jwt failed", zap.Error(err))
		return nil, ErrGenerateTokenFailure
	} else {
		sess := &pb.Session{
			UserId: user.ProfileId,
			Token:  jwtToken,
		}
		resp.Session = sess
	}

	return resp, nil
}

func (s *Service) AuthenticateRoom(ctx context.Context, request *pb.AuthenticateRoomRequest) (*pb.AuthenticateRoomResponse, error) {
	uid, err := tools.GenerateUUID()
	if err != nil {
		s.logger.Error("generate uuid failed", zap.Error(err))
		return nil, ErrGenerateTokenFailure
	}
	request.RoomInfo.RoomId = uid
	jwtToken, err := tools.CreateRoomToken(
		request.RoomInfo, s.jwtSecret,
	)
	if err != nil {
		s.logger.Error("creat jwt failed", zap.Error(err))
		return nil, ErrGenerateTokenFailure
	}
	return &pb.AuthenticateRoomResponse{Token: jwtToken}, nil
}

func (s *Service) ValidateRoomToken(ctx context.Context, request *pb.ValidateRoomTokenRequest) (*pb.ValidateRoomTokenResponse, error) {
	roomInfo, err := tools.ParseRoomToken(request.GetToken(), s.jwtSecret)
	if err != nil {
		s.logger.Error("parse room token failed", zap.Error(err))
		return nil, ErrParseJwtTokenFailure
	}
	return &pb.ValidateRoomTokenResponse{RoomInfo: roomInfo}, nil
}

func (s *Service) Bind(ctx context.Context, req *pb.BindRequest) (response *pb.BindResponse, err error) {
	// 1、获取Token
	token, err := tools.ParseJwtToken(ctx)
	if err != nil {
		s.logger.Error("bind parse jwt token err", zap.Error(err))
		return nil, ErrParseJwtTokenFailure
	}
	// 2、解析token，获取uid
	profileId, err := tools.DecodeToken(token)
	if err != nil {
		s.logger.Error("bind jwt token decode to uid err", zap.Error(err))
		return nil, InvalidProfileId
	}

	// 2、获取AccessToken
	var platformUid string
	var platformType string
	switch req.Token.(type) {
	case *pb.BindRequest_FacebookToken:
		platformType = PlatformTypeFacebook

		// 获取FaceBookID
		accessToken := req.GetFacebookToken()

		app := fb.New(s.facebookAppId, s.facebookSecret)
		token := app.AppAccessToken()

		session := app.Session(token)
		resp, err := session.Get(DebugToken, fb.Params{
			InputToken: accessToken,
		})
		if err != nil {
			s.logger.Error("get facebook info error", zap.Error(err))
			return nil, ErrGetFacebookInfoFailure
		}

		var userInfo FacebookUserInfo
		err = resp.Decode(&userInfo)
		if err != nil {
			s.logger.Error("decode facebook info error", zap.Error(err))
			return nil, ErrDecodeFacebookInfoFailure
		}
		if !userInfo.Data.IsValid {
			s.logger.Error("facebook info invalid", zap.Error(err))
			return nil, ErrFacebookInvalid
		}

		platformUid = userInfo.Data.UserID
		//platformUid, err = tools.GetFacebookId(s.facebookAppId, s.facebookSecret, accessToken)
		//if err != nil {
		//	s.logger.Error("get facebook id error", zap.Error(err))
		//	return nil, ErrGeneralFailure
		//}
	case *pb.BindRequest_GoogleToken:
		platformType = PlatformTypeGoogle

		// 获取GoogleID
		idToken := req.GetGoogleToken()

		oauth2Service, err := oauth2.NewService(nil, option.WithAPIKey(s.googleSecret))
		if err != nil {
			s.logger.Error("oauth2 init error", zap.Error(err))
			return nil, ErrInitGoogleServiceFailure
		}
		tokenInfo := oauth2Service.Tokeninfo()
		token := tokenInfo.IdToken(idToken)
		do, err := token.Do()
		if err != nil {
			s.logger.Error("get google info error", zap.Error(err))
			return nil, ErrGetGoogleInfoFailure
		}
		// TODO 需要对不同的app验证Audience字段，如tata为743906483241-idjfg26hdialusrrf10v7l3ddkpf32jk.apps.googleusercontent.com
		// 参考：https://developers.google.com/identity/sign-in/web/backend-auth#verify-the-integrity-of-the-id-token
		if do.Audience == "" {
			s.logger.Error("google info audience error", zap.Error(err))
			return nil, ErrGoogleAudienceInvalid
		}
		if do.ExpiresIn <= 0 {
			s.logger.Error("google token expired !", zap.Error(err))
			return nil, ErrGoogleTokenExpired
		}
		platformUid = do.UserId

		//platformUid, err = tools.GetGoogleId(s.googleSecret, idToken)
		//if err != nil {
		//	s.logger.Error("get google id error", zap.Error(err))
		//	return nil, ErrGeneralFailure
		//}
	}

	auth, err := s.doc.LoadOrBindUserAuth(req.AppId, platformType, platformUid, profileId)
	if err != nil {
		s.logger.Error("load or bind user auth error", zap.Error(err))
		return nil, ErrGeneralFailure
	}

	resp := &pb.BindResponse{
		Uid:       platformUid,
		IsReLogin: false,
	}
	if auth.ProfileId != profileId {
		resp.IsReLogin = true
	}
	return resp, nil
}

func (s *Service) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (response *pb.ValidateTokenResponse, err error) {
	jwtToken := req.JwtToken
	//返回三种报错：1.超时 2.格式错误 3.其他错误
	uid, err := tools.ParseProfileToken(jwtToken, s.jwtSecret)
	if err != nil {
		return nil, err
	}
	sess := &pb.Session{
		Token:  jwtToken,
		UserId: uid,
	}
	return &pb.ValidateTokenResponse{Session: sess}, nil
}

func (s *Service) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	if fullMethodName == MethodValidateToken || fullMethodName == MethodAuthenticate {
		return ctx, nil
	}
	s.logger.Debug("authfun run , client is calling method:" + fullMethodName)
	token, err := tools.ParseJwtToken(ctx)
	if err != nil {
		return ctx, ErrMiddlewareParseJwtFailure
	}

	_, err = s.ValidateToken(ctx, &pb.ValidateTokenRequest{
		JwtToken: token,
	})
	// ValidateToken 中已经对err进行封装
	if err != nil {
		return ctx, err
	}
	s.logger.Debug("platform middleware authfun run over , check passed !")
	return ctx, nil
}
