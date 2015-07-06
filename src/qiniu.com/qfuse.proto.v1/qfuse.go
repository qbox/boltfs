package qfuse

import (
	"bazil.org/fuse"
	"os"
	"time"
)

// ---------------------------------------------------------------------------
// A `init` request is the first request sent on a FUSE file system.

type InitRequest struct {
	Major  uint32
	Minor  uint32
	// Maximum readahead in bytes that the kernel plans to use.
	MaxReadahead uint32
	Flags        fuse.InitFlags
}

type InitResponse struct {
	// Maximum readahead in bytes that the kernel can use. Ignored if
	// greater than InitRequest.MaxReadahead.
	MaxReadahead uint32
	Flags        fuse.InitFlags
	// Maximum size of a single write operation.
	// Linux enforces a minimum of 4 KiB.
	MaxWrite uint32
}

// ---------------------------------------------------------------------------
// A `destroy` request is sent by the kernel when unmounting the file system.
// No more requests will be received after this one, but it should still be responded to.

// ---------------------------------------------------------------------------
// A `statfs` request requests information about the mounted file system.

type StatfsResponse struct {
	Blocks  uint64 // Total data blocks in file system.
	Bfree   uint64 // Free blocks in file system.
	Bavail  uint64 // Free blocks in file system if you're not root.
	Files   uint64 // Total files in file system.
	Ffree   uint64 // Free files in file system.
	Bsize   uint32 // Block size
	Namelen uint32 // Maximum file name length?
	Frsize  uint32 // Fragment size, smallest addressable data size in the file system.
}

// ---------------------------------------------------------------------------
// An `access` request asks whether the file can be accessed for the purpose specified by the mask.

type AccessRequest struct {
	Inode uint64
	Mask  uint32
}

// ---------------------------------------------------------------------------
// A `getattr` request asks for the metadata for the file denoted by Inode.
// An Attr is the metadata for a single file or directory.

type GetattrRequest struct {
	Inode uint64
}

type GetattrResponse struct {
	Attr fuse.Attr
}

// ---------------------------------------------------------------------------
// A `listxattr` request asks to list the extended attributes associated with Inode.

type ListattrRequest struct {
	Inode    uint64 // inode number
	Size     uint32 // maximum size to return
	Position uint32 // offset within attribute list
}

type ListattrResponse struct {
	XattrNames []byte // 以 '\0' 为分割
}

// ---------------------------------------------------------------------------
// A `getxattr` request asks for the extended attributes associated with Inode.

type GetxattrRequest struct {
	// Inode number
	Inode  uint64

	// Name of the attribute requested.
	Name string

	// Maximum size to return.
	Size uint32

	// Offset within extended attributes.
	//
	// Only valid for OS X, and then only with the resource fork
	// attribute.
	Position uint32
}

type GetxattrResponse struct {
	Xattr []byte
}

// ---------------------------------------------------------------------------
// A `removexattr` request asks to remove an extended attribute associated with Inode.

type RemoveattrRequest struct {
	Inode uint64
	Name  string
}

// ---------------------------------------------------------------------------
// A `setxattr` request asks to set an extended attribute associated with a file.

type SetxattrRequest struct {
	Inode uint64

	// Flags can make the request fail if attribute does/not already
	// exist. Unfortunately, the constants are platform-specific and
	// not exposed by Go1.2. Look for XATTR_CREATE, XATTR_REPLACE.
	//
	// TODO improve this later
	//
	// TODO XATTR_CREATE and exist -> EEXIST
	//
	// TODO XATTR_REPLACE and not exist -> ENODATA
	Flags uint32

	// Offset within extended attributes.
	//
	// Only valid for OS X, and then only with the resource fork
	// attribute.
	Position uint32

	Name  string
	Xattr []byte
}

// ---------------------------------------------------------------------------
// A `lookup` request asks to look up the given name in the directory named by Inode.

type LookupRequest struct {
	Inode uint64
	Name  string
}

type LookupResponse struct {
	Inode      uint64
	Generation uint64
	EntryValid time.Duration
	Attr       fuse.Attr
}

// ---------------------------------------------------------------------------
// An `open` request asks to open a file or directory.

type OpenRequest struct {
	Inode  uint64
	Flags  fuse.OpenFlags
	Dir    bool // is this Opendir?
}

type OpenResponse struct {
	Handle uint64
	Flags  fuse.OpenResponseFlags
}

