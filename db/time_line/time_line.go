package dbTimeLine

import (
	"go-net/db"
	"go-net/model"
)

/*
* 获取可见的时间线.
*
* */
func GetTimeLines(limit, offset int) (result []*model.TimeLine) {
	err := db.DB.Table("time_lines t").
		Select("t.*,u.nickname,u.avatar").
		Joins(`
			left join users u on
			t.owner_id = u.id
		`).
		Where("status in (0,2)").
		Limit(limit).
		Offset(offset).
		Find(&result).
		Error
	if err != nil {
		panic(err)
	}

	return
}

func PostTimeLines(data *model.TimeLine) (result *model.TimeLine) {
	t := &model.TimeLine{
		OwnerID:  data.OwnerID,
		Elements: data.Elements,
	}

	err := db.DB.
		Create(t).
		Error
	if err != nil {
		panic(err)
	}

	return t
}

func PutTimeLines(data *model.TimeLine) (result *model.TimeLine) {
	t := &model.TimeLine{
		OwnerID:  data.OwnerID,
		Elements: data.Elements,
	}

	err := db.DB.Table("time_lines").
		Update("elements", data.Elements).
		Where("id = ?", data.ID).
		Error
	if err != nil {
		panic(err)
	}

	return t
}

func DeleteTimeLines(id uint, status uint8) (result *model.TimeLine) {
	t := &model.TimeLine{}

	err := db.DB.Table("time_lines").
		Update("status", status).
		Where("id = ?", id).
		Find(t).
		Error
	if err != nil {
		panic(err)
	}

	return t
}
