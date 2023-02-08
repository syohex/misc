#include "server.h"

#include <stdlib.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <errno.h>
#include <unistd.h>

server_error ServerCreate(server_t **server, uint16_t port, int backlog) {
    server_error error = ERROR_OK;

    server_t *s = (server_t *)malloc(sizeof(server_t));
    if (s == NULL) {
        return ERROR_MEMORY;
    }

    s->port = port;

    s->sock = socket(AF_INET, SOCK_STREAM, 0);
    if (s->sock == -1) {
        error = ServerErrorFromErrno(errno);
        goto error;
    }

    int yes = 1;
    int ret = setsockopt(s->sock, SOL_SOCKET, SO_REUSEADDR, (const void *)&yes, sizeof(yes));
    if (ret == -1) {
        error = ServerErrorFromErrno(errno);
        goto error_sock;
    }

    struct sockaddr_in addr;
    addr.sin_family = AF_INET;
    addr.sin_port = htons(port);
    addr.sin_addr.s_addr = INADDR_ANY;

    ret = bind(s->sock, (struct sockaddr *)&addr, sizeof(addr));
    if (ret == -1) {
        error = ServerErrorFromErrno(errno);
        goto error_sock;
    }

    ret = listen(s->sock, backlog);
    if (ret == -1) {
        error = ServerErrorFromErrno(errno);
        goto error_sock;
    }

    *server = s;
    return error;

error_sock:
    close(s->sock);
error:
    free(s);
    return error;
}

void ServerDestroy(server_t *server) {
    if (server == NULL) {
        return;
    }

    close(server->sock);
    free(server);
}
