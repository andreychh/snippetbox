package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type Config struct {
	tcpAddr   string
	staticDir string
}

// #3.1 parseFlags creates Config from command-line flags
func parseFlags() Config {
	var tcpAddr = flag.String("tcpAddr", ":4000", "HTTP network address")
	var staticDir = flag.String("staticDir", "./ui/static/", "Directory for static files")

	flag.Parse()
	return Config{
		tcpAddr:   *tcpAddr,
		staticDir: *staticDir,
	}
}

// App #3.3 Инъекция зависимостей посредством структуры
type App struct {
	logger *slog.Logger
}

func NewApp(logger *slog.Logger) *App {
	return &App{
		logger: logger,
	}
}

func main() {
	var config = parseFlags()

	// #3.2 Создание структурированного логера.
	// Пользовательские регистраторы, созданные slog.New(), безопасны для параллелизма.
	// Уровень логирования по умолчанию - Info
	var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true, // #3.2 Добавляет атрибут source=fileAbsolutePath:LineNumber
		Level:     slog.LevelDebug,
	}))

	var app = NewApp(logger)
	// #3.2 Ключи атрибутов всегда должны быть строками, но значения могут быть любого типа.
	logger.Info("starting server", "tcpAddr", config.tcpAddr)
	// #3.2 slog.Logger.Error() не останавливает выполнение
	// #3.3 app.routes(config) изолирует создание mux
	logger.Error(http.ListenAndServe(config.tcpAddr, app.routes(config)).Error())
	os.Exit(1)
}
