package cgo

/*
#include <stdlib.h>
#include "filestat.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type FileInfo struct {
	Size  uint64
	Mode  uint32
	Nlink uint32
	UID   uint32
	GID   uint32
}

func getFileStat(path string) (*FileInfo, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var cInfo C.struct_FileInfo

	if C.MyGetFileStat(cPath, &cInfo) != 0 {
		return nil, fmt.Errorf("failed to get file status for %s", path)
	}

	info := &FileInfo{
		Size:  uint64(cInfo.size),
		Mode:  uint32(cInfo.mode),
		Nlink: uint32(cInfo.nlink),
		UID:   uint32(cInfo.uid),
		GID:   uint32(cInfo.gid),
	}

	return info, nil
}

func RunCGo() {
	//path := "/Users/kealuya/mywork/my_git/renhao_go2/cgo/filestat.c" // 替换为你想检查的文件路径
	path := "./cgo/test.txt" // 替换为你想检查的文件路径
	info, err := getFileStat(path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("File size: %d bytes\n", info.Size)
	fmt.Printf("File mode: %o\n", info.Mode)
	fmt.Printf("Number of links: %d\n", info.Nlink)
	fmt.Printf("Owner UID: %d\n", info.UID)
	fmt.Printf("Group GID: %d\n", info.GID)
}
