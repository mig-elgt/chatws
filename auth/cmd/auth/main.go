package main

import (
	"flag"
	"net"

	"github.com/mig-elgt/chatws/auth/grpc"
	"github.com/mig-elgt/jwt"
	"github.com/sirupsen/logrus"
)

func main() {
	grpcAddr := flag.String("listen", ":8080", "address:port to listen on")
	flag.Parse()

	lis, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		logrus.Fatalf("could not listen to port %v: %v", *grpcAddr, err)
	}
	logrus.Infof("GRPC listening on %s", *grpcAddr)

	tv := jwt.New("secret_key")
	s := grpc.NewAPI(tv)

	if err := s.Serve(lis); err != nil {
		logrus.Fatalf("grpc Server failed: %v", err)
	}
}
