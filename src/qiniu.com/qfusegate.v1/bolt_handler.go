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

	args := &InitRequest{
		Major: req.Major,
		Minor: req.Minor,
		MaxReadahead: req.MaxReadahead,
		Flags: req.Flags,
	}

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

	args := &AccessRequest{
		Inode: uint64(req.Node),
		Mask: req.Mask,
	}

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

	args := &GetattrRequest{
		Inode: uint64(req.Node),
	}

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

	args := &ListxattrRequest{
		Inode: uint64(req.Node),
		Size: req.Size,
		Position: req.Position,
	}

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

	args := &GetxattrRequest{
		Inode: uint64(req.Node),
		Size: req.Size,
		Position: req.Position,
		Name: req.Name,
	}

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

	args := &RemovexattrRequest{
		Inode: uint64(req.Node),
		Name: req.Name,
	}

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

	args := &SetxattrRequest{
		Inode: uint64(req.Node),
		Flags: req.Flags,
		Position: req.Position,
		Name: req.Name,
		Xattr: req.Xattr,
	}

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

	args := &LookupRequest{
		Inode: uint64(req.Node),
		Name: req.Name,
	}

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

	args := &OpenRequest{
		Inode: uint64(req.Node),
		Flags: req.Flags,
		Dir: req.Dir,
	}

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

	args := &CreateRequest{
		Inode: uint64(req.Node),
		Flags: req.Flags,
		Mode: req.Mode,
		Name: req.Name,
	}

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

	args := &MkdirRequest{
		Inode: uint64(req.Node),
		Mode: req.Mode,
		Name: req.Name,
	}

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

	args := &SymlinkRequest{
		Inode: uint64(req.Node),
		NewName: req.NewName,
		Target: req.Target,
	}

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

	args := &ReadlinkRequest{
		Inode: uint64(req.Node),
	}

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

	args := &LinkRequest{
		Inode: uint64(req.Node),
		OldInode: uint64(req.OldNode),
		NewName: req.NewName,
	}

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

	args := &MknodRequest{
		Inode: uint64(req.Node),
		Mode: req.Mode,
		Rdev: req.Rdev,
		Name: req.Name,
	}

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

	args := &RenameRequest{
		NewDirInode: uint64(req.NewDir),
		OldName: req.OldName,
		NewName: req.NewName,
	}

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

	args := &RemoveRequest{
		Inode: uint64(req.Node),
		Dir: req.Dir,
		Name: req.Name,
	}

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

	args := &ReadRequest{
		Handle: uint64(req.Handle),
		Offset: req.Offset,
		Size: req.Size,
		Dir: req.Dir,
	}

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

	args := &WriteRequest{
		Handle: uint64(req.Handle),
		Offset: req.Offset,
		Flags: req.Flags,
		Data: req.Data,
	}

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

	args := &SetattrRequest{
		Valid: req.Valid,
		Handle: uint64(req.Handle),
		Size: req.Size,
		Atime: Time(req.Atime.UnixNano()),
		Mtime: Time(req.Mtime.UnixNano()),
		Mode: req.Mode,
		Uid: req.Uid,
		Gid: req.Gid,
		Bkuptime: Time(req.Bkuptime.UnixNano()),
		Chgtime: Time(req.Chgtime.UnixNano()),
		Crtime: Time(req.Crtime.UnixNano()),
		Flags: req.Flags,
	}

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

	args := &FlushRequest{
		Handle: uint64(req.Handle),
		LockOwner: req.LockOwner,
		Flags: req.Flags,
	}

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

	args := &FsyncRequest{
		Handle: uint64(req.Handle),
		Flags: req.Flags,
		Dir: req.Dir,
	}

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

	args := &ReleaseRequest{
		Handle: uint64(req.Handle),
		Flags: req.Flags,
		ReleaseFlags: req.ReleaseFlags,
		LockOwner: req.LockOwner,
		Dir: req.Dir,
	}

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

	args := &ForgetRequest{
		Inode: uint64(req.Node),
		LookupReqid: uint64(req.N),
	}

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

	args := &InterruptRequest{
		IntrReqId: uint64(req.IntrID),
	}

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

