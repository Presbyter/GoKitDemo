package main

import (
	"GoKitDemo/cmd/user/endpoints"
	userservice "GoKitDemo/cmd/user/service"
	"GoKitDemo/cmd/user/transport"
	"GoKitDemo/pb"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/go-kit/kit/log"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

func main() {
	fs := flag.NewFlagSet("usersvc", flag.ExitOnError)
	var (
		grpcAddr = fs.String("grpc-addr", ":8082", "gRPC listen address")
	)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var (
		errors chan error

		service     = userservice.New(logger)
		endpointSet = endpoints.New(service, logger)
		grpcServer  = transport.NewGRPCServer(endpointSet, logger)
	)

	grpcListener, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		logger.Log("transport", "gRPC", "during", "Listen", "err", err)
		panic(err)
	}

	baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
	pb.RegisterUserServiceServer(baseServer, grpcServer)
	fmt.Println("开启 gRPC. 监听:8082")
	errors <- baseServer.Serve(grpcListener)
	fmt.Println(<-errors)
}
