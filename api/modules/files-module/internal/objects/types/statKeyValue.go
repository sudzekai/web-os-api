package types

type StatKeyPattern struct {
	Key     string
	Pattern string
}

var (
	Name       = StatKeyPattern{"Name", "%n"}
	FullName   = StatKeyPattern{"FullName", "%N"}
	Type       = StatKeyPattern{"Type", "%F"}
	Size       = StatKeyPattern{"Size", "%s"}
	Inode      = StatKeyPattern{"Inode", "%i"}
	Links      = StatKeyPattern{"Links", "%h"}
	DeviceID   = StatKeyPattern{"DeviceID", "%d"}
	RDevice    = StatKeyPattern{"RDevice", "%r"}
	BlockSize  = StatKeyPattern{"BlockSize", "%B"}
	Blocks     = StatKeyPattern{"Blocks", "%b"}
	Owner      = StatKeyPattern{"Owner", "%U"}
	OwnerID    = StatKeyPattern{"OwnerID", "%u"}
	Group      = StatKeyPattern{"Group", "%G"}
	GroupID    = StatKeyPattern{"GroupID", "%g"}
	LinkTarget = StatKeyPattern{"LinkTarget", "%N"}

	Permissions = StatKeyPattern{"Permissions", "%a"}

	SpecialPermissions = StatKeyPattern{"SpecialPermissions", "%a"}
	UserPermissions    = StatKeyPattern{"UserPermissions", "%a"}
	GroupPermissions   = StatKeyPattern{"GroupPermissions", "%a"}
	OthersPermissions  = StatKeyPattern{"OthersPermissions", "%a"}

	BirthDateTime        = StatKeyPattern{"BirthDateTime", "%w"}
	ModificationDateTime = StatKeyPattern{"ModificationDateTime", "%y"}
	ChangeDateTime       = StatKeyPattern{"ChangeDateTime", "%z"}
	AccessDateTime       = StatKeyPattern{"AccessDateTime", "%x"}
)
