-- auto-generated definition
create table allmg
(
    id2          int auto_increment
        primary key,
    id           int          not null,
    works        varchar(200) not null,
    translation  varchar(200) not null,
    proofreading varchar(200) not null,
    repairer     varchar(200) not null,
    producer     varchar(200) not null,
    settime      int          null,
    dir          varchar(20)  not null
)
    comment '20220330同步表单' engine = MyISAM
                               collate = utf8_unicode_ci;