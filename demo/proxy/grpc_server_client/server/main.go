package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	"google.golang.org/grpc/metadata"

	pb "github.com/captainlee1024/gateway-demo/demo/proxy/grpc_server_client/proto"
)

var (
	port = flag.Int("port", 50055, "the port to server on")
)

const (
	streamingCount = 10
)

type server struct {
	// gateway-grpc 不许嵌入这个结构体才能向前兼容
	pb.UnimplementedEchoServer
}

func (s *server) ServerStreamingEcho(in *pb.EchoRequest, stream pb.Echo_ServerStreamingEchoServer) error {
	fmt.Printf("=== ServerStreamingEcho===\n")
	fmt.Printf("request received:%v\n", in)
	// Read request and send response
	for i := 0; i < streamingCount; i++ {
		fmt.Printf("echo message%v\n", in.Message)
		err := stream.Send(&pb.EchoResponse{Message: in.Message})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *server) ClientStreamingEcho(stream pb.Echo_ClientStreamingEchoServer) error {
	fmt.Printf("=== ClientStreamingEcho ===\n")
	// Read request and send response
	var message string
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("echo last received message\n")
			return stream.SendAndClose(&pb.EchoResponse{Message: message})
		}
		message = in.Message
		fmt.Printf("request received: %v, building echo\n", in)
		if err != nil {
			return err
		}
	}
}

func (s *server) BidirectionalStreamingEcho(stream pb.Echo_BidirectionalStreamingEchoServer) error {
	fmt.Printf("=== BidirectionalStreamingEcho ===\n")
	// Read reuqete and send response.
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("request received %v, send echo\n", in)
		if err := stream.Send(&pb.EchoResponse{Message: in.Message}); err != nil {
			return err
		}
	}
}

func (s *server) UnaryEcho(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	fmt.Printf("=== UnaryEcho ===\n")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("miss metadata from context")
	}
	fmt.Printf("md ========> %#v\n", md)
	fmt.Printf("request received:%v, sending echo\n", in)
	return &pb.EchoResponse{Message: in.Message}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen:%v", err)
	}
	fmt.Printf("server listening at %v\n", lis.Addr())
	s := grpc.NewServer()
	pb.RegisterEchoServer(s, &server{})
	s.Serve(lis)
}
