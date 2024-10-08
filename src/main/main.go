package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"net/http"
	"os"
	"scrape-backend/src/main/internal/database"
	_ "time"
)

type apiConfig struct {
	DB    *database.Queries
	Minio *minio.Client
}

func main() {
	minioServerEndpoint := "127.0.0.1:9000"

	accessKey := "username123"
	secretKey := "password123"

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	minioClient, err := minio.New(minioServerEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		fmt.Println("Error initializing Minio client:", err)
		return
	}

	db := database.New(conn)
	apiCfg := apiConfig{
		DB:    db,
		Minio: minioClient,
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users/register", apiCfg.registerUserHandler)
	v1Router.Post("/users/login", apiCfg.loginUserHandler)
	v1Router.Post("/blog-posts/create", apiCfg.middlewareAuth(apiCfg.handlerCreateBlogPost))
	v1Router.Post("/my-blog-posts", apiCfg.middlewareAuth(apiCfg.handlerGetBlogPostsByUser))
	v1Router.Post("/all-blog-posts", apiCfg.handlerGetAllBlogPosts)
	v1Router.Post("/bug-reports/create", apiCfg.middlewareAuth(apiCfg.handlerCreateBugReport))
	v1Router.Post("/my-bug-reports", apiCfg.middlewareAuth(apiCfg.handlerGetBugReportsByUser))
	v1Router.Post("/all-bug-reports", apiCfg.handlerGetAllBugReports)
	v1Router.Get("/patch-notes", apiCfg.handlerGetPatchNotes)
	v1Router.Post("/files/create", apiCfg.handlerCreateFile)
	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Server starting on port  %v", portString)
	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
