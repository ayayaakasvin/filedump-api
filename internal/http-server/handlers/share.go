package handlers

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"up-down-server/internal/http-server/ctx"
	"up-down-server/internal/lib/bindjson"
	"up-down-server/internal/lib/linkgeneration"
	"up-down-server/internal/models"
	"up-down-server/internal/models/dto"
	"up-down-server/internal/repository/postgresql"
)

const shareLinkKey = "share link:%s"

// Example: /download/shared/<hash>
const prefix = "/download/shared/"

func (h *Handlers) CreateShareLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var slr dto.ShareLinkRequest
		if err := bindjson.BindJson(r.Body, &slr); err != nil {
			models.SendErrorJson(w, http.StatusBadRequest, "failed to bind request %s", err.Error())
			return
		}

		reqUserID := r.Context().Value(ctx.CtxUserIDKey).(int)

		filemeta, err := h.fileRepo.GetFileMeta(r.Context(), slr.FileUUID)
		if err != nil {
			if err.Error() != postgresql.NotFound {
				h.logger.WithError(err).Error("failed to get filemeta data")
			}
			models.SendErrorJson(w, http.StatusInternalServerError, "failed to get filemeta data")
			return
		}

		if filemeta.UserID != reqUserID {
			models.SendErrorJson(w, http.StatusUnauthorized, "access denied")
			return
		}

		link, err := linkgeneration.GenerateRandomLink(slr.FileUUID, slr.Duration)
		if err != nil {
			switch err.Error() {
			case linkgeneration.InvalidInput, linkgeneration.WayTooLong:
				models.SendErrorJson(w, http.StatusBadRequest, "%s", err.Error())
			default:
				h.logger.WithError(err).Error("internal error during link generation")
				models.SendErrorJson(w, http.StatusInternalServerError, "internal error")
			}

			return
		}

		key := fmt.Sprintf(shareLinkKey, link)

		set := h.cache.SetNX(r.Context(), key, slr.FileUUID, slr.Duration)
		if set.Err() != nil {
			h.logger.WithError(set.Err()).Error("cache error during share link creation")
			models.SendErrorJson(w, http.StatusInternalServerError, "cache error")
			return
		}
		if !set.Val() {
			models.SendErrorJson(w, http.StatusConflict, "share link already exists")
			return
		}

		data := models.NewData()
		data["link"] = link

		models.SendSuccessJson(w, http.StatusOK, data)
	}
}

func (h *Handlers) DownloadFileViaSharedLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, prefix) || len(r.URL.Path) <= len(prefix) {
			http.NotFound(w, r)
			return
		}
		hash := strings.TrimPrefix(r.URL.Path, prefix)
		h.logger.Info("link", hash)

		var fileuuid string
		if fileuuidAny, err := h.cache.Get(r.Context(), fmt.Sprintf(shareLinkKey, hash)); err != nil {
			models.SendErrorJson(w, http.StatusNotFound, "no such link available")
			return
		} else {
			fileuuid = fileuuidAny.(string)
		}

		if fileuuid == "" {
			models.SendErrorJson(w, http.StatusInternalServerError, "internal error")
			return
		}

		filemeta, err := h.fileRepo.GetFileMeta(r.Context(), fileuuid)
		if err != nil {
			switch err.Error() {
			case postgresql.NotFound:
				models.SendErrorJson(w, http.StatusNotFound, "no such file")
			default:
				models.SendErrorJson(w, http.StatusInternalServerError, "failed to get filemeta data")
			}

			return
		}

		if _, err := os.Stat(filepath.Join("./", filemeta.FilePath)); os.IsNotExist(err) {
			models.SendErrorJson(w, http.StatusNotFound, "file not found")
			return
		} else if err != nil {
			h.logger.Errorf("Failed to stat file: %v", err)
			models.SendErrorJson(w, http.StatusInternalServerError, "server error")
			return
		}

		mimeType := mime.TypeByExtension(filepath.Ext(filemeta.FileName))
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		w.Header().Set("Content-Type", mimeType)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filemeta.FileName))
		w.Header().Set("Content-Transfer-Encoding", "binary") // Optional, helps in some clients
		w.Header().Set("Cache-Control", "no-cache")           // Optional

		http.ServeFile(w, r, filemeta.FilePath)
	}
}
