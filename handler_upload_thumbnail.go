package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	video, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Problem with getting video", err)
		return
	}
	if video.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "Video not belonging to you", err)
		return
	}
	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	const maxMemory = 10 << 20

	r.ParseMultipartForm(maxMemory)

	// "thumbnailk" should match the HTML form input name
	file, header, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse form file", err)
		return
	}
	defer file.Close()

	// fmt.Println("content_typ", header.Header.Get("Content-Type"))
	contentType := header.Header.Get("Content-Type")
	// fileFormat := strings.Split(contentType, `/`)[1]
	// fileName := videoIDString + "." + fileFormat
	// fmt.Println("filename", fileName)
	// path := filepath.Join(cfg.assetsRoot, fileName)
	// fmt.Println("path", path)
	// newFile, err := os.Create(path)
	// neviem, err := io.Copy(newFile, file)
	// fmt.Println("neviem", neviem, err)
	// s := "http://localhost:" + cfg.port + "/assets/" + fileName
	mimeType, params, err := mime.ParseMediaType(contentType)
	if mimeType != "image/png" && mimeType != "image/jpg" {
		respondWithError(w, http.StatusInternalServerError, "We need jpg or png", err)
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Mime type error", err)
		return
	}

	fmt.Println("kkkk", params, mimeType)
	assetPath := getAssetPath(contentType)
	assetDiskPath := cfg.getAssetDiskPath(assetPath)

	dst, err := os.Create(assetDiskPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create file on server", err)
		return
	}
	defer dst.Close()
	if _, err = io.Copy(dst, file); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error saving file", err)
		return
	}
	url := cfg.getAssetURL(assetPath)
	video.ThumbnailURL = &url
	fmt.Println(*video.ThumbnailURL)
	if err := cfg.db.UpdateVideo(video); err != nil {
		respondWithError(w, http.StatusBadRequest, "Problem with updating video", err)
		return
	}
	respondWithJSON(w, http.StatusOK, struct{}{})
}
