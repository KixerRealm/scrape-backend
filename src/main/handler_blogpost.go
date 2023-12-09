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
		Title       string `json:"title"`
		Description string `json:"description"`
		Files       []struct {
			Filename   string `json:"file_name"`
			FolderName string `json:"folder_name"`
			Content    string `json:"content"`
		}
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	blogPost, err := apiCfg.DB.CreateBlogPost(r.Context(), database.CreateBlogPostParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Title:       params.Title,
		Description: params.Description,
		UserID:      user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create blog_post: %s", err))
		return
	}

	if len(params.Files) > 0 && params.Files != nil {
		fileIDs := apiCfg.saveFile(params.Files, w, r)
		for _, fileID := range fileIDs {
			_, err := apiCfg.DB.CreateBlogPostFile(r.Context(), database.CreateBlogPostFileParams{
				BlogPostID: blogPost.ID,
				FileID:     fileID,
			})
			if err != nil {
				respondWithError(w, 500, fmt.Sprintf("Error associating file with blog post: %s", err))
				return
			}
		}
	}

	respondWithJSON(w, 200, databaseBlogPostToBlogPost(blogPost))
}

func (apiCfg *apiConfig) handlerGetBlogPostsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	blogPosts, err := apiCfg.DB.GetBlogPostsByUserWithFiles(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get blog posts by user: %s", err))
		return
	}
	respondWithJSON(w, 200, apiCfg.databaseBlogPostsWithFilesToBlogPostsWithFiles(blogPosts, r))
}

func (apiCfg *apiConfig) handlerGetAllBlogPosts(w http.ResponseWriter, r *http.Request) {
	blogPosts, err := apiCfg.DB.GetBlogPosts(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get blog posts: %s", err))
		return
	}
	respondWithJSON(w, 200, apiCfg.databaseAllBlogPostsWithFilesToAllBlogPostsWithFiles(blogPosts, r))
}

func (apiCfg *apiConfig) handlerGetPatchNotes(w http.ResponseWriter, r *http.Request) {
	blogPosts, err := apiCfg.DB.GetBlogPostsByCreatedAt(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get patch notes: %s", err))
		return
	}
	respondWithJSON(w, 200, databasePatchNotesToPatchNotes(blogPosts))
}
