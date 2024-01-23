//
//  common.h
//  xDagWallet
//
//  Created by Rui Xie on 7/11/18.
//  Copyright © 2018 xrdavies. All rights reserved.
//

#ifndef common_h
#define common_h

#include <time.h>
#include "errno.h"

#define COINNAME "XDAG"

// This is for timeval redefinition issue for pthread
#define HAVE_STRUCT_TIMESPEC

// This is to disable security warning from Visual C++
#ifdef _MSC_VER
#define _CRT_SECURE_NO_WARNINGS
#endif


#ifdef __cplusplus
extern "C" {
#endif

    extern int g_xdag_testnet;
    /* see dnet_user_crypt_action */
    extern int xdag_user_crypt_action(unsigned *data, unsigned long long data_id, unsigned size, int action);

#ifdef __cplusplus
}
#endif

#endif /* common_h */
