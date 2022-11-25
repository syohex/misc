#include <stdio.h>
#include <sys/socket.h>
#include <netdb.h>
#include <unistd.h>

int main() {
    int sock = socket(AF_INET, SOCK_STREAM, 0);
    if (sock == -1) {
        perror("socket");
        return 1;
    }
    printf("open socket\n");

    close(sock);
    printf("close socket\n");

    return 0;
}
