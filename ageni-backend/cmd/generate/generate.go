package main

import (
	"fmt"
	"github.com/readonme/open-studio/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

// generate code
func main() {
	// specify the output directory (default: "./query")
	// ### if you want to query without context constrain, set mode gen.WithoutContext ###
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("failed get wd:", err)
		return
	}
	conf.Init()
	projectRootDir := findRootDir(currentDir)
	g := gen.NewGenerator(gen.Config{
		OutPath: projectRootDir + "/dal/query",
		/* Mode: gen.WithoutContext|gen.WithDefaultQuery*/
		//if you want the nullable field generation property to be pointer type, set FieldNullable true
		/* FieldNullable: true,*/
		//if you want to assign field which has default value in `Create` API, set FieldCoverable true, reference: https://gorm.io/docs/create.html#Default-Values
		/* FieldCoverable: true,*/
		// if you want generate field with unsigned integer type, set FieldSignable true
		/* FieldSignable: true,*/
		//if you want to generate index tags from database, set FieldWithIndexTag true
		/* FieldWithIndexTag: true,*/
		//if you want to generate type tags from database, set FieldWithTypeTag true
		/* FieldWithTypeTag: true,*/
		//if you need unit tests for query code, set WithUnitTest true
		/* WithUnitTest: true, */
	})

	// reuse the database connection in Project or create a connection here
	// if you want to use GenerateModel/GenerateModelAs, UseDB is necessary or it will panic
	db, _ := gorm.Open(mysql.Open(conf.Conf.StudioDB.DSN))
	g.UseDB(db)

	// apply basic crud api on structs or table models which is specified by table name with function
	// GenerateModel/GenerateModelAs. And generator will generate table models' code when calling Excute.

	g.ApplyBasic(
		g.GenerateModel("bots"),
		g.GenerateModel("conversations"),
		g.GenerateModel("messages"),
		g.GenerateModel("plugins"),
		g.GenerateModel("methods"),
		g.GenerateModel("bot_plugin"),
		g.GenerateModel("bot_drafts"),
		g.GenerateModel("users"),
	)

	// execute the action of code generation
	g.Execute()
}

func findRootDir(currentDir string) string {
	for {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			return currentDir
		}
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}
	return ""
}
