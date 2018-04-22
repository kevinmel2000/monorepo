package main

import (
	"context"
	"log"
	"net/http"
	"time"

	pb "github.com/lab46/monorepo/experiment/envoy/grpc/svc2/pb"
	"google.golang.org/grpc"
)

func main() {
	// lbBuilder := grpc.NewLBBuilderWithFallbackTimeout(time.Second*1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "svc2-1:9001", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	// lbBuilder.Build(*conn, )
	client := pb.NewSvc2Client(conn)

	http.HandleFunc("/pingserver", func(w http.ResponseWriter, r *http.Request) {
		req := pb.PingRequest{Message: "hola"}
		resp, err := client.Ping(context.Background(), &req)
		if err != nil {
			log.Printf("[ERROR] %s", err.Error())
			w.Write([]byte(err.Error()))
			return
		}
		if resp == nil {
			log.Println("Weird, response is nil")
			w.Write([]byte("Weird response is nil"))
			return
		}
		log.Printf("Response from: %s", resp.Message)
		w.Write([]byte(resp.Message))
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	log.Println("service is running")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
