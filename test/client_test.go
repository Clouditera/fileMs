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
		Uuid: "da4879ff-90f8-4000-8ab6-62ecd59555df",
	}
	ret, err := Client.GetFile(context.Background(), in)
	fmt.Println(ret, err)
}

func TestGetFileContent(t *testing.T) {
	in := &grpcClient.FileRequest{
		Uuid: "a94f93f4-0c22-45ba-a0e2-a81514727d27",
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
		Uuid: "196b3b43-d100-4b8e-987b-3020f6db4665",
	}
	ret, err := Client.ListFileVersions(context.Background(), in)
	fmt.Println(ret, err)
}
