/* Библиотека для вычисления CRC32-кода данной последовательности байт;
   исходный текст взят со страницы CRC из wikiпедии (www.wikipedia.ru)
*/

/*
  Name  : CRC-32
  Poly  : 0x04C11DB7	x^32 + x^26 + x^23 + x^22 + x^16 + x^12 + x^11 +
			x^10 + x^8  + x^7  + x^5  + x^4  + x^2  + x    + 1
  Init  : 0xFFFFFFFF
  Revert: true
  XorOut: 0xFFFFFFFF
  Check : 0xCBF43926 ("123456789")
  MaxLen: 268 435 455 байт (2 147 483 647 бит) - обнаружение
   одинарных, двойных, пакетных и всех нечетных ошибок
*/

#include <stdlib.h> /* T5.245 */
#include <stdio.h>
#ifdef __DuS__
#include <dus/error.h>
#else
#define _errn(n) return ((n) << 1 | 1)
#include "crc.h"
#endif

//static const char version[] = "CRC library, ...-T4.046-T11.609"; /* $DVS:time$ */

unsigned *crc_table = NULL;

int crc_makeTable(unsigned table[256])
{
	unsigned long crc; int i, j;
	for(i = 0; i < 256; i++) {
		crc = i;
		for(j = 0; j < 8; j++)
			crc = (crc >> 1) ^ (crc & 1 ? 0xEDB88320UL : 0);
		table[i] = (unsigned)crc;
	}
	return 0;
}

int crc_init(void)
{
	if(crc_table) return 0;
	crc_table = malloc(256 * sizeof(unsigned));
	if(!crc_table) _errn(0);
	crc_makeTable(crc_table);
	return 0;
}

#define crc_addChar(crc,c) ((crc)=crc_table[(unsigned char)(crc)^(c)]^((crc)>>8))

unsigned crc_addArray(unsigned char *buf, unsigned len, unsigned crc)
{
	crc = ~crc;
	while(len--) crc_addChar(crc, *buf++);
	return ~crc;
}

unsigned crc_addFile(FILE *f, unsigned len, unsigned crc)
{
	int c;
	crc = ~crc;
	while(len-- && (c = fgetc(f)) != EOF) crc_addChar(crc, c);
	return ~crc;
}
