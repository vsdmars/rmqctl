package pkg

import (
	"os"
	"path"
	"syscall"

	"github.com/urfave/cli/altsrc"
	"go.uber.org/zap"
	cli "gopkg.in/urfave/cli.v1"
)

const (
	configYaml = "rmqctl.conf"
)

const (
	generalAction = iota
	createQueueAction
	createExchangeAction
	createBindAction
	createPolicyAction
	createUserAction
	createVhostAction
	listQueueAction
	listExchangeAction
	listBindAction
	listVhostAction
	listNodeAction
	listPolicyAction
	listUserAction
	deleteQueueAction
	deleteExchangeAction
	deleteBindAction
	deletePolicyAction
	deleteUserAction
	deleteVhostAction
	updateUserAction
	updateVhostAction
	publishAction
	consumeAction
)

var (
	configDir = func() string {
		currentDir, _ := os.Getwd()
		return path.Join(currentDir, configYaml)
	}()
)

var (
	// add action:job here.
	actionMap = map[int]cli.ActionFunc{
		createQueueAction:    createQueueJob,
		createExchangeAction: createExchangeJob,
		createBindAction:     createBindJob,
		createUserAction:     createUserJob,
		createVhostAction:    createVhostJob,
		listQueueAction:      listQueueJob,
		listExchangeAction:   listExchangeJob,
		listBindAction:       listBindJob,
		listVhostAction:      listVhostJob,
		listNodeAction:       listNodeJob,
		listPolicyAction:     listPolicyJob,
		listUserAction:       listUserJob,
		deleteQueueAction:    deleteQueueJob,
		deleteExchangeAction: deleteExchangeJob,
		deleteBindAction:     deleteBindJob,
		deletePolicyAction:   deletePolicyJob,
		deleteUserAction:     deleteUserJob,
		deleteVhostAction:    deleteVhostJob,
		updateUserAction:     updateUserJob,
		updateVhostAction:    updateVhostJob,
		publishAction:        publishJob,
		consumeAction:        consumeJob,
	}
)

func actions(action int) cli.ActionFunc {
	if f, ok := actionMap[action]; ok {
		return f
	}

	return noSuchJob
}

