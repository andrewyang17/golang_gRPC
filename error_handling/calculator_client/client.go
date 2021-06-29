package main

import (
	"context"
	"fmt"
	"goGrpc/error_handling/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("could not connect: %w", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	doErrorUnary(c)
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a SquareRoot Unary RPC...")
	// Correct call
	doErrorCall(c, 10)
	// Error call
	doErrorCall(c, -2)
}

func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})
	if err != nil {
		resErr, ok := status.FromError(err)
		if ok {
			fmt.Printf("Error message from server: %v\n", resErr.Message())
			fmt.Printf("Error code from server: %v\n", resErr.Code())
			if resErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
				return
			}
			return
		} else {
			log.Fatalf("error calling SquareRoot: %v", err)
		}
	}
	fmt.Printf("Result of SquareRoot of %v:  %v \n", n, res.GetNumber())
}