package public

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"go.uber.org/zap"

	pb "github.com/donglei1234/platform/services/notification/generated/grpc/go/notification/api"
)

type message struct {
	APNS        string `json:"APNS"`
	APNSSandbox string `json:"APNS_SANDBOX"`
	GCM         string `json:"GCM"`
}

type Data struct {
	Sound *string     `json:"sound,omitempty"`
	Data  interface{} `json:"custom_data"`
	Badge *int        `json:"badge,omitempty"`
}

type iosPush struct {
	APS Data `json:"aps"`
}

type gcmPush struct {
	Custom interface{} `json:"custom"`
	Badge  *int        `json:"badge,omitempty"`
}

type gcmPushWrapper struct {
	Data gcmPush `json:"data"`
}

// 消息封装
func newMessageJSON(data *Data) (m string, err error) {
	b, err := json.Marshal(iosPush{
		APS: *data,
	})
	if err != nil {
		return
	}
	payload := string(b)

	b, err = json.Marshal(gcmPushWrapper{
		Data: gcmPush{
			Custom: data.Data,
			Badge:  data.Badge,
		},
	})
	if err != nil {
		return
	}
	gcm := string(b)

	pushData, err := json.Marshal(message{
		APNS:        payload,
		APNSSandbox: payload,
		GCM:         gcm,
	})
	if err != nil {
		return
	}
	m = string(pushData)
	return
}

func (s *Service) RegisterArn(ctx context.Context, request *pb.RegisterArnRequest) (*pb.NothingResponse, error) {
	s.logger.Info("register arn")
	profileId := request.GetProfileId().ProfileId //用户ID
	//创建终端节点并订阅默认主题,同一设备不同token时会覆盖之前的数据
	applicationArn := s.awsSNSApplicationARN[request.GetAppId()][request.GetDeviceType()]
	if err := s.CreateEndPointArnAndStore(applicationArn, request.GetDeviceToken(),
		profileId, request.GetAppId(), request.GetDeviceId(), request.GetRegion(),
		request.GetDeviceType()); err != nil {
		s.logger.Error("crate endpoint arn err ")
		return &pb.NothingResponse{}, err
	}
	s.logger.Info("register arn success ")
	return &pb.NothingResponse{}, nil
}

func (s *Service) PublishMessage(ctx context.Context, request *pb.PublishMessageRequest) (*pb.NothingResponse, error) {
	go s.GoPublish(ctx, request)
	return &pb.NothingResponse{}, nil
}

func (s *Service) DeleteArn(ctx context.Context, request *pb.DeleteArnRequest) (*pb.NothingResponse, error) {
	var deleteArnOutPut string
	if request.GetDelType() == pb.DeleteArnRequest_EndPointArn { //删除设备终端令牌
		//获取endPointArn
		userInfo, err := s.GetUserInfoByProfileId(request.GetPublishId(), request.GetAppId())
		if err != nil {
			s.logger.Error("get user info by profile id err", zap.Error(err))
			return &pb.NothingResponse{}, err
		}
		//删除本地数据库存储的相关数据
		if err := s.db.DelUserInfoByProfileId(request.GetPublishId(),
			request.GetAppId(), request.GetDeviceID()); err != nil {
			s.logger.Error("del user info by profile id err", zap.Error(err))
			return &pb.NothingResponse{}, err
		}
		//删除aws服务的endPointArn
		deleteEndPointOutPut, err := s.sns.DeleteEndpoint(&sns.DeleteEndpointInput{
			EndpointArn: aws.String(userInfo["endPointArn"]),
		})
		if err != nil {
			s.logger.Error("aws delete end point err", zap.Error(err))
			return &pb.NothingResponse{}, err
		}
		deleteArnOutPut = deleteEndPointOutPut.String()
	} else { //删除主题
		//通过topicName获取topicArn
		topicArn, err := s.db.GetTopicArn(request.GetAppId(), request.GetPublishId())
		if err != nil {
			s.logger.Error("delete topic:get topic arn err ", zap.Error(err))
			return &pb.NothingResponse{}, nil
		}
		//删除本地数据库存储的相关数据
		if err := s.db.DelTopicArn(request.GetAppId(), request.GetPublishId()); err != nil {
			s.logger.Error("delete topic arn err ", zap.Error(err))
			return &pb.NothingResponse{}, nil
		}
		//删除aws服务的topic
		deleteTopicOutPut, err := s.sns.DeleteTopic(&sns.DeleteTopicInput{
			TopicArn: aws.String(topicArn),
		})
		if err != nil {
			s.logger.Error("aws delete topic err ", zap.Error(err))
			return &pb.NothingResponse{}, nil
		}
		deleteArnOutPut = deleteTopicOutPut.String()
	}
	s.logger.Info("delete arn success ",
		zap.String("publish id ", request.GetPublishId()),
		zap.String("delete end point output ", deleteArnOutPut))
	return &pb.NothingResponse{}, nil
}

