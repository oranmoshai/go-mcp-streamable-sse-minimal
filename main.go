package main

import "github.com/oranmoshai/go-mcp-streamable-sse-minimal/service"

func main() {

	s := service.New()
	err := s.Run()
	if err != nil {
		panic(err)
	}

}
