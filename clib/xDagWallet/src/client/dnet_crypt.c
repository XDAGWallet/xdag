#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include "../dus/dfslib_random.h"
#include "../dus/crc.h"
#include "dnet_crypt.h"
#include "./utils/utils.h"

#define KEYFILE     "xdagj_dat" DELIMITER "dnet_key.dat"
#define PWDLEN      64
#define SECTOR_LOG  9
#define SECTOR_SIZE (1 << SECTOR_LOG)
#define KEYLEN_MIN	(DNET_KEYLEN / 4)

struct dnet_keys {
    struct dnet_key priv;
    struct dnet_key pub;
};

struct dnet_keys *g_dnet_keys;
static struct dfslib_crypt *g_dnet_user_crypt = 0;

struct dnet_session {
    struct dnet_key key;
    uint32_t sector_write[SECTOR_SIZE / 4];
    uint32_t sector_read[SECTOR_SIZE / 4];
    struct dfslib_crypt crypt_write;
    struct dfslib_crypt crypt_read;
    uint64_t pos_write, pos_read;
    struct dnet_host *host;
    const struct dnet_session_ops *ops;
    void *private_data;
    uint32_t route_ip;
    uint32_t route_port;
};

static int g_keylen = 0;

static int input_password(const char *prompt, char *buf, unsigned len) {
    return 0;
}

static int(*g_input_password)(const char *prompt, char *buf, unsigned size) = &input_password;

//static int(*g_input_password)(const char *prompt, char *buf, unsigned size) = NULL;


static void dnet_make_key(dfsrsa_t *key, int keylen) {
    unsigned i;
    for (i = keylen; i < DNET_KEYLEN; i += keylen)
        memcpy(key + i, key, keylen * sizeof(dfsrsa_t));
}

static int dnet_detect_keylen(dfsrsa_t *key, int keylen) {
    if (g_keylen && (key == g_dnet_keys->priv.key || key == g_dnet_keys->pub.key))
        return g_keylen;
    while (keylen >= 8) {
        if (memcmp(key, key + keylen / 2, keylen * sizeof(dfsrsa_t) / 2)) break;
        keylen /= 2;
    }
    return keylen;
}

static int dnet_rsa_crypt(dfsrsa_t *data, int datalen, dfsrsa_t *key, int keylen) {
    return dfsrsa_crypt(data, datalen, key, dnet_detect_keylen(key, keylen));
}

#define dfsrsa_crypt dnet_rsa_crypt

static void dnet_sector_to_password(uint32_t sector[SECTOR_SIZE / 4], char password[PWDLEN + 1]) {
    int i;
    for (i = 0; i < PWDLEN / 8; ++i) {
        unsigned crc = crc_of_array((unsigned char *)(sector + i * SECTOR_SIZE / 4 / (PWDLEN / 8)), SECTOR_SIZE / (PWDLEN / 8));
        sprintf(password + 8 * i, "%08X", crc);
    }
}

static void dnet_random_sector(uint32_t sector[SECTOR_SIZE / 4]) {
    char password[PWDLEN + 1] = "Iyf&%d#$jhPo_t|3fgd+hf(s@;)F5D7gli^kjtrd%.kflP(7*5gt;Y1sYRC4VGL&";
    int i, j;
    for (i = 0; i < 3; ++i) {
        struct dfslib_string str;
        dfslib_utf8_string(&str, password, PWDLEN);
        dfslib_random_sector(sector, 0, &str, &str);
        for (j = KEYLEN_MIN / 8; j <= SECTOR_SIZE / 4; j += KEYLEN_MIN / 8)
            sector[j - 1] &= 0x7FFFFFFF;
        if (i == 2) break;
        dfsrsa_crypt((dfsrsa_t *)sector, SECTOR_SIZE / sizeof(dfsrsa_t), g_dnet_keys->priv.key, DNET_KEYLEN);
        dnet_sector_to_password(sector, password);
    }
}

int dnet_generate_random_array(void *array, unsigned long size) {
    uint32_t sector[SECTOR_SIZE / 4];
    unsigned long i;
    if (size < 4 || size & (size - 1)) return -1;
    if (size >= 512) {
        for (i = 0; i < size; i += 512) dnet_random_sector((uint32_t *)((uint8_t *)array + i));
    } else {
        dnet_random_sector(sector);
        for (i = 0; i < size; i += 4) {
            *(uint32_t *)((uint8_t *)array + i) = crc_of_array((unsigned char *)sector + i * 512 / size, 512 / size);
        }
    }
    return 0;
}

