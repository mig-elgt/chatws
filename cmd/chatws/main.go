package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/mig-elgt/chatws/authgrpc"
	"github.com/mig-elgt/chatws/handler"
	"github.com/mig-elgt/chatws/rabbitmq"
	"github.com/sirupsen/logrus"
)

var (
	rabbitmqHost = os.Getenv("RABBITMQ_HOST")
	rabbitmqPort = os.Getenv("RABBITMQ_PORT")
	rabbitmqUser = os.Getenv("RABBITMQ_USER")
	rabbitmqPass = os.Getenv("RABBITMQ_PASSWORD")
)

func main() {
	restAddr := flag.Int("addr", 8080, "web socket server port")
	authSvcHost := flag.String("auth_host", "auth:8080", "address:port to listen on")
	flag.Parse()
	broker, err := rabbitmq.New(rabbitmqUser, rabbitmqPass, rabbitmqHost, rabbitmqPort)
	if err != nil {
		panic(err)
	}
	authSvc, err := authgrpc.New(*authSvcHost)
	if err != nil {
		panic(err)
	}
	h := handler.New(authSvc, broker)
	logrus.Infof("WebSocket server listing on %v", *restAddr)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", *restAddr), h); err != nil {
		panic(err)
	}
}
