package public

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediaconvert"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	pb "github.com/donglei1234/platform/services/storage/generated/grpc/go/storage/api"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"regexp"
	"strings"
)

func (s *Service) UploadFile(ctx context.Context, request *pb.UploadRequest) (*pb.NothingResponse, error) {
	appId := request.AppId
	fileName := request.FileName
	profileId := request.ProfileId
	var aclFlag string
	var key string
	// appId/profileId/key_a0d99f20-1dd1-459b-b516-dfeca4005203
	r, _ := regexp.Compile("[^0-9a-zA-Z_]")
	fileNameAfterRegex := r.ReplaceAllString(fileName, "_") + "_" + uuid.New().String()
	if profileId == "" {
		aclFlag = "authenticated-read"
		key = makeFilePath(appId, "config", fileNameAfterRegex)
	} else {
		aclFlag = "public-read"
		key = makeFilePath(appId, profileId, fileNameAfterRegex)
	}

	_, err := s.db.SavePathKey(profileId, appId, fileName, key)
	if err != nil {
		s.logger.Error("save path key failed!", zap.Any("err", err))
	}
	uploader := s3manager.NewUploader(s.sess)
	file := bytes.NewReader(request.Content)
	_, err = uploader.Upload(&s3manager.UploadInput{
		ACL:    &aclFlag,
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		s.logger.Error("Unable to upload.",
			zap.Any("filePath", key),
			zap.Any("bucketName", s.bucketName))
		return &pb.NothingResponse{}, err
	}
	return &pb.NothingResponse{}, nil
}

func (s *Service) DeleteFile(ctx context.Context, request *pb.DeleteRequest) (*pb.NothingResponse, error) {
	bucket := s.bucketName
	s3StorageKey, _ := s.db.GetPathByKey(request.ProfileId, request.AppId, request.FileName)
	s.logger.Info("get path by key.",
		zap.Any("profiledId", request.ProfileId),
		zap.Any("appId", request.AppId),
		zap.Any("key", s3StorageKey))
	client := s3.New(s.sess)
	_, err := client.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(s3StorageKey)})
	if err != nil {
		s.logger.Error("Unable to delete object from bucket",
			zap.Any("bucketName", bucket),
			zap.Any("key", s3StorageKey))
		return &pb.NothingResponse{}, err
	}

	err = client.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(s3StorageKey),
	})

	return &pb.NothingResponse{}, nil
}

func (s *Service) GetFileUrl(ctx context.Context, request *pb.DownloadRequest) (*pb.DownloadResponse, error) {
	s3StorageKey, _ := s.db.GetPathByKey(request.ProfileId, request.AppId, request.FileName)
	s.logger.Info("get path by key.",
		zap.Any("profiledId", request.ProfileId),
		zap.Any("appId", request.AppId),
		zap.Any("key", request.FileName))

	url := "https://" + s.bucketName + ".s3." + *s.sess.Config.Region + ".amazonaws.com/" + s3StorageKey
	return &pb.DownloadResponse{
		Url: url,
	}, nil
}

func (s *Service) GetFileContent(ctx context.Context, request *pb.DownloadRequest) (*pb.GetFileContentResponse, error) {
	bucket := s.bucketName
	s3StorageKey, _ := s.db.GetPathByKey(request.ProfileId, request.AppId, request.FileName)
	if s3StorageKey == "" {
		s3StorageKey = request.AppId + "/config/" + request.FileName
	}
	s.logger.Info("get path by key.",
		zap.Any("profiledId", request.ProfileId),
		zap.Any("appId", request.AppId),
		zap.Any("key", request.FileName))

	downloader := s3manager.NewDownloader(s.sess)
	buffer := aws.NewWriteAtBuffer([]byte{})
	numBytes, err := downloader.Download(buffer,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(s3StorageKey),
		})

	s.logger.Info("download result",
		zap.Any("bucketName", bucket),
		zap.Any("key", s3StorageKey),
		zap.Any("numBytes", numBytes))

	if err != nil {
		s.logger.Error("Unable to download item",
			zap.Any("bucketName", bucket),
			zap.Any("key", s3StorageKey),
			zap.Any("err", err))
	}
	return &pb.GetFileContentResponse{
		Content: buffer.Bytes(),
	}, nil
}

