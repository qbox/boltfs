package qfusegate

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

func handleInitRequest(ctx Context, host string, req *fuse.InitRequest) {

	client := rpc.DefaultClient

	args := &InitRequest{}
	_ = args

	n := unsafe.Sizeof(args)
	body := toReader(unsafe.Pointer(&args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/init", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	ret := &InitResponse{}
	_ = ret

	req.Respond(&fuse.InitResponse{})
}

func handleDestroyRequest(ctx Context, host string, req *fuse.DestroyRequest) {

	client := rpc.DefaultClient

	resp, err := client.DoRequest(ctx, "POST", host + "/v1/destroy")
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	req.Respond()
}

func handleStatfsRequest(ctx Context, host string, req *fuse.StatfsRequest) {

	client := rpc.DefaultClient

	resp, err := client.DoRequest(ctx, "POST", host + "/v1/statfs")
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	ret := &StatfsResponse{}
	_ = ret

	req.Respond(&fuse.StatfsResponse{})
}

func handleAccessRequest(ctx Context, host string, req *fuse.AccessRequest) {

	client := rpc.DefaultClient

	args := &AccessRequest{}
	_ = args

	n := unsafe.Sizeof(args)
	body := toReader(unsafe.Pointer(&args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/access", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	req.Respond()
}

func handleGetattrRequest(ctx Context, host string, req *fuse.GetattrRequest) {

	client := rpc.DefaultClient

	args := &GetattrRequest{}
	_ = args

	n := unsafe.Sizeof(args)
	body := toReader(unsafe.Pointer(&args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/getattr", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	ret := &GetattrResponse{}
	_ = ret

	req.Respond(&fuse.GetattrResponse{})
}

func handleListxattrRequest(ctx Context, host string, req *fuse.ListxattrRequest) {

	client := rpc.DefaultClient

	args := &ListxattrRequest{}
	_ = args

	n := unsafe.Sizeof(args)
	body := toReader(unsafe.Pointer(&args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/listxattr", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	ret := &ListxattrResponse{}
	_ = ret

	req.Respond(&fuse.ListxattrResponse{})
}

func handleGetxattrRequest(ctx Context, host string, req *fuse.GetxattrRequest) {

	client := rpc.DefaultClient

	args := &GetxattrRequest{}
	_ = args

	n := unsafe.Offsetof(args.Name)
	body1 := toReader(unsafe.Pointer(&args), n)
	body := io.MultiReader(body1, strings.NewReader(args.Name))
	resp, err := client.DoRequestWith(
		ctx, "POST", host + "/v1/getxattr", "application/fuse", body, int(n)+len(args.Name))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	ret := &GetxattrResponse{}
	_ = ret

	req.Respond(&fuse.GetxattrResponse{})
}

func handleRemovexattrRequest(ctx Context, host string, req *fuse.RemovexattrRequest) {

	client := rpc.DefaultClient

	args := &RemovexattrRequest{}
	_ = args

	n := unsafe.Offsetof(args.Name)
	body1 := toReader(unsafe.Pointer(&args), n)
	body := io.MultiReader(body1, strings.NewReader(args.Name))
	resp, err := client.DoRequestWith(
		ctx, "POST", host + "/v1/removexattr", "application/fuse", body, int(n)+len(args.Name))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	req.Respond()
}

func handleSetxattrRequest(ctx Context, host string, req *fuse.SetxattrRequest) {

	client := rpc.DefaultClient

	args := &SetxattrRequest{}
	_ = args

	n := unsafe.Offsetof(args.Name)
	body1 := toReader(unsafe.Pointer(&args), n)

	var encoder stringEncoder
	encoder.PutString(args.Name)
	encoder.PutBytes(args.Xattr)

	body := io.MultiReader(body1, &encoder.Buffer)
	resp, err := client.DoRequestWith(
		ctx, "POST", host + "/v1/setxattr", "application/fuse", body, int(n)+encoder.Len())
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	req.Respond()
}

func handleLookupRequest(ctx Context, host string, req *fuse.LookupRequest) {

	client := rpc.DefaultClient

	args := &LookupRequest{}
	_ = args

	n := unsafe.Offsetof(args.Name)
	body1 := toReader(unsafe.Pointer(&args), n)
	body := io.MultiReader(body1, strings.NewReader(args.Name))
	resp, err := client.DoRequestWith(
		ctx, "POST", host + "/v1/lookup", "application/fuse", body, int(n)+len(args.Name))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	ret := &LookupResponse{}
	_ = ret

	req.Respond(&fuse.LookupResponse{})
}

func handleOpenRequest(ctx Context, host string, req *fuse.OpenRequest) {

	client := rpc.DefaultClient

	args := &OpenRequest{}
	_ = args

	n := unsafe.Sizeof(args)
	body := toReader(unsafe.Pointer(&args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/open", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	ret := &OpenResponse{}
	_ = ret

	req.Respond(&fuse.OpenResponse{})
}

func handleCreateRequest(ctx Context, host string, req *fuse.CreateRequest) {

	client := rpc.DefaultClient

	args := &CreateRequest{}
	_ = args

	n := unsafe.Offsetof(args.Name)
	body1 := toReader(unsafe.Pointer(&args), n)
	body := io.MultiReader(body1, strings.NewReader(args.Name))
	resp, err := client.DoRequestWith(
		ctx, "POST", host + "/v1/create", "application/fuse", body, int(n)+len(args.Name))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	ret := &CreateResponse{}
	_ = ret

	req.Respond(&fuse.CreateResponse{})
}

func handleMkdirRequest(ctx Context, host string, req *fuse.MkdirRequest) {

	client := rpc.DefaultClient

	args := &MkdirRequest{}
	_ = args

	n := unsafe.Offsetof(args.Name)
	body1 := toReader(unsafe.Pointer(&args), n)
	body := io.MultiReader(body1, strings.NewReader(args.Name))
	resp, err := client.DoRequestWith(
		ctx, "POST", host + "/v1/mkdir", "application/fuse", body, int(n)+len(args.Name))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	ret := &MkdirResponse{}
	_ = ret

	req.Respond(&fuse.MkdirResponse{})
}

func handleSymlinkRequest(ctx Context, host string, req *fuse.SymlinkRequest) {

	client := rpc.DefaultClient

	args := &SymlinkRequest{}
	_ = args

	n := unsafe.Offsetof(args.NewName)
	body1 := toReader(unsafe.Pointer(&args), n)

	var encoder stringEncoder
	encoder.PutString(args.NewName)
	encoder.PutString(args.Target)

	body := io.MultiReader(body1, &encoder.Buffer)
	resp, err := client.DoRequestWith(
		ctx, "POST", host + "/v1/symlink", "application/fuse", body, int(n)+encoder.Len())
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	ret := &SymlinkResponse{}
	_ = ret

	req.Respond(&fuse.SymlinkResponse{})
}

func handleReadlinkRequest(ctx Context, host string, req *fuse.ReadlinkRequest) {

	client := rpc.DefaultClient

	args := &ReadlinkRequest{}
	_ = args

	n := unsafe.Sizeof(args)
	body := toReader(unsafe.Pointer(&args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/readlink", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	ret := &ReadlinkResponse{}
	_ = ret

	req.Respond(ret.Target)
}

func handleLinkRequest(ctx Context, host string, req *fuse.LinkRequest) {

	client := rpc.DefaultClient

	args := &LinkRequest{}
	_ = args

	n := unsafe.Offsetof(args.NewName)
	body1 := toReader(unsafe.Pointer(&args), n)
	body := io.MultiReader(body1, strings.NewReader(args.NewName))
	resp, err := client.DoRequestWith(
		ctx, "POST", host + "/v1/link", "application/fuse", body, int(n)+len(args.NewName))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	ret := &LinkResponse{}
	_ = ret

	req.Respond(&fuse.LookupResponse{})
}

func handleMknodRequest(ctx Context, host string, req *fuse.MknodRequest) {

	client := rpc.DefaultClient

	args := &MknodRequest{}
	_ = args

	n := unsafe.Offsetof(args.Name)
	body1 := toReader(unsafe.Pointer(&args), n)
	body := io.MultiReader(body1, strings.NewReader(args.Name))
	resp, err := client.DoRequestWith(
		ctx, "POST", host + "/v1/mknod", "application/fuse", body, int(n)+len(args.Name))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	ret := &MknodResponse{}
	_ = ret

	req.Respond(&fuse.LookupResponse{})
}

func handleRenameRequest(ctx Context, host string, req *fuse.RenameRequest) {

	client := rpc.DefaultClient

	args := &RenameRequest{}
	_ = args

	n := unsafe.Offsetof(args.OldName)
	body1 := toReader(unsafe.Pointer(&args), n)

	var encoder stringEncoder
	encoder.PutString(args.OldName)
	encoder.PutString(args.NewName)

	body := io.MultiReader(body1, &encoder.Buffer)
	resp, err := client.DoRequestWith(
		ctx, "POST", host + "/v1/rename", "application/fuse", body, int(n)+encoder.Len())
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	req.Respond()
}

func handleRemoveRequest(ctx Context, host string, req *fuse.RemoveRequest) {

	client := rpc.DefaultClient

	args := &RemoveRequest{}
	_ = args

	n := unsafe.Offsetof(args.Name)
	body1 := toReader(unsafe.Pointer(&args), n)
	body := io.MultiReader(body1, strings.NewReader(args.Name))
	resp, err := client.DoRequestWith(
		ctx, "POST", host + "/v1/remove", "application/fuse", body, int(n)+len(args.Name))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	req.Respond()
}

func handleReadRequest(ctx Context, host string, req *fuse.ReadRequest) {

	client := rpc.DefaultClient

	args := &ReadRequest{}
	_ = args

	n := unsafe.Sizeof(args)
	body := toReader(unsafe.Pointer(&args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/read", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	ret := &ReadResponse{}
	_ = ret

	req.Respond(&fuse.ReadResponse{})
}

func handleWriteRequest(ctx Context, host string, req *fuse.WriteRequest) {

	client := rpc.DefaultClient

	args := &WriteRequest{}
	_ = args

	n := unsafe.Offsetof(args.Data)
	body1 := toReader(unsafe.Pointer(&args), n)
	body := io.MultiReader(body1, bytes.NewReader(args.Data))
	resp, err := client.DoRequestWith(
		ctx, "POST", host + "/v1/write", "application/fuse", body, int(n)+len(args.Data))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	ret := &WriteResponse{}
	_ = ret

	req.Respond(&fuse.WriteResponse{})
}

func handleSetattrRequest(ctx Context, host string, req *fuse.SetattrRequest) {

	client := rpc.DefaultClient

	args := &SetattrRequest{}
	_ = args

	n := unsafe.Sizeof(args)
	body := toReader(unsafe.Pointer(&args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/setattr", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	ret := &SetattrResponse{}
	_ = ret

	req.Respond(&fuse.SetattrResponse{})
}

func handleFlushRequest(ctx Context, host string, req *fuse.FlushRequest) {

	client := rpc.DefaultClient

	args := &FlushRequest{}
	_ = args

	n := unsafe.Sizeof(args)
	body := toReader(unsafe.Pointer(&args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/flush", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	req.Respond()
}

func handleFsyncRequest(ctx Context, host string, req *fuse.FsyncRequest) {

	client := rpc.DefaultClient

	args := &FsyncRequest{}
	_ = args

	n := unsafe.Sizeof(args)
	body := toReader(unsafe.Pointer(&args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/fsync", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	req.Respond()
}

func handleReleaseRequest(ctx Context, host string, req *fuse.ReleaseRequest) {

	client := rpc.DefaultClient

	args := &ReleaseRequest{}
	_ = args

	n := unsafe.Sizeof(args)
	body := toReader(unsafe.Pointer(&args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/release", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	req.Respond()
}

func handleForgetRequest(ctx Context, host string, req *fuse.ForgetRequest) {

	client := rpc.DefaultClient

	args := &ForgetRequest{}
	_ = args

	n := unsafe.Sizeof(args)
	body := toReader(unsafe.Pointer(&args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/forget", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	req.Respond()
}

func handleInterruptRequest(ctx Context, host string, req *fuse.InterruptRequest) {

	client := rpc.DefaultClient

	args := &InterruptRequest{}
	_ = args

	n := unsafe.Sizeof(args)
	body := toReader(unsafe.Pointer(&args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/interrupt", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	req.Respond()
}