func flags(f int) []cli.Flag {
	switch f {
	case generalAction:
		return []cli.Flag{
			altsrc.NewStringFlag(
				cli.StringFlag{
					Name:   "username",
					Hidden: true,
				},
			),
			altsrc.NewStringFlag(
				cli.StringFlag{
					Name:   "password",
					Hidden: true,
				},
			),
			altsrc.NewStringFlag(
				cli.StringFlag{
					Name:   "host",
					Hidden: true,
				},
			),
			altsrc.NewStringFlag(
				cli.StringFlag{
					Name:   "vhost",
					Hidden: true,
				},
			),
			altsrc.NewIntFlag(
				cli.IntFlag{
					Name:   "port",
					Hidden: true,
				},
			),
			altsrc.NewIntFlag(
				cli.IntFlag{
					Name:   "apiport",
					Hidden: true,
				},
			),
			cli.StringFlag{
				Name:  "u",
				Usage: "username"},
			cli.StringFlag{
				Name:  "p",
				Usage: "password"},
			cli.StringFlag{
				Name:  "H",
				Usage: "host"},
			cli.StringFlag{
				Name:  "V",
				Value: "/",
				Usage: "vhost"},
			cli.IntFlag{
				Name:  "P",
				Value: 5672,
				Usage: "amqp port"},
			cli.IntFlag{
				Name:  "AP",
				Value: 15672,
				Usage: "api port"},
			cli.StringFlag{
				Name:  "config,c",
				Value: configDir,
				Usage: "configuration file"},
			cli.BoolFlag{
				Name:  "debug,d",
				Usage: "prints out debug log"},
		}
	case createQueueAction:
		return []cli.Flag{
			cli.BoolFlag{
				Name:  "durable,d",
				Usage: "queue survives after cluster restarts",
			},
			cli.BoolFlag{
				Name:  "autodelete,a",
				Usage: "queue will be deleted if no active consumers",
			},
			cli.BoolFlag{
				Name:  "HA",
				Usage: "HA enabled, creates policy name 'queuename_HA' under vhost -v (default: /)",
			},
			cli.StringFlag{
				Name:  "HAMODE",
				Value: "all",
				Usage: "HA mode, all|exactly|nodes",
			},
			cli.StringFlag{
				Name:  "HAPARAM",
				Usage: "HA parameters",
			},
			cli.StringFlag{
				Name:  "HASYNC",
				Value: "automatic",
				Usage: "HA SYNC mode, automatic|manual",
			},
		}
	case createExchangeAction:
		return []cli.Flag{
			cli.StringFlag{
				Name:  "type,t",
				Value: "direct",
				Usage: "kind, direct|fanout|topic|headers",
			},
			cli.BoolFlag{
				Name:  "durable,d",
				Usage: "survives after cluster restarts",
			},
			cli.BoolFlag{
				Name:  "autodelete,a",
				Usage: "exchange will be deleted if no active consumers",
			},
		}
	case createBindAction:
		return []cli.Flag{
			cli.StringFlag{
				Name:  "type,t",
				Value: "queue",
				Usage: "type, queue|exchange",
			},
		}
	case createUserAction:
		return []cli.Flag{
			cli.StringFlag{
				Name:  "tag,t",
				Value: "",
				Usage: "tag, none|management|policymaker|monitoring|administrator",
			},
		}
	case createVhostAction:
		return []cli.Flag{
			cli.BoolFlag{
				Name:  "trace,t",
				Usage: "enable vhost tracing",
			},
		}
	case listQueueAction:
		return []cli.Flag{
			cli.StringFlag{
				Name:  "o",
				Value: "plain",
				Usage: "output format, plain|json|rawjson|bash",
			},
		}
	case listExchangeAction:
		return []cli.Flag{
			cli.StringFlag{
				Name:  "o",
				Value: "plain",
				Usage: "output format, plain|json|rawjson|bash",
			},
		}
	case listNodeAction:
		return []cli.Flag{
			cli.StringFlag{
				Name:  "o",
				Value: "plain",
				Usage: "output format, plain|json|rawjson|bash",
			},
		}
	case listVhostAction:
		return []cli.Flag{
			cli.StringFlag{
				Name:  "o",
				Value: "plain",
				Usage: "output format, plain|json|rawjson|bash",
			},
		}
	case listBindAction:
		return []cli.Flag{
			cli.BoolFlag{
				Name:  "all,a",
				Usage: "list bindings from all vhosts",
			},
			cli.StringFlag{
				Name:  "o",
				Value: "plain",
				Usage: "output format, plain|json|rawjson|bash",
			},
		}
	case listPolicyAction:
		return []cli.Flag{
			cli.BoolFlag{
				Name:  "all,a",
				Usage: "list policies from all vhosts",
			},
			cli.StringFlag{
				Name:  "o",
				Value: "plain",
				Usage: "output format, plain|json|rawjson|bash",
			},
		}
	case listUserAction:
		return []cli.Flag{
			cli.StringFlag{
				Name:  "o",
				Value: "plain",
				Usage: "output format, plain|json|rawjson|bash",
			},
		}
	case deleteBindAction:
		return []cli.Flag{
			cli.StringFlag{
				Name:  "type,t",
				Value: "queue",
				Usage: "type, queue|exchange",
			},
		}
	case publishAction:
		return []cli.Flag{
			cli.IntFlag{
				Name:  "burst,b",
				Value: 1,
				Usage: "burst publishing",
			},
			cli.StringFlag{
				Name:  "mode",
				Value: "transient",
				Usage: "publish deliver mode, transient|persistent",
			},
			cli.BoolFlag{
				Name:  "mandatory,m",
				Usage: "if no queues bound to exchange, delivery fails",
			},
			cli.BoolFlag{
				Name:  "immediate,i",
				Usage: "if no consumers on the matched queue, delivery fails",
			},
		}
	case consumeAction:
		return []cli.Flag{
			cli.BoolFlag{
				Name:  "daemon,d",
				Usage: "daemon mode",
			},
			cli.StringFlag{
				Name:  "acktype,t",
				Value: "ack",
				Usage: "acknowledge type, ack|nack|reject",
			},
			cli.BoolFlag{
				Name:  "autoack,a",
				Usage: "acknowledge by default once receives message",
			},
			cli.BoolFlag{
				Name:  "nowait,nw",
				Usage: "begins without waiting cluster to confirm",
			},
			cli.StringFlag{
				Name:  "o",
				Value: "plain",
				Usage: "output format, plain|json",
			},
		}
	case updateUserAction:
		return []cli.Flag{
			cli.StringFlag{
				Name:  "tag,t",
				Value: "",
				Usage: "tag, none|management|policymaker|monitoring|administrator",
			},
		}
	case updateVhostAction:
		return []cli.Flag{
			cli.BoolFlag{
				Name:  "trace,t",
				Usage: "enable vhost tracing",
			},
		}
	default:
		return nil
	}
}

