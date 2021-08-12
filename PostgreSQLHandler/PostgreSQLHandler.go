package PostgreSQLHandler

import (
	"database/sql"
	"fmt"
)

const (
	host     = "appliancestatesdb.cyebc6nm0xm9.eu-west-2.rds.amazonaws.com"
	port     = 5432
	user     = "postgres"
	password = "asdbpassword"
	dbname   = "appliancestatesdb"
)

const (
	querySQLStatement        = `SELECT mode FROM HomeAppliances WHERE name = $1;`
	updateSingleSQLStatement = `UPDATE HomeAppliances SET mode = $2 WHERE name = $1;`
)

func Connect() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	return db
}

func TestConnection(db *sql.DB) {
	result := db.Ping()
	if result != nil {
		panic(result)
	}
}

func CloseConnection(db *sql.DB) {
	db.Close()
}

func UpdateMode(db *sql.DB, applianceData []string) {

	_, err := db.Exec(updateSingleSQLStatement, applianceData[0], applianceData[1])
	if err != nil {
		panic(err)
	}
}

func QueryModeProp(db *sql.DB, appliance string) (mode string) {

	row := db.QueryRow(querySQLStatement, appliance)
	switch err := row.Scan(&mode); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		return mode
	default:
		panic(err)
	}
	return mode
}