func (s *Service) SubscribeTopic(ctx context.Context, request *pb.SubscribeTopicRequest) (*pb.NothingResponse, error) {
	var endPointArn string
	if request.GetSubType() == pb.SubscribeTopicRequest_External { //订阅aws topic
		profileId := request.GetProfileId().ProfileId
		//通过ProfileId获取userInfo
		deviceId, err := s.db.GetDeviceIdByProfileId(profileId, request.GetAppId())
		if err != nil {
			s.logger.Error("get device id by profile id err", zap.Error(err))
			return &pb.NothingResponse{}, err
		}
		userInfo, err := s.GetUserInfoByProfileId(profileId, request.GetAppId())
		if err != nil {
			s.logger.Error("get endpoint arn by profile id err", zap.Error(err))
			return &pb.NothingResponse{}, err
		}
		if _, ok := userInfo["endPointArn"]; !ok { //因某种原因没得到endPointArn,则取出deviceToken再次创建终端节点并存储
			applicationArn := s.awsSNSApplicationARN[request.GetAppId()][userInfo["deviceType"]]
			if err := s.CreateEndPointArnAndStore(applicationArn, userInfo["deviceToken"],
				profileId, request.GetAppId(), deviceId, userInfo["region"],
				userInfo["deviceType"]); err != nil {
				return &pb.NothingResponse{}, err
			}
			userInfo, err = s.GetUserInfoByProfileId(profileId, request.GetAppId())
		}
		endPointArn = userInfo["endPointArn"]
		//判断topicName是否已创建，若已创建则取出arn,若未创建则创建之后再订阅
		//通过topicName获取topicArn
		topicArn, err := s.db.GetTopicArn(request.GetAppId(), request.GetTopicName())
		if err != nil { //数据库中未查询到topicArn
			s.logger.Error("get topic arn by topic name err", zap.Error(err))
		}
		if topicArn == "" {
			//创建topic
			if topicArn, err = s.CreateTopicArnAndStore(request.GetAppId(), request.GetTopicName()); err != nil {
				return &pb.NothingResponse{}, err
			}
			topicArn, err = s.db.GetTopicArn(request.GetAppId(), request.GetTopicName())
		}
		//直接在数据库中查询到该topicName的topicArn
		subscribeOutPut, err := s.sns.Subscribe(&sns.SubscribeInput{
			Endpoint: aws.String(endPointArn),
			TopicArn: aws.String(topicArn),
			Protocol: aws.String("application"),
		})
		if err != nil {
			s.logger.Error("subscribe err", zap.Error(err))
			return &pb.NothingResponse{}, err
		}
		s.logger.Info("subscribe success", zap.String("subscription arn", *subscribeOutPut.SubscriptionArn))
	} else { //TODO 订阅内部topic，返回流式grpc

	}
	return &pb.NothingResponse{}, nil
}

func (s *Service) GetUserInfoByProfileId(profileId, appId string) (userInfo map[string]string, err error) {
	deviceId, err := s.db.GetDeviceIdByProfileId(profileId, appId)
	if err != nil {
		s.logger.Error("get device id by profile id err", zap.Error(err))
		return userInfo, err
	}
	userInfo, err = s.db.GetUserInfoByDeviceId(profileId, appId, deviceId)
	if err != nil {
		s.logger.Error("get user info by device id err", zap.Error(err))
		return userInfo, err
	}
	return userInfo, nil
}

func (s *Service) PubMsg(message []byte, arn string) error {
	data := &Data{
		Sound: aws.String("default"),
		Badge: aws.Int(1),
		Data:  message,
	}
	msg, err := newMessageJSON(data)
	if err != nil {
		return err
	}
	s.logger.Info(msg)
	input := &sns.PublishInput{
		Message:          aws.String(msg),
		MessageStructure: aws.String("json"),
		TargetArn:        aws.String(arn),
	}
	publish, err := s.sns.Publish(input)
	if err != nil {
		s.logger.Error("sns.publish err", zap.Error(err))
		return err
	}
	s.logger.Info("publish success ", zap.String("message id : ", *publish.MessageId))
	return nil
}

