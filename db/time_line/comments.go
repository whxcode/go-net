package dbTimeLine

import (
	"go-net/db"
	"go-net/model"
)

func Comments(timeID uint) []*model.TimeComment {
	var result []*model.TimeComment

	err := db.DB.Table("time_comments t").
		Select("t.*,u.nickname,u.avatar").
		Joins("left join users u on t.owner_id = u.id").
		Where("t.time_id = ?", timeID).
		Order("t.created_at desc").
		Find(&result).Error
	if err != nil {
		panic(err)
	}

	return result
}

func PostComments(comment *model.TimeComment) *model.TimeComment {
	m := &model.TimeComment{
		OwnerID:  comment.OwnerID,
		Elements: comment.Elements,
	}

	err := db.DB.Create(m).Error

	if err == nil {
		panic(err)
	}

	return m
}

func PutComment(comment *model.TimeComment) *model.TimeComment {
	return nil
}

func DeleteComment(commentID uint) {
}
