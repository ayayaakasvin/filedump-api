package handlers

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"up-down-server/internal/http-server/ctx"
	"up-down-server/internal/lib/bindjson"
	"up-down-server/internal/lib/validinput"
	"up-down-server/internal/models"
	"up-down-server/internal/models/dto"
	"up-down-server/internal/repository/postgresql"

	"github.com/google/uuid"
)

const (
	XWWWFormApplication = "application/x-www-form-urlencoded"
	downloadDir         = "files"
	maxSizeForFile      = 10 << 20
)

// File upload handler with maxSizeForFile of ~100mb, better to request with Form holding file param with file itself
func (h *Handlers) UploadFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(maxSizeForFile)
		if err != nil {
			h.logger.Errorf("Parse error: %v", err)
			models.SendErrorJson(w, http.StatusInternalServerError, "Failed to parse MultipartForm")
			return
		}

		f, handler, err := r.FormFile("file")
		if err != nil {
			h.logger.Errorf("FormFile error: %v", err)
			models.SendErrorJson(w, http.StatusInternalServerError, "Failed to parse file from Form")
			return
		}
		defer f.Close()

		uuidOfNewFIle := uuid.New().String()
		dir := filepath.Join(".", "files")
		fullPath := filepath.Join(dir, uuidOfNewFIle)
		userId := r.Context().Value(ctx.CtxUserIDKey).(int)

		metadata := models.NewFileMetaData(uuidOfNewFIle, handler, fullPath, userId)
		if err := h.fileRepo.InsertFileName(r.Context(), metadata); err != nil {
			h.logger.Errorf("InsertFileName error: %v", err)
			models.SendErrorJson(w, http.StatusInternalServerError, "Failed to insert recors")
			return
		}

		dst, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			h.logger.Errorf("OpenFile error: %v", err)
			models.SendErrorJson(w, http.StatusInternalServerError, "Failed to save file")
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, f); err != nil {
			h.logger.Errorf("Copy error: %v", err)
			models.SendErrorJson(w, http.StatusInternalServerError, "Failed to write file")
			return
		}

		h.logger.Infof("Uploaded file saved as: %s\n", fullPath)
		resp := models.NewData()
		resp["file_id"] = uuidOfNewFIle
		models.SendSuccessJson(w, http.StatusCreated, resp)
	}
}

// File Download, with single param as "file_id" as uuid of file(string)
// File served via ServeFile with necessary headers
func (h *Handlers) DownloadFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file_id := r.URL.Query().Get("file_id")
		if file_id == "" {
			models.SendErrorJson(w, http.StatusBadRequest, "file_id is required")
			return
		}

		reqUserID := r.Context().Value(ctx.CtxUserIDKey).(int)

		fileMeta, err := h.fileRepo.GetFileMeta(r.Context(), file_id)
		if err != nil {
			switch err.Error() {
			case postgresql.NotFound:
				models.SendErrorJson(w, http.StatusNotFound, "file not found")
			default:
				h.logger.Errorf("Failed to fetch file metadata: %v", err)
				models.SendErrorJson(w, http.StatusInternalServerError, "failed to retrieve metadata")
			}
			return
		} else if fileMeta.UserID != reqUserID {
			models.SendErrorJson(w, http.StatusUnauthorized, "access denied")
		}

		if _, err := os.Stat(filepath.Join("./", fileMeta.FilePath)); os.IsNotExist(err) {
			models.SendErrorJson(w, http.StatusNotFound, "file not found")
			return
		} else if err != nil {
			h.logger.Errorf("Failed to stat file: %v", err)
			models.SendErrorJson(w, http.StatusInternalServerError, "server error")
			return
		}

		mimeType := mime.TypeByExtension(filepath.Ext(fileMeta.FileName))
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		w.Header().Set("Content-Type", mimeType)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileMeta.FileName))
		w.Header().Set("Content-Transfer-Encoding", "binary") // Optional, helps in some clients
		w.Header().Set("Cache-Control", "no-cache")           // Optional

		http.ServeFile(w, r, fileMeta.FilePath)
	}
}

// File metadata list as JSON array of FileMetaData
func (h *Handlers) ListFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIdInt := r.Context().Value(ctx.CtxUserIDKey).(int)

		records, err := h.fileRepo.GetUserRecords(r.Context(), userIdInt)
		if err != nil {
			h.logger.Errorf("Failed to retrieve records: %v", err)
			models.SendErrorJson(w, http.StatusInternalServerError, "Record retrieve error")
			return
		}

		if len(records) == 0 {
			models.SendErrorJson(w, http.StatusNotFound, "No record retrieved")
			return
		}

		data := models.NewData()
		data["records"] = records

		models.SendSuccessJson(w, http.StatusOK, data)
	}
}

