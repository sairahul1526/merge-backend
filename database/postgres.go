package database

import (
	"database/sql"
	CONSTANT "merge-backend/constant"
	"time"

	LOGGER "merge-backend/logger"
	UTIL "merge-backend/util"
	"strconv"
	"strings"

	_ "github.com/lib/pq" // for postgres driver
)

// postgres db

type Postgres struct {
	dbConfig       string
	connectionPool int
	maxLifeTime    time.Duration
	db             *sql.DB
}

// Connect - connects to postgres db
func (pg *Postgres) Connect() error {
	LOGGER.Log(pg.dbConfig, pg.connectionPool)
	var err error
	pg.db, err = sql.Open("postgres", pg.dbConfig)
	if err != nil {
		LOGGER.Warn(pg.dbConfig, pg.connectionPool, err)
		return err
	}

	// database conection pooling
	pg.db.SetMaxOpenConns(pg.connectionPool)
	pg.db.SetMaxIdleConns(pg.connectionPool)
	pg.db.SetConnMaxLifetime(pg.maxLifeTime)

	return nil
}

// InsertWithUniqueID - insert data into table with unique id
func (pg *Postgres) InsertWithUniqueID(table string, body map[string]string, key string) (string, sql.Result, error) {
	LOGGER.Log(table, body, key)
	var (
		result sql.Result
		err    error
	)
	for i := 0; i < 10; i++ { // try to insert with unqiue id for certain number of times; if no limit, server crashes in certain conditions
		body[key] = generateRandomID()
		result, err = pg.InsertSQL(table, body)
		if err == nil {
			break
		}
	}
	if err != nil {
		LOGGER.Warn(table, body, key, err)
		return "", nil, err
	}
	return body[key], result, nil
}

// RowCount - get number of items in database with specified query
func (pg *Postgres) RowCount(tableName string, where string, args ...interface{}) (int, error) {
	LOGGER.Log(tableName, where, args)
	data, err := pg.SelectProcess("select count(*) as ctn from "+tableName+" where "+where, args...)
	if err != nil {
		LOGGER.Warn(tableName, where, args, err)
		return 0, err
	}
	if len(data) == 0 {
		return 0, nil
	}
	count, _ := strconv.Atoi(data[0]["ctn"])
	return count, nil
}

// CheckIfExists - check if data exists in table
func (pg *Postgres) CheckIfExists(table string, params map[string]string) error {
	LOGGER.Log(table, params)
	data, err := pg.SelectSQL(table, []string{"1"}, params)
	if err != nil {
		LOGGER.Warn(table, params, err)
		return err
	}
	if len(data) > 0 {
		return nil
	}
	return CONSTANT.SQLCheckIfExistsEmptyError
}

// ExecuteSQL - execute statement with defined values, with all the input as params to prevent sql injection
func (pg *Postgres) ExecuteSQL(SQLQuery string, params ...interface{}) (sql.Result, error) {
	LOGGER.Log(SQLQuery, params)
	result, err := pg.db.Exec(SQLQuery, params...)
	if err != nil {
		LOGGER.Warn(SQLQuery, params, err)
		return nil, err
	}

	return result, nil
}

// QueryRowSQL - get single data with defined values, with all the input as params to prevent sql injection
func (pg *Postgres) QueryRowSQL(SQLQuery string, params ...interface{}) (string, error) {
	LOGGER.Log(SQLQuery, params)

	var value string
	err := pg.db.QueryRow(SQLQuery, params...).Scan(&value)
	if err != nil {
		LOGGER.Warn(SQLQuery, params, err)
		return "", err
	}

	return value, nil
}

// UpdateSQL - update data with defined values
func (pg *Postgres) UpdateSQL(tableName string, params map[string]string, body map[string]string) (sql.Result, error) {
	LOGGER.Log(tableName, params, body)
	if len(body) == 0 {
		LOGGER.Warn(tableName, params, body, CONSTANT.SQLUpdateBodyEmptyError)
		return nil, CONSTANT.SQLUpdateBodyEmptyError
	}

	args := []interface{}{}
	SQLQuery := "update " + tableName + " set "

	init := false
	i := 1
	for key, val := range body {
		if init {
			SQLQuery += ","
		}
		SQLQuery += `"` + key + `" = $` + strconv.Itoa(i)
		args = append(args, val)
		init = true
		i++
	}
	// add updated_at
	SQLQuery += `, "updated_at" = $` + strconv.Itoa(i)
	args = append(args, UTIL.GetCurrentTime())
	i++

	SQLQuery += " where "
	init = false
	for key, val := range params {
		if init {
			SQLQuery += " and "
		}
		SQLQuery += `"` + key + `" = $` + strconv.Itoa(i)
		args = append(args, val)
		init = true
		i++
	}

	LOGGER.Log(SQLQuery, args)
	result, err := pg.db.Exec(SQLQuery, args...)
	if err != nil {
		LOGGER.Warn(tableName, params, body, err)
		return nil, err
	}

	return result, nil
}

