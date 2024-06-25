# Alby Hub <3 phoenixd

Run your Alby Hub with phoenixd as a backend.

## Requirements

+ Linux distribution
+ Runs pretty much on any VPS or server

## Docker

To run Alby Hub with phoenixd use [docker-compose](https://docs.docker.com/compose/) using the docker-compose.yml file.

    $ cd folder/with/docker-compose.yml
    $ docker-compose up

### Backup

Make sure to backup the `albyhub-phoenixd` which is used as volume for albyhub and phoenixd data files.

## Non Docker

### Installation (non-Docker)

    $ ./install.sh

The install script will prompt you for a installation folder and will install phoenixd and Alby Hub there.

Optionally it also creates a systemd services.

### Running the services

Either use systemd:

    $ sudo systemctl [start|stop] phoenixd.service
    $ sudo systemctl [start|stop] albyhub.service

Or us the start scripts:

    $ [your install path]/phoenixd/start.sh
    $ [your install path]/albyhub/start.sh


### Backup

Make sure to backup your data directories:

+ `[your install path]/phoenixd/data`
+ `[your install path]/albyhub/data`
