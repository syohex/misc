all: getaddrinfo.out client01.out client02.out

getaddrinfo.out: getaddrinfo.c
	gcc -std=c17 -Wall -fsanitize=address,undefined -o $@ $^

client01.out: client01.c
	gcc -std=c17 -Wall -fsanitize=address,undefined -o $@ $^

client02.out: client02.c
	gcc -std=c17 -Wall -fsanitize=address,undefined -o $@ $^

.PHONY: clean
clean:
	-rm -f *.o *.out
