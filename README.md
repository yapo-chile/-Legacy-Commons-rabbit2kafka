# rabbit2kafka

rabbit2kafka  is a microservice on golang for Yapo to move message from a RabbitMQ queue binded to an exchange to Kafka with a given topic


## How to run the service

* Create the dir: `~/go/src/github.mpi-internal.com/Yapo`

* Set the go path: `export GOPATH=~/go` or add the line on your file `.bash_rc`

* Clone this repo:

  ```shell
  cd ~/go/src/github.mpi-internal.com/Yapo
  git clone git@github.mpi-internal.com:Yapo/rabbit2kafka.git
  ```

* On the top dir execute the make instruction to clean and start:

  ```shell
  cd rabbit2kafka
  make start
  ```

* If you change the code:

  `make start`

* To create the Docker image

  `make docker-build`

* To run the Docker image

  `make docker-compose-up-dev`

* To stop the Docker image

  `make docker-compose-down-dev`

* To attach to the running Docker image

  `make docker-attach`
  
  ## Available vars

|Variable               |Default       |Description                |
|-----------------------|--------------|---------------------------|
|KAFKA_HOST             |kafka         |Broker host                |
|KAFKA_TOPIC            |events-queue  |Topic to send messages     |
|KAFKA_PORT             |9092          |Broker port                |
|RABBITMQ_HOST          |rabbit        |Rabbit host                |
|RABBITMQ_PORT          |5672          |Rabbit port                |
|RABBITMQ_QUEUE         |backend_event |Queue with message to read |
|RABBITMQ_VHOST         |backend_event |Virtual host to connect    |
|RABBITMQ_EXCHANGE      |backend_event |Distibute the messages     |
|RABBITMQ_EXCHANGE_TYPE |topic         |Type of distribution       |
|RABBITMQ_USER          |yapo          |User to connect            |
|RABBITMQ_PASSWORD      |yapo2014      |Password of the user       |
|RABBITMQ_CONSUMER_TAG  |              |Id for consumer client     |
|LOGGER_SYSLOG_ENABLED  |false         |System log on/off          |
|LOGGER_SYSLOG_IDENTITY |              |Identify this ms           |
|LOGGER_STDLOG_ENABLED  |true          |Standar log on/off         |
|LOGGER_LOG_LEVEL       |0             |Level of log (0 = debug)   |

