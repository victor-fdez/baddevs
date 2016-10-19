
# BADdevs developer guide

## Requirements

Setup both NodeJS 6 & GOlang

* NodeJS6 : is used to develop the client application on Angular 2, the whole NodeJS6 application is contained in the client/ folder
* GOLang : is used to develop the http server to serve generated Node application html, js, and static files + to fullfil API requests etc

### Debian

#### NodeJS 6

The following commands will setup NodeJS6 on your debian based distribution. For other distribution please check [guide](https://nodejs.org/en/download/package-manager/)
 
```bash
#Install NodeJS
curl -sL https://deb.nodesource.com/setup_6.x | sudo -E bash -
sudo apt-get install -y nodejs
#Install build tools for some NodeJS packages
sudo apt-get install -y build-essential
```

##### Extra tools for Webpack and Typescript

The following commands will install all the required packages and tools for the project. If you would like to get more information about the project please visit it's [page](https://github.com/AngularClass/angular2-webpack-starter/tree/material2)

```bash
npm install --save
export PATH=./node_modules/.bin/:$PATH
```

##### Extra Linux Configs

Increase inotify watchers, else during development webpack won't be able to watch changes to most of your files, and so it won't automatically recompile. Check the [guide](https://github.com/guard/listen/wiki/Increasing-the-amount-of-inotify-watchers)

```bash
echo fs.inotify.max_user_watches=524288 | sudo tee -a /etc/sysctl.conf && sudo sysctl -p
```

#### GOlang

```bash
# Install the GOLang binaries
VERSION=1.7.1
OS=linux
ARCH=amd64
tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
#Install GO
#Setup some ENV variables for compiling GO programs
export GOROOT=$HOME/go
export GOBIN=$GOROOT/bin
export GOPATH=$GOROOT
export PATH=$PATH:/usr/local/go/bin:$GOROOT/bin
```
##### Glide

Glide is required in order to install the correct dependencies for this go project. This [website](https://github.com/Masterminds/glide) contains all of it's information.

```bash
curl https://glide.sh/get | sh
```

##### CompileDaemon

Install the compile daemon so that whenever you make a change to the go programs, it will automatically detect the change and compile the GO programs.

```bash
#Install CompileDaemon locally
go get github.com/githubnemo/CompileDaemon
```


## Run in Development

### Run in Development Mode

#### Run the Server

```bash
# get dependencies
glide install
CompileDaemon -color -command='./baddevs --port 8081 --host 0.0.0.0'
```

#### Run Client Generator using Web Pack

```bash
# Move to client/ folder
cd client/
# Run npm generator that regenerates based on changes
npm run watch:dev
```

### Run in Production Mode

#### Run the Server

```bash
# Build
glide install
go build
# Run
./baddevs --port 80 --host 0.0.0.0
# Helps
./baddevs
```

#### Run Client Generator using Web Pack

```bash
# Move to client/ folder
cd client/
# Run npm generator that regenerates based on changes
npm run build:prod 
```
