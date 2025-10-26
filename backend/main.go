package main

import (
	"shuttlesync/migrations"
)

func main() {

	// 初始化数据库
	InitDb()

	// 运行数据库迁移（创建 RefreshToken 表）
	migrations.RunMigrations()

	router()

}
