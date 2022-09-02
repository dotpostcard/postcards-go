package compiler

import (
	"context"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ds-flatfs"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipfs/go-merkledag"
	car2 "github.com/ipld/go-car"
	"io"
)

func (pc *Postcard) Compile(output io.WriteCloser) error {
	ctx := context.Background()

	ds := flatfs.Create("/tmp/jp", flatfs.NextToLast(2))
	//???
	bst := blockstore.NewBlockstore()
	bse := blockservice.New(bst, nil)
	das := merkledag.NewDAGService(bse)

	return car2.WriteCar(ctx, das, []cid.Cid{root}, output)
}
