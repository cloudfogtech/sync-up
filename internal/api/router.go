package api

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/catfishlty/sync-up/internal/check"
	"github.com/catfishlty/sync-up/internal/common"
	"github.com/catfishlty/sync-up/internal/env"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

const identityKey = "username"

func (s *Server) InitRouter() {
	secretKey := s.envs.GetGlobalEnvWithDefault(env.SecretKeyTpl, common.DefaultSecretKey)
	authMiddleware := s.newAuthMiddleware(secretKey, common.JwtRealm)
	r := gin.New()
	apiGorup := r.Group("api", commonErrorHandler())
	authApiGroup := apiGorup.Group("auth")
	{
		authApiGroup.POST("login", authMiddleware.LoginHandler)
		authApiGroup.GET("refresh", authMiddleware.RefreshHandler)
		authApiGroup.GET("logout", authMiddleware.LogoutHandler)
	}
	syncApiGroup := apiGorup.Group("sync", authMiddleware.MiddlewareFunc())
	{
		syncApiGroup.GET("", s.GetSyncList)
		syncApiGroup.GET(":id", s.GetSyncDetail)
		syncApiGroup.POST(":id/run", s.RunSync)
	}
	logApiGroup := apiGorup.Group("log", authMiddleware.MiddlewareFunc())
	{
		logApiGroup.GET("", s.GetLogAll)
		logApiGroup.GET(":id", s.GetLogList)
		logApiGroup.GET(":id/:log-id", s.GetLogDetail)
	}
	s.router = r
}

func (s *Server) newAuthMiddleware(secretKey, jwtRealm string) *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       jwtRealm,
		Key:         []byte(secretKey),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if u, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: u.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				Username: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVal User
			if err := c.ShouldBind(&loginVal); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVal.Username
			password := loginVal.Password
			err := check.UsernamePassword(username, password)
			if err != nil {
				log.Debugf("login failed: %s, %s", username, err.Error())
				return nil, jwt.ErrFailedAuthentication
			}
			user, err := s.CheckUser(username, password)
			if err != nil || user == nil {
				log.Infof("login failed: %s from ip='%s', %s", username, c.Request.RemoteAddr, err.Error())
				return nil, jwt.ErrFailedAuthentication
			}
			log.Debugf("login success: %s", username)
			return user, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if u, ok := data.(*User); ok {
				log.Debugf("auth success: %s", u.Username)
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			log.Debugf("Unauthorized: %d, %s", code, message)
			c.JSON(code, MsgResponse(message))
		},
		TokenLookup:   "header: Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}
	return authMiddleware
}
