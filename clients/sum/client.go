package main

import (
	"context"
	"github.com/panytsch/grpc-go-course-tasks/go_course_pb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	s := &SumService{go_course_pb.NewGoCourseServiceClient(cc)}
	sum := s.Sum(45, 22)
	log.Printf("sum: %v", sum)
}

type SumService struct {
	go_course_pb.GoCourseServiceClient
}

func (s *SumService) Sum(fn int32, sn int32) int64 {
	req := &go_course_pb.SumRequest{
		FirstNumber:  fn,
		SecondNumber: sn,
	}
	log.Printf("Trying to send request to sum grpc: %v", req)
	res, err := s.GoCourseServiceClient.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Cannot get answer from sum grpc: %v", err)
	}
	log.Printf("Got answer: %v", *res)
	return res.GetSum()
}
