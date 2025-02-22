package receipt

import (
	"math"
	"regexp"
	"strings"
	"ticket-processor/internal/models"
	"time"
)

func CalculatePoints(receipt models.Receipt) int {
	points := 0

	points += calculateRetailerPoints(receipt.Retailer)
	points += calculateRoundDollarBonus(receipt.Total)
	points += calculateQuarterMultipleBonus(receipt.Total)
	points += calculateItemCountBonus(len(receipt.Items))
	points += calculateItemDescriptionPoints(receipt.Items)
	points += calculateOddDayBonus(receipt.PurchaseDate)
	points += calculatePurchaseTimeBonus(receipt.PurchaseTime)

	return points
}

func calculateRetailerPoints(retailer string) int {
	re := regexp.MustCompile("[a-zA-Z0-9]")
	matches := re.FindAllString(retailer, -1)
	return len(matches)
}

func calculateRoundDollarBonus(amount float64) int {
	if amount > 0 && amount == float64(int(amount)) {
		return 50
	}
	return 0
}

// very tricky one, as working with float numbers very painful
// I have experience with this kind of task, in e-commerce project
// for this particular case, lets assume our processor will always receive
// already correctly rounded numbers
// in other case solution will be more complex with comparing reminder with some epsilon
func calculateQuarterMultipleBonus(amount float64) int {
	if amount > 0 && math.Mod(amount*4, 1) == 0 {
		return 25
	}
	return 0
}

func calculateItemCountBonus(itemCount int) int {
	return (itemCount / 2) * 5
}

func calculateItemDescriptionPoints(items []models.Item) int {
	points := 0
	for _, item := range items {
		trimmedDescription := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDescription)%3 == 0 {
			points += int(math.Ceil(item.Price * 0.2))
		}
	}
	return points
}

func calculateOddDayBonus(purchaseDate string) int {
	d, err := time.Parse(time.DateOnly, purchaseDate)
	if err != nil {
		return 0
	}
	day := d.Day()

	if day%2 != 0 {
		return 6
	}
	return 0
}

func calculatePurchaseTimeBonus(purchaseTime string) int {
	t, err := time.Parse("15:04", purchaseTime)
	if err != nil {
		return 0
	}

	after2pm := t.Hour() >= 14 && t.Hour() < 16
	if after2pm {
		return 10
	}
	return 0
}
