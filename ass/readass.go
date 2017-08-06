package ass
import (
	"os"
	"fmt"
	"strings"
	"bufio"
	"strconv"
)

//The main struct for a .ass subtitle file. 


func floatit(in string)(out float64) {
	if in == "" {
		out = 0.0
	} else {
		outa, err := strconv.ParseFloat(in, 64)
		if err != nil {
			panic(err)
		}
		out = outa
	}
	return
}


func reversesplit(in []string)(out string) {
	out = strings.Join(in, ":")
	return
}

func intit(in string)(out int) {
	if in == "" {
		out = 0
	} else {
		outa, err := strconv.Atoi(in)
		if err != nil {
			panic(err)
		}
		out = outa
	}
	return
}


//Loads the .ass file and parses out the various possible valid lines. 
func Loadass(v *Ass, filepath string) {
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Cannot read file")
		os.Exit(1)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		splitter := strings.Split(scanner.Text(), ":")
		prefix := splitter[0]
		suffix := strings.Trim(reversesplit(splitter[1:]), " ")
		switch prefix {
		case "Title":
			v.ScriptInfo.Body.Title = suffix
		case "ScriptType":
			v.ScriptInfo.Body.ScriptType = suffix
		case "WrapStyle":
			v.ScriptInfo.Body.WarpStyle = suffix
		case "ScaledBorderAndShadow":
			v.ScriptInfo.Body.SBaShadow = suffix
		case "YCbCr Matrix":
			v.ScriptInfo.Body.YCbCrMatrix = suffix
		case "PlayResX":
			v.ScriptInfo.Body.PlayResX = intit(suffix)
		case "PlayResY":
			v.ScriptInfo.Body.PlayResY = intit(suffix)
		case "Audio File":
			v.PGarbage.Body.AudioFile = suffix
		case "Video File":
			v.PGarbage.Body.VideoFile = suffix
		case "Video AR Mode":
			v.PGarbage.Body.VideoARMode = suffix
		case "Video AR Value":
			v.PGarbage.Body.VideoARValue = floatit(suffix)
		case "Video Zoom Percent" :
			v.PGarbage.Body.VideoZoomPercent = floatit(suffix)
		case "Scroll Position" :
			v.PGarbage.Body.ScrollPosition = intit(suffix)
		case "Active Line" :
			v.PGarbage.Body.ActiveLine = intit(suffix)
		case "Video Position" :
			v.PGarbage.Body.VideoPos = intit(suffix)
		case "Style" :
			v.Styles.Body = append(v.Styles.Body, *Createstyle(suffix))
		case "Dialogue":
			v.Events.Body = append(v.Events.Body, *Createevent(suffix, prefix))
		case "Comment":
			v.Events.Body = append(v.Events.Body, *Createevent(suffix, prefix))
		}
	}
}

//Creates the event, takes a full .ass line par the Comment/Dialogue portion, and the event type as an argument.
//For example `Dialogue: 0,0:03:20.10,0:03:21.36,Default,,0,0,0,,` would parse to
//in = `0,0:03:20.10,0:03:21.36,Default,,0,0,0,`
//etype = Dialogue
func Createevent(in string, etype string) *Event{
	split := strings.Split(in, ",")
	return &Event {
		Format: etype,
		Layer: intit(split[0]),
		Start: split[1],
		End: split[2],
		Style: split[3],
		Name: split[4],
		MarginL: intit(split[5]),
		MarginR: intit(split[6]),
		MarginV: intit(split[7]),
		Effect: split[8],
		Text: split[9],
	}
}
//Takes a full .ass style line as an argument.
//Similar to Createevent, except Styles don't have multiple Formats, so we only take the format-less style string. 
func Createstyle(in string) *Style{
	split := strings.Split(in, ",")
	return &Style {
		Format: "Style: ",
		Name: split[0],
		Fontname: split[1],
		Fontsize: intit(split[2]),
		PrimaryColour: split[3],
		SecondaryColour: split[4],
		OutlineColour: split[5],
		Backcolour: split[6],
		Bold: intit(split[7]),
		Italic: intit(split[8]),
		Underline: intit(split[9]),
		StrikeOut: intit(split[10]),
		ScaleX: intit(split[11]),
		ScaleY: intit(split[12]),
		Spacing: intit(split[13]),
		Angle: intit(split[14]),
		BorderStyle: intit(split[15]),
		Outline: intit(split[16]),
		Shadow: intit(split[17]),
		Alignment: intit(split[18]),
		MarginL: intit(split[19]),
		MarginR: intit(split[20]),
		MarginV: intit(split[21]),
		Encoding: intit(split[22]),
	}
}

//Sets default headers for the various fields. 
func Setheaders(v *Ass) {
	v.ScriptInfo.Header = "[Script Info]"
	v.PGarbage.Header = "[Aegisub Project Garbage]"
	v.Styles.Header = "[V4+ Styles]"
	v.Events.Header = "[Events]"
	v.Styles.Format = "Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding"
	v.Events.Format = "Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text"
		
}

//Takes file path as argument. 
func ParseAss(filename string) *Ass{
	v := &Ass{}
	Setheaders(v)
	Loadass(v, filename)
	return v
}