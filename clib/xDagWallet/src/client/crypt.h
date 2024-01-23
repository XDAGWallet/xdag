/* cryptography, T13.654-T13.826 $DVS:time$ */

#ifndef XDAG_CRYPT_H
#define XDAG_CRYPT_H

#include <stdint.h>
typedef uint64_t xdag_hash_t[4];
#ifdef __cplusplus
extern "C" {
#endif

    // initialization of the encryption system
    extern int xdag_crypt_init(int withrandom);

    // creates a new pair of private and public keys
    extern void *xdag_create_key(xdag_hash_t privkey, xdag_hash_t pubkey, uint8_t *pubkey_bit);

    // returns the internal representation of the key and the public key by the known private key
    extern void *xdag_private_to_key(const xdag_hash_t privkey, xdag_hash_t pubkey, uint8_t *pubkey_bit);

    // returns the internal representation of the key by the known public key
    extern void *xdag_public_to_key(const xdag_hash_t pubkey, uint8_t pubkey_bit);

    // removes the internal key representation
    extern void xdag_free_key(void *key);

    // verify that the signature (sign_r, sign_s) corresponds to a hash 'hash', a version for its own key
    extern int xdag_verify_signature(const void *key, const xdag_hash_t hash, const xdag_hash_t sign_r, const xdag_hash_t sign_s);

#ifdef __cplusplus
};
#endif

#endif
