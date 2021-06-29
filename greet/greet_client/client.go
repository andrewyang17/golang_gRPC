package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"goGrpc/greet/greetpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello from client")

	tls := true
	opts := grpc.WithInsecure()

	if tls {
		certFile := "ssl/ca.crt"  // CA trust certificate
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loaoding CA trust certificate: %v", sslErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	cc, err := grpc.Dial("localhost:50051", opts)

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBiDirectionalStreaming(c)
	//doUnaryWithDeadline(c, 1*time.Second)  // should complete
	//doUnaryWithDeadline(c, 5*time.Second)  // should timeout
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Andrew",
			LastName:  "Yang",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalln("error while calling Greet RPC: %w", err)
	}
	fmt.Printf("Response from Greet: %v \n", res.GetResult())
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Andrew",
			LastName:  "Yang",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalln("error while calling GreetManyTimes RPC: %w", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("error while reading stream: %w", err)
		}
		fmt.Printf("Response from GreetManytimes: %v \n", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")
	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Andrew",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Andrew1",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Andrew2",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Andrew3",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Andrew4",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalln("error while calling LongGreet RPC: %w", err)
	}
	for _, req := range requests {
		fmt.Printf("Sending req: %v \n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln("error while receiving response from LongGreet RPC: %w", err)
	}
	fmt.Printf("LongGreet Response: %v\n", res)
}

func doBiDirectionalStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Bi-Directional Streaming RPC...")
	var wg sync.WaitGroup

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Andrew",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Andrew1",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Andrew2",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Andrew3",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Andrew4",
			},
		},
	}
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while creating stream: %v", err)
		return
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, req := range requests {
			fmt.Printf("Sending message: %v \n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while receiving: %v", err)
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
	}()
	wg.Wait()
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting to do a UnaryWithDeadline RPC...")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Andrew",
			LastName:  "Yang",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline was exceeded")
			} else {
				fmt.Printf("unexpected error: %v \n", statusErr)
			}
		} else {
			log.Fatalln("error while calling UnaryWithDeadline RPC: %w", err)
		}
		return
	}
	fmt.Printf("Response from UnaryWithDeadline: %v \n", res.GetResult())
}