package pkg

import (
	"fmt"
	"net"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
	cli "gopkg.in/urfave/cli.v1"
)

// ---- Jobs ----

// miss C++ template...
// waiting for generic in Golang2 :-)

func noSuchJob(ctx *cli.Context) error {
	return cli.NewExitError("No such command!", 1)
}

func createQueueJob(ctx *cli.Context) error {
	data := createQueueType{}

	err := validateCreateQueue(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connect(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = createQueue(conn, &data)
	if err != nil {
		return err
	}

	return nil
}

func createExchangeJob(ctx *cli.Context) error {
	data := createExchangeType{}

	err := validateCreateExchange(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connect(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = createExchange(conn, &data)
	if err != nil {
		return err
	}

	return nil
}

func createBindJob(ctx *cli.Context) error {
	data := createBindType{}

	err := validateCreateBind(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connect(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = createBind(conn, &data)
	if err != nil {
		return err
	}

	return nil
}

func createBindExJob(ctx *cli.Context) error {
	data := createBindExType{}

	err := validateCreateBindEx(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connect(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = createBindEx(conn, &data)
	if err != nil {
		return err
	}

	return nil
}

func deleteQueueJob(ctx *cli.Context) error {
	data := deleteQueueType{}

	err := validateDeleteQueue(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connect(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = deleteQueue(conn, &data)
	if err != nil {
		return err
	}

	return nil
}

func deleteExchangeJob(ctx *cli.Context) error {
	data := deleteExchangeType{}

	err := validateDeleteExchange(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connect(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = deleteExchange(conn, &data)
	if err != nil {
		return err
	}

	return nil
}

func deleteBindJob(ctx *cli.Context) error {
	data := deleteBindType{}

	err := validateDeleteBind(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connect(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = deleteBind(conn, &data)
	if err != nil {
		return err
	}

	return nil
}

func deleteBindExJob(ctx *cli.Context) error {
	data := deleteBindExType{}

	err := validateDeleteBindEx(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connect(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = deleteBindEx(conn, &data)
	if err != nil {
		return err
	}

	return nil
}

func publishJob(ctx *cli.Context) error {
	data := publishType{}

	err := validatePublish(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connect(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = publishMsg(conn, &data)
	if err != nil {
		return err
	}

	return nil
}

func consumeJob(ctx *cli.Context) error {
	data := consumeType{}

	err := validateConsume(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connect(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = consumeMsg(conn, &data)
	if err != nil {
		return err
	}

	return nil
}

// ---- End Jobs ----

func connect(amqpconn *amqpConnectionType) (*amqp.Connection, error) {
	amqpURL := amqp.URI{Scheme: "amqp",
		Host:     amqpconn.Host,
		Username: amqpconn.Username,
		Password: "XXXXX",
		Port:     amqpconn.Port,
		Vhost:    amqpconn.Vhost,
	}

	logger.Debug("amqp connect URL", zap.String("amqp", amqpURL.String()))

	amqpURL.Password = amqpconn.Password

	// tcp connection timeout in 3 seconds.
	connection, err := amqp.DialConfig(amqpURL.String(),
		amqp.Config{
			Vhost: amqpconn.Vhost,
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, 3*time.Second)
			},
			Heartbeat: 10 * time.Second,
			Locale:    "en_US"})
	if err != nil {
		logger.Debug("Opening amqp connection failed.",
			zap.String("error", err.Error()))

		return nil, cli.NewExitError(err.Error(), 1)
	}

	return connection, nil
}

func createQueue(conn *amqp.Connection, data *createQueueType) error {
	channel, err := conn.Channel()
	if err != nil {
		logger.Debug("Opening amqp channel failed.",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	_, err = channel.QueueDeclare(
		data.QueueName,
		data.Durable,
		data.Autodelete,
		data.Exclusive,
		data.NoWait,
		data.Args,
	)
	if err != nil {
		logger.Debug("QueueDeclare failed, channel closed.",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	if data.Ha {
		return createQueueHaJob(data)
	}

	return nil
}

func createExchange(conn *amqp.Connection, data *createExchangeType) error {
	channel, err := conn.Channel()
	if err != nil {
		logger.Debug("Opening channel failed.",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	err = channel.ExchangeDeclare(
		data.ExchangeName,
		data.Kind,
		data.Durable,
		data.Autodelete,
		data.Internal,
		data.NoWait,
		data.Args,
	)
	if err != nil {
		logger.Debug("ExchangeDeclare failed, channel closed.",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func createBind(conn *amqp.Connection, data *createBindType) error {
	channel, err := conn.Channel()
	if err != nil {
		logger.Debug("Opening channel failed.",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	err = channel.QueueBind(
		data.QueueName,
		data.Key,
		data.ExchangeName,
		data.NoWait,
		data.Args,
	)
	if err != nil {
		logger.Debug("QueueBind failed, channel closed.",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func createBindEx(conn *amqp.Connection, data *createBindExType) error {
	channel, err := conn.Channel()
	if err != nil {
		logger.Debug("Opening channel failed.",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	err = channel.ExchangeBind(
		data.ToExchange,
		data.Key,
		data.FromExchange,
		data.NoWait,
		data.Args,
	)
	if err != nil {
		logger.Debug("ExchangeBind failed, channel closed.",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func deleteQueue(conn *amqp.Connection, data *deleteQueueType) error {
	channel, err := conn.Channel()
	if err != nil {
		logger.Debug("Opening channel failed.",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	_, err = channel.QueueInspect(data.QueueName)
	if err != nil {
		logger.Debug("QueueDelete failed, queue does not exist.",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	cnt, err := channel.QueueDelete(
		data.QueueName,
		data.IfUnuse,
		data.IfEmpty,
		data.NoWait,
	)
	if err != nil {
		logger.Debug("QueueDelete failed, channel closed.",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Printf("Number of message left: %d\n", cnt)
	return nil
}

func deleteExchange(conn *amqp.Connection, data *deleteExchangeType) error {
	channel, err := conn.Channel()
	if err != nil {
		logger.Debug("Opening channel failed.",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	err = channel.ExchangeDelete(
		data.ExchangeName,
		data.IfUnuse,
		data.NoWait,
	)
	if err != nil {
		logger.Debug("ExchangeDelete failed, channel closed.",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func deleteBind(conn *amqp.Connection, data *deleteBindType) error {
	channel, err := conn.Channel()
	if err != nil {
		logger.Debug("Opening channel failed.",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	_, err = channel.QueueInspect(data.QueueName)
	if err != nil {
		logger.Debug("QueueUnbind failed, queue does not exist.",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	err = channel.QueueUnbind(
		data.QueueName,
		data.Key,
		data.ExchangeName,
		data.Args,
	)
	if err != nil {
		logger.Debug("QueueUnbind failed, channel closed.",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func deleteBindEx(conn *amqp.Connection, data *deleteBindExType) error {
	channel, err := conn.Channel()
	if err != nil {
		logger.Debug("Opening channel failed.",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	err = channel.ExchangeUnbind(
		data.ToExchange,
		data.Key,
		data.FromExchange,
		data.NoWait,
		data.Args,
	)
	if err != nil {
		logger.Debug("ExchangeUnBind failed, channel closed.",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func publishMsg(conn *amqp.Connection, data *publishType) error {
	channel, err := conn.Channel()
	if err != nil {
		logger.Debug("Opening channel failed.",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	err = channel.Publish(
		data.ExchangeName,
		data.Key, // routing key
		data.Mandatory,
		data.Immediate,
		data.Message)
	if err != nil {
		logger.Debug("Publish failed, channel closed.",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

// consumeMsg consumes message from the cluster.
// it uses it's own formatter instead of generic ones due
// to looping.
func consumeMsg(conn *amqp.Connection, data *consumeType) error {
	channel, err := conn.Channel()
	if err != nil {
		logger.Debug("Opening channel failed.",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

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
