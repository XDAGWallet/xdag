/* header of dfslib string functions, T10.259-T10.267; $DVS:time$ */

#ifndef DFSLIB_STRING_H_INCLUDED
#define DFSLIB_STRING_H_INCLUDED

#include "dfslib_types.h"
#include <stddef.h>

enum dfslib_string_types {
	DFSLIB_STRING_NONE,
	DFSLIB_STRING_UTF8,
	DFSLIB_STRING_UNICODE,
};

/* a structure describing the utf8 string or an array of unicode characters */
struct dfslib_string {
	int type;
	size_t len;
	union {
	    const char *utf8;
	    const dfs16 *unicode;
	};
};

#ifdef __cplusplus
extern "C" {
#endif

    extern int dfslib_unicode_to_utf8(dfs16 uni, char **pbuf, unsigned *psize);
    extern int dfslib_utf8_to_unicode(const char **utf, unsigned *plen);

    extern int dfslib_unicode_read(const struct dfslib_string *str, unsigned *ptr);
    extern int dfslib_unicode_cmp(const struct dfslib_string *str, unsigned *ptr, int unicode);
    extern int dfslib_unicode_strlen(const struct dfslib_string *str);
    extern int dfslib_string_to_unicode(struct dfslib_string *str, dfs16 *buf, unsigned len);
    extern int dfslib_string_to_utf8(struct dfslib_string *str, char *buf, unsigned len);
    extern int dfslib_substring(const struct dfslib_string *str, struct dfslib_string *substr, unsigned begin, unsigned end);
    extern int dfslib_unicode_strchr(const struct dfslib_string *str, int unicode);
    extern int dfslib_unicode_strtok(const struct dfslib_string *str, struct dfslib_string *token, const struct dfslib_string *limits, unsigned *ptr);

    static inline struct dfslib_string *dfslib_utf8_string(struct dfslib_string *str,
            const char *utf8, size_t len) {
        str->type = DFSLIB_STRING_UTF8;
        str->len = len;
        str->utf8 = utf8;
        return str;
    }

    static inline struct dfslib_string *dfslib_unicode_string(struct dfslib_string *str, const dfs16 *unicode, unsigned len) {
        str->type = DFSLIB_STRING_UNICODE;
        str->len = len;
        str->unicode = unicode;
        return str;
    }

#ifdef __cplusplus
};
#endif
#endif
