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
	// Output: hello.postcard has checksum 7a9bc161678b39e2841f7e8f08026119
}
