package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"mailcat/internal/config"
	"mailcat/internal/database"
	"mailcat/internal/router"
)

func main() {
	// 加载配置
	configPath := "config/config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 确保数据库目录存在
	dbDir := filepath.Dir(cfg.Database.Path)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	// 初始化数据库
	db, err := database.NewDB(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// 设置路由
	r := router.SetupRouter(db, cfg.API.AuthToken, cfg.Admin.Password)

	// 启动服务器
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting MailCat server on %s", addr)
	log.Printf("Health check: http://%s/health", addr)
	log.Printf("Admin panel: http://%s/admin/login", addr)
	log.Printf("API endpoints:")
	log.Printf("  POST /api/v1/emails - Receive email")
	log.Printf("  GET  /api/v1/emails - List emails")
	log.Printf("  GET  /api/v1/emails/:id - Get email by ID")
	log.Printf("Admin endpoints:")
	log.Printf("  GET  /admin/login - Admin login page")
	log.Printf("  GET  /admin/dashboard - Admin dashboard")
	
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}