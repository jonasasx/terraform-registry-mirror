package hash

import (
	"bufio"
	"golang.org/x/mod/sumdb/dirhash"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func GetHashes(downloadUrl string, shaUrl string) ([]string, error) {
	hashes := []string{}
	hash, err := downloadAndHash(downloadUrl)
	if err != nil {
		return nil, err
	}
	hashes = append(hashes, hash)
	shaUrlHashes, err := getShaUrlHashes(shaUrl)
	if err != nil {
		return nil, err
	}
	hashes = slices.Concat(hashes, shaUrlHashes)
	return hashes, nil
}

func downloadAndHash(url string) (string, error) {
	filename := filepath.Base(url)
	out, err := os.CreateTemp("", filename)
	if err != nil {
		return "", err
	}
	defer out.Close()
	resp, err := http.Get(url)
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	hash, err := dirhash.HashZip(out.Name(), dirhash.Hash1)
	return hash, nil
}

func getShaUrlHashes(url string) ([]string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var hashes []string
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()
		hash, _, _ := strings.Cut(line, " ")
		hashes = append(hashes, "zh:"+hash)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return hashes, nil
}
