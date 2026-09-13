package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/sudzekai/web-os-api/executor"
	logging "github.com/sudzekai/web-os-api/logging/core"
	"github.com/sudzekai/web-os-api/modules/files-module/internal/objects"
	"github.com/sudzekai/web-os-api/modules/files-module/internal/objects/types"
	"github.com/sudzekai/web-os-api/server/abstractions"
	serverTypes "github.com/sudzekai/web-os-api/server/types"
)

type filesController struct{}

var FilesController = filesController{}

func (ctr filesController) Connect(srv abstractions.IServer) {
	srv.AddHandler("GET /files", ctr.GetFiles)
}

func (filesController) GetFiles(r *http.Request) (result serverTypes.MethodResult) {
	log := logging.NewLogger("FilesModule:getFiles")

	path := r.URL.Query().Get("path")

	if path == "" {
		path = "."
	}

	cmdResult := executor.ExecuteInDirectory(path, "ls", "-a")

	log.LogDebug("вывод: %s", cmdResult.Stdout)

	if cmdResult.Error != nil {
		result.StatusCode = 500
		result.Error = cmdResult.Error
		return
	}

	files := make([]objects.FileInfo, 0)

	for file := range strings.SplitSeq(strings.Replace(cmdResult.Stdout, "\n", " ", -1), "\n") {
		log.LogDebug("найден файл %s", file)

		pattern := fmt.Sprintf("$\"%s\"", strings.Join(getSimpleFileInfoPattern(path), ""))
		args := []string{file, "-c", pattern}
		stat := executor.ExecuteInDirectory("stat", path, args...)

		files = append(files, *objects.NewFileInfo(stat.Stdout))
	}

	result.StatusCode = 200
	result.Data = files
	return
}

func getSimpleFileInfoPattern(path string) []string {
	patterns := make([]string, 0)

	patterns = append(patterns, fmt.Sprintf("%s:%s\\n", types.Type.Key, types.Type.Pattern))
	patterns = append(patterns, fmt.Sprintf("%s:%s\\n", types.Name.Key, types.Name.Pattern))
	patterns = append(patterns, fmt.Sprintf("%s:%s\\n", types.FullName.Key, filepath.Join(path, types.Name.Pattern)))
	patterns = append(patterns, fmt.Sprintf("%s:%s\\n", types.Size.Key, types.Size.Pattern))

	return patterns
}
