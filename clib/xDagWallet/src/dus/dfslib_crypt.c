/* cryptography engine for dfslib, T4.846-T10.734; $DVS:time$ */

#ifdef __LDuS__
#include <string.h>
#endif

#include "dfslib_crypt.h"

#define DFS_MAGIC0	 572035291u
#define DFS_MAGIC1	2626708081u
#define DFS_MAGIC2	2471573851u
#define DFS_MAGIC3	3569250857u
#define DFS_MAGIC4	1971772241u
#define DFS_MAGIC5	1615037507u
#define DFS_MAGIC6	  43385317u
#define DFS_MAGIC7	1229426917u
#define DFS_MAGIC8	3433359571u


int dfslib_crypt_set_password(struct dfslib_crypt *dfsc,
			      const struct dfslib_string *password) {
	unsigned ptr = 0;
	dfsc->ispwd = 0;
	dfsc->pwd[0] = DFS_MAGIC0;
	dfsc->pwd[1] = DFS_MAGIC1;
	dfsc->pwd[2] = DFS_MAGIC2;
	dfsc->pwd[3] = DFS_MAGIC3;
	if (password) {
	    dfs64 res;
	    int i, err;
	    while ((err = dfslib_unicode_read(password, &ptr)) >= 0) {
		for (i = 0, res = err; i < 4; ++i, res >>= 32)
		    dfsc->pwd[i] = (dfs32)(res += (dfs64)dfsc->pwd[i] * DFS_MAGIC4);
	    }
	    if (err != DFSLIB_NAME_TOO_LONG) return err;
	}
	return dfsc->ispwd = !!ptr;
}

int dfslib_crypt_copy_password(struct dfslib_crypt *to,
		const struct dfslib_crypt *from) {
	int i;
	for (i = 0; i < DFSLIB_CRYPT_PWDLEN; ++i)
	    to->pwd[i] = from->pwd[i];
	return to->ispwd = from->ispwd;
}

int dfslib_crypt_is_password(struct dfslib_crypt *dfsc) {
	return dfsc->ispwd;
}

#define dfs_enmix2(c,d) c = *--d ^= c * DFS_MAGIC6
#define dfs_enmix4(c,d) dfs_enmix2(c,d),dfs_enmix2(c,d),dfs_enmix2(c,d),dfs_enmix2(c,d)
#define dfs_enmix6(c,d) dfs_enmix4(c,d),dfs_enmix4(c,d),dfs_enmix4(c,d),dfs_enmix4(c,d)

static inline void dfs_enmixSector(register dfs32 *d) {
	int i = 8; register dfs32 c = DFS_MAGIC5; d += 128;
	while (i--) dfs_enmix6(c, d);
}

static inline void dfs_enmixArray(register dfs32 *d, int size) {
	register dfs32 c = DFS_MAGIC5; d += size;
	while (size--) dfs_enmix2(c, d);
}

#define dfs_unmix2(c,d) c *= DFS_MAGIC6, c ^= *--d ^= c
#define dfs_unmix4(c,d) dfs_unmix2(c,d),dfs_unmix2(c,d),dfs_unmix2(c,d),dfs_unmix2(c,d)
#define dfs_unmix6(c,d) dfs_unmix4(c,d),dfs_unmix4(c,d),dfs_unmix4(c,d),dfs_unmix4(c,d)

static inline void dfs_unmixSector(register dfs32 *d) {
	int i = 8; register dfs32 c = DFS_MAGIC5;
	while (i--) dfs_unmix6(c, d);
}

static inline void dfs_unmixArray(register dfs32 *d, int size) {
	register dfs32 c = DFS_MAGIC5;
	while (size--) dfs_unmix2(c, d);
}

#define dfs_crypt0(a,x,y,z,t) a = ((dfs32)(((dfs64)y * \
		(z + dfsc->regs[x >> 16])) >> 16) ^ dfsc->regs[(dfs16)t])

