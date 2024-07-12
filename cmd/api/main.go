package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var apiToken = os.Getenv("API_TOKEN")

const (
	brokerURL   = "http://localhost:8080"
	compilerURL = "http://compiler:8000"
	cmsURL      = "http://cms:8000"
)

func main() {
	proxy1 := createReverseProxy(compilerURL)
	proxy2 := createReverseProxy(cmsURL)

	r := gin.Default()

	r.Use(authenticate())
	r.Use(rateLimit())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Broker is running",
		})
	})

	r.Any("/compiler/*path", proxy1)
	r.Any("/cms/*path", proxy2)

	log.Printf("API Gateway listening on %s", brokerURL)
	log.Fatal(r.Run(":8080"))
}

func createReverseProxy(targetURL string) func(*gin.Context) {
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("Error parsing target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	return func(c *gin.Context) {
		log.Printf("Request received for %s", c.Request.URL.Path)

		if strings.HasPrefix(c.Request.URL.Path, "/compiler") {
			c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/compiler")
		} else if strings.HasPrefix(c.Request.URL.Path, "/cms") {
			c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/cms")
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Authorization header missing",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Invalid Authorization header format",
			})
			return
		}

		receivedAPIToken := parts[1]
		if receivedAPIToken != apiToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Unauthorized",
			})
			return
		}

		c.Next()
	}
}

func rateLimit() gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(1*time.Minute/50), 50)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}
		c.Next()
	}
}
