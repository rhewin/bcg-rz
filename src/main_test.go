package main

import (
	"math"
	"testing"
)

func TestCalculateDisc(t *testing.T) {	
	var subtotal float64
	var total []float64
	var disc float64
	var grandtotal float64
	// Scenario 1: Test Macbook Pro Promotion
	cart := map[string]interface{}{
		"Post_qty": []int{0, 1, 0, 1}, 
		"Post_skuqty": []string{"120P90-0", "43N23P-1", "A304SD-0", "234234-1"}, // Must in format: sku-qty
		"Post_total": []float64{49.99, 5399.99, 109.50, 30.00},
	}		

	subtotal = 0.0	
	total = cart["Post_total"].([]float64)	
	for key, value := range cart["Post_qty"].([]int) {				
		subtotal = subtotal + float64(value) * total[key]
	}					
	disc = math.Round(calculateDisc(cart)*100)/100
	grandtotal = subtotal - disc
	if (grandtotal != 5399.99) {
		t.Errorf("Promotion Buy MacBook Pro, Free Raspberry Pi B Failed! Expected Grand Total 5399.99, but got %v", grandtotal)
	} else {
		t.Logf("Promotion Buy MacBook Pro, Free Raspberry Pi B Success! Expected Grand Total 5399.99, but got %v", grandtotal);
	}
	// fmt.Println(subtotal)
	// fmt.Println(disc)
	// fmt.Println(grandtotal)

	// Scenario 2: Test Google Homes Promotion
	cart = map[string]interface{}{
		"Post_qty": []int{3, 0, 0, 0}, 
		"Post_skuqty": []string{"120P90-3", "43N23P-0", "A304SD-0", "234234-0"}, 
		"Post_total": []float64{49.99, 5399.99, 109.50, 30.00},
	}

	subtotal = 0.0	
	total = cart["Post_total"].([]float64)	
	for key, value := range cart["Post_qty"].([]int) {				
		subtotal = subtotal + float64(value) * total[key]
	}					
	disc = math.Round(calculateDisc(cart)*100)/100
	grandtotal = math.Round((subtotal - disc)*100)/100
	if (grandtotal != 99.98) {
		t.Errorf("Promotion Buy 3 Google Homes, for the price 2 Failed! Expected Grand Total 99.98, but got %v", grandtotal)
	} else {
		t.Logf("Promotion Buy 3 Google Homes, for the price 2 Success! Expected Grand Total 99.98, but got %v", grandtotal);
	}
	// fmt.Println(subtotal)
	// fmt.Println(disc)
	// fmt.Println(grandtotal)

	// Scenario 3: Test Alexa Speaker Promotion
	cart = map[string]interface{}{
		"Post_qty": []int{0, 0, 3, 0}, 
		"Post_skuqty": []string{"120P90-0", "43N23P-0", "A304SD-3", "234234-0"}, 
		"Post_total": []float64{49.99, 5399.99, 109.50, 30.00},
	}

	subtotal = 0.0	
	total = cart["Post_total"].([]float64)	
	for key, value := range cart["Post_qty"].([]int) {				
		subtotal = subtotal + float64(value) * total[key]
	}					
	disc = math.Round(calculateDisc(cart)*100)/100
	grandtotal = math.Round((subtotal - disc)*100)/100
	if (grandtotal != 295.65) {
		t.Errorf("Promotion Buy more than 3 Alexa Speakers will have 10 percent discount on all Alexa Speakers Failed! Expected Grand Total 295.65, but got %v", grandtotal)
	} else {
		t.Logf("Promotion Buy more than 3 Alexa Speakers will have 10 percent discount on all Alexa Speakers Success! Expected Grand Total 295.65, but got %v", grandtotal)
	}
	// fmt.Println(subtotal)
	// fmt.Println(disc)
	// fmt.Println(grandtotal)
}