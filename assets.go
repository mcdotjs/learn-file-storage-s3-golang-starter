package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type D interface {
	Write(p []byte) (n int, err error)
}

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

// maybe refactor with this functions
func getAssetPath(mediaType string) string {
	ext := mediaTypeToExt(mediaType)
	key := make([]byte, 32)
	rand.Read(key)

	dest := bytes.NewBuffer(make([]byte, 0))
	encoder := base64.NewEncoder(base64.RawURLEncoding, dest)
	encoder.Write(key)
	encoder.Close()

	//return dest.String() + ext
	return fmt.Sprintf("%s%s", dest.String(), ext)
	//NOTE: ORRRRRRR
	// base := make([]byte, 32)
	// 	_, err := rand.Read(base)
	// 	if err != nil {
	// 		panic("failed to generate random bytes")
	// 	}
	// 	id := base64.RawURLEncoding.EncodeToString(base)
	// 	ext := mediaTypeToExt(mediaType)
	// 	return fmt.Sprintf("%s%s", id, ext)
}

func (cfg apiConfig) getAssetDiskPath(assetPath string) string {
	return filepath.Join(cfg.assetsRoot, assetPath)
}

func (cfg apiConfig) getAssetURL(assetPath string) string {
	return fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, assetPath)
}

func mediaTypeToExt(mediaType string) string {
	parts := strings.Split(mediaType, "/")
	if len(parts) != 2 {
		return ".bin"
	}
	return "." + parts[1]
}
