package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"
	//protocol buffersライブラリー
	proto "github.com/kcwebapply/grpc-stream-sample/proto"
	"google.golang.org/grpc"
)

/*func main() {
	fmt.Println("Hello. I'm a client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewGreetServiceClient(conn)
	fmt.Printf("Created client %f", c)
}*/

func main() {
	fmt.Println("Hello. I'm a client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewGreetServiceClient(conn)

	doServerStreaming(c)
	doClientStreaming(c)

}

func doServerStreaming(c proto.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC ...")
	req := &proto.GreetManyTimesRequest{
		Greeting: &proto.Greeting{
			FirstName: "tanaka",
			LastName:  "tarou",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTImes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTImes: %v", msg.GetResult())
	}

}

func doClientStreaming(c proto.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC")

	requests := []*proto.LongGreetRequest{
		&proto.LongGreetRequest{
			Greeting: &proto.Greeting{
				FirstName: "one",
			},
		},
		&proto.LongGreetRequest{
			Greeting: &proto.Greeting{
				FirstName: "two",
			},
		},
		&proto.LongGreetRequest{
			Greeting: &proto.Greeting{
				FirstName: "three",
			},
		},
		&proto.LongGreetRequest{
			Greeting: &proto.Greeting{
				FirstName: "four",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet %v", err)
	}
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v", err)
	}
	fmt.Printf("LongGreet Response: %v", res)

}
