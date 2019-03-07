rmqctl_
-------

.. All external links are here
.. _rmqctl: https://github.com/vsdmars/rmqctl
.. _rmqctl.conf: ./rmqctl.conf
.. _rabbitmq_tls.config: ./example/rabbitmq_tls.config
.. |travis| image:: https://api.travis-ci.org/vsdmars/rmqctl.svg?branch=v1
  :target: https://travis-ci.org/vsdmars/rmqctl
.. |go report| image:: https://goreportcard.com/badge/github.com/vsdmars/rmqctl
  :target: https://goreportcard.com/report/github.com/vsdmars/rmqctl
.. |go doc| image:: https://godoc.org/github.com/vsdmars/rmqctl?status.svg
  :target: https://godoc.org/github.com/vsdmars/rmqctl
.. |license| image:: https://img.shields.io/github/license/mashape/apistatus.svg?style=flat
  :target: ./LICENSE
.. |release| image:: https://img.shields.io/badge/release-v1.0.12-blue.svg
  :target: https://github.com/vsdmars/rmqctl/tree/v1.0.11
.. _binary release v1.0.0: https://github.com/vsdmars/rmqctl/releases/tag/v1.0.0
.. _binary release v1.0.3: https://github.com/vsdmars/rmqctl/releases/tag/v1.0.3
.. _binary release v1.0.7: https://github.com/vsdmars/rmqctl/releases/tag/v1.0.7
.. _binary release v1.0.8: https://github.com/vsdmars/rmqctl/releases/tag/v1.0.8
.. _binary release v1.0.9: https://github.com/vsdmars/rmqctl/releases/tag/v1.0.9
.. _binary release v1.0.10: https://github.com/vsdmars/rmqctl/releases/tag/v1.0.10
.. _binary release v1.0.11: https://github.com/vsdmars/rmqctl/releases/tag/v1.0.11
.. _binary release v1.0.12: https://github.com/vsdmars/rmqctl/releases/tag/v1.0.12

.. ;; And now we continue with the actual content

|travis| |go report| |go doc| |license| |release|

----

rmqctl is *the* swiss-army knife tool for rabbitmq with kubectl like commands.

----


Binary Release:
---------------

`binary release v1.0.12`_
 - fix issues for rabbit-hole.DeleteBinding uses BindingInfo.PropertiesKey as routing key

 instead of BindingInfo.RoutingKey

 - reference:
      https://cdn.rawgit.com/rabbitmq/rabbitmq-management/v3.7.12/priv/www/api/index.html
      /api/bindings/vhost/e/exchange/q/queue/props
      https://github.com/michaelklishin/rabbit-hole/blob/master/bindings.go#L193

`binary release v1.0.11`_
 - Logging bug fix

`binary release v1.0.10`_
 - Purge queue / purge queue with prompt [y/n]
 - Consume queue with numbers, e.g only consumes 10 messages

`binary release v1.0.9`_
 - honors -a, -d in create queue/exchange

`binary release v1.0.8`_
 - Now supports TLS connection for AMQP and HTTPS
 - New 'tls' entry in rmqctl.conf_
 - New flag '-T' indicates using TLS connection.
 - Bug fix.

`binary release v1.0.7`_
 - Now supports burst message publish mode.

   Alone with daemon mode, rmqctl is used as a stress test tool for rabbitmq.

   e.g.
    $ rmqctl publish exchange_name routing_key "MESSAGE" -b 1000000

   Publish with other payload
    $ rmqctl publish exchange_name routing_key "$(cat payload.json)" -b 1000000

 - Now supports publish mode: Transient, Persistent
 - Change default config file name to *rmqctl.conf*
 - Change load config file name flag to '-c'
 - Formalize debug log message.

`binary release v1.0.3`_
 - Publish/Consume use amqp protocol for performance.
   Other actions using rabbitmq REST API calls.
 - Now supports bash/rawjson output format.

`binary release v1.0.0`_
 - init. release


rmqctl.conf_
-------------

rmqctl loads rmqctl.conf (yaml) under working directory if there is one.
Command arguments have higher precedence if provided.

