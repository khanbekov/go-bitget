package main

import (
	"context"
	"fmt"
	"log"

	"github.com/khanbekov/go-bitget/futures"
	"github.com/khanbekov/go-bitget/futures/trading"
)

func main() {
	// Initialize the futures client with your API credentials
	// NOTE: Use demo API keys for testing, production API keys for live trading
	client := futures.NewClient("your_api_key", "your_secret_key", "your_passphrase")

	// Example 1: Cancel orders by order IDs
	fmt.Println("Example 1: Batch cancel orders by order IDs")
	batchCancelService := trading.NewBatchCancelOrdersService(client)

	// Configure the service for USDT-M Futures
	result, err := batchCancelService.
		Symbol("BTCUSDT").
		ProductType(trading.ProductTypeUSDTFutures).
		MarginCoin("USDT").
		AddOrderId("order_id_1").
		AddOrderId("order_id_2").
		AddOrderId("order_id_3").
		Do(context.Background())

	if err != nil {
		log.Printf("Error cancelling orders: %v", err)
		return
	}

	fmt.Printf("Successfully cancelled %d orders\n", len(result.SuccessList))
	for _, success := range result.SuccessList {
		fmt.Printf("  - Cancelled order: %s (client: %s)\n", success.OrderId, success.ClientOrderId)
	}

	if len(result.FailureList) > 0 {
		fmt.Printf("Failed to cancel %d orders:\n", len(result.FailureList))
		for _, failure := range result.FailureList {
			fmt.Printf("  - Failed order: %s, reason: %s (code: %s)\n",
				failure.OrderId, failure.ErrorMsg, failure.ErrorCode)
		}
	}

	fmt.Println()

	// Example 2: Cancel orders by client order IDs
	fmt.Println("Example 2: Batch cancel orders by client order IDs")
	result2, err := trading.NewBatchCancelOrdersService(client).
		Symbol("ETHUSDT").
		ProductType(trading.ProductTypeUSDTFutures).
		AddClientOid("my_client_order_1").
		AddClientOid("my_client_order_2").
		Do(context.Background())

	if err != nil {
		log.Printf("Error cancelling orders by client ID: %v", err)
		return
	}

	fmt.Printf("Successfully cancelled %d orders by client ID\n", len(result2.SuccessList))

	fmt.Println()

	// Example 3: Mixed order IDs and client order IDs
	fmt.Println("Example 3: Mixed order IDs and client order IDs")

	// Create order items manually for more control
	orderItems := []trading.BatchCancelOrderItem{
		{OrderId: "system_order_1"},
		{ClientOid: "my_client_order_3"},
		{OrderId: "system_order_2", ClientOid: "my_client_order_4"}, // Both provided - orderId takes precedence
	}

	result3, err := trading.NewBatchCancelOrdersService(client).
		Symbol("ADAUSDT").
		ProductType(trading.ProductTypeUSDTFutures).
		OrderIdList(orderItems).
		Do(context.Background())

	if err != nil {
		log.Printf("Error cancelling mixed orders: %v", err)
		return
	}

	fmt.Printf("Successfully cancelled %d mixed orders\n", len(result3.SuccessList))

	fmt.Println()

	// Example 4: Cancel all orders for a product type (without specifying orderIdList)
	fmt.Println("Example 4: Cancel all orders for COIN-M Futures")
	result4, err := trading.NewBatchCancelOrdersService(client).
		ProductType(trading.ProductTypeCoinFutures).
		Do(context.Background())

	if err != nil {
		log.Printf("Error cancelling all COIN-M orders: %v", err)
		return
	}

	fmt.Printf("Successfully cancelled %d COIN-M orders\n", len(result4.SuccessList))

	fmt.Println("Batch cancel examples completed!")
}
