/* header for dfslib cryptography engine, T4.846-T10.273; $DVS:time$ */

#ifndef DFSLIB_CRYPT_H_INCLUDED
#define DFSLIB_CRYPT_H_INCLUDED

#include "dfslib_types.h"
#include "dfslib_string.h"

#define DFSLIB_CRYPT_PWDLEN 4

struct dfslib_crypt {
	dfs32	regs[0x10000];
	dfs32	pwd[DFSLIB_CRYPT_PWDLEN];
	dfs32	ispwd;
};

#ifdef __cplusplus
extern "C" {
#endif

    extern int dfslib_crypt_set_password(struct dfslib_crypt *dfsc, const struct dfslib_string *password);
    extern int dfslib_crypt_copy_password(struct dfslib_crypt *to, const struct dfslib_crypt *from);
    extern int dfslib_crypt_is_password(struct dfslib_crypt *dfsc);
    extern int dfslib_crypt_set_sector0(struct dfslib_crypt *dfsc, const void *sector);
    extern int dfslib_encrypt_sector(struct dfslib_crypt *dfsc, dfs32 *sector, dfs64 sectorNo);
    extern int dfslib_uncrypt_sector(struct dfslib_crypt *dfsc, dfs32 *sector, dfs64 sectorNo);
    extern int dfslib_encrypt_array(struct dfslib_crypt *dfsc, dfs32 *data, unsigned size, dfs64 sectorNo);
    extern int dfslib_uncrypt_array(struct dfslib_crypt *dfsc, dfs32 *data, unsigned size, dfs64 sectorNo);

#ifdef __cplusplus
};
#endif

#endif
