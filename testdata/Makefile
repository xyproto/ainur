.PHONY:all clean

CC = tcc
CXX = clang++
CXXFLAGS ?= -O2

all: clang_hello tcc_hello rust_hello rust_hello_stripped go_hello

clang_hello: clang_hello.cpp
	${CXX} ${CXXFLAGS} $< -o $@
	chmod -x $@

tcc_hello: tcc_hello.c
	${CC} ${CXXFLAGS} $< -o $@
	chmod -x $@

rust_hello: rust_hello.rs
	rustc $<
	chmod -x $@

rust_hello_stripped: rust_hello.rs
	rustc $< -o $@
	strip $@
	chmod -x $@

go_hello: main.go
	GOOS=linux GOARCH=amd64 go build -o go_hello main.go

clean:
	rm -f clang_hello tcc_hello rust_hello go_hello
