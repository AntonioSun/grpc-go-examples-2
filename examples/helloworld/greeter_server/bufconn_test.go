package main

import (
	"context"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/AntonioSun/grpc-go-examples-2/examples/helloworld/helloworld"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
}

func bufDialer(ctx context.Context, address string) (net.Conn, error) {
	return lis.Dial()
}

func TestHello_bufconn(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	// set up test cases
	helloTests := []struct {
		name     string
		expected string
	}{
		{
			name:     "world",
			expected: "Hello world",
		},
		{
			name:     "bob",
			expected: "Hello bob",
		},
	}

	for _, tt := range helloTests {
		req := &pb.HelloRequest{Name: tt.name}
		resp, err := client.SayHello(ctx, req)
		if err != nil {
			t.Errorf("HelloTest(%v) got unexpected error", err)
		}
		if resp.Message != tt.expected {
			t.Errorf("HelloText(%v)=%v, expected %v", tt.name, resp.Message, tt.expected)
		}
	}
}
