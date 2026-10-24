package postgres

import (
	"context"
	"errors"
	"subly-v2/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChangeRepository struct{
	pool *pgxpool.Pool
}

func NewChangeRepository(pool *pgxpool.Pool) *ChangeRepository{
	return &ChangeRepository{pool: pool}
}

func (c *ChangeRepository) SaveChanges(changes domain.ProposedChange) (*domain.ProposedChange, error){
	ctx := context.Background()
	query := "INSERT INTO proposed_changes (user_id, slot_id, action, status, raw_text, new_start_time, new_end_time, new_room, created_add) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;"

	err := c.pool.QueryRow(ctx, query,
		changes.UserID, changes.SlotID,
		changes.Action, changes.Status,
		changes.RawText, changes.NewStartTime,
		changes.NewEndTime, changes.NewRoom,
		changes.CreatedAdd,
	).Scan(
		&changes.ID,
	)

	if err != nil{
		return nil, err
	}

	return &changes, nil
}

func (c *ChangeRepository) GetPendingByUserID (user_id int) ([]domain.ProposedChange, error){
	ctx := context.Background()
	query := "SELECT id, user_id, slot_id, action, status, raw_text, new_start_time, new_end_time, new_room, created_add FROM proposed_changes WHERE user_id = $1 AND status = 'pending'"

	cSlice := make([]domain.ProposedChange, 0)

	row, err := c.pool.Query(ctx, query, user_id)
	if err != nil{
		return nil, err
	}

	defer row.Close()

	for row.Next(){
		var r domain.ProposedChange

		err := row.Scan(
			&r.ID,
			&r.UserID, &r.SlotID,
			&r.Action, &r.Status,
			&r.RawText, &r.NewStartTime,
			&r.NewEndTime, &r.NewRoom,
			&r.CreatedAdd,
		)

		if err != nil {
			return nil, err
		}

		cSlice = append(cSlice, r)
	}

	if err = row.Err(); err != nil{
		return nil, err
	}

	return cSlice, nil
}

func (c *ChangeRepository) GetChangeByID(id int) (*domain.ProposedChange, error){
	ctx := context.Background()
	query := "SELECT id, user_id, slot_id, action, status, raw_text, new_start_time, new_end_time, new_room, created_add FROM proposed_changes WHERE id = $1;"

	var t domain.ProposedChange

	err := c.pool.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.UserID,
		&t.SlotID, &t.Action,
		&t.Status, &t.RawText,
		&t.NewStartTime, &t.NewEndTime,
		&t.NewRoom,  &t.CreatedAdd,
	)
	if err != nil{
		if errors.Is(err, pgx.ErrNoRows) == true{
			return nil, nil
		}
		
		return nil, err
	}
	return &t, nil
}

func (c *ChangeRepository) UpdateStatusChanges(id int, user_id int, new_status domain.ChangeStatus) (*domain.ProposedChange, error){
	ctx := context.Background()
	query := "UPDATE proposed_changes SET status = $1 WHERE user_id = $2 AND id = $3 RETURNING id, user_id, slot_id, action, status, raw_text, new_start_time, new_end_time, new_room, created_add"

	var t domain.ProposedChange

	err := c.pool.QueryRow(ctx, query, new_status, user_id, id).Scan(
		&t.ID, &t.UserID,
		&t.SlotID, &t.Action,
		&t.Status, &t.RawText,
		&t.NewStartTime, &t.NewEndTime,
		&t.NewRoom, &t.CreatedAdd,
	)

	if err != nil{
		return nil, err
	}

	return &t, nil
}