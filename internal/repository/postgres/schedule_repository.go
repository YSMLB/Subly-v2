package postgres

import (
	"context"
	"subly-v2/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ScheduleRepository struct{
	pool *pgxpool.Pool
}

func NewScheduleRepository(pool *pgxpool.Pool) *ScheduleRepository{
	return &ScheduleRepository{pool: pool}
}

func (as *ScheduleRepository)AddSlot(schedule domain.ScheduleSlot) error{
	ctx := context.Background()
	query := "INSERT INTO schedule_slots (user_id, subject, teacher, room, day_of_week, start_time, end_time, type, parity, is_cancelled) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"

	_, err := as.pool.Exec(
		ctx, query, 
		schedule.UserID, schedule.Subject, 
		schedule.Teacher, schedule.Room,
		schedule.DayOfWeek, schedule.StartTime,
		schedule.EndTime, schedule.Type,
		schedule.Parity, schedule.IsCancelled,
	)
	if err != nil{
		return err
	}

	return nil
}