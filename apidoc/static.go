package apidoc

import (
	"net/http"
	"strings"

	"github.com/yusufsyaifudin/go-bindata-assetfs"
)

// see this tutorial https://medium.com/@erinus/go-my-way-day-3-9c9b420ed43e
type apidocGoBindata struct {
	FileSystem http.FileSystem
}

func (apidocGoBindata *apidocGoBindata) Open(name string) (http.File, error) {
	return apidocGoBindata.FileSystem.Open(name)
}

func (apidocGoBindata *apidocGoBindata) Exists(prefix string, filepath string) bool {
	var err error
	var url string
	url = strings.TrimPrefix(filepath, prefix)
	if len(url) < len(filepath) {
		_, err = apidocGoBindata.FileSystem.Open(url)
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func Static() *apidocGoBindata {
	var fs *assetfs.AssetFS
	fs = &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
	}
	return &apidocGoBindata{fs}
}
