package broker

import (
	"context"
	"fmt"
	"micro/config"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/gnatsd/server"
	natsserver "github.com/nats-io/nats-server/test"
	nats "github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// So that we can pass tests and benchmarks...
type tLogger interface {
	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// TestLogger
type TestLogger tLogger

// RunDefaultServer will run a server on the default port.
func RunDefaultServer() *server.Server {
	return RunServerOnPort(nats.DefaultPort)
}

// RunServerOnPort will run a server on the given port.
func RunServerOnPort(port int) *server.Server {
	opts := natsserver.DefaultTestOptions
	opts.Port = port
	return RunServerWithOptions(opts)
}

// RunServerWithOptions will run a server with the given options.
func RunServerWithOptions(opts server.Options) *server.Server {
	return natsserver.RunServer(&opts)
}

// RunServerWithConfig will run a server with the given configuration file.
func RunServerWithConfig(configFile string) (*server.Server, *server.Options) {
	return natsserver.RunServerWithConfig(configFile)
}

func NewEncodedConnection(t tLogger) *nats.Conn {
	conn := NewDefaultConnection(t)
	nc, _ = nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	return conn
}

// NewDefaultConnection
func NewDefaultConnection(t tLogger) *nats.Conn {
	return NewConnection(t, nats.DefaultPort)
}

// NewConnection forms connection on a given port.
func NewConnection(t tLogger, port int) *nats.Conn {
	url := fmt.Sprintf("nats://127.0.0.1:%d", port)
	nc, err := nats.Connect(url)
	if err != nil {
		t.Fatalf("Failed to create default connection: %v\n", err)
		return nil
	}
	return nc
}
func TestDefaultConnection(t *testing.T) {
	s := RunDefaultServer()
	defer s.Shutdown()

	nc := NewDefaultConnection(t)
	nc.Close()
}

func TestConnect(t *testing.T) {
	s := RunDefaultServer()
	defer s.Shutdown()

	tests := []struct {
		step string
		conf config.Config
		err  error
	}{
		{
			step: "A",
			conf: config.Config{
				Nats: config.NATS{
					AllowReconnect: true,
					MaxReconnect:   5,
					ReconnectWait:  5,
					Timeout:        5,
					Endpoints:      nats.DefaultOptions.Servers,
					Username:       "root",
					Password:       "123456789",
					Encoder:        nats.JSON_ENCODER,
				},
			},
			err: nil,
		},
		{
			step: "B",
			conf: config.Config{
				Nats: config.NATS{
					Endpoints: []string{"https://worngIP:6666"},
					Encoder:   nats.JSON_ENCODER,
				},
			},
			err: nil,
		},
		{
			step: "C",
			conf: config.Config{
				Nats: config.NATS{
					AllowReconnect: true,
					MaxReconnect:   5,
					ReconnectWait:  5,
					Timeout:        5,
					Endpoints:      nats.DefaultOptions.Servers,
					Username:       "root",
					Password:       "123456789",
					Encoder:        "",
				},
			},
			err: fmt.Errorf("no encoder registered for ''"),
		},
		{
			step: "D",
			conf: config.Config{
				Nats: config.NATS{
					AllowReconnect: true,
					MaxReconnect:   -1,
					ReconnectWait:  -1,
					Timeout:        50,
					Endpoints:      nats.DefaultOptions.Servers,
					Username:       "root",
					Password:       "123456789",
					Encoder:        "json",
				},
			},
			err: fmt.Errorf("nats: no servers available for connection"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.step, func(t *testing.T) {

			// clear synced value
			once = sync.Once{}

			// try to connect
			err := Nats.Connect(tc.conf)
			if err != nil && tc.err != nil {
				assert.Equal(t, err.Error(), tc.err.Error())
			}
		})
	}

}

func TestGetConnection(t *testing.T) {
	s := RunDefaultServer()
	defer s.Shutdown()

	NewEncodedConnection(t)
	// defer conn.Close()
	defer nc.Close()

	tests := []struct {
		step string
		err  error
	}{
		{
			step: "A",
			err:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.step, func(t *testing.T) {

			conn := Nats.Conn()
			if conn == nil && tc.err == nil {
				assert.Errorf(t, fmt.Errorf("error in get connection"), "error in get connection")
				return
			}
		})
	}
}

func TestSendChan(t *testing.T) {
	s := RunDefaultServer()
	defer s.Shutdown()

	NewEncodedConnection(t)
	// defer conn.Close()
	defer nc.Close()

	type person struct {
		Name    string
		Address string
		Age     int
	}

	ch := make(chan interface{})

	if err := Nats.SendChan("hello", ch); err != nil {
		assert.Equal(t, err.Error(), "nats: argument needs to be a channel type")
	}
	me := &person{Name: "derek", Age: 22, Address: "85 Second St"}
	ch <- me

}

func TestSendByContext(t *testing.T) {
	subject := "test"
	s := RunDefaultServer()
	defer s.Shutdown()

	NewEncodedConnection(t)
	// defer conn.Close()
	defer nc.Close()

	Nats.Subscribe(subject, func(m *nats.Msg) {
		time.Sleep(200 * time.Millisecond)
		nc.Publish(m.Reply, []byte("NG"))
	})

	var resp interface{}
	if err := Nats.SendByContext(context.Background(), subject, []byte("world"), &resp); err != nil {
		assert.Error(t, fmt.Errorf("invalid error for send data with chan"))
		return
	}

	ctx, c := context.WithCancel(context.Background())
	c()
	if err := Nats.SendByContext(ctx, subject, []byte("world"), &resp); err != nil {
		assert.Contains(t, err.Error(), "context canceled")
	}

}

func TestPublish(t *testing.T) {
	subject := "test"

	s := RunDefaultServer()
	defer s.Shutdown()

	NewEncodedConnection(t)
	// defer conn.Close()
	defer nc.Close()

	if err := Nats.Publish(context.Background(), subject, "publish"); err != nil {
		assert.Error(t, fmt.Errorf("invalid error for publish data"))
		return
	}

	ch := make(chan string, 2)
	if err := Nats.Publish(context.Background(), subject, ch); err != nil {
		assert.Contains(t, err.Error(), "json: unsupported type: chan string")
	}
}

func TestRequestWithReply(t *testing.T) {
	subject := "test"
	testString := "Hello World!"
	chanTestString := make(chan string, 2)
	oReply := "foobar"

	s := RunDefaultServer()
	defer s.Shutdown()

	NewEncodedConnection(t)
	// defer conn.Close()
	defer nc.Close()

	nc.Subscribe(subject, func(subj, reply, s string) {
		if s != testString {
			t.Fatalf("Received test string of '%s', wanted '%s'\n", s, testString)
		}
		if subj != subject {
			t.Fatalf("Received subject of '%s', wanted '%s'\n", subj, subject)
		}
		if reply != oReply {
			t.Fatalf("Received reply of '%s', wanted '%s'\n", reply, oReply)
		}
	})

	if err := Nats.RequestWithReply(subject, testString, oReply); err != nil {
		assert.Contains(t, err.Error(), "connection closed")
	}

	if err := Nats.RequestWithReply(subject, oReply, testString); err != nil {
		assert.Contains(t, err.Error(), "connection closed")
	}

	if err := Nats.RequestWithReply(subject, chanTestString, oReply); err != nil {
		assert.Contains(t, err.Error(), "json: unsupported type: chan string")
	}

}

func TestSubscribe(t *testing.T) {
	subject := "test"

	s := RunDefaultServer()
	defer s.Shutdown()

	NewEncodedConnection(t)
	// defer conn.Close()
	defer nc.Close()

	if _, err := Nats.Subscribe(subject, nil); err != nil {
		assert.Error(t, err)
		return
	}

	if _, err := Nats.Subscribe("", func(resp *nats.Msg) {}); err != nil {
		assert.Contains(t, err.Error(), "invalid subject")
	}
}

func TestRecvChan(t *testing.T) {
	subject := "test"

	s := RunDefaultServer()
	defer s.Shutdown()

	NewEncodedConnection(t)
	// defer conn.Close()
	defer nc.Close()

	ch := make(chan interface{}, 2)
	if _, err := Nats.RecvChan(subject, ch); err != nil {
		assert.Error(t, err)
	}

	if _, err := Nats.RecvChan("", ch); err != nil {
		assert.Contains(t, err.Error(), "invalid subject")
	}

}

func TestRecvGroup(t *testing.T) {
	subject := "test"
	queue := "test"

	s := RunDefaultServer()
	defer s.Shutdown()

	NewEncodedConnection(t)
	// defer conn.Close()
	defer nc.Close()

	if _, err := Nats.RecvGroup(subject, queue, func(resp interface{}) {}); err != nil {
		assert.Error(t, err)
	}

	if _, err := Nats.RecvGroup(subject, queue, nil); err != nil {
		assert.Contains(t, err.Error(), "Handler required for EncodedConn Subscription")
	}

	if _, err := Nats.RecvGroup(subject, "", func(resp interface{}) {}); err != nil {
		assert.Contains(t, err.Error(), "nats: invalid queue name")
	}

	if _, err := Nats.RecvGroup("", queue, func(resp interface{}) {}); err != nil {
		assert.Contains(t, err.Error(), "invalid subject")
	}

}

func TestErrorReporter(t *testing.T) {
	subject := "test"

	s := RunDefaultServer()
	defer s.Shutdown()

	url := fmt.Sprintf("nats://127.0.0.1:%d", nats.DefaultPort)

	options := nats.Options{
		Url:          url,
		AsyncErrorCB: Nats.errorReporter(zap.NewExample()),
	}

	conn, err := options.Connect()
	if err != nil {
		t.Fatalf("Failed to create default connection: %v\n", err)
	}
	nc, _ = nats.NewEncodedConn(conn, nats.JSON_ENCODER)

	if err := Nats.Publish(context.Background(), subject, "publish"); err != nil {
		assert.Error(t, fmt.Errorf("invalid error for publish data"))
		return
	}

	if err := nc.PublishRequest("foo", "bar", "foo"); err != nil {
		assert.Error(t, err, "Expected an error")
	}

	if err := nc.Request("foo", "foo", nil, 2*time.Second); err != nil {
		assert.Error(t, err, "Expected an error")
	}

	hand := Nats.errorReporter(zap.NewExample())
	hand(conn, &nats.Subscription{
		Subject: subject,
		Queue:   "test",
	}, fmt.Errorf("ERROR"))
}
