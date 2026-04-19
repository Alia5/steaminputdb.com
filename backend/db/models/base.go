package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Timestamps struct {
	CreatedAt time.Time `bun:",notnull,default:current_timestamp,skipupdate"`
	UpdatedAt time.Time `bun:",notnull,default:current_timestamp"`
}

func (t *Timestamps) BeforeAppendModel(_ context.Context, q bun.Query) error {
	now := time.Now()
	switch q.(type) {
	case *bun.InsertQuery:
		if t.CreatedAt.IsZero() {
			t.CreatedAt = now
		}
		t.UpdatedAt = now
	case *bun.UpdateQuery:
		t.UpdatedAt = now
	}
	return nil
}

type Base struct {
	ID         uuid.UUID  `bun:",pk,type:uuid,notnull"`
	Timestamps Timestamps `bun:",embed"`
}

func (b *Base) BeforeAppendModel(_ context.Context, q bun.Query) error {
	if _, ok := q.(*bun.InsertQuery); ok && b.ID == uuid.Nil {
		b.ID = uuid.New()
	}

	return nil
}
