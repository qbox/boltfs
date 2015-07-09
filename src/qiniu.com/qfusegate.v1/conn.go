package qfusegate

import (
	"bytes"
	"io"
	"strings"
	"sync"
	"unsafe"

	"bazil.org/fuse"
	"golang.org/x/net/context"
	"qiniupkg.com/x/rpc.v7"
)

// ---------------------------------------------------------------------------

type Conn struct {
	target   string
	c        *fuse.Conn
	readOnly bool
}

func NewConn(c *fuse.Conn, args *MountArgs) (p *Conn, err error) {

	p = &Conn{
		c:        c,
		target:   args.TargetFSHost,
		readOnly: args.ReadOnly != 0,
	}
	return
}

func (p *Conn) Serve() (err error) {

	var wg sync.WaitGroup
	defer wg.Wait()

	for {
		req, err := p.c.ReadRequest()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			p.serveRequest(req)
		}()
	}
	return nil
}

func (p *Conn) serveRequest(r fuse.Request) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	switch r := r.(type) {
	// Handle operations.
	case *fuse.ReadRequest:
		handleReadRequest(ctx, p.target, r)
	case *fuse.WriteRequest:
		handleWriteRequest(ctx, p.target, r)
	case *fuse.FlushRequest:
		handleFlushRequest(ctx, p.target, r)
	case *fuse.FsyncRequest:
		handleFsyncRequest(ctx, p.target, r)
	case *fuse.ReleaseRequest:
		handleReleaseRequest(ctx, p.target, r)

	// Node operations.
	case *fuse.AccessRequest:
		handleAccessRequest(ctx, p.target, r)
	case *fuse.GetattrRequest:
		handleGetattrRequest(ctx, p.target, r)
	case *fuse.SetattrRequest:
		handleSetattrRequest(ctx, p.target, r)
	case *fuse.SymlinkRequest:
		handleSymlinkRequest(ctx, p.target, r)
	case *fuse.ReadlinkRequest:
		handleReadlinkRequest(ctx, p.target, r)
	case *fuse.LinkRequest:
		handleLinkRequest(ctx, p.target, r)
	case *fuse.RemoveRequest:
		handleRemoveRequest(ctx, p.target, r)
	case *fuse.LookupRequest:
		handleLookupRequest(ctx, p.target, r)
	case *fuse.MkdirRequest:
		handleMkdirRequest(ctx, p.target, r)
	case *fuse.OpenRequest:
		handleOpenRequest(ctx, p.target, r)
	case *fuse.CreateRequest:
		handleCreateRequest(ctx, p.target, r)
	case *fuse.GetxattrRequest:
		handleGetxattrRequest(ctx, p.target, r)
	case *fuse.ListxattrRequest:
		handleListxattrRequest(ctx, p.target, r)
	case *fuse.SetxattrRequest:
		handleSetxattrRequest(ctx, p.target, r)
	case *fuse.RemovexattrRequest:
		handleRemovexattrRequest(ctx, p.target, r)
	case *fuse.RenameRequest:
		handleRenameRequest(ctx, p.target, r)
	case *fuse.MknodRequest:
		handleMknodRequest(ctx, p.target, r)
	case *fuse.ForgetRequest:
		handleForgetRequest(ctx, p.target, r)

	// FS operations.
	case *fuse.InterruptRequest:
		handleInterruptRequest(ctx, p.target, r)
	case *fuse.InitRequest:
		handleInitRequest(ctx, p.target, r)
	case *fuse.StatfsRequest:
		handleStatfsRequest(ctx, p.target, r)
	case *fuse.DestroyRequest:
		handleDestroyRequest(ctx, p.target, r)
	}

	// Note: To FUSE, ENOSYS means "this server never implements this request."
	// It would be inappropriate to return ENOSYS for other operations in this
	// switch that might only be unavailable in some contexts, not all.
	/*
	case *FsyncdirRequest:
		done(ENOSYS)
		r.RespondError(ENOSYS)

	case *GetlkRequest, *SetlkRequest, *SetlkwRequest:
		done(ENOSYS)
		r.RespondError(ENOSYS)

	case *BmapRequest:
		done(ENOSYS)
		r.RespondError(ENOSYS)

	case *SetvolnameRequest, *GetxtimesRequest, *ExchangeRequest:
		done(ENOSYS)
		r.RespondError(ENOSYS)
	*/
	r.RespondError(fuse.ENOSYS)
}

func replyError(r fuse.Request, err error) {

	if e, ok := err.(*rpc.ErrorInfo); ok && e.Errno != 0 {
		r.RespondError(fuse.Errno(e.Errno))
	} else {
		r.RespondError(err)
	}
}

func toReader(p unsafe.Pointer, n uintptr) (r io.Reader) {

	b := ((*[1<<30]byte)(p))[:n]
	return bytes.NewReader(b)
}

// ---------------------------------------------------------------------------

var (
	stringReplacer = strings.NewReplacer("\\", "\\\\", "\n", "\\n")
)

type stringEncoder struct {
	bytes.Buffer
}

func (p *stringEncoder) PutString(data string) {

	stringReplacer.WriteString(&p.Buffer, data)
	p.WriteByte('\n')
}

func (p *stringEncoder) PutBytes(data []byte) {

	stringReplacer.WriteString(&p.Buffer, string(data))
	p.WriteByte('\n')
}

// ---------------------------------------------------------------------------

