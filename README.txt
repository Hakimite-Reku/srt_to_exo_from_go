Srt to Exo from Go

・About
SRTファイルをAviutlに読み込めるようなEXOファイルに変換します。
MITライセンスの元、ご自由にお使いください。

使い方
1. zipファイルを解凍し、好きな場所に中身の3種類からご使用のWindowsパソコン環境に合わせたexeファイル(SRTtoEXO.exe)とconfig.sample.yamlファイルを配置してください。
補足：exeファイルは3種類ありますが、大抵の場合386かamd64のどちらかで動作すると思われます

2. config.sample.yamlファイルをconfig.yamlに名前を変更してください。

3. config.yamlファイル内のFilePath: 以降に変換してほしいSRTファイルへのパスを入力してください。

4. exeファイルをダブルクリックすると、SRTファイルと同じフォルダ内にEXOファイルが生成されます。EXOファイルはAviUtlにドラッグ＆ドロップで配置することができます。

- 注意：動画ファイルのサイズとフレームレートを合わせないとうまく動作しないので、ご自分の動画ファイルに合わせた設定をお願いします。

- config.yamlファイルの中身は配置するテキストオブジェクトの設定となっています。いろいろ試してみてください。

## Credits
- example.srt from [Fileformat](https://docs.fileformat.com/ja/video/srt/)
