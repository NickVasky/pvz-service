package repo

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
)

var _psq = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func TestSubquery(t *testing.T) {
	sql := _psq.
		Select(
			"r.id",
			"r.date_time",
			"r.pvz_id",
			"s.name").
		From("receptions r").
		Join("statuses s ON r.status_id = s.id").
		Where(sq.And{
			sq.Eq{"r.pvz_id": 11},
			sq.Eq{"s.id": 1},
		})

	t.Log(sql.ToSql())
}
