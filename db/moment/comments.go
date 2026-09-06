package dbMoment

import (
	"go-net/db"
	"go-net/model"
)

type commentDB struct{}

var CommentDB = &commentDB{}

func (c *commentDB) GetMomentComments(momentID uint) (result *[]model.MomentCommentsResponse) {
	err := db.DB.Debug().Table("moments_comments m").
		Select(`m.*,
		(case when u.nickname is  null then u.username else u.nickname end) as nickname,
		u.avatar`).
		Joins("left join users u on u.id =  m.user_id").
		Where("moment_id = ?", momentID).Find(&result).Error
	if err != nil {
		panic(err)
	}

	return
}

func (c *commentDB) PostMomentComments(comment *model.MomentComments) *model.MomentComments {
	data := &model.MomentComments{
		MomentID: comment.MomentID,
		UserID:   comment.UserID,
		Elements: comment.Elements,
	}

	err := db.DB.Debug().Create(data).Error
	if err != nil {
		panic(err)
	}

	return data
}

func (c *commentDB) PutMomentComments(comment *model.MomentComments) *model.MomentComments {
	data := &model.MomentComments{
		ID:       comment.ID,
		Elements: comment.Elements,
	}

	err := db.DB.Debug().Updates(data).Error
	if err != nil {
		panic(err)
	}

	return data
}

func (c *commentDB) DeleteMomentComments(commentID *model.MomentComments) *model.MomentComments {
	data := &model.MomentComments{
		ID: commentID.ID,
	}

	err := db.DB.Debug().Delete(data).Error
	if err != nil {
		panic(err)
	}

	return data
}
