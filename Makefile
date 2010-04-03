
INSTALL_PREFIX = "$(SOURCEQL_HOME)/qplanner/tdpp"

test_build:
	echo "building"
	gobuild -I . -a

clean_install:
	-rm -r build
	-rm -r $(INSTALL_PREFIX)/parser
	-rm $(INSTALL_PREFIX)/*.a

install_build: clean_install
	mkdir build
	mkdir build/parser

	6g parser/set/set.go
	gopack crg set.a set.6
	6g parser/stack/stack.go
	gopack crg stack.a stack.6
	6g parser/token/token.go
	gopack crg token.a token.6
	find . -name "*.6" | xargs -I "%s" rm %s
	cp *.a build/parser
	rm *.a

	6g -I "build/" -o grammar.6 parser/grammar/gram.go parser/grammar/build.go
	gopack crg grammar.a grammar.6
	find . -name "*.6" | xargs -I "%s" rm %s
	cp *.a build/parser
	rm *.a

	6g -I "build/" -o stack.6 stack/tokenstack.go
	gopack crg stack.a stack.6
	find . -name "*.6" | xargs -I "%s" rm %s
	cp *.a build
	rm *.a

	6g -I "build/" parser/parser.go parser/processor.go
	gopack crg parser.a parser.6
	find . -name "*.6" | xargs -I "%s" rm %s
	cp *.a build
	rm *.a

install: install_build
	mkdir $(INSTALL_PREFIX)/parser
	cp -r build/* $(INSTALL_PREFIX)

_buildtest:
	gobuild -t

_runtest:
	./_testmain

test: _buildtest _runtest clean

.PHONY : clean
clean :
	find . -regextype posix-egrep -regex "(.*\.6)|(.*\.a)" | xargs -I"%s" rm %s
	-rm -f time _testmain 2> /dev/null
	ls
