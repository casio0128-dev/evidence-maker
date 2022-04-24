package utils

import (
	"evidence-maker/conf"
	_const "evidence-maker/const"
	"fmt"
	"github.com/xuri/excelize/v2"
	"image"
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
func OpenExcel(filePath string) (evd *excelize.File, err error) {
	if IsExist(filePath) {
		if evd, err = excelize.OpenFile(filePath); err != nil {
			return
		}
	} else {
		evd = excelize.NewFile()
	}
	return
}

func OutputExcelFile(wg *sync.WaitGroup, cf conf.Config) error {
	dirPath, err := os.MkdirTemp(_const.OutputDirectory, time.Now().Format(_const.OutputDirectoryPattern))
	if err != nil {
		return err
	}
	files, err := GetExcelFileNames()
	if err != nil {
		return err
	}

	for _, bn := range files {
		bn = bn
		wg.Add(1)
		go func(bookName string) {
			defer func() {
				// メソッド終了時またはパニック発生時に、waitGroupを終了にしてサブゴルーチンの終了を記録
				switch recover() {
				default:
					wg.Done()
				}
			}()

			name := fmt.Sprintf(_const.OutputExcelPattern, strings.Join([]string{dirPath, bookName}, string(os.PathSeparator)))
			book, err := OpenExcel(cf.Template.FilePath)
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

			sheetNames, err := GetSheetNames(bookName)
			if err != nil {
				panic(err)
			}

			for _, sheetName := range sheetNames {
				imagePath := strings.Join([]string{_const.InputDirectory, bookName, sheetName}, string(os.PathSeparator))
				if !IsExistSheetName(book, sheetName) {
					book.NewSheet(sheetName)
					if cf.Template.IsSheetSpecification() {
						if err := book.CopySheet(book.GetSheetIndex(cf.Template.SheetName), book.GetSheetIndex(sheetName)); err != nil {
							panic(err)
						}
					}
				}
				if err := PastePictures(book, imagePath, sheetName, cf.TargetCol, cf.TargetRow, cf.Offset); err != nil {
					panic(err)
				}
			}

			if err := book.SaveAs(name); err != nil {
				panic(err)
			}
		}(bn)
	}
	return nil
}

func IsExistSheetName(book *excelize.File, name string) bool {
	for _, sheetName := range book.GetSheetList() {
		if strings.EqualFold(sheetName, name) {
			return true
		}
	}
	return false
}

func PastePictures(file *excelize.File, path, sheetName, targetCol string, targetRow, imageOffset int) error {
	pictures, err := GetDirNames(path, func(de os.DirEntry) bool {
		return de.IsDir()
	})
	if err != nil {
		return err
	}

	var currentRow = targetRow
	for _, picture := range pictures {
		picture = strings.Join([]string{path, picture}, string(os.PathSeparator))
		if !IsExist(picture) {
			continue
		}

		targetCell := fmt.Sprintf("%s%d", targetCol, currentRow)
		if err := file.AddPicture(sheetName, targetCell, picture, _const.PictureOption); err != nil {
			return err
		}

		pict, err := os.Open(picture)
		if err != nil {
			return err
		}

		img, _, err := image.Decode(pict)
		if err != nil {
			return err
		}

		rowHeightPoint, err := file.GetRowHeight(sheetName, 1)
		if err != nil {
			return err
		}

		rowHeightPixel := Point2Pixel(rowHeightPoint)
		currentRow += int(RoundUp(float64(img.Bounds().Max.Y)/rowHeightPixel, 0)) + imageOffset
	}
	return nil
}

func GetExcelFileNames() ([]string, error) {
	return GetDirNames(_const.InputDirectory+string(os.PathSeparator), func(de os.DirEntry) bool {
		// src直下のディレクトリ名がエビデンスファイル名となるため、ディレクトリ以外はスキップ
		return !de.IsDir()
	})
}

func GetSheetNames(path string) ([]string, error) {
	return GetDirNames(strings.Join([]string{_const.InputDirectory, path}, string(os.PathSeparator)), func(de os.DirEntry) bool {
		return !de.IsDir()
	})
}
