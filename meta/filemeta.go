package meta

import "sort"

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

func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
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
