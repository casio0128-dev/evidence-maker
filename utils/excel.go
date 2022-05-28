package utils

import (
	"evidence-maker/conf"
	"evidence-maker/consts"
	"fmt"
	"github.com/xuri/excelize/v2"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"
	"sync"
	"time"
)

/*
	エクセルファイルのポインタを取得する。
	・指定したパスにエクセルファイルがあれば、そのポインタを返却する
	・指定したパスにエクセルファイルがなければ、新規にファイルポインタを作成し返却する。
*/
func openExcel(filePath string) (evd *excelize.File, err error) {
	if isExist(filePath) {
		if evd, err = excelize.OpenFile(filePath); err != nil {
			return
		}
	} else {
		evd = excelize.NewFile()
	}
	return
}

func OutputExcelFile(wg *sync.WaitGroup, cf *conf.Config) error {
	dirPath, err := os.MkdirTemp(consts.OutputDirectory, time.Now().Format(consts.OutputDirectoryPattern))
	if err != nil {
		return err
	}
	files, err := getExcelFileNames()
	if err != nil {
		return err
	}

	for _, bn := range files {
		bn := bn
		wg.Add(1)
		go func(bookName string) {
			defer func() {
				// メソッド終了時またはパニック発生時に、waitGroupを終了にしてサブゴルーチンの終了を記録
				switch recover() {
				default:
					wg.Done()
				}
			}()

			name := fmt.Sprintf(consts.OutputExcelPattern, strings.Join([]string{dirPath, bookName}, string(os.PathSeparator)))
			book, err := openExcel(cf.Template.FilePath)
			if err != nil {
				panic(err)
			}
			defer func(f *excelize.File) {
				// メソッド終了時またはパニック発生時に、ファイルポインタをクローズする
				switch recover() {
				default:
					if err := f.Close(); err != nil {
						panic(err)
					}
				}
			}(book)

			sheetNames, err := getSheetNames(bookName)
			if err != nil {
				panic(err)
			}

			wg := &sync.WaitGroup{}
			for _, sheetName := range sheetNames {
				wg.Add(1)
				sheetName := sheetName
				pastePictureOnSheet(wg, book, cf, bookName, sheetName)
			}
			wg.Wait()
			if err := book.SaveAs(name); err != nil {
				panic(err)
			}
		}(bn)
	}
	return nil
}

func isExistSheetName(book *excelize.File, name string) bool {
	for _, sheetName := range book.GetSheetList() {
		if strings.EqualFold(sheetName, name) {
			return true
		}
	}
	return false
}

func pastePictureOnSheet(wg *sync.WaitGroup, book *excelize.File, cf *conf.Config, bookName, sheetName string) {
	if !isExistSheetName(book, sheetName) {
		book.NewSheet(sheetName)
		if cf.Template.IsSheetSpecification() {
			if err := book.CopySheet(book.GetSheetIndex(cf.Template.SheetName), book.GetSheetIndex(sheetName)); err != nil {
				panic(err)
			}
		}
	}

	go func() {
		if err := pastePictures(wg, book, bookName, sheetName, cf.TargetCol, cf.TargetRow, cf.Offset); err != nil {
			panic(err)
		}
	}()
}

func pastePictures(wg *sync.WaitGroup, file *excelize.File, bookName, sheetName, targetCol string, targetRow, imageOffset int) error {
	defer func() {
		switch recover() {
		default:
			wg.Done()
		}
	}()

	imagePath := strings.Join([]string{consts.InputDirectory, bookName, sheetName}, string(os.PathSeparator))
	pictures, err := getDirNames(imagePath, func(de os.DirEntry) bool {
		return de.IsDir()
	})
	if err != nil {
		return err
	}

	var currentRow = targetRow
	for _, picture := range pictures {
		picture := strings.Join([]string{imagePath, picture}, string(os.PathSeparator))
		if !isExist(picture) {
			continue
		}

		targetCell := fmt.Sprintf("%s%d", targetCol, currentRow)
		if err := file.AddPicture(sheetName, targetCell, picture, consts.PictureOption); err != nil {
			return err
		}

		rowHeightPoint, err := file.GetRowHeight(sheetName, 1)
		if err != nil {
			return err
		}

		imgHeight, _, err := getImageSize(picture)
		if err != nil {
			return err
		}

		rowHeightPixel := point2Pixel(rowHeightPoint)
		currentRow += int(roundUp(float64(imgHeight)/rowHeightPixel, 0)) + imageOffset
	}
	return nil
}

func getExcelFileNames() ([]string, error) {
	return getDirNames(consts.InputDirectory+string(os.PathSeparator), func(de os.DirEntry) bool {
		// src直下のディレクトリ名がエビデンスファイル名となるため、ディレクトリ以外はスキップ
		return !de.IsDir()
	})
}

func getSheetNames(path string) ([]string, error) {
	return getDirNames(strings.Join([]string{consts.InputDirectory, path}, string(os.PathSeparator)), func(de os.DirEntry) bool {
		return !de.IsDir()
	})
}
