package utils

import (
	"crypto/md5"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
	"testing"
)

func TestCompress(t *testing.T) {
	testRoot := ".syncup_test"
	fileRoot := testRoot + "/file"
	_ = os.MkdirAll(fileRoot, 0755)
	compressRoot := testRoot + "/compress"
	decompressRoot := testRoot + "/decompress"
	_ = os.MkdirAll(compressRoot, 0755)
	_ = os.MkdirAll(decompressRoot, 0755)
	fileList := createTestFile(fileRoot)
	log.Infof("%v", fileList)
	source := fileRoot
	target := compressRoot
	file, err := CompressToDir(source, target)
	if err != nil {
		panic(err)
	}
	err = Decompress(file, decompressRoot)
	if err != nil {
		panic(err)
	}
	result := compareDir(fileRoot, decompressRoot+"/file")
	log.Infof("compare result: %t", result)
	defer os.RemoveAll(testRoot)
}

func compareDir(dir1 string, dir2 string) bool {
	log.Infof("compare dir1: %s, dir2: %s", dir1, dir2)
	file1, _ := os.Open(dir1)
	file2, _ := os.Open(dir2)
	fileInfos1, _ := file1.Readdir(-1)
	fileInfos2, _ := file2.Readdir(-1)
	if len(fileInfos1) != len(fileInfos2) {
		return false
	}
	for i := 0; i < len(fileInfos1); i++ {
		if fileInfos1[i].IsDir() != fileInfos2[i].IsDir() {
			return false
		}

		if fileInfos1[i].Name() != fileInfos2[i].Name() {
			return false
		}
		if fileInfos1[i].IsDir() {
			result := compareDir(dir1+"/"+fileInfos1[i].Name(), dir2+"/"+fileInfos2[i].Name())
			if !result {
				return false
			}
		} else {
			log.Infof("compare file1: %s, file2: %s", dir1+"/"+fileInfos1[i].Name(), dir2+"/"+fileInfos2[i].Name())

			if fileInfos1[i].Size() != fileInfos2[i].Size() {
				return false
			}
			cmpFile1 := dir1 + "/" + fileInfos1[i].Name()
			cmpFile2 := dir2 + "/" + fileInfos2[i].Name()
			fileMD51 := getFileMD5(cmpFile1)
			fileMD52 := getFileMD5(cmpFile2)
			if !strings.EqualFold(fileMD51, fileMD52) {
				log.Infof("md5 diff error: %s-%s, %s-%s", cmpFile1, fileMD51, cmpFile2, fileMD52)
				return false
			}
		}
	}
	return true
}

func createTestFile(fileRoot string) []string {
	_ = os.MkdirAll(fileRoot+"/sub1", 0755)
	testFile1, _ := os.Create(fileRoot + "/sub1/test1.txt")
	_, _ = testFile1.WriteString("test1 for compress/decompress")
	defer testFile1.Close()
	testFile2, _ := os.Create(fileRoot + "/test2.txt")
	_, _ = testFile2.WriteString("test2 for compress/decompress")
	defer testFile2.Close()
	fileList := make([]string, 0)
	fileList = append(fileList, testFile1.Name())
	fileList = append(fileList, testFile2.Name())
	return fileList
}

func getFileMD5(filepath string) string {
	file, _ := os.Open(filepath)
	defer file.Close()
	md5hash := md5.New()
	_, _ = io.Copy(md5hash, file)
	hash := md5hash.Sum(nil)
	return fmt.Sprintf("%x", hash)
}
