package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"test_tages/example_client/constants"
	c "test_tages/proto"
	"time"
)

func UploadFile(client c.FileServiceClient) {
	fmt.Print("Введите путь к файлу для загрузки: ")
	var filePath string
	_, err := fmt.Scan(&filePath)
	if err != nil {
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("не удалось открыть файл: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("не удалось получить информацию о файле: %v", err)
	}

	stream, err := client.UploadFile(context.Background())
	if err != nil {
		log.Fatalf("не удалось создать поток: %v", err)
	}

	buffer := make([]byte, constants.UploadBuffer)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("не удалось прочитать файл: %v", err)
		}

		if err := stream.Send(&c.FileChunk{
			Data:     buffer[:n],
			Filename: fileInfo.Name(),
		}); err != nil {
			log.Fatalf("не удалось отправить часть файла: %v", err)
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("не удалось получить ответ: %v", err)
	}

	fmt.Printf("Файл '%s' успешно загружен, размер: %d байт\n", response.Filename, response.Size)
}

func DownloadFile(client c.FileServiceClient) {
	fmt.Print("Введите имя файла для скачивания: ")
	var filename string
	_, err := fmt.Scan(&filename)
	if err != nil {
		return
	}

	stream, err := client.DownloadFile(context.Background(), &c.DownloadRequest{Filename: filename})
	if err != nil {
		log.Fatalf("не удалось скачать файл: %v", err)
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
		chunk, err := stream.Recv()
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

	fmt.Printf("Файл '%s' успешно скачан\n", filename)
}

func ListFiles(client c.FileServiceClient) {
	fileList, err := client.ListFiles(context.Background(), &c.Empty{})
	if err != nil {
		log.Fatalf("не удалось получить список файлов: %v", err)
	}

	tf := func(t int64) string {
		return time.Unix(t, 0).Format(time.RFC3339)
	}
	for _, fileInfo := range fileList.Files {
		fmt.Printf("Имя файла: `%s` | Дата создания: `%s` | Дата обновления: `%s`\n",
			fileInfo.Filename,
			tf(fileInfo.CreatedAt),
			tf(fileInfo.UpdatedAt),
		)
	}
}
