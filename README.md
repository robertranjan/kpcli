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
    export BACKUP_DIR="./bkups"
    export DATABASE_BACKUP="${BACKUP_DIR}/master-db.kdbx"

## create database

syntax:
    kpcli \
        --keyfile {keyfile} \
        --name {xyx.kdbx} \
        --pass {password to encrypt/decrypt kdbx} \
        create

eg: 1

    bin/kpcli \
        --keyfile ${KDBX_KEYFILE} \
        --database ${KDBX_DATABASE} \
        --pass ${KDBX_PASSWORD} \
        create

## list entries

syntax:

    ./kpcli --keyfile <keyfile> \
            --database <database-filename> \
        ls  [--reverse] [--days 10] [--sort-by-col 1|2|3|4]
            ; --reverse -> reverse order
            ; --sort-by-col N
                1 -> title
                2 -> history count
                3 -> creation time
                4 -> mod time
            ; --days 10 --> shows entries created or modified in the last 10 days

eg: 1

    bin/kpcli \
        --keyfile ${KDBX_KEYFILE} \
        --database ${KDBX_DATABASE} \
        --pass ${KDBX_PASSWORD} \
        ls 

eg: 2

    kpcli ls --sortby-col 4 -d 2
        ; shows entries modified in last 2 days ORDER by col 4(modified time)

## diff databases

example

    bin/kpcli \
        --keyfile ${KDBX_KEYFILE} \
        --database ${KDBX_DATABASE} \
        --pass ${KDBX_PASSWORD} \
        diff \
            --database2 ${DATABASE_BACKUP}

output

        Running diff between
            ${DATABASE_BACKUP} and
            ${KDBX_DATABASE}

        here are the diffs:
        ${BACKUP_DIR}/master-db.kdbx to ./tmp/master-db.kdbx
        ----------------------------------------------------------------------
        (removed)  Root/example entry 1
        ( added )  Root/example entry 2
        ( added )  Root/example entry 3

### steps

1. List all entries from ${KDBX_DATABASE} to ./database1.out
2. List all entries from ${DATABASE_BACKUP} to ./database2.out
3. Run a diff between ./database1.out and ./database2.out
