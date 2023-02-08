#pragma once

typedef enum {
    ERROR_OK,
    ERROR_MEMORY,
    ERROR_INVALID,
    ERROR_ACCESS,
    ERROR_UNKNOWN,
} server_error;

server_error ServerErrorFromErrno(int error);
