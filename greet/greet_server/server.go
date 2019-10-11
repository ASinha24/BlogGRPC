package main

import (
	greetpb "GO/GRPC/greet/greetpb"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type server struct{}

//Unary API
func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet Function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	result := "Hello " + firstName + " " + lastName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

//server streaming API
func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes Function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

//Client streaming API
func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet Function was invoked with client streaming \n")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//reach end of client streaming
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("error while client streaming...%v ", err)
		}
		firstname := req.GetGreeting().GetFirstName()
		result += "Hello " + firstname + "!"
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone Function was invoked with client streaming \n")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while receiving the request from client %v", err)
			return err
		}
		firstname := req.GetGreeting().GetFirstName()
		result := "Hello " + firstname + "! "
		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if err != nil {
			log.Fatalf("Error while sending data to client %v", err)
		}
	}
}

//commands to kill the running process , run in cmd prompt
// netstat -ano | findstr :50051

// taskkill/pid 4712 /F

func main() {
	//default port for grpc network tcp
	fmt.Println("Hello GRPC server started on port 50051")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("failed to listen  %v", err)
	}

	s := grpc.NewServer()

	greetpb.RegisterGreetServiceServer(s, &server{})

	//binding the port to grpc server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve %v ", err)
	}
}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Printf("GreetWithDeadline Function was invoked with %v\n", req)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("The client cancelled the requset")
			return nil, status.Errorf(codes.Canceled, "The client cancelled the request\n")
		}
		time.Sleep(1 * time.Second)
	}
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	result := "Hello " + firstName + " " + lastName
	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}
	return res, nil
}
