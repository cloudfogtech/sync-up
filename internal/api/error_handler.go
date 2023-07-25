package api

import (
	"errors"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/cloudfogtech/sync-up/internal/exp"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func commonErrorHandler() gin.HandlerFunc {
	errorType := "api"
	return customRecovery(func(c *gin.Context, recovered interface{}) {
		if customError := recovered.(*exp.CommonError); customError != nil {
			msg := customError.Message
			if msg == "" {
				if customError.Err != nil {
					msg = customError.Err.Error()
				} else {
					msg = "Unknown error"
				}
			}
			c.JSON(customError.Code, gin.H{
				"message": msg,
			})
			if customError.IsSystemError {
				log.Errorf("[%s] code=%d, msg=%s", errorType, customError.Code, msg)
			}
			return
		}
		if baseErr := recovered.(error); baseErr != nil {
			log.Errorf("[%s] code=%d, msg=%s, %v", errorType, http.StatusInternalServerError, baseErr.Error(), baseErr)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": baseErr.Error(),
			})
			return
		}
		log.Errorf("[%s] Unhandled error: %v", errorType, recovered)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unhandled error",
		})
	})
}

func customRecovery(recoveryFunc gin.RecoveryFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne, &se) {
						seStr := strings.ToLower(se.Error())
						if strings.Contains(seStr, "broken pipe") ||
							strings.Contains(seStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				if brokenPipe {
					log.Errorf("broken pipe: %v", err)
					c.AbortWithStatus(http.StatusInternalServerError)
					c.Abort()
				} else {
					recoveryFunc(c, err)
				}
			}
		}()
		c.Next()
	}
}
