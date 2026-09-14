package dbMoment

import (
	"go-net/db"
	"go-net/model"

	"gorm.io/gorm"
)

type privacyDB struct{}

var PrivacyDB = &privacyDB{}

func (m *privacyDB) MomentPrivacyTargetID(userID model.UserID, targetID model.UserID) *model.MomentPrivacy {
	var reuslt *model.MomentPrivacy
	res := db.DB.Debug().Where("user_id = ? AND target_id = ?", userID, targetID).First(&reuslt)

	if res.Error == gorm.ErrRecordNotFound {
		return nil
	}

	return reuslt
}

func (m *privacyDB) SetMomentPrivacy(privacy *model.MomentPrivacy) *model.MomentPrivacy {
	var result *model.MomentPrivacy = &model.MomentPrivacy{}

	if privacy.ID != 0 {
		db.DB.Model(&model.MomentPrivacy{}).
			Where("id = ?", privacy.ID).
			Updates(map[string]interface{}{
				"hide_their": privacy.HideTheir,
				"hide_mine":  privacy.HideMine,
			}).Find(result)
	} else {

		result.UserID = privacy.UserID
		result.TargetID = privacy.TargetID
		result.HideTheir = privacy.HideTheir
		result.HideMine = privacy.HideMine

		db.DB.Debug().Create(&result)
	}

	return result
}
