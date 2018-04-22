package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	pb "github.com/lab46/monorepo/experiment/grpc/svc2/pb"
	"google.golang.org/grpc"
)

const svcEnv = "SVCENV"

var senv string

type service struct{}

func (s *service) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	resp := pb.PingResponse{
		Message: fmt.Sprintf("hello from %s", senv),
	}
	log.Printf("response from %s", senv)
	return &resp, nil
}

func main() {
	senv = os.Getenv(svcEnv)
	listener, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	svc := &service{}
	pb.RegisterSvc2Server(s, svc)
	log.Println("gRPC service started")
	go func() {
		log.Fatal(s.Serve(listener))
	}()

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		resp, err := svc.Ping(context.Background(), &pb.PingRequest{Message: "holla"})
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte(resp.Message))
	})
	go func() {
		http.ListenAndServe(":9000", nil)
	}()
	select {}
}
