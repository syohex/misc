#define _POSIX_C_SOURCE 200809L

#include <sys/types.h>
#include <sys/socket.h>
#include <ifaddrs.h>
#include <netdb.h>
#include <stdio.h>

int main() {
    struct ifaddrs *adapters;

    int ret = getifaddrs(&adapters);
    if (ret == -1) {
        perror("getifaddrs");
        return 1;
    }

    struct ifaddrs *it = adapters;
    while (it != NULL) {
        int family = it->ifa_addr->sa_family;
        if (family == AF_INET || family == AF_INET6) {
            char buf[128];
            size_t size = family == AF_INET ? sizeof(struct sockaddr_in) : sizeof(struct sockaddr_in6);

            getnameinfo(it->ifa_addr, size, buf, 128, 0, 0, NI_NUMERICHOST);

            printf("## name=%s, IP%s, %s\n", it->ifa_name, family == AF_INET ? "v4" : "v6", buf);
        }

        it = it->ifa_next;
    }

    freeifaddrs(adapters);
    return 0;
}
