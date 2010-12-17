include $(GOROOT)/src/Make.inc
TARG=cpu

CGOFILES=\
	cpu.go\

CGO_OFILES=\
	cpu.o\

CLEANFILES+=example
include $(GOROOT)/src/Make.pkg

ifeq ($(GOOS),darwin)
CGO_CFLAGS +=-D__DARWIN__
else ifeq ($(GOOS),freebsd)
CGO_CFLAGS +=-D__FREEBSD__
else ifeq ($(GOOS),linux)
CGO_CFLAGS +=-D__LINUX__
else ifeq ($(GOOS),windows)
CGO_CFLAGS +=-D__WINDOWS__
endif

ifeq ($(GOARCH),amd64)
CGO_CFLAGS +=-D__AMD64__
else ifeq ($(GOARCH),386)
CGO_CFLAGS +=-D__386__
endif

example: install example.go
	$(GC) example.go
	$(LD) -o $@ example.$O
