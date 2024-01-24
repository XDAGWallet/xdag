//
// Created by swordlet on 2021/3/24.
//
#include "xdag_runtime.h"

int xdag_set_password_callback_wrap(password_callback callback, int is_testnet)
{
    xdag_user_crypt_action((uint32_t *)(callback), 0, 0, 6);
    //// return xdag_set_password_callback(callback);
    return client_init(is_testnet);
}

void *xdag_get_default_key()
{
    return xdag_default_key();
}