func commands() []cli.Command {
	cmds := []cli.Command{
		{
			Name:        "create",
			Usage:       "create resource",
			UsageText:   "rmqctl [global options] create resource [resource options] [arguments...]",
			Description: "rmqctl create queue/exchange/bind/bindex/user/vhost",
			Subcommands: []cli.Command{
				{
					Name:        "queue",
					Usage:       "rmqctl [global options] create queue [queue options] QUEUE_NAME",
					UsageText:   "create resource queue",
					Description: "create queue under vhost -v (default: /)",
					Category:    "create",
					Flags:       flags(createQueueAction),
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(createQueueAction),
					After:  doneAfter,
				},
				{
					Name:        "exchange",
					Usage:       "rmqctl [global options] create exchange [exchange options] EXCHANGE_NAME",
					UsageText:   "create resource exchange",
					Description: "create exchange",
					Category:    "create",
					Flags:       flags(createExchangeAction),
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(createExchangeAction),
					After:  doneAfter,
				},
				{
					Name:        "bind",
					Usage:       "rmqctl [global options] create bind [bind options] SOURCE(Exchange Name) DESTINATION ROUTING_KEY",
					UsageText:   "binds DESTINATION(queue or exchange, --type queue|exchange) to SOURCE(exchange)",
					Description: "create binding for Queue or Exchange to source exchange",
					Category:    "create",
					Flags:       flags(createBindAction),
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(createBindAction),
					After:  doneAfter,
				},
				{
					Name:        "user",
					Usage:       "rmqctl [global options] create user [user options] USERNAME PASSWORD",
					UsageText:   "create user with tag none(if no '-t')|management|policymaker|monitoring|administrator",
					Description: "tag information: https://www.rabbitmq.com/management.html#permissions",
					Category:    "create",
					Flags:       flags(createUserAction),
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(createUserAction),
					After:  doneAfter,
				},
				{
					Name:        "vhost",
					Usage:       "rmqctl [global options] create vhost [vhost options] VHOST_NAME",
					UsageText:   "create vhost with tracing using '-t'",
					Description: "tracing information: https://www.rabbitmq.com/firehose.html",
					Category:    "create",
					Flags:       flags(createVhostAction),
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(createVhostAction),
					After:  doneAfter,
				},
			},
		},
		{
			Name:        "list",
			Usage:       "rmqctl [global options] list resource [resource options] [arguments...]",
			UsageText:   "list resource",
			Description: "rmqctl list queue/exchange/bind/node/policy/user/vhost",
			Subcommands: []cli.Command{
				{
					Name:        "queue",
					Usage:       "rmqctl [global options] list queue [queue options] [QUEUE_NAME optional]",
					UsageText:   "list resource queue",
					Description: "list queue under vhost(default: /)",
					Category:    "list",
					Flags:       flags(listQueueAction),
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(listQueueAction),
				},
				{
					Name:        "exchange",
					Usage:       "rmqctl [global options] list exchange [exchange options] [EXCHANGE_NAME optional]",
					UsageText:   "list resource exchange",
					Description: "list exchange under vhost(default: /)",
					Category:    "list",
					Flags:       flags(listExchangeAction),
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(listExchangeAction),
				},
				{
					Name:        "bind",
					Usage:       "rmqctl [global options] list bind [bind options]",
					UsageText:   "list resource bind",
					Description: "list binding",
					Category:    "list",
					Flags:       flags(listBindAction),
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(listBindAction),
				},
				{
					Name:        "vhost",
					Usage:       "rmqctl [global options] list vhost [vhost options] [VHOST_NAME optional]",
					UsageText:   "list resource vhost",
					Description: "list vhost",
					Category:    "list",
					Flags:       flags(listVhostAction),
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(listVhostAction),
				},
				{
					Name:        "node",
					Usage:       "rmqctl [global options] list node [node options] [NODE_NAME optional]",
					UsageText:   "list nodes",
					Description: "list node",
					Category:    "list",
					Flags:       flags(listNodeAction),
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(listNodeAction),
				},
				{
					Name:        "policy",
					Usage:       "rmqctl [global options] list policy [policy options] [POLICY_NAME optional]",
					UsageText:   "list resource policy",
					Description: "list policy",
					Category:    "list",
					Flags:       flags(listPolicyAction),
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(listPolicyAction),
				},
				{
					Name:        "user",
					Usage:       "rmqctl [global options] list user [user options] [USERNAME optional]",
					UsageText:   "list resource user",
					Description: "list users/specified user",
					Category:    "list",
					Flags:       flags(listUserAction),
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(listUserAction),
				},
			},
		},
		{
			Name:        "delete",
			Usage:       "rmqctl [global options] delete resource [resource options] [arguments...]",
			UsageText:   "delete resource",
			Description: "rmqctl delete queue/exchange/bind/bindex/policy/user/vhost",
			Subcommands: []cli.Command{
				{
					Name:        "queue",
					Usage:       "rmqctl [global options] delete queue [queue options] QUEUE_NAME",
					UsageText:   "delete resource queue",
					Description: "delete queue under vhost(default: /)",
					Category:    "delete",
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(deleteQueueAction),
					After:  doneAfter,
				},
				{
					Name:        "exchange",
					Usage:       "rmqctl [global options] delete exchange [exchange options] EXCHANGE_NAME",
					UsageText:   "delete resource exchange",
					Description: "delete exchange under vhost(default: /)",
					Category:    "delete",
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(deleteExchangeAction),
					After:  doneAfter,
				},
				{
					Name:        "bind",
					Usage:       "rmqctl [global options] delete bind [bind options] SOURCE(Exchange Name) DESTINATION ROUTING_KEY",
					UsageText:   "delete bind DESTINATION(queue or exchange, --type queue|exchange) to SOURCE(exchange)",
					Description: "delete binding from Queue or Exchange to source exchange",
					Category:    "delete",
					Flags:       flags(deleteBindAction),
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(deleteBindAction),
					After:  doneAfter,
				},
				{
					Name:        "policy",
					Usage:       "rmqctl [global options] delete policy POLICY_NAME",
					UsageText:   "delete resource policy",
					Description: "delete policy under vhost -v (default: /)",
					Category:    "delete",
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(deletePolicyAction),
					After:  doneAfter,
				},
				{
					Name:        "user",
					Usage:       "rmqctl [global options] delete user USERNAME",
					UsageText:   "delete resource user",
					Description: "delete user",
					Category:    "delete",
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(deleteUserAction),
					After:  doneAfter,
				},
				{
					Name:        "vhost",
					Usage:       "rmqctl [global options] delete vhost VHOST_NAME",
					UsageText:   "delete resource vhost",
					Description: "delete vhost",
					Category:    "delete",
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(deleteVhostAction),
					After:  doneAfter,
				},
			},
		},
		{
			Name:        "publish",
			Usage:       "rmqctl [global options] publish [publish options] EXCHANGE_NAME KEY MESSAGE",
			UsageText:   "publish exchange key message",
			Description: "rmqctl publish EXCHANGE_NAME KEY MESSAGE",
			Category:    "publish",
			Flags:       flags(publishAction),
			Action:      actions(publishAction),
			After:       doneAfter,
		},
		{
			Name:        "consume",
			Usage:       "rmqctl [global options] consume [consume options] QUEUE_NAME",
			UsageText:   "consume queue",
			Description: "rmqctl consume QUEUE_NAME",
			Category:    "consume",
			Flags:       flags(consumeAction),
			Action:      actions(consumeAction),
		},
		{
			Name:        "update",
			Usage:       "update resource",
			UsageText:   "rmqctl [global options] update resource [resource options] [arguments...]",
			Description: "rmqctl update user/vhost",
			Subcommands: []cli.Command{
				{
					Name:        "user",
					Usage:       "rmqctl [global options] update user [user options] USERNAME PASSWORD",
					UsageText:   "update resource user",
					Description: "update existing user",
					Category:    "update",
					Flags:       flags(updateUserAction),
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(updateUserAction),
					After:  doneAfter,
				},
				{
					Name:        "vhost",
					Usage:       "rmqctl [global options] update vhost [vhost options] VHOST_NAME",
					UsageText:   "update resource vhost",
					Description: "update existing vhost, enable tracing by providing '-t'",
					Category:    "update",
					Flags:       flags(updateVhostAction),
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(updateVhostAction),
					After:  doneAfter,
				},
			},
		},
	}

	return cmds
}

