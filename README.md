# rabbit2kafka

rabbit2kafka  is a microservice on golang for Yapo to move message from a RabbitMQ queue binded to an exchange to Kafka with a given topic

<!-- Badger start badges -->
[![Status of the build](https://badger.spt-engprod-pro.mpi-internal.com/badge/travis/Yapo/rabbit2kafka)](https://travis.mpi-internal.com/Yapo/rabbit2kafka)
<!-- Badger end badges -->

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

* How to check format

  `make validate`

* To create the Docker image

  `make docker-build`

* To run the Docker image

  `make docker-compose-up-dev`

* To stop the Docker image

  `make docker-compose-down-dev`

* To attach to the running Docker image

  `make docker-attach`

