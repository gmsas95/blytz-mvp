package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.uber.org/zap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	// Import your Firebase client (adjust path as needed)
	package integration

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	firebase "github.com/gmsas95/blytz-mvp/services/auction-service/pkg/firebase"
)

func TestAuctionFlow(t *testing.T) {
	ctx := context.Background()

	// Initialize Firebase
	app, err := firebase.NewFirebaseApp(ctx)
	assert.NoError(t, err, "Failed to initialize Firebase")

	// Create a test auction
	auctionData := firebase.AuctionData{
		Title:         "Test Auction",
		Description:   "This is a test auction",
		StartingPrice: 10.0,
		Duration:      1, // 1 hour
	}
	auction, err := app.CreateAuction(ctx, auctionData)
	assert.NoError(t, err, "Failed to create auction")
	assert.NotNil(t, auction, "Auction should not be nil")

	log.Printf("Created auction with ID: %s", auction.AuctionID)

	// Place a bid on the auction
	bidData := firebase.BidData{
		AuctionID: auction.AuctionID,
		Amount:    15.0,
	}
	bid, err := app.PlaceBid(ctx, bidData)
	assert.NoError(t, err, "Failed to place bid")
	assert.NotNil(t, bid, "Bid should not be nil")

	log.Printf("Placed bid with ID: %s", bid.BidID)

	// End the auction
	endAuctionResponse, err := app.EndAuction(ctx, auction.AuctionID)
	assert.NoError(t, err, "Failed to end auction")
	assert.NotNil(t, endAuctionResponse, "End auction response should not be nil")

	log.Printf("Ended auction with winner: %s", *endAuctionResponse.WinnerID)
}

)

// TestAuctionFlow tests the complete auction lifecycle
func TestAuctionFlow(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	// Initialize Firebase client
	firebaseClient := firebase.NewClient(logger)

	t.Run("Complete Auction Lifecycle", func(t *testing.T) {
		// Step 1: Create a user
		t.Run("Create User", func(t *testing.T) {
			userResp, err := firebaseClient.CreateUser(ctx,
				"testuser@example.com",
				"password123",
				"Test User",
				"+1234567890")

			require.NoError(t, err)
			assert.True(t, userResp.Success)
			assert.NotEmpty(t, userResp.UID)

			t.Logf("✅ User created: %s", userResp.UID)
		})

		// Step 2: Create an auction
		t.Run("Create Auction", func(t *testing.T) {
			auctionData := firebase.AuctionData{
				Title:         "Vintage Rolex Watch",
				Description:   "Authentic vintage Rolex Submariner",
				StartingPrice: 5000.00,
				Duration:      24, // 24 hours
				Category:      "watches",
				Images:        []string{"watch1.jpg", "watch2.jpg"},
			}

			auctionResp, err := firebaseClient.CreateAuction(ctx, auctionData)
			require.NoError(t, err)
			assert.True(t, auctionResp.Success)
			assert.NotEmpty(t, auctionResp.AuctionID)

			t.Logf("✅ Auction created: %s", auctionResp.AuctionID)
		})

		// Step 3: Get auction details
		t.Run("Get Auction Details", func(t *testing.T) {
			// This would use the auction ID from previous step
			details, err := firebaseClient.GetAuctionDetails(ctx, "test-auction-id")
			require.NoError(t, err)
			assert.True(t, details.Success)
			assert.NotNil(t, details.Auction)

			t.Logf("✅ Auction details retrieved: %d bids", details.TotalBids)
		})

		// Step 4: Place a bid
		t.Run("Place Bid", func(t *testing.T) {
			bidData := firebase.BidData{
				AuctionID: "test-auction-id",
				Amount:    5500.00,
			}

			bidResp, err := firebaseClient.PlaceBid(ctx, bidData)
			require.NoError(t, err)
			assert.True(t, bidResp.Success)
			assert.NotEmpty(t, bidResp.BidID)

			t.Logf("✅ Bid placed: %s for $%.2f", bidResp.BidID, bidData.Amount)
		})

		// Step 5: Create payment intent
		t.Run("Create Payment Intent", func(t *testing.T) {
			paymentResp, err := firebaseClient.CreatePaymentIntent(ctx,
				100.00, // $100 bid amount
				"test-auction-id",
				"test-bid-id")

			require.NoError(t, err)
			assert.True(t, paymentResp.Success)
			assert.NotEmpty(t, paymentResp.ClientSecret)
			assert.NotEmpty(t, paymentResp.PaymentIntentID)

			t.Logf("✅ Payment intent created: %s", paymentResp.PaymentIntentID)
		})

		// Step 6: Send notification
		t.Run("Send Notification", func(t *testing.T) {
			notificationData := map[string]string{
				"auctionId": "test-auction-id",
				"type":      "bid_placed",
			}

			notifResp, err := firebaseClient.SendNotification(ctx,
				"test-user-id",
				"New Bid Placed!",
				"Someone placed a bid on your auction",
				notificationData)

			require.NoError(t, err)
			assert.True(t, notifResp.Success)

			t.Logf("✅ Notification sent: %s", notifResp.MessageID)
		})

		// Step 7: End auction
		t.Run("End Auction", func(t *testing.T) {
			endResp, err := firebaseClient.EndAuction(ctx, "test-auction-id")
			require.NoError(t, err)
			assert.True(t, endResp.Success)

			if endResp.WinnerID != nil {
				t.Logf("✅ Auction ended with winner: %s, winning bid: $%.2f",
					*endResp.WinnerID, endResp.WinningBid)
			} else {
				t.Logf("✅ Auction ended with no bids")
			}
		})

		// Step 8: Send auction update notification
		t.Run("Send Auction Update", func(t *testing.T) {
			updateResp, err := firebaseClient.SendAuctionUpdate(ctx,
				"test-auction-id",
				"auction_ended",
				"Auction has ended!")

			require.NoError(t, err)
			assert.True(t, updateResp.Success > 0)

			t.Logf("✅ Auction update sent to %d participants", updateResp.NotificationsSent)
		})
	})
}

