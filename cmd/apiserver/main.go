package main

import (
	"github.com/yiran15/api-server/cmd"
)

// @title           Swagger API
// @version         1.0
// @description     api-server api docs.
// @host      10.0.0.10:8080
func main() {
	if err := cmd.NewCmd().Execute(); err != nil {
		panic(err)
	}
}
