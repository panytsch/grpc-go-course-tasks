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

	service := &PrimeNumberService{go_course_pb.NewGoCourseServiceClient(cc)}
	result := service.PrimeNumber(120)
	log.Printf("result for number 120: %v", result)
}

type PrimeNumberService struct {
	go_course_pb.GoCourseServiceClient
}

func (s *PrimeNumberService) PrimeNumber(number int32) []int32 {
	var result []int32
	req := &go_course_pb.PrimeNumberRequest{
		Number: number,
	}
	stream, err := s.GoCourseServiceClient.PrimeNumber(context.Background(), req)
	if err != nil {
		log.Fatalf("Cannot receive response from service: %v", err)
	}
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Error while fetching responses: %v", err)
		}
		result = append(result, response.GetDivisor())
	}
	return result
}
