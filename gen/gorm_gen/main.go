package main

import (
	inf "auth/internal/infrastructure"

	"github.com/joho/godotenv"
	"gorm.io/gen"
)

func main() {
	godotenv.Load(".env")

	inf.Init()

	generate((*inf.DB)(inf.GetDB()), gen.Config{
		OutPath:      "internal/infrastructure/query",
		ModelPkgPath: "internal/infrastructure/model",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})
}

func generate(db *inf.DB, conf gen.Config) {
	g := gen.NewGenerator(conf)

	g.UseDB(db.Client)

	all := g.GenerateAllTable()
	g.ApplyBasic(all...)

	g.Execute()
}