// TestPaymentFlow tests the complete payment process
func TestPaymentFlow(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	firebaseClient := firebase.NewClient(logger)

	t.Run("Complete Payment Process", func(t *testing.T) {
		// Step 1: Create payment intent
		paymentIntent, err := firebaseClient.CreatePaymentIntent(ctx,
			1000.00, // $1000
			"auction-123",
			"bid-456")

		require.NoError(t, err)
		assert.True(t, paymentIntent.Success)

		t.Logf("✅ Payment intent created: %s", paymentIntent.PaymentIntentID)

		// Step 2: Confirm payment (simulating successful payment)
		confirmResp, err := firebaseClient.ConfirmPayment(ctx,
			paymentIntent.PaymentIntentID,
			"auction-123",
			"bid-456")

		require.NoError(t, err)
		assert.True(t, confirmResp.Success)

		t.Logf("✅ Payment confirmed: %s", confirmResp.Status)
	})
}

// TestErrorHandling tests error scenarios
func TestErrorHandling(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	firebaseClient := firebase.NewClient(logger)

	t.Run("Invalid Auction Creation", func(t *testing.T) {
		// Try to create auction with invalid data
		invalidData := firebase.AuctionData{
			Title:         "", // Empty title should fail
			Description:   "Test",
			StartingPrice: -100, // Negative price should fail
			Duration:      0,    // Zero duration should fail
			Category:      "test",
		}

		_, err := firebaseClient.CreateAuction(ctx, invalidData)
		assert.Error(t, err)
		t.Logf("✅ Correctly rejected invalid auction data: %v", err)
	})

	t.Run("Invalid Bid Placement", func(t *testing.T) {
		// Try to place bid with invalid data
		invalidBid := firebase.BidData{
			AuctionID: "", // Empty auction ID should fail
			Amount:    -50, // Negative amount should fail
		}

		_, err := firebaseClient.PlaceBid(ctx, invalidBid)
		assert.Error(t, err)
		t.Logf("✅ Correctly rejected invalid bid: %v", err)
	})
}

// BenchmarkAuctionCreation benchmarks auction creation performance
func BenchmarkAuctionCreation(b *testing.B) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	firebaseClient := firebase.NewClient(logger)

	auctionData := firebase.AuctionData{
		Title:         "Benchmark Auction",
		Description:   "Performance testing auction",
		StartingPrice: 100.00,
		Duration:      24,
		Category:      "benchmark",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := firebaseClient.CreateAuction(ctx, auctionData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Helper function to run all tests
func RunAllIntegrationTests(t *testing.T) {
	t.Log("🚀 Starting Firebase Integration Tests")

	TestAuctionFlow(t)
	TestPaymentFlow(t)
	TestErrorHandling(t)

	t.Log("✅ All integration tests completed successfully!")
}