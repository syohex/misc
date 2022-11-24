#define _POSIX_C_SOURCE 200809L

#include <stdio.h>
#include <string.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netdb.h>
#include <arpa/inet.h>

int main() {
    struct addrinfo hints;
    memset(&hints, 0, sizeof(struct addrinfo));
    hints.ai_family = AF_INET;
    hints.ai_socktype = SOCK_STREAM;

    struct addrinfo *resource;
    int ret = getaddrinfo("syohex.org", "443", &hints, &resource);
    if (ret != 0) {
        perror("getaddrinfo");
        return 1;
    }

    printf("## syohex.org=%s\n", inet_ntoa(((struct sockaddr_in *)(resource->ai_addr))->sin_addr));
    freeaddrinfo(resource);
    return 0;
}
