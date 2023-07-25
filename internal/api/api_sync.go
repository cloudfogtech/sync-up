package api

import (
	"fmt"
	"net/http"

	"github.com/cloudfogtech/sync-up/internal/exp"
	"github.com/cloudfogtech/sync-up/internal/sync"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetSyncList(c *gin.Context) {
	list := make([]sync.Info, 0)
	for _, syncer := range s.s.GetSyncers() {
		list = append(list, syncer.Info())
	}
	c.JSON(http.StatusOK, ListResponse(list, int64(len(list))))
}

func (s *Server) GetSyncDetail(c *gin.Context) {
	id := c.Param("id")
	for _, syncer := range s.s.GetSyncers() {
		info := syncer.Info()
		if id == info.Id {
			c.JSON(http.StatusOK, info)
			return
		}
	}
	c.AbortWithStatus(http.StatusNotFound)
}

func (s *Server) RunSync(c *gin.Context) {
	id := c.Param("id")
	for _, syncer := range s.s.GetSyncers() {
		info := syncer.Info()
		if id == info.Id {
			err := syncer.Sync()
			if err != nil {
				panic(&exp.CommonError{
					Code:    http.StatusBadRequest,
					Message: fmt.Sprintf("syncer[%s] run error", info.Id),
					Err:     err,
				})
			}
			c.JSON(http.StatusOK, info)
			return
		}
	}
	c.AbortWithStatus(http.StatusNotFound)
}
