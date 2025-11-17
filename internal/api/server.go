package api

import (
	"gohome/internal/core"
	"gohome/internal/devices"
	"gohome/internal/state"
	"gohome/utils"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  state.Store
	Router *gin.Engine
}

func NewServer(store state.Store) *Server {

	return &Server{store: store, Router: gin.Default()}
}

func (s *Server) Start() {

	s.Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	s.Router.GET("/adapters", func(c *gin.Context) {
		c.JSON(200, utils.Map(utils.SortedValuesByKey(s.store.GetAdapters()), func(e core.Adapter) map[string]any {
			return e.ToJson()
		}))
	})

	s.Router.GET("/parsers", func(c *gin.Context) {
		c.JSON(200, utils.Map(utils.SortedValuesByKey(s.store.GetParsers()), func(e core.Parser) map[string]any {
			return map[string]any{"name": e.Name()}
		}))
	})

	s.Router.GET("/devices", func(c *gin.Context) {
		c.JSON(200, utils.Map(utils.SortedValuesByKey(s.store.GetDevices()), func(e core.Device) map[string]any {
			return e.ToJson()
		}))
	})

	s.Router.GET("/device-types", func(c *gin.Context) {
		c.JSON(200, core.DeviceTypes)
	})

	s.Router.POST("/devices", func(ctx *gin.Context) {
		var req struct {
			Type       string `json:"type"`
			Addr       string `json:"addr"`
			Name       string `json:"name"`
			ParserType string `json:"parser_type"`
		}
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		parser, ok := s.store.GetParsers()[req.ParserType]
		if !ok {
			ctx.JSON(400, gin.H{"error": "parser not found"})
			return
		}

		if _, ok := s.store.GetDevices()[req.Addr]; ok {
			ctx.JSON(400, gin.H{"error": "device already exists"})
			return
		}
		var d core.Device
		switch core.DeviceType(req.Type) {
		case core.TemperatureSensorType:
			d = devices.NewTemperatureSensor(req.Addr, req.Name, parser)

		case core.HumiditySensorType:
			d = devices.NewHumiditySensor(req.Addr, req.Name, parser)

		default:
			ctx.JSON(400, gin.H{"error": "unknown device type"})
			return
		}

		s.store.SaveDevice(d)

		ctx.JSON(201, gin.H{"status": "created"})
	})

	s.Router.POST("/devices/:addr/link", func(ctx *gin.Context) {
		addr := ctx.Param("addr")

		var body struct {
			AdapterId string `json:"adapter_id"`
		}
		if err := ctx.BindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		device, ok := s.store.GetDevice(addr)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "device does not exists"})
			return
		}

		adapter, ok := s.store.GetAdapter(body.AdapterId)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "adapter does not exists"})
			return
		}

		if err := adapter.RegisterDevice(device); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to link adapter to device"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "linked"})
	})

	s.Router.POST("/devices/:addr/unlink", func(ctx *gin.Context) {
		addr := ctx.Param("addr")

		var body struct {
			AdapterId string `json:"adapter_id"`
		}
		if err := ctx.BindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		device, ok := s.store.GetDevice(addr)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "device does not exists"})
			return
		}

		adapter, ok := s.store.GetAdapter(body.AdapterId)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "adapter does not exists"})
			return
		}

		if err := adapter.UnregisterDevice(device.Addr()); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "unlinked"})
	})

	s.Router.POST("/adapter/:id/start", func(ctx *gin.Context) {
		id := ctx.Param("id")

		adapter, ok := s.store.GetAdapter(id)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "adapter does not exists"})
			return
		}

		if err := adapter.Start(ctx); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "started"})
	})

	s.Router.POST("/adapter/:id/stop", func(ctx *gin.Context) {
		id := ctx.Param("id")

		adapter, ok := s.store.GetAdapter(id)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "adapter does not exists"})
			return
		}

		if err := adapter.Stop(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "started"})
	})

	s.Router.POST("/adapter/:id/restart", func(ctx *gin.Context) {
		id := ctx.Param("id")

		adapter, ok := s.store.GetAdapter(id)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "adapter does not exists"})
			return
		}

		if err := adapter.Stop(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := adapter.Start(ctx); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "started"})
	})

	s.Router.POST("/devices/:addr/start", func(ctx *gin.Context) {
		addr := ctx.Param("addr")

		device, ok := s.store.GetDevice(addr)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "device does not exists"})
			return
		}

		if err := device.Start(ctx); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "started"})
	})

	s.Router.Run(":8080")
}
