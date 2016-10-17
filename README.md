# Bluckdb

[![Build Status](https://travis-ci.org/BenJoyenConseil/bluckdb.svg?branch=master)](https://travis-ci.org/BenJoyenConseil/bluckdb) [![Stories in Ready](https://badge.waffle.io/BenJoyenConseil/bluckdb.png?label=ready&title=Ready)](https://waffle.io/BenJoyenConseil/bluckdb) [![Go Report Card](https://goreportcard.com/badge/github.com/BenJoyenConseil/bluckdb)](https://goreportcard.com/report/github.com/BenJoyenConseil/bluckdb) [![GoDoc](https://godoc.org/github.com/BenJoyenConseil/bluckdb?status.svg)](https://godoc.org/github.com/BenJoyenConseil/bluckdb)

It is a Key/Value store that implements bucketing based on [extendible hashing](https://en.wikipedia.org/wiki/Extendible_hashing)

The ``server.go`` file is a simple http server that answers on the 8080 port.


There are 4 endpoints :

    http://hostname:2233/get?key=<some_key>
    http://hostname:2233/put?key=<some_key>&value=<some_value>
    http://hostname:2233/meta
    http://hostname:2233/debug?page_id=<id_of_the_page_to_display>


## the goal

The goal of this project is to explore and to reinvent the wheel of well known, state of the art, algorithms and data structures.
For experimental and learning purpose only, not production ready.


## design

A Directory is a table of buckets called "Page". 

A Page is a byte array of 4096 bytes length, append only. 
Trailer : It stores actual usage of the Page at 4094 bytes (unint16), and local depth at 9092 bytes (unint16)

A Record is a byte array with a key, a value and the headers :
 
    type Record interface {
        Key() []byte
        Val() []byte
        KeyLen() uint16
        ValLen() uint16
    }
    type ByteRecord []byte
         
Actual public methods :

* put : append the record at the offset given by `Page.use()` value
* get : read in reverse way, starting from the end and iterating until the key is found, or the beginning

This design allows updating values for a given key without doing lookup before inserting (put is O(1) if the Page is not full). When the Page is full, the `Directory.split()` method skips the old values of the same key and re-insert just the latest

# How to start

## Get the package

    go get github.com/BenJoyenConseil/bluckdb

If you run a go program for the first time, do not forget to setup your GOPATH : export GOPATH=$HOME/Dev/go

## Run the server

    go run server.go


Note : Go will silently exit if a process is already using **port 2233**

## Benchmarks
    
    BenchmarkBluckDBPut                    1000000	            2786 ns/op   -> 2,7 µs
    BenchmarkBluckDBGet                    1000000	            1406 ns/op   -> 1,4 µs
    BenchmarkPutNaiveDiskKVStore-4          200000              6250 ns/op   -> 6,2 µs
    BenchmarkGetNaiveDiskKVStore-4              30          44017416 ns/op   ->  44 ms
    BenchmarkPutHashMap-4                  1000000              1385 ns/op   -> 1,3 µs
    BenchmarkGetHashMap-4                  2000000               711 ns/op   -> 0,7 µs
