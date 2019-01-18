package main

import (
	"context"
	"io"
	"log"

	fmpb "github.com/matheustp/find-maximum-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Panicf("Error dialing: %v", err)
	}
	c := fmpb.NewFindMaximumServiceClient(cc)
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Panicf("Error calling function: %v", err)
	}
	wchan := make(chan struct{})
	go func() {
		numbers := []int32{1, 5, 3, 6, 2, 20}
		for _, n := range numbers {
			log.Printf("Sending %v\n", n)
			stream.Send(&fmpb.FindMaximumRequest{
				Num: n,
			})
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving message: %v", err)
				break
			}
			log.Printf("New maximum value: %v\n", res.GetMax())
		}
		close(wchan)
	}()
	<-wchan
}
