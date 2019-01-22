rmqctl_
=====================

.. All external links are here
.. _rmqctl: https://github.com/vsdmars/rmqctl
.. _rabbitmqadmin: https://www.rabbitmq.com/management-cli.html
.. |travis| image:: https://api.travis-ci.org/vsdmars/rmqctl.svg?branch=v1
  :target: https://travis-ci.org/vsdmars/rmqctl
.. |go report| image:: https://goreportcard.com/badge/github.com/vsdmars/rmqctl
  :target: https://goreportcard.com/report/github.com/vsdmars/rmqctl
.. |go doc| image:: https://godoc.org/github.com/vsdmars/rmqctl?status.svg
  :target: https://godoc.org/github.com/vsdmars/rmqctl
.. |license| image:: https://img.shields.io/github/license/mashape/apistatus.svg?style=flat
  :target: ./LICENSE
.. |release| image:: https://img.shields.io/badge/release-v1.0.0-blue.svg
  :target: https://github.com/vsdmars/rmqctl/tree/v1.0.0


.. ;; And now we continue with the actual content

|travis| |go report| |go doc| |license| |release|

----

rmqctl is a Golang version of `rabbitmqadmin`_ with simular
commands like kubectl.

----

## rmqctl_config.yaml

rmqctl_config.yaml contains connection information to
rabbitmq cluster.

.. code:: yaml

   username: guest
   password: guest
   port: 5672
   apiport: 15672
   host: "127.0.0.1"
   vhost: "/"


rmqctl by default loads rmqctl_config.yaml under the working directory.
rmqctl can loads rmqctl_config.yaml from other location by using
--load path_to_rmqctl_config.yaml


## Usage

```
$ rmqctl --help
```



```
NAME:
   rmqctl - tool for controlling rabbitmq cluster.

USAGE:
   rmqctl [global options] command subcommand [subcommand options] [arguments...]

VERSION:
   v1.0.0

DESCRIPTION:
   rmqctl is a swiss-knife for rabbitmq cluster.

AUTHOR:
   verbalsaint <vsdmars@gmail.com>

COMMANDS:
     create   create resource
     list     rmqctl [global options] list resource [resource options] [arguments...]
     delete   rmqctl [global options] delete resource [resource options] [arguments...]
     update   update resource
     help, h  Shows a list of commands or help for one command
   consume:
     consume  rmqctl [global options] consume [consume options] QUEUE_NAME
   publish:
     publish  rmqctl [global options] publish [publish options] EXCHANGE_NAME KEY MESSAGE

GLOBAL OPTIONS:
   --username value  cluster username
   --password value  cluster password
   --host value      cluster host
   --vhost value     cluster vhost (default: "/")
   --port value      cluster port (default: 5672)
   --apiport value   cluster api port (default: 15672)
   --load value      config file location (default: "~/rmqctl_config.yaml")
   --debug, -d       run in debug mode
   --help, -h        show help
   --version, -v     print the version

COPYRIGHT:
   LICENSE information on https://github.com/vsdmars/rmqctl

```


## Consume message in daemon mode

```
$ rmqctl consume --help
```

```
NAME:
   rmqctl consume - rmqctl [global options] consume [consume options] QUEUE_NAME

USAGE:
   consume queue

CATEGORY:
   consume

DESCRIPTION:
   rmqctl consume QUEUE_NAME

OPTIONS:
   --daemon, -d               daemon mode
   --acktype value, -t value  acknowledge type, ack|nack|reject (default: "ack")
   --autoack, -a              acknowledge by default once receives message
   --nowait, --nw             begins without waiting cluster to confirm
   -o value                   output format, plain|json (default: "plain")
```

Example:
```
$ rmqctl consume -d QUEUE_NAME
```

----

Kubernetes is an open source system for managing [containerized applications]
across multiple hosts; providing basic mechanisms for deployment, maintenance,
and scaling of applications.

Kubernetes builds upon a decade and a half of experience at Google running
production workloads at scale using a system called [Borg],
combined with best-of-breed ideas and practices from the community.

Kubernetes is hosted by the Cloud Native Computing Foundation ([CNCF]).
If you are a company that wants to help shape the evolution of
technologies that are container-packaged, dynamically-scheduled
and microservices-oriented, consider joining the CNCF.
For details about who's involved and how Kubernetes plays a role,
read the CNCF [announcement].

----

## To start using Kubernetes

See our documentation on [kubernetes.io].

Try our [interactive tutorial].

Take a free course on [Scalable Microservices with Kubernetes].

## To start developing Kubernetes

The [community repository] hosts all information about
building Kubernetes from source, how to contribute code
and documentation, who to contact about what, etc.

If you want to build Kubernetes right away there are two options:

##### You have a working [Go environment].

```
$ go get -d k8s.io/kubernetes
$ cd $GOPATH/src/k8s.io/kubernetes
$ make
```

##### You have a working [Docker environment].

```
$ git clone https://github.com/kubernetes/kubernetes
$ cd kubernetes
$ make quick-release
```

For the full story, head over to the [developer's documentation].

## Support

If you need support, start with the [troubleshooting guide],
and work your way through the process that we've outlined.
