// Code generated by "genprotocol -ver=213afa194ef0e682076c6a0cbf801946c13d343cc54330be7c4557e46057a498 -basedir=. -prefix=w3d -statstype=int"

package w3d_statnoti

import (
	"fmt"
	"net/http"
	"sync"
	"text/template"

	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_idnoti"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_packet"
)

func (ns *StatNotification) String() string {
	return fmt.Sprintf("StatNotification[%v]", len(ns))
}

type StatNotification [w3d_idnoti.NotiID_Count]StatRow

func New() *StatNotification {
	ns := new(StatNotification)
	for i := 0; i < w3d_idnoti.NotiID_Count; i++ {
		ns[i].Name = w3d_idnoti.NotiID(i).String()
	}
	return ns
}
func (ns *StatNotification) Add(hd w3d_packet.Header) {
	if int(hd.Cmd) >= w3d_idnoti.NotiID_Count {
		return
	}
	ns[hd.Cmd].add(hd)
}
func (ns *StatNotification) ToWeb(w http.ResponseWriter, r *http.Request) error {
	tplIndex, err := template.New("index").Parse(`
	<html><head><title>Notification packet stat Info</title></head><body>
	<table border=1 style="border-collapse:collapse;">` +
		HTML_tableheader +
		`{{range $i, $v := .}}` +
		HTML_row +
		`{{end}}` +
		HTML_tableheader +
		`</table><br/>
	</body></html>`)
	if err != nil {
		return err
	}
	if err := tplIndex.Execute(w, ns); err != nil {
		return err
	}
	return nil
}

const (
	HTML_tableheader = `<tr>
	<th>Name</th>
	<th>Count</th>
	<th>Total Byte</th>
	<th>Max Byte</th>
	<th>Avg Byte</th>
	</tr>`
	HTML_row = `<tr>
	<td>{{$v.Name}}</td>
	<td>{{$v.Count }}</td>
	<td>{{$v.TotalByte }}</td>
	<td>{{$v.MaxByte }}</td>
	<td>{{printf "%10.3f" $v.Avg }}</td>
	</tr>
	`
)

type StatRow struct {
	mutex     sync.Mutex
	Name      string
	Count     int
	TotalByte int
	MaxByte   int
}

func (ps *StatRow) add(hd w3d_packet.Header) {
	ps.mutex.Lock()
	ps.Count++
	n := int(hd.BodyLen()) + w3d_packet.HeaderLen
	ps.TotalByte += n
	if n > ps.MaxByte {
		ps.MaxByte = n
	}
	ps.mutex.Unlock()
}
func (ps *StatRow) Avg() float64 {
	return float64(ps.TotalByte) / float64(ps.Count)
}
