# Music organisation and streaming system

Tchaik is an open source music organisation and streaming system.  The backend is written in [Go](http://golang.org), the frontend is built using [React](https://facebook.github.io/react/), [Flux](https://facebook.github.io/flux/) and [PostCSS](https://github.com/postcss/postcss).

![Tchaik UI](https://s3-ap-southeast-2.amazonaws.com/dhowden-pictures/tchaik-may.jpg "Tchaik UI")
![Tchaik UI (Classical Music)](https://s3-ap-southeast-2.amazonaws.com/dhowden-pictures/tchaik-july.png "Tchaik UI (Classical Music)")

# Features

* Automatic prefix grouping and enumeration detection (ideal for classical music: properly groups big works together).
* Multiplatform web-based UI and REST-like API for controlling player.
* Multiple storage and caching options: Amazon S3, Google Cloud Storage, local and remote file stores.
* Import music library from iTunes or build from directory-tree of audio tracks.

# Requirements

* [Go 1.4+](http://golang.org/dl/) (recent changes have only been tested on 1.4.2).
* [NodeJS](https://nodejs.org/), [NPM](https://www.npmjs.com/) and [Gulp](http://gulpjs.com/) installed globally (for building the UI).
* Recent version of Chrome (Firefox may also work, though hasn't been fully tested).

# Building

If you haven't setup Go before, you need to first set a `GOPATH` (see [https://golang.org/doc/code.html#GOPATH](https://golang.org/doc/code.html#GOPATH)).

To fetch and build the code for Tchaik:

    $ go get tchaik.com/cmd/...

This will fetch the code and build the command line tools into `$GOPATH/bin` (assumed to be in your `PATH` already).

Building the UI:

    $ cd $GOPATH/src/tchaik.com/cmd/tchaik/ui
    $ npm install
    $ gulp

Alternatively, if you want the JS and CSS to be recompiled and have the browser refreshed as you change the source files:

    $ WS_URL="ws://localhost:8080/socket" gulp serve

Then browse to `http://localhost:3000/` to use tchaik.

# Starting the UI

To start Tchaik you first need to move into the `cmd/tchaik` directory:

    $ cd $GOPATH/src/tchaik.com/cmd/tchaik

## Importing an iTunes Library

The easiest way to begin is to build a Tchaik library on-the-fly and start the UI in one command:

    $ tchaik -itlXML ~/path/to/iTunesLibrary.xml

You can also convert the iTunes Library into a Tchaik library using the `tchimport` tool:

    $ tchimport -itlXML ~/path/to/iTunesLibrary.xml -out lib.tch
    $ tchaik -lib lib.tch

NB: A Tchaik library will generally be smaller than its corresponding iTunes Library.  Tchaik libraries are stored as gzipped-JSON (rather than Apple plist) and contain a subset of the metadata used by iTunes.

## Importing Audio Files

Alternatively you can build a Tchaik library on-the-fly from a directory-tree of audio files. Only files with supported metadata (see [github.com/dhowden/tag](https://github.com/dhowden/tag)) will be included in the index:

    $ tchaik -path /all/my/music

To avoid rescanning your entire collection every time you restart, you can build a Tchaik library using the `tchimport` tool:

    $ tchimport -path /all/my/music -out lib.tch
    $ tchaik -lib lib.tch

# More Advanced Options

A full list of command line options is available from the `--help` flag:

    $ tchaik --help
    Usage of ./tchaik:
      -add-path-prefix string
        	add prefix to every path
      -artwork-cache string
        	path to local artwork cache (content addressable)
      -auth-password string
        	password to use for HTTP authentication
      -auth-user string
        	username to use for HTTP authentication (set to enable)
      -checklist string
        	path to checklist file (default "checklist.json")
      -debug
        	print debugging information
      -favourites string
        	path to favourites file (default "favourites.json")
      -itlXML string
        	path to iTunes Library XML file
      -lib string
        	path to Tchaik library file
      -listen string
        	bind address to http listen (default "localhost:8080")
      -local-store string
        	local media store, full local path /path/to/root (default "/")
      -media-cache string
        	path to local media cache
      -path string
        	path to directory containing music files (to build index from)
      -play-history string
        	path to play history file (default "history.json")
      -remote-store string
        	remote media store, tchstore server address <hostname>:<port>, s3://<region>:<bucket>/path/to/root for S3, or gs://<project-id>:<bucket>/path/to/root for Google Cloud Storage
      -static-dir string
        	Path to the static asset directory (default "ui/static")
      -tls-cert string
        	path to a certificate file, must also specify -tls-key
      -tls-key string
        	path to a certificate key file, must also specify -tls-cert
      -trace-listen string
        	bind address for trace HTTP server
      -trim-path-prefix string
        	remove prefix from every path

### -local-store

Set `-local-store` to the local path that contains your media files.  You can use `trim-path-prefix` and `add-path-prefix` to rewrite paths used in the Tchaik library so that file locations can still be correctly resolved.

### -remote-store

Set `-remote-store` to the URI of a running [tchstore](http://godoc.org/tchaik.com/cmd/tchstore) server  (`hostname:port`).  Instead, S3 paths can be used: `s3://<region>:<bucket>/path/to/root` (set the environment variables `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` to pass credentials to the S3 client), or Google Cloud Storage paths: `gs://<project-id>:<bucket>/path/to/root` (set the GOOGLE_APPLICATION_CREDENTIALS to point at your JSON credentials file).

### -media-cache

Set `-media-cache` to cache all files loaded from `-remote-store` (or `-local-store` if set).

### -artwork-cache

Set `-artwork-cache` to create/use a content addressable filesystem for track artwork.  An index file will be created in the path on first use.  The folder should initially be empty to ensure that no other files interfere with the system.

### -trace-listen

Set `-trace-listen` to a suitable bind address (i.e. `localhost:4040`) to start an HTTP server which defines the `/debug/requests` endpoint used to inspect server requests.  Currently we only support tracing for media (track/artwork/icon) requests.  See [https://godoc.org/golang.org/x/net/trace](https://godoc.org/golang.org/x/net/trace) for more details. 

# Windows Support

The default value for parameter `-local-store` is `/` which does not work on Windows.  When all library music is organised under a common path you can set `-local-store` and `-trim-path-prefix` to get around this (for instance `-local-store C:\Path\To\Music -trim-path-prefix C:\Path\To\Music`).

# Get Involved!

Development is on-going and the codebase is changing very quickly.  If you're interested in contributing then it's safest to jump into our gitter room and chat to people before getting started!

[![Join the chat at https://gitter.im/tchaik/tchaik](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/tchaik/tchaik?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
