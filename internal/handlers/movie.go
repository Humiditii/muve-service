package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"movie-api-backend/internal/db"
	"movie-api-backend/internal/models"

	"github.com/gin-gonic/gin"
)

type MovieRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	Genre       string  `json:"genre"`
	Director    string  `json:"director"`
	ReleaseYear int     `json:"release_year"`
	Duration    int     `json:"duration"`
	Rating      float64 `json:"rating"`
}

func UploadMovie(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	// Parse multipart form
	err := c.Request.ParseMultipartForm(32 << 20) // 32MB max memory
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	// Get movie file
	file, header, err := c.Request.FormFile("movie")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movie file is required"})
		return
	}
	defer file.Close()

	// Validate file type
	if !isValidVideoFile(header.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video file format"})
		return
	}

	// Create upload directory if it doesn't exist
	uploadDir := os.Getenv("UPLOAD_PATH")
	if uploadDir == "" {
		uploadDir = "./uploads/movies"
	}
	os.MkdirAll(uploadDir, 0755)

	// Generate unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
	filepath := filepath.Join(uploadDir, filename)

	// Save file
	dst, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Get file size
	fileInfo, _ := dst.Stat()
	fileSize := fileInfo.Size()

	// Parse movie metadata
	title := c.Request.FormValue("title")
	description := c.Request.FormValue("description")
	genre := c.Request.FormValue("genre")
	director := c.Request.FormValue("director")
	releaseYear, _ := strconv.Atoi(c.Request.FormValue("release_year"))
	duration, _ := strconv.Atoi(c.Request.FormValue("duration"))
	rating, _ := strconv.ParseFloat(c.Request.FormValue("rating"), 64)

	// Create movie record
	movie := models.Movie{
		Title:       title,
		Description: description,
		Genre:       genre,
		Director:    director,
		ReleaseYear: releaseYear,
		Duration:    duration,
		Rating:      rating,
		FilePath:    filepath,
		FileSize:    fileSize,
		UploadedBy:  userID.(uint),
	}

	if err := db.GetDB().Create(&movie).Error; err != nil {
		// Clean up file if database save fails
		os.Remove(filepath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save movie"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Movie uploaded successfully",
		"movie":   movie,
	})
}

func GetMovies(c *gin.Context) {
	var movies []models.Movie
	
	// Get query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	genre := c.Query("genre")
	search := c.Query("search")

	// Build query
	query := db.GetDB().Preload("User")
	
	if genre != "" {
		query = query.Where("genre ILIKE ?", "%"+genre+"%")
	}
	
	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Count total records
	var total int64
	query.Model(&models.Movie{}).Count(&total)

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&movies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movies": movies,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func GetMovie(c *gin.Context) {
	movieID := c.Param("id")
	
	var movie models.Movie
	if err := db.GetDB().Preload("User").First(&movie, movieID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"movie": movie})
}

func StreamMovie(c *gin.Context) {
	movieID := c.Param("id")
	userID, _ := c.Get("user_id")
	
	var movie models.Movie
	if err := db.GetDB().First(&movie, movieID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	// Record view history
	viewHistory := models.ViewHistory{
		UserID:    userID.(uint),
		MovieID:   movie.ID,
		WatchedAt: time.Now(),
	}
	db.GetDB().Create(&viewHistory)

	// Check if file exists
	if _, err := os.Stat(movie.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie file not found"})
		return
	}

	// Handle range requests for video streaming
	rangeHeader := c.GetHeader("Range")
	if rangeHeader != "" {
		serveVideoRange(c, movie.FilePath, rangeHeader)
		return
	}

	// Serve entire file
	c.Header("Content-Type", "video/mp4")
	c.Header("Accept-Ranges", "bytes")
	c.File(movie.FilePath)
}

func DeleteMovie(c *gin.Context) {
	movieID := c.Param("id")
	
	var movie models.Movie
	if err := db.GetDB().First(&movie, movieID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	// Delete file
	if err := os.Remove(movie.FilePath); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to delete file: %v\n", err)
	}

	// Delete from database
	if err := db.GetDB().Delete(&movie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted successfully"})
}

func GetViewHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var history []models.ViewHistory
	if err := db.GetDB().Preload("Movie").Where("user_id = ?", userID).Order("watched_at DESC").Find(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch view history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}

func isValidVideoFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".mp4", ".avi", ".mkv", ".mov", ".wmv", ".flv", ".webm"}
	
	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}

func serveVideoRange(c *gin.Context, filePath, rangeHeader string) {
	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
		return
	}

	fileSize := stat.Size()
	
	// Parse range header
	ranges := strings.Split(strings.TrimPrefix(rangeHeader, "bytes="), "-")
	if len(ranges) != 2 {
		c.Status(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	start, err := strconv.ParseInt(ranges[0], 10, 64)
	if err != nil {
		start = 0
	}

	end := fileSize - 1
	if ranges[1] != "" {
		if endVal, err := strconv.ParseInt(ranges[1], 10, 64); err == nil {
			end = endVal
		}
	}

	if start > end || start < 0 || end >= fileSize {
		c.Status(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	// Set response headers
	c.Header("Content-Type", "video/mp4")
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	c.Header("Content-Length", strconv.FormatInt(end-start+1, 10))
	c.Status(http.StatusPartialContent)

	// Seek to start position
	file.Seek(start, io.SeekStart)

	// Copy range to response
	io.CopyN(c.Writer, file, end-start+1)
}
