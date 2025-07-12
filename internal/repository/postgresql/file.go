package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"

	"up-down-server/internal/models"
)

const (
	NotFound     = "not found"
	UnAuthorized = "unauthorized"
)

func (p *PostgreSQL) InsertFileName(ctx context.Context, file *models.FileMetaData) error {
	stmt, err := p.conn.PrepareContext(ctx, `INSERT INTO files (file_uuid, filename, filepath, size, mime_type, user_id) VALUES ($1, $2, $3, $4, $5, $6)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, file.FileUUID, file.FileName, file.FilePath, file.Size, file.MimeType, file.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgreSQL) DeleteFileByUUID(ctx context.Context, uuidOfFile string, userId int) error {
	var ownerID int
	err := p.conn.QueryRowContext(ctx,
		`SELECT user_id FROM files WHERE file_uuid = $1`, uuidOfFile).Scan(&ownerID)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New(NotFound) // 404
		}
		return err
	}

	if ownerID != userId {
		return errors.New(UnAuthorized) // 401
	}

	stmt, err := p.conn.PrepareContext(ctx, `DELETE FROM files WHERE file_uuid = $1 AND user_id = $2`)
	if err != nil {
		return err // 500
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, uuidOfFile, userId)
	if err != nil {
		return err // 500
	}

	return nil
}

func (p *PostgreSQL) GetFileMeta(ctx context.Context, uuidOfFile string, userId int) (*models.FileMetaData, error) {
	stmt, err := p.conn.PrepareContext(ctx, `SELECT file_uuid, filename, filepath, uploaded_at, size, mime_type, user_id FROM files WHERE file_uuid = $1 LIMIT 1`)
	if err != nil {
		return nil, err // 500
	}
	defer stmt.Close()

	var metadata *models.FileMetaData = &models.FileMetaData{}
	err = stmt.QueryRowContext(ctx, uuidOfFile).Scan(
		&metadata.FileUUID,
		&metadata.FileName,
		&metadata.FilePath,
		&metadata.UploadedAt,
		&metadata.Size,
		&metadata.MimeType,
		&metadata.UserID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(NotFound) // 404
		}
		return nil, err // 500
	}

	if metadata.UserID != userId {
		return nil, errors.New(UnAuthorized) // 401
	}

	return metadata, nil
}

func (p *PostgreSQL) GetUUID(ctx context.Context, filename string) (string, error) {
	stmt, err := p.conn.PrepareContext(ctx, `SELECT * FROM files WHERE file_uuid = $1`)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var uuidOfFile string
	err = stmt.QueryRowContext(ctx, filename).Scan(&uuidOfFile)
	if err != nil {
		return "", err
	}

	return uuidOfFile, nil
}

func (p *PostgreSQL) GetAllRecords(ctx context.Context) ([]*models.FileMetaData, error) {
	stmt, err := p.conn.PrepareContext(ctx, `SELECT file_uuid, filename, filepath, uploaded_at, size, mime_type, user_id FROM files`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var result []*models.FileMetaData
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		fmd := new(models.FileMetaData)
		if err := rows.Scan(&fmd.FileUUID, &fmd.FileName, &fmd.FilePath, &fmd.UploadedAt, &fmd.Size, &fmd.MimeType, &fmd.UserID); err != nil {
			return nil, err
		}

		fmd.FileExt = filepath.Ext(fmd.FileName)

		result = append(result, fmd)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *PostgreSQL) GetUserRecords(ctx context.Context, userId int) ([]*models.FileMetaData, error) {
	stmt, err := p.conn.PrepareContext(ctx, `SELECT file_uuid, filename, filepath, uploaded_at, size, mime_type, user_id FROM files WHERE user_id = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var result []*models.FileMetaData
	rows, err := stmt.QueryContext(ctx, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		fmd := new(models.FileMetaData)
		if err := rows.Scan(&fmd.FileUUID, &fmd.FileName, &fmd.FilePath, &fmd.UploadedAt, &fmd.Size, &fmd.MimeType, &fmd.UserID); err != nil {
			return nil, err
		}

		fmd.FileExt = filepath.Ext(fmd.FileName)

		result = append(result, fmd)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *PostgreSQL) RenameFileName(ctx context.Context, updatedFilename, uuidOfFile string, userId int) error {
	var ownerID int
	err := p.conn.QueryRowContext(ctx,
		`SELECT user_id FROM files WHERE file_uuid = $1`, uuidOfFile).Scan(&ownerID)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New(NotFound) // 404
		}
		return err // 500
	}

	if ownerID != userId {
		return errors.New(UnAuthorized) // 401
	}

	stmt, err := p.conn.PrepareContext(ctx, `UPDATE files SET filename = $1 WHERE file_uuid = $2 AND user_id = $3`)
	if err != nil {
		return err // 500
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, updatedFilename, uuidOfFile, userId)
	if err != nil {
		return err // 500
	}

	return nil
}
