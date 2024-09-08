package service

import (
	"context"
	file "test_tages/proto"
	"test_tages/server/constants"
)

func NewFileService(dataDir string) *FileService {
	return &FileService{
		uploadLimiter:   make(chan struct{}, constants.MaxUpload),
		downloadLimiter: make(chan struct{}, constants.MaxDownload),
		listLimiter:     make(chan struct{}, constants.MaxList),
		dataDir:         dataDir,
	}
}

type Interface interface {
	UploadFile(stream file.FileService_UploadFileServer) error
	DownloadFile(req *file.DownloadRequest, stream file.FileService_DownloadFileServer) error
	ListFiles(_ context.Context, _ *file.Empty) (*file.FileList, error)
}
