package pkg

import (
	"fmt"
	"strconv"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"
	cli "gopkg.in/urfave/cli.v1"
)

func logValidation(err error) {
	for _, err := range err.(validator.ValidationErrors) {
		logger.Debug("validation error",
			zap.String("namespace", err.Namespace()),
			zap.String("field", err.Field()),
			zap.String("structNamespace", err.StructNamespace()),
			zap.String("structField", err.StructField()),
			zap.String("tag", err.Tag()),
			zap.String("actualTag", err.ActualTag()),
			zap.String("kind", err.Kind().String()),
			zap.String("type", err.Type().String()))
	}
}

func validateAmqp(ctx *cli.Context) (amqpConnectionType, error) {
	amqpData := amqpConnectionType{
		Username: ctx.GlobalString("username"),
		Password: ctx.GlobalString("password"),
		Host:     ctx.GlobalString("host"),
		Vhost:    ctx.GlobalString("vhost"),
		Port:     ctx.GlobalInt("port"),
		APIPort:  ctx.GlobalInt("apiport"),
	}

	v := validator.New()
	err := v.Struct(amqpData)
	if err != nil {
		logValidation(err)
		return amqpData, cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return amqpData, nil
}

func validateCreateQueue(ctx *cli.Context, d *createQueueType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.Durable = ctx.Bool("du")
	d.Autodelete = ctx.Bool("ad")
	d.Exclusive = ctx.Bool("exc")
	d.NoWait = ctx.Bool("nw")
	d.QueueName = ctx.Args().First()

	// HA variables
	d.Ha = ctx.Bool("ha")
	d.HaMode = ctx.String("hm")
	d.HaParam = ctx.String("hp")
	d.HaSyncMode = ctx.String("sm")

	// Check HA variables setting
	if d.Ha {
		if d.HaMode == "exactly" {
			if _, err := strconv.Atoi(d.HaParam); err != nil {
				logger.Debug("validation error, 'exactly' HA mode should have interger parameter",
					zap.String("HA MODE", d.HaMode),
					zap.String("HA Param", d.HaParam))
				return cli.NewExitError("command error, use --help to see the proper usage.", 1)
			}
		}
		if d.HaMode == "nodes" {
			if d.HaParam == "" {
				logger.Debug("validation error, 'nodes' HA mode should have node's name parameter",
					zap.String("HA MODE", d.HaMode),
					zap.String("HA Param", d.HaParam))
				return cli.NewExitError("command error, use --help to see the proper usage.", 1)
			}
		}
		if d.HaMode == "all" {
			if d.HaParam != "" {
				logger.Debug("validation error, 'all' HA mode should not have parameter",
					zap.String("HA MODE", d.HaMode),
					zap.String("HA Param", d.HaParam))
				return cli.NewExitError("command error, use --help to see the proper usage.", 1)
			}
		}
	}

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validateCreateExchange(ctx *cli.Context, d *createExchangeType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.Durable = ctx.Bool("du")
	d.Autodelete = ctx.Bool("ad")
	d.Internal = ctx.Bool("exc")
	d.NoWait = ctx.Bool("nw")
	d.Kind = ctx.String("type")
	d.ExchangeName = ctx.Args().First()

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validateCreateBind(ctx *cli.Context, d *createBindType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.NoWait = ctx.Bool("nw")
	d.QueueName = ctx.Args().First()
	d.ExchangeName = ctx.Args().Get(1)
	d.Key = ctx.Args().Get(2)

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validateCreateBindEx(ctx *cli.Context, d *createBindExType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.NoWait = ctx.Bool("nw")
	d.FromExchange = ctx.Args().First()
	d.ToExchange = ctx.Args().Get(1)
	d.Key = ctx.Args().Get(2)

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validateCreateUser(ctx *cli.Context, d *createUserType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.RmqUsername = ctx.Args().First()
	d.RmqPassword = ctx.Args().Get(1)
	d.Tag = ctx.String("tag")

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validateCreateVhost(ctx *cli.Context, d *createVhostType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.VhostName = ctx.Args().First()
	d.Tracing = ctx.Bool("trace")

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validateListQueue(ctx *cli.Context, d *listQueueType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.Formatter = ctx.String("o")
	d.QueueName = ctx.Args().First()

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validateListExchange(ctx *cli.Context, d *listExchangeType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.Formatter = ctx.String("o")
	d.ExchangeName = ctx.Args().First()

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validateListBind(ctx *cli.Context, d *listBindType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.Formatter = ctx.String("o")
	d.All = ctx.Bool("all")

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validateListVhost(ctx *cli.Context, d *listVhostType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.Formatter = ctx.String("o")
	d.VhostName = ctx.Args().First()

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validateListNode(ctx *cli.Context, d *listNodeType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.Formatter = ctx.String("o")
	d.NodeName = ctx.Args().First()

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validateListPolicy(ctx *cli.Context, d *listPolicyType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.All = ctx.Bool("all")
	d.PolicyName = ctx.Args().First()
	d.Formatter = ctx.String("o")

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validateListUser(ctx *cli.Context, d *listUserType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.RmqUsername = ctx.Args().First()
	d.Formatter = ctx.String("o")

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validateDeleteQueue(ctx *cli.Context, d *deleteQueueType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.NoWait = ctx.Bool("nw")
	d.QueueName = ctx.Args().First()

	if !ctx.Bool("force") {
		d.IfUnuse = true
		d.IfEmpty = true
	}

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validateDeleteExchange(ctx *cli.Context, d *deleteExchangeType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.NoWait = ctx.Bool("nw")
	d.ExchangeName = ctx.Args().First()

	if !ctx.Bool("force") {
		d.IfUnuse = true
	}

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validateDeleteBind(ctx *cli.Context, d *deleteBindType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.QueueName = ctx.Args().First()
	d.ExchangeName = ctx.Args().Get(1)
	d.Key = ctx.Args().Get(2)

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validateDeleteBindEx(ctx *cli.Context, d *deleteBindExType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.NoWait = ctx.Bool("nw")
	d.FromExchange = ctx.Args().First()
	d.ToExchange = ctx.Args().Get(1)
	d.Key = ctx.Args().Get(2)

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validateDeletePolicy(ctx *cli.Context, d *deletePolicyType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.PolicyName = ctx.Args().First()

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validateDeleteUser(ctx *cli.Context, d *deleteUserType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.RmqUsername = ctx.Args().First()

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validateDeleteVhost(ctx *cli.Context, d *deleteVhostType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.VhostName = ctx.Args().First()

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validatePublish(ctx *cli.Context, d *publishType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.Immediate = ctx.Bool("immediate")
	d.Mandatory = ctx.Bool("mandatory")
	d.ExchangeName = ctx.Args().First()
	d.Key = ctx.Args().Get(1)

	if msg := ctx.Args().Get(2); msg != "" {
		d.Message = amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(fmt.Sprintf("%s", msg)),
			DeliveryMode: amqp.Persistent,
			MessageId:    "!42!",
		}
	}

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validateConsume(ctx *cli.Context, d *consumeType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.AutoAck = ctx.Bool("autoack")
	d.QueueName = ctx.Args().First()
	d.AckType = ctx.String("acktype")
	d.NoWait = ctx.Bool("nw")
	d.Daemon = ctx.Bool("daemon")
	d.Formatter = ctx.String("o")

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validateUpdateUser(ctx *cli.Context, d *updateUserType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.RmqUsername = ctx.Args().First()
	d.RmqPassword = ctx.Args().Get(1)
	d.Tag = ctx.String("tag")

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validateUpdateVhost(ctx *cli.Context, d *updateVhostType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return nil
	}

	d.amqpConnectionType = amqpData
	d.VhostName = ctx.Args().First()
	d.Tracing = ctx.Bool("trace")

	v := validator.New()
	err = v.Struct(d)
	if err != nil {
		logValidation(err)
		return cli.NewExitError("command error, use --help to see the proper usage.", 1)
	}

	return nil
}
