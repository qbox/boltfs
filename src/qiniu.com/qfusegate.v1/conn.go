package qfusegate

import (
	"io"
	"sync"

	"bazil.org/fuse"
	"golang.org/x/net/context"
)

// ---------------------------------------------------------------------------

type Conn struct {
	c    *fuse.Conn
	args *MountArgs
}

func NewConn(c *fuse.Conn, args *MountArgs) (p *Conn, err error) {

	p = &Conn{c, args}
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
		_, _ = r, ctx
	case *fuse.WriteRequest:
	case *fuse.FlushRequest:
	case *fuse.FsyncRequest:
	case *fuse.ReleaseRequest:

	// Node operations.
	case *fuse.AccessRequest:
	case *fuse.GetattrRequest:
	case *fuse.SetattrRequest:
	case *fuse.SymlinkRequest:
	case *fuse.ReadlinkRequest:
	case *fuse.LinkRequest:
	case *fuse.RemoveRequest:
	case *fuse.LookupRequest:
	case *fuse.MkdirRequest:
	case *fuse.OpenRequest:
	case *fuse.CreateRequest:
	case *fuse.GetxattrRequest:
	case *fuse.ListxattrRequest:
	case *fuse.SetxattrRequest:
	case *fuse.RemovexattrRequest:
	case *fuse.RenameRequest:
	case *fuse.MknodRequest:
	case *fuse.ForgetRequest:

	// FS operations.
	case *fuse.InterruptRequest:
	case *fuse.InitRequest:
	case *fuse.StatfsRequest:
	case *fuse.DestroyRequest:
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

// ---------------------------------------------------------------------------

