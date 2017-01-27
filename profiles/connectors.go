package profiles

import "github.com/jinzhu/gorm"

type ProfileGorm struct {
	DB *gorm.DB
}

func (o *ProfileGorm) Get(id uint) (*Profile, error) {
	obj := &Profile{}
	err := o.DB.First(obj, id).Error

	return obj, err
}

func NewPageManagerGorm(db *gorm.DB) *ProfileGorm {
	return &ProfileGorm{DB: db}
}

