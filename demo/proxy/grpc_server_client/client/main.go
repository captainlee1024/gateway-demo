package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	pb "github.com/captainlee1024/gateway-demo/demo/proxy/grpc_server_client/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	//addr = flag.String("addr", "localhost:50051", "the address to connect to")
	addr = flag.String("addr", "localhost:8012", "the address to connect to")
)

const (
	timestampFormat = time.StampNano // "jan_2 15:04:05.000"
	streamingCount  = 10
	// AccessToken xx
	//AccessToken = " some-secret-token1"
	AccessToken = " eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTA4ODc0MzAsImlzcyI6ImZyb250ZW5kX3Rlc3RfMSJ9.PIWxtmzcSXbgAG9RWpPIJZhYjub61i1A90ajHA8SapM"
)

func unaryCallWithMetadata(c pb.EchoClient, message string) {
	fmt.Printf("=== unary ===\n")

	// Create metadata and conetxt
	md := metadata.Pairs("timetamp", time.Now().Format(timestampFormat))
	md.Append("authorization", "Bearer"+AccessToken)

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	r, err := c.UnaryEcho(ctx, &pb.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("failed to call UnaryEcho: %v\n", err)
	}
	fmt.Printf("response: %v\n", r.Message)
}

func serverStreamingWithMetadata(c pb.EchoClient, message string) {
	fmt.Printf("=== server streaming ===\n")

	md := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
	md.Append("authorization", "Bearer"+AccessToken)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	stream, err := c.ServerStreamingEcho(ctx, &pb.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("failed to call ServerStreamingEcho: %v\n", err)
	}

	// Read all the response.
	var rpcStatus error
	fmt.Printf("response:\n")
	for {
		r, err := stream.Recv()
		if err != nil {
			rpcStatus = err
			break
		}
		fmt.Printf("- %s\n", r.Message)
	}
	if rpcStatus != io.EOF {
		log.Fatalf("failed finish server streaming: %v\n", rpcStatus)
	}
}

func clientStreamWithMetadta(c pb.EchoClient, message string) {
	fmt.Printf("=== client streaming ===\n")
	md := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
	md.Append("authorization", "Bearer"+AccessToken)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := c.ClientStreamingEcho(ctx)
	if err != nil {
		log.Fatalf("failed to call ClientStreamingEcho: %v\n", err)
	}

	// Send all requests to the server.
	for i := 0; i < streamingCount; i++ {
		if err := stream.Send(&pb.EchoRequest{Message: message}); err != nil {
			log.Fatalf("failed to send streaming: %v\n", err)
		}
	}
	// Read the response.
	r, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to CloseAndRecv: %v\n", err)
	}
	fmt.Printf("response: %v\n", r.Message)
}

func bidirectionalWithMetadata(c pb.EchoClient, message string) {
	fmt.Printf("=== bidirectional ===\n")
	md := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
	md.Append("authorization", "Bearer"+AccessToken)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := c.BidirectionalStreamingEcho(ctx)
	if err != nil {
		log.Fatalf("failed to call bidirectionalStreamingEcho: %v\n", err)
	}

	go func() {
		// Send all requests to the server
		for i := 1; i < streamingCount; i++ {
			err := stream.Send(&pb.EchoRequest{Message: message})
			if err != nil {
				log.Fatalf("failed to send streaming: %v\n", err)
			}
		}
		stream.CloseSend()
	}()

	// Read all the responses.
	var rpcStatus error
	fmt.Printf("response:\n")
	for {
		r, err := stream.Recv()
		if err != nil {
			rpcStatus = err
			break
		}
		fmt.Printf("- %s\n", r.Message)
	}
	if rpcStatus != nil {
		log.Fatalf("failed to finish server streaming: %v\n", rpcStatus)
	}
}

const message = "this is examplse/metadata"

func main() {
	flag.Parse()
	wg := sync.WaitGroup{}
	//for i := 0; i < 1; i++ {
	// 两个并发测试限流中间件
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := grpc.Dial(*addr, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("did not connect: %v\n", err)
			}
			defer conn.Close()

			c := pb.NewEchoClient(conn)

			// 调用一元方法
			// for i := 0; i < 10; i++ {
			unaryCallWithMetadata(c, message)
			time.Sleep(time.Millisecond * 400)
			// }

			// 服务端流式
			serverStreamingWithMetadata(c, message)
			//time.Sleep(time.Second * 1)
			// 限流测试，缩小睡眠时间
			time.Sleep(time.Millisecond * 400)

			// 客户端流式
			clientStreamWithMetadta(c, message)
			//time.Sleep(time.Millisecond * 400)
			// 限流测试，缩小睡眠时间
			time.Sleep(time.Second * 1)

			// 双向流式
			bidirectionalWithMetadata(c, message)
			// time.Sleep(time.Second * 1)
		}()
	}
	wg.Wait()
	time.Sleep(time.Second * 1)
}
