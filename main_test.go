package main

import (
	"bytes"
	"math"
	"testing"
)

func TestStripNonAlphanumericChars(t *testing.T) {
	expected := []byte("MMCornerMarket")
	actual := stripNonAlphanumericChars([]byte("M&M Corner Market"))
	if !bytes.Equal(actual, expected) {
		t.Fatalf("Result is incorrect. Expected: %s, Actual: %s", expected, actual)
	}
}

func TestStripNonAlphanumericChars2(t *testing.T) {
	expected := []byte("TraderJoes")
	actual := stripNonAlphanumericChars([]byte("Trader Joe's"))
	if !bytes.Equal(actual, expected) {
		t.Fatalf("Result is incorrect. Expected: %s, Actual: %s", expected, actual)
	}
}

func TestCalculatePointsFromReceipt(t *testing.T) {
	// mock a receipt for testing
	tmpReceipt := Receipt{
		Retailer:     "test retailer",
		PurchaseDate: "2023-10-07",
		PurchaseTime: "15:00",
		Items: []Item{
			{ShortDescription: "test description", Price: "1.00"},
			{ShortDescription: "test description", Price: "1.00"},
		},
		Total: "2.00",
	}

	expected := 108
	result := calculatePointsFromReceipt(&tmpReceipt)
	if result != expected {
		t.Fatalf("Result is incorrect. Expected: %d, Actual: %d", expected, result)
	}
}

func TestCalculatePointsFromReceipt2(t *testing.T) {
	// mock a receipt for testing
	tmpReceipt := Receipt{
		Retailer:     "test retailer",
		PurchaseDate: "2023-10-07",
		PurchaseTime: "15:00",
		Items: []Item{
			{ShortDescription: "divisible by three", Price: "1.00"},
			{ShortDescription: "divisible by three", Price: "1.00"},
		},
		Total: "2.00",
	}

	expected := 110
	result := calculatePointsFromReceipt(&tmpReceipt)
	if result != expected {
		t.Fatalf("Result is incorrect. Expected: %d, Actual: %d", expected, result)
	}
}

func TestGetRewardsByIdFromMap(t *testing.T) {
	// mock an id and receipt for testing
	tmpReceipt := Receipt{
		Retailer:     "test retailer",
		PurchaseDate: "2023-10-07",
		PurchaseTime: "15:00",
		Items:        []Item{{ShortDescription: "test description", Price: "1.00"}},
		Total:        "1.00",
	}

	tmpId := postReceipt(tmpReceipt)
	expected := calculatePointsFromReceipt(&tmpReceipt)
	pointMap[tmpId] = expected
	result := getRewardsById(tmpId)

	if result != expected {
		t.Fatalf("Result is incorrect. Expected: %d, Actual: %d", expected, result)
	} else {
		t.Logf("Result from function is %d and has been sourced from the point map", result)
	}

	receiptMap = make(map[string]Receipt) // reset receipt map
	pointMap = make(map[string]int)       // reset point map
}

func TestGetRewardsByIdCalculated(t *testing.T) {
	// mock an id and receipt for testing
	tmpReceipt := Receipt{
		Retailer:     "test retailer",
		PurchaseDate: "2023-10-07",
		PurchaseTime: "15:00",
		Items:        []Item{{ShortDescription: "test description", Price: "1.00"}},
		Total:        "1.00",
	}

	tmpId := postReceipt(tmpReceipt)
	result := getRewardsById(tmpId)
	expected := calculatePointsFromReceipt(&tmpReceipt)

	if result != expected {
		t.Fatalf("Result is incorrect. Expected: %d, Actual: %d", expected, result)
	} else {
		t.Logf("Result from function is %d and has been calculated", result)
	}

	receiptMap = make(map[string]Receipt) // reset receipt map
	pointMap = make(map[string]int)       // reset point map
}

func TestGetRewardsByIdNoId(t *testing.T) {
	// mock an id and receipt for testing
	tmpReceipt := Receipt{
		Retailer:     "test retailer",
		PurchaseDate: "2023-10-07",
		PurchaseTime: "15:00",
		Items:        []Item{{ShortDescription: "test description", Price: "1.00"}},
		Total:        "1.00",
	}

	tmpId := postReceipt(tmpReceipt)

	receiptMap = make(map[string]Receipt) // clear receipt map
	pointMap = make(map[string]int)       // clear point map

	result := getRewardsById(tmpId)

	if result != math.MinInt {
		t.Fatal("Result is invalid. Generated UUID and receipt is not present in receipt map")
	} else {
		t.Logf("Result from function is %d", math.MaxInt)
	}

	receiptMap = make(map[string]Receipt) // reset receipt map
	pointMap = make(map[string]int)       // reset point map
}

func TestGenerateNewUUID(t *testing.T) {
	result := generateNewUUID()
	if result == "" {
		t.Fatalf("Result is incorrect. Generated UUID is nil")
	} else {
		t.Logf("Result from function is %s", result)
	}
}

func TestPostReceipt(t *testing.T) {
	// mock a receipt for testing
	tmpReceipt := Receipt{
		Retailer:     "test retailer",
		PurchaseDate: "2023-10-07",
		PurchaseTime: "15:00",
		Items:        []Item{},
		Total:        "1.00",
	}
	result := postReceipt(tmpReceipt)

	if result == "" {
		t.Fatalf("Result is incorrect. Generated UUID is nil")
	}
	if len(receiptMap) == 0 {
		t.Fatalf("Result is invalid. Generated UUID and receipt is not present in receipt map")
	}

	receiptMap = make(map[string]Receipt) // reset receipt map
}
