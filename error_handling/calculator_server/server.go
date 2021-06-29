package main

import (
	"context"
	"fmt"
	"goGrpc/error_handling/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"math"
	"net"
)

type server struct {}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Println("Received SquareRoot RPC")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("received a negative number: %v", number))
	}
	return &calculatorpb.SquareRootResponse{Number: math.Sqrt(float64(number))}, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalln("Failed to listen: %w", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalln("Failed to serve: %w", err)
	}
	fmt.Println("gRPC server is listening!")
}
