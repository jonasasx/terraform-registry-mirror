package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/mod/sumdb/dirhash"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func GetHashes(downloadUrl string) ([]string, error) {
	hashes, err := downloadAndHash(downloadUrl)
	if err != nil {
		return nil, err
	}
	return hashes, nil
}

func downloadAndHash(url string) ([]string, error) {
	filename := filepath.Base(url)
	out, err := os.CreateTemp("", filename)
	if err != nil {
		return nil, err
	}
	defer out.Close()
	resp, err := http.Get(url)
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, err
	}
	hashes := []string{}
	hash, err := dirhash.HashZip(out.Name(), dirhash.Hash1)
	hashes = append(hashes, hash)
	f, err := os.Open(out.Name())
	if err != nil {
		return nil, err
	}
	defer f.Close()
	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return nil, err
	}
	hash = hex.EncodeToString(hasher.Sum(nil))
	hashes = append(hashes, "zh:"+hash)
	return hashes, nil
}
