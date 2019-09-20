package main

import (
	"github.com/art-frela/HW3/karpov/infra"
)

func main(){
	server := infra.NewBlogServer()
	server.Run(":8888")
}
