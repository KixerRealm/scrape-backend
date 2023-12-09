package main

import (
	"github.com/google/uuid"
	"net/http"
	"scrape-backend/src/main/internal/database"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Username  string    `json:"username"`
	Token     string    `json:"token"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
		Username:  dbUser.Username,
		Token:     dbUser.Token,
	}
}

type BlogPost struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	ImageFilename string    `json:"image_filename"`
	UserID        uuid.UUID `json:"user_id"`
}

type BlogPostResponse struct {
	BlogPostID          uuid.UUID `json:"id"`
	BlogPostCreatedAt   time.Time `json:"created_at"`
	BlogPostUpdatedAt   time.Time `json:"updated_at"`
	BlogPostTitle       string    `json:"title"`
	BlogPostDescription string    `json:"description"`
	Files               []File    `json:"files"`
}

type BugReportResponse struct {
	BugReportID          uuid.UUID `json:"id"`
	BugReportCreatedAt   time.Time `json:"created_at"`
	BugReportUpdatedAt   time.Time `json:"updated_at"`
	BugReportTitle       string    `json:"title"`
	BugReportDescription string    `json:"description"`
	Files                []File    `json:"files"`
}

type File struct {
	FileID     uuid.UUID `json:"id"`
	FileName   string    `json:"file_name"`
	FolderName string    `json:"folder_name"`
	Content    string    `json:"content"`
}

func (apiCfg *apiConfig) databaseBlogPostsWithFilesToBlogPostsWithFiles(dbBlogPosts []database.GetBlogPostsByUserWithFilesRow, r *http.Request) []BlogPostResponse {
	responseMap := make(map[uuid.UUID]BlogPostResponse)

	for _, dbBlogPost := range dbBlogPosts {

		if response, ok := responseMap[dbBlogPost.BlogPostID]; ok {
			response.Files = append(response.Files, File{
				FileID:     dbBlogPost.FileID,
				FileName:   dbBlogPost.FileName,
				FolderName: dbBlogPost.FolderName,
				Content:    apiCfg.getFile(r, dbBlogPost.FileName),
			})
			responseMap[dbBlogPost.BlogPostID] = response
		} else {
			responseMap[dbBlogPost.BlogPostID] = BlogPostResponse{
				BlogPostID:          dbBlogPost.BlogPostID,
				BlogPostCreatedAt:   dbBlogPost.BlogPostCreatedAt,
				BlogPostUpdatedAt:   dbBlogPost.BlogPostUpdatedAt,
				BlogPostTitle:       dbBlogPost.BlogPostTitle,
				BlogPostDescription: dbBlogPost.BlogPostDescription,
				Files: []File{
					{
						FileID:     dbBlogPost.FileID,
						FileName:   dbBlogPost.FileName,
						FolderName: dbBlogPost.FolderName,
						Content:    apiCfg.getFile(r, dbBlogPost.FileName),
					},
				},
			}
		}
	}

	var result []BlogPostResponse
	for _, response := range responseMap {
		result = append(result, response)
	}

	return result
}

func (apiCfg *apiConfig) databaseAllBlogPostsWithFilesToAllBlogPostsWithFiles(dbBlogPosts []database.GetBlogPostsRow, r *http.Request) []BlogPostResponse {
	responseMap := make(map[uuid.UUID]BlogPostResponse)

	for _, dbBlogPost := range dbBlogPosts {

		if response, ok := responseMap[dbBlogPost.BlogPostID]; ok {
			response.Files = append(response.Files, File{
				FileID:     dbBlogPost.FileID,
				FileName:   dbBlogPost.FileName,
				FolderName: dbBlogPost.FolderName,
				Content:    apiCfg.getFile(r, dbBlogPost.FileName),
			})
			responseMap[dbBlogPost.BlogPostID] = response
		} else {
			responseMap[dbBlogPost.BlogPostID] = BlogPostResponse{
				BlogPostID:          dbBlogPost.BlogPostID,
				BlogPostCreatedAt:   dbBlogPost.BlogPostCreatedAt,
				BlogPostUpdatedAt:   dbBlogPost.BlogPostUpdatedAt,
				BlogPostTitle:       dbBlogPost.BlogPostTitle,
				BlogPostDescription: dbBlogPost.BlogPostDescription,
				Files: []File{
					{
						FileID:     dbBlogPost.FileID,
						FileName:   dbBlogPost.FileName,
						FolderName: dbBlogPost.FolderName,
						Content:    apiCfg.getFile(r, dbBlogPost.FileName),
					},
				},
			}
		}
	}

	var result []BlogPostResponse
	for _, response := range responseMap {
		result = append(result, response)
	}

	return result
}

func databaseBlogPostToBlogPost(dbBlogPost database.BlogPost) BlogPost {
	return BlogPost{
		ID:          dbBlogPost.ID,
		CreatedAt:   dbBlogPost.CreatedAt,
		UpdatedAt:   dbBlogPost.UpdatedAt,
		Title:       dbBlogPost.Title,
		Description: dbBlogPost.Description,
		UserID:      dbBlogPost.UserID,
	}
}

func databaseBlogPostsToBlogPosts(dbBlogPosts []database.BlogPost) []BlogPost {
	blogPosts := []BlogPost{}
	for _, dbBlogPost := range dbBlogPosts {
		blogPosts = append(blogPosts, databaseBlogPostToBlogPost(dbBlogPost))
	}
	return blogPosts
}

type BugReport struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      uuid.UUID `json:"user_id"`
}

func (apiCfg *apiConfig) databaseBugReportsWithFilesToBugReportsWithFiles(dbBugReports []database.GetBugReportsByUserWithFilesRow, r *http.Request) []BugReportResponse {
	responseMap := make(map[uuid.UUID]BugReportResponse)

	for _, dbBugReport := range dbBugReports {

		if response, ok := responseMap[dbBugReport.BugReportID]; ok {
			response.Files = append(response.Files, File{
				FileID:     dbBugReport.FileID,
				FileName:   dbBugReport.FileName,
				FolderName: dbBugReport.FolderName,
				Content:    apiCfg.getFile(r, dbBugReport.FileName),
			})
			responseMap[dbBugReport.BugReportID] = response
		} else {
			responseMap[dbBugReport.BugReportID] = BugReportResponse{
				BugReportID:          dbBugReport.BugReportID,
				BugReportCreatedAt:   dbBugReport.BugReportCreatedAt,
				BugReportUpdatedAt:   dbBugReport.BugReportUpdatedAt,
				BugReportTitle:       dbBugReport.BugReportTitle,
				BugReportDescription: dbBugReport.BugReportDescription,
				Files: []File{
					{
						FileID:     dbBugReport.FileID,
						FileName:   dbBugReport.FileName,
						FolderName: dbBugReport.FolderName,
						Content:    apiCfg.getFile(r, dbBugReport.FileName),
					},
				},
			}
		}
	}

	var result []BugReportResponse
	for _, response := range responseMap {
		result = append(result, response)
	}

	return result
}

func (apiCfg *apiConfig) databaseAllBugReportsWithFilesToAllBugReportsWithFiles(dbBugReports []database.GetBugReportsRow, r *http.Request) []BugReportResponse {
	responseMap := make(map[uuid.UUID]BugReportResponse)

	for _, dbBugReport := range dbBugReports {

		if response, ok := responseMap[dbBugReport.BugReportID]; ok {
			response.Files = append(response.Files, File{
				FileID:     dbBugReport.FileID,
				FileName:   dbBugReport.FileName,
				FolderName: dbBugReport.FolderName,
				Content:    apiCfg.getFile(r, dbBugReport.FileName),
			})
			responseMap[dbBugReport.BugReportID] = response
		} else {
			responseMap[dbBugReport.BugReportID] = BugReportResponse{
				BugReportID:          dbBugReport.BugReportID,
				BugReportCreatedAt:   dbBugReport.BugReportCreatedAt,
				BugReportUpdatedAt:   dbBugReport.BugReportUpdatedAt,
				BugReportTitle:       dbBugReport.BugReportTitle,
				BugReportDescription: dbBugReport.BugReportDescription,
				Files: []File{
					{
						FileID:     dbBugReport.FileID,
						FileName:   dbBugReport.FileName,
						FolderName: dbBugReport.FolderName,
						Content:    apiCfg.getFile(r, dbBugReport.FileName),
					},
				},
			}
		}
	}

	var result []BugReportResponse
	for _, response := range responseMap {
		result = append(result, response)
	}

	return result
}

func databaseBugReportToBugReport(dbBugReport database.BugReport) BugReport {
	return BugReport{
		ID:          dbBugReport.ID,
		CreatedAt:   dbBugReport.CreatedAt,
		UpdatedAt:   dbBugReport.UpdatedAt,
		Title:       dbBugReport.Title,
		Description: dbBugReport.Description,
		UserID:      dbBugReport.UserID,
	}
}

func databaseBugReportsToBugReports(dbBlogPosts []database.BugReport) []BugReport {
	bugReports := []BugReport{}
	for _, dbBlogPost := range dbBlogPosts {
		bugReports = append(bugReports, databaseBugReportToBugReport(dbBlogPost))
	}
	return bugReports
}

type PatchNote struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	ImageFilename string    `json:"image_filename"`
	UserID        uuid.UUID `json:"user_id"`
	WeekNumber    string    `json:"week_number"`
}

func databasePatchNoteToPatchNote(dbPatchNote database.GetBlogPostsByCreatedAtRow) PatchNote {
	return PatchNote{
		ID:          dbPatchNote.ID,
		CreatedAt:   dbPatchNote.CreatedAt,
		UpdatedAt:   dbPatchNote.UpdatedAt,
		Title:       dbPatchNote.Title,
		Description: dbPatchNote.Description,
		UserID:      dbPatchNote.UserID,
		WeekNumber:  dbPatchNote.WeekNumber,
	}
}

func databasePatchNotesToPatchNotes(dbPatchNotes []database.GetBlogPostsByCreatedAtRow) []PatchNote {
	patchNotes := []PatchNote{}
	for _, dbPatchNote := range dbPatchNotes {
		patchNotes = append(patchNotes, databasePatchNoteToPatchNote(dbPatchNote))
	}
	return patchNotes
}
