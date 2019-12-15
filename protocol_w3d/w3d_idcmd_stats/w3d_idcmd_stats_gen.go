// Code generated by "genprotocol -ver=f5b1d289172cf84ad5d01b91533408be6b17961cf28ddd6fe767224298a8aedd -basedir=. -prefix=w3d -statstype=int"

package w3d_idcmd_stats

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_idcmd"
)

type CommandIDStat [w3d_idcmd.CommandID_Count]int

func (es *CommandIDStat) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "CommandIDStats[")
	for i, v := range es {
		fmt.Fprintf(&buf,
			"%v:%v ",
			w3d_idcmd.CommandID(i), v)
	}
	buf.WriteString("]")
	return buf.String()
}
func (es *CommandIDStat) Inc(e w3d_idcmd.CommandID) {
	es[e] += 1
}
func (es *CommandIDStat) Add(e w3d_idcmd.CommandID, v int) {
	es[e] += v
}
func (es *CommandIDStat) SetIfGt(e w3d_idcmd.CommandID, v int) {
	if es[e] < v {
		es[e] = v
	}
}
func (es *CommandIDStat) Get(e w3d_idcmd.CommandID) int {
	return es[e]
}

func (es *CommandIDStat) ToWeb(w http.ResponseWriter, r *http.Request) error {
	tplIndex, err := template.New("index").Funcs(IndexFn).Parse(`
		<html>
		<head>
		<title>CommandID statistics</title>
		</head>
		<body>
		<table border=1 style="border-collapse:collapse;">` +
		HTML_tableheader +
		`{{range $i, $v := .}}` +
		HTML_row +
		`{{end}}` +
		HTML_tableheader +
		`</table>
	
		<br/>
		</body>
		</html>
		`)
	if err != nil {
		return err
	}
	if err := tplIndex.Execute(w, es); err != nil {
		return err
	}
	return nil
}

func Index(i int) string {
	return w3d_idcmd.CommandID(i).String()
}

var IndexFn = template.FuncMap{
	"CommandIDIndex": Index,
}

const (
	HTML_tableheader = `<tr>
		<th>Name</th>
		<th>Value</th>
		</tr>`
	HTML_row = `<tr>
		<td>{{CommandIDIndex $i}}</td>
		<td>{{$v}}</td>
		</tr>
		`
)