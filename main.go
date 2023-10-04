package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"net/http"
	"receipt-processor/schemas"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
)

var receiptStore map[string]schemas.Receipt
var recipetPoints map[string]int

func main() {
	receiptStore = make(map[string]schemas.Receipt)
	recipetPoints = make(map[string]int)

	r := gin.Default()

	r.POST("/receipts/process", processReceipt)

	r.GET("/receipts/:id/points", getPoints)

	r.Run("localhost:8080")
}

func processReceipt(context *gin.Context) {
	var receipt schemas.Receipt

	if err := context.ShouldBindJSON(&receipt); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "The receipt is invalid"})
		return
	}

	receiptID := generateReceiptID()

	receiptStore[receiptID] = receipt
	calculatePoints(receiptID)

	context.JSON(http.StatusOK, schemas.ProcessResponse{ID: receiptID})
}

func getPoints(context *gin.Context) {
	// Get the receipt ID from the URL parameter
	receiptID := context.Param("id")

	// Calculate the points based on the rules
	points := recipetPoints[receiptID[1:]]

	// Return the points as a response
	context.JSON(http.StatusOK, schemas.PointsResponse{Points: points})
}

// Implement your logic to generate a unique receipt ID here
func generateReceiptID() string {
	keyBytes := make([]byte, 36)
	if _, err := rand.Read(keyBytes); err != nil {
		return ""
	}

	key := hex.EncodeToString(keyBytes)
	return key

}

// Implement your logic to calculate points based on the rules here
func calculatePoints(receiptID string) {
	receipt := receiptStore[receiptID]
	if _, exists := receiptStore[receiptID]; exists {
		fmt.Println("Receipt ID exists in receiptStore")
	} else {
		fmt.Println("Receipt ID does not exist in receiptStore")
		fmt.Println(receiptID)
	}

	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name.
	countAlphanumeric := 0
	for _, char := range receipt.Retailer {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
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
	}

	// Rule 4: 5 points for every two items on the receipt.
	numItems := len(receipt.Items)
	points += (numItems / 2) * 5

	// Rule 5: If the trimmed length of the item description is a multiple of 3,
	// multiply the price by 0.2 and round up to the nearest integer.
	for _, item := range receipt.Items {
		trimmedLen := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLen%3 == 0 {
			priceFloat, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				points += int(math.Ceil(priceFloat * 0.2))
			}
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd.
	purchaseDate, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err == nil && purchaseDate.Day()%2 != 0 {
		points += 6
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
	if err == nil {
		if purchaseTime.After(time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC)) &&
			purchaseTime.Before(time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)) {
			points += 10
		}
	}

	recipetPoints[receiptID] = points
}
