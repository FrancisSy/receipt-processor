package main

import (
	"encoding/json"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Global map to store id and receipt information in memory
var receiptMap = make(map[string]Receipt)

/*
* Global map to store id and points information in memory
* The service is fast enough as it is but this is a nice to have
 */
var pointMap = make(map[string]int)

// Struct that represents item values in request body
type Item struct {
	ShortDescription string `json:"shortDescription" binding:"required"`
	Price            string `json:"price" binding:"required"`
}

// Struct that represents request body
type Receipt struct {
	Retailer     string `json:"retailer" binding:"required"`
	PurchaseDate string `json:"purchaseDate" binding:"required"`
	PurchaseTime string `json:"purchaseTime" binding:"required"`
	Items        []Item `json:"items" binding:"dive"`
	Total        string `json:"total" binding:"required"`
}

// Function to strip non-alphanumeric characters from byte array
func stripNonAlphanumericChars(s []byte) []byte {
	var i int
	for _, j := range s {
		if ('a' <= j && j <= 'z') ||
			('A' <= j && j <= 'Z') ||
			('0' <= j && j <= '9') {
			s[i] = j
			i++
		}
	}

	return s[:i]
}

/*
* Function to count the number of rewards points
* from the receipt based on the set rules of the challenge
*
* Using Decimals here as float values are not
* appropriate when dealing with currency representation
 */
func calculatePointsFromReceipt(r *Receipt) int {
	var points int

	// One point for every alphanumeric character in the retailer name.
	points += len(stripNonAlphanumericChars([]byte(r.Retailer)))

	// 50 points if the total is a round dollar amount with no cents.
	currTotal, _ := decimal.NewFromString(r.Total)
	floorTotal := decimal.Decimal.Floor(currTotal)

	if currTotal.Equals(floorTotal) {
		points += 50
	}

	// 25 points if the total is a multiple of 0.25.
	multiple := decimal.NewFromFloat(0.25)
	if currTotal.Mod(multiple).IsZero() {
		points += 25
	}

	// 5 points for every two items on the receipt.
	points += (len(r.Items) / 2) * 5

	// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	for _, i := range r.Items {
		trimmedStr := strings.TrimSpace(i.ShortDescription)

		if len(trimmedStr)%3 == 0 {
			price, err := decimal.NewFromString(i.Price)

			if err == nil {
				points += int(decimal.Decimal.Ceil(price.Mul(decimal.NewFromFloat(0.2))).IntPart())
			}
		}
	}

	// 6 points if the day in the purchase date is odd.
	dataAndTime := r.PurchaseDate + " " + r.PurchaseTime
	date, _ := time.Parse("2006-01-02 15:04", dataAndTime)
	if date.Day()%2 == 1 {
		points += 6
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	if 14 <= date.Hour() && 0 < date.Minute() && 59 >= date.Minute() && date.Hour() < 16 {
		points += 10
	}

	return points
}

/*
* Function to handle GET request from /receipts/{id}/points URI
*
* First, the function makes a check to the point map
* to see if there is already a store rewards point
* calculation for the provided ID. If so, that
* value is returned.
*
* Else if there is no value, the points for the ID
* is calculated, added to the point map, and returned
 */
func getRewardsById(id string) int {
	points, ok := pointMap[id] // checks point map first

	if ok { // if value is in point map, return that value
		return points
	} else { // else calculate the value, store in point map
		val, ok2 := receiptMap[id]

		if ok2 {
			points = calculatePointsFromReceipt(&val)
			pointMap[id] = points
		} else { // return -999 if id is not present in memory
			return math.MinInt
		}
	}

	return points
}

// Function to generate a new UUID for POST responses
func generateNewUUID() string {
	return uuid.New().String()
}

// Function to handle POST requests from /receipts/process URI
func postReceipt(r Receipt) string {
	id := generateNewUUID()
	receiptMap[id] = r
	return id
}

func main() {
	// initialize router to handle service requests at port 8080
	router := gin.Default()

	// GET handler for obtaining rewards points by ID
	router.GET("/receipts/:id/points", func(c *gin.Context) {
		id := c.Param("id") // obtain id parameter
		points := getRewardsById(id)

		if points == math.MinInt {
			c.String(http.StatusBadRequest, "No receipt found for that id")
		} else {
			c.JSON(http.StatusOK, gin.H{
				"points": points,
			})
		}
	})

	// POST handler for adding receipt information and returning a corresponding ID
	router.POST("/receipts/process", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body) // read body request parameter

		if err != nil { // if error occurs, return 400 message
			c.String(http.StatusBadRequest, "The receipt is invalid")
		} else {
			// parse json body to receipt struct
			var receipt Receipt
			err := json.Unmarshal(body, &receipt)

			if err != nil { // if error occurs, return 400 message
				c.String(http.StatusBadRequest, "The receipt is invalid")
			} else {
				c.JSON(http.StatusOK, gin.H{
					"id": postReceipt(receipt),
				})
			}
		}
	})

	// listen to requests on port 8080
	router.Run(":8080")
}
