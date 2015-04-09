# Music organisation and streaming system

Tchaik is an open source music organisation and streaming system.

![Tchaik UI](https://s3-ap-southeast-2.amazonaws.com/dhowden-pictures/screenshot.jpg "Tchaik UI")

# Getting up and running

### Requirements:

* Go 1.4 (recent changes have only been tested on 1.4.2)
* NodeJS + NPM, SASS/Compass (for building the UI)

### Fetching the code

    $ go get github.com/dhowden/tchaik

Building the command-line tools:

    $ export TCH=$GOPATH/src/github.com/dhowden/tchaik
    $ cd $TCH/cmd/tchaik
    $ go build

Building the UI:

    $ cd ui  # from inside $TCH/cmd/tchaik
    $ npm install
    $ NODE_ENV=development grunt

Building the CSS:

    $ cd static
    $ compass compile

Now you can start Tchaik:

    $ cd $TCH/cmd/tchaik
    $ ./tchaik -itlXML ~/path/to/iTunesLibrary.xml
    Parsing ~/path/to/iTunesLibrary.xml...done.
    Building Tchaik Library...done.
    Building root collection...done.
    Building search index...done.
    Web server is running on http://localhost:8080
    Quit the server with CTRL-C.

