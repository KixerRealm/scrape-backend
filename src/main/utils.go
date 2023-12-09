package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"scrape-backend/src/main/internal/database"
)

func generateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func compareHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type FileRequest []struct {
	Filename   string `json:"file_name"`
	FolderName string `json:"folder_name"`
	Content    string `json:"content"`
}

func (apiCfg *apiConfig) saveFile(files FileRequest, w http.ResponseWriter, r *http.Request) []uuid.UUID {
	var fileIDs []uuid.UUID
	for _, fileParam := range files {
		decodedBytes, err := base64.StdEncoding.DecodeString(fileParam.Content)
		if err != nil {
			fmt.Println("Error decoding Base64:", err)
			return nil
		}

		_, err = apiCfg.Minio.PutObject(r.Context(), "files", fileParam.Filename, bytes.NewReader(decodedBytes), int64(len(decodedBytes)), minio.PutObjectOptions{})
		if err != nil {
			fmt.Println("Error saving to Minio:", err)
			return nil
		}

		fileID := uuid.New()
		_, error1 := apiCfg.DB.CreateFile(r.Context(), database.CreateFileParams{
			ID:         fileID,
			FileName:   fileParam.Filename,
			FolderName: fileParam.FolderName,
		})
		if error1 != nil {
			respondWithError(w, 500, fmt.Sprintf("Error creating file: %s", err))
			return nil
		}

		fileIDs = append(fileIDs, fileID)
	}
	return fileIDs
}

func (apiCfg *apiConfig) getFile(r *http.Request, filename string) string {
	object, err := apiCfg.Minio.GetObject(r.Context(), "files", filename, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println("Error getting object:", err)
	}

	objectContent, err := io.ReadAll(object)
	if err != nil {
		fmt.Println("Error reading object content:", err)
	}
	defer object.Close()

	base64String := base64.StdEncoding.EncodeToString(objectContent)
	return base64String
}
