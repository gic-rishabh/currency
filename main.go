package main

import (
	"currency-Cnvt/model"
	"currencyCnvt/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var con = service.Currency{}

func init() {
	con.Server ="mongodb+srv://brks:iambrks@cluster0.iskflou.mongodb.net/?retryWrites=true&w=majority"
	con.Database = "currency"
    con.Collection = "currencyrates"

	con.Connect()
}

func addCurrencyExchange(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid method")
		return
	}

	var currency []model.Currency

	if err := json.NewDecoder(r.Body).Decode(&currency); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	if docs, err := con.Insert(currency); err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable To Insert Record")
	} else {
		respondWithJson(w, http.StatusAccepted, map[string]string{
			"message": strconv.Itoa(docs) + " Record Inserted Successfully",
		})
	}
}

func convertCurrency(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid method")
		return
	}

	var convert model.Converter

	if err := json.NewDecoder(r.Body).Decode(&convert); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	if convert.Amount == 0 {
		respondWithError(w, http.StatusBadRequest, "Please provide amount greater than 0 for conversion")
		return
	}
	//fmt.Println(convert)
	if docs, err := con.Convert(convert); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
	} else {
		respondWithJson(w, http.StatusAccepted, map[string]string{
			"message": "Amount In Dollar : " + fmt.Sprintf("%v", docs),
		})
	}
}
func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func main() {
	http.HandleFunc("/add-currency/", addCurrencyExchange)
	http.HandleFunc("/convert-currency/", convertCurrency)
	fmt.Println("server started at 9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}