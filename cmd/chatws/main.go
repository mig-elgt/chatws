package main

import (
	"net/http"

	"github.com/mig-elgt/chatws/handler"
)

func main() {
	h := handler.New(nil, nil)
	if err := http.ListenAndServe(":8080", h); err != nil {
		panic(err)
	}
}
