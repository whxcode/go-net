package dbTimeLine

import (
	"go-net/db"
	"go-net/model"

	"gorm.io/gorm"
)

func Likes(timeID uint) []*model.TimeLike {
	var result []*model.TimeLike

	err := db.DB.
		Table("time_likes t").
		Select("t.*,u.nickname,u.avatar").
		Joins(`left join users u on
			t.owner_id = u.id
		`).
		Where("time_id = ?", timeID).
		Find(&result).
		Error
	if err != nil {
		panic(err)
	}

	return result
}

func PostLikes(timeID uint, userID uint) *model.Moment {
	var m *model.Moment

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var result *model.TimeLike = &model.TimeLike{
			TimeLineID: timeID,
			OwnerID:    userID,
		}

		err := tx.
			Create(result).
			Error
		if err != nil {
			return err
		}

		return tx.
			Table("time_likes").
			Update("like_count", gorm.Expr("like_count + ?", 1)).
			Find(&m).
			Error
	})
	if err != nil {
		panic(err)
	}

	return m
}

func DeleteLikes(timeID uint, userID uint) *model.Moment {
	var m *model.Moment

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Raw(`delete from time_likes where id = ? and owner_id`, timeID, userID).
			Scan(nil).
			Error
		if err != nil {
			panic(err)
		}

		return tx.Model(&model.Moment{}).
			Where("id = ?", timeID).
			Update("like_count", gorm.Expr("like_count - ?", 1)).
			Find(&m).
			Error
	})
	if err != nil {
		panic(err)
	}

	return m
}
