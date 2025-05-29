package repo

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

var _psq = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func TestSubquery(t *testing.T) {
	IDs := []uuid.UUID{uuid.New(), uuid.New()}

	sql := _psq.
		Select(
			"r.id",
			"r.date_time",
			"r.pvz_id",
			"s.name").
		From("receptions r").
		Join("statuses s ON r.status_id = s.id").
		Where(sq.Eq{"r.id": IDs})

	t.Log(sql.ToSql())
}
