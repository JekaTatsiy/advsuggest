package advsuggest_test

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	Repository "github.com/JekaTatsiy/advsuggest/advsuggest"
	"github.com/lib/pq"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"gorm.io/driver/postgres"
	gorm "gorm.io/gorm"
)

func TestAdvProductRepository(t *testing.T) {
	format.MaxLength = 0

	RegisterFailHandler(Fail)
	RunSpecs(t, "AdvProductRepository")
}

func BuildSuggest(ID int) *Repository.Item {
	return &Repository.Item{
		ID:        ID,
		LinkURL:   "link" + strconv.Itoa(ID),
		Title:     "title" + strconv.Itoa(ID),
		Queries:   []string{"q1", "q2"},
		Active:    true,
		UpdateAT:  time.Now(),
		CreatedAT: time.Now(),
	}
}

func SuggestArray(suggest *Repository.Item) []*Repository.Item {
	return append(make([]*Repository.Item, 0), suggest)
}
func CreateSuggests(start, count int) ([]*Repository.Item, int) {
	prod := make([]*Repository.Item, 0)
	lastID := 1
	for i := start; i < start+count; i++ {
		prod = append(prod, BuildSuggest(i))
	}
	return prod, lastID + 1
}
func ToStringRow(item *Repository.Item) []string {
	layout := "2016-02-02T15:04:05.000Z"

	id := strconv.Itoa(item.ID)
	UpdateAT := item.UpdateAT.Format(layout)
	CreatedAT := item.UpdateAT.Format(layout)
	active := strconv.FormatBool(item.Active)
	queries := fmt.Sprintf("{%s}", strings.Join(item.Queries, ","))
	return []string{id, item.LinkURL, item.Title, queries, active, UpdateAT, CreatedAT}
}

func ToMockRows(items []*Repository.Item) *sqlmock.Rows {
	rows := sqlmock.
		NewRows([]string{"id", "link_url", "title", "queries", "active", "updated_at", "created_at"})

	for _, x := range items {
		rows.AddRow(x.ID, x.LinkURL, x.Title, pq.Array(x.Queries), x.Active, x.UpdateAT, x.CreatedAT)
	}
	return rows
}

var _ = Describe("AdvProductRepository", func() {
	var repository Repository.Repository
	var mock sqlmock.Sqlmock // mock
	var ctx context.Context

	BeforeEach(func() {
		var dbmock *sql.DB
		var err error

		dbmock, mock, err = sqlmock.New() // mock sql.DB
		Expect(err).ShouldNot(HaveOccurred())

		gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: dbmock})) // open gorm db
		Expect(err).ShouldNot(HaveOccurred())

		repository = Repository.New(gdb)
	})
	AfterEach(func() {
		err := mock.ExpectationsWereMet() // make sure all expectations were met
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("Public functions", func() {
		When("add new suggest", func() {
			It("Success", func() {
				next := 1
				suggests, next := CreateSuggests(next, 3)

				err := repository.Add(ctx, suggests, false)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("add new suggest", func() {
			It("Success", func() {
				next := 4
				suggests, next := CreateSuggests(next, 3)

				mock.ExpectBegin()
				mock.ExpectQuery("INSERT").WillReturnRows(sqlmock.NewRows(nil))
				mock.ExpectCommit()
				
				mock.ExpectBegin()
				mock.ExpectExec("DELETE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT").WillReturnRows(sqlmock.NewRows(nil))
				mock.ExpectCommit()

				err := repository.Add(ctx, suggests, false)
				Expect(err).ShouldNot(HaveOccurred())

				err = repository.Add(ctx, suggests, true)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		When("get list suggest", func() {
			It("Success", func() {
				next := 1
				suggests, next := CreateSuggests(next, 10)

				mock.ExpectBegin()
				mock.ExpectQuery("INSERT").WillReturnRows(sqlmock.NewRows(nil))
				mock.ExpectCommit()

				mock.ExpectQuery("SELECT").WillReturnRows(ToMockRows(suggests))

				err := repository.Add(ctx, suggests, false)
				Expect(err).ShouldNot(HaveOccurred())


				iter, err := repository.GetListAdvSuggest(ctx)
				Expect(err).ShouldNot(HaveOccurred())


				item := Repository.Item{}
				cnt := 0
				for iter.Next(&item) {
					cnt++
					Expect(item.ID).Should(Equal(cnt))
				}
				Expect(cnt).Should(Equal(10))

			})
		})
		When("suggests by id", func() {
			It("Success", func() {
				next := 1
				suggests1, next := CreateSuggests(next, 3)
				suggests2, next := CreateSuggests(next, 3)
				suggests3, next := CreateSuggests(next, 3)

				mock.ExpectBegin()
				mock.ExpectQuery("INSERT").WillReturnRows(sqlmock.NewRows(nil))
				mock.ExpectCommit()

				mock.ExpectBegin()
				mock.ExpectQuery("INSERT").WillReturnRows(sqlmock.NewRows(nil))
				mock.ExpectCommit()

				mock.ExpectBegin()
				mock.ExpectQuery("INSERT").WillReturnRows(sqlmock.NewRows(nil))
				mock.ExpectCommit()

				mock.ExpectQuery("SELECT").WillReturnRows(ToMockRows(
					[]*Repository.Item{suggests1[1], suggests2[1], suggests3[1]}))

				var err error
				err = repository.Add(ctx, suggests1, false)
				Expect(err).ShouldNot(HaveOccurred())
				err = repository.Add(ctx, suggests2, false)
				Expect(err).ShouldNot(HaveOccurred())
				err = repository.Add(ctx, suggests3, false)
				Expect(err).ShouldNot(HaveOccurred())

				items, err := repository.GetAdvSuggestByIDs(ctx, []int{2, 4, 6})
				Expect(err).Should(BeNil())
				Expect(len(items)).Should(Equal(3))

				err = repository.ChangeStateAdvSuggestByID(ctx, false, 1)
				Expect(err).Should(HaveOccurred())

			})
		})
	})
})
