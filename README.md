# build and install

- clone and install

        git clone
        cd kpcli
        make build install

- Install using go install

        go install @latest

    this get installed to ${GOBIN}

- Usage

        kpcli --version

## setup

    export KEYFILE="./tmp/master-db.key"
    export DATABASE_BACKUP="./tmp/master-db.2023-01-01.kdbx"
    export DATABASE="./tmp/master-db.kdbx"
    export MY_DB_PASS="super_secret"

## createdb

    mkdir ./tmp
    export TEST_KEYFILE="./tmp/master-db.key"
    export TEST_DATABASE="./tmp/master-db.kdbx"
    export TEST_PASS="super_secret"

    bin/kpcli \
        --keyfile ${TEST_KEYFILE} \
        --database ${TEST_DATABASE} \
        --pass ${TEST_PASS} \
        createdb

## ls

    kpcli ls

## diff

example 1

    kpcli \
        --keyfile ${KEYFILE} \
        --database ${DATABASE} \
        --pass ${MY_DB_PASS} \
        diff \
            --database2 ${DATABASE_BACKUP}

output

        Running diff between
            ${DATABASE_BACKUP} and
            ${DATABASE}

        here are the diffs:
        ${DATABASE_BACKUP} to ${DATABASE}
        ----------------------------------------------------------------------
        (removed)  Root/H&R handr block
        ( added )  Root/H&R handr block 2022
        ( added )  Root/Income Tax

### diff between kdbx

1. Find the previous backup and decrupt to tmp/
2. List all entries from tmp/{}.kdbx
3. List all entries from actual/{}.kdbx
4. Run a diff and show the diffs
