package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/catfishlty/sync-up/internal/exp"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (s *Server) GetLogAll(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		exp.HandleRequestInvalid(errors.New("invalid page"))
	}
	logs, total, err := s.db.GetAllLogList(page)
	if err != nil {
		log.Errorf("[DB] GetAllLogList error: %v", err)
		exp.HandleRequestInvalid(errors.New("can't get logs"))
	}
	c.JSON(http.StatusOK, ListResponse(logs, total))
}

func (s *Server) GetLogList(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		exp.HandleRequestInvalid(errors.New("invalid page"))
	}
	serviceId := c.Param("id")
	logs, total, err := s.db.GetLogListByServiceId(serviceId, page)
	if err != nil {
		log.Errorf("[DB] GetLogListByServiceId error: %v", err)
		exp.HandleRequestInvalid(errors.New("can't get logs"))
	}
	c.JSON(http.StatusOK, ListResponse(logs, total))
}

func (s *Server) GetLogDetail(c *gin.Context) {
	serviceId := c.Param("id")
	logId := c.Param("log-id")
	logs, err := s.db.GetLog(serviceId, logId)
	if err != nil {
		log.Errorf("[DB] GetLog error: %v", err)
		exp.HandleRequestInvalid(errors.New("can't get logs"))
	}
	c.JSON(http.StatusOK, logs)
}
