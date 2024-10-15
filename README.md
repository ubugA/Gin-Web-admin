## 环境配置
```sql
-- 1. MySQL 数据库连接信息，配置到 ./configs/fat_configs.toml 中 --
[mysql.read]
addr = '127.0.0.1:3306'
name = 'gin_api_admin'
pass = '123456789'
user = 'root'

[mysql.write]
addr = '127.0.0.1:3306'
name = 'gin_api_admin'
pass = '123456789'
user = 'root'

-- 2. 创建数据库 --
CREATE DATABASE gin_api_admin DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;

-- 创建管理员表 --
CREATE TABLE `admin` (
     `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
     `username` varchar(32) NOT NULL DEFAULT '' COMMENT '用户名',
     `password` varchar(100) NOT NULL DEFAULT '' COMMENT '密码',
     `nickname` varchar(60) NOT NULL DEFAULT '' COMMENT '昵称',
     `mobile` varchar(20) NOT NULL DEFAULT '' COMMENT '手机号',
     `is_used` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用(1:是 -1:否)',
     `created_user` varchar(60) NOT NULL DEFAULT '' COMMENT '创建人',
     `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
     `updated_user` varchar(60) NOT NULL DEFAULT '' COMMENT '更新人',
     `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
     PRIMARY KEY (`id`),
     UNIQUE KEY `uniq_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='管理员表';
```

## 启动

```
$ go run main.go -env fat  

// -env 表示设置哪个环境，主要是区分使用哪个配置文件，默认为 fat
// -env dev 表示为本地开发环境，使用的配置信息为：configs/dev_configs.toml
// -env fat 表示为测试环境，使用的配置信息为：configs/fat_configs.toml
// -env uat 表示为预上线环境，使用的配置信息为：configs/uat_configs.toml
// -env pro 表示为正式环境，使用的配置信息为：configs/pro_configs.toml
```

## 接口文档

- 接口文档：http://127.0.0.1:9999/swagger/index.html
- 心跳检测：http://127.0.0.1:9999/system/health

