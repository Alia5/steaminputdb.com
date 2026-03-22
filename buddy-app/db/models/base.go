package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Base struct {
	ID        uuid.UUID `bun:",pk,type:uuid,notnull"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func (b *Base) BeforeAppendModel(_ context.Context, q bun.Query) error {
	if _, ok := q.(*bun.InsertQuery); ok && b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	// required for sqlite...
	if _, ok := q.(*bun.UpdateQuery); ok {
		b.UpdatedAt = time.Now()
	}
	return nil
}
