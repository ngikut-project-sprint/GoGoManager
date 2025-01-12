package services

import (
	"mime/multipart"

	"github.com/ngikut-project-sprint/GoGoManager/internal/utils"
)

type FileService interface {
	Create(file multipart.FileHeader) *utils.GoGoError
}

type ValidFileSizeFunc func(file multipart.FileHeader) error

type ValidFileTypeFunc func(file multipart.FileHeader) error

type fileService struct {
	validateFileSize ValidFileSizeFunc
	validateFileType ValidFileTypeFunc
}

func NewFileService(
	validateFileSize ValidFileSizeFunc,
	validateFileType ValidFileTypeFunc,
) FileService {
	return &fileService{
		validateFileSize: validateFileSize,
		validateFileType: validateFileType,
	}
}

func (s *fileService) Create(file multipart.FileHeader) *utils.GoGoError {
	fileSizeErr := s.validateFileSize(file)
	if fileSizeErr != nil {
		return utils.WrapError(fileSizeErr, utils.InvalidFileSize, "Invalid file size")
	}

	fileTypeErr := s.validateFileType(file)
	if fileTypeErr != nil {
		return utils.WrapError(fileTypeErr, utils.InvalidFileType, "Invalid file type")
	}

	return nil
}
