package qfusegate

import (
	"bazil.org/fuse"
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"unsafe"

	. "golang.org/x/net/context"
	. "qiniu.com/qfuse.proto.v1"
)

func handleInitRequest(ctx Context, host string, req *fuse.InitRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &InitRequest{
		Major: req.Major,
		Minor: req.Minor,
		MaxReadahead: req.MaxReadahead,
		Flags: req.Flags,
	}

	n := unsafe.Sizeof(*args)
	body := toReader(unsafe.Pointer(args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/init", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	ret := new(InitResponse)
	err = fromReader(unsafe.Pointer(ret), unsafe.Sizeof(*ret), resp.Body)
	if err != nil {
		replyError(req, err)
		return
	}

	fuseResp := new(fuse.InitResponse)
	fuseResp.MaxReadahead = ret.MaxReadahead
	fuseResp.Flags = ret.Flags
	fuseResp.MaxWrite = ret.MaxWrite
	req.Respond(fuseResp)
}

func handleDestroyRequest(ctx Context, host string, req *fuse.DestroyRequest) {

	client := newBoltClient(&req.Header, nil)

	resp, err := client.DoRequest(ctx, "POST", host + "/v1/destroy")
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	req.Respond()
}

func handleStatfsRequest(ctx Context, host string, req *fuse.StatfsRequest) {

	client := newBoltClient(&req.Header, nil)

	resp, err := client.DoRequest(ctx, "POST", host + "/v1/statfs")
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	ret := new(StatfsResponse)
	err = fromReader(unsafe.Pointer(ret), unsafe.Sizeof(*ret), resp.Body)
	if err != nil {
		replyError(req, err)
		return
	}

	fuseResp := new(fuse.StatfsResponse)
	fuseResp.Blocks = ret.Blocks
	fuseResp.Bfree = ret.Bfree
	fuseResp.Bavail = ret.Bavail
	fuseResp.Files = ret.Files
	fuseResp.Ffree = ret.Ffree
	fuseResp.Bsize = ret.Bsize
	fuseResp.Namelen = ret.Namelen
	fuseResp.Frsize = ret.Frsize
	req.Respond(fuseResp)
}

func handleAccessRequest(ctx Context, host string, req *fuse.AccessRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &AccessRequest{
		Inode: uint64(req.Node),
		Mask: req.Mask,
	}

	n := unsafe.Sizeof(*args)
	body := toReader(unsafe.Pointer(args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/access", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	req.Respond()
}

func handleGetattrRequest(ctx Context, host string, req *fuse.GetattrRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &GetattrRequest{
		Inode: uint64(req.Node),
	}

	n := unsafe.Sizeof(*args)
	body := toReader(unsafe.Pointer(args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/getattr", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	ret := new(GetattrResponse)
	err = fromReader(unsafe.Pointer(ret), unsafe.Sizeof(*ret), resp.Body)
	if err != nil {
		replyError(req, err)
		return
	}

	fuseResp := new(fuse.GetattrResponse)
	assignAttr(&fuseResp.Attr, &ret.Attr)
	req.Respond(fuseResp)
}

func handleListxattrRequest(ctx Context, host string, req *fuse.ListxattrRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &ListxattrRequest{
		Inode: uint64(req.Node),
		Size: req.Size,
		Position: req.Position,
	}

	n := unsafe.Sizeof(*args)
	body := toReader(unsafe.Pointer(args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/listxattr", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	ret := new(ListxattrResponse)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		replyError(req, err)
		return
	}
	ret.XattrNames = b

	fuseResp := new(fuse.ListxattrResponse)
	fuseResp.Xattr = ret.XattrNames
	req.Respond(fuseResp)
}

func handleGetxattrRequest(ctx Context, host string, req *fuse.GetxattrRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &GetxattrRequest{
		Inode: uint64(req.Node),
		Size: req.Size,
		Position: req.Position,
		Name: req.Name,
	}

	n := unsafe.Offsetof(args.Name)
	body1 := toReader(unsafe.Pointer(args), n)
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

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	ret := new(GetxattrResponse)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		replyError(req, err)
		return
	}
	ret.Xattr = b

	fuseResp := new(fuse.GetxattrResponse)
	fuseResp.Xattr = ret.Xattr
	req.Respond(fuseResp)
}

func handleRemovexattrRequest(ctx Context, host string, req *fuse.RemovexattrRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &RemovexattrRequest{
		Inode: uint64(req.Node),
		Name: req.Name,
	}

	n := unsafe.Offsetof(args.Name)
	body1 := toReader(unsafe.Pointer(args), n)
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

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	req.Respond()
}

func handleSetxattrRequest(ctx Context, host string, req *fuse.SetxattrRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &SetxattrRequest{
		Inode: uint64(req.Node),
		Flags: req.Flags,
		Position: req.Position,
		Name: req.Name,
		Xattr: req.Xattr,
	}

	n := unsafe.Offsetof(args.Name)
	body1 := toReader(unsafe.Pointer(args), n)

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

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	req.Respond()
}

func handleLookupRequest(ctx Context, host string, req *fuse.LookupRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &LookupRequest{
		Inode: uint64(req.Node),
		Name: req.Name,
	}

	n := unsafe.Offsetof(args.Name)
	body1 := toReader(unsafe.Pointer(args), n)
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

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	ret := new(LookupResponse)
	err = fromReader(unsafe.Pointer(ret), unsafe.Sizeof(*ret), resp.Body)
	if err != nil {
		replyError(req, err)
		return
	}

	fuseResp := new(fuse.LookupResponse)
	fuseResp.Node = fuse.NodeID(ret.Inode)
	fuseResp.Generation = ret.Generation
	fuseResp.EntryValid = ret.EntryValid
	assignAttr(&fuseResp.Attr, &ret.Attr)
	req.Respond(fuseResp)
}

func handleOpenRequest(ctx Context, host string, req *fuse.OpenRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &OpenRequest{
		Inode: uint64(req.Node),
		Flags: req.Flags,
		Dir: req.Dir,
	}

	n := unsafe.Sizeof(*args)
	body := toReader(unsafe.Pointer(args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/open", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	ret := new(OpenResponse)
	err = fromReader(unsafe.Pointer(ret), unsafe.Sizeof(*ret), resp.Body)
	if err != nil {
		replyError(req, err)
		return
	}

	fuseResp := new(fuse.OpenResponse)
	fuseResp.Handle = fuse.HandleID(ret.Handle)
	fuseResp.Flags = ret.Flags
	req.Respond(fuseResp)
}

func handleCreateRequest(ctx Context, host string, req *fuse.CreateRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &CreateRequest{
		Inode: uint64(req.Node),
		Flags: req.Flags,
		Mode: req.Mode,
		Name: req.Name,
	}

	n := unsafe.Offsetof(args.Name)
	body1 := toReader(unsafe.Pointer(args), n)
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

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	ret := new(CreateResponse)
	err = fromReader(unsafe.Pointer(ret), unsafe.Sizeof(*ret), resp.Body)
	if err != nil {
		replyError(req, err)
		return
	}

	fuseResp := new(fuse.CreateResponse)
	fuseResp.Node = fuse.NodeID(ret.Inode)
	fuseResp.Generation = ret.Generation
	fuseResp.EntryValid = ret.EntryValid
	assignAttr(&fuseResp.Attr, &ret.Attr)
	fuseResp.Handle = fuse.HandleID(ret.Handle)
	fuseResp.Flags = ret.Flags
	req.Respond(fuseResp)
}

func handleMkdirRequest(ctx Context, host string, req *fuse.MkdirRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &MkdirRequest{
		Inode: uint64(req.Node),
		Mode: req.Mode,
		Name: req.Name,
	}

	n := unsafe.Offsetof(args.Name)
	body1 := toReader(unsafe.Pointer(args), n)
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

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	ret := new(MkdirResponse)
	err = fromReader(unsafe.Pointer(ret), unsafe.Sizeof(*ret), resp.Body)
	if err != nil {
		replyError(req, err)
		return
	}

	fuseResp := new(fuse.MkdirResponse)
	fuseResp.Node = fuse.NodeID(ret.Inode)
	fuseResp.Generation = ret.Generation
	fuseResp.EntryValid = ret.EntryValid
	assignAttr(&fuseResp.Attr, &ret.Attr)
	req.Respond(fuseResp)
}

func handleSymlinkRequest(ctx Context, host string, req *fuse.SymlinkRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &SymlinkRequest{
		Inode: uint64(req.Node),
		NewName: req.NewName,
		Target: req.Target,
	}

	n := unsafe.Offsetof(args.NewName)
	body1 := toReader(unsafe.Pointer(args), n)

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

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	ret := new(SymlinkResponse)
	err = fromReader(unsafe.Pointer(ret), unsafe.Sizeof(*ret), resp.Body)
	if err != nil {
		replyError(req, err)
		return
	}

	fuseResp := new(fuse.SymlinkResponse)
	fuseResp.Node = fuse.NodeID(ret.Inode)
	fuseResp.Generation = ret.Generation
	fuseResp.EntryValid = ret.EntryValid
	assignAttr(&fuseResp.Attr, &ret.Attr)
	req.Respond(fuseResp)
}

func handleReadlinkRequest(ctx Context, host string, req *fuse.ReadlinkRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &ReadlinkRequest{
		Inode: uint64(req.Node),
	}

	n := unsafe.Sizeof(*args)
	body := toReader(unsafe.Pointer(args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/readlink", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	ret := new(ReadlinkResponse)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		replyError(req, err)
		return
	}
	ret.Target = string(b)
	req.Respond(ret.Target)
}

func handleLinkRequest(ctx Context, host string, req *fuse.LinkRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &LinkRequest{
		Inode: uint64(req.Node),
		OldInode: uint64(req.OldNode),
		NewName: req.NewName,
	}

	n := unsafe.Offsetof(args.NewName)
	body1 := toReader(unsafe.Pointer(args), n)
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

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	ret := new(LinkResponse)
	err = fromReader(unsafe.Pointer(ret), unsafe.Sizeof(*ret), resp.Body)
	if err != nil {
		replyError(req, err)
		return
	}

	fuseResp := new(fuse.LookupResponse)
	fuseResp.Node = fuse.NodeID(ret.Inode)
	fuseResp.Generation = ret.Generation
	fuseResp.EntryValid = ret.EntryValid
	assignAttr(&fuseResp.Attr, &ret.Attr)
	req.Respond(fuseResp)
}

func handleMknodRequest(ctx Context, host string, req *fuse.MknodRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &MknodRequest{
		Inode: uint64(req.Node),
		Mode: req.Mode,
		Rdev: req.Rdev,
		Name: req.Name,
	}

	n := unsafe.Offsetof(args.Name)
	body1 := toReader(unsafe.Pointer(args), n)
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

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	ret := new(MknodResponse)
	err = fromReader(unsafe.Pointer(ret), unsafe.Sizeof(*ret), resp.Body)
	if err != nil {
		replyError(req, err)
		return
	}

	fuseResp := new(fuse.LookupResponse)
	fuseResp.Node = fuse.NodeID(ret.Inode)
	fuseResp.Generation = ret.Generation
	fuseResp.EntryValid = ret.EntryValid
	assignAttr(&fuseResp.Attr, &ret.Attr)
	req.Respond(fuseResp)
}

func handleRenameRequest(ctx Context, host string, req *fuse.RenameRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &RenameRequest{
		NewDirInode: uint64(req.NewDir),
		OldName: req.OldName,
		NewName: req.NewName,
	}

	n := unsafe.Offsetof(args.OldName)
	body1 := toReader(unsafe.Pointer(args), n)

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

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	req.Respond()
}

func handleRemoveRequest(ctx Context, host string, req *fuse.RemoveRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &RemoveRequest{
		Inode: uint64(req.Node),
		Dir: req.Dir,
		Name: req.Name,
	}

	n := unsafe.Offsetof(args.Name)
	body1 := toReader(unsafe.Pointer(args), n)
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

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	req.Respond()
}

func handleReadRequest(ctx Context, host string, req *fuse.ReadRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &ReadRequest{
		Handle: uint64(req.Handle),
		Offset: req.Offset,
		Size: req.Size,
		Dir: req.Dir,
	}

	n := unsafe.Sizeof(*args)
	body := toReader(unsafe.Pointer(args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/read", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	ret := new(ReadResponse)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		replyError(req, err)
		return
	}
	ret.Data = b

	fuseResp := new(fuse.ReadResponse)
	fuseResp.Data = ret.Data
	req.Respond(fuseResp)
}

func handleWriteRequest(ctx Context, host string, req *fuse.WriteRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &WriteRequest{
		Handle: uint64(req.Handle),
		Offset: req.Offset,
		Flags: req.Flags,
		Data: req.Data,
	}

	n := unsafe.Offsetof(args.Data)
	body1 := toReader(unsafe.Pointer(args), n)
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

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	ret := new(WriteResponse)
	err = fromReader(unsafe.Pointer(ret), unsafe.Sizeof(*ret), resp.Body)
	if err != nil {
		replyError(req, err)
		return
	}

	fuseResp := new(fuse.WriteResponse)
	fuseResp.Size = ret.Size
	req.Respond(fuseResp)
}

func handleSetattrRequest(ctx Context, host string, req *fuse.SetattrRequest) {

	client := newBoltClient(&req.Header, nil)

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

	n := unsafe.Sizeof(*args)
	body := toReader(unsafe.Pointer(args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/setattr", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	ret := new(SetattrResponse)
	err = fromReader(unsafe.Pointer(ret), unsafe.Sizeof(*ret), resp.Body)
	if err != nil {
		replyError(req, err)
		return
	}

	fuseResp := new(fuse.SetattrResponse)
	assignAttr(&fuseResp.Attr, &ret.Attr)
	req.Respond(fuseResp)
}

func handleFlushRequest(ctx Context, host string, req *fuse.FlushRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &FlushRequest{
		Handle: uint64(req.Handle),
		LockOwner: req.LockOwner,
		Flags: req.Flags,
	}

	n := unsafe.Sizeof(*args)
	body := toReader(unsafe.Pointer(args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/flush", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	req.Respond()
}

func handleFsyncRequest(ctx Context, host string, req *fuse.FsyncRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &FsyncRequest{
		Handle: uint64(req.Handle),
		Flags: req.Flags,
		Dir: req.Dir,
	}

	n := unsafe.Sizeof(*args)
	body := toReader(unsafe.Pointer(args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/fsync", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	req.Respond()
}

func handleReleaseRequest(ctx Context, host string, req *fuse.ReleaseRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &ReleaseRequest{
		Handle: uint64(req.Handle),
		Flags: req.Flags,
		ReleaseFlags: req.ReleaseFlags,
		LockOwner: req.LockOwner,
		Dir: req.Dir,
	}

	n := unsafe.Sizeof(*args)
	body := toReader(unsafe.Pointer(args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/release", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	req.Respond()
}

func handleForgetRequest(ctx Context, host string, req *fuse.ForgetRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &ForgetRequest{
		Inode: uint64(req.Node),
		LookupReqid: uint64(req.N),
	}

	n := unsafe.Sizeof(*args)
	body := toReader(unsafe.Pointer(args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/forget", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	req.Respond()
}

func handleInterruptRequest(ctx Context, host string, req *fuse.InterruptRequest) {

	client := newBoltClient(&req.Header, nil)

	args := &InterruptRequest{
		IntrReqId: uint64(req.IntrID),
	}

	n := unsafe.Sizeof(*args)
	body := toReader(unsafe.Pointer(args), n)
	resp, err := client.DoRequestWith(ctx, "POST", host + "/v1/interrupt", "application/fuse", body, int(n))
	if err != nil {
		replyError(req, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		respondError(req, resp.Body)
		return
	}

	req.Respond()
}

