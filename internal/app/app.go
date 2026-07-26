package app

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ZhuJincheng-git/stride-backend/internal/config"
	"github.com/ZhuJincheng-git/stride-backend/internal/handler"
	"github.com/ZhuJincheng-git/stride-backend/internal/repository"
	"github.com/ZhuJincheng-git/stride-backend/internal/router"
	"github.com/ZhuJincheng-git/stride-backend/internal/service"
	"github.com/ZhuJincheng-git/stride-backend/pkg/jwt"
)

type App struct {
	Cfg    *config.Config
	DB     *gorm.DB
	Engine *gin.Engine
}

func New(cfg *config.Config, db *gorm.DB) *App {
	tokens := jwt.New(cfg.JWTSecret, cfg.JWTExpires(), "stride-backend")

	// --- repositories ---
	userRepo := repository.NewUserRepository(db)

	// --- services ---
	authSvc := service.NewAuthService(userRepo, tokens)

	// --- handlers ---
	h := &router.Handlers{
		Auth: handler.NewAuthHandler(authSvc),
	}

	engine := router.Build(h, tokens)
	return &App{Cfg: cfg, DB: db, Engine: engine}
}
