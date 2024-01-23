#ifndef XDAG_MINER_H
#define XDAG_MINER_H

#include <stdio.h>

#ifdef __cplusplus
extern "C" {
#endif

    extern struct dfslib_crypt *g_crypt;

    extern int client_init(int is_testnet);

#ifdef __cplusplus
};
#endif

#endif
