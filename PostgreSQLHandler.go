package main

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

type PostgreSQLHandler struct {
	db *sql.DB
}

const (
	host     = "appliancestatesdb.cyebc6nm0xm9.eu-west-2.rds.amazonaws.com"
	port     = 5432
	user     = "postgres"
	password = "asdbpassword"
	dbname   = "appliancestatesdb"
)

const updateSingleSQLStatement = `UPDATE HomeAppliances SET current_mode = $2 WHERE name = $1;`
const createToySQLStatement = `INSERT INTO HomeAppliances (name, mode) VALUES ($1, $2) ON CONFLICT DO NOTHING;`
const getNumberOfToys = `SELECT COUNT(id) FROM HomeAppliances;`
const toyDataQuery = `SELECT name, command_with_name, id, publish_topic, subscribe_topic FROM HomeAppliances WHERE id=$1;`

func (postgreHandler *PostgreSQLHandler) Connect() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s",
		host, port, user, password, dbname)

	var err error
	postgreHandler.db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("db connection established")
}

func (postgreHandler *PostgreSQLHandler) TestConnection() {
	result := postgreHandler.db.Ping()
	if result != nil {
		panic(result)
	}
}

func (postgreHandler *PostgreSQLHandler) Disconnect() {
	err := postgreHandler.db.Close()
	if err != nil {
		panic(err)
	}
}

func (postgreHandler *PostgreSQLHandler) UpdateToyMode(toyName string, toyMode string) {
	_, err := postgreHandler.db.Exec(updateSingleSQLStatement, toyName, toyMode)
	if err != nil {
		fmt.Println("Couldnt update mode", err)
	}
}

func (postgreHandler *PostgreSQLHandler) CreateToy(toyName string, toyMode string) {
	_, err := postgreHandler.db.Exec(createToySQLStatement, toyName, toyMode)
	if err != nil {
		fmt.Println("unable to create toy object in db", err)
	}
}

func (postgreHandler *PostgreSQLHandler) GetNumberOfToys() (number int) {
	row := postgreHandler.db.QueryRow(getNumberOfToys)
	switch err := row.Scan(&number); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		fmt.Println(number)
	default:
		panic(err)
	}
	return number
}

func (postgreHandler *PostgreSQLHandler) PullToyData(toyId int) (toyAttributes ToyAttributes) {

	row := postgreHandler.db.QueryRow(toyDataQuery, toyId)
	switch err := row.Scan(&toyAttributes.name, pq.Array(&toyAttributes.commandWithName), &toyAttributes.id, &toyAttributes.publishTopic, &toyAttributes.subscribeTopic); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		fmt.Println(toyAttributes.name, toyAttributes.commandWithName, toyAttributes.id, toyAttributes.publishTopic, toyAttributes.subscribeTopic)
	default:
		panic(err)
	}
	return toyAttributes
}
