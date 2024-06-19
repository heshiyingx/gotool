create table app
(
    id            bigint auto_increment comment '主键ID',
    name          varchar(255)                       null comment 'app名称',
    state         tinyint  default 0                 not null comment '状态',
    app_key       bigint                             not null comment 'appKey',
    app_secret    varchar(100)                       not null comment 'appSecret',
    app_login_url varchar(200)                       not null comment 'app_login_url',
    login_policy  tinyint  default 0                 not null comment 'login_policy',
    created_at    datetime default CURRENT_TIMESTAMP null comment '创建时间，创建时自动插入当前时间',
    updated_at    datetime default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP comment '更新时间，创建时自动插入当前时间，更新时自动更新为当前时间',
    primary key id (id),
    constraint idx_appkey unique (app_key)
) comment 'app表';