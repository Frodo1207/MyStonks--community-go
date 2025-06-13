package v1

import (
	"MyStonks-go/internal/service"
	"github.com/gin-gonic/gin"
)

type EventApi struct {
	eventSrv *service.EventSrv
}

func NewEventApi(eventSrv *service.EventSrv) *EventApi {
	return &EventApi{
		eventSrv: eventSrv,
	}
}

func (e *EventApi) GetEvents(c *gin.Context) {
	events, err := e.eventSrv.GetEvents()
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, map[string]interface{}{
		"code": 200,
		"data": events,
	})
}
