package pkg

import "database/sql"

func CloseRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {

		PanicIfError(err)
	}
}
