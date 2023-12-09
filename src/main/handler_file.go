package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"scrape-backend/src/main/internal/database"
)

func (apiCfg *apiConfig) handlerCreateFile(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Filename   string `json:"file_name"`
		FolderName string `json:"folder_name"`
		Content    string `json:"content"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(params.Content)
	if err != nil {
		fmt.Println("Error decoding Base64:", err)
		return
	}

	err = ioutil.WriteFile(params.Filename, []byte(string(decodedBytes)), 0644)
	if err != nil {
		fmt.Println("Error saving to file:", err)
		return
	}

	file, err := apiCfg.DB.CreateFile(r.Context(), database.CreateFileParams{
		ID:         uuid.New(),
		FileName:   params.Filename,
		FolderName: params.FolderName,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create file: %s", err))
		return
	}
	respondWithJSON(w, 200, file)
}
