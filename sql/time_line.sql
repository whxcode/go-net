/**
* 时间线功能
* */

create table time_lines (
  id bigint not null comment "时间线ID" primary key auto_increment,
  owner_id int not null comment "时间线发布者id",
  elements json not null comment "时间线内容",
  like_count int default 0 comment "点赞数",
  status tinyint default 0 comment "时间线状态，0-正常，1-删除,2-被举报中,3-举报成功,4-举报失败",
  created_at datetime default current_timestamp comment "创建时间",
  updated_at datetime default current_timestamp on update current_timestamp comment "更新时间",
  -- 点赞数不能小于 0
  constraint chk_like_count check (like_count >= 0),
  constraint fk_time_line_owner_id
  foreign key (owner_id)
  references  users(id)
  on delete cascade,

  -- 使用用户 id 作为索引
  INDEX idx_owner_id (owner_id)
    
) ENGINE=InnoDB default charset=utf8mb4 comment "时间线表;记录用户的时间线信息";

-- 
create table time_comments (
  id bigint not null comment "评论id" primary key auto_increment,
  time_line_id bigint not null comment "时间线id",
  owner_id int not null comment "评论者id",
  elements json not null comment "时间线内容",
  status tinyint default 0 comment "评论状态，0-正常，1-删除",
  created_at datetime default current_timestamp comment "创建时间",
  updated_at datetime default current_timestamp on update current_timestamp comment "更新时间",

  index idx_time_line_id (time_line_id),
  index idx_owner_id (owner_id),

  constraint fk_comments_time_line
  foreign key (time_line_id)
  references time_lines(id)
  on delete cascade,

  constraint fk_commenets_owner
  foreign key (owner_id)
  references users(id),

  constraint chk_comments_status check (status in (0, 1))

) ENGINE=InnoDB default charset=utf8mb4 comment "时间线表;记录用户的时间线信息";


--
create table time_likes (
  id bigint not null comment "点赞id" primary key auto_increment,
  time_line_id bigint not null comment "时间线id",
  owner_id int not null comment "点赞者id",
  created_at datetime default current_timestamp comment "创建时间",

  index idx_time_line_id (time_line_id),
  index idx_owner_id (owner_id),

  constraint fk_likes_time_line
  foreign key (time_line_id)
  references time_lines(id)
  on delete cascade,

  constraint fk_likes_owner
  foreign key (owner_id)
  references users(id)

) ENGINE=InnoDB default charset=utf8mb4 comment "时间线;点赞表";
