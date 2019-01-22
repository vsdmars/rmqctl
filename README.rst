rmqctl_
-------

.. All external links are here
.. _rmqctl: https://github.com/vsdmars/rmqctl
.. _rabbitmqadmin: https://www.rabbitmq.com/management-cli.html
.. _rmqctl_config.yaml: ./rmqctl_config.yaml
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
.. _binary release v1.0.0: https://github.com/vsdmars/rmqctl/releases/tag/v1.0.0

.. ;; And now we continue with the actual content

|travis| |go report| |go doc| |license| |release|

----

rmqctl is a Golang version of `rabbitmqadmin`_ with simular
commands like kubectl.

----


Binary Release:
---------------

`binary release v1.0.0`_


rmqctl_config.yaml_
-------------------

rmqctl_config.yaml contains connection information to
rabbitmq cluster.

.. code:: yaml

   username: guest
   password: guest
   port: 5672
   apiport: 15672
   host: "127.0.0.1"
   vhost: "/"


rmqctl by default loads *rmqctl_config.yaml* under the working directory.

rmqctl loads *rmqctl_config.yaml* from other location by using --load :

::

 rmqctl --load path_to_rmqctl_config.yaml COMMANDS


=========
Supports
=========

Create
------
- queue
- exchange
- queue binding
- exchange binding
- user
- vhost


List
----
- queue
- exchange
- queue binding
- exchange binding
- user
- vhost
- node
- policy


Delete
------
- queue
- exchange
- queue binding
- exchange binding
- user
- vhost
- policy


Update
------
- vhost
- user


Pulish
------
- Publish message to exchange with routing key


Consume
-------
- Consume message with specified acknowledge mode
- Run as daemon, consuming message to STDOUT


=====
Usage
=====

Create queue
------------

::

   $ rmqctl create queue TEST_QUEUE_1
   done
   $ rmqctl create queue TEST_QUEUE_2
   done


List all queues
---------------

::

   $ rmqctl list queue
   |Name         |Vhost |Durable |AutoDelete |MasterNode |Status |Consumers |Policy |Messages
   |TEST_QUEUE_1 |/     |false   |false      |rabbit@r1  |       |0         |       |0
   |TEST_QUEUE_2 |/     |false   |false      |rabbit@r1  |       |0         |       |0



List particular queue in json
-----------------------------

::

   $ rmqctl list queue TEST_QUEUE_1 -o json

.. code:: json

   [
     {
       "name": "TEST_QUEUE_1",
       "vhost": "/",
       "durable": false,
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

  $ rmqctl create exchange TEST_EXCHANGE_1 --durable -t fanout
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


Create queue binding
--------------------

::

  $ rmqctl create bind TEST_QUEUE_1 TEST_EXCHANGE_1 RUN
  done
  $ rmqctl create bind TEST_QUEUE_2 TEST_EXCHANGE_1 RUN
  done


List queue binding
------------------

::

  $ rmqctl list bind
  |Source          |Destination  |Vhost |Key          |DestinationType
  |                |TEST_QUEUE_1 |/     |TEST_QUEUE_1 |queue
  |                |TEST_QUEUE_2 |/     |TEST_QUEUE_2 |queue
  |TEST_EXCHANGE_1 |TEST_QUEUE_1 |/     |RUN          |queue


Publish message to exchange
---------------------------

Publish message to a fanout exchange, we'll see queues bounded to the

exchange *TEST_EXCHANGE_1* received the message.

::

   $ rmqctl publish TEST_EXCHANGE_1 RUN "This is a test message"
   done

   $ rmqctl list queue
   |Name         |Vhost |Durable |AutoDelete |MasterNode |Status |Consumers |Policy |Messages
   |TEST_QUEUE_1 |/     |false   |false      |rabbit@r1  |       |0         |       |1
   |TEST_QUEUE_2 |/     |false   |false      |rabbit@r1  |       |0         |       |1


Consume queue's messages
------------------------

::

   $ rmqctl consume TEST_QUEUE_1
   |Message
   |This is a test message



Consume queue's messages in daemon mode
---------------------------------------

::

   $ rmqctl consume TEST_QUEUE_2 -d
   |Message
   |This is a test message


Create user/vhost/exchang bind, update user info/vhost tracing etc.
-------------------------------------------------------------------
Use --help for specific details.

::

   $ rmqctl --help


Contact
-------
Bug, feature request, welcome to shoot me a email at:

**vsdmars<at>gmail.com**
