package scanner

import (
	"context"
	"gohome/internal/core"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpScanner struct {
	router *gin.Engine
	core.BaseScanner
}

func NewHttpScanner(addresses []string, router *gin.Engine) *HttpScanner {
	s := &HttpScanner{
		router: router,
	}
	s.BaseScanner = *core.NewBaseScanner(addresses, s)
	return s
}

func (s *HttpScanner) OnStart(ctx context.Context) error {
	// Une route par adresse connue
	for _, addr := range s.Addresses {
		localAddr := addr // capture
		s.router.POST("/scan/"+localAddr, func(c *gin.Context) {
			var body struct {
				Value float64 `json:"value"`
				Type  string  `json:"type"`
			}

			if err := c.BindJSON(&body); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
				return
			}

			adv := &core.BasicAdvertisment{
				Value: body.Value,
				Type:  body.Type,
			}

			select {
			case s.BaseScanner.ChOut[localAddr] <- adv:
			default:
				// queue full → drop
			}

			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	// Shutdown when ctx closes
	go func() {
		<-ctx.Done()
		_ = s.OnStop()
	}()

	return nil
}

func (s *HttpScanner) OnStop() error {
	// rien à arrêter pour le HTTP
	return nil
}
