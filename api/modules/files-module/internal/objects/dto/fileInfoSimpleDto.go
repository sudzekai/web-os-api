package dto

import "time"

type FileInfoSimpleDto struct {
	Name             string
	Size             int64
	CreationDateTime time.Time
	ChangeDateTime   time.Time
}
