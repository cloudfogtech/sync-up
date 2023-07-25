package api

import (
	"errors"
	"github.com/cloudfogtech/sync-up/internal/common"
	"github.com/cloudfogtech/sync-up/internal/data"
	"github.com/cloudfogtech/sync-up/internal/env"
	"github.com/cloudfogtech/sync-up/internal/sync"
	"github.com/cloudfogtech/sync-up/internal/utils"
	"github.com/gin-gonic/gin"
	gorm_logrus "github.com/onrik/gorm-logrus"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Server struct {
	s      *sync.Manager
	router *gin.Engine
	user   *User
	envs   *env.Manager
	db     *data.DB
}

func NewServer() (*Server, error) {
	sm, err := sync.NewManager()
	if err != nil {
		return nil, err
	}
	s := &Server{
		s:    sm,
		envs: env.NewManager(),
	}
	s.InitUser()
	err = s.InitDB()
	if err != nil {
		return nil, err
	}
	s.InitRouter()
	return s, nil
}

func (s *Server) Serve() (*http.Server, error) {
	envMgr := env.NewManager()
	port, err := strconv.Atoi(envMgr.GetGlobalEnvWithDefault(env.PortTpl, common.DefaultPort))
	if err != nil {
		return nil, err
	}
	log.Infof("[DEBUG] port=%d", port)
	return &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: s.router,
	}, nil
}

func (s *Server) InitUser() {
	pass := s.envs.GetGlobalEnvWithDefault(env.PasswordTpl, common.DefaultPassword)
	encodePassword := utils.EncodePasswordFE(pass, common.DefaultFESalt)
	encodePassword = utils.EncodePasswordBE(encodePassword, common.DefaultBESalt)
	s.user = &User{
		Username: s.envs.GetGlobalEnvWithDefault(env.UsernameTpl, common.DefaultUsername),
		Password: encodePassword,
	}
}

func (s *Server) CheckUser(username string, feEncodePassword string) (*User, error) {
	if s.user.Username != username || s.user.Password != utils.EncodePasswordBE(feEncodePassword, common.DefaultBESalt) {
		return nil, errors.New("'username' or 'password' is not correct")
	}
	return s.user, nil
}

func (s *Server) InitDB() error {
	dbType := s.envs.GetGlobalEnvWithDefault(env.DBTypeTpl, common.DefaultDBType)
	dbDsn := s.envs.GetGlobalEnvWithDefault(env.DBDsnTpl, common.DefaultDBDsnSqlite)
	dbInstance, err := data.GetDatabase(dbType, dbDsn)
	if err != nil {
		return err
	}
	orm, err := gorm.Open(dbInstance, &gorm.Config{
		Logger: gorm_logrus.New(),
	})
	s.db = data.NewDB(orm)
	s.db.Migrate()
	return nil
}
