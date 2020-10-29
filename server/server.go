package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

// Server define a http server
type Server struct {
	dbCon string
}

// Answer is the json answer to http request
type Answer struct {
	CarP1kWh         float64 `json:"carP1kWh"`
	CarP2kWh         float64 `json:"carP2kWh"`
	CarP3kWh         float64 `json:"carP3kWh"`
	CarTotalkWh      float64 `json:"carTotalkWh"`
	TotalConsumption float64 `json:"totalkWh"`
	OtherConsumption float64 `json:"otherTotalkWh"`

	CarP1Cost       float64 `json:"carP1Cost"`
	CarP2Cost       float64 `json:"carP2Cost"`
	CarP3Cost       float64 `json:"carP3Cost"`
	CarTotalCost    float64 `json:"carTotalCost"`
	CarTotalCostTax float64 `json:"carTotalCostTax"`
	OtherCost       float64 `json:"otherCost"`
	OtherCostTax    float64 `json:"otherCostTax"`
	TotalCost       float64 `json:"totalCost"`
	TotalCostTax    float64 `json:"totalCostTax"`
}

// New : create a new server
func New(con string) *Server {

	return &Server{
		dbCon: con,
	}
}

// Run : run the server
func (s *Server) Run() {

	http.HandleFunc("/getCarConsumption", s.getCarConsumptionRequest)
	http.HandleFunc("/getOtherConsumption", s.getOtherConsumptionRequest)
	http.HandleFunc("/getBillRequest", s.getBillRequest)

	http.ListenAndServe("0.0.0.0:8080", nil)
}

func (s *Server) getCarConsumptionRequest(w http.ResponseWriter, r *http.Request) {

	// Getting parameters

	dateStartParam, ok := r.URL.Query()["dateStart"]
	if !ok || len(dateStartParam[0]) < 1 {
		log.Println("Url parameter missing: dateStart")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"type\" : \"error\", \"message\":\"Url parameter missing: dateStart\"}"))
		return
	}
	dateStart := dateStartParam[0]

	dateEndParam, ok := r.URL.Query()["dateEnd"]
	if !ok || len(dateEndParam[0]) < 1 {
		log.Println("Url parameter missing: dateEnd")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"type\" : \"error\", \"message\":\"Url parameter missing: dateEnd\"}"))
		return
	}
	dateEnd := dateEndParam[0]

	// Query database

	ans := s.getCarConsumption(dateStart, dateEnd)

	// Send result

	js, err := json.Marshal(&ans)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"type\" : \"error\", \"message\":\"Can't parse json answer\"}"))
		return
	}
	w.Write(js)
}

func (s *Server) getOtherConsumptionRequest(w http.ResponseWriter, r *http.Request) {

	// Getting parameters

	dateStartParam, ok := r.URL.Query()["dateStart"]
	if !ok || len(dateStartParam[0]) < 1 {
		log.Println("Url parameter missing: dateStart")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"type\" : \"error\", \"message\":\"Url parameter missing: dateStart\"}"))
		return
	}
	dateStart := dateStartParam[0]

	dateEndParam, ok := r.URL.Query()["dateEnd"]
	if !ok || len(dateEndParam[0]) < 1 {
		log.Println("Url parameter missing: dateEnd")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"type\" : \"error\", \"message\":\"Url parameter missing: dateEnd\"}"))
		return
	}
	dateEnd := dateEndParam[0]

	startCounterParam, ok := r.URL.Query()["startCounter"]
	if !ok || len(startCounterParam[0]) < 1 {
		log.Println("Url parameter missing: startCounter")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"type\" : \"error\", \"message\":\"Url parameter missing: startCounter\"}"))
		return
	}
	startCounter, _ := strconv.ParseFloat(startCounterParam[0], 64)

	endCounterParam, ok := r.URL.Query()["endCounter"]
	if !ok || len(endCounterParam[0]) < 1 {
		log.Println("Url parameter missing: endCounter")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"type\" : \"error\", \"message\":\"Url parameter missing: endCounter\"}"))
		return
	}
	endCounter, _ := strconv.ParseFloat(endCounterParam[0], 64)

	// Query database

	ans := s.getCarConsumption(dateStart, dateEnd)
	ans = s.calculateTotalConsumption(ans, startCounter, endCounter)

	// Send result

	js, err := json.Marshal(&ans)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"type\" : \"error\", \"message\":\"Can't parse json answer\"}"))
		return
	}
	w.Write(js)
}

