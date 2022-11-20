#define _POSIX_C_SOURCE 200809L

#include <sys/types.h>
#include <sys/socket.h>
#include <netdb.h>
#include <string.h>
#include <stdio.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <errno.h>

int main(void) {
    struct addrinfo hints;
    memset(&hints, 0, sizeof(struct addrinfo));
    hints.ai_family = AF_INET;
    hints.ai_socktype = SOCK_DGRAM;

    const char *dns_server = "8.8.8.8";
    const char *dns_port = "53";
    struct addrinfo *peer_addr = NULL;

    int ret = getaddrinfo(dns_server, dns_port, &hints, &peer_addr);
    if (ret != 0) {
        perror("getaddrinfo");
        return 1;
    }

    int sock = socket(peer_addr->ai_family, peer_addr->ai_socktype, peer_addr->ai_protocol);
    if (sock == -1) {
        perror("socket");
        goto error;
    }

    ret = connect(sock, peer_addr->ai_addr, peer_addr->ai_addrlen);
    if (ret == -1) {
        perror("connect");
        goto error;
    }

    struct sockaddr_in name;
    socklen_t len = sizeof(name);
    ret = getsockname(sock, (struct sockaddr *)&name, &len);
    if (ret == -1) {
        perror("getsockname");
        goto error;
    }

    char buffer[256];
    const char *p = inet_ntop(peer_addr->ai_family, &name.sin_addr, buffer, 255);
    if (p != NULL) {
        printf("Local IP address = %s\n", buffer);
    } else {
        printf("Failed to get IP address:%s\n", strerror(errno));
    }

error:
    if (sock != -1) {
        close(sock);
    }

    if (peer_addr != NULL) {
        freeaddrinfo(peer_addr);
    }

    return 0;
}
