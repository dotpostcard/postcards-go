package postcard_test

import (
	"crypto/md5"
	"fmt"

	"github.com/jphastings/postcard-go"
)

func ExampleCompileFiles() {
	filename, data, err := postcard.CompileFiles("fixtures/hello-meta.yaml")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s has checksum %x", filename, md5.Sum(data))
	// Output: hello.postcard has checksum ff4a9cfa149751cf0a9d68592da92f3c
}
