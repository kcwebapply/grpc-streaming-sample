package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	proto "github.com/kcwebapply/grpc-stream-sample/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//server greet.pb.go GreetServiceClient
type server struct{}

func (*server) Greet(ctx context.Context, req *proto.GreetRequest) (*proto.GreetResponse, error) {
	fmt.Printf("Greet func was invoked with %v", req)
	firstName := req.GetGreeting().GetFirstName()
	if firstName == "" {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Recived a empty string",
		)
	}

	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "the client canceld the request")
	}

	result := "Hello " + firstName
	//client config deadline
	/*
	   res := &proto.GreetWithDeadlineResponse{
	       Result: result,
	   }
	*/
	res := &proto.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *proto.GreetManyTimesRequest, stream proto.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes func was invoked with %v", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello" + firstName + " number " + strconv.Itoa(i)
		res := &proto.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream proto.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet function was invoked with a streaming request")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//we have finished reading the client stream
			return stream.SendAndClose(&proto.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "

	}
}

func main() {
	fmt.Println("Hello World!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
