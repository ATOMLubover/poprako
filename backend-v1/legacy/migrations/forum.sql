-- auto-generated definition
create table forum
(
    id       int auto_increment
        primary key,
    username varchar(60)            not null,
    msg      varchar(200)           not null,
    hide     varchar(2) default '1' not null
)
    engine = MyISAM
    collate = utf8_unicode_ci;