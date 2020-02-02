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

	service := &AverageService{go_course_pb.NewGoCourseServiceClient(cc)}
	result := service.AverageNumber([]int64{23234, 3543, 45647, 94395, 4, 7787})
	log.Printf("result : %v", result)
}

type AverageService struct {
	go_course_pb.GoCourseServiceClient
}

func (s *AverageService) AverageNumber(numbers []int64) float64 {
	var result float64
	if len(numbers) == 0 {
		return result
	}
	stream, err := s.GoCourseServiceClient.LongAverage(context.Background())
	if err != nil {
		log.Fatalf("Cannot open connectionError: %v", err)
	}
	for _, number := range numbers {
		request := &go_course_pb.AverageRequest{
			Number: number,
		}
		err = stream.Send(request)
		if err != nil {
			log.Fatalf("Cannot send request. Error: %v, Request: %v", err, request)
		}
	}
	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Cannot recieve response from service: %v", err)
	}
	return response.GetAverage()
}
