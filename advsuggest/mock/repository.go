package mock

import (
	"github.com/JekaTatsiy/advsuggest/advsuggest"
	"context"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

var _ advsuggest.Repository = (*RepositoryMock)(nil)

// Add implements advsuggest.Repository
func (r *RepositoryMock) Add(ctx context.Context, data []*advsuggest.Item, clean bool) error {
	args := r.Called(ctx, data, clean)
	return args.Error(0)
}

// ChangeStateAdvSuggestByID implements advsuggest.Repository
func (r *RepositoryMock) ChangeStateAdvSuggestByID(ctx context.Context, active bool, id int) error {
	args := r.Called(ctx, active, id)
	return args.Error(0)
}

// GetAdvSuggestByIDs implements advsuggest.Repository
func (r *RepositoryMock) GetAdvSuggestByIDs(ctx context.Context, ids []int) ([]*advsuggest.Item, error) {
	args := r.Called(ctx, ids)
	return args.Get(0).([]*advsuggest.Item), args.Error(1)
}

// GetListAdvSuggest implements advsuggest.Repository
func (r *RepositoryMock) GetListAdvSuggest(ctx context.Context) (advsuggest.Iterator, error) {
	args := r.Called(ctx)
	return args.Get(0).(advsuggest.Iterator), args.Error(1)
}
