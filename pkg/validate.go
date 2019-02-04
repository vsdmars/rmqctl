package pkg

import (
	"strconv"

	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"
	cli "gopkg.in/urfave/cli.v1"
)

func logValidation(err error) {
	for _, err := range err.(validator.ValidationErrors) {
		logger.Debug(
			"validation error",
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

func validates(d interface{}) error {
	if err := validator.New().Struct(d); err != nil {
		logValidation(err)

		return cli.NewExitError(
			"command error, use --help to see the proper usage.", 1)
	}

	return nil
}

func validateAmqp(ctx *cli.Context) (amqpConnectionType, error) {
	var username string
	var password string
	var host string
	var vhost string
	var port int
	var aport int
	var tls bool

	if username = ctx.GlobalString("u"); len(username) == 0 {
		username = ctx.GlobalString("username")
	}
	if password = ctx.GlobalString("p"); len(password) == 0 {
		password = ctx.GlobalString("password")
	}
	if host = ctx.GlobalString("H"); len(host) == 0 {
		host = ctx.GlobalString("host")
	}
	if vhost = ctx.GlobalString("V"); len(vhost) == 0 {
		vhost = ctx.GlobalString("vhost")
	}
	if port = ctx.GlobalInt("P"); port == 0 {
		port = ctx.GlobalInt("port")
	}
	if aport = ctx.GlobalInt("A"); aport == 0 {
		aport = ctx.GlobalInt("apiport")
	}
	if tls = ctx.GlobalBool("T"); !tls {
		tls = ctx.GlobalBool("tls")
	}

	amqpData := amqpConnectionType{
		Username: username,
		Password: password,
		Host:     host,
		Vhost:    vhost,
		Port:     port,
		APIPort:  aport,
		TLS:      tls,
	}

	if err := validates(amqpData); err != nil {
		return amqpData, err
	}

	return amqpData, nil
}

func validateCreateQueue(ctx *cli.Context, d *createQueueType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return err
	}

	d.amqpConnectionType = amqpData
	d.Durable = ctx.Bool("du")
	d.Autodelete = ctx.Bool("ad")
	d.QueueName = ctx.Args().First()

	// HA variables
	d.Ha = ctx.Bool("HA")
	d.HaMode = ctx.String("HAMODE")
	d.HaParam = ctx.String("HAPARAM")
	d.HaSyncMode = ctx.String("HASYNC")

	// Check HA variables setting
	if d.Ha {
		if d.HaMode == "exactly" {
			if _, err := strconv.Atoi(d.HaParam); err != nil {
				logger.Debug(
					"validation error, 'exactly' HA mode should have interger parameter",
					zap.String("HA MODE", d.HaMode),
					zap.String("HA Param", d.HaParam),
				)

				return cli.NewExitError("command error, use --help to see the proper usage.", 1)
			}
		}
		if d.HaMode == "nodes" {
			if d.HaParam == "" {
				logger.Debug(
					"validation error, 'nodes' HA mode should have node's name parameter",
					zap.String("HA MODE", d.HaMode),
					zap.String("HA Param", d.HaParam),
				)

				return cli.NewExitError("command error, use --help to see the proper usage.", 1)
			}
		}
		if d.HaMode == "all" {
			if d.HaParam != "" {
				logger.Debug(
					"validation error, 'all' HA mode should not have parameter",
					zap.String("HA MODE", d.HaMode),
					zap.String("HA Param", d.HaParam),
				)

				return cli.NewExitError("command error, use --help to see the proper usage.", 1)
			}
		}
	}

	return validates(d)
}

func validateCreateExchange(ctx *cli.Context, d *createExchangeType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return err
	}

	d.amqpConnectionType = amqpData
	d.Durable = ctx.Bool("du")
	d.Autodelete = ctx.Bool("ad")
	d.Kind = ctx.String("type")
	d.ExchangeName = ctx.Args().First()

	return validates(d)
}

func validateCreateBind(ctx *cli.Context, d *createBindType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return err
	}

	d.amqpConnectionType = amqpData
	d.Type = ctx.String("type")
	d.SourceExchangeName = ctx.Args().First()
	d.DestinationName = ctx.Args().Get(1)
	d.Key = ctx.Args().Get(2)

	return validates(d)
}

func validateCreateUser(ctx *cli.Context, d *createUserType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return err
	}

	d.amqpConnectionType = amqpData
	d.RmqUsername = ctx.Args().First()
	d.RmqPassword = ctx.Args().Get(1)
	d.Tag = ctx.String("tag")

	return validates(d)
}

func validateCreateVhost(ctx *cli.Context, d *createVhostType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return err
	}

	d.amqpConnectionType = amqpData
	d.VhostName = ctx.Args().First()
	d.Tracing = ctx.Bool("trace")

	return validates(d)
}

