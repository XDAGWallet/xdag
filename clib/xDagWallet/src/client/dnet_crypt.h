#ifndef DNET_CRYPT_H_INCLUDED
#define DNET_CRYPT_H_INCLUDED

#include <sys/types.h>
#include "../dus/dfsrsa.h"

#define DNET_VERSION "T11.231-T13.714" /* $DVS:time$ */
#define DNET_KEY_SIZE	4096
#define DNET_KEYLEN	((DNET_KEY_SIZE * 2) / (sizeof(dfsrsa_t) * 8))

struct dnet_key {
    dfsrsa_t key[DNET_KEYLEN];
};

#ifdef __cplusplus
extern "C" {
#endif
    extern int dnet_crypt_init(const char *version);

    extern int dnet_generate_random_array(void *array, unsigned long size);

    /* выполнить действие с паролем пользователя:
     * 1 - закодировать данные (data_id - порядковый номер данных, size - размер данных, измеряется в 32-битных словах)
     * 2 - декодировать -//-
     * 3 - ввести пароль и проверить его, возвращает 0 при успехе
     * 4 - ввести пароль и записать его отпечаток в массив data длины 16 байт
     * 5 - проверить, что отпечаток в массиве data соответствует паролю
     * 6 - setup callback function to input password, data is pointer to function
     *     int (*)(const char *prompt, char *buf, unsigned size);
     */
    extern int dnet_user_crypt_action(unsigned *data, unsigned long long data_id, unsigned size, int action);
#ifdef __cplusplus
};
#endif

#endif


