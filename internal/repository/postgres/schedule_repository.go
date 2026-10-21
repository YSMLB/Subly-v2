package postgres

import (
	"context"
	"errors"
	"subly-v2/internal/domain"

	"github.com/jackc/pgx/v5"
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


func (gst *ScheduleRepository)GetScheduleToday(user_id int, dayOfWeek int) (*domain.ScheduleSlot, error){
	ctx := context.Background()
	query := "SELECT user_id, subject, teacher, room, day_of_week, start_time, end_time, type, parity, is_cancelled FROM schedule_slots WHERE user_id = $1 AND day_of_week = $2"

	var t domain.ScheduleSlot

	err := gst.pool.QueryRow(ctx, query, user_id, dayOfWeek).Scan(
		&t.UserID, &t.Subject, &t.Teacher,
		&t.Room, &t.DayOfWeek, &t.StartTime,
		&t.EndTime, &t.Type, &t.Parity,
		&t.IsCancelled,
	)

	if err != nil{
		if errors.Is(err, pgx.ErrNoRows) == true{
			return nil, nil
		}
		return nil, err
	}
	
	return &t, nil
}