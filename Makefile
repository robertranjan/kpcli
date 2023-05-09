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
BIN = bin/kpcli
PASSWORD ?= super_secret


define banner
@printf "############################################\n"
@printf "# $@ # \n"
@printf "############################################\n"
endef

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
	GOOS=darwin GOARCH=amd64 go build \
		-o bin/kpcli-darwin-amd64.bin \
		-ldflags "-X main.Version=$(versionDetail)" \
		-trimpath "github.com/robertranjan/kpcli" .
	cp bin/kpcli-darwin-amd64.bin $(BIN)
#	GOOS=darwin GOARCH=arm64 go build \
#		-o bin/kpclidarwin-arm64.bin \
#		-ldflags "-X main.Version=$(versionDetail)" \
#		-trimpath "github.com/robertranjan/kpcli/cmd/kpcli" .
#	GOOS=linux GOARCH=amd64 go build \
#		-o bin/kpcli-linux-amd64.bin \
#		-ldflags "-X main.Version=$(versionDetail)" \
#		-trimpath "github.com/robertranjan/kpcli/cmd/kpcli" .

run: build
	$(call banner, $@)
	$(BIN) ls

install:
	$(call banner, $@)
	cp $(BIN) ~/go/bin/.
	cp $(BIN) ~/bin/.

createdb: build
	rm -rf ./tmp && mkdir tmp
	$(BIN) \
		--keyfile $(KDBX_KEYFILE) \
		--database $(KDBX_DATABASE) \
		--pass $(KDBX_PASSWORD) \
		createdb

ls: build
	$(BIN) \
		--keyfile $(KDBX_KEYFILE) \
		--database $(KDBX_DATABASE) \
		--pass $(KDBX_PASSWORD) \
		ls

diff: build
	$(BIN) \
		--keyfile $(KDBX_KEYFILE) \
		--database $(KDBX_DATABASE) \
		--pass $(KDBX_PASSWORD) \
		diff \
			--database2 ${DATABASE_BACKUP}

local-test:
	$(call banner, $@)
	@rm -rf ./tmp && mkdir tmp

	@echo "---- create db - success ----"
	$(BIN) -p '${PASSWORD}' --kf ${KDBX_KEYFILE} --db ${KDBX_DATABASE} createdb || true

	@echo "---- create db - failure ----"
	$(BIN) -p  '${PASSWORD}' -k ${KDBX_KEYFILE} --db ${KDBX_DATABASE} createdb || true

	@echo "---- List entries sorted by creation time ----"
	@$(BIN) -k ${KDBX_KEYFILE} --dbfile ${KDBX_DATABASE} --pass '${PASSWORD}' ls --sortby-col 3 -d 30

	@echo "---- List entries sorted by modification time ----"
	@$(BIN) -k ${KDBX_KEYFILE}  --dbfile ${KDBX_DATABASE} --pass '${PASSWORD}' ls --sortby-col 4 -d 3

	@echo "---- List all fields ----"
	@$(BIN) -k ${KDBX_KEYFILE} --dbfile ${KDBX_DATABASE} --pass "${PASSWORD}" ls --sortby-col 1 -f all -d 3
