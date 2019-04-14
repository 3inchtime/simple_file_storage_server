package meta

import (
	"fmt"
	"simple_file_storage_server/dbops"
	"sort"
)

type FileMeta struct {
	FileSha1   string
	FileName   string
	FileSize   int64
	Location   string
	UploadTime string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

//保存文件信息至数据库
func FileMetaUploadDB(fmeta FileMeta) bool {
	filehash := fmeta.FileSha1
	filename := fmeta.FileName
	filesize := fmeta.FileSize
	fileaddr := fmeta.Location
	return dbops.UploadFileMetaDB(filehash, filename, filesize, fileaddr)
}

func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

//从数据库获取文件信息
func GetFileMetaDB(fileSha1 string) (FileMeta, error) {
	fresult, err := dbops.GetFileMetaDB(fileSha1)
	if err != nil {
		fmt.Printf("Query file err: %s\n", err.Error())
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		FileSha1: fresult.FileHash,
		FileName: fresult.FileName.String,
		FileSize: fresult.FileSize.Int64,
		Location: fresult.FileAddr.String,
	}
	return fmeta, nil

}

func GetLastFileMetas(limit int) []FileMeta {
	fileMetaArray := make([]FileMeta, limit)

	for _, file := range fileMetas {
		fileMetaArray = append(fileMetaArray, file)
	}
	// 文件按创建时间排序
	sort.Sort(ByUploadTime(fileMetaArray))

	return fileMetaArray[:limit]
}

func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
