package gen

import (
	"fmt"
	"github.com/heshiyingx/gotool/dbext/sql/parser"
	"github.com/heshiyingx/gotool/dbext/sql/template"
	"github.com/heshiyingx/gotool/util"
	"github.com/heshiyingx/gotool/util/console"
	"github.com/heshiyingx/gotool/util/format"
	"github.com/heshiyingx/gotool/util/pathext"
	stringx "github.com/heshiyingx/gotool/util/stringext"
	goformat "go/format"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const pwd = "."

type Config struct {
	// NamingFormat is used to define the naming format of the generated file name.
	// just like time formatting, you can specify the formatting style through the
	// two format characters go, and zero. for example: snake format you can
	// define as go_zero, camel case format you can it is defined as goZero,
	// and even split characters can be specified, such as go#zero. in theory,
	// any combination can be used, but the prerequisite must meet the naming conventions
	// of each operating system file name.
	// Note: NamingFormat is based on snake or camel string
	NamingFormat string `yaml:"namingFormat"`
}
type defaultGenerator struct {
	console.Console
	// source string
	dir           string
	pkg           string
	cfg           *Config
	ignoreColumns []string
}
type Option func(generator *defaultGenerator)
type codeTuple struct {
	modelCode       string
	modelCustomCode string
}
type code struct {
	importsCode string
	varsCode    string
	typesCode   string
	newCode     string
	opCode      string
	//insertCode     string
	//findCode       []string
	//updateCode     string
	//deleteCode     string
	cacheExtra     string
	tableName      string
	customizedCode string
}

// Key describes cache key
type Key struct {
	// VarLeft describes the variable of cache key expression which likes cacheUserIdPrefix
	VarLeft string
	// VarRight describes the value of cache key expression which likes "cache:user:id:"
	VarRight string
	// VarExpression describes the cache key expression which likes cacheUserIdPrefix = "cache:user:id:"
	VarExpression string
	// KeyLeft describes the variable of key definition expression which likes userKey
	KeyLeft string
	// KeyRight describes the value of key definition expression which likes fmt.Sprintf("%s%v", cacheUserPrefix, user)
	KeyRight string
	// DataKeyRight describes data key likes fmt.Sprintf("%s%v", cacheUserPrefix, data.User)
	DataKeyRight string
	// KeyExpression describes key expression likes userKey := fmt.Sprintf("%s%v", cacheUserPrefix, user)
	KeyExpression string
	// DataKeyExpression describes data key expression likes userKey := fmt.Sprintf("%s%v", cacheUserPrefix, data.User)
	DataKeyExpression string
	// FieldNameJoin describes the filed slice of table
	FieldNameJoin Join
	// Fields describes the fields of table
	Fields []*parser.Field
}

type Table struct {
	parser.Table
	PrimaryCacheKey        Key
	UniqueCacheKey         []Key
	ContainsUniqueCacheKey bool
	ignoreColumns          []string
}

// NewDefaultGenerator creates an instance for defaultGenerator
func NewDefaultGenerator(dir string, cfg *Config, opt ...Option) (*defaultGenerator, error) {
	if cfg == nil {
		cfg = &Config{NamingFormat: "go-zero"}
	}
	if dir == "" {
		dir = pwd
	}
	dirAbs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	dir = dirAbs
	pkg := util.SafeString(filepath.Base(dirAbs))
	err = pathext.MkdirIfNotExist(dir)
	if err != nil {
		return nil, err
	}

	generator := &defaultGenerator{dir: dir, cfg: cfg, pkg: pkg}
	var optionList []Option
	optionList = append(optionList, newDefaultOption())
	optionList = append(optionList, opt...)
	for _, fn := range optionList {
		fn(generator)
	}

	return generator, nil
}
func newDefaultOption() Option {
	return func(generator *defaultGenerator) {
		generator.Console = console.NewColorConsole()
	}
}
func (g *defaultGenerator) StartFromDDL(filename string, withCache, strict bool, database string) error {
	modeList, err := g.genFromDDL(filename, withCache, strict, database)
	if err != nil {
		return err
	}
	return g.createFile(modeList)
}
func (g *defaultGenerator) genFromDDL(filename string, withCache, strict bool, database string) (map[string]*codeTuple, error) {
	m := make(map[string]*codeTuple)

	tables, err := parser.Parse(filename, database, strict)
	if err != nil {
		return nil, err
	}
	for _, e := range tables {
		gencode, customerCode, err := g.genModel(*e, withCache)
		if err != nil {
			return nil, err
		}
		m[e.Name.Source()] = &codeTuple{
			modelCode:       gencode,
			modelCustomCode: customerCode,
		}
	}
	return m, nil
}
func (g *defaultGenerator) createFile(modelList map[string]*codeTuple) error {
	dirAbs, err := filepath.Abs(g.dir)
	if err != nil {
		return err
	}

	g.dir = dirAbs
	log.Println("代码生成位置:", dirAbs)
	g.pkg = util.SafeString(filepath.Base(dirAbs))
	err = pathext.MkdirIfNotExist(dirAbs)
	if err != nil {
		return err
	}

	for tableName, codes := range modelList {
		tn := stringx.From(tableName)
		modelFilename, err := format.FileNamingFormat(g.cfg.NamingFormat,
			fmt.Sprintf("%s_model", tn.Source()))
		if err != nil {
			return err
		}

		name := util.SafeString(modelFilename) + "_gen.go"
		filename := filepath.Join(dirAbs, name)
		err = os.WriteFile(filename, []byte(codes.modelCode), os.ModePerm)
		if err != nil {
			return err
		}

		name = util.SafeString(modelFilename) + ".go"
		filename = filepath.Join(dirAbs, name)
		if pathext.FileExists(filename) {
			g.Warning("%s already exists, ignored.", name)
			continue
		}
		err = os.WriteFile(filename, []byte(codes.modelCustomCode), os.ModePerm)
		if err != nil {
			return err
		}
	}

	// generate error file
	varFilename, err := format.FileNamingFormat(g.cfg.NamingFormat, "vars")
	if err != nil {
		return err
	}

	filename := filepath.Join(dirAbs, varFilename+".go")
	text := template.VarFile
	//text, err := pathx.LoadTemplate(category, errTemplateFile, template.Error)
	//if err != nil {
	//	return err
	//}

	err = util.With("vars").Parse(text).SaveTo(map[string]any{
		"pkg": g.pkg,
	}, filename, false)
	if err != nil {
		return err
	}

	g.Success("Done.")
	return nil
}

func (g *defaultGenerator) genModel(in parser.Table, withCache bool) (string, string, error) {
	primaryKey, uniqueKey := genCacheKeys(in)
	dbModelOpCode := make([]string, 0)
	var table Table
	table.Table = in
	table.PrimaryCacheKey = primaryKey
	table.UniqueCacheKey = uniqueKey
	table.ContainsUniqueCacheKey = len(uniqueKey) > 0
	table.ignoreColumns = g.ignoreColumns

	importsCode, err := genImports(table, withCache, in.ContainsTime())
	if err != nil {
		return "", "", err
	}
	varsCode, err := genVars(table, withCache)
	if err != nil {
		return "", "", err
	}

	insertCode, insertCodeInterface, err := genInsert(table, withCache)
	if err != nil {
		return "", "", err
	}

	findByPKCode, findByPKInterface, err := genFindPK(table, withCache)
	if err != nil {
		return "", "", err
	}
	updateByPKCode, updateByPKInterface, err := genUpdateByPK(table, withCache)
	if err != nil {
		return "", "", err
	}

	deleteCode, deleteInterface, err := genDeleteByPK(table, withCache)
	if err != nil {
		return "", "", err
	}
	uniqueKeyCode, err := genFindAndUpdateOneByUniqueKey(table, withCache)
	if err != nil {
		return "", "", err
	}
	var list []string
	list = append(list, insertCodeInterface, findByPKInterface, updateByPKInterface, deleteInterface, uniqueKeyCode.findOneInterfaceMethod)
	typesCode, err := genTypes(table, strings.Join(list, "\n"), withCache)
	if err != nil {
		return "", "", err
	}
	defaultModelNewCode, err := genNew(table, withCache)
	if err != nil {
		return "", "", err
	}
	customCode, err := genModelCustom(in, withCache, g.pkg)
	if err != nil {
		return "", "", err
	}
	tableName, err := genTableName(table)
	if err != nil {
		return "", "", err
	}
	dbModelOpCode = append(dbModelOpCode, insertCode, findByPKCode, updateByPKCode, deleteCode, uniqueKeyCode.findOneMethod)
	codeInfo := &code{
		importsCode: importsCode,
		varsCode:    varsCode,
		typesCode:   typesCode,
		newCode:     defaultModelNewCode,
		opCode:      strings.Join(dbModelOpCode, "\n"),
		//insertCode:     insertCode,
		//findCode:       findCode,
		//updateCode:     updateCode,
		//deleteCode:     deleteCode,
		//cacheExtra:     ret.cacheExtra,
		tableName: tableName,
		//customizedCode: customizedCode,
	}
	//fmt.Println(codeInfo)
	genCode, err := g.genGenCode(table, codeInfo)
	if err != nil {
		return "", "", err
	}

	return genCode, customCode, nil
}
func (g *defaultGenerator) genGenCode(table Table, code *code) (string, error) {
	text, err := pathext.LoadTemplate(category, modelGenTemplateFile, template.ModelGen)
	if err != nil {
		return "", err
	}
	t := util.With("model").
		Parse(text).
		GoFmt(true)
	output, err := t.Execute(map[string]any{
		"pkg":         g.pkg,
		"imports":     code.importsCode,
		"vars":        code.varsCode,
		"types":       code.typesCode,
		"new":         code.newCode,
		"opCode":      code.opCode,
		"extraMethod": code.cacheExtra,
		"tableName":   code.tableName,
		"data":        table,
		"customized":  code.customizedCode,
	})
	if err != nil {
		return "", err
	}
	source, err := goformat.Source(output.Bytes())
	if err != nil {
		return "", err
	}
	return string(source), nil
}

func genCacheKeys(table parser.Table) (Key, []Key) {
	var primaryKey Key
	var uniqueKey []Key
	primaryKey = genCacheKey(table.Db, table.Name, table.PrimaryKey.Fields)
	for _, each := range table.UniqueIndex {
		uniqueKey = append(uniqueKey, genCacheKey(table.Db, table.Name, each))
	}
	sort.Slice(uniqueKey, func(i, j int) bool {
		return uniqueKey[i].VarLeft < uniqueKey[j].VarLeft
	})

	return primaryKey, uniqueKey
}

func genCacheKey(db, table stringx.String, in []*parser.Field) Key {
	var (
		varLeftJoin, varRightJoin, fieldNameJoin Join
		varLeft, varRight, varExpression         string

		keyLeftJoin, keyRightJoin, keyRightArgJoin, dataRightJoin         Join
		keyLeft, keyRight, dataKeyRight, keyExpression, dataKeyExpression string
	)

	dbName, tableName := util.SafeString(db.Source()), util.SafeString(table.Source())
	if len(dbName) > 0 {
		varLeftJoin = append(varLeftJoin, "cache", dbName, tableName)
		varRightJoin = append(varRightJoin, "cache", dbName, tableName)
		keyLeftJoin = append(keyLeftJoin, dbName, tableName)
	} else {
		varLeftJoin = append(varLeftJoin, "cache", tableName)
		varRightJoin = append(varRightJoin, "cache", tableName)
		keyLeftJoin = append(keyLeftJoin, tableName)
	}

	for _, each := range in {
		varLeftJoin = append(varLeftJoin, each.Name.Source())
		varRightJoin = append(varRightJoin, each.Name.Source())
		keyLeftJoin = append(keyLeftJoin, each.Name.Source())
		keyRightJoin = append(keyRightJoin, util.EscapeGolangKeyword(stringx.From(each.Name.ToCamel()).Untitle()))
		keyRightArgJoin = append(keyRightArgJoin, "%v")
		dataRightJoin = append(dataRightJoin, "data."+each.Name.ToCamel())
		fieldNameJoin = append(fieldNameJoin, each.Name.Source())
	}
	varLeftJoin = append(varLeftJoin, "prefix")
	keyLeftJoin = append(keyLeftJoin, "key")

	varLeft = util.SafeString(varLeftJoin.Camel().With("").Untitle())
	varRight = fmt.Sprintf(`"%s"`, varRightJoin.Camel().Untitle().With(":").Source()+":")
	varExpression = fmt.Sprintf(`%s = %s`, varLeft, varRight)

	keyLeft = util.SafeString(keyLeftJoin.Camel().With("").Untitle())
	keyRight = fmt.Sprintf(`fmt.Sprintf("%s%s", %s, %s)`, "%s", keyRightArgJoin.With(":").Source(), varLeft, keyRightJoin.With(", ").Source())
	dataKeyRight = fmt.Sprintf(`fmt.Sprintf("%s%s", %s, %s)`, "%s", keyRightArgJoin.With(":").Source(), varLeft, dataRightJoin.With(", ").Source())
	keyExpression = fmt.Sprintf("%s := %s", keyLeft, keyRight)
	dataKeyExpression = fmt.Sprintf("%s := %s", keyLeft, dataKeyRight)

	return Key{
		VarLeft:           varLeft,
		VarRight:          varRight,
		VarExpression:     varExpression,
		KeyLeft:           keyLeft,
		KeyRight:          keyRight,
		DataKeyRight:      dataKeyRight,
		KeyExpression:     keyExpression,
		DataKeyExpression: dataKeyExpression,
		Fields:            in,
		FieldNameJoin:     fieldNameJoin,
	}
}
func (t Table) isIgnoreColumns(columnName string) bool {
	for _, v := range t.ignoreColumns {
		if v == columnName {
			return true
		}
	}
	return false
}
