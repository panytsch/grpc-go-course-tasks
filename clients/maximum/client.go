package main

import (
	"context"
	"github.com/panytsch/grpc-go-course-tasks/go_course_pb"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	service := &MaximumService{go_course_pb.NewGoCourseServiceClient(cc)}
	service.ShowMaximum([]int64{23234, 3543, 45647, 94395, 4, 7787})
}

type MaximumService struct {
	go_course_pb.GoCourseServiceClient
}

func (s *MaximumService) ShowMaximum(numbers []int64) {
	if len(numbers) == 0 {
		return
	}
	stream, err := s.GoCourseServiceClient.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Cannot open connectionError: %v", err)
	}
	waitCh := make(chan bool)
	go sendMaximumRequests(numbers, stream)

	go receiveMaximumResponses(stream, waitCh)

	<-waitCh
}

func sendMaximumRequests(numbers []int64, stream go_course_pb.GoCourseService_FindMaximumClient) {
	for _, number := range numbers {
		request := &go_course_pb.MaximumRequest{
			Number: number,
		}
		err := stream.Send(request)
		if err != nil {
			log.Fatalf("Cannot send request. Error: %v, Request: %v", err, request)
		}
	}
	_ = stream.CloseSend()
}

func receiveMaximumResponses(stream go_course_pb.GoCourseService_FindMaximumClient, waitCh chan bool) {
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Error while receiving responses: %v", err)
		}
		log.Printf("Received maximum: %v", res.GetCurrentMaximum())
	}
	close(waitCh)
}
