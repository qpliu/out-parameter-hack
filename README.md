oph is a hack to use out parameters in MySQL stored procedures with Go sql
drivers.

[![GoDoc](https://godoc.org/github.com/qpliu/out-parameter-hack?status.svg)](https://godoc.org/github.com/qpliu/out-parameter-hack)
[![Build Status](https://travis-ci.org/qpliu/out-parameter-hack.svg?branch=master)](https://travis-ci.org/qpliu/out-parameter-hack)

# Example

```go
	db, err := sql.Open("mymysql", "DBNAME/USER/PASSWD")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// CALL EXAMPLE('example',1,NULL,@outString)
	var outString sql.NullString
	if err := oph.Call(db, func(resultSetIndex int, rows *sql.Rows) error {
		for rows.Next() {
			var id int64
			var name sql.NullString
			rows.Scan(&id, &name)
		}
		return rows.Err()
	}, "EXAMPLE", "example", 1, nil, &outString); err != nil {
		panic(err)
	}

```
