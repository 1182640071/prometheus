# SQL

    目前只对postgresql进行了适配测试，mysql的支持后续会添加

## postgresql:
    创建数据库：
    CREATE DATABASE prometheus encoding='utf8';

    创建用户表：
    create table if not exists platform_user (  
    name varchar(50) not null,  
    password varchar(50) not null,  
    role varchar(50) not null,  
    PRIMARY KEY ( name ) ) ;

    初始化一个用户admin 密码Test123!@#
    insert into platform_user values ('admin', '1de52961c3a824453fd18895eba7a9b0', 'admin');

    创建prometheus.yml配置表:
    create table if not exists configuration (  
    name varchar(50) not null,   
    finterval varchar(10) not null,  
    rinterval varchar(10) not null,  
    aurl varchar(100) not null,  
    rpath varchar(100) not null,  
    jpath varchar(100) not null,  
    timeout varchar(10) not null,  
    PRIMARY KEY ( name ) ) ;

    创建组管理(Job)表
    CREATE TABLE group_config (  
    id SERIAL PRIMARY KEY,  
    name varchar(50) NOT NULL UNIQUE,  
    finterval varchar(10) NOT NULL,  
    scheme varchar(10) NOT NULL,  
    insecure_skip_verify varchar(10) NOT NULL,  
    metrics_path varchar(100) NOT NULL,  
    match_regulation varchar(500) NOT NULL,  
    federalid varchar(50) DEFAULT NULL,  
    honor_labels varchar(10) NOT NULL
    ) ;


    创建主机管理(Target)表
    CREATE TABLE host_config (  
    id varchar(50) PRIMARY KEY,  
    name varchar(50) NOT NULL,  
    ip varchar(50) NOT NULL,  
    port varchar(6) NOT NULL,  
    group_id int NOT NULL,  
    label varchar(500),  
    status varchar(5) NOT NULL,  
    CONSTRAINT unique_matches unique (ip, port, group_id)  
    ) ;

## mysql:

    创建数据库  
    CREATE DATABASE IF NOT EXISTS prometheus DEFAULT CHARACTER SET utf8;

    创建用户表
    create table if not exists platform_user (  
    name varchar(50) not null,  
    password varchar(50) not null,  
    role varchar(50) not null,  
    PRIMARY KEY ( name ) ) ENGINE=InnoDB DEFAULT CHARSET=utf8;  
    
    初始化数据，创建一个管理员用户admin 密码Test123!@#    
    insert into platform_user values ('admin', '1de52961c3a824453fd18895eba7a9b0', 'admin');


    创建prometheus.yml配置表
    create table if not exists configuration (  
    name varchar(50) not null,   
    finterval varchar(10) not null,  
    rinterval varchar(10) not null,  
    aurl varchar(100) not null,  
    rpath varchar(100) not null,  
    jpath varchar(100) not null,  
    timeout varchar(10) not null,  
    PRIMARY KEY ( name ) ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    
    创建组管理(Job)表
    CREATE TABLE `group_config` (  
    `id` int(11) NOT NULL AUTO_INCREMENT,  
    `name` varchar(50) NOT NULL,  
    `finterval` varchar(10) NOT NULL,  
    `scheme` varchar(10) NOT NULL,  
    `insecure_skip_verify` varchar(10) NOT NULL,  
    `metrics_path` varchar(100) NOT NULL,  
    `match_regulation` varchar(500) NOT NULL,  
    `federalid` varchar(50) DEFAULT NULL,  
    `honor_labels` varchar(10) NOT NULL,  
    PRIMARY KEY (`id`),  
    UNIQUE KEY `name` (`name`)  
    ) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8  
    
    
    创建主机管理(Target)表
    CREATE TABLE `host_config` (  
    `id` varchar(50) NOT NULL,  
    `name` varchar(50) NOT NULL,  
    `ip` varchar(50) NOT NULL,  
    `port` varchar(6) NOT NULL,  
    `group_id` int(11) NOT NULL,  
    `label` varchar(500),  
    `status` varchar(5) NOT NULL,  
    PRIMARY KEY (`id`),  
    UNIQUE KEY `target` (`ip`, `port`, `group_id`)  
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8