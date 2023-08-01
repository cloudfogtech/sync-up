package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/cloudfogtech/sync-up/internal/common"
	"github.com/cloudfogtech/sync-up/internal/data"
	"github.com/cloudfogtech/sync-up/internal/env"
	"github.com/cloudfogtech/sync-up/internal/sync"
	"github.com/cloudfogtech/sync-up/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	gorm_logrus "github.com/onrik/gorm-logrus"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Server struct {
	s         *sync.Manager
	router    *gin.Engine
	user      *User
	envs      *env.Manager
	db        *data.DB
	scheduler *gocron.Scheduler
}

func NewServer() (*Server, error) {
	sm, err := sync.NewManager()
	if err != nil {
		return nil, err
	}
	s := &Server{
		s:         sm,
		envs:      env.NewManager(),
		scheduler: gocron.NewScheduler(time.Local),
	}
	debug := s.envs.GetGlobalEnvWithDefault(env.DebugTpl, common.DefaultDebug)
	if debug == "true" {
		gin.SetMode(gin.DebugMode)
		log.SetLevel(log.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
		log.SetLevel(log.InfoLevel)
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
	log.Infof("Start Sync Up service at port: %d", port)
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

func (s *Server) StartScheduler() error {
	for _, syncer := range s.s.GetSyncers() {
		if syncer.Info().Cron == "" {
			continue
		}
		job, err := s.scheduler.Tag("system").Cron("0 0 1 * *").Do(func() {
		})
		if err != nil {
			log.Errorf("init system cron job [%s] error: %v", syncer.Info().Id, err)
			return err
		}
		job, err = s.scheduler.Tag("syncer").CronWithSeconds(syncer.Info().Cron).Do(func() error {
			err := syncer.Sync()
			if err != nil {
				log.Errorf("run syncer cron job [syncerId=%s] error: %v", syncer.Info().Id, err)
				return err
			}
			log.Infof("run syncer cron job [syncerId=%s] success at %s", syncer.Info().Id, time.Now().Format(time.RFC3339))
			return nil
		})
		if err != nil {
			log.Errorf("init syncer cron job [%s] error: %v", syncer.Info().Id, err)
			return err
		}
		job.RegisterEventListeners(gocron.WhenJobReturnsError(func(_ string, err error) {
			log.Errorf("job [%s] run failed, error: %v", syncer.Info().Id, err)
		}))
		log.Debugf("init syncer cron job [%s] success", syncer.Info().Id)
	}
	s.scheduler.StartAsync()
	return nil
}
