
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

The following global tools to be able to generate the client project. If you would like to get more information about the project please visit it's [page](https://github.com/AngularClass/angular2-webpack-starter/tree/material2)

```bash
npm install --global webpack
npm install --global webpack-dev-server
npm install --global karma-cli
npm install --global protractor
npm install --global typescript
npm install --global rimraf 
```

##### Extra Linux Configs

Increase inotify watchers, else during development webpack won't be able to watch changes to most of your files, and so it won't automatically recompile. Check the [guide](https://github.com/guard/listen/wiki/Increasing-the-amount-of-inotify-watchers)

```bash
echo fs.inotify.max_user_watches=524288 | sudo tee -a /etc/sysctl.conf && sudo sysctl -p
```

#### GOlang

```
```
