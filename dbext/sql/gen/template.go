package gen

import "github.com/heshiyingx/gotool/dbext/sql/template"

const (
	category                              = "model"
	customizedTemplateFile                = "customized.tpl"
	deleteTemplateFile                    = "delete.tpl"
	deleteMethodTemplateFile              = "interface-delete.tpl"
	fieldTemplateFile                     = "field.tpl"
	findOneTemplateFile                   = "find-by-pk.tpl"
	findOneMethodTemplateFile             = "interface-find-by-pk.tpl"
	findOneByFieldTemplateFile            = "find-one-by-field.tpl"
	findOneByFieldMethodTemplateFile      = "interface-find-one-by-field.tpl"
	findOneByFieldExtraMethodTemplateFile = "find-one-by-field-extra-method.tpl"
	importsTemplateFile                   = "import.tpl"
	importsWithNoCacheTemplateFile        = "import-no-cache.tpl"
	insertTemplateFile                    = "insert.tpl"
	insertTemplateMethodFile              = "interface-insert.tpl"
	modelGenTemplateFile                  = "model-gen.tpl"
	modelCustomTemplateFile               = "model.tpl"
	modelNewTemplateFile                  = "model-new.tpl"
	tableNameTemplateFile                 = "table-name.tpl"
	tagTemplateFile                       = "tag.tpl"
	typesTemplateFile                     = "types.tpl"
	updateTemplateFile                    = "update.tpl"
	updateMethodTemplateFile              = "interface-update.tpl"
	varTemplateFile                       = "var.tpl"
	VarFileTemplateFile                   = "var_file.tpl"
)

var templates = map[string]string{
	customizedTemplateFile:                template.Customized,
	deleteTemplateFile:                    template.Delete,
	deleteMethodTemplateFile:              template.DeleteMethod,
	fieldTemplateFile:                     template.Field,
	findOneTemplateFile:                   template.FindByPK,
	findOneMethodTemplateFile:             template.FindOneMethod,
	findOneByFieldTemplateFile:            template.FindOneByField,
	findOneByFieldMethodTemplateFile:      template.FindOneByFieldMethod,
	findOneByFieldExtraMethodTemplateFile: template.FindOneByFieldExtraMethod,
	importsTemplateFile:                   template.Imports,
	importsWithNoCacheTemplateFile:        template.ImportsNoCache,
	insertTemplateFile:                    template.Insert,
	insertTemplateMethodFile:              template.InsertMethod,
	modelGenTemplateFile:                  template.ModelGen,
	modelCustomTemplateFile:               template.ModelCustom,
	modelNewTemplateFile:                  template.New,
	tableNameTemplateFile:                 template.TableName,
	tagTemplateFile:                       template.Tag,
	typesTemplateFile:                     template.Types,
	updateTemplateFile:                    template.Update,
	updateMethodTemplateFile:              template.UpdateMethod,
	varTemplateFile:                       template.Vars,
	VarFileTemplateFile:                   template.VarFile,
}