static int dnet_test_keys(void) {
    uint32_t src[SECTOR_SIZE / 4], dest[SECTOR_SIZE / 4];
    dnet_random_sector(src);
    memcpy(dest, src, SECTOR_SIZE);
    if (dfsrsa_crypt((dfsrsa_t *)dest, SECTOR_SIZE / sizeof(dfsrsa_t), g_dnet_keys->priv.key, DNET_KEYLEN)) return 1;
    if (dfsrsa_crypt((dfsrsa_t *)dest, SECTOR_SIZE / sizeof(dfsrsa_t), g_dnet_keys->pub.key, DNET_KEYLEN)) return 2;
    if (memcmp(dest, src, SECTOR_SIZE)) return 3;
    memcpy(dest, src, SECTOR_SIZE);
    if (dfsrsa_crypt((dfsrsa_t *)dest, SECTOR_SIZE / sizeof(dfsrsa_t), g_dnet_keys->pub.key, DNET_KEYLEN)) return 4;
    if (dfsrsa_crypt((dfsrsa_t *)dest, SECTOR_SIZE / sizeof(dfsrsa_t), g_dnet_keys->priv.key, DNET_KEYLEN)) return 5;
    if (memcmp(dest, src, SECTOR_SIZE)) return 6;
    return 0;
}

static int set_user_crypt(struct dfslib_string *pwd) {
    uint32_t sector0[128];
    int i;
    g_dnet_user_crypt = malloc(sizeof(struct dfslib_crypt));
    if (!g_dnet_user_crypt) return -1;
    memset(g_dnet_user_crypt->pwd, 0, sizeof(g_dnet_user_crypt->pwd));
    dfslib_crypt_set_password(g_dnet_user_crypt, pwd);
    for (i = 0; i < 128; ++i) sector0[i] = 0x4ab29f51u + i * 0xc3807e6du;
    for (i = 0; i < 128; ++i) {
        dfslib_crypt_set_sector0(g_dnet_user_crypt, sector0);
        dfslib_encrypt_sector(g_dnet_user_crypt, sector0, 0x3e9c1d624a8b570full + i * 0x9d2e61fc538704abull);
    }
    return 0;
}

/* 
 * 1 - to encode data (data _ id, serial number data and size data measured in 32 - bit words)
 * 2 - decode
 * 3 - password and check it returns 0 on success
 * 4 - password and write print array data of length 16 bytes
 * 5 - to verify that the pattern in the data corresponds to the password
 * 6 - setup callback function to input password, data is pointer to function 
 *     int (*)(const char *prompt, char *buf, unsigned size);
 */
int dnet_user_crypt_action(unsigned *data, unsigned long long data_id, unsigned size, int action) {
    if (action != 4 && action != 6 && !g_dnet_user_crypt) return 0;
    switch (action) {
        case 1:
            dfslib_encrypt_array(g_dnet_user_crypt, data, size, data_id);
            break;
        case 2:
            dfslib_uncrypt_array(g_dnet_user_crypt, data, size, data_id);
            break;
        case 3:
            {
                struct dfslib_crypt *crypt = malloc(sizeof(struct dfslib_crypt));
                struct dfslib_string str;
                char pwd[256];
                int res;
                if (!crypt) return -1;
                memset(pwd, 0, 256);
                memset(&str, 0, sizeof(struct dfslib_string));
                (*g_input_password)("Password", pwd, 256);
                dfslib_utf8_string(&str, pwd, strlen(pwd));
                memset(crypt->pwd, 0, sizeof(crypt->pwd));
                crypt->ispwd = 0;
                dfslib_crypt_set_password(crypt, &str);
                res = (g_dnet_user_crypt->ispwd == crypt->ispwd
                        && !memcmp(g_dnet_user_crypt->pwd, crypt->pwd, sizeof(crypt->pwd)));
                free(crypt);
                return res ? 0 : -1;
            }
        case 4:
            {
                struct dfslib_crypt *crypt = malloc(sizeof(struct dfslib_crypt));
                struct dfslib_string str;
                char pwd[256];
                memset(pwd, 0, 256);
                memset(&str, 0, sizeof(struct dfslib_string));
                (*g_input_password)("Password", pwd, 256);
                dfslib_utf8_string(&str, pwd, strlen(pwd));
                memset(crypt->pwd, 0, sizeof(crypt->pwd));
                dfslib_crypt_set_password(crypt, &str);
                memcpy(data, crypt->pwd, sizeof(crypt->pwd));
                free(crypt);
                return 0;
            }
        case 5:
            return memcmp(g_dnet_user_crypt->pwd, data, sizeof(g_dnet_user_crypt->pwd)) ? -1 : 0;
        case 6:
            g_input_password = (int(*)(const char *, char *, unsigned))(void *)data;
            return 0;
        default: return -1;
    }
    return 0;
}

