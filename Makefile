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
LOG_LEVEL ?= error
BACKUP_DIR = ./bkups

GitTagLocal = $(shell git tags | tail -n1 | awk '{print $$1}')
GitTagRemote = $(shell git ls-remote --tags 2>/dev/null | tail -n1 | awk -F'/' '{print $$3}')
appVersionConfig = $(shell awk '/"version":/ {print $$2}' butler.json | tr -d '",')
appVersionSrc = $(shell awk -F'=' '/var Version/ {gsub(/"/,DD,$$2);gsub(/ /,"",$$2);print $$2}' version/version.go)
DOT_FILE = ./tmp/flow-dia.dot
GOTRACE_BIN = ~/go/bin/gotrace
DOT_OUTFILE = ./tmp/gotrace.png


define banner
@printf "############################################\n"
@printf "# $@ # \n"
@printf "############################################\n"
endef

# for kpcli diff
DATABASE_BACKUP = $(BACKUP_DIR)/master-db.kdbx

# for kpcli create & kpcli ls
KDBX_KEYFILE="./tmp/master-db.key"
KDBX_DATABASE="./tmp/master-db.kdbx"
KDBX_PASSWORD='$(PASSWORD)'

# for kpcli create & kpcli ls
KDBX_KEYFILE2="./tmp/master-db.key"
KDBX_DATABASE2="$(BACKUP_DIR)/master-db.kdbx"
KDBX_PASSWORD2='$(PASSWORD)'

# .PHONY: help
# help:
# 	$(call banner, $@)
# 	@printf "\e[1;33mHere are the available targets:\e[32m\n"
# 	@make -qp 2> /dev/null | \
# 		awk -F':' '/^[a-zA-Z0-9][^$$#\/\t=]*:([^=]|$$)/ {split($$1,A,/ /);for(i in A)print A[i]}' | \
# 		sort -u | egrep -v "(Makefile)" | column
# 	@printf "\e[0m\n"
# 	echo versionDetail: $(versionDetail)

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

install: build
	$(call banner, $@)
	cp $(BIN) ~/go/bin/.
	cp $(BIN) ~/bin/.

create: build
	$(call banner, $@)
	rm -rf ./tmp && mkdir tmp
	$(BIN) \
		--log-level $(LOG_LEVEL) \
		--keyfile $(KDBX_KEYFILE) \
		--database $(KDBX_DATABASE) \
		--pass $(KDBX_PASSWORD) \
		create

ls: build
	$(call banner, $@)
	$(BIN) \
		--log-level $(LOG_LEVEL) \
		--keyfile $(KDBX_KEYFILE) \
		--database $(KDBX_DATABASE) \
		--pass $(KDBX_PASSWORD) \
		ls

diff: build
	$(call banner, $@)
	$(BIN) \
		--log-level $(LOG_LEVEL) \
		--keyfile $(KDBX_KEYFILE) \
		--database $(KDBX_DATABASE) \
		--pass $(KDBX_PASSWORD) \
		diff \
			--keyfile2 $(KDBX_KEYFILE2) \
			--database2 $(KDBX_DATABASE2) \
			--pass2 $(KDBX_PASSWORD2)

local-test:
	$(call banner, $@)
	@rm -rf ./tmp && mkdir tmp

	@echo "---- create db - success ----"
	$(BIN) --log-level $(LOG_LEVEL) -p '${PASSWORD}' --kf ${KDBX_KEYFILE}
		--db ${KDBX_DATABASE} create || true

	@echo "---- create db - failure ----"
	$(BIN) --log-level $(LOG_LEVEL) -p  '${PASSWORD}' -k ${KDBX_KEYFILE}
		--db ${KDBX_DATABASE} create || true

	@echo "---- List entries sorted by creation time ----"
	@$(BIN) --log-level $(LOG_LEVEL) -k ${KDBX_KEYFILE}
		--db ${KDBX_DATABASE} --pass '${PASSWORD}' ls --sortby-col 3 -d 30

	@echo "---- List entries sorted by modification time ----"
	@$(BIN) --log-level $(LOG_LEVEL) -k ${KDBX_KEYFILE}
		--db ${KDBX_DATABASE} --pass '${PASSWORD}' ls --sortby-col 4 -d 3

	@echo "---- List all fields ----"
	@$(BIN) --log-level $(LOG_LEVEL) -k ${KDBX_KEYFILE}
		--db ${KDBX_DATABASE} --pass "${PASSWORD}" ls --sortby-col 1 -f all -d 3

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##/ /p' ${MAKEFILE_LIST} | column -t -s ':'

## release-check: check local & remoate git tags and suggest to create and or push tags to remote
release-check:
	@printf "\n - Git tag:\n"
	@printf "   - $(GitTagLocal): local\n"
	@printf "   - $(GitTagRemote): remote \n\n"
	@printf " - Version from files:\n"
	@printf "   - $(appVersionConfig): butler.json \n"
	@printf "   - $(appVersionSrc): version/version.go\n\n"
ifeq ($(GitTagRemote),$(GitTagLocal))
	@echo " - Version: $(appVersionConfig) is already on upstream. Either create new tag or overwrite and push forcefully."
	@printf "     create new tag      : git tag $(appVersionConfig)+1 && git push origin --tags \n"
	@printf "     overwrite forcefully: git tag --force $(appVersionConfig) && git push origin --tags --force \n"
else
	@printf " - tag: $(GitTagLocal) is not in remote, push it to remote using below cmd\n"
	@printf "     - merge PR on GitHub UI\n       git co main && git pull\n"
	@printf "       git tag $(appVersionConfig) && git push origin --tags \n"
endif
