package types

type FileType string

const (
	Regular         FileType = "file"
	Directory       FileType = "directory"
	Symlink         FileType = "symlink"
	Socket          FileType = "socket"
	Pipe            FileType = "pipe"
	BlockDevice     FileType = "block device"
	CharacterDevice FileType = "character device"
	Unknown         FileType = "unknown"
)
