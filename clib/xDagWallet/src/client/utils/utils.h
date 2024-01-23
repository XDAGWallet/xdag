//
//  utils.h
//  xdag
//
//  Copyright © 2018 xdag contributors.
//

#ifndef XDAG_UTILS_HEADER_H
#define XDAG_UTILS_HEADER_H

#include <stdio.h>
#include <stdint.h>

#ifdef _WIN32
#define DELIMITER "\\"
#else
#define DELIMITER "/"
#endif
typedef uint64_t xdag_time_t;
#ifdef __cplusplus
extern "C" {
#endif

    extern FILE* xdag_open_file(const char *path, const char *mode);
    extern void xdag_close_file(FILE *f);

#ifdef __cplusplus
};
#endif

#endif /* utils_h */
