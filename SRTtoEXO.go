package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	srt "github.com/suapapa/go_subtitle"
	"golang.org/x/text/encoding/japanese"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"gopkg.in/yaml.v3"
)

const HeaderFormat = `[exedit]
width={{.UserConfig.Width}}
height={{.UserConfig.Height}}
rate={{.UserConfig.FrameRate}}
scale=1
length={{.UserConfig.Length}}
audio_rate={{.UserConfig.AudioRate}}
audio_ch={{.UserConfig.AudioCh}}
`

const TextObjFormat = `
[{{.SrtScript.Idx}}]
start={{.SrtScript.Start}}
end={{.SrtScript.End}}
layer={{.UserConfig.Layer}}
overlay=1
camera=0
[{{.SrtScript.Idx}}.0]
_name=テキスト
サイズ={{.ObjConfig.Size}}
表示速度=0.0
文字毎に個別オブジェクト=0
移動座標上に表示する=0
自動スクロール=0
B={{.TextConfig.Bold}}
I={{.TextConfig.Italic}}
type={{.TextConfig.EffectType}}
autoadjust={{.TextConfig.AutoAdjust}}
soft={{.TextConfig.Soft}}
monospace={{.TextConfig.Monospace}}
align={{.TextConfig.Align}}
spacing_x={{.TextConfig.SpacingX}}
spacing_y={{.TextConfig.SpacingY}}
precision={{.TextConfig.Precision}}
color={{.TextConfig.Color}}
color2={{.TextConfig.Color2}}
font={{.TextConfig.Font}}
text={{.SrtScript.Text}}
[{{.SrtScript.Idx}}.1]
_name=標準描画
X={{.ObjConfig.X}}
Y={{.ObjConfig.Y}}
Z={{.ObjConfig.Z}}
拡大率={{.ObjConfig.Zoom}}
透明度={{.ObjConfig.AlphaBlend}}
回転={{.ObjConfig.Rotate}}
blend={{.ObjConfig.Blend}}
`
type Config struct {
	UserConfig `yaml:",inline"`
	ObjConfig  `yaml:",inline"`
	TextConfig `yaml:",inline"`
	SrtScript  `yaml:",inline"`
}
type UserConfig struct {
	FilePath  string `yaml:"FilePath"`
	MovieSize []int  `yaml:"MovieSize"`
	FrameRate int    `yaml:"FrameRate"`
	AudioRate int    `yaml:"AudioRate"`
	AudioCh   int    `yaml:"AudioCh"`
	Layer     int    `yaml:"Layer"`
	Length    int
	Width     int
	Height    int
}
type ObjConfig struct {
	X          float32 `yaml:"X"`
	Y          float32 `yaml:"Y"`
	Z          float32 `yaml:"Z"`
	Zoom       float32 `yaml:"Zoom"`
	AlphaBlend float32 `yaml:"AlphaBlend"`
	Rotate     float32 `yaml:"Rotate"`
	Size       float32 `yaml:"Size"`
	Blend      int     `yaml:"Blend"`
}
type TextConfig struct {
	Font       string  `yaml:"Font"`
	Bold       int    `yaml:"Bold"`
	Italic     int    `yaml:"Italic"`
	EffectType int     `yaml:"EffectType"`
	AutoAdjust int    `yaml:"AutoAdjust"`
	Soft       int    `yaml:"Soft"`
	Monospace  int    `yaml:"Monospace"`
	Align      int     `yaml:"Align"`
	SpacingX   float32 `yaml:"SpacingX"`
	SpacingY   float32 `yaml:"SpacingY"`
	Precision  int    `yaml:"Precision"`
	Color      string  `yaml:"Color"`
	Color2     string  `yaml:"Color2"`
}
type SrtScript struct {
	Idx   int
	Start int
	End   int
	Text  string
}

func main() {
	// read config.yaml
	var conf Config
	confFile, _ := os.ReadFile("config.yaml")
	yaml.Unmarshal(confFile, &conf)
	// .srt file open
	file, err := os.Open(conf.UserConfig.FilePath)
	// error
	if err != nil {
		panic(err)
	}
	fileName := file.Name()
	book, err := srt.ReadSrt(file)
	// error
	if err != nil {
		panic(err)
	}

	// header init
	conf.UserConfig.Length = int(book[len(book)-1].End.Seconds() * float64(conf.UserConfig.FrameRate))
	conf.UserConfig.Width = conf.UserConfig.MovieSize[0]
	conf.UserConfig.Height = conf.UserConfig.MovieSize[1]
	header := template.Must(template.New("header").Parse(HeaderFormat))
	str := new(strings.Builder)

	// header write
	exoFile, err := os.Create(fileName[:len(fileName)-4] + ".exo")
	// error
	if err != nil {
		panic(err)
	}
	defer exoFile.Close()
	// default part init
	sfjsEncoder := japanese.ShiftJIS.NewEncoder()
	Encoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()


	writer := transform.NewWriter(exoFile, sfjsEncoder)
	defer exoFile.Close()

	header.Execute(str, conf)
	writer.Write([]byte(ConvertCRLF(str.String())))
	str.Reset()

	// default part init
	textObj := template.Must(template.New("textObj").Parse(TextObjFormat))
	for _, script := range book {
		conf.SrtScript.Idx = script.Idx - 1
		conf.SrtScript.Start = TimeToFrame(script.Start, conf.FrameRate) + 1
		conf.SrtScript.End = TimeToFrame(script.End, conf.FrameRate)
		conf.SrtScript.Text = ConvertExoText(script.Text, Encoder)
		textObj.Execute(str, conf)
		writer.Write([]byte(ConvertCRLF(str.String())))
		str.Reset()
	}
}

func TimeToFrame(t time.Duration, fr int) int {
	return int(t.Seconds()) * fr
}

func ConvertExoText(str string, Encoder *encoding.Encoder) string {
	encStr, _, _ := transform.String(Encoder, str)
	ret := fmt.Sprintf("%x", encStr)
	ret = ret + strings.Repeat("0", 4096 - len(ret))
	return ret
}

func ConvertCRLF(str string) string {
    return strings.NewReplacer(
        "\r", "\r\n",
        "\n", "\r\n",
    ).Replace(str)
}
