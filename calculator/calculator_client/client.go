package main

import (
	calculatorpb "GO/GRPC/calculator/calculatorpb"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello I am a client")
	clientConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect : %v", err)
	}
	defer clientConn.Close()

	c := calculatorpb.NewCalculatorServiceClient(clientConn)

	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBiDiStreaming(c)
	doErrorUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Unary RPC Operation started......")
	req := &calculatorpb.SumRequest{
		FirstNum:  10,
		SecondNum: 3,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling sum RPC....%v", err)
	}

	log.Printf("response from Sum Server %v ", res.GetSumResult())

}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("PrimeDecomposition Server Streaming RPC Operation started......")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 120,
	}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Printf("error while callingPrimeNumberDecomposition RPC %v ", err)
	}
	for {
		res, err := resStream.Recv()
		if err == io.EOF {
			//reached end of stream
			break
		}
		if err != nil {
			log.Printf("received error in streaming the msg %v", err)
		}
		log.Printf("Response from PrimeFactorDecomposition %v", res.GetPrimeFactor())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("ComputeAverage client Streaming RPC Operation started......")
	numbers := []int64{3, 5, 9, 54, 23}
	Stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling computeAverage client function %v ", err)
	}
	for _, number := range numbers {
		fmt.Printf("sending number %v\n", number)
		Stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}
	res, err := Stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving the data %v", err)
	}
	fmt.Printf("The average of the number is %v", res.GetAverage())
}

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Bi Direction RPC operation started.....")
	//call the function
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("error while calling the FindMaximum function %v", err)
	}

	waitc := make(chan struct{})
	numbers := []int64{4, 7, 2, 19, 4, 6, 32}
	//send go routing
	go func() {
		for _, number := range numbers {
			fmt.Printf("sending number %v\n", number)
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	//receive goroutine
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while reading client %v", err)
				break
			}
			log.Printf("The maxmim number from the responses %v\n", res.GetMaximum())
		}
		close(waitc)
	}()
	<-waitc
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("error handling in client")
	doErrorCall(c, 10)
	doErrorCall(c, -2)
}

func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})

	if err != nil {
		resErr, ok := status.FromError(err)
		if ok {
			//actual error from user
			fmt.Printf("error msg from server %v\n", resErr.Message())
			fmt.Println(resErr.Code())
			if resErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
				return
			}
		} else {
			log.Fatalf("big Error calling square root %v\n", resErr)
			return
		}
	}
	fmt.Printf("The square root of number %v = %v\n", n, res.GetNumberRoot())

}
