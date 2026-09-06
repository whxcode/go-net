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
func (m *momentDB) GetMoments(userID model.UserID, friendID model.UserID, limit, offset int) (result []*model.MomentResponse) {
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
	if friendID == 0 {
		db.DB.Exec("set @friend_id = ?", nil)
	} else {
		db.DB.Exec("set @friend_id = ?", friendID)
	}

	sql := `
-- 获取未屏蔽我的好友或我为屏幕好友的 id 
with 
friend_tb  as (
select 
(case when f.user_id = @user_id then f.user_id else f.friend_id end) as user_id,
(case when f.user_id = @user_id then f.friend_id else f.user_id end) as friend_id
from friends f
where f.status = 1 and ( f.user_id = @user_id or f.friend_id = @user_id )
),

friend_ids as (
select 
f.friend_id
,coalesce(mp.hide_their,0) as hide_their,
mp.user_id,
mp.target_id

from friend_tb f
left join moments_privacy mp 
-- 好友设置了不让我看他的朋友圈
on mp.user_id = f.friend_id
where 
 ( coalesce(mp.hide_mine,0) = 0

and not exists (
select id from moments_privacy y 
where 
-- 我屏蔽了该好用
y.user_id = @user_id
and y.target_id =  f.friend_id
and y.hide_their = 1
))
)

-- select * from friend_ids;

select m.id,m.owner_id, 
u.nickname,
u.avatar
from moments m
left join users u on u.id = m.owner_id
where 
 (case when @friend_id is null then True else m.owner_id = @friend_id end) and
(m.owner_id = @user_id or -- 仅自己可见
(case
  when m.visible = 0 then 1 -- 公开
  when m.visible = 1 then -- 好友可见
    m.owner_id in (select friend_ids.friend_id from friend_ids) 
  when m.visible = 2 then
    m.owner_id = @user_id -- 仅自己可见
  when m.visible = 3 then -- 部分好友可见
    m.owner_id in (select friend_ids.friend_id from friend_ids) 
    and EXISTS ( select visible from moments_visible mv where mv.moment_id = m.id and mv.user_id = @user_id and mv.visible = 0)
  else 0
end) )
order by m.created_at desc 
limit 10 offset 0;
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

func (m *momentDB) GetMomentsWithUserID(userID, friendID model.UserID, limit, offset int) (result []*model.MomentResponse) {
	db.DB.Exec("set @user_id = ?", userID)
	db.DB.Exec("set @friend_id = ?", friendID)

	sql := `
	select 
	m.*, 
	u.nickname,
	u.avatar
	from moments m
	left join users u
	on owner_id = u.id
	where m.owner_id = @friend_id
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
