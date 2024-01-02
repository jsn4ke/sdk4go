package main

import (
	raft_ext "github.com/jsn4ke/sdk4go/cmd/protoc-gen-raftext/internal"
	"google.golang.org/protobuf/compiler/protogen"
)

// protoc -I? --plugin=protoc-gen-go=$HOME/go/bin/protoc-gen-go --go_out=? ?.proto
// protoc -I? --plugin=./protoc-gen-raftext --raftext_out=? ?.proto
func main() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		for _, v := range gen.Files {
			if v.Generate {
				raft_ext.Gen(gen, v)
			}
		}
		return nil
	})
}
