/* dfslib library for random data generation, T10.273-T10.717; $DVS:time$ */

#include <stdlib.h>
#include <time.h>
#if !defined(_WIN32) && !defined(_WIN64)
#include <unistd.h>
#include <sys/time.h>
#include <sys/times.h>
#define USE_RAND48
#else
#include <Windows.h>
#include <process.h>
#define getpid _getpid
#endif
#include "dfslib_random.h"

unsigned dfslib_random_get(unsigned limit) {
	unsigned res;
#ifdef USE_RAND48
	res = (unsigned)mrand48();
#else
	res = rand();
#endif
	if (limit) res %= limit;
	return res;
}

void dfslib_random_fill(void *buf, unsigned long len, int xor, struct dfslib_string *tip) {
	unsigned res = 0, bytes3 = 0, rnd, tmp = 0;
	while (len) {
	    rnd = dfslib_random_get(0);
	    if (tip) {
		int uni = dfslib_unicode_read(tip, &tmp);
		if (uni < 0) {
		    tmp = 0;
		    uni = dfslib_unicode_read(tip, &tmp);
		}
		rnd += uni;
	    }
		(void)(res *= 41), (void)(res += rnd % 41), bytes3 += 2;
	    if (bytes3 >= 10) {
		if (xor) *(unsigned char *)buf ^= (unsigned char)res;
		else *(unsigned char *)buf = (unsigned char)res;
			(void)(res >>= 8), bytes3 -= 3;
			(void)(buf = (unsigned char *)buf + 1), --len;
	    }
	}
}

void dfslib_random_sector(dfs32 *sector, struct dfslib_crypt *crypt0,
		struct dfslib_string *password, struct dfslib_string *tip) {
	struct dfslib_crypt crypt[1];
	struct dfslib_string tip0[1];
	char tips[6 * DFSLIB_CRYPT_PWDLEN];
	dfs64 nsector;
	int i;
	if (crypt0) dfslib_crypt_copy_password(crypt, crypt0);
	else dfslib_crypt_set_password(crypt, password);
	if (!tip) {
	    for (i = 0; i < DFSLIB_CRYPT_PWDLEN; ++i) {
		dfs32 res = crypt->pwd[i];
		int j;
		for (j = 0; j < 6; ++j)
			(void)(tips[i * 6 + j] = (res % 41) + 0x21), res /= 41;
	    }
	    dfslib_utf8_string(tip0, tips, 6 * DFSLIB_CRYPT_PWDLEN);
	    tip = tip0;
	}
	for (i = 0; i < 3; ++i) {
	    dfslib_random_fill(sector, 512, i, tip);
	    dfslib_random_fill(&nsector, 8, i, tip);
	    dfslib_crypt_set_sector0(crypt, sector);
	    dfslib_encrypt_sector(crypt, sector, nsector);
	}
}

void dfslib_random_init(void) {
	dfs64 seed = 0, time1, time2, clock, pid;
#if !defined(_WIN32) && !defined(_WIN64)
	struct timeval tv[1];
	struct tms tms[1];
	gettimeofday(tv, NULL);
	time1 = tv->tv_sec;
	time2 = tv->tv_usec;
	clock = times(tms);
#else
	FILETIME ft[1];
	GetSystemTimeAsFileTime(ft);
	time1 = ft->dwHighDateTime;
	time2 = ft->dwLowDateTime;
	clock = GetTickCount();
#endif
	pid = getpid();
	seed ^= time1;	seed *= 0x8E230615u;
	seed ^= time2;	seed *= 0x40D95A7Bu;
	seed ^= clock;	seed *= 0x0EE493B1u;
	seed ^= pid;	seed *= 0xB8204941u;
	dfslib_random_fill(&seed, 8, 1, 0);
#if 0
	printf("dfslib_random: time=(%lX, %lX), times=%lX, pid=%lX, seed=%llX\n",
	    (long)tv->tv_sec, (long)tv->tv_usec, (long)clock, (long)pid, seed);
#endif
#ifdef USE_RAND48
	{
	    unsigned short xsubi[3];
	    int i;
	    for (i = 0; i < 3; ++i)
			(void)(xsubi[i] = seed & 0xFFFF), seed >>= 16;
	    seed48(xsubi);
	}
#else
	srand((unsigned)seed);
#endif
}
