package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/ngikut-project-sprint/GoGoManager/internal/config"
	"github.com/ngikut-project-sprint/GoGoManager/internal/services"
	"github.com/ngikut-project-sprint/GoGoManager/internal/utils"
)

type FileHandler struct {
	fileService services.FileService
}

func NewFileHandler(fileService services.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}

func (h *FileHandler) PostFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, header, err := r.FormFile("file")
	if err != nil {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fileError := h.fileService.Create(*header)
	if fileError != nil {
		switch fileError.Type {
		case utils.InvalidFileSize:
			utils.SendErrorResponse(w, "Invalid file size", http.StatusBadRequest)
			return
		case utils.InvalidFileType:
			utils.SendErrorResponse(w, "Invalid file type", http.StatusBadRequest)
			return
		}
	}

	cfg, err := config.Get()
	if err != nil {
		utils.SendErrorResponse(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	awsCreds := credentials.NewStaticCredentialsProvider(
		cfg.AWS.AccessKey,
		cfg.AWS.SecretKey,
		"",
	)

	awsCfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
		awsConfig.WithRegion(cfg.AWS.Region),
		awsConfig.WithCredentialsProvider(awsCreds))
	if err != nil {
		utils.SendErrorResponse(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	client := s3.NewFromConfig(awsCfg)

	f, err := header.Open()
	if err != nil {
		utils.SendErrorResponse(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	uploader := manager.NewUploader(client)
	id := uuid.New()
	pathName := fmt.Sprintf("ngikut/gogomanager/%s.png", &id)

	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(cfg.AWS.BucketName),
		Key:    aws.String(pathName),
		Body:   f,
		ACL:    "public-read",
	})
	if err != nil {
		utils.SendErrorResponse(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	region := awsCfg.Region
	uri := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", cfg.AWS.BucketName, region, pathName)
	fileResponse := utils.FileResponse{
		Uri: uri,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(fileResponse); err != nil {
		log.Println("Failed to write response:", err)
		utils.SendErrorResponse(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *FileHandler) File(w http.ResponseWriter, r *http.Request) {
	http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				h.PostFile(w, r) // Handle POST requests
			default:
				utils.SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}).ServeHTTP(w, r)
}
