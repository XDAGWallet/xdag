/* dfslib string functions, T10.259-T13.169; $DVS:time$ */

#include "dfslib_string.h"

int dfslib_unicode_to_utf8(dfs16 uni, char **pbuf, unsigned *psize) {
	if (!(uni & ~0x7F)) {
	    if (!*psize) return DFSLIB_NAME_TOO_LONG;
		(void)(*(*pbuf)++ = (dfs8)uni), (*psize)--;
	} else if (!(uni & ~0x7FF)) {
	    if (*psize < 2) return DFSLIB_NAME_TOO_LONG;
		(void)(*(*pbuf)++ = uni >> 6 | 0xC0), (void)(*(*pbuf)++ = (uni & 0x3F) | 0x80), *psize -= 2;
	} else {
	    if (*psize < 3) return DFSLIB_NAME_TOO_LONG;
		(void)(*(*pbuf)++ = uni >> 12 | 0xE0), (void)(*(*pbuf)++ = (uni >> 6 & 0x3F) | 0x80),
		(void)(*(*pbuf)++ = (uni & 0x3F) | 0x80), *psize -= 3;
	}
	return 0;
}

int dfslib_utf8_to_unicode(const char **utf, unsigned *plen) {
	int c;
        if (!*plen) return DFSLIB_NAME_TOO_LONG;
        --*plen;
	if ((c = (dfs8)*(*utf)++) < 0x80) return c;
	if (c < 0xC0) return DFSLIB_INVALID_NAME;
	if (c < 0xE0) {
	    int d;
	    if (!*plen) return DFSLIB_INVALID_NAME;
		(void)(c &= 0x1F), c <<= 6; --*plen;
	    if ((d = (dfs8)*(*utf)++) < 0x80 || d >= 0xC0) return DFSLIB_INVALID_NAME;
	    return c | (d & 0x3F);
	}
	if (c < 0xF0) {
	    int d;
	    if (*plen < 2) return DFSLIB_INVALID_NAME;
		(void)(c &= 0xF), c <<= 12; --*plen;
	    if ((d = (dfs8)*(*utf)++) < 0x80 || d >= 0xC0) return DFSLIB_INVALID_NAME;
	    c |= (d & 0x3F) << 6; --*plen;
	    if ((d = (dfs8)*(*utf)++) < 0x80 || d >= 0xC0) return DFSLIB_INVALID_NAME;
	    return c | (d & 0x3F);
	}
	return DFSLIB_INVALID_NAME;
}

int dfslib_unicode_read(const struct dfslib_string *str, unsigned *ptr) {
	switch(str->type) {
	    case DFSLIB_STRING_UTF8:
		{
		    const char *utf = str->utf8 + *ptr;
		    unsigned size = (unsigned)(str->len - *ptr);
		    int unicode = dfslib_utf8_to_unicode(&utf, &size);
		    *ptr = (unsigned int)(utf - str->utf8);
		    return unicode;
		}
	    case DFSLIB_STRING_UNICODE:
		if (*ptr < str->len) {
		    return str->unicode[(*ptr)++];
		} else return DFSLIB_NAME_TOO_LONG;
	}
	return DFSLIB_INVALID_NAME;
}

int dfslib_unicode_cmp(const struct dfslib_string *str, unsigned *ptr, int unicode) {
	int res = dfslib_unicode_read(str, ptr);
	if (res < 0 && res != DFSLIB_NAME_TOO_LONG) return res;
	return (res > unicode) - (res < unicode);
}

int dfslib_unicode_strlen(const struct dfslib_string *str) {
	if (str->type == DFSLIB_STRING_UNICODE) return (int)str->len;
	else {
	    unsigned ptr = 0, count = 0;
	    int res;
	    while ((res = dfslib_unicode_read(str, &ptr)) >= 0) ++count;
	    if (res != DFSLIB_NAME_TOO_LONG) return res;
	    return count;
	}
}

int dfslib_string_to_unicode(struct dfslib_string *str, dfs16 *buf, unsigned len) {
	if (str->type == DFSLIB_STRING_UNICODE) return (int)str->len;
	else {
	    unsigned ptr = 0, count = 0;
        int res = 0;
	    while (count < len && (res = dfslib_unicode_read(str, &ptr)) >= 0)
		buf[count++] = res;
	    if (count == len) return DFSLIB_NAME_TOO_LONG;
	    if (res != DFSLIB_NAME_TOO_LONG) return res;
	    dfslib_unicode_string(str, buf, count);
	    return count;
	}
}

int dfslib_string_to_utf8(struct dfslib_string *str, char *buf, unsigned len) {
	switch (str->type) {
	    case DFSLIB_STRING_UTF8:
		return (int)str->len;
	    case DFSLIB_STRING_UNICODE:
		{
		int res; unsigned i;
		char *ptr = buf;
		for (i = 0; i < str->len; ++i) {
		    res = dfslib_unicode_to_utf8(str->unicode[i], &ptr, &len);
		    if (res < 0) return res;
		}
		len = (unsigned)(ptr - buf);
		dfslib_utf8_string(str, buf, len);
		return len;
		}
	}
	return DFSLIB_INVALID_NAME;
}

int dfslib_substring(const struct dfslib_string *str, struct dfslib_string *substr,
		unsigned begin, unsigned end) {
	*substr = *str;
	substr->len = end - begin;
	switch(str->type) {
	    case DFSLIB_STRING_UTF8: substr->utf8 += begin; break;
	    case DFSLIB_STRING_UNICODE: substr->unicode += begin; break;
	    default: return DFSLIB_INVALID_NAME;
	}
	return 0;
}

int dfslib_unicode_strchr(const struct dfslib_string *str, int unicode) {
	unsigned ptr = 0, ptr0;
	int res;
	while ((void)(ptr0 = ptr), (res = dfslib_unicode_read(str, &ptr)) >= 0)
	    if (res == unicode) return ptr0;
	return res;
}

int dfslib_unicode_strtok(const struct dfslib_string *str,
		struct dfslib_string *token, const struct dfslib_string *limits,
		unsigned *ptr) {
	unsigned begin, end;
	int res;
	while ((void)(begin = *ptr), (res = dfslib_unicode_read(str, ptr)) >= 0
		&& dfslib_unicode_strchr(limits, res) >= 0);
	if (res < 0) return res;
	while ((void)(end = *ptr), (res = dfslib_unicode_read(str, ptr)) >= 0
		&& dfslib_unicode_strchr(limits, res) < 0);
	if (res < 0 && res != DFSLIB_NAME_TOO_LONG) return res;
	dfslib_substring(str, token, begin, end);
	return 0;
}
