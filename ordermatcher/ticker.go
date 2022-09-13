package ordermatcher

import (
	"fmt"
	"main/kiteconnectsimulator/models"
	kiteticker "main/kiteconnectsimulator/ticker"
	"os"
	"strconv"
	"time"
)

// Triggered when any error is raised
func onError(err error) {
	fmt.Println("Error: ", err)
}

// Triggered when websocket connection is closed
func onClose(code int, reason string) {
	fmt.Println("Close: ", code, reason)
}

// Triggered when connection is established and ready to send and accept data
func onConnect() {
	fmt.Println("Connected")
}

// Triggered when tick is recevived
func onTick(tick models.Tick) {
	fmt.Println("Tick: ", tick)
	handleTick(tick)
}

// Triggered when reconnection is attempted which is enabled by default
func onReconnect(attempt int, delay time.Duration) {
	fmt.Printf("Reconnect attempt %d in %fs\n", attempt, delay.Seconds())
}

// Triggered when maximum number of reconnect attempt is made and the program is terminated
func onNoReconnect(attempt int) {
	fmt.Printf("Maximum no of reconnect attempt reached: %d", attempt)
}

func Runticker(ticker *kiteticker.Ticker) {
	// Assign callbacks
	ticker.OnError(onError)
	ticker.OnClose(onClose)
	ticker.OnConnect(onConnect)
	ticker.OnReconnect(onReconnect)
	ticker.OnNoReconnect(onNoReconnect)
	ticker.OnTick(onTick)

	// Start the connection
	ticker.Serve()
}

// getEnv returns the value of the environment variable provided.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getEnvUint32 returns the value of the environment variable provided converted as Uint32.
func getEnvUint32(key string, fallback int) uint32 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.Atoi(value)
		if err != nil {
			return uint32(fallback)
		}
		return uint32(i)
	}
	return uint32(fallback)
}
