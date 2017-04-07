// Package oph is a hack to use out parameter hack in MySQL stored procedures.
package oph

import (
	"bytes"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"
)

// Queryer can run queries.  Examples include *sql.DB and *sql.Tx.
type Queryer interface {
	Query(string, ...interface{}) (*sql.Rows, error)
}

// CallString is the SQL that would be sent to the queryer calling the
// stored procedure with the given arguments.
func CallString(proc string, params ...interface{}) (string, error) {
	callString, _, err := callParameters(proc, params...)
	return callString, err
}

// Call calls the stored procedure from the queryer with the given parameters.
//
// readResultSet is called for each result set returned by the stored
// procedure.
//
// Parameters of types implementing sql.Scanner are passed to the stored
// procedure as out parameters and scanned after the result sets returned
// by the stored procedure are read.
//
// Parameters with type nil, bool, *bool, int, int8, int16, int32, int64,
// *int, *int8, *int16, *int32, *int64, uint, uint8, uint16, uint32, uint64,
// *uint, *uint8, *uint16, *uint32, *uint64, float32, float64, *float32,
// *float64, string, *string, time.Time, *time.Time are passed to the stored
// procedure as in parameters.
//
// Parameters of all other types result in error.
func Call(queryer Queryer, readResultSet func(resultSetIndex int, rows *sql.Rows) error, proc string, params ...interface{}) error {
	callString, outParams, err := callParameters(proc, params...)
	if err != nil {
		return err
	}

	rows, err := queryer.Query(callString)
	if err != nil {
		return err
	}
	defer rows.Close()

	resultSetIndex := 0
	for {
		cols, err := rows.Columns()
		if err != nil {
			return err
		}
		if len(cols) == len(outParams) && len(cols) > 0 && cols[0] == "@1" {
			if rows.Next() {
				if err := rows.Scan(outParams...); err != nil {
					return err
				}
			}
			return rows.Err()
		}
		if err := readResultSet(resultSetIndex, rows); err != nil {
			return err
		}
		if rows.NextResultSet() {
			resultSetIndex++
		} else {
			return rows.Err()
		}
	}
}

func callParameters(proc string, params ...interface{}) (string, []interface{}, error) {
	var buf bytes.Buffer
	var outParams []interface{}

	buf.WriteString("CALL ")
	buf.WriteString(proc)
	buf.WriteString("(")

	for i, param := range params {
		if i > 0 {
			buf.WriteString(",")
		}
		if _, ok := param.(sql.Scanner); ok {
			outParams = append(outParams, param)
			buf.WriteString("@")
			buf.WriteString(strconv.Itoa(len(outParams)))
			continue
		}
		switch p := param.(type) {
		case nil:
			buf.WriteString("NULL")
		case bool:
			if p {
				buf.WriteString("1")
			} else {
				buf.WriteString("0")
			}
		case *bool:
			if p == nil {
				buf.WriteString("NULL")
			} else if *p {
				buf.WriteString("1")
			} else {
				buf.WriteString("0")
			}
		case int:
			buf.WriteString(strconv.FormatInt(int64(p), 10))
		case int8:
			buf.WriteString(strconv.FormatInt(int64(p), 10))
		case int16:
			buf.WriteString(strconv.FormatInt(int64(p), 10))
		case int32:
			buf.WriteString(strconv.FormatInt(int64(p), 10))
		case int64:
			buf.WriteString(strconv.FormatInt(int64(p), 10))
		case *int:
			if p == nil {
				buf.WriteString("NULL")
			} else {
				buf.WriteString(strconv.FormatInt(int64(*p), 10))
			}
		case *int8:
			if p == nil {
				buf.WriteString("NULL")
			} else {
				buf.WriteString(strconv.FormatInt(int64(*p), 10))
			}
		case *int16:
			if p == nil {
				buf.WriteString("NULL")
			} else {
				buf.WriteString(strconv.FormatInt(int64(*p), 10))
			}
		case *int32:
			if p == nil {
				buf.WriteString("NULL")
			} else {
				buf.WriteString(strconv.FormatInt(int64(*p), 10))
			}
		case *int64:
			if p == nil {
				buf.WriteString("NULL")
			} else {
				buf.WriteString(strconv.FormatInt(int64(*p), 10))
			}
		case uint:
			buf.WriteString(strconv.FormatUint(uint64(p), 10))
		case uint8:
			buf.WriteString(strconv.FormatUint(uint64(p), 10))
		case uint16:
			buf.WriteString(strconv.FormatUint(uint64(p), 10))
		case uint32:
			buf.WriteString(strconv.FormatUint(uint64(p), 10))
		case uint64:
			buf.WriteString(strconv.FormatUint(uint64(p), 10))
		case *uint:
			if p == nil {
				buf.WriteString("NULL")
			} else {
				buf.WriteString(strconv.FormatUint(uint64(*p), 10))
			}
		case *uint8:
			if p == nil {
				buf.WriteString("NULL")
			} else {
				buf.WriteString(strconv.FormatUint(uint64(*p), 10))
			}
		case *uint16:
			if p == nil {
				buf.WriteString("NULL")
			} else {
				buf.WriteString(strconv.FormatUint(uint64(*p), 10))
			}
		case *uint32:
			if p == nil {
				buf.WriteString("NULL")
			} else {
				buf.WriteString(strconv.FormatUint(uint64(*p), 10))
			}
		case *uint64:
			if p == nil {
				buf.WriteString("NULL")
			} else {
				buf.WriteString(strconv.FormatUint(uint64(*p), 10))
			}
		case float32:
			buf.WriteString(strconv.FormatFloat(float64(p), 'G', -1, 32))
		case *float32:
			if p == nil {
				buf.WriteString("NULL")
			} else {
				buf.WriteString(strconv.FormatFloat(float64(*p), 'G', -1, 32))
			}
		case float64:
			buf.WriteString(strconv.FormatFloat(float64(p), 'G', -1, 64))
		case *float64:
			if p == nil {
				buf.WriteString("NULL")
			} else {
				buf.WriteString(strconv.FormatFloat(float64(*p), 'G', -1, 64))
			}
		case string:
			writeEscapedString(&buf, p)
		case *string:
			if p == nil {
				buf.WriteString("NULL")
			} else {
				writeEscapedString(&buf, *p)
			}
		case time.Time:
			buf.WriteString(p.UTC().Format("'2006-01-02 15:04:05'"))
		case *time.Time:
			if p == nil {
				buf.WriteString("NULL")
			} else {
				buf.WriteString(p.UTC().Format("'2006-01-02 15:04:05'"))
			}
		default:
			return "", nil, errors.New("oph: Unsupported type for parameter " + strconv.Itoa(i))
		}
	}

	buf.WriteString(")")
	if len(outParams) > 0 {
		buf.WriteString(";SELECT ")
		for i := 0; i < len(outParams); i++ {
			if i > 0 {
				buf.WriteString(",")
			}
			buf.WriteString("@")
			buf.WriteString(strconv.Itoa(i + 1))
		}
	}

	return buf.String(), outParams, nil
}

func writeEscapedString(buf *bytes.Buffer, str string) {
	buf.WriteString("'")
	for len(str) > 0 {
		i := strings.IndexAny(str, "'\\")
		if i < 0 {
			buf.WriteString(str)
			break
		}
		buf.WriteString(str[:i])
		buf.WriteByte(str[i])
		buf.WriteByte(str[i])
		str = str[i+1:]
	}
	buf.WriteString("'")
}