int dnet_crypt_init(const char *version) {
    FILE *f;
    struct dnet_keys *keys;
    int i;
    g_dnet_keys = malloc(sizeof(struct dnet_keys));
    if (!g_dnet_keys) return 1;
    keys = g_dnet_keys;
    dfslib_random_init();
    if (crc_init()) return 2;
    f = xdag_open_file(KEYFILE, "rb");
    if (f) {
        if (fread(keys, sizeof(struct dnet_keys), 1, f) != 1) {
            xdag_close_file(f);
            f = 0;
        } else {
            g_keylen = dnet_detect_keylen(keys->pub.key, DNET_KEYLEN);
            if (dnet_test_keys()) {
                struct dfslib_string str;
                char pwd[256];
                (*g_input_password)("Password", pwd, 256);
                dfslib_utf8_string(&str, pwd, strlen(pwd));
                set_user_crypt(&str);
                if (g_dnet_user_crypt) for (i = 0; i < (sizeof(struct dnet_keys) >> 9); ++i)
                    dfslib_uncrypt_sector(g_dnet_user_crypt, (uint32_t *)keys + 128 * i, ~(uint64_t)i);
                g_keylen = 0;
                g_keylen = dnet_detect_keylen(keys->pub.key, DNET_KEYLEN);
            }
        }
    }
    if (!f) {
        char buf[256];
        struct dfslib_string str;
        f = xdag_open_file(KEYFILE, "wb");
        if (!f) return 3;

        {
            struct dfslib_string str, str1;
            char pwd[256], pwd1[256];
            (*g_input_password)("Set password", pwd, 256);
            dfslib_utf8_string(&str, pwd, strlen(pwd));
            (*g_input_password)("Re-type password", pwd1, 256);
            dfslib_utf8_string(&str1, pwd1, strlen(pwd1));
            if (str.len != str1.len || memcmp(str.utf8, str1.utf8, str.len)) {
                //xdag_err(error_pwd_inconsistent, "Passwords differ.");
                //xdag_wrapper_event(event_id_promot, error_none,"Passwords differ.\n");
                return 4;
            }
            if (str.len) set_user_crypt(&str);
            (*g_input_password)("Type random keys", buf, 256);
        }

        dfslib_random_fill(keys->pub.key, DNET_KEYLEN * sizeof(dfsrsa_t), 0, dfslib_utf8_string(&str, buf, strlen(buf)));
        //xdag_wrapper_event(event_id_promot, error_none, "Generating host keys... ");
#ifdef __arm__
        g_keylen = KEYLEN_MIN;
#else
        g_keylen = DNET_KEYLEN;
#endif
        dfsrsa_keygen(keys->priv.key, keys->pub.key, g_keylen);
        dnet_make_key(keys->priv.key, g_keylen);
        dnet_make_key(keys->pub.key, g_keylen);
        //xdag_wrapper_event(event_id_promot, error_none,"OK.\n");
        if (g_dnet_user_crypt) for (i = 0; i < (sizeof(struct dnet_keys) >> 9); ++i)
            dfslib_encrypt_sector(g_dnet_user_crypt, (uint32_t *)keys + 128 * i, ~(uint64_t)i);
        if (fwrite(keys, sizeof(struct dnet_keys), 1, f) != 1) return 5;
        if (g_dnet_user_crypt) for (i = 0; i < (sizeof(struct dnet_keys) >> 9); ++i)
            dfslib_uncrypt_sector(g_dnet_user_crypt, (uint32_t *)keys + 128 * i, ~(uint64_t)i);
    }
    xdag_close_file(f);
    return -dnet_test_keys();
}
