#include "callback.h"
#include "../../clib/xdag_runtime.h"
int init_password_callback(int is_testnet) {
    return xdag_set_password_callback_wrap(goPasswordCallback, is_testnet);
}