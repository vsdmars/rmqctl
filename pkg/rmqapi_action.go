package pkg

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	rh "github.com/michaelklishin/rabbit-hole"
	"go.uber.org/zap"
	cli "gopkg.in/urfave/cli.v1"
)

const httpTLSTimeout = 5
const httpReqTimeout = 3

func connectAPI(amqpconn *amqpConnectionType) (*rh.Client, error) {
	logurl := func(u url.URL) {
		logger.Debug(
			"api URL",
			zap.String("service", "api"),
			zap.String("URL", u.String()),
		)
	}

	var apiURL url.URL
	var transport *http.Transport

	if amqpconn.TLS {
		certfile, keyfile, err := getCertPath()
		if err != nil {
			return nil, err
		}

		apiURL = url.URL{
			Scheme: "https",
			Host: fmt.Sprintf(
				"%s:%d",
				amqpconn.Host,
				amqpconn.APIPort,
			),
		}

		logurl(apiURL)

		cfg := &tls.Config{}

		cert, err := tls.LoadX509KeyPair(certfile, keyfile)
		if err != nil {
			logger.Debug(
				"load x509 key pair failed",
				zap.String("service", "api"),
			)

			return nil, err
		}

		cfg.Certificates = append(cfg.Certificates, cert)

		transport = &http.Transport{
			TLSClientConfig:     cfg,
			TLSHandshakeTimeout: time.Duration(httpTLSTimeout) * time.Second,
		}

		logger.Debug(
			"timeout",
			zap.String("service", "api"),
			zap.String("tls handshake timeout", (time.Duration(httpTLSTimeout)*time.Second).String()),
		)

	} else {
		apiURL = url.URL{
			Scheme: "http",
			Host: fmt.Sprintf(
				"%s:%d",
				amqpconn.Host,
				amqpconn.APIPort,
			),
		}

		logurl(apiURL)
	}

	conn, err := rh.NewClient(
		apiURL.String(),
		amqpconn.Username,
		amqpconn.Password,
	)
	if err != nil {
		logger.Debug(
			"api create connection failed.",
			zap.String("service", "api"),
		)

		return nil, cli.NewExitError(err.Error(), 1)
	}

	// http client request timeout in 3 seconds
	conn.SetTimeout(time.Duration(httpReqTimeout) * time.Second)

	logger.Debug(
		"timeout",
		zap.String("service", "api"),
		zap.String("request timeout", (time.Duration(httpReqTimeout)*time.Second).String()),
	)

	if amqpconn.TLS {
		conn.SetTransport(transport)
	}

	return conn, nil
}