func validateListQueue(ctx *cli.Context, d *listQueueType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return err
	}

	d.amqpConnectionType = amqpData
	d.Formatter = ctx.String("o")
	d.QueueName = ctx.Args().First()

	return validates(d)
}

func validateListExchange(ctx *cli.Context, d *listExchangeType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return err
	}

	d.amqpConnectionType = amqpData
	d.Formatter = ctx.String("o")
	d.ExchangeName = ctx.Args().First()

	return validates(d)
}

func validateListBind(ctx *cli.Context, d *listBindType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return err
	}

	d.amqpConnectionType = amqpData
	d.Formatter = ctx.String("o")
	d.All = ctx.Bool("all")

	return validates(d)
}

func validateListVhost(ctx *cli.Context, d *listVhostType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return err
	}

	d.amqpConnectionType = amqpData
	d.Formatter = ctx.String("o")
	d.VhostName = ctx.Args().First()

	return validates(d)
}

func validateListNode(ctx *cli.Context, d *listNodeType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return err
	}

	d.amqpConnectionType = amqpData
	d.Formatter = ctx.String("o")
	d.NodeName = ctx.Args().First()

	return validates(d)
}

func validateListPolicy(ctx *cli.Context, d *listPolicyType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return err
	}

	d.amqpConnectionType = amqpData
	d.All = ctx.Bool("all")
	d.PolicyName = ctx.Args().First()
	d.Formatter = ctx.String("o")

	return validates(d)
}

func validateListUser(ctx *cli.Context, d *listUserType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return err
	}

	d.amqpConnectionType = amqpData
	d.RmqUsername = ctx.Args().First()
	d.Formatter = ctx.String("o")

	return validates(d)
}

func validateDeleteQueue(ctx *cli.Context, d *deleteQueueType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return err
	}

	d.amqpConnectionType = amqpData
	d.QueueName = ctx.Args().First()

	return validates(d)
}

func validateDeleteExchange(ctx *cli.Context, d *deleteExchangeType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return err
	}

	d.amqpConnectionType = amqpData
	d.ExchangeName = ctx.Args().First()

	return validates(d)
}

func validateDeleteBind(ctx *cli.Context, d *deleteBindType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return err
	}

	d.amqpConnectionType = amqpData
	d.Type = ctx.String("type")
	d.SourceExchangeName = ctx.Args().First()
	d.DestinationName = ctx.Args().Get(1)
	d.Key = ctx.Args().Get(2)

	return validates(d)
}

func validateDeletePolicy(ctx *cli.Context, d *deletePolicyType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return err
	}

	d.amqpConnectionType = amqpData
	d.PolicyName = ctx.Args().First()

	return validates(d)
}

func validateDeleteUser(ctx *cli.Context, d *deleteUserType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return err
	}

	d.amqpConnectionType = amqpData
	d.RmqUsername = ctx.Args().First()

	return validates(d)
}

func validateDeleteVhost(ctx *cli.Context, d *deleteVhostType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return err
	}

	d.amqpConnectionType = amqpData
	d.VhostName = ctx.Args().First()

	return validates(d)
}

func validatePublish(ctx *cli.Context, d *publishType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return err
	}

	d.amqpConnectionType = amqpData
	d.Immediate = ctx.Bool("immediate")
	d.Mandatory = ctx.Bool("mandatory")
	d.ExchangeName = ctx.Args().First()
	d.Burst = ctx.Int("burst")
	d.Key = ctx.Args().Get(1)
	d.Message = ctx.Args().Get(2)

	mode := ctx.String("mode")

	if mode == "transient" {
		d.Mode = 0
	} else if mode == "persistent" {
		d.Mode = 2
	} else {
		d.Mode = 3
	}

	return validates(d)
}

func validateConsume(ctx *cli.Context, d *consumeType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return err
	}

	d.amqpConnectionType = amqpData
	d.AutoAck = ctx.Bool("autoack")
	d.QueueName = ctx.Args().First()
	d.AckType = ctx.String("acktype")
	d.NoWait = ctx.Bool("nw")
	d.Daemon = ctx.Bool("daemon")
	d.Formatter = ctx.String("o")

	return validates(d)
}

func validateUpdateUser(ctx *cli.Context, d *updateUserType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return err
	}

	d.amqpConnectionType = amqpData
	d.RmqUsername = ctx.Args().First()
	d.RmqPassword = ctx.Args().Get(1)
	d.Tag = ctx.String("tag")

	return validates(d)
}

func validateUpdateVhost(ctx *cli.Context, d *updateVhostType) error {
	amqpData, err := validateAmqp(ctx)
	if err != nil {
		return err
	}

	d.amqpConnectionType = amqpData
	d.VhostName = ctx.Args().First()
	d.Tracing = ctx.Bool("trace")

	return validates(d)
}