func (s *Server) getBillRequest(w http.ResponseWriter, r *http.Request) {

	// Getting parameters

	dateStartParam, ok := r.URL.Query()["dateStart"]
	if !ok || len(dateStartParam[0]) < 1 {
		log.Println("Url parameter missing: dateStart")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"type\" : \"error\", \"message\":\"Url parameter missing: dateStart\"}"))
		return
	}
	dateStart := dateStartParam[0]

	dateEndParam, ok := r.URL.Query()["dateEnd"]
	if !ok || len(dateEndParam[0]) < 1 {
		log.Println("Url parameter missing: dateEnd")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"type\" : \"error\", \"message\":\"Url parameter missing: dateEnd\"}"))
		return
	}
	dateEnd := dateEndParam[0]

	startCounterParam, ok := r.URL.Query()["startCounter"]
	if !ok || len(startCounterParam[0]) < 1 {
		log.Println("Url parameter missing: startCounter")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"type\" : \"error\", \"message\":\"Url parameter missing: startCounter\"}"))
		return
	}
	startCounter, _ := strconv.ParseFloat(startCounterParam[0], 64)

	endCounterParam, ok := r.URL.Query()["endCounter"]
	if !ok || len(endCounterParam[0]) < 1 {
		log.Println("Url parameter missing: endCounter")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"type\" : \"error\", \"message\":\"Url parameter missing: endCounter\"}"))
		return
	}
	endCounter, _ := strconv.ParseFloat(endCounterParam[0], 64)

	pE1Param, ok := r.URL.Query()["pE1"]
	if !ok || len(pE1Param[0]) < 1 {
		log.Println("Url parameter missing: pE1")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"type\" : \"error\", \"message\":\"Url parameter missing: pE1\"}"))
		return
	}
	pE1, _ := strconv.ParseFloat(pE1Param[0], 64)

	pE2Param, ok := r.URL.Query()["pE2"]
	if !ok || len(pE2Param[0]) < 1 {
		log.Println("Url parameter missing: pE2")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"type\" : \"error\", \"message\":\"Url parameter missing: pE2\"}"))
		return
	}
	pE2, _ := strconv.ParseFloat(pE2Param[0], 64)

	pE3Param, ok := r.URL.Query()["pE3"]
	if !ok || len(pE3Param[0]) < 1 {
		log.Println("Url parameter missing: pE3")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"type\" : \"error\", \"message\":\"Url parameter missing: pE3\"}"))
		return
	}
	pE3, _ := strconv.ParseFloat(pE3Param[0], 64)

	cE1Param, ok := r.URL.Query()["cE1"]
	if !ok || len(pE1Param[0]) < 1 {
		log.Println("Url parameter missing: cE1")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"type\" : \"error\", \"message\":\"Url parameter missing: cE1\"}"))
		return
	}
	cE1, _ := strconv.ParseFloat(cE1Param[0], 64)

	cE2Param, ok := r.URL.Query()["cE2"]
	if !ok || len(cE2Param[0]) < 1 {
		log.Println("Url parameter missing: cE2")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"type\" : \"error\", \"message\":\"Url parameter missing: cE2\"}"))
		return
	}
	cE2, _ := strconv.ParseFloat(cE2Param[0], 64)

	cE3Param, ok := r.URL.Query()["cE3"]
	if !ok || len(cE3Param[0]) < 1 {
		log.Println("Url parameter missing: cE3")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"type\" : \"error\", \"message\":\"Url parameter missing: cE3\"}"))
		return
	}
	cE3, _ := strconv.ParseFloat(cE3Param[0], 64)

	taxParam, ok := r.URL.Query()["tax"]
	if !ok || len(taxParam[0]) < 1 {
		log.Println("Url parameter missing: tax")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"type\" : \"error\", \"message\":\"Url parameter missing: tax\"}"))
		return
	}
	tax, _ := strconv.ParseFloat(taxParam[0], 64)

	fixParam, ok := r.URL.Query()["fix"]
	if !ok || len(fixParam[0]) < 1 {
		log.Println("Url parameter missing: fix")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"type\" : \"error\", \"message\":\"Url parameter missing: fix\"}"))
		return
	}
	fix, _ := strconv.ParseFloat(fixParam[0], 64)

	// Query database

	ans := s.getCarConsumption(dateStart, dateEnd)
	ans = s.calculateTotalConsumption(ans, startCounter, endCounter)
	ans = s.getCarConsumptionPrice(ans, pE1, pE2, pE3, cE1, cE2, cE3, tax, fix)

	// Send result

	js, err := json.Marshal(&ans)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"type\" : \"error\", \"message\":\"Can't parse json answer\"}"))
		return
	}
	w.Write(js)
}

// Internal

