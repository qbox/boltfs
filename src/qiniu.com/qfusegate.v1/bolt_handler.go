package qfusegate

import (
	"bazil.org/fuse"
	. "golang.org/x/net/context"
)

func handleInitRequest(ctx Context, target string, r *fuse.InitRequest) {
}

func handleDestroyRequest(ctx Context, target string, r *fuse.DestroyRequest) {
}

func handleStatfsRequest(ctx Context, target string, r *fuse.StatfsRequest) {
}

func handleAccessRequest(ctx Context, target string, r *fuse.AccessRequest) {
}

func handleGetattrRequest(ctx Context, target string, r *fuse.GetattrRequest) {
}

func handleListxattrRequest(ctx Context, target string, r *fuse.ListxattrRequest) {
}

func handleGetxattrRequest(ctx Context, target string, r *fuse.GetxattrRequest) {
}

func handleRemovexattrRequest(ctx Context, target string, r *fuse.RemovexattrRequest) {
}

func handleSetxattrRequest(ctx Context, target string, r *fuse.SetxattrRequest) {
}

func handleLookupRequest(ctx Context, target string, r *fuse.LookupRequest) {
}

func handleOpenRequest(ctx Context, target string, r *fuse.OpenRequest) {
}

func handleCreateRequest(ctx Context, target string, r *fuse.CreateRequest) {
}

func handleMkdirRequest(ctx Context, target string, r *fuse.MkdirRequest) {
}

func handleSymlinkRequest(ctx Context, target string, r *fuse.SymlinkRequest) {
}

func handleReadlinkRequest(ctx Context, target string, r *fuse.ReadlinkRequest) {
}

func handleLinkRequest(ctx Context, target string, r *fuse.LinkRequest) {
}

func handleMknodRequest(ctx Context, target string, r *fuse.MknodRequest) {
}

func handleRenameRequest(ctx Context, target string, r *fuse.RenameRequest) {
}

func handleRemoveRequest(ctx Context, target string, r *fuse.RemoveRequest) {
}

func handleReadRequest(ctx Context, target string, r *fuse.ReadRequest) {
}

func handleWriteRequest(ctx Context, target string, r *fuse.WriteRequest) {
}

func handleSetattrRequest(ctx Context, target string, r *fuse.SetattrRequest) {
}

func handleFlushRequest(ctx Context, target string, r *fuse.FlushRequest) {
}

func handleFsyncRequest(ctx Context, target string, r *fuse.FsyncRequest) {
}

func handleReleaseRequest(ctx Context, target string, r *fuse.ReleaseRequest) {
}

func handleForgetRequest(ctx Context, target string, r *fuse.ForgetRequest) {
}

func handleInterruptRequest(ctx Context, target string, r *fuse.InterruptRequest) {
}

