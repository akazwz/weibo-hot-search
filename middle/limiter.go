package middle

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"
	"time"
)

func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	bucket := ratelimit.NewBucket(fillInterval, cap)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"code": 5030,
				"msg":  "系统繁忙",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
