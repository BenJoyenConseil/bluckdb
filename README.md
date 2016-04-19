# Bluckdb

[![Build Status](https://travis-ci.org/BenJoyenConseil/bluckdb.svg?branch=master)](https://travis-ci.org/BenJoyenConseil/bluckdb) [![Stories in Ready](https://badge.waffle.io/BenJoyenConseil/bluckdb.png?label=ready&title=Ready)](https://waffle.io/BenJoyenConseil/bluckdb)

A Key/Value store implementation using Golang

The ``server.go`` file is a simple http server that answers on the 8080 port.


There are two endpoints :

    http://hostname:2233/get/?key=<some_key>
    http://hostname:2233/put/?key=<some_key>&value=<some_value>


## the goal

The goal of this project is to explore and reinvent the wheel of well known, state of the art, algorithms and data structures.
For experimental and learning purpose only, not production ready.


# How to start

## Get the package
* go get github.com/BenJoyenConseil/bluckdb
* If you run a go program for the first time, do not forget to setup your GOPATH : export GOPATH=$HOME/Dev/go

## Run the server
* go run server.go
* Go will silently exit if a process is already using port 2233
