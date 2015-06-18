package qfusegate

import (
	"bazil.org/fuse"
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

	return
}

// ---------------------------------------------------------------------------

