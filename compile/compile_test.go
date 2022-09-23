package compile_test

import (
	"crypto/md5"
	"fmt"

	"github.com/jphastings/postcard-go/compile"
)

func ExampleFromFiles() {
	filename, data, err := compile.FromFiles("../fixtures/hello-meta.yaml")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s has checksum %x", filename, md5.Sum(data))
	// Output: hello.postcard has checksum 3b701cc7001611bac858d8b824573c22
}
