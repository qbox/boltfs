QBOLT 网络协议
==============

# 约定

## 请求通用头

```
Authorization: QBolt base64(<Uid/Gid/Pid:uint32>)
X-Reqid: base36(<Reqid:uint64>)
```

## 出错返回包

```
Errno int32
```

# 协议

## 初始化(/v1/init)

* A `init` request is the first request sent on a FUSE file system.

请求体：

```
Major  uint32
Minor  uint32
// Maximum readahead in bytes that the kernel plans to use.
MaxReadahead uint32
Flags        InitFlags
```

返回体：

```
// Maximum readahead in bytes that the kernel can use. Ignored if
// greater than InitRequest.MaxReadahead.
MaxReadahead uint32
Flags        InitFlags
// Maximum size of a single write operation.
// Linux enforces a minimum of 4 KiB.
MaxWrite uint32
```

其中 InitFlags：

```
type InitFlags uint32

const (
	InitAsyncRead       InitFlags = 1 << 0
	InitPosixLocks      InitFlags = 1 << 1
	InitFileOps         InitFlags = 1 << 2
	InitAtomicTrunc     InitFlags = 1 << 3
	InitExportSupport   InitFlags = 1 << 4
	InitBigWrites       InitFlags = 1 << 5
	InitDontMask        InitFlags = 1 << 6
	InitSpliceWrite     InitFlags = 1 << 7
	InitSpliceMove      InitFlags = 1 << 8
	InitSpliceRead      InitFlags = 1 << 9
	InitFlockLocks      InitFlags = 1 << 10
	InitHasIoctlDir     InitFlags = 1 << 11
	InitAutoInvalData   InitFlags = 1 << 12
	InitDoReaddirplus   InitFlags = 1 << 13
	InitReaddirplusAuto InitFlags = 1 << 14
	InitAsyncDIO        InitFlags = 1 << 15
	InitWritebackCache  InitFlags = 1 << 16
	InitNoOpenSupport   InitFlags = 1 << 17

	InitCaseSensitive InitFlags = 1 << 29 // OS X only
	InitVolRename     InitFlags = 1 << 30 // OS X only
	InitXtimes        InitFlags = 1 << 31 // OS X only
)
```

## 终止（/v1/destroy）

* A `destroy` request is sent by the kernel when unmounting the file system.
* No more requests will be received after this one, but it should still be responded to.

请求体：无

返回体：无

## 文件系统信息（/v1/statfs）

* A `statfs` request requests information about the mounted file system.

请求体：无

返回体：

```
Blocks  uint64 // Total data blocks in file system.
Bfree   uint64 // Free blocks in file system.
Bavail  uint64 // Free blocks in file system if you're not root.
Files   uint64 // Total files in file system.
Ffree   uint64 // Free files in file system.
Bsize   uint32 // Block size
Namelen uint32 // Maximum file name length?
Frsize  uint32 // Fragment size, smallest addressable data size in the file system.
```

## 可访问性（/v1/access）

* An `access` request asks whether the file can be accessed for the purpose specified by the mask.

请求体：

```
Inode uint64
Mask  uint32
```

返回体：

* 如果可访问则无需返回任何信息。如果不能访问，返回错误包。

## 取属性（/v1/getattr）

* A `getattr` request asks for the metadata for the file denoted by Inode.
* An Attr is the metadata for a single file or directory.

请求体：

```
Inode uint64
```

返回体：

```
Attr Attr
```

其中 Attr：

```
type Attr struct {
	Valid time.Duration // how long Attr can be cached

	Inode  uint64      // inode number
	Size   uint64      // size in bytes
	Blocks uint64      // size in blocks
	Atime  time.Time   // time of last access
	Mtime  time.Time   // time of last modification
	Ctime  time.Time   // time of last inode change
	Crtime time.Time   // time of creation (OS X only)
	Mode   os.FileMode // file mode
	Nlink  uint32      // number of links
	Uid    uint32      // owner uid
	Gid    uint32      // group gid
	Rdev   uint32      // device numbers
	Flags  uint32      // chflags(2) flags (OS X only)
}
```

## 取扩展属性列表（/v1/listxattr)

* A `listxattr` request asks to list the extended attributes associated with Inode.

请求体：

```
Inode    uint64 // inode number
Size     uint32 // maximum size to return
Position uint32 // offset within attribute list
```

