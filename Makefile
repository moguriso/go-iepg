BUILD		:= release
HASH=$(shell git rev-parse --short HEAD)
BUILDDATE=$(shell date '+%Y/%m/%d %H:%M:%S %Z')
GOVERSION=$(shell go version)
VERSION     := 1.0.1

SRCS		:= $(shell find ./ -type f -name '*.go')
LDFLAGS		:= -ldflags="-s -w -X \"main.version=$(VERSION)\" -X \"main.hash=$(HASH)\" -X \"main.builddate=$(BUILDDATE)\" -X \"main.goversion=$(GOVERSION)\" -extldflags "

BuildCommand = GOARCH=$(1) GOOS=windows go build -a -v -tags '$(2)' $(3) 

ifeq ($(BUILD),debug)
TAGS	:= debug
IS_DBG	:= -debug
else
TAGS	:= 
IS_DBG	:= 
endif

.PHONY: all
all: amd64

.PHONY: clean
clean:
	@rm -rf go-iepg.exe
	@rm -rf go-iepg

386: $(SRCS)
	$(call BuildCommand,386,$(TAGS),$(LDFLAGS))

arm: $(SRCS)
	$(call BuildCommand,arm,$(TAGS),$(LDFLAGS))

mips: $(SRCS)
	$(call BuildCommand,mips,$(TAGS),$(LDFLAGS))

amd64: $(SRCS)
	$(call BuildCommand,amd64,$(TAGS),$(LDFLAGS))

