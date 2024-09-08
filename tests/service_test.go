package tests

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"net"
	"os"
	file "test_tages/proto"
	"test_tages/server/service"
	"testing"
)

func TestFileService(t *testing.T) {
	dataDir := t.TempDir()
	fileService := service.NewFileService(dataDir)
	s := grpc.NewServer()
	file.RegisterFileServiceServer(s, fileService)

	lis, err := net.Listen("tcp", ":7777")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	go func(lis net.Listener) {
		if err := s.Serve(lis); err != nil {
			log.Fatalf(err.Error())
		}
		defer s.Stop()
	}(lis)

	creds := insecure.NewCredentials()
	conn, err := grpc.NewClient(":7777", grpc.WithTransportCredentials(creds))
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	client := file.NewFileServiceClient(conn)

	filename := "test.txt"
	data := []byte("Hello, world!")

	stream, err := client.UploadFile(context.Background())
	if err != nil {
		t.Fatalf("failed to upload file: %v", err)
	}
	err = stream.Send(&file.FileChunk{
		Data:     data,
		Filename: filename,
	})
	if err != nil {
		t.Fatalf("failed to send file chunk: %v", err)
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatalf("failed to close and receive: %v", err)
	}
	if resp.Filename != filename {
		t.Errorf("expected filename %s, got %s", filename, resp.Filename)
	}
	if resp.Size != int64(len(data)) {
		t.Errorf("expected size %d, got %d", len(data), resp.Size)
	}

	st, err := client.DownloadFile(context.Background(), &file.DownloadRequest{Filename: filename})
	if err != nil {
		t.Fatalf("failed to upload file: %v", err)
	}
	outFile, err := os.Create(filename)
	if err != nil {
		log.Fatalf("не удалось создать файл: %v", err)
	}
	defer func(outFile *os.File) {
		err := outFile.Close()
		if err != nil {
			return
		}
	}(outFile)

	for {
		chunk, err := st.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("не удалось получить данные: %v", err)
		}

		if _, err := outFile.Write(chunk.Data); err != nil {
			log.Fatalf("не удалось записать данные в файл: %v", err)
		}
	}
	if outFile.Name() != filename {
		t.Errorf("не совпадают навзвания файлов")
	}

	defer func(filename string) { _ = os.Remove(filename) }(filename)

	fileList, err := client.ListFiles(context.Background(), &file.Empty{})
	if err != nil {
		t.Fatalf("failed to list files: %v", err)
	}
	if len(fileList.Files) != 1 {
		t.Errorf("expected 1 file, got %d", len(fileList.Files))
	}
	fileInfo := fileList.Files[0]
	if fileInfo.Filename != filename {
		t.Errorf("expected filename %s, got %s", filename, fileInfo.Filename)
	}
	if fileInfo.CreatedAt > fileInfo.UpdatedAt {
		t.Errorf("created at %d should be less than or equal to updated at %d", fileInfo.CreatedAt, fileInfo.UpdatedAt)
	}
}
