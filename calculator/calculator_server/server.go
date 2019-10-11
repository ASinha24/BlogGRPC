package main

import (
	calculatorpb "GO/GRPC/calculator/calculatorpb"
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum Function was invoked with %v\n ", req)
	result := req.GetFirstNum() + req.GetSecondNum()
	res := &calculatorpb.SumResponse{
		SumResult: result,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition Function was invoked with %v\n ", req)
	number := req.GetNumber()
	divisor := int64(2)

	for number > 1 {
		if number%divisor == 0 {
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			}
			stream.Send(res)
			// stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
			// 	PrimeFactor: divisor,
			// })
			number = number / divisor
			time.Sleep(1000 * time.Millisecond)
		} else {
			divisor++
			fmt.Printf("divisor has increased to %v\n", divisor)
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage Function was invoked with client streaming\n")
	sum := int64(0)
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//resched end
			average := float64(sum) / float64(count)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})
		}
		if err != nil {
			log.Fatalf("error while streaming the client %v ", err)
		}
		sum += req.GetNumber()
		count++
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Println("FindMaximum Function invoked!")
	max := int64(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while reading client stream %v", err)
			return err
		}

		number := req.GetNumber()
		if number > max {
			max = number
			err = stream.Send(&calculatorpb.FindMaximumResponse{
				Maximum: max,
			})
			if err != nil {
				log.Fatalf("error while sending data to client %v", err)
				return err
			}
		}
	}

}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	number := req.GetNumber()

	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("The number is negative %v...", number),
		)
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	fmt.Println("Hello GRPC calculator server on port 50051")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("failed to listen  %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	reflection.Register(s)

	//binding the port to grpc server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve %v ", err)
	}

}
