#define _POSIX_C_SOURCE 200809L

#include <stdio.h>
#include <unistd.h>
#include <netdb.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <arpa/inet.h>

int main() {
    char host_name[1024];

    int ret = gethostname(host_name, 1023);
    if (ret == -1) {
        perror("gethostname");
        return 1;
    }

    printf("HostName is '%s'\n", host_name);

    struct hostent *host = gethostbyname(host_name);
    if (host == NULL) {
        perror("gethostbyname");
        return 1;
    }

    char **addresses = host->h_addr_list;
    while (*addresses != NULL) {
        printf("## %s\n", inet_ntoa(*(struct in_addr *)*addresses));
        ++addresses;
    }

    return 0;
}
