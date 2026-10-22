package postgres

import (
	"context"
	"subly-v2/internal/domain"

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