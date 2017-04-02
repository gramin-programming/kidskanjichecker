package main

import (
	"testing"
	"bytes"
	"time"
	"archive/zip"
	"os"
	"bufio"
	"fmt"
	"strings"
	"io/ioutil"
)

func TestIsMatchTwoDimensionalArray(t *testing.T) {
	learning := [][]string{
		[]string{"一", "右", "雨", "円", "王", "音", "下", "火", "花"},
	}
	isMatch := isMatchTwoDimensionalArray(learning, "右")
	if ! isMatch {
		t.Error("Not Match")
	}
	isMatch = isMatchTwoDimensionalArray(learning, "右右")
	if isMatch {
		t.Error("Not Match")
	}
	isMatch = isMatchTwoDimensionalArray(learning, "")
	if isMatch {
		t.Error("Not Match")
	}
}

func TestGetMaxYearLearning(t *testing.T){
	checkLearning := getMaxYearLearning(0)
	if l := len(checkLearning);l != 0{
		t.Error("Length is " + string(l))
	}
	checkLearning = getMaxYearLearning(1)
	if len(checkLearning) != 1{
		t.Error("Not Match")
	}
	checkLearning = getMaxYearLearning(3)
	if len(checkLearning) != 3{
		t.Error("Not Match")
	}
	checkLearning = getMaxYearLearning(7)
	if len(checkLearning) != 7{
		t.Error("Not Match")
	}
	// 7より上は 7扱いになる
	checkLearning = getMaxYearLearning(8)
	if l := len(checkLearning);l != 7{
		t.Error("Not Match")
	}
}

func TestShowLine(t *testing.T){
	buf := &bytes.Buffer{}

	line := "あい玉うえお足かき早く"
	showLine(buf, line, true, 1)
	output := buf.String()
	if output != "あい玉うえお足かき早く\n" {
		t.Error("Actual is " + output)
	}
	buf.Reset()

	line = "あいうえお確刊かきく"
	showLine(buf, line, true, 3)
	output = buf.String()
	if output != "あいうえお_確__刊_かきく\n" {
		t.Error("Actual is " + output)
	}
	buf.Reset()

	line = "あいう側えお解格か胸き雪く"
	showLine(buf, line, true, 5)
	output = buf.String()
	if output != "あいう側えお解格か_胸_き雪く\n" {
		t.Error("Actual is " + output)
	}
	buf.Reset()

	line = "あいう側えお解格か胸き粋雪く"
	showLine(buf, line, true, 7)
	output = buf.String()
	if output != "あいう側えお解格か胸き粋雪く\n" {
		t.Error("Actual is " + output)
	}
	buf.Reset()

	line = "あいう側えお解格か胸き粋雪く"
	showLine(buf, line, false, 3)
	output = buf.String()
	if !strings.Contains(output, "あいう") {
		t.Error("Actual is " + output)
	}
	buf.Reset()
}

func TestCheckKidsKanjiForOdp(t *testing.T){
	odpContent := []byte(`<?xml version="1.0" encoding="UTF-8"?>
	<office:document-content>
		<office:body>
			<office:presentation>
				<draw:page>
					<draw:frame>
						<draw:text-box>
							<text:p>
								<text:span>
								あい解格か胸き
								</text:span>
							</text:p>
						</draw:text-box>
					</draw:frame>
						<draw:custom-shape>
							<text:p>
								<text:span>
									側えお
								</text:span>
							</text:p>
						</draw:custom-shape>
						...
				</draw:page>
					...
			</office:presentation>
		</office:body>
	</office:document-content></xml>`)
	tmpFile, err := ioutil.TempFile("", "kids-kanji-checker-")
	if err != nil{
		panic(err)
	}
	defer os.Remove(tmpFile.Name()) // clean up
	fmt.Println("-----------")
	if _, err := tmpFile.Write(odpContent); err != nil {
		panic(err)
	}
	if err := tmpFile.Close(); err != nil {
		panic(err)
	}
	
	buf := &bytes.Buffer{}	
	checkKidsKanjiForOdp(buf, tmpFile.Name(), true, 3)
	output := buf.String()
	expect := "あい_解__格_か_胸_き_側_えお"
	o := strings.Trim(output, " 	\r\n")
	o = strings.Replace(o, "\t", "", -1)
	o = strings.Replace(o, " ", "", -1)
	o = strings.Replace(o, "\r", "", -1)
	o = strings.Replace(o, "\n", "", -1)
	if o != expect {
		t.Error("Actual is \"" + o + "\"")
	}
	buf.Reset()
}


