/* header of dfslib library for random data generation, T10.273-T10.273; $DVS:time$ */

#ifndef DFSLIB_RANDOM_H_INCLUDED
#define DFSLIB_RANDOM_H_INCLUDED

#include "dfslib_string.h"
#include "dfslib_crypt.h"

#ifdef __cplusplus
extern "C" {
#endif

    extern void dfslib_random_init(void);
    extern unsigned dfslib_random_get(unsigned limit);
    extern void dfslib_random_fill(void *buf, unsigned long len, int xor, struct dfslib_string *tip);
    extern void dfslib_random_sector(dfs32 *sector, struct dfslib_crypt *crypt, struct dfslib_string *password, struct dfslib_string *tip);

#ifdef __cplusplus
};
#endif

#endif
