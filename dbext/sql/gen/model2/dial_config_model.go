package model2

import (
	"github.com/heshiyingx/gotool/dbext/gormdb"
)

var _ DialConfigDBInterface = (*customDialConfigDB)(nil)

type (
	// DialConfigDBInterface is an interface to be customized, add more methods here,
	// and implement the added methods in customDialConfigModel.
	DialConfigDBInterface interface {
		dialConfigModel
	}

	customDialConfigDB struct {
		*defaultDialConfigModel
	}
)

// NewDialConfigDB returns a model for the database table.
func NewDialConfigDB(config gormdb.Config) DialConfigDBInterface {
	return &customDialConfigDB{
		defaultDialConfigModel: newDefaultDialConfigModel(config),
	}
}
