//
//  utils.c
//  xdag
//
//  Copyright © 2018 xdag contributors.
//

#include "utils.h"

#if defined (__MACOS__) || defined (__APPLE__)
#include <libgen.h>
#define PATH_MAX 4096
#include <sys/time.h>
#include <unistd.h>
#endif

FILE* xdag_open_file(const char *path, const char *mode)
{
    char abspath[1024] = {0};
    // sprintf(abspath, "%s%s", g_xdag_current_path, path);
    sprintf(abspath, "%s", path);
    FILE* f = fopen(abspath, mode);
    return f;
}

void xdag_close_file(FILE *f)
{
    fclose(f);
}