// File delete handler, deletion via "file_id" param and further delete from db and etc
func (h *Handlers) DeleteFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqUserID := r.Context().Value(ctx.CtxUserIDKey).(int)

		fileuuid := r.URL.Query().Get("file_id")
		if fileuuid == "" {
			models.SendErrorJson(w, http.StatusBadRequest, "file_id is required")
			return
		}

		if filemeta, err := h.fileRepo.GetFileMeta(r.Context(), fileuuid); err != nil {
			switch err.Error() {
			case postgresql.NotFound:
				models.SendErrorJson(w, http.StatusNotFound, "file not found")
			default:
				h.logger.Errorf("Failed to fetch file metadata: %v", err)
				models.SendErrorJson(w, http.StatusInternalServerError, "failed to retrieve metadata")
			}
			return
		} else if filemeta.UserID != reqUserID {
			models.SendErrorJson(w, http.StatusUnauthorized, "access denied")
			return 
		}

		if err := h.fileRepo.DeleteFileByUUID(r.Context(), fileuuid); err != nil {
			h.logger.Errorf("Failed to delete record: %v", err)
			models.SendErrorJson(w, http.StatusInternalServerError, "Failed to delete record")
			return
		}

		pathToFile := path.Join(downloadDir, fileuuid)
		if err := os.Remove(pathToFile); err != nil {
			h.logger.Errorf("Failed to delete file: %v", err)
			models.SendErrorJson(w, http.StatusInternalServerError, "Failed to delete file")
			return
		}

		models.SendSuccessJson(w, http.StatusOK, nil)
	}
}

// Just like Listing, but with solo data as FileMetaData JSON
func (h *Handlers) GetFileMetaData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqUserID := r.Context().Value(ctx.CtxUserIDKey).(int)

		fileuuid := r.URL.Query().Get("file_id")
		if fileuuid == "" {
			models.SendErrorJson(w, http.StatusBadRequest, "file_id is required")
			return
		}

		// call to db
		filemeta, err := h.fileRepo.GetFileMeta(r.Context(), fileuuid)
		if err != nil {
			switch err.Error() {
			case postgresql.NotFound:
				models.SendErrorJson(w, http.StatusNotFound, "file not found")
			default:
				h.logger.Errorf("Failed to fetch file metadata: %v", err)
				models.SendErrorJson(w, http.StatusInternalServerError, "failed to retrieve metadata")
			}
			return
		} else if filemeta.UserID != reqUserID {
			models.SendErrorJson(w, http.StatusUnauthorized, "access denied")
		}

		data := models.NewData()
		data["metadata"] = filemeta
		models.SendSuccessJson(w, http.StatusOK, data)
	}
}

// Filename update handler, changing filename
func (h *Handlers) UpdateFileName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqUserID := r.Context().Value(ctx.CtxUserIDKey).(int)

		fileuuid := r.URL.Query().Get("file_id")
		if fileuuid == "" {
			models.SendErrorJson(w, http.StatusBadRequest, "file_id is required")
			return
		}

		if filemeta, err := h.fileRepo.GetFileMeta(r.Context(), fileuuid); err != nil {
			switch err.Error() {
			case postgresql.NotFound:
				models.SendErrorJson(w, http.StatusNotFound, "file not found")
			default:
				h.logger.Errorf("Failed to fetch file metadata: %v", err)
				models.SendErrorJson(w, http.StatusInternalServerError, "failed to retrieve metadata")
			}
			return
		} else if filemeta.UserID != reqUserID {
			models.SendErrorJson(w, http.StatusUnauthorized, "access denied")
			return 
		}

		var updateReq dto.UpdateFileNameRequest
		if err := bindjson.BindJson(r.Body, &updateReq); err != nil {
			models.SendErrorJson(w, http.StatusBadGateway, "failed to bind request")
			return
		}

		if !validinput.IsValidFileName(updateReq.Filename) {
			models.SendErrorJson(w, http.StatusBadRequest, "invalid filename request body")
			return
		}

		if err := h.fileRepo.RenameFileName(r.Context(), updateReq.Filename, fileuuid); err != nil {
			switch err.Error() {
			case postgresql.NotFound:
				models.SendErrorJson(w, http.StatusNotFound, "file not found")
			case postgresql.UnAuthorized:
				models.SendErrorJson(w, http.StatusUnauthorized, "access denied")
			default:
				h.logger.Errorf("Failed to update filename: %v", err)
				models.SendErrorJson(w, http.StatusInternalServerError, "failed to update filename")
			}
			return
		}

		models.SendSuccessJson(w, http.StatusOK, nil)
	}
}