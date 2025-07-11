package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	requests = make(map[string]int)
	mu       sync.Mutex
)

// func RateLimiter() gin.HandlerFunc {
// 	go cleanup()

// 	return func(c *gin.Context) {
// 		userID, _ := c.Get("user_id")
// 		if userID == nil {
// 			c.Next()
// 			return
// 		}

// 		mu.Lock()
// 		defer mu.Unlock()

// 		id := userID.(string)
// 		requests[id]++

// 		if requests[id] > 10 {
// 			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
// 			return
// 		}

// 		c.Next()
// 	}
// }

// func cleanup() {
// 	for {
// 		time.Sleep(time.Minute)
// 		mu.Lock()
// 		requests = make(map[string]int)
// 		mu.Unlock()
// 	}
// }

func RateLimiter() gin.HandlerFunc {
	go cleanup()

	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		if userID == nil {
			c.Next()
			return
		}

		mu.Lock()
		defer mu.Unlock()

		id := userID.(string)
		requests[id]++
		count := requests[id]

		if count > 10 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			return
		}

		// Optional log
		c.Writer.Header().Set("X-RateLimit-Remaining", strconv.Itoa(10-count))
		c.Next()
	}
}

func cleanup() {
	for {
		time.Sleep(10 * time.Second) // reduce for testing
		mu.Lock()
		requests = make(map[string]int)
		mu.Unlock()
	}
}
