# Bluckdb

[![Build Status](https://travis-ci.org/BenJoyenConseil/bluckdb.svg?branch=master)](https://travis-ci.org/BenJoyenConseil/bluckdb) [![Stories in Ready](https://badge.waffle.io/BenJoyenConseil/bluckdb.png?label=ready&title=Ready)](https://waffle.io/BenJoyenConseil/bluckdb) [![Go Report Card](https://goreportcard.com/badge/github.com/BenJoyenConseil/bluckdb)](https://goreportcard.com/report/github.com/BenJoyenConseil/bluckdb) [![GoDoc](https://godoc.org/github.com/BenJoyenConseil/bluckdb?status.svg)](https://godoc.org/github.com/BenJoyenConseil/bluckdb)

It is a Key/Value store that implements bucketing based on [extendible hashing](https://en.wikipedia.org/wiki/Extendible_hashing)

The ``server.go`` file is a simple http server that answers on the 2233 port.


There are 3 endpoints :

    curl -XGET http://hostname:2233/v1/data/path/to/fs/?id=<key>
    curl -XPUT -d 'value' http://hostname:2233/v1/data/path/to/fs/?id=<key>
    curl -XGET http://hostname:2233/v1/meta/path/to/fs/
    curl -XGET http://hostname:2233/v1/debug/path/to/fs/?page_id=<id_of_the_page_to_display>


+ **/v1** is the version number of the API
+ **/path/to/fs/** is a filesystem absolute path to a folder (existing or not, it will create if not)
+ GET **/data** ```/path/to/fs/?id=key``` is the url to get the value of the specified key. It returns a JSON
+ PUT **/data** ```/path/to/fs/?id=key``` is the url to put the value of the specified key. It returns StatusOk if it is inserted
+ GET **/meta** ```/path/to/fs/``` is the url to print metadata of the table. It returns a JSON formatted response
+ PUT **/debug** ```/path/to/fs/?id_page=key``` is the url to print the dump of a specific page

## the goal

The goal of this project is to explore and to reinvent the wheel of well known, state of the art, algorithms and data structures.
For experimental and learning purpose only, not production ready.


## design

The datastructure is a persistent hashtable. 
There is no separation between the index and the data, they are in the same file
That is the reason why a record cannot be bigger than a bucket currently (no handling for big record)

#### meta directory 

A Directory is a table of buckets called "Page". 
It is the structure used to know where a bucket points physicaly to the data file 

    dir := &Directory{
          Table:      []int{0, 1, 3, 2, 0, 1, 3, 2},
          Gd:         2,
          LastPageId: 3,
       	 data: []byte{...},
    }

The ```Table``` array stores the ids of the pages (i.e pointers). To retrieve a page, the algorithm
hash(key) and lookups to the ```Table``` what index it points to. Then it does a subset
of the ```data``` byte slice with ```Table[ hash(key) ] * 4096```, and casts it to a Page structure.

When a new Page is required because there is no free space, ```LastPageId + 1``` is use to be the next 
pointer filled in the ```Table``` array to point to the new Page (so it points at the end of the File => ```(LastPageId + 1) * 4096```).

#### page layout
A Page is a byte array of 4096 bytes length, append only. 
Trailer : It stores actual usage of the Page at 4094 bytes (unint16), and local depth at 4092 bytes (unint16)

| Record 1 | Record 2 | Record 3 | Record 1 v2 | **LOCAL_DEPTH** | **PAGE_USE** |  

#### record layout

A Record is a byte array with a key, a value and the headers :
 
    type Record interface {
        Key() []byte
        Val() []byte
        KeyLen() uint16
        ValLen() uint16
    }
    
    type ByteRecord []byte


| ... | *k* | *e* | *y* | **v** | **a** | **l** | **u** | **e** | **0x5** | **0x0** | *0x3* | *0x0* | ... |


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

It runs an httpserver with an instance of MmapKVStore

## Benchmarks
    
    BenchmarkBluckDBPut                    1000000	            1359 ns/op   -> 1,3 µs
    BenchmarkBluckDBGet                    1000000	            1406 ns/op   -> 1,4 µs
    BenchmarkPutNaiveDiskKVStore-4          200000              6250 ns/op   -> 6,2 µs
    BenchmarkGetNaiveDiskKVStore-4              30          44017416 ns/op   ->  44 ms
    BenchmarkPutHashMap-4                  1000000              1385 ns/op   -> 1,3 µs
    BenchmarkGetHashMap-4                  2000000               711 ns/op   -> 0,7 µs


## Projects used by BluckDB

 * [Iris web framework](https://github.com/kataras/iris)
 * [Gommon logger](https://github.com/labstack/gommon/log)
 * [Mmap-go](https://github.com/edsrzf/mmap-go)