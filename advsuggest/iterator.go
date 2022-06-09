package advsuggest

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type advSuggestIterator struct {
	rows gorm.Rows
	err  error
}

type Iterator = (*advSuggestIterator)

func (g *advSuggestIterator) Next(dest interface{}) bool {
	item, ok := dest.(*Item)
	if !ok {
		g.err = fmt.Errorf("unknown target type")
		return false
	}

	if !g.rows.Next() {
		return false
	}

	var updatedAT sql.NullTime

	if err := g.rows.Scan(&item.ID, &item.LinkURL, &item.Title,
		&item.Queries, &item.Active, &updatedAT, &item.CreatedAT); err != nil {
		g.err = err
		return false
	}

	item.UpdateAT = updatedAT.Time

	return true
}

func (g *advSuggestIterator) Err() error {
	return errors.WithStack(g.err)
}

func (g *advSuggestIterator) Release() {
	g.rows.Close()
}

func newADVSuggestIterator(rows gorm.Rows) Iterator {
	return &advSuggestIterator{
		rows: rows,
	}
}
