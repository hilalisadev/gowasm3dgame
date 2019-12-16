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

package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kasworld/gowasm3dgame/game/gameconst"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_authorize"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_gob"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_idcmd"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_packet"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_serveconnbyte"
	"github.com/kasworld/uuidstr"
)

func (svr *Server) initServiceWeb(ctx context.Context) {
	webMux := http.NewServeMux()
	webMux.Handle("/",
		http.FileServer(http.Dir(svr.config.ClientDataFolder)),
	)
	webMux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		svr.serveWebSocketClient(ctx, w, r)
	})
	svr.clientWeb = &http.Server{
		Handler: webMux,
		Addr:    svr.config.ServicePort,
	}
	svr.marshalBodyFn = w3d_gob.MarshalBodyFn
	svr.unmarshalPacketFn = w3d_gob.UnmarshalPacket
	svr.DemuxReq2BytesAPIFnMap = [...]func(
		me interface{}, hd w3d_packet.Header, rbody []byte) (
		w3d_packet.Header, interface{}, error){
		w3d_idcmd.Invalid:   svr.bytesAPIFn_ReqInvalid,
		w3d_idcmd.MakeTeam:  svr.bytesAPIFn_ReqMakeTeam,
		w3d_idcmd.Act:       svr.bytesAPIFn_ReqAct,
		w3d_idcmd.Heartbeat: svr.bytesAPIFn_ReqHeartbeat,
	}

}

func CheckOrigin(r *http.Request) bool {
	return true
}

func (svr *Server) serveWebSocketClient(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: CheckOrigin,
	}
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("upgrade %v\n", err)
		return
	}

	stg := svr.stage
	connID := uuidstr.New()
	c2sc := w3d_serveconnbyte.NewWithStats(
		connID, // connid
		gameconst.SendBufferSize,
		w3d_authorize.NewAllSet(),
		svr.SendStat, svr.RecvStat,
		svr.apiStat,
		svr.notiStat,
		svr.errorStat,
		svr.DemuxReq2BytesAPIFnMap)

	// add to conn manager
	svr.connManager.Add(connID, c2sc)
	stg.Conns.Add(connID, c2sc)

	// start client service
	c2sc.StartServeWS(ctx, wsConn,
		gameconst.ReadTimeoutSec, gameconst.WriteTimeoutSec, svr.marshalBodyFn)
	wsConn.Close()
	// connection cleanup here

	// del from conn manager
	svr.connManager.Del(connID)
	stg.Conns.Del(connID)
}