package data

import (
	"io/fs"
	"ls/utils"
	"strings"
	"time"
)

const (
	blue    = "\033[1;34m" // Bright Blue for directories
	blue2   = "\033[34m"   // Blue
	green   = "\033[1;32m" // Bright Green for executable files
	cyan    = "\033[36m"   // Cyan for audio files
	cyan2   = "\033[1;36m" // Bright Cyan for symbolic links
	red     = "\033[1;31m" // Bright Red for broken symbolic links
	magenta = "\033[1;35m" // Bright Magenta for image files
	yellow  = "\033[1;33m" // Bright Yellow
	black   = "\033[30m"   // Black
	reset   = "\033[0m"    // Resets to the default terminal color

	// Background colors
	bgBlack  = "\033[40m"
	bgRed    = "\033[41m"
	bgYellow = "\033[43m"
	bgBlue   = "\033[44m"
	bgGreen  = "\033[42m"
)

type MyLSFiles struct {
	Name            string
	IsDir           bool
	IsExec          bool
	IsLink          bool
	LinkTarget      string
	TargetFile      *MyLSFiles
	IsBroken        bool
	IsBlockDevice   bool
	IsCharDevice    bool
	MajorNumber     uint32
	MinorNumber     uint32
	IsSocket        bool
	IsPipe          bool
	IsSetuid        bool
	IsSetgid        bool
	IsStickyDir     bool
	IsOtherWritable bool
	Size            int64
	ModTime         time.Time
	Mode            fs.FileMode
	OwnerName       string
	GroupName       string
	NLink           uint64
}

func (file MyLSFiles) GetColor() string {
	if file.IsBroken {
		return bgBlack + red
	}
	if file.IsStickyDir && file.IsOtherWritable {
		return bgGreen + black
	}
	if file.IsOtherWritable {
		return bgGreen + blue2
	}
	if file.IsStickyDir {
		return bgBlue
	}
	if file.IsSetuid {
		return bgRed
	}
	if file.IsSetgid {
		return bgYellow + black
	}
	if file.IsLink {
		return cyan2
	}
	if file.IsDir {
		return blue
	}
	if file.IsExec {
		return green
	}
	if file.IsBlockDevice || file.IsCharDevice || file.IsPipe {
		return bgBlack + yellow
	}
	if file.IsSocket {
		return magenta
	}

	ext := strings.ToLower(utils.Ext(file.Name))

	switch ext {
	case // Image files
		".jpg", ".jpeg", ".gif", ".bmp", ".pbm", ".pgm", ".ppm", ".tga",
		".xbm", ".xpm", ".tif", ".tiff", ".png", ".svg", ".svgz", ".mng",
		".pcx",
		// Video files
		".mov", ".mpg", ".mpeg", ".m2v", ".mkv", ".webm", ".ogm", ".mp4", ".m4v",
		".mp4v", ".vob", ".qt", ".nuv", ".wmv", ".asf", ".rm", ".rmvb",
		".flc", ".avi", ".fli", ".flv", ".gl", ".dl", ".xcf", ".xwd",
		".yuv", ".cgm", ".emf", ".axv", ".anx", ".ogv", ".ogx":
		return magenta
	case // audio files
		".aac", ".au", ".flac", ".mid", ".midi", ".mka", ".mp3", ".mpc",
		".ogg", ".ra", ".wav", ".axa", ".oga", ".spx", ".xspf":
		return cyan
	case // compressed files
		".tar", ".tgz", ".arj", ".taz", ".lzh", ".lzma", ".tlz", ".txz",
		".zip", ".z", ".Z", ".dz", ".gz", ".lz", ".xz", ".bz2", ".bz", ".tbz", ".tbz2",
		".tz", ".deb", ".rpm", ".jar", ".rar", ".ace", ".zoo", ".cpio",
		".7z", ".rz", ".cab", ".war", ".ear", ".sar":
		return red
	}
	return reset
}
