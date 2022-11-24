all: getaddrinfo.out

getaddrinfo.out: getaddrinfo.c
	gcc -std=c17 -Wall -fsanitize=address,undefined -o $@ $^

.PHONY: clean
clean:
	-rm -f *.o *.out
