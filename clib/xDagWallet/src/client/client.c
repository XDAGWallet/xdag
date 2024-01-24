#include <stdio.h>


#if defined(_WIN32) || defined(_WIN64)
#else
#include <unistd.h>
#endif

#include "crypt.h"
#include "wallet.h"
#include "common.h"
#include "client.h"
#include "dnet_crypt.h"


int client_init(int is_testnet)
{
    g_xdag_testnet = is_testnet;

    printf("Starting xdag, version 0.1.0\n");
    if (dnet_crypt_init(DNET_VERSION)) {
        sleep(3);
        printf("Password incorrect.\n");
        return -1;

    }

    printf("Initializing cryptography...\n");
    if (xdag_crypt_init(1)) {
        printf("Init crypto failed.\n");
        return -2;
    }

    printf("Reading wallet...\n");
    if (xdag_wallet_init()) {
        printf("Init wallet failed.\n");
        return -3;
    }

    return 0;
}
