// filestat.c
#include "filestat.h"
#include <sys/types.h>
#include <sys/stat.h>
#include <unistd.h>



// 获取文件状态信息
int MyGetFileStat(const char *path, struct FileInfo *info) {
    struct stat st;
    if (stat(path, &st) != 0) {
        return -1;
    }

    info->size = st.st_size;
    info->mode = st.st_mode;
    info->nlink = st.st_nlink;
    info->uid = st.st_uid;
    info->gid = st.st_gid;

    return 0;
}
