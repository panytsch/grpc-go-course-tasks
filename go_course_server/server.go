package main

import (
	"context"
	"github.com/panytsch/grpc-go-course-tasks/go_course_pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	go_course_pb.RegisterGoCourseServiceServer(s, &server{})

	log.Printf("listening.. \n")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Cannot serve: %v", err)
	}
}

type server struct{}

func (*server) LongAverage(stream go_course_pb.GoCourseService_LongAverageServer) error {
	log.Printf("LongAverage method has invoked\n")
	sum := int64(0)
	count := int64(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			result := float64(sum) / float64(count)
			return stream.SendAndClose(&go_course_pb.AverageResponse{
				Average: result,
			})
		} else if err != nil {
			log.Fatalf("Error while receiving requests: %v", err)
		}
		log.Printf("Received request: %v\n", req)
		sum += req.GetNumber()
		count++
	}
}

func (*server) Sum(_ context.Context, request *go_course_pb.SumRequest) (*go_course_pb.SumResponse, error) {
	log.Printf("Got new sum request: %v\n", request)
	sum := int64(request.GetFirstNumber() + request.GetSecondNumber())
	return &go_course_pb.SumResponse{
		Sum: sum,
	}, nil
}

func (*server) PrimeNumber(req *go_course_pb.PrimeNumberRequest, srv go_course_pb.GoCourseService_PrimeNumberServer) error {
	log.Printf("Received new request: %v", req)
	divisor := int32(2)
	number := req.GetNumber()
	for number > 1 {
		if number%divisor == 0 {
			err := srv.Send(&go_course_pb.PrimeNumberResponse{
				Divisor: divisor,
			})
			if err != nil {
				log.Fatalf("Cannot send response: %v", err)
			}
			log.Printf("Sending number %v", divisor)
			number = number / divisor
		} else {
			divisor++
		}
	}
	return nil
}
