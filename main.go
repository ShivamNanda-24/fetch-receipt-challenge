package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"net/http"
	"receipt-processor-challenge/schemas"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
)

var receiptStore map[string]schemas.Receipt // Map of receipt IDs to receipts
var recipetPoints map[string]int            // Map of receipt IDs to points

func main() {
	// Initialize the receipt store and receipt points
	receiptStore = make(map[string]schemas.Receipt)
	recipetPoints = make(map[string]int)

	r := gin.Default()

	// Implemented the two endpoints and run the server at localhost:8080
	r.POST("/receipts/process", processReceipt)

	r.GET("/receipts/:id/points", getPoints)

	r.Run("localhost:8080")
}

/***
 * Implemented the logic to process the receipt and return the receipt ID as a response
 * @param context The context of the request
 * @return The receipt ID as a response
 */
func processReceipt(context *gin.Context) {
	var receipt schemas.Receipt // The receipt variable to store the JSON data

	// Bind the JSON data from the request body to the receipt variable
	if err := context.ShouldBindJSON(&receipt); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "The receipt is invalid"}) // Return a 400 error if the receipt is invalid
		return
	}

	// Generate a unique receipt ID
	receiptID := generateReceiptID()

	// Store the recipt in a map with the receipt ID as the key
	receiptStore[receiptID] = receipt

	// Calculate points for the receipt
	recipetPoints[receiptID] = calculatePoints(receiptID)

	// Respond with a success message and the receipt ID
	context.JSON(http.StatusOK, schemas.ProcessResponse{ID: receiptID})
}

/***
 * Implemented the logic to get the points for a receipt and return the points as a response
 * @param context The context of the request
 * @return The points as a response
 */

func getPoints(context *gin.Context) {
	// Get the receipt ID from the URL parameter
	receiptID := context.Param("id")

	// Get the points for the receipt ID
	points := recipetPoints[receiptID[1:]]

	// Return the points as a response
	context.JSON(http.StatusOK, schemas.PointsResponse{Points: points})
}

/*** Helper Functions ***/

/***
 * Implemented the logic to generate a unique receipt ID
 * @return The receipt ID
 */

func generateReceiptID() string {

	// Generate a random 36 byte key based on the example
	keyBytes := make([]byte, 36)
	if _, err := rand.Read(keyBytes); err != nil { // Use rand to generate a random 36 byte key based on the pattern ""
		return ""
	}

	key := hex.EncodeToString(keyBytes)
	return key

}

/***
 * Implements the rules to calculate the points for a receipt.
 * @param receiptID The receipt ID
 * @return The points for the receipt
 */
func calculatePoints(receiptID string) int {

	receipt := receiptStore[receiptID]
	points := 0

	if _, exists := receiptStore[receiptID]; !exists { // Check if the receipt exists
		fmt.Println("Receipt does not exist")
		return 0
	}

	// Rule 1: One point for every alphanumeric character in the retailer name.
	countAlphanumeric := 0
	for _, char := range receipt.Retailer {
		if unicode.IsLetter(char) || unicode.IsDigit(char) { // Check if the character is alphanumeric
			countAlphanumeric++
		}
	}
	points += countAlphanumeric

	// Rule 2: 50 points if the total is a round dollar amount with no cents.
	if strings.HasSuffix(receipt.Total, ".00") {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25.
	totalFloat, err := strconv.ParseFloat(receipt.Total, 64)
	if err == nil && math.Mod(totalFloat, 0.25) == 0 {
		points += 25
	} else {
		fmt.Println("Error parsing total")
	}

	// Rule 4: 5 points for every two items on the receipt.
	numItems := len(receipt.Items)
	points += (numItems / 2) * 5

	// Rule 5: If the trimmed length of the item description is a multiple of 3,
	// multiply the price by 0.2 and round up to the nearest integer.
	for _, item := range receipt.Items {
		trimmedLen := len(strings.TrimSpace(item.ShortDescription)) // Trim the item description and get the length
		if trimmedLen%3 == 0 {
			priceFloat, err := strconv.ParseFloat(item.Price, 64) // Parse the price to a float
			if err == nil {
				points += int(math.Ceil(priceFloat * 0.2)) // Multiply the price by 0.2 and round up to the nearest integer
			}
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd.
	purchaseDate, err := time.Parse("2006-01-02", receipt.PurchaseDate) // Parse the purchase date to a time
	if err == nil && purchaseDate.Day()%2 != 0 {                        // Check if the day is odd
		points += 6
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
	if err == nil {
		if purchaseTime.After(time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC)) && // Check if the time is after 2:00pm and before 4:00pm
			purchaseTime.Before(time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)) {
			points += 10
		}
	}

	return points
}