返回体：

```
XattrNames []byte // 以 '\0' 为分割
```

## 取某扩展属性（/v1/getxattr）

* A `getxattr` request asks for the extended attributes associated with Inode.

请求体：

```
// Inode number
Inode  uint64

// Maximum size to return.
Size uint32

// Name of the attribute requested.
Name string

// Offset within extended attributes.
//
// Only valid for OS X, and then only with the resource fork
// attribute.
Position uint32
```

返回体：

```
Xattr []byte
```

## 删除扩展属性（/v1/removexattr）

* A `removexattr` request asks to remove an extended attribute associated with Inode.

请求体：

```
Inode uint64
Name  string
```

返回体：无

## 设置扩展属性（/v1/setxattr）

* A `setxattr` request asks to set an extended attribute associated with a file.

请求体：

```
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
```

返回体：无

## 查找节点（/v1/lookup）

* A `lookup` request asks to look up the given name in the directory named by Inode.

请求体：

```
Inode uint64
Name  string
```

返回体：

```
Inode      uint64
Generation uint64
EntryValid time.Duration
Attr       Attr
```

## 打开文件/目录（/v1/open）

* An `open` request asks to open a file or directory.

请求体：

```
Inode  uint64
Dir    bool // is this Opendir?
Flags  OpenFlags
```

返回体：

```
Handle uint64
Flags  OpenResponseFlags
```

其中 OpenFlags：

```
// OpenFlags are the O_FOO flags passed to open/create/etc calls. For
// example, os.O_WRONLY | os.O_APPEND.
type OpenFlags uint32

const (
	// Access modes. These are not 1-bit flags, but alternatives where
	// only one can be chosen. See the IsReadOnly etc convenience
	// methods.
	OpenReadOnly  OpenFlags = syscall.O_RDONLY
	OpenWriteOnly OpenFlags = syscall.O_WRONLY
	OpenReadWrite OpenFlags = syscall.O_RDWR

	OpenAppend    OpenFlags = syscall.O_APPEND
	OpenCreate    OpenFlags = syscall.O_CREAT
	OpenExclusive OpenFlags = syscall.O_EXCL
	OpenSync      OpenFlags = syscall.O_SYNC
	OpenTruncate  OpenFlags = syscall.O_TRUNC
)

// OpenAccessModeMask is a bitmask that separates the access mode
// from the other flags in OpenFlags.
const OpenAccessModeMask OpenFlags = syscall.O_ACCMODE
```

其中 OpenResponseFlags：

```
// The OpenResponseFlags are returned in the OpenResponse.
type OpenResponseFlags uint32

const (
	OpenDirectIO    OpenResponseFlags = 1 << 0 // bypass page cache for this open file
	OpenKeepCache   OpenResponseFlags = 1 << 1 // don't invalidate the data cache on open
	OpenNonSeekable OpenResponseFlags = 1 << 2 // (Linux?)

	OpenPurgeAttr OpenResponseFlags = 1 << 30 // OS X
	OpenPurgeUBC  OpenResponseFlags = 1 << 31 // OS X
)
```

## 创建文件（/v1/create）

* A `create` request asks to create and open a file (not a directory).

请求体：

```
Inode  uint64
Name   string
Flags  OpenFlags
Mode   os.FileMode
```

返回体：

```
// lookup response
//
Inode      uint64
Generation uint64
EntryValid time.Duration
Attr       Attr

// open response
//
Handle uint64
Flags  OpenResponseFlags
```

## 创建目录（/v1/mkdir）

* A `mkdir` request asks to create (but not open) a directory.

请求体：

```
Inode  uint64
Name   string
Mode   os.FileMode
```

返回体：

```
// lookup response
//
Inode      uint64
Generation uint64
EntryValid time.Duration
Attr       Attr
```

## 创建软链接（/v1/symlink）

* A `symlink` request is a request to create a symlink making NewName point to Target.

请求体：

```
Inode   uint64
NewName string
Target  string
```

返回体：

```
// lookup response
//
Inode      uint64
Generation uint64
EntryValid time.Duration
Attr       Attr
```

## 读软链接（/v1/readlink）

* A `readlink` request is a request to read a symlink's target.

请求体：

```
Inode  uint64
```

返回体：

```
Target string
```

