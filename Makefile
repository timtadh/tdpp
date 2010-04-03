
INSTALL_PREFIX = "$(SOURCEQL_HOME)/qplanner"

build:
	echo "building"
	gobuild -a

install:
# 	go build library building broken
# 	gobuild -lib=true
	6g set/set.go
	gopack crg set.a set.6
	rm set.6
	6g stack/stack.go
	gopack crg stack.a stack.6
	rm stack.6
	6g -I . parser/parser.go parser/gram.go parser/build.go parser/token.go
	gopack crg parser.a parser.6
	rm parser.6
# 	cp *.a $(GOROOT)/pkg/$(GOOS)_$(GOARCH)
	cp *.a $(INSTALL_PREFIX)

_buildtest:
	gobuild -t

_runtest:
	./_testmain

test: _buildtest _runtest clean

.PHONY : clean
clean :
	-find . -name "*.6" | xargs -I"%s" rm %s
	-rm -f time _testmain *.6 *.a 2> /dev/null
	ls
