package client

import (
	"context"
	"fmt"
	"time"

	"github.com/donglei1234/platform/services/common/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/donglei1234/platform/services/storage/generated/grpc/go/storage/api"
)

type StorageClient interface {
	Close() error
	GetFiles(
		ctx context.Context,
		profileId string,
		appId string) (res *pb.GetFilesRequest, err error)
	UploadFile(
		ctx context.Context,
		path []byte,
		key string,
		profileId string,
		appId string) (*pb.NothingResponse, error)
	Delete(
		ctx context.Context,
		key string,
		profileId string,
		appId string) (*pb.NothingResponse, error)
	DownloadForUrl(
		ctx context.Context,
		key string,
		profileId string,
		appId string) (*pb.DownloadResponse, error)
	DownloadForItem(
		ctx context.Context,
		key string,
		profileId string,
		appId string) (res *pb.GetFileContentResponse, err error)
	GetFileACL(
		ctx context.Context,
		key string,
		profileId string,
		appId string) (res *pb.ACLResponse, err error)
	SetFileACL(
		ctx context.Context,
		key string,
		profileId string,
		aclType int32,
		appId string) (res *pb.NothingResponse, err error)
}

type publicClient struct {
	l *zap.Logger
	*client
}

func (c *publicClient) GetFiles(ctx context.Context, profileId string, appId string) (res *pb.GetFilesRequest, err error) {
	cli := pb.NewStorageClient(c.conn)
	request := &pb.GetFilesRequest{
		ProfileId: profileId,
		AppId:     appId,
	}
	name, err := cli.GetFiles(ctx, request)
	data := name.Data
	for _, item := range data {
		fmt.Println("Name:         ", item.Name)
		fmt.Println("Last modified:", item.LastModified)
		fmt.Println("Size:         ", item.Size)
		fmt.Println("Storage class:", item.StorageClass)
	}
	return nil, nil
}

func (c *publicClient) UploadFile(ctx context.Context, path []byte, key string, profileid string, appId string) (*pb.NothingResponse, error) {
	cli := pb.NewStorageClient(c.conn)
	request := &pb.UploadRequest{
		Content:   path,
		FileName:  key,
		ProfileId: profileid,
		AppId:     appId,
	}
	return cli.UploadFile(ctx, request)
}

func (c *publicClient) Delete(ctx context.Context, key string, profileid string, appId string) (*pb.NothingResponse, error) {
	cli := pb.NewStorageClient(c.conn)
	request := &pb.DeleteRequest{
		FileName:  key,
		ProfileId: profileid,
		AppId:     appId,
	}
	return cli.DeleteFile(ctx, request)
}

func (c *publicClient) DownloadForUrl(ctx context.Context, key string, profileid string, appId string) (*pb.DownloadResponse, error) {
	cli := pb.NewStorageClient(c.conn)
	request := &pb.DownloadRequest{
		FileName:  key,
		ProfileId: profileid,
		AppId:     appId,
	}
	return cli.GetFileUrl(ctx, request)
}

func (c *publicClient) DownloadForItem(ctx context.Context, key string, profileid string, appId string) (res *pb.GetFileContentResponse, err error) {
	cli := pb.NewStorageClient(c.conn)
	request := &pb.DownloadRequest{
		FileName:  key,
		ProfileId: profileid,
		AppId:     appId,
	}
	return cli.GetFileContent(ctx, request)
}

func (c *publicClient) GetFileACL(ctx context.Context, key string, profileid string, appId string) (res *pb.ACLResponse, err error) {
	cli := pb.NewStorageClient(c.conn)
	request := &pb.GetACLRequest{
		FileName:  key,
		ProfileId: profileid,
		AppId:     appId,
	}
	return cli.GetFileACL(ctx, request)
}

func (c *publicClient) SetFileACL(ctx context.Context, key string, profileid string, aclType int32, appId string) (*pb.NothingResponse, error) {
	cli := pb.NewStorageClient(c.conn)
	request := &pb.SetACLRequest{
		FileName:  key,
		ProfileId: profileid,
		AppId:     appId,
		AclType:   aclType,
	}
	return cli.SetFileACL(ctx, request)
}

func NewStorageClient(l *zap.Logger, target string, secure bool) (client *publicClient, err error) {
	if c, e := newClient(target, secure); e != nil {
		err = e
	} else {
		client = &publicClient{client: c, l: l}
	}

	return
}

type client struct {
	conn *grpc.ClientConn
}

func newClient(target string, secure bool) (cli *client, err error) {
	if conn, e := utils.Dial(
		target,
		utils.TransportSecurity(secure),
		grpc.WithBackoffMaxDelay(5*time.Second),
	); e != nil {
		err = e
	} else {
		cli = &client{conn: conn}
	}

	return
}

func (c *client) Close() error {
	return c.conn.Close()
}
