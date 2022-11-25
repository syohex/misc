#define _POSIX_C_SOURCE 200809L

#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <sys/socket.h>
#include <netdb.h>

int main() {
    struct addrinfo hints;
    memset(&hints, 0, sizeof(struct addrinfo));

    hints.ai_family = AF_INET;
    hints.ai_socktype = SOCK_STREAM;

    struct addrinfo *server;
    int ret = getaddrinfo("0.0.0.0", "8000", &hints, &server);
    if (ret != 0) {
        fprintf(stderr, "getaddrinfo: %d", ret);
        return 1;
    }

    int sock = socket(server->ai_family, server->ai_socktype, server->ai_protocol);
    if (sock == -1) {
        perror("socket");
        return 1;
    }

    freeaddrinfo(server);
    close(sock);
    printf("## done\n");
    return 0;
}
