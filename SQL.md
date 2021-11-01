# SQL

创建数据库  
> CREATE DATABASE IF NOT EXISTS prometheus DEFAULT CHARACTER SET utf8;

创建用户表
>create table if not exists platform_user (  
name varchar(50) not null,  
password varchar(50) not null,  
role varchar(50) not null,  
PRIMARY KEY ( name ) ) ENGINE=InnoDB DEFAULT CHARSET=utf8;  

>初始化数据，创建一个管理员用户：   
insert into platform_user values ('admin', '1de52961c3a824453fd18895eba7a9b0', 'admin');

创建prometheus.yml配置表
>create table if not exists configuration (  
name varchar(50) not null,   
finterval varchar(10) not null,  
rinterval varchar(10) not null,  
aurl varchar(100) not null,  
rpath varchar(100) not null,  
jpath varchar(100) not null,  
timeout varchar(10) not null,  
PRIMARY KEY ( name ) ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

创建组管理(Job)表
>CREATE TABLE `group_config` (  
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
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8  

