package models

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
	"recipebook/helpers"
)

type Users struct {
	GormModel
	Username string    `gorm:"unique;not null;uniqueIndex" json:"username" form:"username" valid:"required~Your username is required"`
	Email    string    `gorm:"unique;not null;uniqueIndex" json:"email" form:"email" valid:"required~Your email is required"`
	Password string    `gorm:"not null" json:"password" form:"password" valid:"required~Your password is required,min_string_length(6)~Your password should be at least 6 characters long."`
	Bio      string    `gorm:"not null" json:"bio" form:"bio" valid:"required~Your bio is required."`
	Recipes  []Recipes `gorm:"foreignKey:UserID" json:"recipes"`
}

func (u *Users) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(u)

	if errCreate != nil {
		err = errCreate
		return
	}

	u.Password = helpers.HasHPass(u.Password)
	err = nil
	return
}