#define dfs_crypt2(a,b,c,d,x,y,z,t) \
		dfs_crypt0(a,x,y,z,t), dfs_crypt0(b,y,z,t,x),\
		dfs_crypt0(c,z,t,x,y), dfs_crypt0(d,t,x,y,z)

#define dfs_crypt3(a,b,c,d,x,y,z,t) \
		dfs_crypt2(a,b,c,d,x,y,z,t), dfs_crypt2(x,y,z,t,a,b,c,d)

#define dfs_crypt6(a,b,c,d,x,y,z,t) \
		dfs_crypt3(a,b,c,d,x,y,z,t), dfs_crypt3(a,b,c,d,x,y,z,t), \
		dfs_crypt3(a,b,c,d,x,y,z,t), dfs_crypt3(a,b,c,d,x,y,z,t), \
		dfs_crypt3(a,b,c,d,x,y,z,t), dfs_crypt3(a,b,c,d,x,y,z,t), \
		dfs_crypt3(a,b,c,d,x,y,z,t), dfs_crypt3(a,b,c,d,x,y,z,t)

static inline dfs16 dfs_mod(dfs64 big_, dfs16 small_) {
	if ((unsigned long)big_ == big_) {
		return (unsigned long)big_ % small_;
	} else {
	    dfs32 tmp = (dfs32)(big_ >> 32) % small_;
		tmp <<= 16;
		tmp |= (dfs32)big_ >> 16;
		tmp %= small_;
		tmp <<= 16;
		tmp |= (dfs16)big_; return tmp % small_;
	}
}

static void dfs_prepare(struct dfslib_crypt *dfsc, dfs64 sectorNo,
			dfs32 *px, dfs32 *py, dfs32 *pz, dfs32 *pt){
	dfs32 a,b,c,d,x,y,z,t;
	sectorNo *= (dfs64)DFS_MAGIC7 << 32 | DFS_MAGIC8;
	x = dfsc->pwd[0] ^ dfsc->regs[dfs_mod(sectorNo, 65479) + 31];
	y = dfsc->pwd[1] ^ dfsc->regs[dfs_mod(sectorNo, 65497) + 11];
	z = dfsc->pwd[2] ^ dfsc->regs[dfs_mod(sectorNo, 65519) +  5];
	t = dfsc->pwd[3] ^ dfsc->regs[dfs_mod(sectorNo, 65521) +  3];
	dfs_crypt6(a,b,c,d,x,y,z,t);
	*px = x;
	*py = y;
	*pz = z;
	*pt = t;
}

#define dfs_encrypt2(a,b,c,d,x,y,z,t,data) \
		dfs_crypt0(a,x,y,z,t)^~*data,	dfs_crypt0(b,y,z,t,x),\
		dfs_crypt0(c,z,t,x,y), *data++-=dfs_crypt0(d,t,x,y,z)

#define dfs_encrypt3(a,b,c,d,x,y,z,t,data) \
		dfs_encrypt2(a,b,c,d,x,y,z,t,data), dfs_encrypt2(x,y,z,t,a,b,c,d,data)

#define dfs_encrypt6(a,b,c,d,x,y,z,t,data) \
		dfs_encrypt3(a,b,c,d,x,y,z,t,data), dfs_encrypt3(a,b,c,d,x,y,z,t,data), \
		dfs_encrypt3(a,b,c,d,x,y,z,t,data), dfs_encrypt3(a,b,c,d,x,y,z,t,data), \
		dfs_encrypt3(a,b,c,d,x,y,z,t,data), dfs_encrypt3(a,b,c,d,x,y,z,t,data), \
		dfs_encrypt3(a,b,c,d,x,y,z,t,data), dfs_encrypt3(a,b,c,d,x,y,z,t,data)

