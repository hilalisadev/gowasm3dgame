// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wasmclient

import (
	"bytes"
	"fmt"
	"syscall/js"

	"github.com/kasworld/goguelike2/enum/achievetype"
	"github.com/kasworld/gowasm3dgame/enums/acttype"
	"github.com/kasworld/gowasm3dgame/game/gameconst"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_version"
)

func (app *WasmClient) updataServiceInfo() {
	msgCopyright := `</hr>Copyright 2019 SeukWon Kang 
		<a href="https://github.com/kasworld/gowasm3dgame" target="_blank">gowasm3dgame</a>`

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "gowasm3dgame webclient<br/>")
	fmt.Fprintf(&buf, "Protocol %v<br/>", w3d_version.ProtocolVersion)
	fmt.Fprintf(&buf, "Data %v<br/>", gameconst.DataVersion)
	fmt.Fprintf(&buf, "%v<br/>", msgCopyright)
	div := js.Global().Get("document").Call("getElementById", "serviceinfo")
	div.Set("innerHTML", buf.String())
}

func (app *WasmClient) updateDebugInfo() {
	var buf bytes.Buffer
	fmt.Fprintf(&buf,
		"%v<br/>Ping %v<br/>ServerClientTickDiff %v<br/>",
		app.DispInterDur, app.PingDur, app.ServerClientTictDiff,
	)
	fmt.Fprintf(&buf,
		"scene obj count %v<br/>geomatry cache count %v<br/>material cache count %v<br/>",
		len(app.vp.jsSceneObjs),
		len(app.vp.geometryCache),
		len(app.vp.materialCache),
	)

	div := js.Global().Get("document").Call("getElementById", "debuginfo")
	div.Set("innerHTML", buf.String())
}

func (app *WasmClient) updateSysmsg() {
	app.systemMessage = app.systemMessage.GetLastN(achievetype.AchieveType_Count + 100)
	div := js.Global().Get("document").Call("getElementById", "sysmsg")
	div.Set("innerHTML", app.systemMessage.ToHtmlStringRev())
}

func (app *WasmClient) updateTeamStatsInfo() {
	var buf bytes.Buffer

	if stats := app.statsInfo; stats != nil {
		fmt.Fprintf(&buf, "Stage %v<br/>", stats.UUID)

		buf.WriteString(`<table border=1 style="border-collapse:collapse;">`)

		buf.WriteString(`<colgroup>`)
		fmt.Fprintf(&buf, `<col >`)
		for _, v := range stats.Stats {
			fmt.Fprintf(&buf, `<col style="background-color:%v">`,
				v.Color24.ToHTMLColorString())
		}
		buf.WriteString(`</colgroup>`)

		buf.WriteString(`<tr><th>act\team</th>`)
		for ti, _ := range stats.Stats {
			fmt.Fprintf(&buf, "<th >%v</th>",
				ti)
		}
		buf.WriteString(`</tr>`)

		buf.WriteString(`<tr><td>UUID</td>`)
		for _, v := range stats.Stats {
			fmt.Fprintf(&buf, "<td>%v</td>", v.UUID)
		}
		buf.WriteString(`</tr>`)

		buf.WriteString(`<tr><td>AP</td>`)
		for _, v := range stats.Stats {
			fmt.Fprintf(&buf, "<td>%v</td>", v.AP)
		}
		buf.WriteString(`</tr>`)

		buf.WriteString(`<tr><td>Alive</td>`)
		for _, v := range stats.Stats {
			fmt.Fprintf(&buf, "<td>%v</td>", v.Alive)
		}
		buf.WriteString(`</tr>`)

		for acti := 0; acti < acttype.ActType_Count; acti++ {
			fmt.Fprintf(&buf, "<tr><td>%v</td>", acttype.ActType(acti))
			for ti, _ := range stats.Stats {
				fmt.Fprintf(&buf, "<td>%v</td>",
					stats.Stats[ti].ActStats[acti])
			}
			buf.WriteString(`</tr>`)
		}
		buf.WriteString(`</table>`)
	}

	div := js.Global().Get("document").Call("getElementById", "teamstatsinfo")
	div.Set("innerHTML", buf.String())
}
