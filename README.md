# MrMeeseeksBot

## Project Setup
1. Install `docker` and `docker-composer`

[https://docs.docker.com/compose/install/](https://docs.docker.com/compose/install/)

2. Clone this repo and build containers
```
$ git clone https://github.com/melzareix/MrMeeseeksBot.git
$ cd MrMeeseeksBot
$ docker-compose build
```

3. To add golang dependencies edit `Docker/dev/golang/dockerfile` and rebuild with `docker-compose build`
```
RUN go get -u PACKAGE_GITHUB_URL
```

4. To run the container - add `-d` flag to run as daemon
```
$ docker-compose up
```
