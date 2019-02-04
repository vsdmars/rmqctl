package pkg

import (
	"os"
	"os/user"
	"path"

	"go.uber.org/zap"
)

func getCertPath() (string, string, error) {
	usr, err := user.Current()
	if err != nil {
		logger.Debug(
			"failed to locate user's home dir",
			zap.String("service", "amqp"),
		)

		return "", "", err
	}

	certfile := path.Join(usr.HomeDir, ".ssh", "rmq_cert.pem")
	keyfile := path.Join(usr.HomeDir, ".ssh", "rmq_key.pem")

	logger.Debug(
		"cert path",
		zap.String("service", "amqp"),
		zap.String("certfile", certfile),
		zap.String("keyfile", keyfile),
	)

	if _, err := os.Stat(certfile); err != nil {
		logger.Debug(
			"rmq_cert.pem missing",
			zap.String("service", "amqp"),
		)

		return "", "", err
	}

	if _, err := os.Stat(keyfile); err != nil {
		logger.Debug(
			"rmq_key.pem missing",
			zap.String("service", "amqp"),
		)

		return "", "", err
	}

	return certfile, keyfile, nil
}
