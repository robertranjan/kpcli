# kpcli

[![Go Report Card](https://goreportcard.com/badge/github.com/robertranjan/kpcli)](https://goreportcard.com/report/github.com/robertranjan/kpcli)

## build and install

- clone and install

        git clone
        cd kpcli
        make build install

- Install using go install

        go install github.com/robertranjan/kpcli@latest

    this will install the tool at ${GOBIN} dir(usually ~/go/bin)

- Usage

        kpcli --version

## setup env vars

    export KDBX_KEYFILE="./tmp/master-db.key"
    export KDBX_DATABASE="./tmp/master-db.kdbx"
    export KDBX_PASSWORD="super_secret"
    export DATABASE_BACKUP="./bkup1/master-db.kdbx"

## create database

    bin/kpcli \
        --keyfile ${KDBX_KEYFILE} \
        --database ${KDBX_DATABASE} \
        --pass ${KDBX_PASSWORD} \
        createdb

## list entries

    bin/kpcli \
        --keyfile $(KDBX_KEYFILE) \
        --database $(KDBX_DATABASE) \
        --pass $(KDBX_PASSWORD) \
        ls 

## diff databases

example

    bin/kpcli \
        --keyfile $(KDBX_KEYFILE) \
        --database $(KDBX_DATABASE) \
        --pass $(KDBX_PASSWORD) \
        diff \
            --database2 ${DATABASE_BACKUP}

output

        Running diff between
            ${DATABASE_BACKUP} and
            ${KDBX_DATABASE}

        here are the diffs:
        ./bkup1/master-db.kdbx to ./tmp/master-db.kdbx
        ----------------------------------------------------------------------
        (removed)  Root/H&R handr block
        ( added )  Root/H&R handr block 2022
        ( added )  Root/Income Tax

### steps

1. List all entries from ${KDBX_DATABASE} to ./database1.out
2. List all entries from ${DATABASE_BACKUP} to ./database2.out
3. Run a diff between ./database1.out and ./database2.out
