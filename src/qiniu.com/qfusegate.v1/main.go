package qfusegate

import (
	"encoding/json"
	"io"
	"os"
	"sync"
	"syscall"

	"bazil.org/fuse"
	"github.com/qiniu/errors"
	"github.com/qiniu/http/httputil.v1"
	"github.com/qiniu/log.v1"
)

var (
	ErrInvalidAllowMode = httputil.NewError(
		400, "invalid argument `allow`: value can be `allow_root` or `allow_other`")
)

// ---------------------------------------------------------------------------

type Config struct {
	SaveToFile string `json:"save_to"`
	BackupFile string `json:"backup_to"`
}

type Service struct {
	Config

	mounts []*MountArgs
	conns  map[string]*Conn // mountpoint => conn
	mutex  sync.Mutex
}

func New(cfg *Config) (p *Service, err error) {

	var mounts []*MountArgs

	f, err2 := os.Open(cfg.SaveToFile)
	if err2 == nil {
		json.NewDecoder(f).Decode(&mounts)
		f.Close()
	}

	p = &Service{
		Config: *cfg,
		conns:  make(map[string]*Conn),
	}
	for _, args := range mounts {
		fuse.Unmount(args.MountPoint)
		err = p.mount(args)
		if err != nil {
			return
		}
	}
	return
}

// ---------------------------------------------------------------------------

type MountArgs struct {
	// 加载点
	//
	MountPoint  string `json:"mountpoint"`

	// 目标文件系统位置(Host)。如 "http://127.0.0.1:7777"
	//
	TargeFSHost string `json:"target"`

	// the file system name (also called source) that is visible in the list of mounted file systems
	//
	FSName      string `json:"fsname"`

	// Subtype sets the subtype of the mount. The main type is always `fuse`.
	// The type in a list of mounted file systems will look like `fuse.foo`.
	//
	Subtype     string `json:"subtype"`

	// VolumeName sets the volume name shown in Finder.
	// OS X only. Others ignore this option.
	//
	Name        string `json:"name"`

	// "allow_other" allows other users to access the file system.
	// "allow_root" allows other users to access the file system.
	//
	AllowMode   string `json:"allow"`

	// ReadOnly makes the mount read-only.
	//
	ReadOnly    int    `json:"readonly"`
}

func (p *Service) PostMount(args *MountArgs) (err error) {

	p.mutex.Lock()
	_, ok := p.conns[args.MountPoint]
	if ok {
		err = syscall.EEXIST
	} else {
		err = p.mount(args)
		if err == nil {
			p.mounts = append(p.mounts, args)
			err = p.save()
		}
	}
	p.mutex.Unlock()
	return
}

func (p *Service) save() (err error) {

	os.Remove(p.BackupFile)

	backup, err := os.Create(p.BackupFile)
	if err != nil {
		err = errors.Info(err, "os.Create:", p.BackupFile).Detail(err)
		return
	}
	defer backup.Close()

	old, err2 := os.Open(p.SaveToFile)
	if err2 == nil {
		io.Copy(backup, old)
		old.Close()
	}

	f, err := os.Create(p.SaveToFile)
	if err != nil {
		err = errors.Info(err, "os.Create:", p.SaveToFile).Detail(err)
		return
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(p.mounts)
	if err != nil {
		err = errors.Info(err, "json.Marshal:", p.mounts).Detail(err)
	}
	return
}

func (p *Service) mount(args *MountArgs) (err error) {

	options, err := mountOptions(args)
	if err != nil {
		err = errors.Info(err, "parse mount options failed").Detail(err)
		return
	}

	log.Info("To mount:", args.MountPoint)
	c, err := fuse.Mount(args.MountPoint, options...)
	if err != nil {
		err = errors.Info(err, "fuse.Mount:", args.MountPoint, options).Detail(err)
		return
	}

	log.Info("NewConn:", *args)
	conn, err := NewConn(c, args)
	if err != nil {
		err = errors.Info(err, "qfusegate.NewConn failed").Detail(err)
		return
	}
	p.conns[args.MountPoint] = conn

	go func() {
		err := conn.Serve()
		if err != nil {
			log.Error("Serve failed:", err, "mount:", *args)
		}
	}()
	return
}

func mountOptions(args *MountArgs) (options []fuse.MountOption, err error) {

	if args.FSName != "" {
		options = append(options, fuse.FSName(args.FSName))
	}
	if args.Subtype != "" {
		options = append(options, fuse.Subtype(args.Subtype))
	}
	if args.Name != "" {
		options = append(options, fuse.VolumeName(args.Name))
	}
	switch args.AllowMode {
	case "":
	case "allow_root":
		options = append(options, fuse.AllowRoot())
	case "allow_other":
		options = append(options, fuse.AllowOther())
	default:
		err = ErrInvalidAllowMode
		return
	}
	if args.ReadOnly != 0 {
		options = append(options, fuse.ReadOnly())
	}
	return
}

// ---------------------------------------------------------------------------

type unmountArgs struct {
	MountPoint string `json:"mountpoint"`
}

func (p *Service) PostUnmount(args *unmountArgs) (err error) {

	return fuse.Unmount(args.MountPoint)
}

// ---------------------------------------------------------------------------