func TestCheckKidsKanjiForDocx(t *testing.T){
	odpContent := []byte(`<w:document>
		<w:body>
			<w:p>
				<w:r>
					<w:t>あい解格か胸き側えお</w:t>
				</w:r>
				<w:r>
					<w:tab/>
				</w:r>
			</w:p>
			<w:p>....</w:p>
		</w:body>
	</w:document>`)

	tmpFile, err := ioutil.TempFile("", "kids-kanji-checker-")
	if err != nil{
		panic(err)
	}
	defer os.Remove(tmpFile.Name()) // clean up
	fmt.Println("-----------")
	if _, err := tmpFile.Write(odpContent); err != nil {
		panic(err)
	}
	if err := tmpFile.Close(); err != nil {
		panic(err)
	}

	buf := &bytes.Buffer{}
	checkKidsKanjiForDocx(buf, tmpFile.Name(), true, 3)
	output := buf.String()
	expect := "あい_解__格_か_胸_き_側_えお"
	o := strings.Trim(output, " 	\r\n")
	o = strings.Replace(o, "\t", "", -1)
	o = strings.Replace(o, " ", "", -1)
	o = strings.Replace(o, "\r", "", -1)
	o = strings.Replace(o, "\n", "", -1)
	if o != expect {
		t.Error("Actual is \"" + o + "\"")
	}
	buf.Reset()
}

func TestCheckKidsKanjiForTxt(t *testing.T){
	contents := []byte(`あい解格か胸き側えお`)

	tmpFile, err := ioutil.TempFile("", "kids-kanji-checker-")
	if err != nil{
		panic(err)
	}
	name := tmpFile.Name()
	defer os.Remove(tmpFile.Name()) // clean up
	fmt.Println("-----------")
	if _, err := tmpFile.Write(contents); err != nil {
		panic(err)
	}
	if err := tmpFile.Close(); err != nil {
		panic(err)
	}
	fd, _ := os.Open(name)

	buf := &bytes.Buffer{}
	checkKidsKanjiForTxt(buf, bufio.NewScanner(fd),true, 3)
	output := buf.String()
	expect := "あい_解__格_か_胸_き_側_えお"
	o := strings.Trim(output, " 	\r\n")
	o = strings.Replace(o, "\t", "", -1)
	o = strings.Replace(o, " ", "", -1)
	o = strings.Replace(o, "\r", "", -1)
	o = strings.Replace(o, "\n", "", -1)
	if o != expect {
		t.Error("Actual is \"" + o + "\"")
	}
	buf.Reset()
}

func TestDumpFileInZip(t *testing.T){
	contents := []byte("あいうえお")
	b := new(bytes.Buffer)
	w := zip.NewWriter(b)
	header := &zip.FileHeader{
		Name:         "content.xml",
		Method:       zip.Store,
		ModifiedTime: uint16(time.Now().UnixNano()),
		ModifiedDate: uint16(time.Now().UnixNano()),
	}
	f, err := w.CreateHeader(header)
	if err != nil{
		panic(err)
	}
	f.Write(contents)
	w.Close()
	// buffer を file に dump
	tmpFile, err := ioutil.TempFile("", "kids-kanji-checker-")
	if err != nil {
		panic(err)
	}
	tmpFile.Write(b.Bytes())
	tmpFile.Close()

	path, err := dumpFileInZip(tmpFile.Name(), "content.xml")
	if err != nil {
		panic(err)
	}
	if !strings.HasSuffix(path, "content.xml") {
		t.Error("Actual is \"" + path + "\"")
	}

}