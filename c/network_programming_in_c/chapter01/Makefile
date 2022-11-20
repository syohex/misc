all: get_ip

get_ip: get_ip.c
	gcc -std=c17 -Wall -fsanitize=address,undefined -o $@ $^

clean:
	-rm -f *.o get_ip