// Cmd starts the application.
func Cmd() error {
	defer logCleanUp()

	regSignalHandler(
		func() {
			os.Exit(128 + int(syscall.SIGQUIT))
		},
		syscall.SIGQUIT,
	)

	regSignalHandler(
		func() {
			os.Exit(128 + int(syscall.SIGTERM))
		},
		syscall.SIGTERM,
	)

	cliapp := cli.NewApp()
	cliapp.Authors = []cli.Author{
		{
			Name:  "shuo-huan chang/verbalsaint",
			Email: "vsdmars@gmail.com",
		},
	}
	cliapp.Copyright = "LICENSE information on https://github.com/vsdmars/rmqctl"

	cliapp.Name = "rmqctl"
	cliapp.Usage = "tool for controlling rabbitmq cluster."
	cliapp.UsageText = "rmqctl [global options] command subcommand [subcommand options] [arguments...]"
	cliapp.Description = "rmqctl is a swiss-knife for rabbitmq cluster."
	cliapp.Version = rmqctlVersion.string()

	flags := flags(generalAction)
	commands := commands()
	cliapp.Before = mainBefore(flags)
	cliapp.Flags = flags
	cliapp.Commands = commands

	if err := cliapp.Run(os.Args); err != nil {
		logger.Debug("encountered error",
			zap.String("application", "rmqctl"),
			zap.String("version", cliapp.Version),
			zap.String("error", err.Error()))
		return err
	}

	return nil
}
