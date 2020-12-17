## goenv のインストール

git clone https://github.com/syndbg/goenv.git ~/.goenv
echo 'export GOENV_ROOT="$HOME/.goenv"' >> ~/.bashrc
echo 'export PATH="$GOENV_ROOT/bin:$PATH"' >> ~/.bashrc

echo 'eval "$(goenv init -)"' >> ~/.bashrc
echo 'export PATH="$GOROOT/bin:$PATH"' >> ~/.bashrc
echo 'export PATH="$PATH:$GOPATH/bin"' >> ~/.bashrc

## go install

goenv install -l
goenv install "version"
goenv global "version"

## 参考にしたサイト
### ナチュラルハーモニー
http://blog.livedoor.jp/mtsk44h6-004/archives/3514537.html

### wasm + canvas 
https://golangtokyo.github.io/codelab/go-webassembly/
https://tech.basicinc.jp/articles/197
https://undersourcecode.hatenablog.com/entry/2019/08/02/194749
http://shinimae.hatenablog.com/entry/2016/09/29/193858
https://blog.narumium.net/2019/03/09/%E3%80%90go%E3%80%91ver1-12%E3%81%A7%E3%81%AEwebassembly/
https://beta-notes.way-nifty.com/blog/2020/04/post-4897c6.html
https://beta-notes.way-nifty.com/blog/2020/04/post-09a624.html

## HCL
https://www.konicaminolta.jp/instruments/knowledge/color/section2/02.html

## TODO:

### web app として提供

- フロントで file read とリサイズ, 処理結果受け取り
- gae で nh 処理するサーバ + 処理後イメージの gcs 格納 + パス返す
