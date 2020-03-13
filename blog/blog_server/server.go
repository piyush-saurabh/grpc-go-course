package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/bson"

	
	"google.golang.org/grpc/codes"
	//"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/roguesecurity/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
)

var collection *mongo.Collection

type server struct{}

// this is a model
// this will hold the data in our mongodb
// map the struct with the database (bson)
type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

// Function to create blog
func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error){
	fmt.Println("Create blog request")

	blog := req.GetBlog()

	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Title: blog.GetTitle(),
		Content: blog.GetContent(),
	}

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil{
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot convert to OID"),
		)
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id: oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title: blog.GetTitle(),
			Content: blog.GetContent(),
		},
	}, nil
}

// Function to read blog
func (*server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error){
	fmt.Println("Read blog request")

	blogID := req.GetBlogId()

	oid, err := primitive.ObjectIDFromHex(blogID)

	if err !=nil{
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse blog ID"),
		)
	}

	// Create an empty struct
	data := &blogItem{}

	filter := bson.M{"_id" : oid}
	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil{
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with specified ID: %v", err),
		)
	}

	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			Id: data.ID.Hex(),
			AuthorId: data.AuthorID,
			Content: data.Content,
			Title: data.Title,
		},
	}, nil
}

func main(){
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("*** Inside main() ***")

	// Mongodb Connection
	fmt.Println("Connecting to mongodb ...")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil{
		log.Fatal(err)
	}

	err = client.Connect(context.TODO())
	if err != nil{
		log.Fatal(err)
	}

	// run a database
	// collection is equivalent to table. 
	collection = client.Database("mydb").Collection("blog")


	//port binding
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	//Create gRPC server
	s:= grpc.NewServer()
	blogpb.RegisterBlogServiceServer(s, &server{})

	// Start the server
	go func(){
		fmt.Println("*** Starting the Blog server ***")
		if err:=s.Serve(lis); err!=nil {
			log.Fatalf("Failed to start the server: %v", err)
		}
	}()

	// Waiting for CTRL+C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()

	fmt.Println("Closing MongoDB Connection")
	client.Disconnect(context.TODO())

	fmt.Println("End of program")

	
}