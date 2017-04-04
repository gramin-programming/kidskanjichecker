# kids-kanji-checker
日本の小学生・中学生が使う常用漢字以外の漢字が使われていないかをチェックするツールです。

![イメージ](https://raw.githubusercontent.com/gramin-programming/kidskanjichecker/master/kids-kanji-checker-image.png)

# インストール

[リリース](https://github.com/gramin-programming/kidskanjichecker/releases) からダウンロードして実行するだけです。

- kids-kanji-checker_darwin_386: Mac 32bit 向け
- kids-kanji-checker_darwin_amd64: Mac 64bit 向け。最近の Mac であればこれを使って下さい
- kids-kanji-checker_linux_386: Linux 32bit 向け
- kids-kanji-checker_linux_amd64: Linux 64bit 向け。最近のデスクトップ Linux であればこれを使って下さい
- kids-kanji-checker_linux_arm: Linux ARM 向け。
- kids-kanji-checker_windows_386.exe: Windows 32bit 向け。
- kids-kanji-checker_windows_amd64.exe: Windows 64bit 向け。最近の Windows であればこれを使って下さい

ダウンロード後に kids-kanji-checker 等に名前を変えても問題ありません。

この後、Linux・Mac 環境では権限を付与します。
```
$ chmod 0755 kids-kanji-checker_darwin_amd64
```
どこからでも実行したい場合は /usr/local/bin あたりに移動して下さい。

```
$ sudo mv -v kids-kanji-checker_darwin_amd64 /usr/local/bin/.
```

# 使い方
```
Usage: ./kids-kanji-checker_linux_amd64 [OPTIONS] argument ...
  -fileType string
        ファイルタイプ(odp か docx のみ)。基本は拡張子で判断
  -input-file string
        Input file
  -is-quiet
        余計な表示をしない
  -max-year int
        何年生までの常用漢字をチェックするか（中学生以上は 7） (default 3)
  -no-color
        チェック時の Color 表示をやめるか
  -stdin
        stdin
  -version
        Version
```

# 利用例
この例ではリリース後に kids-kanji-checker というファイルに名前を変更しています。

## テキスト文書に対して 3年生までに習わない漢字を着色して表示する

```
$ ./kids-kanji-checker -input-file sample/sample.txt -max-year 3
```

3年生までに習わない漢字が着色されて表示されます

## Microsoft Word（Docx）文書に対して小学 1年生までに習わない漢字を着色して表示する
```
$ ./kids-kanji-checker -input-file sample/sample.docx -max-year 1
```

## 中学までに習わない漢字を標準入力から受け取って表示する
```
$ cat sample/sample.docx | ./kids-kanji-checker -stdin -max-year 7
```

## 着色がサポートされていない環境で LibreOffice Impress 文書に対して小学 3年生（-max-yearを指定しない場合はこうなります）までに習わない漢字を表示する
```
$ ./kids-kanji-checker -input-file sample/sample.odp -no-color
```