.. code:: yaml

   username: guest
   password: guest
   port: 5672
   apiport: 15672
   host: localhost
   tls: true
   vhost: "/"


::

 Loads rmqctl.conf from other location
 $ rmqctl -c path/to/rmqctl.conf COMMANDS


=========
Supports
=========

AMQP Protocol
-------------
rmqctl_ uses amqp protocol library for publish/consume message for speed.

rmqctl_ supports burst publish/daemon consume, act as a perfect tool for stress test

and debugging the application.


TLS support
-----------
Place client certificate and private key pair with read only permission (0400)

under $HOME/.ssh/ name as follows:


::

   ~/.ssh/rmq_cert.pem
   ~/.ssh/rmq_key.pem


If rabbitmq server using self-signed certificate,

remember to register self-signed CA into client's host system.

Setting up rabbitmq server TLS support for both

AMQP and API Service config file can refere to example:

rabbitmq_tls.config_



Create
------
- queue, queue in HA mode(with single command)
- exchange
- queue/exchange binding
- user
- vhost
- --help for more features


List
----
- queue
- exchange
- queue/exchange binding
- user
- vhost
- node
- policy
- --help for more features


Delete
------
- queue
- exchange
- queue/exchange binding
- user
- vhost
- policy
- --help for more features


Update
------
- vhost
- user
- --help for more features


Publish
-------
- Publish with routing key
- Burst publishing
- Supports transient|persistent modes
- --help for more features


Consume
-------
- Consume supports ack|nack|reject|auto-ack acknowledge modes.
- Run as daemon, consume on-demand.
- Consume number of messages with flag -c NUMBER
- --help for more features


Purge
-----
- Purge queue with prompt
- --help for more features


=====
Usage
=====

Create queue
------------

::

   // TEST_QUEUE_1 created as durable
   $ rmqctl create queue TEST_QUEUE_1 -d
   done

   // TEST_QUEUE_2 created as durable and autodelete
   $ rmqctl -d create queue TEST_QUEUE_2 -d -a
   done


Create queue in HA mode
-----------------------

rmqctl is able to create queue in HA mode.

Three modes supported: all(default),exactly,nodes

Following command creates TEST_QUEUE_3 queue in HA mode,

which by default it has queue slaves in all other rabbitmq nodes (default: 'all' mode)

rmqctl automatically creates queue's HA policy with name: QueueName_HA

::

   $ rmqctl create queue TEST_QUEUE_3 --HA
   done


List all queues
---------------

::

   $ rmqctl list queue
   |Name         |Vhost |Durable |AutoDelete |MasterNode |Status |Consumers |Policy          |Messages
   |TEST_QUEUE_1 |/     |true    |false      |rabbit@r1  |       |0         |                |0
   |TEST_QUEUE_2 |/     |true    |true       |rabbit@r1  |       |0         |                |0
   |TEST_QUEUE_3 |/     |true    |true       |rabbit@r1  |       |0         |TEST_QUEUE_3_HA |0


List Policy
-----------

::

   $ rmqctl list policy
    Name            |Vhost |Pattern      |Priority |ApplyTo |Definition
   |TEST_QUEUE_3_HA |/     |TEST_QUEUE_3 |0        |queues  |map[ha-mode:all ha-sync-mode:automatic]


List particular queue in json
-----------------------------

::

   $ rmqctl list queue TEST_QUEUE_1 -o json

.. code:: json

   [
     {
       "name": "TEST_QUEUE_1",
       "vhost": "/",
       "durable": true,
       "auto_delete": false,
       "arguments": {},
       "node": "rabbit@r1",
       "status": "",
       "memory": 10576,
       ...
       }
    ]


Create exchange
---------------

::

  $ rmqctl create exchange TEST_EXCHANGE_1 -d -t fanout
  done


List all exchanges
------------------

