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
	configYaml = "rmqctl_config.yaml"
)

const (
	generalAction = iota
	createQueueAction
	createExchangeAction
	createBindAction
	createBindExAction
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
	deleteBindExAction
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
		createBindExAction:   createBindExJob,
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
		deleteBindExAction:   deleteBindExJob,
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
			altsrc.NewStringFlag(cli.StringFlag{Name: "username",
				Usage: "cluster username"}),
			altsrc.NewStringFlag(cli.StringFlag{Name: "password",
				Usage: "cluster password"}),
			altsrc.NewStringFlag(cli.StringFlag{Name: "host",
				Usage: "cluster host"}),
			altsrc.NewStringFlag(cli.StringFlag{Name: "vhost", Value: "/",
				Usage: "cluster vhost"}),
			altsrc.NewIntFlag(cli.IntFlag{Name: "port", Value: 5672,
				Usage: "cluster port"}),
			altsrc.NewIntFlag(cli.IntFlag{Name: "apiport", Value: 15672,
				Usage: "cluster api port"}),
			cli.StringFlag{Name: "load", Value: configDir,
				Usage: "config file location"},
			cli.BoolFlag{Name: "debug, d", Usage: "run in debug mode"},
		}
	case createQueueAction:
		return []cli.Flag{
			cli.BoolFlag{Name: "durable,du",
				Usage: "queue survives after cluster restarts"},
			cli.BoolFlag{Name: "autodelete, ad",
				Usage: "queue will be deleted if no active consumers"},
			cli.BoolFlag{Name: "exclusive, exc",
				Usage: "queue is accessible only from this connection"},
			cli.BoolFlag{Name: "nowait, nw",
				Usage: "queue is assumed to be declared before"},
			cli.BoolFlag{Name: "ha",
				Usage: "HA enabled, creates policy name 'queuename_HA' under vhost -v (default: /)"},
			cli.StringFlag{Name: "hamode, hm", Value: "all",
				Usage: "HA mode, all|exactly|nodes"},
			cli.StringFlag{Name: "haparam, hp",
				Usage: "HA parameters"},
			cli.StringFlag{Name: "syncmode, sm", Value: "automatic",
				Usage: "HA SYNC mode, automatic|manual"},
		}
	case createExchangeAction:
		return []cli.Flag{
			cli.StringFlag{Name: "type, t", Value: "direct",
				Usage: "kind, direct|fanout|topic|headers"},
			cli.BoolFlag{Name: "durable,du",
				Usage: "survives after cluster restarts"},
			cli.BoolFlag{Name: "autodelete, ad",
				Usage: "exchange will be deleted if no active consumers"},
			cli.BoolFlag{Name: "internal, itn",
				Usage: "exchange does not accept publishings"},
			cli.BoolFlag{Name: "nowait, nw",
				Usage: "does not wait for confirmation from the server"},
		}
	case createBindAction:
		return []cli.Flag{
			cli.BoolFlag{Name: "nowait, nw",
				Usage: "does not wait for confirmation from the server"},
		}
	case createBindExAction:
		return []cli.Flag{
			cli.BoolFlag{Name: "nowait, nw",
				Usage: "does not wait for confirmation from the server"},
		}
	case createUserAction:
		return []cli.Flag{
			cli.StringFlag{Name: "tag, t", Value: "",
				Usage: "tag, none|management|policymaker|monitoring|administrator"},
		}
	case createVhostAction:
		return []cli.Flag{
			cli.BoolFlag{Name: "trace, t", Usage: "enable vhost tracing"},
		}
	case listQueueAction:
		return []cli.Flag{
			cli.StringFlag{Name: "o", Value: "plain",
				Usage: "output format, plain|json"},
		}
	case listExchangeAction:
		return []cli.Flag{
			cli.StringFlag{Name: "o", Value: "plain",
				Usage: "output format, plain|json"},
		}
	case listNodeAction:
		return []cli.Flag{

			cli.StringFlag{Name: "o", Value: "plain",
				Usage: "output format, plain|json"},
		}
	case listVhostAction:
		return []cli.Flag{
			cli.StringFlag{Name: "o", Value: "plain",
				Usage: "output format, plain|json"},
		}
	case listBindAction:
		return []cli.Flag{
			cli.BoolFlag{Name: "all, a",
				Usage: "list bindings from all vhosts"},
			cli.StringFlag{Name: "o", Value: "plain",
				Usage: "output format, plain|json"},
		}
	case listPolicyAction:
		return []cli.Flag{
			cli.BoolFlag{Name: "all, a",
				Usage: "list policies from all vhosts"},
			cli.StringFlag{Name: "o", Value: "plain",
				Usage: "output format, plain|json"},
		}
	case listUserAction:
		return []cli.Flag{
			cli.StringFlag{Name: "o", Value: "plain",
				Usage: "output format, plain|json"},
		}
	case deleteQueueAction:
		return []cli.Flag{
			cli.BoolFlag{Name: "force, f",
				Usage: "delete queue even it has consumers or messages"},
			cli.BoolFlag{Name: "nowait, nw",
				Usage: "without waiting for cluster's response"},
		}
	case deleteExchangeAction:
		return []cli.Flag{
			cli.BoolFlag{Name: "force, f",
				Usage: "delete exchange even it has queue binding to it"},
			cli.BoolFlag{Name: "nowait, nw",
				Usage: "without waiting for cluster's response"},
		}
	case deleteBindExAction:
		return []cli.Flag{
			cli.BoolFlag{Name: "nowait, nw",
				Usage: "without waiting for cluster's response"},
		}
	case publishAction:
		return []cli.Flag{
			cli.BoolFlag{Name: "mandatory, m",
				Usage: "if no queues bound to exchange, delivery fails"},
			cli.BoolFlag{Name: "immediate, i",
				Usage: "if no consumers on the matched queue, delivery fails"},
		}
	case consumeAction:
		return []cli.Flag{
			cli.BoolFlag{Name: "daemon, d",
				Usage: "daemon mode"},
			cli.StringFlag{Name: "acktype, t", Value: "ack",
				Usage: "acknowledge type, ack|nack|reject"},
			cli.BoolFlag{Name: "autoack, a",
				Usage: "acknowledge by default once receives message"},
			cli.BoolFlag{Name: "nowait, nw",
				Usage: "begins without waiting cluster to confirm"},
			cli.StringFlag{Name: "o", Value: "plain",
				Usage: "output format, plain|json"},
		}
	case updateUserAction:
		return []cli.Flag{
			cli.StringFlag{Name: "tag, t", Value: "",
				Usage: "tag, none|management|policymaker|monitoring|administrator"},
		}
	case updateVhostAction:
		return []cli.Flag{
			cli.BoolFlag{Name: "trace, t", Usage: "enable vhost tracing"},
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
					Name:        "bindex",
					Usage:       "rmqctl [global options] create bindex [bind-exchange options] FROM_EXCHANGE TO_EXCHANGE KEY",
					UsageText:   "create resource bind-exchange",
					Description: "create bindex",
					Category:    "create",
					Flags:       flags(createBindExAction),
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(createBindExAction),
					After:  doneAfter,
				},
				{
					Name:        "bind",
					Usage:       "rmqctl [global options] create bind [bind options] QUEUE_NAME EXCHANGE_NAME KEY",
					UsageText:   "create resource bind queue exchange key",
					Description: "create binding",
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
					Flags:       flags(deleteQueueAction),
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
					Flags:       flags(deleteExchangeAction),
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(deleteExchangeAction),
					After:  doneAfter,
				},
				{
					Name:        "bind",
					Usage:       "rmqctl [global options] delete bind [bind options] QUEUE_NAME EXCHANGE_NAME KEY",
					UsageText:   "delete resource bind",
					Description: "delete binding",
					Category:    "delete",
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(deleteBindAction),
					After:  doneAfter,
				},
				{
					Name:        "bindex",
					Usage:       "rmqctl [global options] delete bindex [bindex options] FROM_EXCHANGE TO_EXCHANGE KEY",
					UsageText:   "delete resource exchange bind",
					Description: "delete exchange binding",
					Category:    "delete",
					Flags:       flags(deleteBindExAction),
					// UseShortOptionHandling: true, //support in cli.v2
					Action: actions(deleteBindExAction),
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

	regSignalHandler(func() {
		os.Exit(128 + int(syscall.SIGQUIT))
	}, syscall.SIGQUIT)

	regSignalHandler(func() {
		os.Exit(128 + int(syscall.SIGTERM))
	}, syscall.SIGTERM)

	cliapp := cli.NewApp()
	cliapp.Authors = []cli.Author{{Name: "verbalsaint", Email: "vsdmars@gmail.com"}}
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
		logger.Debug("rmqctl encountered error",
			zap.String("error", err.Error()))
		return err
	}

	return nil
}
