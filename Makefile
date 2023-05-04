APPNAME = kpcli
DATETIMESTAMP = $(shell date +%Y-%m-%d_%H%M%S)
DATESTAMP = $(shell date +%Y-%m-%d)
range ?= {1..10}
YELLOW=\033[33m
GREEN=\033[32m
RESET=\033[0m
COLOR=$(YELLOW)
tmpDir = /tmp/$(APPNAME)
versionDetail=$(DATESTAMP).$(shell git rev-parse --short HEAD).$(shell git rev-list HEAD --count)

# for kpcli diff
DATABASE_BACKUP = ./bkup1/master-db.kdbx

# for kpcli createdb & kpcli ls
KDBX_KEYFILE="./tmp/master-db.key"
KDBX_DATABASE="./tmp/master-db.kdbx"
KDBX_PASSWORD='super_secret'

.PHONY: help
help:
	$(call banner, $@)
	@printf "\e[1;33mHere are the available targets:\e[32m\n"
	@make -qp 2> /dev/null | \
		awk -F':' '/^[a-zA-Z0-9][^$$#\/\t=]*:([^=]|$$)/ {split($$1,A,/ /);for(i in A)print A[i]}' | \
		sort -u | egrep -v "(Makefile)" | column
	@printf "\e[0m\n"
	echo versionDetail: $(versionDetail)

fix-link:
	$(call banner, $@)
	ln -sfn ./bin/$(APPNAME).bin $(APPNAME)

all: build fix-link
	$(call banner, $@)

.PHONY: build
build:
	$(call banner, $@)
	GOOS=darwin GOARCH=arm64 go build \
		-o bin/kpclidarwin-arm64.bin \
		-ldflags "-X main.Version=$(versionDetail)" \
		-trimpath "github.com/robertranjan/kpcli/cmd/kpcli" ./cmd/kpcli/
	GOOS=darwin GOARCH=amd64 go build \
		-o bin/kpcli-darwin-amd64.bin \
		-ldflags "-X main.Version=$(versionDetail)" \
		-trimpath "github.com/robertranjan/kpcli/cmd/kpcli" ./cmd/kpcli/
	GOOS=linux GOARCH=amd64 go build \
		-o bin/kpcli-linux-amd64.bin \
		-ldflags "-X main.Version=$(versionDetail)" \
		-trimpath "github.com/robertranjan/kpcli/cmd/kpcli" ./cmd/kpcli/

run: build
	$(call banner, $@)
	./kpcli ls

install:
	$(call banner, $@)
	cp bin/kpcli ~/go/bin/.
	cp bin/kpcli ~/bin/.

createdb: build
	rm -rf ./tmp && mkdir tmp
	bin/kpcli \
		--keyfile $(KDBX_KEYFILE) \
		--database $(KDBX_DATABASE) \
		--pass $(KDBX_PASSWORD) \
		createdb

ls: build
	bin/kpcli \
		--keyfile $(KDBX_KEYFILE) \
		--database $(KDBX_DATABASE) \
		--pass $(KDBX_PASSWORD) \
		ls

diff: build
	bin/kpcli \
		--keyfile $(KDBX_KEYFILE) \
		--database $(KDBX_DATABASE) \
		--pass $(KDBX_PASSWORD) \
		diff \
			--database2 ${DATABASE_BACKUP}

local-test:
	$(call banner, $@)
	@rm -rf ./tmp && mkdir tmp

	$(call banner, "create db - success")
	./kpcli -p "Super_Secret" --kf tmp/test.key --db tmp/test.kdbx createdb || true

	$(call banner, "create db - failure")
	./kpcli -p  "Super_Secret" -k tmp/test.key --db tmp/test.kdbx createdb || true

	$(call banner, "List entries sorted by creation time")
	@./kpcli -k ${KDBX_KEYFILE} --dbfile ${KDBX_DATABASE} --pass '${PASSWORD}' ls -c -d 30

	$(call banner, "List entries sorted by modification time")
	@./kpcli -k ${KDBX_KEYFILE}  --dbfile ${KDBX_DATABASE} --pass '${PASSWORD}' ls -m -d 3

	$(call banner, "List all fields")
	@./kpcli -k ${KDBX_KEYFILE} --dbfile ${KDBX_DATABASE} --pass "${PASSWORD}" ls -s t -f all -d 3
