package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"test_tages/example_client/constants"
	"test_tages/example_client/handlers"
	pb "test_tages/proto"
)

func main() {
	creds := insecure.NewCredentials()
	conn, err := grpc.NewClient(constants.ServerAddress, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	client := pb.NewFileServiceClient(conn)

	// Пример использования
	for {
		fmt.Println("Выберите действие:")
		fmt.Println("1. Загрузить файл")
		fmt.Println("2. Скачать файл")
		fmt.Println("3. Просмотреть список файлов")
		fmt.Println("4. Выход")

		var choice int
		_, err := fmt.Scan(&choice)
		if err != nil {
			return
		}

		switch choice {
		case 1:
			handlers.UploadFile(client)
		case 2:
			handlers.DownloadFile(client)
		case 3:
			handlers.ListFiles(client)
		case 4:
			os.Exit(0)
		default:
			fmt.Println("Неверный выбор. Пожалуйста, попробуйте снова.")
		}
	}
}
