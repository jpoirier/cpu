include $(GOROOT)/src/Make.inc
TARG=proctopo

CGOFILES=\
	proctopo.go

LDPATH_freebsd=-Wl,-R,`pwd`
LDPATH_linux=-Wl,-R,`pwd`
LDPATH_darwin=
LDPATH_windows=

#VERS=uname -r
#OSXVER=$shell(sw_vers | awk '/ProductVersion/ {print $2}')

ifeq ($(GOOS),darwin)
CFLAGS+=-D__DARWIN__
else ifeq ($(GOOS),freebsd)
CFLAGS+=-D__FREEBSD__
else ifeq ($(GOOS),linux)
CFLAGS+=-D__LINUX__
else ifeq ($(GOOS),windows)
CFLAGS+=-D__WINDOWS__
EXT=.exe
else
$(error Invalid $$GOOS '$(GOOS)'; must be darwin, freebsd, linux, or windows)
endif

ifeq ($(GOARCH),amd64)
CFLAGS+=-D__AMD64__
else ifeq ($(GOARCH),386)
CFLAGS+=-D__386__
else
$(error Invalid $$GOARCH '$(GOARCH)'; must be 386 or amd64)
endif


CGO_LDFLAGS=proctopo.so $(LDPATH_$(GOOS))
CGO_DEPS=proctopo.so

CLEANFILES +=main$(EXT)
include $(GOROOT)/src/Make.pkg

proctopo.o: proctopo.c
	gcc $(_CGO_CFLAGS_$(GOARCH)) -g -c -fPIC $(CFLAGS) proctopo.c

proctopo.so: proctopo.o
	gcc $(_CGO_CFLAGS_$(GOARCH)) -o $@ proctopo.o $(_CGO_LDFLAGS_$(GOOS))

main: install main.go
	$(GC) main.go
	$(LD) -o $@ main.$O
