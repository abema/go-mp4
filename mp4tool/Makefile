.PHONY: help build install clean

INST_PATH:=$(GOPATH)/bin/mp4tool
BIN:=build/mp4tool
PKG:=github.com/abema/go-mp4/mp4tool

CP:=cp
RM:=rm

help:
	@echo "USAGE:"
	@echo "  make build"
	@echo "  make install"
	@echo "  make clean"

install: $(INST_PATH)

build: $(BIN)

$(INST_PATH): $(BIN)
	$(CP) $(BIN) $(INST_PATH)

.PHONY: $(BIN)

$(BIN):
	go build -o $(BIN) $(PKG)

clean:
	$(RM) $(BIN)
