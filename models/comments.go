package models

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Comments struct {
	GormModel
	Message  string `gorm:"not null" json:"message" form:"message" valid:"required~Your comment message is required"`
	UserID   uint   `gorm:"not null" json:"user_id" form:"user_id" valid:"required~User Id of your comment is required"`
	RecipeID uint   `gorm:"not null" json:"recipe_id" form:"recipe_id" valid:"required~Recipe Id of your comment is required"`
	Users    *Users
	Recipes  *Recipes
}

func (p *Comments) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(p)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}

func (p *Comments) BeforeUpdate(tx *gorm.DB) (err error) {
	_, errUpdate := govalidator.ValidateStruct(p)

	if errUpdate != nil {
		err = errUpdate
		return
	}

	err = nil
	return
}
