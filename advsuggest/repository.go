package advsuggest

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type Item struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	LinkURL   string
	Title     string
	Queries   pq.StringArray `gorm:"type:varchar(255)[]"`
	Active    bool           `gorm:"index" json:"active"`
	UpdateAT  time.Time
	CreatedAT time.Time
}

func (i *Item) TableName() string {
	return "adv_suggest"
}

type Repository interface {
	GetAdvSuggestByIDs(ctx context.Context, ids []int) ([]*Item, error)
	GetListAdvSuggest(ctx context.Context) (Iterator, error)
	Add(ctx context.Context, data []*Item, clean bool) error
	ChangeStateAdvSuggestByID(ctx context.Context, active bool, id int) error
}

type RepositoryImpl struct {
	db *gorm.DB
}

func (r RepositoryImpl) GetAdvSuggestByIDs(ctx context.Context, ids []int) ([]*Item, error) {

	Items := make([]*Item, 0)
	result := r.db.Table("adv_suggest").
		Where(struct {
			active bool
			id     []int
		}{active: true, id: ids}).
		Find(&Items).Debug()

	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	var updatedAT sql.NullTime
	for _, x := range Items {
		x.UpdateAT = updatedAT.Time
	}

	return Items, nil
}

func (r RepositoryImpl) GetListAdvSuggest(ctx context.Context) (Iterator, error) {
	rows, e := r.db.WithContext(ctx).Table("adv_suggest").Rows()
	if e != nil {
		return nil, errors.WithStack(e)
	}
	return newADVSuggestIterator(rows), nil
}

func (r RepositoryImpl) Add(ctx context.Context, data []*Item, clean bool) error {
	tx := r.db.WithContext(ctx)

	if clean {
		result := tx.Table("adv_suggest").Where("true").Delete(&Item{})
		if result.Error != nil {
			return errors.WithStack(result.Error)
		}
	}

	result := tx.Table("adv_suggest").Create(&data)

	if result.Error != nil {
		return errors.WithStack(result.Error)
	}
	return nil

}

func (r RepositoryImpl) ChangeStateAdvSuggestByID(ctx context.Context, active bool, id int) error {
	result := r.db.WithContext(ctx).Table("adv_suggest").Where(id).Update("active = ?, updated_at = NOW()", active)
	return errors.WithStack(result.Error)
}

func New(db *gorm.DB) Repository {
	repo := &RepositoryImpl{
		db: db,
	}

	return repo
}
