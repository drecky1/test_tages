package service

import file "test_tages/proto"

type FileService struct {
	uploadLimiter   chan struct{}
	downloadLimiter chan struct{}
	listLimiter     chan struct{}
	dataDir         string
	file.UnimplementedFileServiceServer
}
