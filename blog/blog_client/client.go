package main

import (
	//"google.golang.org/grpc/codes"
	"fmt"
	"log"
	"context"
	//"io"
	//"time"

	"google.golang.org/grpc"
	"github.com/roguesecurity/grpc-go-course/blog/blogpb"
	//"google.golang.org/grpc/status"
)

func main()  {

	fmt.Println("*** Inside Client's main() ***")

	//create a connection to server
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	//called at the end
	defer cc.Close()
	c := blogpb.NewBlogServiceClient(cc)

	// Create blog
	fmt.Println("Creating a blog")
	blog := &blogpb.Blog{
		AuthorId: "Piyush",
		Title: "My first blog",
		Content: "content of 1st blog",
	}

	createBlogRequest := &blogpb.CreateBlogRequest{
		Blog: blog,
	}
	createBlogResponse, err := c.CreateBlog(context.Background(), createBlogRequest)

	if err != nil{
		log.Fatalf("Unexpected error: %v", err)
	}

	fmt.Printf("Blog has been created: %v", createBlogResponse)
	blogID := createBlogResponse.GetBlog().GetId()

	// Read blog
	fmt.Println("Reading the blog")

	readBlogRequest := &blogpb.ReadBlogRequest{
		BlogId: blogID,
	} 
	readBlogResponse, err := c.ReadBlog(context.Background(), readBlogRequest)

	if err != nil{
		log.Fatalf("Error while reading blog: %v", err)
	}

	fmt.Printf("Blog has been read: %v", readBlogResponse)
}