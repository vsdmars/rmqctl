package pkg

import (
	"github.com/streadway/amqp"
	cli "gopkg.in/urfave/cli.v1"
)

type (
	amqpConnectionType struct {
		Username string `validate:"required"`
		Password string `validate:"required"`
		Host     string `validate:"required"`
		Vhost    string `validate:"required"`
		Port     int    `validate:"required"`
		APIPort  int    `validate:"required"`
	}

	createQueueType struct {
		amqpConnectionType
		QueueName  string     `validate:"required"`
		Durable    bool       `validate:"-"`
		Autodelete bool       `validate:"-"`
		Ha         bool       `validate:"-"`
		HaMode     string     `validate:"oneof=all exactly nodes"`
		HaParam    string     `validate:"-"`
		HaSyncMode string     `validate:"oneof=manual automatic"`
		Args       amqp.Table `validate:"-"`
	}

	createExchangeType struct {
		amqpConnectionType
		ExchangeName string     `validate:"required"`
		Kind         string     `validate:"oneof=direct fanout topic headers"`
		Durable      bool       `validate:"-"`
		Autodelete   bool       `validate:"-"`
		Args         amqp.Table `validate:"-"`
	}

	createBindType struct {
		amqpConnectionType
		SourceExchangeName string     `validate:"required"`
		DestinationName    string     `validate:"required"`
		Key                string     `validate:"required"`
		Type               string     `validate:"oneof=queue exchange"`
		Args               amqp.Table `validate:"-"`
	}

	createUserType struct {
		amqpConnectionType
		RmqUsername string `validate:"required"`
		Tag         string `validate:"isdefault|oneof=management policymaker monitoring administrator"`
		RmqPassword string `validate:"required"`
	}

	createVhostType struct {
		amqpConnectionType
		VhostName string `validate:"required"`
		Tracing   bool   `validate:"-"`
	}

	listQueueType struct {
		amqpConnectionType
		QueueName string `validate:"-"`
		Formatter string `validate:"oneof=json plain"`
	}

	listExchangeType struct {
		amqpConnectionType
		ExchangeName string `validate:"-"`
		Formatter    string `validate:"oneof=json plain"`
	}

	listBindType struct {
		amqpConnectionType
		All       bool   `validate:"-"`
		Formatter string `validate:"oneof=json plain"`
	}

	listVhostType struct {
		amqpConnectionType
		VhostName string `validate:"-"`
		Formatter string `validate:"oneof=json plain"`
	}

	listNodeType struct {
		amqpConnectionType
		NodeName  string `validate:"-"`
		Formatter string `validate:"oneof=json plain"`
	}

	listPolicyType struct {
		amqpConnectionType
		PolicyName string `validate:"-"`
		All        bool   `validate:"-"`
		Formatter  string `validate:"oneof=json plain"`
	}

	listUserType struct {
		amqpConnectionType
		RmqUsername string `validate:"-"`
		Formatter   string `validate:"oneof=json plain"`
	}

	deleteQueueType struct {
		amqpConnectionType
		QueueName string `validate:"required"`
	}

	deleteExchangeType struct {
		amqpConnectionType
		ExchangeName string `validate:"required"`
	}

	deleteBindType struct {
		amqpConnectionType
		SourceExchangeName string     `validate:"required"`
		DestinationName    string     `validate:"required"`
		Key                string     `validate:"required"`
		Type               string     `validate:"oneof=queue exchange"`
		Args               amqp.Table `validate:"-"`
	}

	deletePolicyType struct {
		amqpConnectionType
		PolicyName string `validate:"required"`
	}

	deleteUserType struct {
		amqpConnectionType
		RmqUsername string `validate:"-"`
	}

	deleteVhostType struct {
		amqpConnectionType
		VhostName string `validate:"required"`
	}

	publishType struct {
		amqpConnectionType
		ExchangeName string          `validate:"required"`
		Key          string          `validate:"required"`
		Mandatory    bool            `validate:"-"`
		Immediate    bool            `validate:"-"`
		Message      amqp.Publishing `validate:"required"`
	}

	consumeType struct {
		amqpConnectionType
		QueueName string     `validate:"required"`
		AckType   string     `validate:"oneof=ack nack reject"`
		AutoAck   bool       `validate:"-"`
		NoWait    bool       `validate:"-"`
		Daemon    bool       `validate:"-"`
		Formatter string     `validate:"oneof=json plain"`
		Args      amqp.Table `validate:"-"`
	}

	updateUserType struct {
		amqpConnectionType
		RmqUsername string `validate:"required"`
		Tag         string `validate:"isdefault|oneof=management policymaker monitoring administrator"`
		RmqPassword string `validate:"required"`
	}

	updateVhostType struct {
		amqpConnectionType
		VhostName string `validate:"required"`
		Tracing   bool   `validate:"-"`
	}
)

func noSuchJob(ctx *cli.Context) error {
	return cli.NewExitError("No such command!", 1)
}
