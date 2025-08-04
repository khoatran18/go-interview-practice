package main

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	visitors       map[string]*visitor
	mu             sync.Mutex
	defaultLimiter *rate.Limiter
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var rl = NewRateLimiter(100, time.Minute)

// Article represents a blog article
type Article struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message,omitempty"`
	Error     string      `json:"error,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

// In-memory storage
var articles = []Article{
	{ID: 1, Title: "Getting Started with Go", Content: "Go is a programming language...", Author: "John Doe", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 2, Title: "Web Development with Gin", Content: "Gin is a web framework...", Author: "Jane Smith", CreatedAt: time.Now(), UpdatedAt: time.Now()},
}
var nextID = 3

func main() {
	// TODO: Create Gin router without default middleware
	// Use gin.New() instead of gin.Default()
	router := gin.New()

	// TODO: Setup custom middleware in correct order
	// 1. ErrorHandlerMiddleware (first to catch panics)
	// 2. RequestIDMiddleware
	// 3. LoggingMiddleware
	// 4. CORSMiddleware
	// 5. RateLimitMiddleware
	// 6. ContentTypeMiddleware
	router.Use(
		ErrorHandlerMiddleware(),
		RequestIDMiddleware(),
		LoggingMiddleware(),
		CORSMiddleware(),
		RateLimitMiddleware(),
		ContentTypeMiddleware(),
	)

	// TODO: Setup route groups
	// Public routes (no authentication required)
	// Protected routes (require authentication)

	// TODO: Define routes
	// Public: GET /ping, GET /articles, GET /articles/:id
	// Protected: POST /articles, PUT /articles/:id, DELETE /articles/:id, GET /admin/stats

	public := router.Group("/")
	{
		public.GET("/ping", ping)
		public.GET("/articles", getArticles)
		public.GET("/articles/:id", getArticle)
	}

	private := router.Group("/").Use(AuthMiddleware())
	{
		private.POST("/articles", createArticle)
		private.PUT("/articles/:id", updateArticle)
		private.DELETE("/articles/:id", deleteArticle)
		private.GET("/admin/stats", getStats)
	}

	// TODO: Start server on port 8080
	router.Run(":8080")
}

// TODO: Implement middleware functions

// RequestIDMiddleware generates a unique request ID for each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Generate UUID for request ID
		// Use github.com/google/uuid package
		// Store in context as "request_id"
		// Add to response header as "X-Request-ID"
		requestID := c.GetHeader("X-Request-Id")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// LoggingMiddleware logs all requests with timing information
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Capture start time
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		// TODO: Calculate duration and log request
		// Format: [REQUEST_ID] METHOD PATH STATUS DURATION IP USER_AGENT
		duration := time.Since(start)
		log.Printf("[%s] %s %s %d %s %s %s",
			c.GetString("request_id"),
			c.Request.Method, path,
			c.Writer.Status(),
			duration,
			c.ClientIP(),
			c.Request.UserAgent(),
		)
	}
}

// AuthMiddleware validates API keys for protected routes
func AuthMiddleware() gin.HandlerFunc {
	// TODO: Define valid API keys and their roles
	// "admin-key-123" -> "admin"
	// "user-key-456" -> "user"

	return func(c *gin.Context) {
		// TODO: Get API key from X-API-Key header
		// TODO: Validate API key
		// TODO: Set user role in context
		// TODO: Return 401 if invalid or missing

		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API Key required"})
			c.Abort()
			return
		}

		isValidApiKey := func(apiKey string) bool {
			if strings.Contains(apiKey, "admin-key") || strings.Contains(apiKey, "user-key") {
				return true
			}
			return false
		}

		if !isValidApiKey(apiKey) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
			c.Abort()
			return
		}

		role := ""
		switch apiKey {
		case "admin-key-123":
			role = "admin"
		case "user-key-456":
			role = "user"
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
			c.Abort()
			return
		}

		c.Set("role", role)

		c.Next()
	}
}

// CORSMiddleware handles cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Set CORS headers
		// Allow origins: http://localhost:3000, https://myblog.com
		// Allow methods: GET, POST, PUT, DELETE, OPTIONS
		// Allow headers: Content-Type, X-API-Key, X-Request-ID
		origin := c.Request.Header.Get("Origin")

		allowOrigins := map[string]bool{
			"http://localhost:3000": true,
			"https://myblog.com":    true,
		}

		if allowOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, X-API-Key, X-Request-ID")
		c.Header("Access-Control-Allow-Credentials", "true")

		// TODO: Handle preflight OPTIONS requests
		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusNoContent)
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitMiddleware implements rate limiting per IP
func RateLimitMiddleware() gin.HandlerFunc {
	// TODO: Implement rate limiting
	// Limit: 100 requests per IP per minute
	// Use golang.org/x/time/rate package
	// Set headers: X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset
	// Return 429 if rate limit exceeded

	return func(c *gin.Context) {
		limit := rl.getVisitor(c.ClientIP())
		if !limit.Allow() {
			c.Header("X-RateLimit-Limit", "100")
			c.Header("X-RateLimit-Remaining", "0")
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		remaining := int(limit.Tokens())
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Next()
	}
}

func NewRateLimiter(requests int, duration time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors:       map[string]*visitor{},
		defaultLimiter: rate.NewLimiter(rate.Every(duration), requests),
	}

	go rl.cleanupVisitors()
	return rl
}

func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exist := rl.visitors[ip]
	if !exist {
		rl.visitors[ip] = &visitor{rl.defaultLimiter, time.Now()}
		return rl.defaultLimiter
	}
	v.lastSeen = time.Now()
	return v.limiter
}

func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)

		rl.mu.Lock()
		for id, visitor := range rl.visitors {
			if time.Since(visitor.lastSeen) > 3*time.Minute {
				delete(rl.visitors, id)
			}
		}
		rl.mu.Unlock()
		fmt.Println("Cleanup Visitors Successfully!")
	}
}

// ContentTypeMiddleware validates content type for POST/PUT requests
func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Check content type for POST/PUT requests
		// Must be application/json
		// Return 415 if invalid content type
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			if c.ContentType() != "application/json" {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{
					"error": "Content-Type must be application/json",
					"code":  http.StatusUnsupportedMediaType,
				})
			}
		}

		c.Next()
	}
}

// ErrorHandlerMiddleware handles panics and errors
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// TODO: Handle panics gracefully
		// Return consistent error response format
		// Include request ID in response
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success:   false,
			Error:     fmt.Sprintf("%v", recovered),
			RequestID: c.GetString("request_id"),
		})
		c.Abort()
	})
}

// TODO: Implement route handlers

func response(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, APIResponse{
		Success:   true,
		Error:     message,
		Data:      data,
		RequestID: c.GetString("request_id"),
	})
}

// ping handles GET /ping - health check endpoint
func ping(c *gin.Context) {
	// TODO: Return simple pong response with request ID
	response(c, http.StatusOK, "pong", nil)
}

// getArticles handles GET /articles - get all articles with pagination
func getArticles(c *gin.Context) {
	// TODO: Implement pagination (optional)
	// TODO: Return articles in standard format
	response(c, http.StatusOK, "", articles)
}

// getArticle handles GET /articles/:id - get article by ID
func getArticle(c *gin.Context) {
	// TODO: Get article ID from URL parameter
	// TODO: Find article by ID
	// TODO: Return 404 if not found
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		response(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	article, aId := findArticleByID(idInt)
	if aId == -1 {
		response(c, http.StatusNotFound, "Article not found", nil)
		return
	}
	response(c, http.StatusOK, "Article", article)

}

// createArticle handles POST /articles - create new article (protected)
func createArticle(c *gin.Context) {
	// TODO: Parse JSON request body
	// TODO: Validate required fields
	// TODO: Add article to storage
	// TODO: Return created article
	articleData := &Article{}

	if err := c.ShouldBindJSON(&articleData); err != nil {
		response(c, http.StatusBadRequest, "Invalid article body", nil)
		return
	}

	if err := validateArticle(*articleData); err != nil {
		response(c, http.StatusBadRequest, "Invalid article body", nil)
		return
	}

	articleData.ID = nextID
	articleData.CreatedAt = time.Now()
	articleData.UpdatedAt = articleData.CreatedAt
	articles = append(articles, *articleData)
	nextID++
	response(c, http.StatusCreated, "Article", articleData)
}

// updateArticle handles PUT /articles/:id - update article (protected)
func updateArticle(c *gin.Context) {
	// TODO: Get article ID from URL parameter
	// TODO: Parse JSON request body
	// TODO: Find and update article
	// TODO: Return updated article
	articleData := &Article{}

	if err := c.ShouldBindJSON(&articleData); err != nil {
		response(c, http.StatusBadRequest, "Invalid article body", nil)
		return
	}

	if err := validateArticle(*articleData); err != nil {
		response(c, http.StatusBadRequest, "Invalid article body", nil)
		return
	}

	_, index := findArticleByID(articleData.ID)
	if index == -1 {
		response(c, http.StatusOK, "Article not found", nil)
		return
	}

	articles[index] = *articleData
	response(c, http.StatusOK, "Update Article Successfully", articleData)
}

func deleteArticle(c *gin.Context) {
	// TODO: Get article ID from URL parameter
	// TODO: Find and remove article
	// TODO: Return success message
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	_, index := findArticleByID(id)
	if index == -1 {
		response(c, http.StatusNotFound, "Not found", nil)
		return
	}

	articles = slices.Delete(articles, index, index+1)
	response(c, http.StatusOK, "Article deleted", nil)
}

// getStats handles GET /admin/stats - get API usage statistics (admin only)
func getStats(c *gin.Context) {
	// TODO: Check if user role is "admin"
	// TODO: Return mock statistics
	if c.GetString("role") != "admin" {
		response(c, http.StatusForbidden, "Admin role required", nil)
		return
	}

	stats := map[string]interface{}{
		"total_articles": len(articles),
		"total_requests": 0, // Could track this in middleware
		"uptime":         time.Since(time.Now().Add(-24 * time.Hour)),
	}

	// TODO: Return stats in standard format
	response(c, http.StatusOK, "Stats", stats)
}

// Helper functions

// findArticleByID finds an article by ID
func findArticleByID(id int) (*Article, int) {
	// TODO: Implement article lookup
	// Return article pointer and index, or nil and -1 if not found
	for i, a := range articles {
		if a.ID == id {
			return &a, i
		}
	}
	return nil, -1
}

// validateArticle validates article data
func validateArticle(article Article) error {
	// TODO: Implement validation
	// Check required fields: Title, Content, Author
	if strings.TrimSpace(article.Title) == "" {
		return fmt.Errorf("Title required")
	}
	if strings.TrimSpace(article.Author) == "" {
		return fmt.Errorf("Author required")
	}
	if strings.TrimSpace(article.Content) == "" {
		return fmt.Errorf("Content required")
	}
	return nil
}
