package main

import (
	"google.golang.org/grpc/codes"
	"context"
	"fmt"
	"net"
	"log"
	"time"
	"io"
	"math"

	"google.golang.org/grpc"
	"github.com/roguesecurity/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc/status"
)

type server struct{}

// Unary
func (*server) Sum(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error){

	fmt.Printf("Sum function was invoked with %v\n", req)

	num1 := req.GetNum1()
	num2 := req.GetNum2()

	result := num1 + num2

	res := &calculatorpb.CalculatorResponse{
		Result: result,
	}

	return res, nil
}

// Client Streaming
func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error{

	fmt.Printf("PrimeNumberDecomposition function was invoked with %v\n", req)

	num := req.GetNumber()
	k := 2

	for num > 1{
		if num % int32(k) == 0{
			//result := "Prime number decomposition: " + strconv.Itoa(k)
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				Result: int32(k),
			}
			stream.Send(res)
			num = num / int32(k)
		} else {
			k = k + 1
		}
		
		time.Sleep(1000 * time.Millisecond)

	}
	return nil
}

//Server Streaming
func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage function was invoked with stream input\n")

	var avg,sum,count float32

	count = 1
	sum = 0

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: float32(avg),
		
			})
		}

		if err != nil{
			log.Fatalf("Error while reading client stream: %v", err)
		}

		num := req.GetNumber()
		sum += float32(num)
		avg = sum/count
		count++
	}

	
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error){

	fmt.Printf("SquareRoot function was invoked with %v\n", req)

	number := req.GetNumber()

	if (number < 0){
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", number),
		)
	}

	res := &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}

	return res, nil
}

func main(){
	fmt.Println("*** Starting the server. Inside main() ***")

	//port binding
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	//Create gRPC server
	s:= grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	if err:=s.Serve(lis); err!=nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}