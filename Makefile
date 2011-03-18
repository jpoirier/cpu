include $(GOROOT)/src/Make.inc
TARG=cpu

CGOFILES=\
	cpu.go\

CGO_OFILES=\
	cpu.o\

CLEANFILES+=example

include $(GOROOT)/src/Make.pkg

EXT=
ifeq ($(GOOS),windows)
EXT=.exe
endif

example: install example.go
	$(GC) example.go
	$(LD) -o $@$(EXT) example.$O
