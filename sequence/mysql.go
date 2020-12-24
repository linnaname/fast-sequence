package sequence

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	conf "fast-sequence/conf"

	_ "github.com/go-sql-driver/mysql"
)

/**
	CREATE TABLE `sequence` (
	 `name` varchar(32) NOT NULL,
	 `value` bigint NOT NULL,
	 `step` bigint NOT NULL,
	 `gmt_modified` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 	 `desc` varchar(1024) DEFAULT '' NOT NULL,
	  PRIMARY KEY (`name`)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;
*/
const (
	TABLE_NAME          = "sequence"
	VALUE_COLUMN        = "value"
	NAME_COLUMN         = "name"
	GMT_MODIFIED_COLUMN = "gmt_modified"
)

var gdb *sql.DB

func initMySql() error {
	config, err := conf.ReadConfig("./../conf/conf.json")
	if err != nil {
		return err
	}
	err = initSqlDB(config)
	if err != nil {
		return err
	}
	return nil
}

func initSqlDB(conf *conf.Config) error {
	if conf == nil {
		return errors.New("conf is nil")
	}

	db, err := sql.Open("mysql", conf.DataSourceName)
	if err != nil {
		return err
	}

	//Configuring sql.DB for Better Performance:https://www.alexedwards.net/blog/configuring-sqldb
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(0)
	gdb = db
	return nil
}

func NextRange(name string) (sr *SequenceRange, err error) {
	if name == "" {
		return nil, errors.New("name can't be empty")
	}

	if gdb == nil {
		err := initMySql()
		if err != nil {
			return nil, errors.New("failed to init mysql connection")
		}
	}

	var (
		oldValue int64
		newValue int64
	)

	selectSql := getSelectSql()
	err = gdb.QueryRow(selectSql, name).Scan(&oldValue)
	if err != nil {
		return nil, err
	}

	if oldValue < 0 {
		return nil, errors.New("Sequence value cannot be less than zero")
	}

	if oldValue > 9223372036754775807 {
		return nil, errors.New("Sequence value overflow")
	}

	newValue = oldValue + int64(getStep())

	updateSql := getUpdateSql()
	rst, err := gdb.Exec(updateSql, newValue, time.Now(), name, oldValue)
	if err != nil {
		return nil, err
	}

	affectedRows, err := rst.RowsAffected()
	if err != nil || affectedRows == 0 {
		return nil, err
	}

	localSequenceRange := NewSequenceRange(oldValue, newValue)
	return localSequenceRange, nil
}

//代码可读性稍差但性能比fmt稍好
func getSelectSql() string {
	builder := strings.Builder{}
	builder.WriteString("select ")
	builder.WriteString(VALUE_COLUMN)
	builder.WriteString(" from ")
	builder.WriteString(TABLE_NAME)
	builder.WriteString(" where ")
	builder.WriteString(NAME_COLUMN)
	builder.WriteString(" = ? ")
	return builder.String()
}

/**
采用乐观锁，避免事务
*/
func getUpdateSql() string {
	builder := strings.Builder{}
	builder.WriteString("update ")
	builder.WriteString(TABLE_NAME)
	builder.WriteString(" set ")
	builder.WriteString(VALUE_COLUMN)
	builder.WriteString(" = ?, ")
	builder.WriteString(GMT_MODIFIED_COLUMN)
	builder.WriteString(" = ? where ")
	builder.WriteString(NAME_COLUMN)
	builder.WriteString(" = ? and ")
	builder.WriteString(VALUE_COLUMN)
	builder.WriteString(" = ? ")
	return builder.String()
}

func getStep() int {
	return 10000
}
