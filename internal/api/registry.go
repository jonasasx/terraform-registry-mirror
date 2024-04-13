package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetVersions(hostname string, namespace string, pkg string) ([]Version, error) {
	requestURL := fmt.Sprintf("https://%s/v1/providers/%s/%s/versions", hostname, namespace, pkg)
	res, err := http.Get(requestURL)
	var response VersionsResponse
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return response.Versions, nil
}

func GetPackage(hostname string, namespace string, pkg string, version string, os string, arch string) (*DownloadResponse, error) {
	requestURL := fmt.Sprintf("https://%s/v1/providers/%s/%s/%s/download/%s/%s", hostname, namespace, pkg, version, os, arch)
	res, err := http.Get(requestURL)
	var response DownloadResponse
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

type VersionsResponse struct {
	Id       string      `json:"id"`
	Versions []Version   `json:"versions"`
	Warnings interface{} `json:"warnings"`
}

type Version struct {
	Version   string     `json:"version"`
	Protocols []string   `json:"protocols"`
	Platforms []Platform `json:"platforms"`
}

type Platform struct {
	Os   string `json:"os"`
	Arch string `json:"arch"`
}

type DownloadResponse struct {
	Protocols           []string `json:"protocols"`
	Os                  string   `json:"os"`
	Arch                string   `json:"arch"`
	Filename            string   `json:"filename"`
	DownloadUrl         string   `json:"download_url"`
	ShasumsUrl          string   `json:"shasums_url"`
	ShasumsSignatureUrl string   `json:"shasums_signature_url"`
	Shasum              string   `json:"shasum"`
	SigningKeys         struct {
		GpgPublicKeys []struct {
			KeyId          string      `json:"key_id"`
			AsciiArmor     string      `json:"ascii_armor"`
			TrustSignature string      `json:"trust_signature"`
			Source         string      `json:"source"`
			SourceUrl      interface{} `json:"source_url"`
		} `json:"gpg_public_keys"`
	} `json:"signing_keys"`
}
