package main

import (
	"io"
	"log"
	"net"

	fmpb "github.com/matheustp/find-maximum-grpc/pb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) FindMaximum(stream fmpb.FindMaximumService_FindMaximumServer) error {
	max := int32(0)
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if res.GetNum() > max {
			max = res.GetNum()
			sendErr := stream.Send(&fmpb.FindMaximumResponse{
				Max: max,
			})
			if sendErr != nil {
				log.Fatalf("Error sending data to client: %v", sendErr)
				return sendErr
			}
		}
		log.Printf("Received: %v. Maximum is: %v\n", res.GetNum(), max)
	}
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Panicf("Fail to listen: %v", err)
	}
	s := grpc.NewServer()
	fmpb.RegisterFindMaximumServiceServer(s, &server{})
	if err := s.Serve(l); err != nil {
		log.Panicf("Failed to serve: %v", err)
	}
}
