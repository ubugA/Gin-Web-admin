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

