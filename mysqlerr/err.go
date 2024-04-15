package mysqlerr

import "github.com/go-sql-driver/mysql"

func IsMysqlNotExist(err error) bool {
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}
	switch mysqlErr.Number {
	case 1:

	}
	return false
}