// DeleteSQL - delete data with defined values
func (pg *Postgres) DeleteSQL(tableName string, params map[string]string) (sql.Result, error) {
	LOGGER.Log(tableName, params)
	if len(params) == 0 {
		// atleast one value should be specified for deleting, cannot delete all values
		LOGGER.Warn(tableName, params, CONSTANT.SQLDeleteAllNotAllowedError)
		return nil, CONSTANT.SQLDeleteAllNotAllowedError
	}

	args := []interface{}{}
	SQLQuery := "delete from " + tableName + " where "

	init := false
	i := 1
	for key, val := range params {
		if init {
			SQLQuery += " and "
		}
		SQLQuery += `"` + key + `" = $` + strconv.Itoa(i)
		args = append(args, val)
		init = true
		i++
	}

	LOGGER.Log(SQLQuery, args)
	result, err := pg.db.Exec(SQLQuery, args...)
	if err != nil {
		LOGGER.Warn(tableName, params, err)
		return nil, err
	}
	return result, nil
}

// InsertSQL - insert data with defined values
func (pg *Postgres) InsertSQL(tableName string, body map[string]string) (sql.Result, error) {
	LOGGER.Log(tableName, body)
	if len(body) == 0 {
		LOGGER.Warn(tableName, body, CONSTANT.SQLInsertBodyEmptyError)
		return nil, CONSTANT.SQLInsertBodyEmptyError
	}

	SQLQuery, args := pg.BuildInsertStatement(tableName, body)

	LOGGER.Log(SQLQuery, args)
	result, err := pg.db.Exec(SQLQuery, args...)
	if err != nil {
		LOGGER.Warn(tableName, body, err)
		return nil, err
	}
	return result, nil
}

// BuildInsertStatement - build insert statement with defined values
func (pg *Postgres) BuildInsertStatement(tableName string, body map[string]string) (string, []interface{}) {
	args := []interface{}{}
	SQLQuery := "insert into " + tableName + " "
	keys := " ("
	values := " ("
	init := false
	i := 1
	for key, val := range body {
		if init {
			keys += ","
			values += ","
		}
		keys += ` "` + key + `" `
		values += " $" + strconv.Itoa(i)
		args = append(args, val)
		init = true
		i++
	}
	// add created_at, updated_at
	keys += `, "created_at", "updated_at" `
	values += ", $" + strconv.Itoa(i) + ", $" + strconv.Itoa(i+1)
	args = append(args, UTIL.GetCurrentTime(), UTIL.GetCurrentTime())

	keys += ")"
	values += ")"
	SQLQuery += keys + " values " + values
	return SQLQuery, args
}

// SelectSQL - query data with defined values
func (pg *Postgres) SelectSQL(tableName string, columns []string, params map[string]string) ([]map[string]string, error) {
	args := []interface{}{}
	SQLQuery := "select " + strings.Join(columns, ",") + " from " + tableName + ""
	if len(params) > 0 {
		where := ""
		init := false
		i := 1
		for key, val := range params {
			if init {
				where += " and "
			}
			where += ` "` + key + `" = $` + strconv.Itoa(i)
			args = append(args, val)
			init = true
			i++
		}
		if strings.Compare(where, "") != 0 {
			SQLQuery += " where " + where
		}
	}
	return pg.SelectProcess(SQLQuery, args...)
}

// SelectProcess - execute raw select statement, with all the input as params to prevent sql injection
func (pg *Postgres) SelectProcess(SQLQuery string, params ...interface{}) ([]map[string]string, error) {
	LOGGER.Log(SQLQuery, params)
	rows, err := pg.db.Query(SQLQuery, params...)
	if err != nil {
		LOGGER.Warn(SQLQuery, params, err)
		return []map[string]string{}, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		LOGGER.Warn(SQLQuery, params, err)
		return []map[string]string{}, err
	}

	rawResult := make([][]byte, len(cols))

	dest := make([]interface{}, len(cols))
	data := []map[string]string{}
	rest := map[string]string{}
	for i := range rawResult {
		dest[i] = &rawResult[i]
	}

	for rows.Next() {
		rest = map[string]string{}
		err = rows.Scan(dest...)
		if err != nil {
			LOGGER.Warn(SQLQuery, params, err)
			return []map[string]string{}, err
		}

		for i, raw := range rawResult {
			if raw == nil {
				rest[cols[i]] = ""
			} else {
				rest[cols[i]] = string(raw)
			}
		}

		data = append(data, rest)
	}

	return data, nil
}
