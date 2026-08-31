
/*
1. GET /api/moments —— 朋友圈列表（首页刷的流，自己+好友的）。前端无参数。
2. GET /api/moments/user/:userId?limit=3 —— 看某个人的朋友圈。前端传 userId（URL里
   ）+ limit（条数，个人主页用3）。
3. POST /api/moments —— 发朋友圈。前端传：text_content（文字）、images（图片数组）
   、video（视频对象，含 url/thumbnail/duration，可选）。
4. POST /api/moments/:id/like —— 点赞。传 id。
5. DELETE /api/moments/:id/like —— 取消点赞。传 id。
6. DELETE /api/moments/:id —— 删自己的朋友圈。传 id。
7. POST /api/moments/:id/comments —— 评论。传 id + content（评论内容）。
8. GET /api/moments/privacy/:targetId —— 查和某个人的屏蔽设置。传 targetId。
9. POST /api/moments/privacy —— 设置屏蔽。传 target_id、hide_their（我不看TA的）、h
   ide_mine（不让TA看我的）。
*/

-- 朋友圈表
create table moments (
  id bigint not null comment "朋友圈id" primary key auto_increment,
  owner_id bigint not null comment "朋友圈发布者id",
  content json not null comment "朋友圈文字内容",
  status tinyint default 0 comment "朋友圈状态，0-正常，1-删除",
  like_count int default 0 comment "点赞数",
  visible int default 0 comment "可见性，0-公开，1-好友可见，2-仅自己可见 3-部分好友可见",
  created_at datetime default current_timestamp comment "创建时间",
  updated_at datetime default current_timestamp on update current_timestamp comment "更新时间"
) ENGINE=InnoDB default charset=utf8mb4 comment "朋友圈表;内容主题";


create table moment_likes (
  id bigint not null comment "点赞id" primary key auto_increment,
  moment_id bigint not null comment "朋友圈id",
  user_id bigint not null comment "点赞者id",
  created_at datetime default current_timestamp comment "创建时间",
  UNIQUE KEY uk_moment_user (moment_id, user_id)
) ENGINE=InnoDB default charset=utf8mb4 comment "朋友圈点赞表;记录谁给谁的朋友圈点赞";

-- 朋友圈表
create table moments_visible (
  id bigint not null comment "朋友圈id" primary key auto_increment,
  moment_id bigint not null comment "朋友圈id",
  user_id bigint not null comment "用户 id",
  visible int default 0 comment "0 该好友可见,1 该好友不可见;需要配合coments.visible = 3 的情况",

  UNIQUE KEY uk_moment_user (moment_id, user_id)
) ENGINE=InnoDB default charset=utf8mb4 comment "该条记录谁可见；谁不可见表";

create table moment_privacy (
 id bigint not null primary key auto_increment comment "隐私设置id",
 user_id bigint not null comment "用户id",
 target_id bigint not null comment "目标用户id",
 hide_their tinyint default 0 comment "我不看TA的朋友圈，0-不屏蔽，1-屏蔽",
 hide_mine tinyint default 0 comment "不让TA看我的朋友圈，0-不屏蔽，1-屏蔽",
 UNIQUE KEY uk_user_target (user_id, user_id)

) ENGINE=InnoDB default charset=utf8mb4 comment "朋友圈隐私表;记录谁屏蔽谁的朋友圈";


-- 朋友圈评论表
create table moments_comments (
  id bigint not null comment "评论id" primary key auto_increment,
  moment_id bigint not null comment "朋友圈id",
  user_id bigint not null comment "评论者id",
  content text not null comment "评论内容",
  status tinyint default 0 comment "评论状态，0-正常，1-删除",
  created_at datetime default current_timestamp comment "创建时间",
  updated_at datetime default current_timestamp on update current_timestamp comment "更新时间"
) ENGINE=InnoDB default charset=utf8mb4 comment "朋友圈评论表";

