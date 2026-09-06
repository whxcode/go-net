package dbMoment

import (
	"go-net/db"
	"go-net/model"
)

type privacyDB struct{}

var PrivacyDB = &privacyDB{}

func (m *privacyDB) MomentPrivacyTargetID(userID model.UserID, targetID model.UserID) *model.MomentPrivacy {
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

func (m *privacyDB) SetMomentPrivacy(privacy *model.MomentPrivacy) *model.MomentPrivacy {
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
