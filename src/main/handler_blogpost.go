package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"scrape-backend/src/main/internal/database"
	"time"
)

func (apiCfg *apiConfig) handlerCreateBlogPost(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Title         string `json:"title"`
		Description   string `json:"description"`
		ImageFilename string `json:"image_filename"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	blogPosts, err := apiCfg.DB.CreateBlogPost(r.Context(), database.CreateBlogPostParams{
		ID:            uuid.New(),
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
		Title:         params.Title,
		Description:   params.Description,
		ImageFilename: params.ImageFilename,
		UserID:        user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed follow: %s", err))
		return
	}
	respondWithJSON(w, 200, databaseBlogPostToBlogPost(blogPosts))
}

func (apiCfg *apiConfig) handlerGetBlogPostsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	blogPosts, err := apiCfg.DB.GetBlogPostsByUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feed follows: %s", err))
		return
	}
	respondWithJSON(w, 200, databaseBlogPostsToBlogPosts(blogPosts))
}

func (apiCfg *apiConfig) handlerGetAllBlogPosts(w http.ResponseWriter, r *http.Request) {
	blogPosts, err := apiCfg.DB.GetBlogPosts(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feed follows: %s", err))
		return
	}
	respondWithJSON(w, 200, databaseBlogPostsToBlogPosts(blogPosts))
}

func (apiCfg *apiConfig) handlerGetPatchNotes(w http.ResponseWriter, r *http.Request) {
	blogPosts, err := apiCfg.DB.GetBlogPostsByCreatedAt(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feed follows: %s", err))
		return
	}
	respondWithJSON(w, 200, databasePatchNotesToPatchNotes(blogPosts))
}
