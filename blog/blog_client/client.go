package main

import (
	"GO/GRPC/blog/blogpb"
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect client %v", err)
		return
	}

	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)
	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Alka",
		Title:    "Go Love",
		Content:  "Alka loves Go",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Blog has been created: %v", createBlogRes)
	blogID := createBlogRes.GetBlog().GetId()

	//Reading the blog
	fmt.Println("Reading the blog")

	//invalid Request, expected error
	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "5da414b3daa19db2cd39221f"})
	if err2 != nil {
		fmt.Printf("Error happened while reading: %v \n", err2)
	}

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBlogErr != nil {
		fmt.Printf("Error happened while reading: %v \n", readBlogErr)
	}

	fmt.Printf("Blog was read: %v \n", readBlogRes)

	//Update the blog
	fmt.Println("Updating the blog")
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Alka Sinha",
		Title:    "I love my India",
		Content:  "Fir v Dil h Hindustani(edited)",
	}

	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if updateErr != nil {
		fmt.Printf("Error while updating the id %v\n", updateErr)
		return
	}
	fmt.Printf("Blog was updated %v\n", updateRes)

	//Deleting the blog
	delres, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})
	if err != nil {
		fmt.Printf("error while deleting the data %v\n", err)
	}
	fmt.Printf("Blogid %v was Deleted %v\n", blogID, delres)

	//List all the blog Items
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling listBlog RPC %v\n", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("something wrong happend %v\n", err)
		}
		fmt.Println(res.GetBlog())
	}
}
