#include "error.h"

#include <errno.h>

server_error ServerErrorFromErrno(int error) {
    switch (errno) {
    case EACCES:
        return ERROR_ACCESS;
    case EINVAL:
        return ERROR_INVALID;
    default:
        return ERROR_UNKNOWN;
    }
}
