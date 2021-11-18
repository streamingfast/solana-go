package ws

import (
	"errors"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"

	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jpillora/backoff"

	"go.uber.org/zap"
)

var (
	ErrNotConnected          = errors.New("websocket not connected")
	ErrUrlEmpty              = errors.New("url can not be empty")
	ErrUrlWrongScheme        = errors.New("websocket uri must start with ws or wss scheme")
	ErrUrlNamePassNotAllowed = errors.New("user name and password are not allowed in websocket uri")
	//ErrCantConnect           = errors.New("websocket can't connect")
)

type WsOpts func(dl *websocket.Dialer)

type Websocket struct {
	// Websocket ID
	Id uint64
	// Websocket Meta
	Meta map[string]interface{}

	Logger *zap.Logger

	Errors chan<- error

	Reconnect bool

	// default to 2 seconds
	ReconnectIntervalMin time.Duration
	// default to 30 seconds
	ReconnectIntervalMax time.Duration
	// interval, default to 1.5
	ReconnectIntervalFactor float64
	// default to 2 seconds
	HandshakeTimeout time.Duration
	// Verbose suppress connecting/reconnecting messages.
	Verbose bool

	// Cal function
	OnConnect    func(ws *Websocket)
	OnDisconnect func(ws *Websocket)

	OnConnectError    func(ws *Websocket, err error)
	OnDisconnectError func(ws *Websocket, err error)
	OnReadError       func(ws *Websocket, err error)
	OnWriteError      func(ws *Websocket, err error)

	dialer        *websocket.Dialer
	url           string
	requestHeader http.Header
	httpResponse  *http.Response
	mu            sync.Mutex
	dialErr       error
	isConnected   bool
	isClosed      bool

	*websocket.Conn
}

func (ws *Websocket) WriteJSON(v interface{}) error {
	err := ErrNotConnected
	if ws.IsConnected() {
		err = ws.Conn.WriteJSON(v)
		if err != nil {
			if ws.OnWriteError != nil {
				ws.OnWriteError(ws, err)
			}
			ws.closeAndReconnect()
		}
	}

	return err
}

func (ws *Websocket) WriteMessage(messageType int, data []byte) error {
	err := ErrNotConnected

	if ws.IsConnected() {
		ws.mu.Lock()
		defer ws.mu.Unlock()

		err = ws.Conn.WriteMessage(messageType, data)
		if err != nil {
			if ws.OnWriteError != nil {
				ws.OnWriteError(ws, err)
			}
			ws.closeAndReconnect()
		}
	}

	return err
}

func (ws *Websocket) ReadMessage() (messageType int, message []byte, err error) {
	err = ErrNotConnected

	if ws.IsConnected() {
		messageType, message, err = ws.Conn.ReadMessage()
		if err != nil {
			if ws.OnReadError != nil {
				ws.OnReadError(ws, err)
			}
			ws.closeAndReconnect()
		}
	}

	return messageType, message, err
}

func (ws *Websocket) Close() {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if ws.Conn != nil {
		err := ws.Conn.Close()
		if err == nil && ws.isConnected && ws.OnDisconnect != nil {
			ws.OnDisconnect(ws)
		}
		if err != nil && ws.OnDisconnectError != nil {
			ws.OnDisconnectError(ws, err)
		}
	}

	ws.isConnected = false
}

func (ws *Websocket) closeAndReconnect() {
	ws.Close()
	ws.Connect()
}

func (ws *Websocket) Dial(urlStr string) error {
	_, err := parseUrl(urlStr)
	if err != nil {
		return err
	}

	ws.url = urlStr
	ws.setDefaults()

	ws.dialer = &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: ws.HandshakeTimeout,
	}

	hs := ws.HandshakeTimeout
	go ws.Connect()

	// wait on first attempt
	time.Sleep(hs)

	return nil
}

func (ws *Websocket) Connect() {
	b := &backoff.Backoff{
		Min:    ws.ReconnectIntervalMin,
		Max:    ws.ReconnectIntervalMax,
		Factor: ws.ReconnectIntervalFactor,
		Jitter: true,
	}

	// seed rand for backoff
	rand.Seed(time.Now().UTC().UnixNano())

	for {
		nextInterval := b.Duration()

		wsConn, httpResp, err := ws.dialer.Dial(ws.url, ws.requestHeader)

		ws.mu.Lock()
		ws.Conn = wsConn
		ws.dialErr = err
		ws.isConnected = err == nil
		ws.httpResponse = httpResp
		ws.mu.Unlock()

		if err == nil {
			if ws.Verbose && ws.Logger != nil {
				ws.Logger.Info(fmt.Sprintf("Websocket[%d].Dial: connection was successfully established with %s\n", ws.Id, ws.url))
			}
			if ws.OnConnect != nil {
				ws.OnConnect(ws)
			}
			return
		} else {
			if ws.Verbose && ws.Logger != nil {
				ws.Logger.Error(fmt.Sprintf("Websocket[%d].Dial: can't connect to %s, will try again in %v\n", ws.Id, ws.url, nextInterval))
			}
			if ws.OnConnectError != nil {
				ws.OnConnectError(ws, err)
			}
		}

		time.Sleep(nextInterval)
	}
}

func (ws *Websocket) GetHTTPResponse() *http.Response {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	return ws.httpResponse
}

func (ws *Websocket) GetDialError() error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	return ws.dialErr
}

func (ws *Websocket) IsConnected() bool {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	return ws.isConnected
}

func (ws *Websocket) setDefaults() {
	if ws.ReconnectIntervalMin == 0 {
		ws.ReconnectIntervalMin = 2 * time.Second
	}

	if ws.ReconnectIntervalMax == 0 {
		ws.ReconnectIntervalMax = 30 * time.Second
	}

	if ws.ReconnectIntervalFactor == 0 {
		ws.ReconnectIntervalFactor = 1.5
	}

	if ws.HandshakeTimeout == 0 {
		ws.HandshakeTimeout = 2 * time.Second
	}
}

func parseUrl(urlStr string) (*url.URL, error) {
	if urlStr == "" {
		return nil, ErrUrlEmpty
	}
	u, err := url.Parse(urlStr)

	if err != nil {
		return nil, err
	}

	if u.Scheme != "ws" && u.Scheme != "wss" {
		return nil, ErrUrlWrongScheme
	}

	if u.User != nil {
		return nil, ErrUrlNamePassNotAllowed
	}

	return u, nil
}