package usecase

import (
	"bytes"
	"context"
	"log"
	"os"
	"time"

	"github.com/Nishad4140/SkillSync_ProjectService/internal/adapter"
	"github.com/Nishad4140/SkillSync_ProtoFiles/pb"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ProjectUsecase struct {
	projectAdapter adapter.AdapterInterface
}

func NewProjectUsecase(projectAdapter adapter.AdapterInterface) *ProjectUsecase {
	return &ProjectUsecase{
		projectAdapter: projectAdapter,
	}
}

func (user *ProjectUsecase) UploadFreelancerFile(req *pb.FileRequest) (string, error) {
	minioClient, err := minio.New(os.Getenv("MINIO_ENDPOINT"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ACCESSKEY"), os.Getenv("MINIO_SECRETKEY"), ""),
		Secure: false,
	})
	if err != nil {
		log.Println("error while initializing minio", err)
		return "", err
	}
	objectName := "files/" + req.ObjectName
	contentType := "application/zip"

	n, err := minioClient.PutObject(context.Background(), os.Getenv("BUCKET_NAME"), objectName, bytes.NewReader(req.ImageData), int64(len(req.ImageData)), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Println("error while uploading to minio", err)
		return "", err
	}
	log.Printf("Successfully uploaded %s of size %v to minio", objectName, n)
	presignedURL, err := minioClient.PresignedGetObject(context.Background(), os.Getenv("BUCKET_NAME"), objectName, time.Second*24*60*60, nil)
	if err != nil {
		log.Println("error while generating the url", err)
		return "", err
	}
	url, err := user.projectAdapter.UploadFile(presignedURL.String(), req.ProjectId)

	return url, err

}
