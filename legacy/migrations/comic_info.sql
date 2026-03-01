-- auto-generated definition
create table comic_info
(
    comic_id                int             not null comment '漫画ID'
        primary key,
    comic_name              varchar(255)    not null comment '漫画名称',
    comic_dir               varchar(255)    not null comment '漫画dir',
    platform_category_id    varchar(64)     not null comment '平台分类及ID',
    last_check_success_time bigint unsigned null comment '最后一次检测成功时间',
    cover_image_url         text            null comment '封面图链接',
    last_updated_chapter    varchar(255)    null comment '最后更新话',
    last_updated_time       bigint unsigned null comment '最后更新时间',
    status                  varchar(32)     null comment '状态',
    extra_json              text            null comment 'JSON（code/data-src）'
)
    collate = utf8mb4_unicode_ci;