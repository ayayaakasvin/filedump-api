package dto

import "time"

type ShareLinkRequest struct {
	FileUUID string        `json:"file_uuid"`
	Duration time.Duration `json:"ttl"`
}
