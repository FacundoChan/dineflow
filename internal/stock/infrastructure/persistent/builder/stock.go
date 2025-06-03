package builder

import (
	format "github.com/FacundoChan/dineflow/common/format"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Stock struct {
	ID        []int64  `json:"id,omitempty"`
	ProductID []string `json:"product_id,omitempty"`
	Quantity  []int32  `json:"quantity,omitempty"`
	Version   []int64  `json:"version,omitempty"`

	// extend fields
	OrderBy       string `json:"order_by,omitempty"`
	ForUpdateLock bool   `json:"for_update,omitempty"`
}

func (s *Stock) FormatArg() (string, error) {
	return format.MarshalString(s)
}

func NewStock() *Stock {
	return &Stock{}
}

func (s *Stock) Fill(db *gorm.DB) *gorm.DB {
	db = s.fillWhere(db)
	if s.OrderBy != "" {
		db = db.Order(s.OrderBy)
	}
	return db
}

func (s *Stock) fillWhere(db *gorm.DB) *gorm.DB {
	if len(s.ID) > 0 {
		db = db.Where("id in (?)", s.ID)
	}
	if len(s.ProductID) > 0 {
		db = db.Where("product_id in (?)", s.ProductID)
	}
	if len(s.Quantity) > 0 {
		db = s.fillQuantityGreaterEqual(db)
	}
	if len(s.Version) > 0 {
		db = db.Where("id in (?)", s.Version)
	}
	if s.ForUpdateLock {
		db = db.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate})
	}

	return db
}

func (s *Stock) fillQuantityGreaterEqual(db *gorm.DB) *gorm.DB {
	db = db.Where("quantity >= ?", s.Quantity)
	return db
}

func (s *Stock) IDs(v ...int64) *Stock {
	s.ID = v
	return s
}

func (s *Stock) ProductIDs(v ...string) *Stock {
	s.ProductID = v
	return s
}

func (s *Stock) QuantityGreaterEqual(v ...int32) *Stock {
	s.Quantity = v
	return s
}

func (s *Stock) Versions(v ...int64) *Stock {
	s.Version = v
	return s
}

func (s *Stock) Order(v string) *Stock {
	s.OrderBy = v
	return s
}

func (s *Stock) ForUpdate() *Stock {
	s.ForUpdateLock = true
	return s
}