func createQueue(conn *rh.Client, data *createQueueType) error {
	setting := rh.QueueSettings{
		Durable:    data.Durable,
		AutoDelete: data.Autodelete,
	}

	res, err := conn.DeclareQueue(
		data.Vhost,
		data.QueueName,
		setting,
	)

	if err := handleHTTPResponse(res, "create", "queue"); err != nil {
		return err
	}

	if err != nil {
		logger.Debug(
			"create queue failed",
			zap.String("service", "api"),
			zap.String("queue", data.QueueName),
			zap.String("vhost", data.Vhost),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	if data.Ha {
		return createQueueHA(conn, data)
	}

	return nil
}

func createQueueHA(conn *rh.Client, data *createQueueType) error {
	policyDefFunc := func() rh.PolicyDefinition {
		if data.HaMode == "all" {
			return rh.PolicyDefinition{
				"ha-mode":      data.HaMode,
				"ha-sync-mode": data.HaSyncMode,
			}
		}

		return rh.PolicyDefinition{
			"ha-mode":      data.HaMode,
			"ha-params":    data.HaParam,
			"ha-sync-mode": data.HaSyncMode,
		}
	}

	policy := rh.Policy{
		Vhost:      data.Vhost,
		Pattern:    data.QueueName,
		ApplyTo:    "queues", // queues, exchanges, all
		Definition: policyDefFunc(),
	}

	res, err := conn.PutPolicy(
		data.Vhost,
		fmt.Sprintf("%s_%s", data.QueueName, "HA"),
		policy,
	)

	if err := handleHTTPResponse(res, "create", "policy"); err != nil {
		return err
	}

	if err != nil {
		logger.Debug(
			"set queue HA policy failed",
			zap.String("service", "api"),
			zap.String("queue", data.QueueName),
			zap.String("vhost", data.Vhost),
			zap.String("HA mode", data.HaMode),
			zap.String("HA param", data.HaParam),
			zap.String("HA sync mode", data.HaSyncMode),
			zap.String("http response", fmt.Sprintf("%v", res)),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func createExchange(conn *rh.Client, data *createExchangeType) error {
	setting := rh.ExchangeSettings{
		Type:       data.Kind,
		Durable:    data.Durable,
		AutoDelete: data.Autodelete,
	}

	res, err := conn.DeclareExchange(
		data.Vhost,
		data.ExchangeName,
		setting,
	)

	if err := handleHTTPResponse(res, "create", "exchange"); err != nil {
		return err
	}

	if err != nil {
		logger.Debug(
			"create exchange failed",
			zap.String("service", "api"),
			zap.String("exchange", data.ExchangeName),
			zap.String("vhost", data.Vhost),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func createBind(conn *rh.Client, data *createBindType) error {
	setting := rh.BindingInfo{
		Source:          data.SourceExchangeName,
		Vhost:           data.Vhost,
		Destination:     data.DestinationName,
		DestinationType: data.Type,
		RoutingKey:      data.Key,
	}

	res, err := conn.DeclareBinding(
		data.Vhost,
		setting,
	)

	if err := handleHTTPResponse(res, "create", "bind"); err != nil {
		return err
	}

	if err != nil {
		logger.Debug(
			"create bind failed",
			zap.String("service", "api"),
			zap.String("source", data.SourceExchangeName),
			zap.String("destination", data.DestinationName),
			zap.String("routing key", data.Key),
			zap.String("type", data.Type),
			zap.String("vhost", data.Vhost),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func createUser(conn *rh.Client, data *createUserType) error {
	// check user existence.
	userInfo, err := conn.GetUser(data.RmqUsername)
	if err == nil {
		logger.Debug(
			"user already exists",
			zap.String("service", "api"),
			zap.String("user", userInfo.Name),
			zap.String("tag", userInfo.Tags),
		)

		return cli.NewExitError("user already exists", 1)
	}

	setting := rh.UserSettings{
		Name:     data.RmqUsername,
		Tags:     data.Tag,
		Password: data.RmqPassword}

	res, err := conn.PutUser(data.RmqUsername, setting)

	if err := handleHTTPResponse(res, "create", "user"); err != nil {
		return err
	}

	if err != nil {
		logger.Debug(
			"create user failed",
			zap.String("service", "api"),
			zap.String("user", data.RmqUsername),
			zap.String("tag", data.Tag),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func createVhost(conn *rh.Client, data *createVhostType) error {
	// check vhost existence.
	vhostInfo, err := conn.GetVhost(data.VhostName)
	if err == nil {
		logger.Debug(
			"vhost already exists",
			zap.String("service", "api"),
			zap.String("vhost", vhostInfo.Name),
			zap.Bool("tracing", vhostInfo.Tracing),
		)

		return cli.NewExitError("vhost already exists", 1)
	}

	setting := rh.VhostSettings{
		Tracing: data.Tracing}

	res, err := conn.PutVhost(data.VhostName, setting)

	if err := handleHTTPResponse(res, "create", "vhost"); err != nil {
		return err
	}

	if err != nil {
		logger.Debug(
			"create vhost failed",
			zap.String("service", "api"),
			zap.String("vhost", data.VhostName),
			zap.Bool("tracing", data.Tracing),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func listQueue(conn *rh.Client, data *listQueueType) error {
	if data.QueueName != "" {
		queueInfo, err := conn.GetQueue(data.Vhost, data.QueueName)
		if err != nil {
			logger.Debug(
				"no such queue",
				zap.String("service", "api"),
				zap.String("vhost", data.Vhost),
				zap.String("queue", data.QueueName),
			)

			return cli.NewExitError(err.Error(), 1)
		}

		return detailedQueueInfoF(data.Formatter, queueInfo)
	}

	queueInfos, err := conn.ListQueuesIn(data.Vhost)
	if err != nil {
		logger.Debug(
			"list queues failed",
			zap.String("service", "api"),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	return detailedQueueInfoSliceF(data.Formatter, queueInfos)
}

func listExchange(conn *rh.Client, data *listExchangeType) error {
	if data.ExchangeName != "" {
		exgInfo, err := conn.GetExchange(data.Vhost, data.ExchangeName)
		if err != nil {
			logger.Debug(
				"no such exchange",
				zap.String("service", "api"),
				zap.String("vhost", data.Vhost),
				zap.String("exchange", data.ExchangeName),
			)

			return cli.NewExitError(err.Error(), 1)
		}

		return detailedExchangeInfoF(data.Formatter, exgInfo)
	}

	exgInfos, err := conn.ListExchangesIn(data.Vhost)
	if err != nil {
		logger.Debug(
			"list exchange failed",
			zap.String("service", "api"),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	return detailedExchangeInfoSliceF(data.Formatter, exgInfos)
}

func listBind(conn *rh.Client, data *listBindType) error {
	var binfos []rh.BindingInfo
	var err error

	if data.All {
		binfos, err = conn.ListBindings()
		if err != nil {
			logger.Debug(
				"list bindings from all vhosts failed",
				zap.String("service", "api"),
			)

			return cli.NewExitError(err.Error(), 1)
		}

	} else {
		binfos, err = conn.ListBindingsIn(data.Vhost)
		if err != nil {
			logger.Debug(
				"list bindings from vhost failed",
				zap.String("service", "api"),
				zap.String("vhost", data.Vhost),
			)

			return cli.NewExitError(err.Error(), 1)
		}
	}

	return bindingInfoF(data.Formatter, binfos)
}

func listVhost(conn *rh.Client, data *listVhostType) error {
	if data.VhostName != "" {
		vhostInfo, err := conn.GetVhost(data.VhostName)
		if err != nil {
			logger.Debug(
				"no such vhost",
				zap.String("service", "api"),
				zap.String("vhost", data.VhostName),
			)

			return cli.NewExitError(err.Error(), 1)
		}

		return vhostInfoF(data.Formatter, vhostInfo)
	}

	vhostInfos, err := conn.ListVhosts()
	if err != nil {
		logger.Debug(
			"list vhosts failed",
			zap.String("service", "api"),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	return vhostInfoSliceF(data.Formatter, vhostInfos)
}

func listNode(conn *rh.Client, data *listNodeType) error {
	if data.NodeName != "" {
		nodeInfo, err := conn.GetNode(data.NodeName)
		if err != nil {
			logger.Debug(
				"no such node",
				zap.String("service", "api"),
				zap.String("node", data.NodeName),
			)

			return cli.NewExitError(err.Error(), 1)
		}

		return nodeInfoF(data.Formatter, nodeInfo)
	}

	nodeInfos, err := conn.ListNodes()
	if err != nil {
		logger.Debug(
			"list nodes failed",
			zap.String("service", "api"),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	return nodeInfoSliceF(data.Formatter, nodeInfos)
}

func listPolicy(conn *rh.Client, data *listPolicyType) error {
	if data.PolicyName != "" {
		policyInfo, err := conn.GetPolicy(data.Vhost, data.PolicyName)
		if err != nil {
			logger.Debug(
				"no such policy",
				zap.String("service", "api"),
				zap.String("policy", data.PolicyName),
			)

			return cli.NewExitError(err.Error(), 1)
		}

		return policyInfoF(data.Formatter, policyInfo)
	}

	if data.All {
		policyInfos, err := conn.ListPolicies()
		if err != nil {
			logger.Debug(
				"list all vhosts' policy failed",
				zap.String("service", "api"),
			)

			return cli.NewExitError(err.Error(), 1)
		}

		return policyInfoSliceF(data.Formatter, policyInfos)
	}

	policyInfos, err := conn.ListPoliciesIn(data.Vhost)
	if err != nil {
		logger.Debug(
			"list vhost's policy failed",
			zap.String("service", "api"),
			zap.String("vhost", data.Vhost),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	return policyInfoSliceF(data.Formatter, policyInfos)
}

func listUser(conn *rh.Client, data *listUserType) error {
	if data.RmqUsername != "" {
		userInfo, err := conn.GetUser(data.RmqUsername)
		if err != nil {
			logger.Debug(
				"no such user",
				zap.String("service", "api"),
				zap.String("user", data.RmqUsername),
			)

			return cli.NewExitError(err.Error(), 1)
		}

		return userInfoF(data.Formatter, userInfo)
	}

	userInfos, err := conn.ListUsers()
	if err != nil {
		logger.Debug(
			"list users failed",
			zap.String("service", "api"),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	return userInfoSliceF(data.Formatter, userInfos)
}

func deleteQueue(conn *rh.Client, data *deleteQueueType) error {
	res, err := conn.DeleteQueue(
		data.Vhost,
		data.QueueName,
	)

	if err := handleHTTPResponse(res, "delete", "queue"); err != nil {
		return err
	}

	if err != nil {
		logger.Debug(
			"delete queue failed",
			zap.String("service", "api"),
			zap.String("queue", data.QueueName),
			zap.String("vhost", data.Vhost),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func deleteExchange(conn *rh.Client, data *deleteExchangeType) error {
	res, err := conn.DeleteExchange(
		data.Vhost,
		data.ExchangeName,
	)

	if err := handleHTTPResponse(res, "delete", "exchange"); err != nil {
		return err
	}

	if err != nil {
		logger.Debug(
			"delete exchange failed",
			zap.String("service", "api"),
			zap.String("exchange", data.ExchangeName),
			zap.String("vhost", data.Vhost),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func deleteBind(conn *rh.Client, data *deleteBindType) error {
	setting := rh.BindingInfo{
		Source:          data.SourceExchangeName,
		Vhost:           data.Vhost,
		Destination:     data.DestinationName,
		DestinationType: data.Type,
		RoutingKey:      data.Key,
	}

	res, err := conn.DeleteBinding(
		data.Vhost,
		setting,
	)

	if err := handleHTTPResponse(res, "delete", "bind"); err != nil {
		return err
	}

	if err != nil {
		logger.Debug(
			"delete bind failed",
			zap.String("service", "api"),
			zap.String("source", data.SourceExchangeName),
			zap.String("destination", data.DestinationName),
			zap.String("routing key", data.Key),
			zap.String("type", data.Type),
			zap.String("vhost", data.Vhost),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func deletePolicy(conn *rh.Client, data *deletePolicyType) error {
	res, err := conn.DeletePolicy(data.Vhost, data.PolicyName)
	if err != nil {
		logger.Debug(
			"delete policy failed",
			zap.String("service", "api"),
			zap.String("policy", data.PolicyName),
			zap.String("http response", fmt.Sprintf("%v", res)),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func deleteUser(conn *rh.Client, data *deleteUserType) error {
	res, err := conn.DeleteUser(data.RmqUsername)

	if err := handleHTTPResponse(res, "delete", "user"); err != nil {
		return err
	}

	if err != nil {
		logger.Debug(
			"delete user failed",
			zap.String("service", "api"),
			zap.String("user", data.RmqUsername),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func deleteVhost(conn *rh.Client, data *deleteVhostType) error {
	res, err := conn.DeleteVhost(data.VhostName)

	if err := handleHTTPResponse(res, "delete", "vhost"); err != nil {
		return err
	}

	if err != nil {
		logger.Debug(
			"delete vhost failed",
			zap.String("service", "api"),
			zap.String("vhost", data.VhostName),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func updateUser(conn *rh.Client, data *updateUserType) error {
	// check user existence.
	_, err := conn.GetUser(data.RmqUsername)
	if err != nil {
		logger.Debug(
			"user doesn't exist",
			zap.String("service", "api"),
			zap.String("user", data.RmqUsername),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	setting := rh.UserSettings{
		Name:     data.RmqUsername,
		Tags:     data.Tag,
		Password: data.RmqPassword}

	res, err := conn.PutUser(data.RmqUsername, setting)

	if err := handleHTTPResponse(res, "update", "user"); err != nil {
		return err
	}

	if err != nil {
		logger.Debug(
			"update user failed",
			zap.String("service", "api"),
			zap.String("user", data.RmqUsername),
			zap.String("tag", data.Tag),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func updateVhost(conn *rh.Client, data *updateVhostType) error {
	// check vhost existence.
	_, err := conn.GetVhost(data.VhostName)
	if err != nil {
		logger.Debug(
			"vhost doesn't exist",
			zap.String("service", "api"),
			zap.String("vhost", data.VhostName),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	setting := rh.VhostSettings{
		Tracing: data.Tracing}

	res, err := conn.PutVhost(data.VhostName, setting)

	if err := handleHTTPResponse(res, "delete", "vhost"); err != nil {
		return err
	}

	if err != nil {
		logger.Debug(
			"create vhost failed",
			zap.String("service", "api"),
			zap.String("vhost", data.VhostName),
			zap.Bool("tracing", data.Tracing),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}
