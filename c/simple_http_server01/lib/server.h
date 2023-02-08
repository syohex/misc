#pragma once

#include <stdint.h>

#include "error.h"

typedef struct {
    int sock;
    uint16_t port;
} server_t;

server_error ServerCreate(server_t **server, uint16_t port, int backlog);
void ServerDestroy(server_t *server);
