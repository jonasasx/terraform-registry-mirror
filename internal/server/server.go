package server

import (
	"fmt"
	"net/http"
	"strings"
	"terraform-registry-mirror/internal/api"
	"terraform-registry-mirror/internal/hash"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/store/go_cache/v4"
	"github.com/gin-gonic/gin"
)

type Server interface {
	Index(c *gin.Context)
	Version(c *gin.Context)
}

type server struct {
	versionsCache         *cache.Cache[[]api.Version]
	downloadResponseCache *cache.Cache[*api.DownloadResponse]
	hashesCache           *cache.Cache[[]string]
	logger                logger
}

func NewServer(cacheStore *go_cache.GoCacheStore) Server {
	return server{
		versionsCache:         cache.New[[]api.Version](cacheStore),
		downloadResponseCache: cache.New[*api.DownloadResponse](cacheStore),
		hashesCache:           cache.New[[]string](cacheStore),
		logger:                logger{writer: gin.DefaultWriter},
	}
}

func (s server) Index(c *gin.Context) {
	hostname := c.Param("hostname")
	namespace := c.Param("namespace")
	pkg := c.Param("pkg")
	pkgPath := fmt.Sprintf("%s/%s/%s", hostname, namespace, pkg)

	cacheKey := fmt.Sprintf("versions:%s:%s:%s", hostname, namespace, pkg)

	versions, err := s.versionsCache.Get(c, cacheKey)
	if err != nil {
		s.logger.warn(err.Error())
	}

	if versions == nil {
		versions, err = api.GetVersions(hostname, namespace, pkg)
		if err != nil {
			s.logger.error(err.Error())
			c.AsciiJSON(http.StatusInternalServerError, ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("unable to fetch versions for %s", pkgPath),
			})
			return
		}

		err := s.versionsCache.Set(c, cacheKey, versions)
		if err != nil {
			s.logger.warn(err.Error())
		}
	}
	response := make(map[string]struct{})
	for _, version := range versions {
		response[version.Version] = struct{}{}
	}

	c.AsciiJSON(http.StatusOK, IndexResponse{Versions: response})
}

func (s server) Version(c *gin.Context) {
	hostname := c.Param("hostname")
	namespace := c.Param("namespace")
	pkg := c.Param("pkg")
	pkgPath := fmt.Sprintf("%s/%s/%s", hostname, namespace, pkg)

	version := strings.Replace(c.Param("version.json"), ".json", "", 1)
	versions, err := api.GetVersions(hostname, namespace, pkg)
	if err != nil {
		s.logger.error(err.Error())
		c.AsciiJSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("unable to fetch versions for %s", pkgPath),
		})
		return
	}

	response := VersionResponse{Archives: map[string]Build{}}
	for _, item := range versions {
		if item.Version == version {
			for _, platform := range item.Platforms {
				cacheKey := fmt.Sprintf("downloadResponse:%s:%s:%s:%s:%s:%s", hostname, namespace, pkg, version, platform.Os, platform.Arch)
				packageItem, _ := s.downloadResponseCache.Get(c, cacheKey)
				if packageItem == nil {
					packageItem, err = api.GetPackage(hostname, namespace, pkg, version, platform.Os, platform.Arch)
					if err != nil {
						s.logger.error(err.Error())
						c.AsciiJSON(http.StatusInternalServerError, ErrorResponse{
							Code:    http.StatusInternalServerError,
							Message: fmt.Sprintf("unable to download package for %s", pkgPath),
						})
						return
					}
					err := s.downloadResponseCache.Set(c, cacheKey, packageItem)
					if err != nil {
						s.logger.warn(err.Error())
					}
				}
				cacheKey = fmt.Sprintf("hashes:%s", packageItem.DownloadUrl)
				hashes, _ := s.hashesCache.Get(c, cacheKey)
				if hashes == nil {
					hashes, err = hash.GetHashes(packageItem.DownloadUrl)
					if err != nil {
						s.logger.error(err.Error())
						c.AsciiJSON(http.StatusInternalServerError, ErrorResponse{
							Code:    http.StatusInternalServerError,
							Message: fmt.Sprintf("unable to get checksums for package %s", pkgPath),
						})
						return
					}
					err := s.hashesCache.Set(c, cacheKey, hashes)
					if err != nil {
						s.logger.warn(err.Error())
					}
				}
				response.Archives[platform.Os+"_"+platform.Arch] = Build{
					Url:    packageItem.DownloadUrl,
					Hashes: hashes,
				}
			}
		}
	}
	c.AsciiJSON(http.StatusOK, response)
}

type ErrorResponse struct {
	Message string
	Code    int
}

type IndexResponse struct {
	Versions map[string]struct{} `json:"versions"`
}

type VersionResponse struct {
	Archives map[string]Build `json:"archives"`
}

type Build struct {
	Url    string   `json:"url"`
	Hashes []string `json:"hashes"`
}