func (s *Server) getCarConsumption(dateStart, dateEnd string) *Answer {

	db, err := sql.Open("postgres", s.dbCon)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	// Query database

	var pkWhStr string

	ans := &Answer{}

	// P1

	p1Sql := `SELECT sum(charge_energy_used) FROM charging_processes
				WHERE start_date >= $1
				AND start_date < $2
				AND CAST(start_date at time zone 'gmt+2' AS time) BETWEEN '13:00' and '23:00'
				AND address_id = '1'`

	err = db.QueryRow(p1Sql, dateStart, dateEnd).Scan(&pkWhStr)
	if err != nil {
		fmt.Println("Failed to execute p1 query: ", err)
	} else {
		ans.CarP1kWh, _ = strconv.ParseFloat(pkWhStr, 64)
	}

	// P2

	p2Sql := `SELECT sum(charge_energy_used) charge_energy_used FROM
				(
					SELECT sum(charge_energy_used) charge_energy_used FROM charging_processes
					WHERE start_date >= $1
					AND start_date < $2
					AND CAST(start_date at time zone 'gmt+2' AS time) BETWEEN '07:00' and '13:00'
					AND address_id = '1'
					UNION
					SELECT sum(charge_energy_used) charge_energy_used FROM charging_processes
					WHERE start_date >= $1
					AND start_date < $2
					AND CAST(start_date at time zone 'gmt+2' AS time) BETWEEN '23:00' and '23:59'
					AND address_id = '1'
					UNION
					SELECT sum(charge_energy_used) charge_energy_used FROM charging_processes
					WHERE start_date >= $1
					AND start_date < $2
					AND CAST(start_date at time zone 'gmt+2' AS time) BETWEEN '00:00' and '01:00'
					AND address_id = '1'
				)s`

	err = db.QueryRow(p2Sql, dateStart, dateEnd).Scan(&pkWhStr)
	if err != nil {
		fmt.Println("Failed to execute p2 query: ", err)
	} else {
		ans.CarP2kWh, _ = strconv.ParseFloat(pkWhStr, 64)
	}

	// P3

	p3Sql := `SELECT sum(charge_energy_used) FROM charging_processes
				WHERE start_date >= $1
				AND start_date < $2
				AND CAST(start_date at time zone 'gmt+2' AS time) BETWEEN '01:00' and '07:00'
				AND address_id = '1'`

	err = db.QueryRow(p3Sql, dateStart, dateEnd).Scan(&pkWhStr)
	if err != nil {
		fmt.Println("Failed to execute p3 query: ", err)
	} else {
		ans.CarP3kWh, _ = strconv.ParseFloat(pkWhStr, 64)
	}

	// Total

	ans.CarTotalkWh = ans.CarP1kWh + ans.CarP2kWh + ans.CarP3kWh
	ans.CarTotalkWh = math.Floor((ans.CarTotalkWh)*100) / 100

	return ans
}

func (s *Server) calculateTotalConsumption(ans *Answer, startCounter, endCounter float64) *Answer {

	// Calculate

	ans.TotalConsumption = endCounter - startCounter
	ans.OtherConsumption = ans.TotalConsumption - ans.CarTotalkWh

	// Round

	ans.TotalConsumption = math.Floor((ans.TotalConsumption)*100) / 100
	ans.OtherConsumption = math.Floor((ans.OtherConsumption)*100) / 100

	return ans
}

func (s *Server) getCarConsumptionPrice(ans *Answer, pE1, pE2, pE3, cE1, cE2, cE3, tax, fix float64) *Answer {

	// Calculate

	ans.CarP1Cost = (ans.CarP1kWh * pE1) + (ans.CarP1kWh * cE1)
	ans.CarP2Cost = (ans.CarP2kWh * pE2) + (ans.CarP2kWh * cE2)
	ans.CarP3Cost = (ans.CarP3kWh * pE3) + (ans.CarP3kWh * cE3)

	ans.CarTotalCost = ans.CarP1Cost + ans.CarP2Cost + ans.CarP3Cost
	ans.CarTotalCostTax = ans.CarTotalCost + ((ans.CarTotalCost / 100) * tax)

	ans.OtherCost = (ans.OtherConsumption * pE2) + (ans.OtherConsumption * cE2)
	ans.OtherCostTax = ans.OtherCost + ((ans.OtherCost / 100) * tax)

	ans.TotalCost = ans.CarTotalCost + ans.OtherCost + fix
	ans.TotalCostTax = ans.TotalCost + ((ans.TotalCost / 100) * tax)

	// Round

	ans.CarP1Cost = math.Floor((ans.CarP1Cost)*100) / 100
	ans.CarP2Cost = math.Floor((ans.CarP2Cost)*100) / 100
	ans.CarP3Cost = math.Floor((ans.CarP3Cost)*100) / 100
	ans.CarTotalCost = math.Floor((ans.CarTotalCost)*100) / 100
	ans.CarTotalCostTax = math.Floor((ans.CarTotalCostTax)*100) / 100
	ans.OtherCost = math.Floor((ans.OtherCost)*100) / 100
	ans.OtherCostTax = math.Floor((ans.OtherCostTax)*100) / 100
	ans.TotalCost = math.Floor((ans.TotalCost)*100) / 100
	ans.TotalCostTax = math.Floor((ans.TotalCostTax)*100) / 100

	return ans
}
