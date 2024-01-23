/* basic types for dfslib, T10.090-T13.418; $DVS:time$ */

#ifndef DFSLIB_TYPES_H_INCLUDED
#define DFSLIB_TYPES_H_INCLUDED

#if defined(_WIN32) || defined(_WIN64)
#ifndef inline
#define inline __inline
#endif
#endif

enum dfslib_errors {
	DFSLIB_ERROR_MAX		= -2,
	DFSLIB_NOT_REALIZED		= -2,
	DFSLIB_IO_ERROR			= -3,
	DFSLIB_CORRUPT_DATA		= -4,
	DFSLIB_INVALID_VOLUME	= -5,
	DFSLIB_INVALID_TYPE		= -6,
	DFSLIB_INVALID_VALUE	= -7,
	DFSLIB_FILE_IS_DIR		= -8,
	DFSLIB_FILE_IS_NOT_DIR	= -9,
	DFSLIB_DIR_NOT_EMPTY	= -10,
	DFSLIB_INVALID_NAME		= -11,
	DFSLIB_NAME_TOO_LONG	= -12,
	DFSLIB_NOT_FOUND		= -13,
	DFSLIB_ACCESS_DENIED	= -14,
	DFSLIB_READ_ONLY		= -15,
	DFSLIB_NO_FREE_MEMORY	= -16,
	DFSLIB_NO_FREE_SPACE	= -17,
	DFSLIB_TEMP_UNAVAILABLE	= -18,
	DFSLIB_TOO_MANY_LOCKED	= -19,
	DFSLIB_INTERNAL_ERROR	= -20,
	DFSLIB_ERROR_MIN		= -20,
};

/* the structure of ino: 0 access bits + 4 nnode bits + nsector */
#define DFSLIB_INO_NNODE_SHIFT		0
#define DFSLIB_INO_NSECTOR_SHIFT	(DFSLIB_INO_NNODE_SHIFT + 4)
#define DFSLIB_SECT_PER_BLOCK_SHIFT	6
#define DFSLIB_SECT_PER_BLOCK		(1 << DFSLIB_SECT_PER_BLOCK_SHIFT)

typedef unsigned long long	dfs64;
typedef unsigned int		dfs32;
typedef unsigned short		dfs16;
typedef unsigned char		dfs8;

#endif
