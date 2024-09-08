package service

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	file "test_tages/proto"
	"test_tages/server/constants"
)

func (s *FileService) UploadFile(stream file.FileService_UploadFileServer) error {
	s.uploadLimiter <- struct{}{}
	defer func() { <-s.uploadLimiter }()

	var chunks [][]byte
	var filename string
	for {
		chunk, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		chunks = append(chunks, chunk.Data)
		filename = chunk.Filename
	}

	data := bytes.Join(chunks, nil)
	path := filepath.Join(s.dataDir, filename)
	if err := os.WriteFile(path, data, os.ModePerm); err != nil {
		return err
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	return stream.SendAndClose(&file.UploadResponse{
		Filename: filename,
		Size:     fileInfo.Size(),
	})
}

func (s *FileService) DownloadFile(req *file.DownloadRequest, stream file.FileService_DownloadFileServer) error {
	s.downloadLimiter <- struct{}{}
	defer func() { <-s.downloadLimiter }()

	path := filepath.Join(s.dataDir, req.Filename)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)

	buf := make([]byte, constants.DownloadBuffer)
	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		err = stream.Send(&file.FileChunk{
			Data: buf[:n],
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *FileService) ListFiles(_ context.Context, _ *file.Empty) (*file.FileList, error) {
	s.listLimiter <- struct{}{}
	defer func() { <-s.listLimiter }()

	var files []*file.FileInfo
	err := filepath.Walk(s.dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, &file.FileInfo{
				Filename:  info.Name(),
				CreatedAt: info.ModTime().Unix(),
				UpdatedAt: info.ModTime().Unix(),
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &file.FileList{
		Files: files,
	}, err
}
