// Code generated by "genprotocol -ver=213afa194ef0e682076c6a0cbf801946c13d343cc54330be7c4557e46057a498 -basedir=. -prefix=w3d -statstype=int"

package w3d_loopwsgorilla

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_const"
	"github.com/kasworld/gowasm3dgame/protocol_w3d/w3d_packet"
)

var bufPool = w3d_packet.NewPool(w3d_const.PacketBufferPoolSize)

func SendControl(
	wsConn *websocket.Conn, mt int, PacketWriteTimeOut time.Duration) error {

	return wsConn.WriteControl(mt, []byte{}, time.Now().Add(PacketWriteTimeOut))
}

func SendPacket(wsConn *websocket.Conn, sendBuffer []byte) error {
	return wsConn.WriteMessage(websocket.BinaryMessage, sendBuffer)
}

func SendLoop(sendRecvCtx context.Context, SendRecvStop func(), wsConn *websocket.Conn,
	timeout time.Duration,
	SendCh chan w3d_packet.Packet,
	marshalBodyFn func(interface{}, []byte) ([]byte, byte, error),
	handleSentPacketFn func(header w3d_packet.Header) error,
) error {

	defer SendRecvStop()
	var err error
loop:
	for {
		select {
		case <-sendRecvCtx.Done():
			err = SendControl(wsConn, websocket.CloseMessage, timeout)
			break loop
		case pk := <-SendCh:
			if err = wsConn.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
				break loop
			}
			oldbuf := bufPool.Get()
			sendBuffer, err := w3d_packet.Packet2Bytes(&pk, marshalBodyFn, oldbuf)
			if err != nil {
				bufPool.Put(oldbuf)
				break loop
			}
			if err = SendPacket(wsConn, sendBuffer); err != nil {
				bufPool.Put(oldbuf)
				break loop
			}
			if err = handleSentPacketFn(pk.Header); err != nil {
				bufPool.Put(oldbuf)
				break loop
			}
			bufPool.Put(oldbuf)
		}
	}
	return err
}

func RecvLoop(sendRecvCtx context.Context, SendRecvStop func(), wsConn *websocket.Conn,
	timeout time.Duration,
	HandleRecvPacketFn func(header w3d_packet.Header, body []byte) error) error {

	defer SendRecvStop()
	var err error
loop:
	for {
		select {
		case <-sendRecvCtx.Done():
			break loop
		default:
			if err = wsConn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
				break loop
			}
			if header, body, lerr := RecvPacket(wsConn); lerr != nil {
				if operr, ok := lerr.(*net.OpError); ok && operr.Timeout() {
					continue
				}
				err = lerr
				break loop
			} else {
				if err = HandleRecvPacketFn(header, body); err != nil {
					break loop
				}
			}
		}
	}
	return err
}

func RecvPacket(wsConn *websocket.Conn) (w3d_packet.Header, []byte, error) {
	mt, rdata, err := wsConn.ReadMessage()
	if err != nil {
		return w3d_packet.Header{}, nil, err
	}
	if mt != websocket.BinaryMessage {
		return w3d_packet.Header{}, nil, fmt.Errorf("message not binary %v", mt)
	}
	return w3d_packet.NewRecvPacketBufferByData(rdata).GetHeaderBody()
}
