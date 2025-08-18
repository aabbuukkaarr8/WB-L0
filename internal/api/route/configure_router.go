package route

import (
	"net/http"

	handlerOrders "L0-arch/internal/api/handler/orders"
	"github.com/gin-gonic/gin"
	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
)

func ConfigureRoutes(r *gin.Engine, orders *handlerOrders.Handler) {
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/order/:order_uid", orders.Get)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.StaticFile("/ui", "web/index.html")
}
