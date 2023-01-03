package service

import (
	"context"
	"fmt"
	pb "github.com/BingguWang/grpc-go-study/server/proto"
	"github.com/BingguWang/grpc-go-study/server/utils"
	"io"
	"log"
	"sync"
	"time"
)

func CallStreamAddScoreByUser(ctx context.Context, client pb.ScoreServiceClient) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	// 客户端获取流
	stream, err := client.AddStreamScoreByUserID(ctx)
	if err != nil {
		log.Fatalf("client.AddStreamScoreByUserID failed: %v", err)
	}
	// 向流内发送请求
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := stream.Send(&pb.AddScoreByUserIDReq{
				UserID: 1,
				Scores: []*pb.Score{
					{
						Type:  uint32(i),
						Value: int32(i),
					},
				},
			}); err != nil {
				log.Fatalf(err.Error())
			}
			fmt.Println("i:", i)
		}(i)
	}
	wg.Wait()
	// 从结果可以看到grpc可以保证有序性，在服务端的接收和这里的发送顺序可以保证一致！
	// 获取响应并关闭流
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("client.AddStreamScoreByUserID failed: %v", err)
	}
	log.Println(utils.ToJsonString(reply))
}

func CallStreamScoreListByUser(ctx context.Context, client pb.ScoreServiceClient, req *pb.GetScoreListByUserIDReq) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	stream, err := client.GetStreamScoreListByUser(ctx, req)
	if err != nil {
		log.Fatalf("client.GetStreamScoreListByUser failed: %v", err)
	}
	for {
		feature, err := stream.Recv() // 从服务端返回的流接收响应
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("client.GetStreamScoreListByUser failed: %v", err)
		}
		log.Println("result :" + utils.ToJsonString(feature))
	}

}

func CallStreamBidirectional(ctx context.Context, client pb.ScoreServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.AddAndGetScore(ctx)
	if err != nil {
		log.Fatalf("client.RouteChat failed: %v", err)
	}
	waitc := make(chan struct{})
	go func() { // 开启协程接收响应流内的响应
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc) // 响应流内的响应接收完毕，CallStreamBidirectional视为调用完毕
				return
			}
			if err != nil {
				log.Fatalf("client.RouteChat failed: %v", err)
			}
			log.Println(utils.ToJsonString(in))
			fmt.Println("client recv time : ", time.Now())
		}
	}()
	for i := 0; i < 5; i++ {
		if err := stream.Send(&pb.AddScoreByUserIDReq{
			UserID: uint64(i), Scores: []*pb.Score{
				{
					Type:  uint32(i),
					Value: int32(i),
				},
			},
		}); err != nil {
			log.Fatalf("stream.Send() failed: %v", err)
		}
		fmt.Println("client send time : ", time.Now())
	}
	stream.CloseSend()
	<-waitc // 等待响应流响应接收完毕
}
