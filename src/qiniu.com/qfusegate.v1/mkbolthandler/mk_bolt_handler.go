package main

import (
	"bazil.org/fuse"
	"fmt"
	"reflect"

	. "qiniu.com/qfuse.proto.v1"
)

// ---------------------------------------------------------------------------

var types = []interface{}{
	new(InitRequest),
	new(fuse.InitRequest),
	new(InitResponse),
	new(fuse.InitResponse),

	nil,
	new(fuse.DestroyRequest),
	nil,
	nil,

	nil,
	new(fuse.StatfsRequest),
	new(StatfsResponse),
	new(fuse.StatfsResponse),

	new(AccessRequest),
	new(fuse.AccessRequest),
	nil,
	nil,

	new(GetattrRequest),
	new(fuse.GetattrRequest),
	new(GetattrResponse),
	new(fuse.GetattrResponse),

	new(ListxattrRequest),
	new(fuse.ListxattrRequest),
	new(ListxattrResponse),
	new(fuse.ListxattrResponse),

	new(GetxattrRequest),
	new(fuse.GetxattrRequest),
	new(GetxattrResponse),
	new(fuse.GetxattrResponse),

	new(RemovexattrRequest),
	new(fuse.RemovexattrRequest),
	nil,
	nil,

	new(SetxattrRequest),
	new(fuse.SetxattrRequest),
	nil,
	nil,

	new(LookupRequest),
	new(fuse.LookupRequest),
	new(LookupResponse),
	new(fuse.LookupResponse),

	new(OpenRequest),
	new(fuse.OpenRequest),
	new(OpenResponse),
	new(fuse.OpenResponse),

	new(CreateRequest),
	new(fuse.CreateRequest),
	new(CreateResponse),
	new(fuse.CreateResponse),

	new(MkdirRequest),
	new(fuse.MkdirRequest),
	new(MkdirResponse),
	new(fuse.MkdirResponse),

	new(SymlinkRequest),
	new(fuse.SymlinkRequest),
	new(SymlinkResponse),
	new(fuse.SymlinkResponse),

	new(ReadlinkRequest),
	new(fuse.ReadlinkRequest),
	new(ReadlinkResponse),
	new(string),

	new(LinkRequest),
	new(fuse.LinkRequest),
	new(LinkResponse),
	new(fuse.LookupResponse),

	new(MknodRequest),
	new(fuse.MknodRequest),
	new(MknodResponse),
	new(fuse.LookupResponse),

	new(RenameRequest),
	new(fuse.RenameRequest),
	nil,
	nil,

	new(RemoveRequest),
	new(fuse.RemoveRequest),
	nil,
	nil,

	new(ReadRequest),
	new(fuse.ReadRequest),
	new(ReadResponse),
	new(fuse.ReadResponse),

	new(WriteRequest),
	new(fuse.WriteRequest),
	new(WriteResponse),
	new(fuse.WriteResponse),

	new(SetattrRequest),
	new(fuse.SetattrRequest),
	new(SetattrResponse),
	new(fuse.SetattrResponse),

	new(FlushRequest),
	new(fuse.FlushRequest),
	nil,
	nil,

	new(FsyncRequest),
	new(fuse.FsyncRequest),
	nil,
	nil,

	new(ReleaseRequest),
	new(fuse.ReleaseRequest),
	nil,
	nil,

	new(ForgetRequest),
	new(fuse.ForgetRequest),
	nil,
	nil,

	new(InterruptRequest),
	new(fuse.InterruptRequest),
	nil,
	nil,
}

// ---------------------------------------------------------------------------

func gen(types []interface{}) {

	req, resp := typeOf(types[0]), typeOf(types[2])
	fuseReq, fuseResp := typeOf(types[1]), typeOf(types[3])

	_ = req
	_ = resp
	_ = fuseReq
	_ = fuseResp

	reqName := fuseReq.Name()
	fmt.Printf("func handle%s(ctx Context, target string, r *fuse.%s) {\n}\n\n", reqName, reqName)
}

func typeOf(v interface{}) reflect.Type {

	if v == nil {
		return nil
	}
	return reflect.TypeOf(v).Elem()
}

func main() {

	n := len(types)
	if n % 4 != 0 {
		panic("invalid types")
	}

	fmt.Printf(`package qfusegate

import (
	"bazil.org/fuse"
	. "golang.org/x/net/context"
)

`)

	for i := 0; i < n; i += 4 {
		gen(types[i:])
	}
}

// ---------------------------------------------------------------------------