func (s *Service) GetFileACL(ctx context.Context, request *pb.GetACLRequest) (*pb.ACLResponse, error) {
	client := s3.New(s.sess)
	s3StorageKey, _ := s.db.GetPathByKey(request.ProfileId, request.AppId, request.FileName)
	s.logger.Info("get path by key.",
		zap.Any("profiledId", request.ProfileId),
		zap.Any("appId", request.AppId),
		zap.Any("key", request.FileName))

	result, err := client.GetObjectAcl(&s3.GetObjectAclInput{
		Bucket: &s.bucketName,
		Key:    &s3StorageKey,
	})

	if len(result.Grants) > 0 {
		aclResult := &pb.ACLResponse{
			Grantee:    *result.Grants[0].Grantee.DisplayName,
			Type:       *result.Grants[0].Grantee.Type,
			Permission: *result.Grants[0].Permission,
		}
		return aclResult, nil
	} else {
		return nil, err
	}
}

func (s *Service) SetFileACL(ctx context.Context, request *pb.SetACLRequest) (*pb.NothingResponse, error) {
	client := s3.New(s.sess)

	aclType := request.AclType
	var aclFlag string
	switch aclType {
	case 1:
		aclFlag = mediaconvert.S3ObjectCannedAclPublicRead
	case 2:
		aclFlag = mediaconvert.S3ObjectCannedAclAuthenticatedRead
	case 3:
		aclFlag = mediaconvert.S3ObjectCannedAclBucketOwnerRead
	case 4:
		aclFlag = mediaconvert.S3ObjectCannedAclBucketOwnerFullControl
	}

	s3StorageKey, _ := s.db.GetPathByKey(request.ProfileId, request.AppId, request.FileName)

	s.logger.Info("get path by key.",
		zap.Any("profiledId", request.ProfileId),
		zap.Any("appId", request.AppId),
		zap.Any("key", request.FileName))

	_, err := client.PutObjectAcl(&s3.PutObjectAclInput{
		ACL:    &aclFlag,
		Bucket: &s.bucketName,
		Key:    &s3StorageKey,
	})
	if err != nil {
		return &pb.NothingResponse{}, err
	}

	return &pb.NothingResponse{}, nil
}

func (s *Service) GetFiles(ctx context.Context, request *pb.GetFilesRequest) (*pb.GetFilesResponse, error) {
	client := s3.New(s.sess)
	bucket := s.bucketName
	resp, err := client.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})

	if err != nil {
		return &pb.GetFilesResponse{
			Data: nil,
		}, err
	}

	var data []*pb.Item
	for _, item := range resp.Contents {
		if request.ProfileId == "" {
			if strings.Contains(*item.Key, request.AppId) && strings.Contains(*item.Key, "config") {
				lastModify, _ := ptypes.TimestampProto(*item.LastModified)
				item := pb.Item{
					Name:         *item.Key,
					LastModified: lastModify,
					Size:         *item.Size,
					StorageClass: *item.StorageClass,
				}
				data = append(data, &item)
			}
		} else {
			if strings.Contains(*item.Key, request.AppId) && strings.Contains(*item.Key, request.ProfileId) {
				lastModify, _ := ptypes.TimestampProto(*item.LastModified)
				item := pb.Item{
					Name:         *item.Key,
					LastModified: lastModify,
					Size:         *item.Size,
					StorageClass: *item.StorageClass,
				}
				data = append(data, &item)
			}
		}
	}

	return &pb.GetFilesResponse{
		Data: data,
	}, nil
}

func makeFilePath(opts ...string) string {
	var filePath string
	for _, item := range opts {
		filePath = filePath + item + "/"
	}
	suffix := strings.TrimSuffix(filePath, "/")
	return suffix
}
