package main

import (
	"context"
	"log"
	"net/http"

	"subly-v2/internal/config"
	"subly-v2/internal/handler"
	"subly-v2/internal/pkg/mailer"
	"subly-v2/internal/repository/postgres"
	"subly-v2/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.LoadConfig()
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.ConnPgrs)
	if err != nil {
		log.Fatalf("Ошибка создания пула соединений: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("БД недоступна: %v", err)
	}
	log.Println("Успешное подключение к PostgreSQL!")

	taskRepo := postgres.NewTaskRepository(pool)
	scheduleRepo := postgres.NewScheduleRepository(pool)
	changeRepo := postgres.NewChangeRepository(pool)
	authRepo := postgres.NewAuthRepository(pool)

	mailSender := mailer.NewSMTPMailer(cfg.SMTPPort, cfg.SMTPHost, cfg.SMTPUser, cfg.SMTPPass)

	changeService := &service.ChangeService{Repo: changeRepo, RepoSchedule: scheduleRepo}
	taskService := &service.TaskService{Repo: taskRepo}
	scheduleService := &service.ScheduleService{Repo: scheduleRepo}
	authService := service.NewAuthService(authRepo, mailSender, cfg.JWTSecret)

	changeHandler := &handler.ChangeHandler{ChangeService: changeService}
	taskHandler := &handler.TaskHandler{ServiceTask: taskService}
	scheduleHandler := &handler.ScheduleHandler{ServiceSchedule: scheduleService}
	authHandler := &handler.AuthHandler{AuthService: authService}

	router := chi.NewRouter()
	router.Use(middleware.Logger, middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Accept", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			http.Error(w, "БД недоступна", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// auth (публичные эндпоинты)
	router.Post("/api/auth/send-code", authHandler.SendCode)
	router.Post("/api/auth/register", authHandler.Register)
	router.Post("/api/auth/login", authHandler.Login)

	// защищенные роуты
	router.Group(func(r chi.Router) {
		r.Use(handler.AuthMiddleware(cfg.JWTSecret))

		// task
		r.Post("/tasks", taskHandler.CreateTask)
		r.Get("/tasks/{id}", taskHandler.GetTaskID)
		r.Get("/tasks", taskHandler.GetList)
		r.Patch("/tasks/{id}/status", taskHandler.UpdateStatus)
		r.Delete("/tasks/{id}", taskHandler.DeleteTask)

		// schedule
		r.Post("/schedule", scheduleHandler.AddSlot)
		r.Get("/schedule/today", scheduleHandler.GetScheduleToday)
		r.Put("/schedule/{id}", scheduleHandler.UpdateTaskSchedule)
		r.Post("/schedule/import", scheduleHandler.ImportSchedule)

		// change
		r.Post("/changes", changeHandler.SaveChanges)
		r.Get("/changes/pending", changeHandler.GetPendingByUserID)
		r.Get("/changes/{id}", changeHandler.GetChangeByID)
		r.Patch("/changes/{id}/status", changeHandler.UpdateStatusChanges)
		r.Post("/changes/import-text", changeHandler.ImportReplacements)
	})

	log.Println("Subly-v2 сервер запущен на порту :8080...")
	http.ListenAndServe(cfg.HTTP, router)
}
