include $(GOROOT)/src/Make.inc
TARG=github.com/jpoirier/cpu

CGOFILES=\
	cpu.go\

CGO_OFILES=\
	cpu.o\

EXT=
ifeq ($(GOOS),windows)
EXT=.exe
endif

CLEANFILES+=example$(EXE)

include $(GOROOT)/src/Make.pkg

example: install example.go
	$(GC) example.go
	$(LD) -o $@$(EXT) example.$O
