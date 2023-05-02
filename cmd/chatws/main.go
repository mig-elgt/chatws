package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/mig-elgt/chatws/handler"
	"github.com/mig-elgt/chatws/rabbitmq"
)

var (
	rabbitmqHost = os.Getenv("RABBITMQ_HOST")
	rabbitmqPort = os.Getenv("RABBITMQ_PORT")
	rabbitmqUser = os.Getenv("RABBITMQ_USER")
	rabbitmqPass = os.Getenv("RABBITMQ_PASSWORD")
)

func main() {
	restAddr := flag.Int("addr", 8080, "web socket server port")
	flag.Parse()
	broker, err := rabbitmq.New(rabbitmqUser, rabbitmqPass, rabbitmqHost, rabbitmqPort)
	if err != nil {
		panic(err)
	}
	h := handler.New(nil, broker)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", *restAddr), h); err != nil {
		panic(err)
	}
}
