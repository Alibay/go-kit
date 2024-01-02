package pg

import (
	"time"

	"github.com/Alibay/go-kit"
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GormDto specifies base attrs for GORM dto
type GormDto struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *gorm.DeletedAt
}

type TotalCount struct {
	TotalCount int `gorm:"column:total"`
}

// StringToNull transforms empty string to nil string, so that gorm stores it as NULL
func StringToNull(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// NullToString transforms NULL to empty string
func NullToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// GetEmptyJson returns empty jsonb
func GetEmptyJson() (*pgtype.JSONB, error) {
	emptyJson := pgtype.JSONB{}
	err := emptyJson.Scan("{}")
	if err != nil {
		return nil, ErrPgEmptyJsonb(err)
	}
	return &emptyJson, nil
}

// MapToJsonb converts map to jsonb
func MapToJsonb[T comparable, K any](payload map[T]K) (*pgtype.JSONB, error) {
	if payload == nil {
		return GetEmptyJson()
	}
	jsonb := pgtype.JSONB{}
	err := jsonb.Set(payload)
	if err != nil {
		return nil, ErrPgSetJsonb(err)
	}
	return &jsonb, nil
}

// ToJsonb converts arbitrary object to jsonb
func ToJsonb[T any](payload *T) (*pgtype.JSONB, error) {
	if payload == nil {
		return GetEmptyJson()
	}
	jsonb := pgtype.JSONB{}
	err := jsonb.Set(payload)
	if err != nil {
		return nil, ErrPgSetJsonb(err)
	}
	return &jsonb, nil
}

func FromJsonb[T any](j *pgtype.JSONB) (*T, error) {
	if j == nil {
		return nil, nil
	}
	var v T
	err := j.AssignTo(&v)
	if err != nil {
		return nil, ErrPgGetJsonb(err)
	}
	return &v, nil
}

const (
	PageSizeMaxLimit = 100
	PageSizeDefault  = 20
)

func PagingLimit(rqLimit int) int {
	if rqLimit <= 0 {
		return PageSizeDefault
	}
	if rqLimit > PageSizeMaxLimit {
		return PageSizeMaxLimit
	}
	return rqLimit
}

func Paging(rq kit.PagingRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// apply sort
		if len(rq.SortBy) == 0 {
			rq.SortBy = []*kit.SortRequest{{
				Field: "updated_at",
				Desc:  true,
			}}
		}
		for _, srt := range rq.SortBy {
			db = db.Order(clause.OrderByColumn{Column: clause.Column{Name: srt.Field}, Desc: srt.Desc})
		}

		// apply paging
		if rq.Index < 0 {
			rq.Index = 0
		}
		offset := (rq.Index - 1) * rq.Size
		if offset < 0 {
			offset = 0
		}
		return db.Limit(PagingLimit(rq.Size)).Offset(offset)
	}
}

func OrderByUpdatedAt(desc bool) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(clause.OrderByColumn{Column: clause.Column{Name: "updated_at"}, Desc: desc})
	}
}

func Merge() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Clauses(clause.OnConflict{UpdateAll: true})
	}
}

func Update() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Omit("created_at")
	}
}

func Single() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(1)
	}
}
