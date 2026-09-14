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
		Where("time_line_id = ?", timeID).
		Find(&result).
		Error
	if err != nil {
		panic(err)
	}

	return result
}

func PostLikes(timeID uint, userID uint) *model.TimeLine {
	var m *model.TimeLine

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Debug().
			Exec("insert into time_likes (time_line_id,owner_id) values (?,?)", timeID, userID).
			Error
		if err != nil {
			return err
		}

		err = tx.
			Debug().
			Table("time_lines").
			Where("id = ?", timeID).
			Update("like_count", gorm.Expr("like_count + ?", 1)).
			Error
		if err != nil {
			return err
		}

		return tx.Debug().
			Table("time_lines").
			Where("id = ?", timeID).
			First(&m).Error
	})
	if err != nil {
		panic(err)
	}

	return m
}

func DeleteLikes(timeID uint, userID uint) *model.TimeLine {
	var m *model.TimeLine

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Exec(`delete from time_likes where id = ? and owner_id = ?`, timeID, userID).
			Error
		if err != nil {
			panic(err)
		}
		err = tx.Model(&model.TimeLine{}).
			Where("id = ?", timeID).
			Update("like_count", gorm.Expr("like_count - ?", 1)).
			Error
		if err != nil {
			panic(err)
		}

		return tx.Debug().
			Where("id = ?", timeID).
			First(&m).Error
	})
	if err != nil {
		panic(err)
	}

	return m
}
