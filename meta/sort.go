package meta

import "time"

const baseTime = "2006-01-02 15:04:05"

type ByUploadTime []FileMeta

func (b ByUploadTime) Len() int {
	return len(b)
}

func (b ByUploadTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByUploadTime) Less(i, j int) bool {
	iTime, _ := time.Parse(baseTime, b[i].UploadTime)
	jTIme, _ := time.Parse(baseTime, b[j].UploadTime)
	return iTime.UnixNano() > jTIme.UnixNano()
}
