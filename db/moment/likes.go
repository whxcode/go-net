package dbMoment

import (
	"fmt"

	"go-net/db"
	"go-net/model"

	"gorm.io/gorm"
)

type likesDB struct{}

var LikesDB = &likesDB{}

func (m *likesDB) MomentLikes(momentID uint) []*model.MomentLikeResponse {
	var result []*model.MomentLikeResponse

	RawSQL := `SELECT m.*, 
	(case when u.nickname is null then u.username else u.nickname end) as nickname,
	u.avatar
	FROM moments_likes m
	LEFT JOIN users u ON u.id = m.user_id
	WHERE m.moment_id = ? order by created_at desc`

	res := db.DB.Debug().Raw(RawSQL, momentID).
		Find(&result)

	if res.Error != nil {
		panic(res.Error)
	}

	return result
}

func (m *likesDB) PostMomentLike(userID model.UserID, momentID uint) *model.MomentLike {
	MomentLike := &model.MomentLike{
		MomentID: momentID,
		UserID:   uint(userID),
	}

	err := db.DB.Debug().Transaction(func(tx *gorm.DB) error {
		err := tx.Debug().Table("moments").
			Where("id = ?", momentID).
			Update("like_count", gorm.Expr("like_count + ?", 1)).
			Error
		if err != nil {
			return err
		}

		err = tx.Debug().Create(MomentLike).Error
		return err
	})
	if err != nil {
		panic(err)
	}

	return MomentLike
}

func (m *likesDB) DeleteMomentLike(userID model.UserID, momentID uint) *model.MomentLike {
	var like *model.MomentLike
	err := db.DB.Debug().Transaction(func(tx *gorm.DB) error {
		result := tx.Debug().Where("moment_id = ? AND user_id = ?", momentID, userID).Delete(&like)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("no like record found for user %d and moment %d", userID, momentID)
		}

		// 2. ✅ 直接更新并返回最新值
		return tx.Model(&model.Moment{}).
			Where("id = ?", momentID).
			Update("like_count", gorm.Expr("like_count - ?", 1)).
			Error // ✅ 链式调用，先更新再查询
	})
	if err != nil {
		panic(err)
	}

	return like
}
