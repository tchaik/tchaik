# Music organisation and streaming system

Tchaik is an open source music organisation and streaming system.  The backend is written in [Go](http://golang.org), the frontend is built using [React](https://facebook.github.io/react/) and [Flux](https://facebook.github.io/flux/).

![Tchaik UI](https://s3-ap-southeast-2.amazonaws.com/dhowden-pictures/screenshot.jpg "Tchaik UI")

# Features

* Automatic prefix grouping and enumeration detection (ideal for classical music: properly groups big works together).
* Multiplatform web-based UI and REST-like API for building remote controls.
* Multiple storage options: Amazon S3, local and remote file stores.
* Imports iTunes Music Library files.

# Getting up and running

### Requirements

* Go 1.4 (recent changes have only been tested on 1.4.2).
* NodeJS, NPM and Gulp installed globally (for building the UI).
* Recent version of Chrome/Firefox/Safari.

### Building

    $ go get github.com/dhowden/tchaik/...

Which will fetch the code and build the command line tools, putting the executables into `$GOPATH/bin` (which is assumed to be in your `PATH` already).

Building the UI:

    $ cd $GOPATH/src/github.com/dhowden/tchaik/cmd/tchaik/ui
    $ npm install
    $ NODE_ENV=development gulp

Alternative, if you want the JS and CSS to be recompiled as you change the source files:

    $ NODE_ENV=development gulp watch

Now you can start Tchaik.  For the moment this means first moving to the cmd/tchaik directory:

    $ cd $GOPATH/src/github.com/dhowden/tchaik/cmd/tchaik
    $ tchaik -itlXML ~/path/to/iTunesLibrary.xml
    Parsing ~/path/to/iTunesLibrary.xml...done.
    Building Tchaik Library...done.
    Building root collection...done.
    Building search index...done.
    Web server is running on http://localhost:8080
    Quit the server with CTRL-C.
