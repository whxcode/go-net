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
		Where("t.time_line_id = ? and t.status = ?", timeID, 0).
		Order("t.created_at desc").
		Find(&result).Error
	if err != nil {
		panic(err)
	}

	return result
}

func PostComments(comment *model.TimeComment) *model.TimeComment {
	err := db.DB.Debug().
		Exec("insert into time_comments (owner_id,time_line_id,elements) values(?,?,?)", comment.OwnerID, comment.TimeLineID, comment.Elements).
		Error
	if err != nil {
		panic(err)
	}

	return comment
}

func PutComment(comment *model.TimeComment) *model.TimeComment {
	err := db.DB.Debug().Table("time_comments").
		Where("id = ?", comment.ID).
		Update("elements", comment.Elements).
		Error
	if err != nil {
		panic(err)
	}

	return nil
}

func DeleteComment(commentID uint) {
	err := db.DB.Table("time_comments").
		Where("id = ?", commentID).
		Update("status", 1).
		Error
	if err != nil {
		panic(err)
	}
}
