#define _XOPEN_SOURCE
#include <assert.h>
#include <stdio.h>
#include <string.h>
#include <time.h>
#include <stdlib.h>

struct Date {
    int year;
    int month;
    int month_day;
};

static int this_month_days(struct tm *t) {
    static int days[] = {31, -1, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31};
    if (t->tm_mon != 1) {
        return days[t->tm_mon];
    }

    if (t->tm_year % 400 == 0) {
        return 29;
    }
    if (t->tm_year % 100 == 0) {
        return 28;
    }
    if (t->tm_year % 4 == 0) {
        return 29;
    }

    return 28;
}

static int calculate_next_monday(struct Date *date) {
    time_t now = time(NULL);
    if (now == -1) {
        perror("time");
        return -1;
    }

    struct tm now_tm;
    if (localtime_r(&now, &now_tm) == NULL) {
        perror("localtime_r");
        return -1;
    }

    int next_monday;
    if (now_tm.tm_wday == 0) {
        // sunday
        next_monday = now_tm.tm_mday + 1;
    } else {
        next_monday = now_tm.tm_mday + 7 + (1 - now_tm.tm_wday);
    }

    int year = now_tm.tm_year + 1900;
    int month = now_tm.tm_mon + 1;
    int month_days = this_month_days(&now_tm);
    if (next_monday > month_days) {
        next_monday -= month_days;
        month += 1;
        if (month > 12) {
            year += 1;
            month = 1;
        }
    }

    date->year = year;
    date->month = month;
    date->month_day = next_monday;

    return 0;
}

void print_date(struct Date *d, const char *msg) {
    printf("%s %d/%02d/%02d\n", msg, d->year, d->month, d->month_day);
}

int main() {
    {
        struct Date d;
        if (calculate_next_monday(&d) != 0) {
            printf("failed\n");
            return -1;
        }
        print_date(&d, "Next monday is ");
    }
    return 0;
}
