// filestat.h

#ifndef FILESTAT_H
#define FILESTAT_H

#include <sys/types.h>
#include <sys/stat.h>
#include <unistd.h>

struct FileInfo {
    unsigned long size;
    unsigned int mode;
    unsigned int nlink;
    unsigned int uid;
    unsigned int gid;
};

int MyGetFileStat(const char *path, struct FileInfo *info);

#endif // FILESTAT_H
