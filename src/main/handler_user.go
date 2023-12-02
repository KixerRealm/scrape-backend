package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"scrape-backend/src/main/internal/database"
	"time"
)

func (apiCfg *apiConfig) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// Define a struct for request parameters
	type registrationParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}

	// Decode JSON request body
	decoder := json.NewDecoder(r.Body)
	var params registrationParams
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	// Check if the user already exists in the database
	_, err := apiCfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error checking user existence: %s", err))
			return
		}
	} else {
		// User already exists, respond with an error
		respondWithError(w, http.StatusConflict, "User with this email already exists")
		return
	}

	var errHash error
	params.Password, errHash = generateHashPassword(params.Password)
	if errHash != nil {
		respondWithError(w, http.StatusConflict, "Could not generate password hash")
		return
	}

	// Create a new user
	newUser, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Email:     params.Email,
		Password:  params.Password,
		Username:  params.Username,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't create user: %s", err))
		return
	}

	// Respond with the newly created user
	respondWithJSON(w, http.StatusOK, databaseUserToUser(newUser))
}

func (apiCfg *apiConfig) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	// Define a struct for login credentials
	type loginCredentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode JSON request body
	decoder := json.NewDecoder(r.Body)
	var credentials loginCredentials
	if err := decoder.Decode(&credentials); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	// Check if the user exists in the database
	existingUser, err := apiCfg.DB.GetUserByEmail(r.Context(), credentials.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error checking user existence: %s", err))
			return
		}
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Compare passwords
	errHash := compareHashPassword(credentials.Password, existingUser.Password)
	if !errHash {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate a JWT token with an expiration time (e.g., 1 hour)
	expirationTime := time.Now().Add(1 * time.Hour)
	token, err := generateJWT(existingUser.Email, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error generating JWT: %s", err))
		return
	}

	// Update the user's token in the database
	err = apiCfg.DB.UpdateUserToken(r.Context(), database.UpdateUserTokenParams{
		ID:    existingUser.ID,
		Token: token,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating user token: %s", err))
		return
	}

	// Respond with the JWT token
	respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}
