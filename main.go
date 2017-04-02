package main

import (
	"flag"
	"fmt"
	"os"
	"bufio"
	"bytes"
	"unicode"
	"github.com/fatih/color"
	"strings"
	"path"
	"archive/zip"
	"path/filepath"
	"io/ioutil"
	"io"
	"encoding/xml"
	"kids-kanji-checker/log"
)

type Node struct {
	XMLName xml.Name
	Name    string `xml:"name,attr"`
	Content []byte `xml:",innerxml"`
	Nodes   []Node `xml:",any"`
}

// Recursive XML Parse
func walk(nodes []Node, f func(Node) bool) {
	for _, n := range nodes {
		if f(n) {
			walk(n.Nodes, f)
		}
	}
}

// Docx 形式のファイルから漢字をチェックする
func checkKidsKanjiForDocx(w io.Writer, contentPath string,
	isNoColorFlag bool, maxYear int) error {
	log.Debug("reading docx file...")
	xmlContent, _ := ioutil.ReadFile(contentPath)
	buf := bytes.NewBuffer(xmlContent)
	dec := xml.NewDecoder(buf)
	var n Node
	err := dec.Decode(&n)
	if err != nil {
		return err
	}
	var para string // paragraph

	walk([]Node{n}, func(n Node) bool {
		// Docx の場合、1文字ずつの <w:t> が大量にあるため、p で区切る
		if n.XMLName.Local == "p" {
			showLine(w, para, isNoColorFlag, maxYear)
			para = ""
		}
		if n.XMLName.Local == "t" {
			para += string(n.Content)
		}
		return true
	})
	return nil
}

// Odp 形式のファイルから漢字をチェックする
func checkKidsKanjiForOdp(w io.Writer, contentPath string,
	  isNoColorFlag bool, maxYear int) error {
	log.Debug("reading odp file...")
	xmlContent, _ := ioutil.ReadFile(contentPath)
	buf := bytes.NewBuffer(xmlContent)
	dec := xml.NewDecoder(buf)
	var n Node
	err := dec.Decode(&n)
	if err != nil {
		return err
	}
	blue := color.New(color.FgBlue)
	pageLine := "--------------------------------------------------------"
	// 再帰的に Xml を parse して span と page を出力する
	walk([]Node{n}, func(n Node) bool {
		if n.XMLName.Local == "page" {
			fmt.Println()
			blue.Println(pageLine)
			blue.Printf("%5s\n", string(n.Name))
			blue.Println(pageLine)
		}
		if n.XMLName.Local == "span" {
			showLine(w, string(n.Content), isNoColorFlag, maxYear)
		}
		return true
	})
	return nil
}

// 2次元配列内の文字列とマッチするかどうかをチェックする
func isMatchTwoDimensionalArray(array [][]string, char string) bool {
	for _, oneArray := range array {
		for _, str := range oneArray {
			if str == char {
				return true
			}
		}
	}
	return false
}

// maxYear に該当する常用漢字を取得
func getMaxYearLearning(maxYear int) [][]string{
	var checkLearning [][]string
	// maxYear に該当する漢字を取得
	for i, char := range Learning {
		if i >= maxYear {
			break
		}
		checkLearning = append(checkLearning, char)
	}
	return checkLearning
}

// 入力文字列に対して漢字をチェックして常用漢字であれば色を変えて出力する
func showLine(w io.Writer, line string, isNoColorFlag bool,
	maxYear int) {
	// maxYear に該当する常用漢字を取得
	checkLearning := getMaxYearLearning(maxYear)

	// 常用漢字かどうかをチェック
	cyan := color.New(color.FgCyan)
	for _, str := range line {
		runeVal := rune(str)
		// 漢字かどうか
		if ! unicode.Is(unicode.Han, runeVal) {
			fmt.Fprintf(w, "%s", string(str))
			continue
		}

		// 常用漢字かどうか
		if isMatchTwoDimensionalArray(checkLearning, string(str)) {
			fmt.Fprintf(w, "%s", string(str))
			continue
		}
		if isNoColorFlag {
			fmt.Fprintf(w, "_%s_", string(str))
		} else {
			cyan.Fprintf(w, "%s", string(str))
		}

	}
	fmt.Fprint(w,"\n")
}

