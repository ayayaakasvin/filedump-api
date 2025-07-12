package models

import (
	"time"
	"mime/multipart"
)

type FileMetaData struct {
	FileUUID   string    `json:"file_uuid"`
	FileName   string    `json:"file_name"`
	FileExt    string    `json:"file_ext,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
	Size       int64     `json:"size"`                // should be in bytes
	FilePath   string    `json:"file_path"`           // absolute path to the file
	MimeType   string    `json:"mime_type,omitempty"` // MIME type of the file
	UserID     int       `json:"user_id"`             // ID of the user who uploaded the file
}

func NewFileMetaData(file_uuid string, file *multipart.FileHeader, filepath string, userID int) *FileMetaData {
	return &FileMetaData{
		FileUUID:   file_uuid,
		FileName:   file.Filename,
		FileExt:    parseExt(file.Filename),
		UploadedAt: time.Now(),
		Size:       file.Size,
		FilePath:   filepath,
		MimeType:   file.Header.Get("Content-Type"),
		UserID:     userID,
	}
}

func parseExt(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i:]
		}
		if filename[i] == '/' || filename[i] == '\\' {
			break
		}
	}
	return ""
}
