package main

import (
	"github.com/google/uuid"
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

func databaseBlogPostToBlogPost(dbBlogPost database.BlogPost) BlogPost {
	return BlogPost{
		ID:            dbBlogPost.ID,
		CreatedAt:     dbBlogPost.CreatedAt,
		UpdatedAt:     dbBlogPost.UpdatedAt,
		Title:         dbBlogPost.Title,
		Description:   dbBlogPost.Description,
		ImageFilename: dbBlogPost.ImageFilename,
		UserID:        dbBlogPost.UserID,
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
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	ImageFilename string    `json:"image_filename"`
	UserID        uuid.UUID `json:"user_id"`
}

func databaseBugReportToBugReport(dbBugReport database.BugReport) BugReport {
	return BugReport{
		ID:            dbBugReport.ID,
		CreatedAt:     dbBugReport.CreatedAt,
		UpdatedAt:     dbBugReport.UpdatedAt,
		Title:         dbBugReport.Title,
		Description:   dbBugReport.Description,
		ImageFilename: dbBugReport.ImageFilename,
		UserID:        dbBugReport.UserID,
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
		ID:            dbPatchNote.ID,
		CreatedAt:     dbPatchNote.CreatedAt,
		UpdatedAt:     dbPatchNote.UpdatedAt,
		Title:         dbPatchNote.Title,
		Description:   dbPatchNote.Description,
		ImageFilename: dbPatchNote.ImageFilename,
		UserID:        dbPatchNote.UserID,
		WeekNumber:    dbPatchNote.WeekNumber,
	}
}

func databasePatchNotesToPatchNotes(dbPatchNotes []database.GetBlogPostsByCreatedAtRow) []PatchNote {
	patchNotes := []PatchNote{}
	for _, dbPatchNote := range dbPatchNotes {
		patchNotes = append(patchNotes, databasePatchNoteToPatchNote(dbPatchNote))
	}
	return patchNotes
}
