package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

// User Structure
type User struct {
	ID int
	Name string
}

func renderView(w http.ResponseWriter, filepath string, data interface{}) {
	var tmpl, err = template.ParseFiles(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}	

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}	
}

func readJSON(filepath string)[]byte{
	jsonFile, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
	}	
	defer jsonFile.Close()

	byteResult, _ := ioutil.ReadAll(jsonFile)
	return byteResult
}

func checkInCart(sku string, cart map[string]interface{}) int {
	for key, value := range cart {
		if (key == "Post_skuqty") {
			for _, value := range value.([]string) {				
				if (strings.Contains(value, sku)) {
					strnumber := strings.Replace(value, sku+"-", "", -1)					
					n, _ := strconv.Atoi(strnumber)
					return n
				}				
			}			
		}		
	}
	return 0
}

func calculateDisc(cart map[string]interface{}) float64 {
	disc := 0.0
	data := map[string]interface{}{}
	bytePromotionRules := readJSON("db/promotion_rules.json")
	json.Unmarshal([]byte(bytePromotionRules), &data)
	byteProduct := readJSON("db/products.json")	
	json.Unmarshal([]byte(byteProduct), &data)
	
	rules := data["promotion_rules"].([]interface{})
	for _, value := range rules {
		rule, _ := value.(map[string]interface{})
		
		var numMustHaveInCart int
		numMustHaveInCart = 0
		var mustHaveQty int
		mustHaveQty = 0
		var mustHaveMinQty int
		mustHaveMinQty = 0
		var mustHaveSku string
		var eligible bool			
		eligible = false
		if (rule["must_have"] != nil) {
			mustHaves := rule["must_have"]
			val, _ := mustHaves.(map[string]interface{})
			mustHaveSku = val["sku"].(string)
			mustHaveQty = 0
			mustHaveMinQty = 0

			if (val["qty"] != nil) {
				mustHaveQty = int(val["qty"].(float64))
			}
			
			if (val["min_qty"] != nil) {
				mustHaveMinQty = int(val["min_qty"].(float64))
			}
			
			numMustHaveInCart = checkInCart(mustHaveSku, cart)
			if (numMustHaveInCart > 0) {					
				if (mustHaveQty > 0 || (mustHaveMinQty > 0 && numMustHaveInCart >= mustHaveMinQty)) {
					eligible = true
				}				
			}
		}

		// Reward Promo Get Free
		if (eligible && rule["get_free"] != nil) {
			getFrees := rule["get_free"]
			val, _ := getFrees.(map[string]interface{})
			freeSku := val["sku"].(string)
			freeQty := 0

			if (val["qty"] != nil) {
				freeQty = int(val["qty"].(float64))
			}

			// Find a get_free product price
			freePrice := 0.0			
			for _, value := range data["products"].([]interface {}) {
				val, _ := value.(map[string]interface{})
				if (val["sku"] == freeSku) {
					freePrice = val["price"].(float64)
				}
			}

			numEligiblePromo := numMustHaveInCart/freeQty			
			numInCart := checkInCart(freeSku, cart)
			if (numInCart > 0 && numInCart >= freeQty) {
				var numPromo int
				if (numInCart <= numEligiblePromo) {
					numPromo = numInCart
				} else {
					numPromo = numEligiblePromo
				}
				disc = disc + float64(freePrice * float64(numPromo));
			} 			
		}
		
		// Reward Promo Count As
		if (eligible && rule["count_as"] != nil) {
			numEligiblePromo := math.Floor(float64(numMustHaveInCart/mustHaveQty))
			// Find a count_as product price
			countAsPrice := 0.0			
			for _, value := range data["products"].([]interface {}) {
				val, _ := value.(map[string]interface{})
				if (val["sku"] == mustHaveSku) {
					countAsPrice = val["price"].(float64)
				}
			}
			disc = disc + float64(countAsPrice * (float64(mustHaveQty) - rule["count_as"].(float64)) * numEligiblePromo);			
		}

		if (eligible && rule["percent_disc"] != nil) {
			// Find a count_as product price
			productPrice := 0.0			
			for _, value := range data["products"].([]interface {}) {
				val, _ := value.(map[string]interface{})
				if (val["sku"] == mustHaveSku) {
					productPrice = val["price"].(float64)
				}
			}

			factor := float64(rule["percent_disc"].(float64)/100) * float64(numMustHaveInCart);
			disc = disc + (productPrice * factor)			
		}
		// fmt.Print("disc: ")
		// fmt.Println(disc)		
	}
	return disc
}

func pageHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	var data = map[string]interface{} {
		"page": map[string]interface{} {
			"Title": "Amazonas - Online Shop",
		},
		"user": User{
			ID: 1,
			Name: "Joe",
		},	
	}
	byteProduct := readJSON("db/products.json")	
	json.Unmarshal([]byte(byteProduct), &data)
		
	renderView(w, path.Join("views", "index.html"), data)
}

func pageCheckout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}	
			
	var cart = map[string]interface{} {
		"Post_skuqty": r.Form["skuqty[]"],
		"Post_qty": r.Form["qty[]"],
		"Post_total": r.Form["total[]"],		
	}	
	
	subtotal := r.FormValue("edt-subtotal-amount")
	iSubtotal, _ := strconv.ParseFloat(subtotal, 64)
	disc := calculateDisc(cart)
	grandtotal := iSubtotal - disc	

	var data = map[string]interface{} {
		"page": map[string]interface{} {
			"Title": "Amazonas - Online Shop",			
		},
		"user": User{
			ID: 1,
			Name: "Joe",
		},
		"cart": cart,				
		"subtotal": subtotal,
		"subtotal_display": r.FormValue("edt-subtotal-amount-disp"),
		"disc": disc,
		"grandtotal": grandtotal,
	}	

	byteProduct := readJSON("db/products.json")
	json.Unmarshal([]byte(byteProduct), &data)	
	
	renderView(w, path.Join("views", "checkout.html"), data)
}

func routes() {
	http.HandleFunc("/", pageHome)
	http.HandleFunc("/checkout", pageCheckout)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	routes()
}