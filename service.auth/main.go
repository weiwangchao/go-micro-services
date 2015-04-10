package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	pb "github.com/harlow/go-micro-services/service.auth/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	port       = flag.Int("port", 10002, "The server port")
	jsonDBFile = flag.String("json_db_file", "./data/customers.json", "A json file containing a list of customers")
	serverName = "service.auth"
)

type authServer struct {
	customers []*pb.Customer
}

// GetCustomer finds a customer by authentication token.
func (s *authServer) GetCustomer(ctx context.Context, req *pb.Req) (*pb.Customer, error) {
	for _, c := range s.customers {
		if c.AuthToken == req.AuthToken {
			return c, nil
		}
	}
	return &pb.Customer{}, errors.New("Invalid Token")
}

// loadCustomers loads customers from a JSON file.
func (s *authServer) loadCustomers(filePath string) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to load file: %v", err)
	}
	if err := json.Unmarshal(file, &s.customers); err != nil {
		log.Fatalf("Failed to load json: %v", err)
	}
}

func newServer() *authServer {
	s := new(authServer)
	s.loadCustomers(*jsonDBFile)
	return s
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}