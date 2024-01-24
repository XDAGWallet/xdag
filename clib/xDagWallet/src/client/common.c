//
//  common.c
//  xDagWallet
//
//  Created by Rui Xie on 7/11/18.
//  Copyright © 2018 xrdavies. All rights reserved.
//

#include "common.h"
#include "dnet_crypt.h"

int g_xdag_testnet = 0;

/* see dnet_user_crypt_action */
int xdag_user_crypt_action(unsigned *data, unsigned long long data_id, unsigned size, int action)
{
    return dnet_user_crypt_action(data, data_id, size, action);
}
