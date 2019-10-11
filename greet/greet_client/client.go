package main

import (
	"GO/GRPC/greet/greetpb"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello I'm a client")
	clientConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect %v ", err)
	}
	defer clientConn.Close()
	c := greetpb.NewGreetServiceClient(clientConn)
	//fmt.Printf("created client %f ", c)
	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBiDiScreaming(c)
	doUnaryWithDeadline(c, 5*time.Second)
	doUnaryWithDeadline(c, 1*time.Second)

}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("starting to do Unary RPC......")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "stephan",
			LastName:  "Mark",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("response from greet %v", res.Result)

}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("starting to do Server Streaming RPC......")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "stephan",
			LastName:  "Mark",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes..%v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			//reached end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream %v", err)
		}
		log.Printf("response from GreetManyTimes : %v ", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	requests := []*greetpb.LongGreetrequest{
		&greetpb.LongGreetrequest{
			Greeting: &greetpb.Greeting{
				FirstName: "stephan",
			},
		},
		&greetpb.LongGreetrequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Mark",
			},
		},
		&greetpb.LongGreetrequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Lucy",
			},
		},
		&greetpb.LongGreetrequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Jonathan",
			},
		},
		&greetpb.LongGreetrequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Donald",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("error while calling LongGreet %v", err)
	}
	for _, req := range requests {
		fmt.Printf("sending the request %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("err while reciving response %v \n", err)
	}
	fmt.Printf("response %v\n", response)
}

func doBiDiScreaming(c greetpb.GreetServiceClient) {
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while creating stream %v", err)
		return
	}
	waitc := make(chan struct{})

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "stephan",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Mark",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Lucy",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Jonathan",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Donald",
			},
		},
	}

	//for sending the messages
	go func() {
		for _, req := range requests {
			fmt.Printf("sending message %v ", req)
			err := stream.Send(req)
			if err != nil {
				log.Fatalf("error while sending the request %v", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	//for receiving the messages
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while receiving the response %v", err)
				break
			}
			log.Printf("received %v", res.GetResult())
		}
		close(waitc)
	}()
	<-waitc
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("starting to doUnaryWithDeadlinery RPC......")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "stephan",
			LastName:  "Maarek",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Time Out was hit, deadline was exceeded")
			} else {
				fmt.Printf("unexpected error %v ", statusErr)
			}
		} else {
			log.Fatalf("error while calling doUnaryWithDeadline RPC: %v", err)
		}
		return
	}
	log.Printf("response from doUnaryWithDeadline %v", res.Result)

}
