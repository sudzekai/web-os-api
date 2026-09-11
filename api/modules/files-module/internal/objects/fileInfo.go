package objects

import (
	"strconv"
	"strings"
	"time"

	"github.com/sudzekai/web-os-api/modules/files-module/internal/objects/types"
)

type FileInfo struct {
	Name     string
	FullName string

	Type types.FileType

	Size      int64
	Inode     uint64
	Links     uint64
	DeviceID  uint64
	RDevice   uint64
	BlockSize int64
	Blocks    int64

	Owner   string
	OwnerID uint32

	Group   string
	GroupID uint32

	LinkTarget string

	SpecialPermissions []types.FilePermission
	UserPermissions    []types.FilePermission
	GroupPermissions   []types.FilePermission
	OthersPermissions  []types.FilePermission

	BirthDateTime        time.Time
	ModificationDateTime time.Time
	ChangeDateTime       time.Time
	AccessDateTime       time.Time

	stats map[string]string
}

func NewFileInfo(stat string) *FileInfo {
	fileInfo := FileInfo{
		stats: make(map[string]string),
	}

	fileInfo.mapStat(stat)
	fileInfo.construct()

	return &fileInfo
}

func (fi *FileInfo) mapStat(stat string) {
	for entry := range strings.SplitSeq(stat, "\n") {
		if idx := strings.Index(entry, ":"); idx != -1 {
			key := strings.TrimSpace(entry[:idx])
			value := strings.TrimSpace(entry[idx+1:])
			fi.stats[key] = value
		}
	}
}

func (fi *FileInfo) construct() {
	fi.Name = fi.getName()
	fi.FullName = fi.getFullName()

	fi.Type = fi.getType()

	fi.Size = fi.getSize()
	fi.Inode = fi.getInode()
	fi.Links = fi.getLinks()
	fi.DeviceID = fi.getDeviceId()
	fi.RDevice = fi.getRDevice()
	fi.BlockSize = fi.getBlockSize()
	fi.Blocks = fi.getBlocks()

	fi.Owner = fi.getOwner()
	fi.OwnerID = fi.getOwnerID()

	fi.Group = fi.getGroup()
	fi.GroupID = fi.getGroupID()

	fi.LinkTarget = fi.getLinkTarget()

	permissions := parsePermissions(fi.getPermissions())

	fi.SpecialPermissions = fi.getSpecialPermissions(permissions)
	fi.UserPermissions = fi.getUserPermissions(permissions)
	fi.GroupPermissions = fi.getGroupPermissions(permissions)
	fi.OthersPermissions = fi.getOthersPermissions(permissions)

	fi.BirthDateTime = fi.getBirthDateTime()
	fi.ModificationDateTime = fi.getModificationDateTime()
	fi.ChangeDateTime = fi.getChangeDateTime()
	fi.AccessDateTime = fi.getAccessDateTime()
}

func (fi *FileInfo) getName() string {
	return fi.stats["Name"]
}

func (fi *FileInfo) getFullName() string {
	return fi.stats["FullName"]
}

func (fi *FileInfo) getType() types.FileType {
	typeStr := fi.stats["Type"]

	switch typeStr {
	case "regular file":
		return types.Regular
	case "directory":
		return types.Directory
	case "symbolic link":
		return types.Symlink
	case "socket":
		return types.Socket
	case "fifo":
		return types.Pipe
	case "block special file":
		return types.BlockDevice
	case "character special file":
		return types.CharacterDevice
	default:
		return types.Unknown
	}
}

func (fi *FileInfo) getSize() int64 {
	sizeStr := fi.stats[types.Size.Key]
	size, _ := strconv.ParseInt(sizeStr, 10, 64)
	return size
}

func (fi *FileInfo) getInode() uint64 {
	inodeStr := fi.stats[types.Inode.Key]
	inode, _ := strconv.ParseUint(inodeStr, 10, 64)
	return inode
}

func (fi *FileInfo) getLinks() uint64 {
	linksStr := fi.stats[types.Links.Key]
	links, _ := strconv.ParseUint(linksStr, 10, 64)
	return links
}

func (fi *FileInfo) getDeviceId() uint64 {
	deviceIDStr := fi.stats[types.DeviceID.Key]
	deviceID, _ := strconv.ParseUint(deviceIDStr, 10, 64)
	return deviceID
}

