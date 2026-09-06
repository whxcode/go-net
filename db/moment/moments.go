package dbMoment

import (
	"fmt"

	"gorm.io/gorm"

	"go-net/db"
	"go-net/model"
)

type momentDB struct{}

var MomentDB = &momentDB{}

/**
* 查询朋友圈列表
*
* */
func (m *momentDB) GetMoments(userID model.UserID, limit, offset int) (result *[]model.Moment) {
	err := db.DB.Debug().Where("owner_id = ?", userID).Limit(limit).Offset(offset).Find(&result).Error
	if err != nil {
		panic(err)
	}

	if result == nil {
		result = &[]model.Moment{}
	}

	return
}

func (m *momentDB) AddMoment(moment *model.Moment) *model.Moment {
	mts := &model.Moment{
		Elements: moment.Elements,
		Visible:  moment.Visible,
		OwnerID:  moment.OwnerID,
	}

	err := db.DB.Debug().Create(mts).Error
	if err != nil {
		panic(err)
	}

	return mts
}

/**
* 1、找到该记录对 like_count + 1
* 2、将用户记录到 moments_likes 表中
*
* */
func (m *momentDB) MomentLike(userID model.UserID, momentID uint) *model.Moment {
	var moment *model.Moment

	err := db.DB.Debug().Transaction(func(tx *gorm.DB) error {
		err := tx.Debug().Where("id = ?", momentID).First(&moment).Error
		if err != nil {
			return err
		}

		moment.LikeCount += 1

		err = tx.Debug().Save(moment).Error
		if err != nil {
			return err
		}

		MomentLike := &model.MomentLike{
			MomentID: momentID,
			UserID:   uint(userID),
		}

		err = tx.Debug().Create(MomentLike).Error
		return err
	})
	if err != nil {
		panic(err)
	}

	return moment
}

/**
* 找到该用户是否点赞该朋友圈
*
* */
func (m *momentDB) MomentUnLike(userID model.UserID, momentID uint) *model.Moment {
	var moment *model.Moment
	err := db.DB.Debug().Transaction(func(tx *gorm.DB) error {
		var like *model.MomentLike
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
			First(&moment).Error // ✅ 链式调用，先更新再查询
	})
	if err != nil {
		panic(err)
	}

	return moment
}

/**
* 找到该用户是否点赞该朋友圈
*
* */
func (m *momentDB) DeleteMoment(momentID uint) {
	err := db.DB.Debug().Transaction(func(tx *gorm.DB) error {
		var moment *model.Moment
		result := tx.Debug().Where("id = ?", momentID).Delete(&moment)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("no moment record found for moment %d", momentID)
		}

		return tx.Model(&model.MomentLike{}).
			Where("moment_id = ?", momentID).
			Delete(nil).
			Error // ✅ 链式调用，先更新再查询
	})
	if err != nil {
		panic(err)
	}

	return
}

/*
*
// 查看和某个人的屏蔽设置。传 targetId。
*
*
*/
func (m *momentDB) MomentPrivacyTargetID(userID model.UserID, targetID model.UserID) *model.MomentPrivacy {
	var reuslt *model.MomentPrivacy
	res := db.DB.Debug().Where("user_id = ? AND target_id = ?", userID, targetID).First(&reuslt)

	if res.RowsAffected == 0 {
		panic("no privacy record found")
	}

	if res.Error != nil {
		panic(res.Error)
	}

	return reuslt
}

func (m *momentDB) SetMomentPrivacy(privacy *model.MomentPrivacy) *model.MomentPrivacy {
	var result *model.MomentPrivacy = &model.MomentPrivacy{}
	result.UserID = privacy.UserID
	result.TargetID = privacy.TargetID
	result.HideTheir = privacy.HideTheir
	result.HideMine = privacy.HideMine

	res := db.DB.Debug().Create(&result)

	if res.Error != nil {
		panic(res.Error)
	}

	return result
}

func (m *momentDB) MomentLikes(momentID uint) []*model.MomentLikeResponse {
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

func (m *momentDB) MomentComments(momentID uint) []*model.MomentCommentsResponse {
	var result []*model.MomentCommentsResponse

	RawSQL := `SELECT m.*, 
	(case when u.nickname is null then u.username else u.nickname end) as nickname,
	u.avatar
	FROM moments_comments m
	LEFT JOIN users u ON u.id = m.user_id
	WHERE m.moment_id = ? order by created_at desc`

	res := db.DB.Debug().Raw(RawSQL, momentID).
		Find(&result)

	if res.Error != nil {
		panic(res.Error)
	}

	return result
}
