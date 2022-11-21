#define _POSIX_C_SOURCE 200809L

#include <stdio.h>
#include <unistd.h>

int main() {
    char host_name[1024];

    int ret = gethostname(host_name, 1023);
    if (ret == -1) {
        perror("gethostname");
        return 1;
    }

    printf("HostName is '%s'\n", host_name);
    return 0;
}
