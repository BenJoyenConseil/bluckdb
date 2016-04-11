# Bluckdb

A Key/Value store implementation using Golang

The ``server.go`` file is a simple http server that answers on the 8080 port.


There are two endpoints :

    http://hostname:8080/get/?key=<some_key>
    http://hostname:8080/put/?key=<some_key>&value=<some_value>


## the goal

The goal of this project is to explore and reinvent the wheel of well known, state of the art, algorithms and data structures.
For experimental and learning purpose only, not production ready.