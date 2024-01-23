//
// Created by swordlet on 2021/3/24.
//

#ifndef DAGWALLET_XDAG_RUNTIME_H
#define DAGWALLET_XDAG_RUNTIME_H


#define _TIMESPEC_DEFINED

#include "xDagWallet/src/client/common.h"
#include "xDagWallet/src/client/client.h"
#include "xDagWallet/src/client/utils/utils.h"
#include "xDagWallet/src/client/dnet_crypt.h"
#include "xDagWallet/src/client/wallet.h"

////---- Duplicated from dnet_crypt.c ----
#define KEYFILE	    "xdagj_dat" DELIMITER "dnet_key.dat"
struct dnet_keys {
    struct dnet_key priv;
    struct dnet_key pub;
};

typedef int(*password_callback)(const char *prompt, char *buf, unsigned size);
////------------------------------------
#ifdef __cplusplus
extern "C" {
#endif
////---- Exporting functions ----

extern int xdag_set_password_callback_wrap(password_callback callback, int is_testnet);
extern void* xdag_get_default_key(void);

#ifdef __cplusplus
};
#endif
#endif //DAGWALLET_XDAG_RUNTIME_H
