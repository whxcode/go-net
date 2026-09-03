package model

import (
	"fmt"

	"gorm.io/gorm"
)

type momentDB struct{}

/**
* 查询朋友圈列表
*
* */
func (m *momentDB) GetMoments(userID UserID, limit, offset int) (result *[]Moment) {
	err := DB.Debug().Where("owner_id = ?", userID).Limit(limit).Offset(offset).Find(&result).Error
	if err != nil {
		panic(err)
	}

	if result == nil {
		result = &[]Moment{}
	}

	return
}

func (m *momentDB) AddMoment(moment *Moment) *Moment {
	mts := &Moment{
		Elements: moment.Elements,
		Visible:  moment.Visible,
		OwnerID:  moment.OwnerID,
	}

	err := DB.Debug().Create(mts).Error
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
func (m *momentDB) MomentLike(userID UserID, momentID uint) *Moment {
	var moment *Moment

	err := DB.Debug().Transaction(func(tx *gorm.DB) error {
		err := tx.Debug().Where("id = ?", momentID).First(&moment).Error
		if err != nil {
			return err
		}

		moment.LikeCount += 1

		err = tx.Debug().Save(moment).Error
		if err != nil {
			return err
		}

		MomentLike := &MomentLike{
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
func (m *momentDB) MomentUnLike(userID UserID, momentID uint) *Moment {
	var moment *Moment
	err := DB.Debug().Transaction(func(tx *gorm.DB) error {
		var like *MomentLike
		result := tx.Debug().Where("moment_id = ? AND user_id = ?", momentID, userID).Delete(&like)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("no like record found for user %d and moment %d", userID, momentID)
		}

		// 2. ✅ 直接更新并返回最新值
		return tx.Model(&Moment{}).
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
	err := DB.Debug().Transaction(func(tx *gorm.DB) error {
		var moment *Moment
		result := tx.Debug().Where("id = ?", momentID).Delete(&moment)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("no moment record found for moment %d", momentID)
		}

		return tx.Model(&MomentLike{}).
			Where("moment_id = ?", momentID).
			Delete(nil).
			Error // ✅ 链式调用，先更新再查询
	})
	if err != nil {
		panic(err)
	}

	return
}

var MomentDB = &momentDB{}
