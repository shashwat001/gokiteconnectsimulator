package kiteticker

import (
	"kiteconnectsimulator/models"
	"log"
	"net/url"
	"time"
)

type Ticker struct {
	callbacks     callbacks
	onSubscribeCB func(token []uint32) error
}

// New creates a new ticker instance.
func New(apiKey string, accessToken string) *Ticker {

	ticker := &Ticker{
		// ticker: NewMainTicker(apiKey, accessToken),
		callbacks: callbacks{},
	}

	return ticker
}

// SetAccessToken set access token.
func (t *Ticker) SetAccessToken(aToken string) {
	// t.ticker.accessToken = aToken
	return
}

// SetRootURL sets ticker root url.
func (t *Ticker) SetRootURL(u url.URL) {
	return
}

// SetConnectTimeout sets default timeout for initial connect handshake
func (t *Ticker) SetConnectTimeout(val time.Duration) {
	return
}

// SetAutoReconnect enable/disable auto reconnect.
func (t *Ticker) SetAutoReconnect(val bool) {
	return
}

// SetReconnectMaxDelay sets maximum auto reconnect delay.
func (t *Ticker) SetReconnectMaxDelay(val time.Duration) error {
	return nil
}

// SetReconnectMaxRetries sets maximum reconnect attempts.
func (t *Ticker) SetReconnectMaxRetries(val int) {
	return
}

// OnConnect callback.
func (t *Ticker) OnConnect(f func()) {
	t.callbacks.onConnect = f
}

// OnError callback.
func (t *Ticker) OnError(f func(err error)) {
	t.callbacks.onError = f
}

// OnClose callback.
func (t *Ticker) OnClose(f func(code int, reason string)) {
	t.callbacks.onClose = f
}

// OnMessage callback.
func (t *Ticker) OnMessage(f func(messageType int, message []byte)) {
	t.callbacks.onMessage = f
}

// OnReconnect callback.
func (t *Ticker) OnReconnect(f func(attempt int, delay time.Duration)) {
	t.callbacks.onReconnect = f
}

// OnNoReconnect callback.
func (t *Ticker) OnNoReconnect(f func(attempt int)) {
	t.callbacks.onNoReconnect = f
}

// OnTick callback.
func (t *Ticker) OnTick(f func(tick models.Tick)) {
	t.callbacks.onTick = f
}

// OnOrderUpdate callback.
func (t *Ticker) OnOrderUpdate(f func(order models.Order)) {
	t.callbacks.onOrderUpdate = f
}

// OnOrderUpdate callback.
func (t *Ticker) OnSubscribe(f func(token []uint32) error) {
	t.onSubscribeCB = f
}

// Close tries to close the connection gracefully. If the server doesn't close it
func (t *Ticker) Close() error {
	return nil
}

// Stop the ticker instance and all the goroutines it has spawned.
func (t *Ticker) Stop() {
	return
}

// Subscribe subscribes tick for the given list of tokens.
func (t *Ticker) Subscribe(tokens []uint32) error {
	return t.onSubscribeCB(tokens)
}

// Unsubscribe unsubscribes tick for the given list of tokens.
func (t *Ticker) Unsubscribe(tokens []uint32) error {
	return nil
}

// SetMode changes mode for given list of tokens and mode.
func (t *Ticker) SetMode(mode Mode, tokens []uint32) error {
	return nil
}

// Resubscribe resubscribes to the current stored subscriptions
func (t *Ticker) Resubscribe() error {
	return nil
}

func (t *Ticker) Serve() {
	select {}
}

func (t *Ticker) TriggerTick(tick models.Tick) {
	if t == nil {
		return
	}
	if t.callbacks.onTick != nil {
		t.callbacks.onTick(tick)
	}
}

func (t *Ticker) TriggerOrderUpdate(order models.Order) {
	if t == nil {
		log.Println("Ticker is nil, callback ignored")
		return
	}
	if t.callbacks.onOrderUpdate != nil {
		t.callbacks.onOrderUpdate(order)
	}
}

func (t *Ticker) TriggerConnect() {
	if t == nil {
		return
	}
	if t.callbacks.onConnect != nil {
		t.callbacks.onConnect()
	}
}
