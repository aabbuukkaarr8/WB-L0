package orders

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) Get(c *gin.Context) {
	orderUID := c.Param("order_uid")
	ms, err := h.srv.Get(c.Request.Context(), orderUID)
	if err != nil {
		logrus.WithError(err).Errorf("[Get By ID]Order Not Found")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	m := Model{}
	m.FillFromService(ms, &ms.Delivery, &ms.Payment, &ms.Items)
	c.JSON(http.StatusOK, m)
}
