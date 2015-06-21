package qfusegate

import (
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

// ---------------------------------------------------------------------------

type FS struct {

}

func NewFS() (p *FS, err error) {

	return &FS{}, nil
}

func (p *FS) Root() (fs.Node, error) {

	return &Dir{}, nil
}

// ---------------------------------------------------------------------------

type Dir struct {

}

func (p *Dir) Attr(ctx context.Context, a *fuse.Attr) error {

	a.Inode = 1
	a.Mode = os.ModeDir | 0555
	return nil
}

func (p *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {

	if name == "hello" {
		return &File{}, nil
	}
	return nil, fuse.ENOENT
}

// ---------------------------------------------------------------------------

type File struct {

}

func (p *File) Attr(ctx context.Context, a *fuse.Attr) error {

	a.Inode = 2
	a.Mode = 0444
	a.Size = uint64(len(greeting))
	return nil
}

func (p *File) ReadAll(ctx context.Context) ([]byte, error) {

	return []byte(greeting), nil
}

const greeting = "hello, world!"

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

	pfs, _ := NewFS()
	fs.Serve(p.c, pfs)
	return
}

// ---------------------------------------------------------------------------

