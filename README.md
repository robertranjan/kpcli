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

## setup env vars

    export KDBX_KEYFILE="./tmp/master-db.key"
    export KDBX_DATABASE="./tmp/master-db.kdbx"
    export KDBX_PASSWORD="super_secret"
    export DATABASE_BACKUP="./bkup1/master-db.kdbx"

## createdb

    bin/kpcli \
        --keyfile ${KDBX_KEYFILE} \
        --database ${KDBX_DATABASE} \
        --pass ${KDBX_PASSWORD} \
        createdb

## ls

    bin/kpcli \
        --keyfile $(KDBX_KEYFILE) \
        --database $(KDBX_DATABASE) \
        --pass $(KDBX_PASSWORD) \
        ls 

## diff

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

### diff between kdbx

1. Find the previous backup and decrupt to tmp/
2. List all entries from tmp/{}.kdbx
3. List all entries from actual/{}.kdbx
4. Run a diff and show the diffs