func (s *Service) CreateEndPointArnAndStore(applicationArn, deviceToken,
	profileId, appId, deviceId, region, deviceType string) error {
	//创建终端节点
	resp, err := s.sns.CreatePlatformEndpoint(&sns.CreatePlatformEndpointInput{
		PlatformApplicationArn: aws.String(applicationArn),
		Token:                  aws.String(deviceToken),
	})
	if err != nil {
		s.logger.Error("create platform endpoint err ", zap.Error(err))
		return err
	}
	endPointArn := *(resp.EndpointArn)
	s.logger.Info("endpoint arn", zap.String("endpoint arn", endPointArn))
	//订阅默认主题
	subResponse, err := s.sns.Subscribe(&sns.SubscribeInput{
		Endpoint: aws.String(endPointArn),
		TopicArn: aws.String(s.topicArn[appId]),
		Protocol: aws.String("application"),
	})
	if err != nil {
		s.logger.Error("subscribe topic err")
		return err
	}
	s.logger.Info("subscribe success", zap.String("subscription arn ", *subResponse.SubscriptionArn))
	if err = s.db.SetDeviceIdAndUserInfo(profileId, appId, deviceId,
		deviceToken, region, deviceType, endPointArn); err != nil {
		return err
	}
	return nil
}

func (s *Service) CreateTopicArnAndStore(appId, topicName string) (topicArn string, err error) {
	createTopicOutput, err := s.sns.CreateTopic(&sns.CreateTopicInput{
		Name: aws.String(topicName),
	})
	if err != nil {
		s.logger.Error("aws create topic err", zap.Error(err))
		return "", err
	}
	if err := s.db.SetTopicArn(appId, topicName, *createTopicOutput.TopicArn); err != nil {
		s.logger.Error("set topic arn err", zap.Error(err))
		return "", err
	}
	return *createTopicOutput.TopicArn, nil
}

func (s *Service) GoPublish(ctx context.Context, request *pb.PublishMessageRequest) {
	var (
		arn string
		err error
	)
	//判断推送类型
	if request.GetPubType() == pb.PublishMessageRequest_Specific { //给特定用户推送消息
		profileId := request.GetPublishId()
		//通过ProfileId获取userInfo
		deviceId, err := s.db.GetDeviceIdByProfileId(profileId, request.GetAppId())
		if err != nil {
			s.logger.Error("get device id by profile id err", zap.Error(err))
			return
		}
		userInfo, err := s.db.GetUserInfoByDeviceId(profileId, request.GetAppId(), deviceId)
		if err != nil {
			s.logger.Error("get user info by device id err", zap.Error(err))
			return
		}
		if _, ok := userInfo["endPointArn"]; !ok { //因某种原因没得到endPointArn,则取出deviceToken再次创建终端节点并存储
			applicationArn := s.awsSNSApplicationARN[request.GetAppId()][userInfo["deviceType"]]
			if err := s.CreateEndPointArnAndStore(applicationArn, userInfo["deviceToken"],
				profileId, request.GetAppId(), deviceId, userInfo["region"],
				userInfo["deviceType"]); err != nil {
				return
			}
			userInfo, err = s.db.GetUserInfoByDeviceId(profileId, request.GetAppId(), deviceId)
		}
		arn = userInfo["endPointArn"]
	} else if request.GetPubType() == pb.PublishMessageRequest_Topic { //给主题推送消息
		topicName := request.GetPublishId()
		if topicName == "default" { //默认主题，从服务中取
			arn = s.topicArn[request.GetAppId()]
		} else {
			if arn, err = s.db.GetTopicArn(request.GetAppId(), topicName); err != nil {
				s.logger.Error("get topic name err", zap.Error(err))
				return
			}
			if arn == "" { //数据库中不存在该topicName,创建新主题
				if arn, err = s.CreateTopicArnAndStore(request.GetAppId(), topicName); err != nil {
					s.logger.Error("create topic err", zap.Error(err))
					return
				}
			}
		}
	}
	//开始推送消息
	err = s.PubMsg(request.GetMessage(), arn)
	if err != nil {
		s.logger.Error("publish message err", zap.Error(err))
		return
	}
	s.logger.Info("publish message success",
		zap.String("message : ", string(request.GetMessage())),
		zap.String("publish id : ", request.GetPublishId()),
	)
}
