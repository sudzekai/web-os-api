package controllers

import (
	"net/http"
	"strings"

	"github.com/sudzekai/web-os-api/executor"
	"github.com/sudzekai/web-os-api/server/abstractions"
	"github.com/sudzekai/web-os-api/server/types"
)

type filesController struct{}

var FilesController = filesController{}

func (ctr filesController) Connect(srv abstractions.IServer) {
	srv.AddHandler("GET /files", ctr.GetFiles)
}

func (filesController) GetFiles(r *http.Request) types.MethodResult {
	dirPath := r.URL.Query().Get("directoryPath")

	if dirPath == "" {
		dirPath = "."
	}

	result := executor.ExecuteInDirectory(dirPath, "ls", "-la")

	files := make([]string, 0)

	for file := range strings.SplitSeq(result.Stdout, " ") {
		files = append(files, file)
	}

	if result.Error != nil {
		return types.MethodResult{
			StatusCode: 500,
			Error:      result.Error,
		}
	}

	return types.MethodResult{
		StatusCode: 200,
		Data:       files,
	}
}
