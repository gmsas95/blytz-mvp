package main

import (
	"context"
	"fmt"
	"log"

	"go.uber.org/zap"
	firebase "github.com/blytz/auction-service/pkg/firebase"
)

// Example: How to integrate Firebase with your auction service
func main() {
	// Initialize logger
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	// Create Firebase client
	firebaseClient := firebase.NewClient(logger)

	fmt.Println("🚀 Firebase Integration Example")
	fmt.Println("================================")

	// Example 1: Create a user
	fmt.Println("\n1. Creating user...")
	userResp, err := firebaseClient.CreateUser(ctx,
		"buyer@example.com",
		"securepassword123",
		"John Buyer",
		"+1234567890")

	if err != nil {
		log.Printf("❌ Failed to create user: %v", err)
	} else {
		fmt.Printf("✅ User created: %s\n", userResp.UID)
	}

	// Example 2: Create an auction
	fmt.Println("\n2. Creating auction...")
	auctionData := firebase.AuctionData{
		Title:         "Vintage Omega Watch",
		Description:   "Authentic vintage Omega Seamaster from 1960s",
		StartingPrice: 1500.00,
		Duration:      48, // 48 hours
		Category:      "watches",
		Images:        []string{"omega1.jpg", "omega2.jpg", "omega3.jpg"},
	}

	auctionResp, err := firebaseClient.CreateAuction(ctx, auctionData)
	if err != nil {
		log.Printf("❌ Failed to create auction: %v", err)
	} else {
		fmt.Printf("✅ Auction created: %s\n", auctionResp.AuctionID)
	}

	// Example 3: Place a bid
	fmt.Println("\n3. Placing bid...")
	bidData := firebase.BidData{
		AuctionID: auctionResp.AuctionID, // Use the auction we just created
		Amount:    1600.00,
	}

	bidResp, err := firebaseClient.PlaceBid(ctx, bidData)
	if err != nil {
		log.Printf("❌ Failed to place bid: %v", err)
	} else {
		fmt.Printf("✅ Bid placed: %s for $%.2f\n", bidResp.BidID, bidData.Amount)
	}

	// Example 4: Create payment intent
	fmt.Println("\n4. Creating payment intent...")
	paymentResp, err := firebaseClient.CreatePaymentIntent(ctx,
		1600.00, // Bid amount
		auctionResp.AuctionID,
		bidResp.BidID)

	if err != nil {
		log.Printf("❌ Failed to create payment intent: %v", err)
	} else {
		fmt.Printf("✅ Payment intent created: %s\n", paymentResp.PaymentIntentID)
		fmt.Printf("   Client secret: %s...\n", paymentResp.ClientSecret[:10])
	}

	// Example 5: Send notification
	fmt.Println("\n5. Sending notification...")
	notifData := map[string]string{
		"auctionId": auctionResp.AuctionID,
		"bidId":     bidResp.BidID,
		"type":      "new_bid",
	}

	notifResp, err := firebaseClient.SendNotification(ctx,
		"auction-host-id",
		"New Bid Alert!",
		"Someone placed a bid on your auction",
		notifData)

	if err != nil {
		log.Printf("❌ Failed to send notification: %v", err)
	} else {
		fmt.Printf("✅ Notification sent: %s\n", notifResp.MessageID)
	}

	// Example 6: End auction
	fmt.Println("\n6. Ending auction...")
	endResp, err := firebaseClient.EndAuction(ctx, auctionResp.AuctionID)
	if err != nil {
		log.Printf("❌ Failed to end auction: %v", err)
	} else {
		if endResp.WinnerID != nil {
			fmt.Printf("✅ Auction ended with winner: %s, winning bid: $%.2f\n",
				*endResp.WinnerID, endResp.WinningBid)
		} else {
			fmt.Println("✅ Auction ended with no bids")
		}
	}

	fmt.Println("\n🎉 Integration example completed!")
	fmt.Println("\n💡 Next steps:")
	fmt.Println("1. Import firebase package in your handlers")
	fmt.Println("2. Create firebase.NewClient(logger) in your service constructors")
	fmt.Println("3. Call Firebase functions when Redis operations need persistence")
	fmt.Println("4. Test with your existing auction service")
}