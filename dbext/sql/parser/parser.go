package parser

import (
	"fmt"
	"github.com/heshiyingx/gotool/dbext/ddl-parser/parser"
	"github.com/heshiyingx/gotool/dbext/sql/converter"
	"github.com/heshiyingx/gotool/util"
	collection2 "github.com/heshiyingx/gotool/util/collection"
	"github.com/heshiyingx/gotool/util/console"
	stringx "github.com/heshiyingx/gotool/util/stringext"
	"path/filepath"
	"strings"
)

const timeImport = "time.Time"

type (
	// Table describes a mysql table
	Table struct {
		Name        stringx.String
		Db          stringx.String
		PrimaryKey  Primary
		UniqueIndex map[string][]*Field
		Fields      []*Field
		ContainsPQ  bool
	}

	// Primary describes a primary key
	Primary struct {
		Fields        []*Field
		AutoIncrement bool
	}

	// Field describes a table field
	Field struct {
		NameOriginal    string
		Name            stringx.String
		ThirdPkg        string
		DataType        string
		Comment         string
		SeqInIndex      int
		OrdinalPosition int
		ContainsPQ      bool
	}

	// KeyType types alias of int
	KeyType int
)

func Parse(filename, database string, strict bool) ([]*Table, error) {
	p := parser.NewParser()
	tables, err := p.From(filename)
	if err != nil {
		return nil, err
	}

	nameOriginals := parseNameOriginal(tables)
	indexNameGen := func(column ...string) string {
		return strings.Join(column, "_")
	}

	prefix := filepath.Base(filename)
	var list []*Table
	primaryIsDefined := false
	for indexTable, e := range tables {
		var (
			primaryColumn    string
			primaryColumnSet = collection2.NewSet[string]()
			uniqueKeyMap     = make(map[string][]string)
			// Unused local variable
			// normalKeyMap     = make(map[string][]string)
			columns = e.Columns
		)

		for _, column := range columns {
			if column.Constraint != nil {
				if column.Constraint.Primary {
					primaryIsDefined = true
					primaryColumnSet.Add(column.Name)
				}

				if column.Constraint.Unique {
					indexName := indexNameGen(column.Name, "unique")
					uniqueKeyMap[indexName] = []string{column.Name}
				}

				if column.Constraint.Key {
					indexName := indexNameGen(column.Name, "idx")
					uniqueKeyMap[indexName] = []string{column.Name}
				}
			}
		}

		for _, e := range e.Constraints {
			//if len(e.ColumnPrimaryKey) > 1 {
			//	return nil, fmt.Errorf("%s: unexpected join primary key", prefix)
			//}
			//
			//if len(e.ColumnPrimaryKey) == 1 {
			//	primaryColumn = e.ColumnPrimaryKey[0]
			//	primaryColumnSet.AddStr(e.ColumnPrimaryKey[0])
			//}
			if primaryIsDefined && len(e.ColumnPrimaryKey) > 0 {
				return nil, fmt.Errorf("primaryIsDefined:%v", primaryColumnSet.Elems())
			}
			for _, pk := range e.ColumnPrimaryKey {
				primaryColumnSet.Add(pk)
			}

			if len(e.ColumnUniqueKey) > 0 {
				list := append([]string(nil), e.ColumnUniqueKey...)
				list = append(list, "unique")
				indexName := indexNameGen(list...)
				uniqueKeyMap[indexName] = e.ColumnUniqueKey
			}
		}

		if primaryColumnSet.Count() > 1 {
			fmt.Printf("%s: unexpected join primary key\n", prefix)
			return nil, fmt.Errorf("%s: unexpected join primary key", prefix)

		}

		delete(uniqueKeyMap, indexNameGen(primaryColumn, "idx"))
		delete(uniqueKeyMap, indexNameGen(primaryColumn, "unique"))
		primaryKey, fieldM, err := convertColumns(columns, primaryColumnSet, strict)
		if err != nil {
			return nil, err
		}

		var fields []*Field
		// sort
		for indexColumn, c := range columns {
			field, ok := fieldM[c.Name]
			if ok {
				field.NameOriginal = nameOriginals[indexTable][indexColumn]
				fields = append(fields, field)
			}
		}

		uniqueIndex := make(map[string][]*Field)

		for indexName, each := range uniqueKeyMap {
			for _, columnName := range each {
				// Prevent a crash if there is a unique key constraint with a nil field.
				if fieldM[columnName] == nil {
					return nil, fmt.Errorf("table %s: unique key with error column name[%s]", e.Name, columnName)
				}
				uniqueIndex[indexName] = append(uniqueIndex[indexName], fieldM[columnName])
			}
		}

		checkDuplicateUniqueIndex(uniqueIndex, e.Name)

		list = append(list, &Table{
			Name:        stringx.From(e.Name),
			Db:          stringx.From(database),
			PrimaryKey:  primaryKey,
			UniqueIndex: uniqueIndex,
			Fields:      fields,
		})
	}

	return list, nil
}
func parseNameOriginal(ts []*parser.Table) (nameOriginals [][]string) {
	var columns []string

	for _, t := range ts {
		columns = []string{}
		for _, c := range t.Columns {
			columns = append(columns, c.Name)
		}
		nameOriginals = append(nameOriginals, columns)
	}
	return
}
func checkDuplicateUniqueIndex(uniqueIndex map[string][]*Field, tableName string) {
	log := console.NewColorConsole()
	uniqueSet := collection2.NewSet[string]()
	//uniqueSet := collection.NewSet()
	for k, i := range uniqueIndex {
		var list []string
		for _, e := range i {
			list = append(list, e.Name.Source())
		}

		joinRet := strings.Join(list, ",")
		if uniqueSet.Contains(joinRet) {
			log.Warning("[checkDuplicateUniqueIndex]: table %s: duplicate unique index %s", tableName, joinRet)
			delete(uniqueIndex, k)
			continue
		}

		uniqueSet.Add(joinRet)
	}
}
func convertColumns(columns []*parser.Column, primaryColumn *collection2.SortSet[string], strict bool) (Primary, map[string]*Field, error) {
	var (
		primaryKey Primary
		fieldM     = make(map[string]*Field)
		//log          = console.NewColorConsole()
		primaryNames = primaryColumn.Elems()
	)
	primaryKey = Primary{
		Fields: make([]*Field, primaryColumn.Count()),
	}
	for _, column := range columns {
		if column == nil {
			continue
		}

		var (
			comment       string
			isDefaultNull bool
		)

		if column.Constraint != nil {
			comment = column.Constraint.Comment
			isDefaultNull = !column.Constraint.NotNull
			if !column.Constraint.NotNull && column.Constraint.HasDefaultValue {
				isDefaultNull = false
			}

			for _, primaryName := range primaryNames {
				if column.Name == primaryName {
					isDefaultNull = false
				}
			}

		}

		dataType, thirdPkg, err := converter.ConvertDataType(column.DataType.Type(), isDefaultNull, column.DataType.Unsigned(), strict)
		if err != nil {
			return Primary{}, nil, err
		}

		//if column.Constraint != nil {
		//	if column.Name == primaryColumn {
		//		if !column.Constraint.AutoIncrement && dataType == "int64" {
		//			log.Warning("[convertColumns]: The primary key %q is recommended to add constraint `AUTO_INCREMENT`", column.Name)
		//		}
		//	} else if column.Constraint.NotNull && !column.Constraint.HasDefaultValue {
		//		log.Warning("[convertColumns]: The column %q is recommended to add constraint `DEFAULT`", column.Name)
		//	}
		//}

		var field Field
		field.Name = stringx.From(column.Name)
		field.ThirdPkg = thirdPkg
		field.DataType = dataType
		field.Comment = util.TrimNewLine(comment)

		for i, primaryName := range primaryNames {
			if field.Name.Source() == primaryName {
				//primaryKey = Primary{
				//	Field: field,
				//}
				primaryKey.Fields[i] = &field
				if column.Constraint != nil {
					primaryKey.AutoIncrement = column.Constraint.AutoIncrement
				}
			}
		}

		fieldM[field.Name.Source()] = &field
	}
	return primaryKey, fieldM, nil
}
func (t *Table) ContainsTime() bool {
	for _, item := range t.Fields {
		if item.DataType == timeImport {
			return true
		}
	}
	return false
}