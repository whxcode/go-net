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
func (m *momentDB) GetMoments(userID model.UserID, limit, offset int) (result []*model.MomentResponse) {
	/*
		err := db.DB.Debug().Table("moments m").
			Select(`m.*,u.nickname,u.avatar`).
			Joins(`left join users u on u.id = m.owner_id`).
			Where("owner_id = ?", userID).
			Limit(limit).
			Offset(offset).
			Find(&result).Error
	*/
	db.DB.Exec("set @user_id = ?", userID)

	sql := `
-- 获取未屏蔽我的好友或我为屏幕好友的 id 
with friend_ids as (
select 
(case when f.user_id = @user_id then f.friend_id else f.user_id end) as friend_id

from friends f
left join moments_privacy mp 
on mp.user_id = ((case when f.user_id = @user_id then f.friend_id else f.user_id end) )

where 
(f.user_id = @user_id or f.friend_id = @user_id)
and (coalesce(mp.hide_mine,0) = 0 and coalesce(mp.hide_their,0) = 0)
)

select 
m.*,
u.nickname,
u.avatar
from moments m
left join users u on u.id = m.owner_id
where 
m.owner_id = @user_id or -- 仅自己可见
(case
  when m.visible = 0 then 1 -- 公开
  when m.visible = 1 then -- 好友可见
    m.owner_id in (select * from friend_ids) 
  when m.visible = 2 then
    m.owner_id = @user_id -- 仅自己可见
  when m.visible = 3 then -- 部分好友可见
    EXISTS ( select visible from moments_visible mv where mv.moment_id = m.id and mv.user_id = @user_id and mv.visible = 0)
  else 0
end) 
order by m.created_at desc 
limit ? offset ?;
	`

	err := db.DB.Debug().
		Raw(sql, limit, offset).
		Scan(&result).
		Error
	if err != nil {
		panic(err)
	}

	if result == nil {
		result = []*model.MomentResponse{}
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
