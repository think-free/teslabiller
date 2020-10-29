package main

import (
	"flag"

	"github.com/jamiealquiza/envy"

	"pricecalculator/server"
)

// https://golangcode.com/generate-a-pdf/

// http://localhost:8080/getCarConsumption?dateStart=2020-10-01&dateEnd=2020-11-01
// http://localhost:8080/getOtherConsumption?dateStart=2020-10-01&dateEnd=2020-11-01&startCounter=90&endCounter=165
// http://localhost:8080/getBillRequest?dateStart=2020-10-01&dateEnd=2020-11-01&startCounter=90&endCounter=180&tax=21&fix=2&pE1=0.062012&pE2=0.002879&pE3=0.000886&cE1=0.051021&cE2=0.045345&cE3=0.038366

// http://172.16.10.110:5000/getBillRequest?dateStart=2020-10-01&dateEnd=2020-11-01&startCounter=90&endCounter=215&tax=21&fix=2&pE1=0.062012&pE2=0.002879&pE3=0.000886&cE1=0.051021&cE2=0.045345&cE3=0.038366

func main() {

	dbUser := flag.String("dbUser", "teslamate", "Database user")
	dbPass := flag.String("dbPass", "secret", "Database password")
	dbServer := flag.String("dbServer", "172.16.10.110", "Database addr")
	dbPort := flag.String("dbPort", "5432", "Database port")

	envy.Parse("PC")
	flag.Parse()

	s := server.New("postgres://" + *dbUser + ":" + *dbPass + "@" + *dbServer + ":" + *dbPort + "/teslamate?sslmode=disable")
	s.Run()
}
