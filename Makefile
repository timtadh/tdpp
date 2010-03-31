build:
	echo "building"
	gobuild -a

_buildtest:
	gobuild -t

_runtest:
	./_testmain

test: _buildtest _runtest clean

.PHONY : clean
clean :
	-find . -name "*.6" | xargs -I"%s" rm %s
	-rm time _testmain *.6 2> /dev/null
	ls

