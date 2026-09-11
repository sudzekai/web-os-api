package types

type FilePermission string

const (
	Read    FilePermission = "read"
	Write   FilePermission = "write"
	Execute FilePermission = "execute"

	Setuid FilePermission = "setuid"
	Setgid FilePermission = "setgid"
	Sticky FilePermission = "sticky"
)