## 创建硬链接（/v1/link）

* A `link` request is a request to create a hard link.

请求体：

```
Inode    uint64
NewName  string
OldInode uint64
```

返回体：

```
// lookup response
//
Inode      uint64
Generation uint64
EntryValid time.Duration
Attr       Attr
```

## 创建设备（/v1/mknod）

请求体：

```
Inode  uint64
Name   string
Mode   os.FileMode
Rdev   uint32
```

返回包：

```
// lookup response
//
Inode      uint64
Generation uint64
EntryValid time.Duration
Attr       Attr
```

## 移动文件（/v1/rename）

* A `rename` request is a request to rename a file.

请求体：

```
NewDirInode uint64
OldName     string
NewName     string
```

返回体：无

## 删除文件/目录（/v1/remove）

* A `remove` request asks to remove a file or directory from the directory Inode.

请求体：

```
Inode  uint64
Name   string // name of the entry to remove
Dir    bool   // is this rmdir?
```

返回体：无

## 读数据（/v1/read）

* A `read` request asks to read from an open file/dir.

请求体：

```
Handle uint64
Offset int64
Size   int
Dir    bool // is this Readdir?
```

返回体：

```
Data []byte
```

## 写数据（/v1/write）

* A `write` request asks to write to an open file.

请求体：

```
Handle uint64
Offset int64
Data   []byte
Flags  WriteFlags
```

返回体：

```
Size int
```

其中 WriteFlags：

```
// The WriteFlags are passed in WriteRequest.
type WriteFlags uint32
```

## 设置属性（/v1/setattr）

* A `setattr` request asks to change one or more attributes associated with a file, as indicated by Valid.

请求体：

```
Valid  SetattrValid
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
```

返回体：

```
Attr Attr
```

其中 SetattrValid：

```
// The SetattrValid are bit flags describing which fields in the SetattrRequest
// are included in the change.
type SetattrValid uint32

const (
	SetattrMode   SetattrValid = 1 << 0
	SetattrUid    SetattrValid = 1 << 1
	SetattrGid    SetattrValid = 1 << 2
	SetattrSize   SetattrValid = 1 << 3
	SetattrAtime  SetattrValid = 1 << 4
	SetattrMtime  SetattrValid = 1 << 5
	SetattrHandle SetattrValid = 1 << 6

	// Linux only(?)
	SetattrAtimeNow  SetattrValid = 1 << 7
	SetattrMtimeNow  SetattrValid = 1 << 8
	SetattrLockOwner SetattrValid = 1 << 9 // http://www.mail-archive.com/git-commits-head@vger.kernel.org/msg27852.html

	// OS X only
	SetattrCrtime   SetattrValid = 1 << 28
	SetattrChgtime  SetattrValid = 1 << 29
	SetattrBkuptime SetattrValid = 1 << 30
	SetattrFlags    SetattrValid = 1 << 31
)
```

## 刷新文件缓存（/v1/flush）

* A `flush` request asks for the current state of an open file to be flushed to storage.
* A single opened Handle may receive multiple `flush` requests over its lifetime.

请求体：

```
Handle    uint64
Flags     uint32
LockOwner uint64
```

返回体：无

## 刷新Inode（/v/fsync）

请求体：

```
Handle uint64
Flags  uint32 // TODO bit 1 is datasync, not well documented upstream
Dir    bool
```

返回体：无

## 关闭文件/目录（/v1/release）

* A `release` request asks to release (close) an open file/dir handle.

请求体：

```
Handle       uint64
Flags        OpenFlags // flags from OpenRequest
ReleaseFlags ReleaseFlags
LockOwner    uint32
Dir          bool // is this Releasedir?
```

返回体：无

其中 ReleaseFlags：

```
// The ReleaseFlags are used in the Release exchange.
type ReleaseFlags uint32

const (
	ReleaseFlush ReleaseFlags = 1 << 0
)
```

## 忘记Inode（/v1/forget）

* A `forget` request is sent by the kernel when forgetting about Inode as returned by `lookup` requests.

请求体：

```
Inode       uint64
LookupReqid uint64
```

返回体：无

## 打断请求（/v1/interrupt）

* An `interrupt` request is a request to interrupt another pending request.
* The response to that request should return an error status of EINTR.

请求体：

```
IntrReqId uint64
```

返回体：无

