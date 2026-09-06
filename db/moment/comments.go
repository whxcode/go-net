package dbMoment

import (
	"go-net/db"
	"go-net/model"
)

type commentDB struct{}

var CommentDB = &commentDB{}

func (c *commentDB) GetMomentComments(momentID uint, limit, offset int) (result *[]model.MomentComments) {
	db.DB.Debug().Where("moment_id = ?", momentID).Limit(limit).Offset(offset).Find(&result)
	return nil
}