// Text 形式のファイルから漢字をチェックする
func checkKidsKanjiForTxt(w io.Writer, scanner *bufio.Scanner,
	isNoColorFlag bool, maxYear int) {
	// http://knightso.hateblo.jp/entry/2014/06/24/090719
	log.Debug("scanning text file...")
	for scanner.Scan() {
		showLine(w, scanner.Text(), isNoColorFlag, maxYear)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

// Zip を解凍して dumpFileName に指定された Temporary File の Path を返す
func dumpFileInZip(zipFilePath string, dumpFileName string) (string, error) {
	tmpDir, err := ioutil.TempDir("", "kids-kanji-checker-")
	if err != nil {
		return "", err
	}

	// https://gist.github.com/hnaohiro/4572580
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return "", err
	}
	defer r.Close()
	var outputPath string
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			rc.Close()
			return "", err
		}

		if f.Name != dumpFileName {
			continue
		}

		outputPath = filepath.Join(tmpDir, f.Name)
		os.MkdirAll(filepath.Dir(outputPath), 0755)
		f, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			f.Close()
			return "", err
		}

		_, err = io.Copy(f, rc)
		if err != nil {
			return "", err
		}

		f.Close()
		rc.Close()
	}
	return outputPath, nil
}

// Docx 形式のファイルから Content の Path を取得する
func getDocxContentPath(filePath string) (string, error) {
	return dumpFileInZip(filePath, filepath.Join("word", "document.xml"))
}

// Odp 形式のファイルから Content の Path を取得する
func getOdpContentPath(filePath string) (string, error) {
	return dumpFileInZip(filePath, "content.xml")
}

func Usage() {
	fmt.Printf("Usage: %s [OPTIONS] argument ...\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	var fp *os.File
	var err error

	versionFlag := flag.Bool("version", false, "Version")
	stdinFlag := flag.Bool("stdin", false, "stdin")
	inputFileFlag := flag.String("input-file", "", "Input file")
	isNoColorFlag := flag.Bool("no-color", false, "チェック時の Color 表示をやめるか")
	maxYearFlag := flag.Int("max-year", 3, "何年生までの常用漢字をチェックするか（中学生以上は 7）")
	fileTypeFlag := flag.String("fileType", "", "ファイルタイプ(odp か docx のみ)。基本は拡張子で判断")
	isQuiet := flag.Bool("is-quiet", false, "余計な表示をしない")
	flag.Usage = Usage
	flag.Parse()

	var fileType string
	var contentPath string

	if *versionFlag {
		fmt.Println("Git Commit:", GitCommit)
		fmt.Println("Version:", Version)
		if VersionPrerelease != "" {
			fmt.Println("Version PreRelease:", VersionPrerelease)
		}
		return
	}

	if ! *isQuiet {
		c := color.New(color.FgGreen)
		c.Println("==================================================")
		c.Println("    kids-kanji-checker")
		c.Println("==================================================")
	}

	if *stdinFlag {
		log.Debug("input file is stdin")
		fp = os.Stdin
		fileType = "txt"
	} else if *inputFileFlag == "" {
		flag.Usage()
		os.Exit(1)
	} else if strings.ToLower(path.Ext(*inputFileFlag)) == ".txt" {
		log.Debug("input file is text")
		fp, err = os.Open(*inputFileFlag)
		if err != nil {
			panic(err)
		}
		defer fp.Close()
		fileType = "txt"

	} else if strings.ToLower(path.Ext(*inputFileFlag)) == ".odp" {
		log.Debug("input file is odp")
		contentPath, err = getOdpContentPath(*inputFileFlag)
		fp, err = os.Open(contentPath)
		if err != nil {
			panic(err)
		}
		defer fp.Close()
		log.Debug("contentPath is %s\n", contentPath)
		fileType = "odp"
	} else if strings.ToLower(path.Ext(*inputFileFlag)) == ".docx" {
		log.Debug("input file is docx")
		contentPath, err = getDocxContentPath(*inputFileFlag)
		fp, err = os.Open(contentPath)
		if err != nil {
			panic(err)
		}
		defer fp.Close()
		log.Debug("contentPath is %s\n", contentPath)
		fileType = "docx"
	}

	if *fileTypeFlag != "" {
		fileType = *fileTypeFlag
	}

	if fileType == "txt" {
		checkKidsKanjiForTxt(os.Stdout,
			bufio.NewScanner(fp), *isNoColorFlag, *maxYearFlag)
	} else if fileType == "odp" {
		checkKidsKanjiForOdp(os.Stdout, contentPath,
			*isNoColorFlag, *maxYearFlag)
	} else if fileType == "docx" {
		checkKidsKanjiForDocx(os.Stdout, contentPath,
			*isNoColorFlag, *maxYearFlag)
	}
}
