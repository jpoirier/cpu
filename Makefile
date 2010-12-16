
all:
	$(MAKE) -f Makefile.obj
	$(MAKE) -f Makefile.pkg install
	$(MAKE) -f Makefile.obj example

clean:
	$(MAKE) -f Makefile.pkg clean




