package main

import (
	"bazil.org/fuse"
	"fmt"
	"reflect"
	"strings"

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

func isFlatType(t reflect.Type) bool {

	kind := t.Kind()
	if kind <= reflect.Complex128 {
		return true
	}
	if kind != reflect.Struct {
		return false
	}

	n := t.NumField()
	for i := 0; i < n; i++ {
		f := t.Field(i)
		if !isFlatType(f.Type) {
			return false
		}
	}
	return true
}

type nonFlatTypeInfo struct {
	Type   string  // type name of first NonFlatType
	Name   string  // field name of first NonFlatType
	Idx    int     // field index of first NonFlatType
	Count  int     // count of NonFlatType
	Offset uintptr // offset of first NonFlatType
}

func nonFlatTypeInfoOf(t reflect.Type) nonFlatTypeInfo {

	if t.Kind() != reflect.Struct {
		panic("t != struct")
	}

	n := t.NumField()
	for i := 0; i < n; i++ {
		f := t.Field(i)
		if !isFlatType(f.Type) {
			return nonFlatTypeInfo{nonFlatTypeOf(f.Type), f.Name, i, n-i, f.Offset}
		}
	}
	panic("nonFlatTypeInfoOf: unexpected")
}

func nonFlatTypeOf(t reflect.Type) string {

	switch t.Kind() {
	case reflect.String:
		return "strings"
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return "bytes"
		}
	}
	println("nonFlatTypeOf:", t.String())
	panic("nonFlatTypeOf: unexpected")
}

func gen(types []interface{}) {

	req, resp := typeOf(types[0]), typeOf(types[2])
	fuseReq, fuseResp := typeOf(types[1]), typeOf(types[3])

	reqName := fuseReq.Name()
	reqPath := "/v1/" + strings.ToLower(strings.TrimSuffix(reqName, "Request"))
	fmt.Printf(`func handle%s(ctx Context, host string, req *fuse.%s) {

	client := rpc.DefaultClient

`, reqName, reqName)

	if req == nil {
		fmt.Printf("\tresp, err := client.DoRequest(ctx, \"POST\", host + \"%s\")\n", reqPath)
	} else {
		argsName := req.Name()
		fmt.Printf("\targs := &%s{}\n\t_ = args\n", argsName)
		if isFlatType(req) {
			fmt.Printf(`
	n := unsafe.Sizeof(args)
	body := toReader(unsafe.Pointer(&args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "%s", "application/fuse", body, int(n))
`, reqPath)
		} else {
			info := nonFlatTypeInfoOf(req)
			if info.Count == 1 && info.Offset == 0 {
				fmt.Printf(`
	body := %s.NewReader(args.%s)
	resp, err := client.DoRequestWith(
		ctx, "POST", host + "%s", "application/fuse", body, len(args.%s))
`, info.Type, info.Name, reqPath, info.Name)
			} else if info.Count == 1 && info.Offset != 0 {
				fmt.Printf(`
	n := unsafe.Offsetof(args.%s)
	body1 := toReader(unsafe.Pointer(&args), n)
	body := io.MultiReader(body1, %s.NewReader(args.%s))
	resp, err := client.DoRequestWith(
		ctx, "POST", host + "%s", "application/fuse", body, int(n)+len(args.%s))
`, info.Name, info.Type, info.Name, reqPath, info.Name)
			} else {
				fmt.Printf(`
	n := unsafe.Offsetof(args.%s)
	body1 := toReader(unsafe.Pointer(&args), n)

	var encoder stringEncoder
`, info.Name)
				for i := 0; i < info.Count; i++ {
					f := req.Field(info.Idx + i)
					switch {
					case f.Type.Kind() == reflect.String:
						fmt.Printf(
`	encoder.PutString(args.%s)
`, f.Name)
					case f.Type.Kind() == reflect.Slice && f.Type.Elem().Kind() == reflect.Uint8:
						fmt.Printf(
`	encoder.PutBytes(args.%s)
`, f.Name)
					default:
						println("type:", req.String())
						panic("field must be string or []byte")
					}
				}
				fmt.Printf(`
	body := io.MultiReader(body1, &encoder.Buffer)
	resp, err := client.DoRequestWith(
		ctx, "POST", host + "%s", "application/fuse", body, int(n)+encoder.Len())
`, reqPath)
			}
		}
	}
	fmt.Printf(`	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

`)

	if resp == nil {
		fmt.Printf("\treq.Respond()\n}\n\n")
		return
	}

	retName := resp.Name()
	fmt.Printf("\tret := &%s{}\n\t_ = ret\n\n", retName)

	if fuseResp.Kind() == reflect.String {
		fmt.Printf("\treq.Respond(ret.Target)\n}\n\n")
		return
	}

	respName := fuseResp.Name()
	fmt.Printf("\treq.Respond(&fuse.%s{})\n}\n\n", respName)
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
	"bytes"
	"io"
	"io/ioutil"
	"qiniupkg.com/x/rpc.v7"
	"strings"
	"unsafe"

	. "golang.org/x/net/context"
	. "qiniu.com/qfuse.proto.v1"
)

`)

	for i := 0; i < n; i += 4 {
		gen(types[i:])
	}
}

// ---------------------------------------------------------------------------

