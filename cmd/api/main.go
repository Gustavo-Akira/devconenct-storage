package main

import (
	"log"
	"os"

	uploadfile "devconnectstorage/internal/application/usecase/upload_file"
	"devconnectstorage/internal/infraestructure/inbound/rest"
	"devconnectstorage/internal/infraestructure/outbound/generator/uuidgen"
	"devconnectstorage/internal/infraestructure/outbound/repository/file/mongodb"
	minioStorage "devconnectstorage/internal/infraestructure/outbound/storage/minio"

	"github.com/gin-gonic/gin"
)

func main() {

	mongoURI := os.Getenv("MONGO_URI")
	mongoDB := os.Getenv("MONGO_DB")
	mongoCollection := os.Getenv("MONGO_COLLECTION")

	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	minioUser := os.Getenv("MINIO_USER")
	minioPassword := os.Getenv("MINIO_PASSWORD")
	minioBucket := os.Getenv("MINIO_BUCKET")
	minioSSL := os.Getenv("MINIO_USE_SSL") == "true"

	fileRepo, err := mongodb.NewMongoFileRepository(mongoURI, "", "", mongoDB, mongoCollection)
	if err != nil {
		log.Fatalf("failed to initialize Mongo repository: %v", err)
	}

	storage, err := minioStorage.NewMinIOStorage(minioEndpoint, minioUser, minioPassword, minioSSL, minioBucket)
	if err != nil {
		log.Fatalf("failed to initialize MinIO storage: %v", err)
	}

	idGenerator := uuidgen.UUIDGenerator{}

	uploadFileUseCase := uploadfile.NewUploadFileUseCase(fileRepo, storage, idGenerator)

	fileController := rest.NewFileRestController(uploadFileUseCase)

	router := gin.Default()
	router.POST("/files", fileController.UploadFile)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
