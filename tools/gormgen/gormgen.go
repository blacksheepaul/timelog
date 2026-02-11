package main

import (
	"os"

	"github.com/ncruces/go-sqlite3/gormlite"
	"gorm.io/gorm"

	"gorm.io/gen"
)

func main() {
	dbPath := os.Getenv("TIMELOG_GEN_DB_PATH")
	if dbPath == "" {
		panic("TIMELOG_GEN_DB_PATH is not set")
	}

	db, err := gorm.Open(gormlite.Open(dbPath))
	if err != nil {
		panic(err)
	}

	// 创建生成器
	g := gen.NewGenerator(gen.Config{
		OutPath:           "../../model/gen",
		ModelPkgPath:      "github.com/blacksheepaul/timelog/model/gen",
		FieldNullable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		WithUnitTest:      false,
		FieldSignable:     false,
		FieldCoverable:    false,
		Mode:              gen.WithoutContext,
	})

	// 使用数据库
	g.UseDB(db)

	// 生成所有模型
	g.ApplyBasic(
		g.GenerateAllTable()...,
	)

	// 执行生成
	g.Execute()
}
