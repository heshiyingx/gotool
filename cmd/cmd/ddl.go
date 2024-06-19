package cmd

import (
	"fmt"
	"github.com/heshiyingx/gotool/dbext/sql/gen"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
	"strconv"
)

var DDlCmd = newDdlCmd()

func newDdlCmd() *cobra.Command {
	return &cobra.Command{
		Use:                    "ddl",
		Short:                  "解析ddl语句，生成db层代码",
		Example:                "",
		ValidArgs:              nil,
		ValidArgsFunction:      nil,
		Args:                   nil,
		ArgAliases:             nil,
		BashCompletionFunction: "",
		Version:                "v0.0.1",
		Run: func(cmd *cobra.Command, args []string) {
			sqlFile := cmd.Flags().Lookup("sql").Value
			dir := cmd.Flags().Lookup("dir").Value
			db := cmd.Flags().Lookup("db").Value
			isCacheString := cmd.Flags().Lookup("cache").Value.String()
			isCache, err := strconv.ParseBool(isCacheString)
			if err != nil {
				log.Println(err)
				return
			}

			generator, err := gen.NewDefaultGenerator(dir.String(), nil)
			if err != nil {
				log.Println(err)
				return
			}
			sqlFilePath, err := filepath.Abs(sqlFile.String())
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Println(sqlFilePath, dir.String(), db.String(), isCache)
			err = generator.StartFromDDL(sqlFilePath, isCache, true, db.String())
			if err != nil {
				log.Println(err)
				return
			}
		},
		RunE: nil,
	}
}

func init() {
	DDlCmd.Flags().String("sql", "", "ddl语句位置")
	DDlCmd.Flags().StringP("dir", "d", ".", "生成代码的存放位置")
	DDlCmd.Flags().String("db", "", "数据库名称")
	DDlCmd.Flags().BoolP("cache", "c", true, "是否启用缓存")
	DDlCmd.MarkFlagRequired("sql")
}
