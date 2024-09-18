package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"snippetbox/internal/models"

	"snippetbox/internal/config"

	// #4.4 пустой идентификатор позволяет выполниться функции
	// init() драйвера для регистрации в пакете "database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// App #3.3 Инъекция зависимостей посредством структуры
type App struct {
	logger  *slog.Logger
	storage *models.SnippetModel
}

// #4.4 инкапсулирует подключение к БД
func openDB(dsn string) (*sql.DB, error) {
	// #4.4 sql.Open не создает никаких подключений, все, что она делает,
	// это инициализирует пул для дальнейшего использования.
	// Фактические подключения к базе данных устанавливаются лениво, по мере необходимости.
	var db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}
	// #4.4  db.Ping подключается к бд и проверяет наличие ошибок.
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("pinging db: %w", err)
	}

	return db, nil
}

func main() {
	var cfg, err = config.Load("config.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	// #4.4 Формат строки подключения (dsn): user:password@protocol(host:port)/dbname?param=value
	// parseTime=true указывает драйверу MySQL автоматически парсить типы данных
	// DATETIME и TIMESTAMP в Go-структуры time.Time.
	// sql.DB - это пул множества подключений к БД. Он безопасен для одновременного доступа.
	db, err := openDB(cfg.DB.DSN())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	var app = App{logger: logger, storage: &models.SnippetModel{DB: db}}

	logger.Info("starting server", "tcpAddr", cfg.App.Addr())
	logger.Error(http.ListenAndServe(cfg.App.Addr(), app.routes(*cfg)).Error())
	os.Exit(1)
}
