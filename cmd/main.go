package cmd

import "github.com/lvkeliang/WHOIM/api"

func main() {
	api.InitRPC()
	api.InitRouter()
}
