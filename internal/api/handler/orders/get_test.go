package orders

import (
    "context"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"

    svc "L0-arch/internal/service/orders"

    "github.com/gin-gonic/gin"
)

type fakeService struct{
    resp *svc.Model
    err error
}

func (f *fakeService) Get(ctx context.Context, orderUID string) (*svc.Model, error) { return f.resp, f.err }

func TestHandler_Get_OK(t *testing.T) {
    gin.SetMode(gin.TestMode)
    m := &svc.Model{ OrderUID: "uid-1" }
    h := NewHandler(&fakeService{resp: m})
    r := gin.New()
    r.GET("/order/:order_uid", h.Get)

    req := httptest.NewRequest(http.MethodGet, "/order/uid-1", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", w.Code)
    }
}

func TestHandler_Get_NotFound(t *testing.T) {
    gin.SetMode(gin.TestMode)
    h := NewHandler(&fakeService{resp: nil, err: fmt.Errorf("not found")})
    r := gin.New()
    r.GET("/order/:order_uid", h.Get)

    req := httptest.NewRequest(http.MethodGet, "/order/does-not-exist", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusNotFound {
        t.Fatalf("expected 404, got %d", w.Code)
    }
}


