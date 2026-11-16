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

	// Проверка пульса
	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			http.Error(w, "БД недоступна", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// auth
	router.Post("/api/auth/send-code", authHandler.SendCode)
	router.Post("/api/auth/register", authHandler.Register)

	// task
	router.Post("/tasks", taskHandler.CreateTask)
	router.Get("/tasks/{id}", taskHandler.GetTaskID)
	router.Get("/tasks", taskHandler.GetList)
	router.Patch("/tasks/{id}/status", taskHandler.UpdateStatus)
	router.Delete("/tasks/{id}", taskHandler.DeleteTask)

	// schedule
	router.Post("/schedule", scheduleHandler.AddSlot)
	router.Get("/schedule/today", scheduleHandler.GetScheduleToday)
	router.Put("/schedule/{id}", scheduleHandler.UpdateTaskSchedule)
	router.Post("/schedule/import", scheduleHandler.ImportSchedule)

	// change
	router.Post("/changes", changeHandler.SaveChanges)
	router.Get("/changes/pending", changeHandler.GetPendingByUserID)
	router.Get("/changes/{id}", changeHandler.GetChangeByID)
	router.Patch("/changes/{id}/status", changeHandler.UpdateStatusChanges)
	router.Post("/changes/import-text", changeHandler.ImportReplacements)

	log.Println("Subly-v2 сервер запущен на порту :8080...")
	http.ListenAndServe(cfg.HTTP, router)
}
