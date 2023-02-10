package test

// test grpc server api -- ../pkg/server
import (
	"context"
	"encoding/base64"
	grpcClient "fileMS/pkg/server"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"testing"
)

var Client grpcClient.FileMsServiceClient

func init() {
	conn, err := grpc.Dial("127.0.0.1:38083", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	//defer conn.Close()
	Client = grpcClient.NewFileMsServiceClient(conn)
}

func TestGetFile(t *testing.T) {
	in := &grpcClient.FileRequest{
		Uuid: "b859e97f-8f98-4fe7-9f0a-dd3c56c1aed8",
	}
	ret, err := Client.GetFile(context.Background(), in)
	fmt.Println(ret, err)
}

func TestGetFileContent(t *testing.T) {
	in := &grpcClient.FileRequest{
		Uuid: "0a25b634-c039-4355-8285-de8ea1620c2c",
	}
	ret, err := Client.GetFileContent(context.Background(), in)
	if nil == err {
		bytes, _ := base64.StdEncoding.DecodeString(ret.Content)
		fmt.Println(string(bytes))
	}

}

func TestDelFile(t *testing.T) {
	in := &grpcClient.FileRequest{
		Uuid: "b859e97f-8f98-4fe7-9f0a-dd3c56c1aed8",
	}
	ret, err := Client.DelFile(context.Background(), in)
	fmt.Println(ret, err)
}

func TestListFileVersions(t *testing.T) {
	in := &grpcClient.FileRequest{
		Uuid: "b10a8c3a-5659-451c-af46-1e4afebac670",
	}
	ret, err := Client.ListFileVersions(context.Background(), in)
	fmt.Println(ret, err)
}
