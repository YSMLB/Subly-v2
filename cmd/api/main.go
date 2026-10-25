package main

import (
	"context"
	"log"
	"net/http"

	"subly-v2/internal/handler"
	"subly-v2/internal/repository/postgres"
	"subly-v2/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

const poolStr = "postgres://postgres:postgres@localhost:5432/subly_db?sslmode=disable"



func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, poolStr)
	if err != nil{
		log.Fatalf("Ошибка создания пул соединений: %v", err)
	}
	defer pool.Close()

	//пингация БД(Fail Fast)
	if err := pool.Ping(ctx); err != nil{
		log.Fatalf("БД недоступна: %v", err)
	}
	log.Println("Успешное подключение к PostrgreSQL!")

	taskRepo := postgres.NewTaskRepository(pool)
	scheduleRepo := postgres.NewScheduleRepository(pool)
	changeRepo := postgres.NewChangeRepository(pool)

	changeService := &service.ChangeService{Repo: changeRepo}
	taskService := &service.TaskService{Repo: taskRepo}
	scheduleService := &service.ScheduleService{Repo: scheduleRepo}

	changeHandler := &handler.ChangeHandler{ChangeService: changeService}
	taskHandler := &handler.TaskHandler{ServiceTask: taskService}
	scheduleHandler := &handler.ScheduleHandler{ServiceSchedule: scheduleService}

	router := chi.NewRouter()
	router.Use(middleware.Logger, middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization", "Accept", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	//проверка пульса
	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		//првоерка жива ли БД прям щас
		if err := pool.Ping(r.Context()); err != nil{
			http.Error(w, "БД недоступна", http.StatusServiceUnavailable)
			return
		} 

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	//task
	router.Post("/tasks", taskHandler.CreateTask)
	router.Get("/tasks/{id}", taskHandler.GetTaskID)
	router.Get("/tasks", taskHandler.GetList)
	router.Patch("/tasks/{id}/status", taskHandler.UpdateStatus)
	router.Delete("/tasks/{id}", taskHandler.DeleteTask)

	//schedule
	router.Post("/schedule", scheduleHandler.AddSlot)
	router.Get("/schedule/today", scheduleHandler.GetScheduleToday)
	router.Put("/schedule/{id}", scheduleHandler.UpdateTaskSchedule)

	//change
	router.Post("/changes", changeHandler.SaveChanges)
	router.Get("/changes/pending", changeHandler.GetPendingByUserID)
	router.Get("/changes/{id}", changeHandler.GetChangeByID)
	router.Patch("/changes/{id}/status", changeHandler.UpdateStatusChanges)

	log.Println("Subly-v2 сервер запущен на порту:8080...")
	http.ListenAndServe(":8080", router)

}
