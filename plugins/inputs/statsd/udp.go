package statsd

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"
)

func (ipt *input) setupUDPServer() error {
	address, err := net.ResolveUDPAddr(ipt.Protocol, ipt.ServiceAddress)
	if err != nil {
		l.Error(err)
		return err
	}

	conn, err := net.ListenUDP(ipt.Protocol, address)
	if err != nil {
		l.Error(err)
		return err
	}

	l.Infof("UDP listening on %q", conn.LocalAddr().String())
	ipt.UDPlistener = conn

	ipt.wg.Add(1)
	go func() {
		defer ipt.wg.Done()
		if err := ipt.udpListen(conn); err != nil {
			l.Warnf("udpListen: %s, ignored", err.Error())
		}
	}()
	return nil
}

// udpListen starts listening for udp packets on the configured port.
func (ipt *input) udpListen(conn *net.UDPConn) error {
	if ipt.ReadBufferSize > 0 {
		if err := ipt.UDPlistener.SetReadBuffer(ipt.ReadBufferSize); err != nil {
			return err
		}
	}

	buf := make([]byte, UDPMaxPacketSize)
	for {
		select {
		case <-ipt.done:
			return nil
		default:
			n, addr, err := conn.ReadFromUDP(buf)
			if err != nil {
				if !strings.Contains(err.Error(), "closed network") {
					l.Errorf("Error reading: %s", err.Error())
					continue
				}
				return err
			}

			l.Debugf("UDP: read %d bytes from %s", n, addr.IP.String())

			b, ok := ipt.bufPool.Get().(*bytes.Buffer)
			if !ok {
				return fmt.Errorf("bufPool is not a bytes buffer")
			}
			b.Reset()
			if _, err := b.Write(buf[:n]); err != nil {
				return err
			}
			select {
			case ipt.in <- job{
				Buffer: b,
				Time:   time.Now(),
				Addr:   addr.IP.String(),
			}:
			default:
				ipt.drops++
				if ipt.drops == 1 || ipt.AllowedPendingMessages == 0 || ipt.drops%ipt.AllowedPendingMessages == 0 {
					l.Errorf("Statsd message queue full. "+
						"We have dropped %d messages so far. "+
						"You may want to increase allowed_pending_messages in the config", ipt.drops)
				}
			}
		}
	}
}
