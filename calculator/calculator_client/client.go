package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/roguesecurity/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func main() {

	fmt.Println("*** Inside Client's main() ***")

	//create a connection to server
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	//called at the end
	defer cc.Close()
	c := calculatorpb.NewCalculatorServiceClient(cc)

	doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)

	//doErrorUnary(c)

}

func doUnary(c calculatorpb.CalculatorServiceClient) {

	fmt.Println("Starting to do a Unary RPC...")

	req := &calculatorpb.CalculatorRequest{
		Num1: 3,
		Num2: 10,
	}

	res, err := c.Sum(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}
	log.Printf("Response from Sum:%v", res.Result)

}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 120,
	}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()

		if err == io.EOF {
			// we have reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		log.Printf("Response from GreetManyTimes:%v", msg.GetResult())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {

	stream, err := c.ComputeAverage(context.Background())

	requests := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Number: 1,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 2,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 3,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 4,
		},
	}

	if err != nil {
		log.Fatalf("error while reading stream: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending request: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving the response: %v\n", res)
	}
	fmt.Printf("Average: %v\n", res.GetAverage())

}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {

	fmt.Println("Starting to do a Error Code Unary RPC...")

	req := &calculatorpb.SquareRootRequest{
		Number: -25,
	}

	res, err := c.SquareRoot(context.Background(), req)

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC (user error)
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())

			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent negative number")
			}
		} else {
			log.Fatalf("Big error while calling square root RPC: %v", err)
		}

	}

	log.Printf("Response from Sum:%v", res.NumberRoot)
}
