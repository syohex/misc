all: get_ip.out get_hostname1.out get_hostname2.out

get_ip.out: get_ip.c
	gcc -std=c17 -Wall -fsanitize=address,undefined -o $@ $^

get_hostname1.out: get_hostname1.c
	gcc -std=c17 -Wall -fsanitize=address,undefined -o $@ $^

get_hostname2.out: get_hostname2.c
	gcc -std=c17 -Wall -fsanitize=address,undefined -o $@ $^

clean:
	-rm -f *.o *.out
