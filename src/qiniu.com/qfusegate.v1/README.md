qfusegate
============

# 协议

## 挂载

请求包：

```
POST /v1/mount
Content-Type: application/json

{
	# 加载点
	#
	"mountpoint": <MountPoint>,

	# 目标文件系统位置(Host)。如 "http://127.0.0.1:7777"
	#
	"target": <TargeFSHost>,

	# the file system name (also called source) that is visible in the list of mounted file systems
	#
	"fsname": <FSName>,

	# Subtype sets the subtype of the mount. The main type is always `fuse`.
	# The type in a list of mounted file systems will look like `fuse.foo`.
	#
	"subtype": <Subtype>,

	# VolumeName sets the volume name shown in Finder.
	# OS X only. Others ignore this option.
	#
	"name": <VolumeName>,

	# "allow_other" allows other users to access the file system.
	# "allow_root" allows other users to access the file system.
	#
	"allow": <AllowMode>,

	# ReadOnly makes the mount read-only.
	#
	"readonly": <ReadOnly>
}
```

返回包：

```
200 OK
```

## 取消挂载

请求包：

```
POST /v1/unmount
Content-Type: application/json

{
	"mountpoint": <MountPoint>
}
```

返回包：

```
200 OK
```

