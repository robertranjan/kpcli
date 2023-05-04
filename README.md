# build and install

- clone and install

        git clone
        cd kpcli
        make build install

- Install using go install

        go install github.com/robertranjan/kpcli@latest

    this get installed to ${GOBIN}

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

## diff 2 databases

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
        ${DATABASE_BACKUP} to ${KDBX_DATABASE}
        ----------------------------------------------------------------------
        (removed)  Root/H&R handr block
        ( added )  Root/H&R handr block 2022
        ( added )  Root/Income Tax

### steps

1. List all entries from KDBX_DATABASE.kdbx to ./database1.out
2. List all entries from DATABASE_BACKUP.kdbx to ./database2.out
3. Run a diff between ./database1.out and ./database2.out and show the diffs
