package main

const (
	dbUser          = "postgres"
	dbPassword      = "1234"
	dbName          = "golang_postgres_rest_api"
	dbHost          = "127.0.0.1"
	dbPort          = "5432"
	dbSSLMode       = "disable"

	applicationPort = ":8100"
)

func main() {
	a := App{}

	a.Initialize(dbUser, dbPassword, dbName, dbHost, dbPort, dbSSLMode)
	a.Run(applicationPort)
}