package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/lib/pq"
)

type PostgreSQL struct {
	Conn string
}

func NewPostgreSQL(conn string) *PostgreSQL {
	return &PostgreSQL{
		Conn: conn,
	}
}

func (p *PostgreSQL) Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", p.Conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (p *PostgreSQL) Insert(db *sql.DB, tableName string, data any) error {
	quoted := pq.QuoteIdentifier(tableName)

	cols, values, placeholders, err := InjectQuery(data)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", quoted, strings.Join(cols, ", "), strings.Join(placeholders, ", "))

	_, err = db.Exec(query, values...)
	if err != nil {
		return err
	}
	return nil
}

func InjectQuery(data any) ([]string, []interface{}, []string, error) {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)

	var (
		cols         []string
		values       []interface{}
		placeholders []string
	)

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		col := field.Tag.Get("json")
		if strings.Contains(col, "omitempty") && v.Field(i).Interface() == "" {
			continue
		}
		if col != "" {
			col = field.Name
		}

		cols = append(cols, pq.QuoteIdentifier(col))

		value := v.Field(i).Interface()
		if field.Type.Kind() == reflect.Struct {
			jsonVal, err := json.Marshal(value)
			if err != nil {
				return nil, nil, nil, err
			}
			value = string(jsonVal)
		}
		if field.Type.Kind() == reflect.Slice {
			value = pq.Array(value)
		}
		values = append(values, value)

		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
	}
	return cols, values, placeholders, nil
}
