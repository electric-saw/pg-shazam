package proxy

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"

	"github.com/electric-saw/pg-shazam/internal/pkg/auth"
	"github.com/electric-saw/pg-shazam/internal/pkg/definitions"
	"github.com/electric-saw/pg-shazam/internal/pkg/log"
	"github.com/electric-saw/pg-shazam/internal/pkg/state"
	"github.com/jackc/pgproto3/v2"
)

func (p *Proxy) startup(client *definitions.FrontendClient) (bool, error) {
	for {
		startupMessage, err := client.Backend.ReceiveStartupMessage()
		if err != nil {
			return false, fmt.Errorf("error receiving startup message: %s", err.Error())
		}

		switch msg := startupMessage.(type) {
		case *pgproto3.StartupMessage:
			if db, ok := msg.Parameters["database"]; ok {
				client.CurrDatabase = db
			}

			return p.handlePass(client, msg)

		case *pgproto3.SSLRequest:
			err := p.handleSSLReq(client)
			if err != nil {
				return false, err
			}
		case *pgproto3.CancelRequest:
			return false, p.sendCancel(client, msg)

		default:
			return false, fmt.Errorf("unknown startup message: %#v", startupMessage)
		}
	}
}

func (p *Proxy) handleSSLReq(client *definitions.FrontendClient) error {
	_, err := client.Conn.Write([]byte("N"))
	if err != nil {
		return fmt.Errorf("error sending deny SSL request: %w", err)
	} else {
		return nil
	}
}

func (p *Proxy) sendCancel(client *definitions.FrontendClient, req *pgproto3.CancelRequest) error {
	if log.IsLevel(log.TraceLevel) {
		buf, err := json.Marshal(req)
		if err != nil {
			return err
		}
		log.Tracef("F -> %s", string(buf))
	}

	session, err := p.stateStore.GetSession(req.ProcessID)
	if err != nil {
		return err
	}

	err = p.stateStore.CancelQuery(session)
	if err != nil {
		return err
	}

	buf := (&pgproto3.NoticeResponse{Message: "Ok"}).Encode(nil)
	_, err = client.Conn.Write(buf)
	return err
}

func (p *Proxy) handlePass(client *definitions.FrontendClient, msg *pgproto3.StartupMessage) (bool, error) {
	if user, ok := msg.Parameters["user"]; !ok {
		buf := (&pgproto3.ErrorResponse{Message: "User not found!"}).Encode(nil)
		_, _ = client.Conn.Write(buf)
		return false, fmt.Errorf("User not found")
	} else {
		buf := (&pgproto3.AuthenticationCleartextPassword{}).Encode(nil)

		_, err := client.Conn.Write(buf)
		if err != nil {
			return false, fmt.Errorf("error sending pass request: %s", err)
		}

		msgFront, err := client.Backend.Receive()
		if err != nil {
			if err.Error() == "EOF" {
				return true, nil
			}
			return false, fmt.Errorf("error receiving password message: %s", err)
		}

		switch passMsg := msgFront.(type) {
		case *pgproto3.PasswordMessage:
			h := md5.New()
			_, _ = h.Write([]byte(passMsg.Password + user))
			hashPass := fmt.Sprintf("md5%x", string(h.Sum(nil)))

			var buf []byte

			pgconn, err := p.shazam.GetROConnection(context.Background())
			if err != nil {
				return false, err
			}

			defer pgconn.Release()

			if ok, msg := auth.ValidateUser(pgconn, user, hashPass); ok {
				buf = (&pgproto3.AuthenticationOk{}).Encode(buf)
				pid, secretKey := state.NewBackendKey(client.Conn)

				err = p.stateStore.SetSession(&state.Session{
					PID:    pid,
					Secret: secretKey,
				})
				if err != nil {
					return false, err
				}

				client.PID = pid
				client.SecretKey = secretKey

				buf = (&pgproto3.BackendKeyData{
					ProcessID: pid,
					SecretKey: secretKey,
				}).Encode(buf)

				buf = (&pgproto3.ReadyForQuery{TxStatus: 'I'}).Encode(buf)
			} else {
				buf = (&pgproto3.ErrorResponse{Message: msg}).Encode(buf)
				buf = (&pgproto3.Close{}).Encode(buf)

			}
			_, err = client.Conn.Write(buf)
			return false, err

		default:
			return false, fmt.Errorf("unknown pass message: %#v", msg)
		}
	}
}
