package logic

import (
	"ls/data"
	"ls/utils"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

const maxSymlinkDepth = 10

func GetFileAttributes(path string, info os.FileInfo, isDirectArgument bool, depth int) data.MyLSFiles {
	stat, _ := info.Sys().(*syscall.Stat_t)
	var nlink uint64 = 1
	var uid, gid uint32

	file := data.MyLSFiles{}

	if stat != nil {
		nlink = uint64(stat.Nlink)
		uid = stat.Uid
		gid = stat.Gid
	}

	ownerName := strconv.Itoa(int(uid))
	if owner, err := user.LookupId(strconv.Itoa(int(uid))); err == nil {
		ownerName = owner.Username
	}

	groupName := strconv.Itoa(int(gid))
	if group, err := user.LookupGroupId(strconv.Itoa(int(gid))); err == nil {
		groupName = group.Name
	}

	var targetFile *data.MyLSFiles
	var targetPath string
	var err error
	if info.Mode()&os.ModeSymlink != 0 {
		targetPath, err = os.Readlink(path)
		if err != nil {
			file.IsBroken = true
		} else {
			absTarget := utils.Join(utils.Dir(path), targetPath)

			targetInfo, err := os.Lstat(absTarget)
			if err != nil {
				file.IsBroken = true
			}

			tf := GetFileAttributes(absTarget, targetInfo, false, 0)
			targetFile = &tf

			if depth < maxSymlinkDepth {
				finalInfo, err := os.Stat(absTarget)
				if err != nil {
					file.IsBroken = true
				} else {
					ft := GetFileAttributes(absTarget, finalInfo, false, depth+1)
					file.FinalTarget = &ft
				}
			}
		}
	}

	var major, minor uint32
	if stat != nil {
		major = uint32((stat.Rdev >> 8) & 0xFF) // Linux/Unix specific
		minor = uint32(stat.Rdev & 0xFF)
	}

	return data.MyLSFiles{
		Name:            GetDisplayName(path, isDirectArgument),
		IsDir:           info.IsDir(),
		IsExec:          !info.IsDir() && (info.Mode().Perm()&0o111 != 0),
		IsLink:          info.Mode()&os.ModeSymlink != 0,
		LinkTarget:      targetPath,
		TargetFile:      targetFile,
		FinalTarget:     file.FinalTarget,
		IsBroken:        info.Mode()&os.ModeSymlink != 0 && !Exists(path),
		IsBlockDevice:   info.Mode()&os.ModeType == os.ModeDevice,
		IsCharDevice:    info.Mode()&os.ModeType == (os.ModeDevice | os.ModeCharDevice),
		MajorNumber:     major,
		MinorNumber:     minor,
		IsSocket:        info.Mode()&os.ModeSocket != 0,
		IsPipe:          info.Mode()&os.ModeNamedPipe != 0,
		IsSetuid:        info.Mode()&os.ModeSetuid != 0,
		IsSetgid:        info.Mode()&os.ModeSetgid != 0,
		IsStickyDir:     info.IsDir() && info.Mode()&os.ModeSticky != 0,
		IsOtherWritable: info.IsDir() && info.Mode()&0o002 != 0,
		Size:            info.Size(),
		ModTime:         info.ModTime(),
		Mode:            info.Mode(),
		OwnerName:       ownerName,
		GroupName:       groupName,
		NLink:           nlink,
	}
}

func GetDisplayName(path string, isDirectArgument bool) string {
	if isDirectArgument {
		return path // Preserve original path for direct arguments
	}
	return utils.Base(path) // Use base name for directory contents
}

// exists checks whether a file or directory exists at the given path.
// Returns true if the file exists, otherwise returns false.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
