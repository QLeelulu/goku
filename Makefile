# include this file which comes with go installation
include ${GOROOT}/src/Make.dist 

TARG=./

# the *_test.go files are not to be included here. Only those your would use to build the actual program.  gotest will figure out the *_test.go files for itself.
GOFILES=\
 route.go\
 utils/utils.go\

# include this file which comes with go installation
# include ${GOROOT}/src/Make.pkg