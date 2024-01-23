#ifndef _DAR_CRC_H_INCLUDED
#define _DAR_CRC_H_INCLUDED

/* CRC library h-file, T4.046-T9.267; $DVS:time$ */

#include <stdio.h>


#ifdef __cplusplus
extern "C" {
#endif
    /* initialization of internal CRC table (with memory allocation) */
    extern int crc_init(void);

    /* constructing a table in an external array of 256 double words */
	extern int crc_makeTable(unsigned table[256]);

    /* add to the accumulated CRC new data contained in the array buf
     len lengths; returns a new CRC; CRC initial value = 0 */
    extern unsigned crc_addArray(unsigned char *buf, unsigned len, unsigned crc);

    /* add to the accumulated CRC new data contained in file f, but not
     more len bytes; returns a new CRC; CRC initial value = 0 */
	extern unsigned crc_addFile	(FILE *f, unsigned len, unsigned crc);
#ifdef __cplusplus
};
#endif

/* calculate the CRC of the array */
#define crc_of_array(buf,len)    crc_addArray(buf,len,0)

/* CRC file count */
#define crc_of_file(f)        crc_addFile(f,-1,0)

#endif
