package pkg

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"os/exec"
	"runtime"
	"syscall"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
	cli "gopkg.in/urfave/cli.v1"
)

const tcpTimeout = 3
const tlsTimeout = 5

func connect(amqpconn *amqpConnectionType) (*amqp.Connection, error) {
	var amqpURI amqp.URI
	var connection *amqp.Connection

	logit := func(amqpURI amqp.URI) {
		logger.Debug(
			"amqp URI",
			zap.String("service", "amqp"),
			zap.String("URI", amqpURI.String()),
		)
	}

	logtimeout := func(t int) {
		logger.Debug(
			"timeout",
			zap.String("service", "amqp"),
			zap.String("timeout", (time.Duration(t)*time.Second).String()),
		)
	}

	if amqpconn.TLS {
		certfile, keyfile, err := getCertPath()
		if err != nil {
			return nil, err
		}

		amqpURI = amqp.URI{Scheme: "amqps",
			Host:     amqpconn.Host,
			Username: amqpconn.Username,
			Password: "XXXXX",
			Port:     amqpconn.Port,
			Vhost:    amqpconn.Vhost,
		}

		logit(amqpURI)

		amqpURI.Password = amqpconn.Password

		cfg := &tls.Config{}

		cert, err := tls.LoadX509KeyPair(certfile, keyfile)
		if err != nil {
			logger.Debug(
				"load x509 key pair failed",
				zap.String("service", "amqp"),
			)

			return nil, err
		}

		cfg.Certificates = append(cfg.Certificates, cert)

		logtimeout(tlsTimeout)

		// tcp connection timeout in 5 seconds.
		conn, err := amqp.DialConfig(
			amqpURI.String(),
			amqp.Config{
				Vhost: amqpconn.Vhost,
				Dial: func(network, addr string) (net.Conn, error) {
					return net.DialTimeout(network, addr, time.Duration(tlsTimeout)*time.Second)
				},
				Heartbeat:       10 * time.Second,
				Locale:          "en_US",
				TLSClientConfig: cfg,
			},
		)
		if err != nil {
			logger.Debug(
				"amqp create connection failed",
				zap.String("service", "amqp"),
			)

			return nil, cli.NewExitError(err.Error(), 1)
		}

		connection = conn
	} else {
		amqpURI = amqp.URI{Scheme: "amqp",
			Host:     amqpconn.Host,
			Username: amqpconn.Username,
			Password: "XXXXX",
			Port:     amqpconn.Port,
			Vhost:    amqpconn.Vhost,
		}

		logit(amqpURI)

		amqpURI.Password = amqpconn.Password

		logtimeout(tcpTimeout)

		// tcp connection timeout in 3 seconds.
		conn, err := amqp.DialConfig(
			amqpURI.String(),
			amqp.Config{
				Vhost: amqpconn.Vhost,
				Dial: func(network, addr string) (net.Conn, error) {
					return net.DialTimeout(network, addr, time.Duration(tcpTimeout)*time.Second)
				},
				Heartbeat: 10 * time.Second,
				Locale:    "en_US"},
		)
		if err != nil {
			logger.Debug(
				"amqp create connection failed",
				zap.String("service", "amqp"),
			)

			return nil, cli.NewExitError(err.Error(), 1)
		}

		connection = conn
	}

	return connection, nil
}

func runExecutable(payload chan<- amqp.Publishing, data *publishType) {
	pr, pw := io.Pipe()
	scanner := bufio.NewScanner(pr)

	go func() {
		defer cancel()

		cmd := exec.CommandContext(ctx, data.Executable)
		cmd.Stdout = pw

		go func() {
			<-ctx.Done()
			pr.Close()
			pw.Close()
			cmd.Process.Signal(syscall.SIGTERM)
			// don't call cmd.Process.Release() due to cmd.Run() calls
			// Wait()
		}()

		if err := cmd.Run(); err != nil {
			logger.Debug(
				"Executable run error",
				zap.String("service", "amqp"),
				zap.String("error", err.Error()),
			)
		}
	}()

	for scanner.Scan() {
		b := scanner.Bytes()

		// beware that scanner.Bytes() will override underlying array
		payload <- amqp.Publishing{
			ContentType:  "text/plain",
			Body:         append(b[:0:0], b...),
			DeliveryMode: data.Mode,
			MessageId:    "!42!",
		}
	}
}

func burstPublish(
	gcount int,
	channel *amqp.Channel,
	payload <-chan amqp.Publishing,
	data *publishType) chan string {

	status := make(chan string)

	go func() {
		defer func() {
			close(status)
		}()

		for p := range payload {
			channel.Publish(
				data.ExchangeName,
				data.Key, // routing key
				data.Mandatory,
				data.Immediate,
				p,
			)
		}

		status <- fmt.Sprintf("publish goroutine number: %v done.", gcount)
	}()

	return status
}

func publishMsg(conn *amqp.Connection, data *publishType) error {
	channel, err := conn.Channel()
	if err != nil {
		logger.Debug(
			"amqp create channel failed",
			zap.String("service", "amqp"),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	var gcount int
	if len(data.Executable) != 0 {
		gcount = 1
	} else {
		gcount = int(float64(runtime.NumCPU()) * float64(0.9))
	}

	logger.Debug(
		"publish using number of goroutines",
		zap.String("service", "amqp"),
		zap.Int("goroutines", gcount),
	)

	payload := make(chan amqp.Publishing, 1000)

	go func() {
		defer close(payload)

		if len(data.Executable) != 0 {
			runExecutable(payload, data)
		} else {
			p := amqp.Publishing{
				ContentType:  "text/plain",
				Body:         []byte(data.Message),
				DeliveryMode: data.Mode,
				MessageId:    "!42!",
			}

			for i := 0; i < data.Burst; i++ {
				select {
				case <-ctx.Done():
					return
				default:
					payload <- p
				}
			}
		}
	}()

	var gatherStatus []chan string
	for i := 0; i < gcount; i++ {
		gatherStatus = append(
			gatherStatus,
			burstPublish(
				i,
				channel,
				payload,
				data),
		)
	}

	// ContextType:
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Complete_list_of_MIME_types

	for _, s := range gatherStatus {
		m := <-s
		logger.Debug(
			"publish goroutine done",
			zap.String("message", m),
		)
	}

	logger.Debug(
		"publish done",
		zap.String("service", "amqp"),
	)

	return nil
}

// consumeMsg consumes message from the cluster.
// it uses it's own formatter instead of generic ones due
// to looping.
func consumeMsg(conn *amqp.Connection, data *consumeType) error {
	channel, err := conn.Channel()
	if err != nil {
		logger.Debug(
			"amqp create channel failed",
			zap.String("service", "amqp"),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	if data.Daemon {
		return daemonConsumeF(channel, data)
	}

	return noneDaemonConsumeF(channel, data)
}

func ackFunction(d *amqp.Delivery, data *consumeType) {
	// Only acknowledge if no autoack.
	if !data.AutoAck {
		switch data.AckType {
		case "ack":
			d.Ack(true)
		case "nack":
			d.Nack(true, false)
		case "reject":
			d.Reject(false)
		}
	}
}
