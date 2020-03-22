# Creating gRPC server and client using golang

Ref: https://www.udemy.com/course/grpc-golang

### Create a protobuf file

```proto
syntax = "proto3";

// This package name will be the part of URL while calling gRCP server
package greet;

// Folder name in which this file resides
option go_package="greetpb";

message Greeting{
    string first_name = 1;
    string last_name = 2;
}

message GreetRequest {
    Greeting greeting = 1;
}

message GreetResponse {
    string result = 1;
}

// This service will bind to ONE port on the server. Analogous to 1 microservice
// This service exposes 3 methods. Analogous to 3 verbs in REST
service GreetService{

    //unary
    rpc Greet(GreetRequest) returns (GreetResponse){};
}
```

### Create the gRPC Server

```go
type server struct{}

func main(){
	//port binding
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	//Create gRPC server
	s:= grpc.NewServer()
  
  // This method is located in the compiled protobuf file
	greetpb.RegisterGreetServiceServer(s, &server{})
	if err:=s.Serve(lis); err!=nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
```

### RPC defination in gRPC server
```go
func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error){
	fmt.Printf("Greet function was invoked with %v", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res,nil
}
```

### Create gRPC Client
```go
func main() {
	//create a connection to server
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	//called at the end
	defer cc.Close()

  // Get the client handler to call RPC methods
	c := greetpb.NewGreetServiceClient(cc)

  // Wrapper method to call RPC methods
	doUnary(c)
}
```

### Calling RPC methods from Client
```go
func doUnary(c greetpb.GreetServiceClient) {

  // Create the response body based on response specification defined in proto file
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Piyush",
			LastName:  "Saurabh",
		},
	}

	// Calling gRPC server method and pass the current context
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet:%v", res.Result)
}
```