int dfslib_encrypt_sector(struct dfslib_crypt *dfsc, dfs32 *sector, dfs64 sectorNo) {
	dfs32 a,b,c,d,x,y,z,t; int i = 8;
	if (!dfsc->ispwd) return 0;
	dfs_prepare(dfsc, sectorNo, &x, &y, &z, &t);
	dfs_enmixSector(sector);
	while(i--) dfs_encrypt6(a,b,c,d,x,y,z,t,sector);
	return 1;
}

int dfslib_encrypt_array(struct dfslib_crypt *dfsc, dfs32 *data, unsigned size, dfs64 sectorNo) {
	dfs32 a,b,c,d,x,y,z,t;
	if (!dfsc->ispwd || (size & 1)) return 0;
	dfs_prepare(dfsc, sectorNo, &x, &y, &z, &t);
	dfs_enmixArray(data, size); size >>= 1;
	while (size--) dfs_encrypt3(a,b,c,d,x,y,z,t,data);
	return 1;
}

#define dfs_uncrypt2(a,b,c,d,x,y,z,t,data) \
		dfs_crypt0(c,z,t,x,y), *data+= dfs_crypt0(d,t,x,y,z),\
		dfs_crypt0(a,x,y,z,t)^~*data++,dfs_crypt0(b,y,z,t,x)

#define dfs_uncrypt3(a,b,c,d,x,y,z,t,data) \
		dfs_uncrypt2(a,b,c,d,x,y,z,t,data), dfs_uncrypt2(x,y,z,t,a,b,c,d,data)

#define dfs_uncrypt6(a,b,c,d,x,y,z,t,data) \
		dfs_uncrypt3(a,b,c,d,x,y,z,t,data), dfs_uncrypt3(a,b,c,d,x,y,z,t,data), \
		dfs_uncrypt3(a,b,c,d,x,y,z,t,data), dfs_uncrypt3(a,b,c,d,x,y,z,t,data), \
		dfs_uncrypt3(a,b,c,d,x,y,z,t,data), dfs_uncrypt3(a,b,c,d,x,y,z,t,data), \
		dfs_uncrypt3(a,b,c,d,x,y,z,t,data), dfs_uncrypt3(a,b,c,d,x,y,z,t,data)

int dfslib_uncrypt_sector(struct dfslib_crypt *dfsc, dfs32 *sector, dfs64 sectorNo) {
	dfs32 a,b,c,d,x,y,z,t; int i = 8;
	if (!dfsc->ispwd) return 0;
	dfs_prepare(dfsc, sectorNo, &x, &y, &z, &t);
	while (i--) dfs_uncrypt6(a,b,c,d,x,y,z,t,sector);
	dfs_unmixSector(sector);
	return 1;
}

int dfslib_uncrypt_array(struct dfslib_crypt *dfsc, dfs32 *data, unsigned size,
			 dfs64 sectorNo) {
	dfs32 a,b,c,d,x,y,z,t; int size0 = size;
	if (!dfsc->ispwd || (size & 1)) return 0;
	dfs_prepare(dfsc, sectorNo, &x, &y, &z, &t);
	size >>= 1;
	while (size--) dfs_uncrypt3(a,b,c,d,x,y,z,t,data);
	dfs_unmixArray(data, size0);
	return 1;
}

static inline void dfs_memcpy(dfs8 *to, const dfs8 *from, unsigned size) {
#ifdef __LDuS__
	memcpy(to, from, size);
#else
	while (size--) *to++ = *from++;
#endif
}

int dfslib_crypt_set_sector0(struct dfslib_crypt *dfsc, const void *sector) {
	const dfs8 *data = (dfs8*)sector;
	int i;
	if (!dfsc->ispwd) return 0;
	for (i = 0; i < 512; ++i) {
	    dfs_memcpy((dfs8*)dfsc->regs + (i << 9),           data + i, 512 - i);
	    dfs_memcpy((dfs8*)dfsc->regs + ((i + 1) << 9) - i, data    , i      );
	}
	for (i = 0; i < 512; ++i)
	    dfslib_encrypt_sector(dfsc, dfsc->regs + (i << 7), i);
	return 1;
}
