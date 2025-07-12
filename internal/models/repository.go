package models

import "context"

type FileMetaRepository interface {
	InsertFileName		(ctx context.Context, file *FileMetaData) 								error
	DeleteFileByUUID	(ctx context.Context, uuidOfFile string, userId int) 					error
	GetFileMeta			(ctx context.Context, uuidOfFile string, userId int) 					(*FileMetaData, error)
	GetUUID				(ctx context.Context, filename string) 									(string, error)
	GetAllRecords		(ctx context.Context) 													([]*FileMetaData, error)
	GetUserRecords		(ctx context.Context, userId int) 										([]*FileMetaData, error)
	RenameFileName		(ctx context.Context, updatedFilename, uuidOfFile string, userId int)	error
}

type UserRepository interface {
	RegisterUser		(ctx context.Context, username, hashedpassword string) 	error
	AuthentificateUser	(ctx context.Context, username, password string) 		(int, error)
}