::

  $ rmqctl list exchange
   |Name               |Vhost |Type    |Durable |AutoDelete
   |                   |/     |direct  |true    |false
   |TEST_EXCHANGE_1    |/     |fanout  |true    |false
   |amq.direct         |/     |direct  |true    |false
   |amq.fanout         |/     |fanout  |true    |false
   |amq.headers        |/     |headers |true    |false
   |amq.match          |/     |headers |true    |false
   |amq.rabbitmq.trace |/     |topic   |true    |false
   |amq.topic          |/     |topic   |true    |false


List particular exchange in json
--------------------------------

::

   $ rmqctl list exchange TEST_EXCHANGE_1 -o json

.. code:: json

   {
     "name": "TEST_EXCHANGE_1",
     "vhost": "/",
     "type": "fanout",
     "durable": true,
     "auto_delete": false,
     "internal": false,
     "arguments": {},
     "incoming": [],
     "outgoing": []
   }


Create queue/exchange binding
-----------------------------

rmqctl is able to create exchange bindings as well.

::

  $ rmqctl create bind TEST_EXCHANGE_1 TEST_QUEUE_1 ROUTING_KEY
  done
  $ rmqctl create bind TEST_EXCHANGE_1 TEST_QUEUE_2 ROUTING_KEY
  done

  Creates exchange binding
  $ rmqctl create bind TEST_EXCHANGE_1 TEST_EXCHANGE_2 ROUTING_KEY -t exchange
  done


List queue/exchange binding
---------------------------

::

  $ rmqctl list bind
  |Source          |Destination     |Vhost |Key          |DestinationType
  |                |TEST_QUEUE_1    |/     |TEST_QUEUE_1 |queue
  |                |TEST_QUEUE_2    |/     |TEST_QUEUE_2 |queue
  |TEST_EXCHANGE_1 |TEST_QUEUE_1    |/     |RUN          |queue
  |TEST_EXCHANGE_1 |TEST_EXCHANGE_2 |/     |RUN          |exchange


Publish message
---------------

Publish to a fanout exchange, observing queues bounded to the

exchange *TEST_EXCHANGE_1* received the message.

::

   $ rmqctl publish TEST_EXCHANGE_1 RUN "This is a test message"
   done

   $ rmqctl list queue
   |Name         |Vhost |Durable |AutoDelete |MasterNode |Status |Consumers |Policy          |Messages
   |TEST_QUEUE_1 |/     |true    |false      |rabbit@r1  |       |0         |                |1
   |TEST_QUEUE_2 |/     |true    |true       |rabbit@r1  |       |0         |                |1
   |TEST_QUEUE_3 |/     |true    |true       |rabbit@r1  |       |0         |TEST_QUEUE_3_HA |0


Publish message in burst mode
-----------------------------

Publish to a fanout exchange in burst mode,

observing queues bounded to the exchange *TEST_EXCHANGE_1* received the message.

::

   $ rmqctl publish TEST_EXCHANGE_1 RUN "This is a test message" -b 424242
   done

   $ rmqctl list queue
   |Name         |Vhost |Durable |AutoDelete |MasterNode |Status |Consumers |Policy          |Messages
   |TEST_QUEUE_1 |/     |true    |false      |rabbit@r1  |       |0         |                |424243
   |TEST_QUEUE_2 |/     |true    |true       |rabbit@r1  |       |0         |                |424243
   |TEST_QUEUE_3 |/     |true    |true       |rabbit@r1  |       |0         |TEST_QUEUE_3_HA |0


Consume message
---------------
Consume 3 messages.

::

   $ rmqctl consume TEST_QUEUE_1 -c 3
   |Message
   This is a test message
   This is a test message
   This is a test message



Consume message in daemon mode
------------------------------

::

   $ rmqctl consume TEST_QUEUE_2 -d
   |Message
   This is a test message
   This is a test message
   ...


Purge queue
-----------
Purge queue without prompt.

::

   $ rmqctl purge TEST_QUEUE_1 -f
   done



Other features including list/update user/vhost/node information, vhost tracing, etc.
-------------------------------------------------------------------------------------
--help for more details.

::

   $ rmqctl --help


Contact
-------
Bug, feature requests, welcome to shoot me an email at:

**vsdmars<at>gmail.com**