func (fi *FileInfo) getRDevice() uint64 {
	rDeviceStr := fi.stats[types.RDevice.Key]
	rDevice, _ := strconv.ParseUint(rDeviceStr, 10, 64)
	return rDevice
}

func (fi *FileInfo) getBlockSize() int64 {
	blockSizeStr := fi.stats[types.BlockSize.Key]
	blockSize, _ := strconv.ParseInt(blockSizeStr, 10, 64)
	return blockSize
}

func (fi *FileInfo) getBlocks() int64 {
	blocksStr := fi.stats[types.Blocks.Key]
	blocks, _ := strconv.ParseInt(blocksStr, 10, 64)
	return blocks
}

func (fi *FileInfo) getOwner() string {
	return fi.stats[types.Owner.Key]
}

func (fi *FileInfo) getOwnerID() uint32 {
	ownerIDStr := fi.stats[types.OwnerID.Key]
	ownerID, _ := strconv.ParseUint(ownerIDStr, 10, 32)
	return uint32(ownerID)
}

func (fi *FileInfo) getGroup() string {
	return fi.stats[types.Group.Key]
}

func (fi *FileInfo) getGroupID() uint32 {
	groupIDStr := fi.stats[types.GroupID.Key]
	groupID, _ := strconv.ParseUint(groupIDStr, 10, 32)
	return uint32(groupID)
}

func (fi *FileInfo) getLinkTarget() string {
	return fi.stats[types.LinkTarget.Key]
}

func (fi *FileInfo) getSpecialPermissions(permissions map[string][]types.FilePermission) []types.FilePermission {
	return permissions[types.SpecialPermissions.Key]
}

func (fi *FileInfo) getUserPermissions(permissions map[string][]types.FilePermission) []types.FilePermission {
	return permissions[types.UserPermissions.Key]
}

func (fi *FileInfo) getGroupPermissions(permissions map[string][]types.FilePermission) []types.FilePermission {
	return permissions[types.GroupPermissions.Key]
}

func (fi *FileInfo) getOthersPermissions(permissions map[string][]types.FilePermission) []types.FilePermission {
	return permissions[types.OthersPermissions.Key]
}

func (fi *FileInfo) getBirthDateTime() time.Time {
	return parseUnixString(fi.stats[types.BirthDateTime.Key])
}

func (fi *FileInfo) getModificationDateTime() time.Time {
	return parseUnixString(fi.stats[types.ModificationDateTime.Key])
}

func (fi *FileInfo) getChangeDateTime() time.Time {
	return parseUnixString(fi.stats[types.ChangeDateTime.Key])
}

func (fi *FileInfo) getAccessDateTime() time.Time {
	return parseUnixString(fi.stats[types.AccessDateTime.Key])
}

func (fi *FileInfo) getPermissions() string {
	return fi.stats[types.Permissions.Key]
}

func parsePermissions(permissionsStr string) map[string][]types.FilePermission {
	chars := []rune(permissionsStr)

	if len(chars) < 4 {
		chars = append([]rune{'0'}, chars...)
	}

	result := make(map[string][]types.FilePermission)

	for i := 0; i < 4; i++ {
		perms := make([]types.FilePermission, 0, 3)

		if chars[i] == '0' {
			continue
		}

		perm := int(chars[i] - '0')

		if i != 0 {
			if perm&4 != 0 {
				perms = append(perms, types.Read)
			}
			if perm&2 != 0 {
				perms = append(perms, types.Write)
			}
			if perm&1 != 0 {
				perms = append(perms, types.Execute)
			}
		} else {
			if perm&4 != 0 {
				perms = append(perms, types.Setuid)
			}

			if perm&2 != 0 {
				perms = append(perms, types.Setgid)
			}

			if perm&1 != 0 {
				perms = append(perms, types.Sticky)
			}
		}

		keys := []string{
			types.SpecialPermissions.Key,
			types.UserPermissions.Key,
			types.GroupPermissions.Key,
			types.OthersPermissions.Key,
		}

		result[keys[i]] = perms
	}

	return result
}

func parseUnixString(unixStr string) time.Time {
	timestamp, _ := strconv.ParseInt(unixStr, 10, 64)

	t := time.Unix(timestamp, 0)
	return t
}