// ---------------------------------------------------------------------------
// A `create` request asks to create and open a file (not a directory).

type CreateRequest struct {
	Inode  uint64
	Name   string
	Flags  fuse.OpenFlags
	Mode   os.FileMode
}

type CreateResponse struct {
	LookupResponse
	OpenResponse
}

// ---------------------------------------------------------------------------
// A `mkdir` request asks to create (but not open) a directory.

type MkdirRequest struct {
	Inode  uint64
	Name   string
	Mode   os.FileMode
}

type MkdirResponse LookupResponse

// ---------------------------------------------------------------------------
// A `symlink` request is a request to create a symlink making NewName point to Target.

type SymlinkRequest struct {
	Inode   uint64
	NewName string
	Target  string
}

type SymlinkResponse LookupResponse

// ---------------------------------------------------------------------------
// A `readlink` request is a request to read a symlink's target.

type ReadlinkRequest struct {
	Inode  uint64
}

type ReadlinkResponse struct {
	Target string
}

// ---------------------------------------------------------------------------
// A `link` request is a request to create a hard link.

type LinkRequest struct {
	Inode    uint64
	NewName  string
	OldInode uint64
}

type LinkResponse LookupResponse

// ---------------------------------------------------------------------------
// 创建设备（/v1/mknod）

type MknodRequest struct {
	Inode  uint64
	Name   string
	Mode   os.FileMode
	Rdev   uint32
}

type MknodResponse LookupResponse

// ---------------------------------------------------------------------------
// A `rename` request is a request to rename a file.

type RenameRequest struct {
	NewDirInode uint64
	OldName     string
	NewName     string
}

// ---------------------------------------------------------------------------
// A `remove` request asks to remove a file or directory from the directory Inode.

type RemoveRequest struct {
	Inode  uint64
	Name   string // name of the entry to remove
	Dir    bool   // is this rmdir?
}

// ---------------------------------------------------------------------------
// A `read` request asks to read from an open file/dir.

type ReadRequest struct {
	Handle uint64
	Offset int64
	Size   int
	Dir    bool // is this Readdir?
}

type ReadResponse struct {
	Data []byte
}

// ---------------------------------------------------------------------------
// A `write` request asks to write to an open file.

type WriteRequest struct {
	Handle uint64
	Offset int64
	Data   []byte
	Flags  fuse.WriteFlags
}

type WriteResponse struct {
	Size int
}

// ---------------------------------------------------------------------------
// A `setattr` request asks to change one or more attributes associated with a file, as indicated by Valid.

type SetattrRequest struct {
	Valid  fuse.SetattrValid
	Handle uint64
	Size   uint64
	Atime  time.Time
	Mtime  time.Time
	Mode   os.FileMode
	Uid    uint32
	Gid    uint32

	// OS X only
	Bkuptime time.Time
	Chgtime  time.Time
	Crtime   time.Time
	Flags    uint32 // see chflags(2)
}

type SetattrResponse struct {
	Attr fuse.Attr
}

// ---------------------------------------------------------------------------
// A `flush` request asks for the current state of an open file to be flushed to storage.
// A single opened Handle may receive multiple `flush` requests over its lifetime.

type FlushRequest struct {
	Handle    uint64
	LockOwner uint64
	Flags     uint32
}

// ---------------------------------------------------------------------------
// 刷新Inode（/v/fsync）

type FsyncRequest struct {
	Handle uint64
	Flags  uint32 // TODO bit 1 is datasync, not well documented upstream
	Dir    bool
}

// ---------------------------------------------------------------------------
// A `release` request asks to release (close) an open file/dir handle.


type ReleaseRequest struct {
	Handle       uint64
	Flags        fuse.OpenFlags // flags from OpenRequest
	ReleaseFlags fuse.ReleaseFlags
	LockOwner    uint32
	Dir          bool // is this Releasedir?
}

// ---------------------------------------------------------------------------
// A `forget` request is sent by the kernel when forgetting about Inode as returned by `lookup` requests.

type ForgetRequest struct {
	Inode       uint64
	LookupReqid uint64
}

// ---------------------------------------------------------------------------
// An `interrupt` request is a request to interrupt another pending request.
// The response to that request should return an error status of EINTR.

type InterruptRequest struct {
	IntrReqId uint64
}

// ---------------------------------------------------------------------------

