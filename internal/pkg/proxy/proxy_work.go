package proxy

import (
	"context"
	"encoding/json"

	"github.com/electric-saw/pg-shazam/internal/pkg/backend"
	"github.com/electric-saw/pg-shazam/internal/pkg/definitions"
	"github.com/electric-saw/pg-shazam/internal/pkg/log"
	"github.com/electric-saw/pg-shazam/internal/pkg/parser"
	"github.com/electric-saw/pg-shazam/internal/pkg/state"
	"github.com/jackc/pgproto3/v2"
)

func (p *Proxy) DeleteSession(pid int64) {
	err := p.stateStore.DeleteSession(pid)
	if err != nil {
		log.Warnf("Fail on delete session %d: %s", pid, err)
	}
}

func (p *Proxy) handleMessages(client *definitions.FrontendClient) error {
	go client.ReadClient(context.Background())
	client.ReadNext()

	for {
		log.Debugf("Handling messages of %s", client.Conn.RemoteAddr().String())
		select {
		case msg := <-client.MsgChan:
			switch msg.(type) {
			case *pgproto3.Terminate:
				return nil
			default:
				err := p.redirectMessage(client, msg)
				if err != nil {
					return err
				}
				client.ReadNext()
			}
		}
	}
}

func (p *Proxy) redirectMessage(client *definitions.FrontendClient, raw pgproto3.FrontendMessage) (err error) {
	switch msg := raw.(type) {
	case *pgproto3.Query:
		qry := parser.ParseQuery(msg.String)

		if qry.Operation == parser.Set {
			conn, err := p.shazam.GetROConnection(context.Background())
			if err != nil {
				return err
			}
			_, err = conn.Exec(context.Background(), msg.String)
			if err != nil {
				return err
			}

			client.SetCommands = append(client.SetCommands, qry.QueryString)
			buf := (&pgproto3.CommandComplete{CommandTag: []byte("SET")}).Encode(nil)
			buf = (&pgproto3.ReadyForQuery{TxStatus: 'I'}).Encode(buf)
			_, err = client.Conn.Write(buf)
			return err
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if qry.DDLOperation {
			errs := p.shazam.RunAllPrimaryHosts(qry.QueryString)
			if len(errs) > 0 {
				log.Errorf("Fail to run '%s' on all primary %v", msg.String, errs)
			}

			if qry.Shards != nil {
				errs := backend.InsertCatalogShardDefinition(p.shazam, qry.TableName, qry.Shards)
				if len(errs) > 0 {
					log.Errorf("Fail to insert shard definition %v", errs)
				}
			}

			// TODO: Alimentar CurrDatabase
			err = p.stateStore.SetHashSet(&state.HashSet{
				Database: client.CurrDatabase,
				Table:    qry.TableName,
				Fields:   qry.Shards,
			})

			if err != nil {
				log.Errorf("Fail on local store hashset cache %v", err)
			}

			buf := (&pgproto3.CommandComplete{CommandTag: []byte("SET")}).Encode(nil)
			buf = (&pgproto3.ReadyForQuery{TxStatus: 'I'}).Encode(buf)
			_, _ = client.Conn.Write(buf)
		} else {
			var hs *state.HashSet

			if len(qry.TableName) > 0 {
				hs, err = p.stateStore.GetHashSet(client.CurrDatabase, qry.TableName)
				if err != nil {
					log.Warnf("Failed on get hashset of table %s: %v", qry.TableName, err)
				}
			}

			var cluster *backend.Cluster
			if hs != nil {
				cluster = p.shazam.ClusterByHash(&qry, hs.Fields)
			} else {
				cluster = p.shazam.GetRandomCluster()
			}

			var conn *backend.Conn

			switch qry.Operation {
			case parser.Select:
				conn, err = cluster.GetROConnection(ctx)
				if err != nil {
					return err
				}
			default:
				conn, err = cluster.GetRWConnection(ctx)
				if err != nil {
					return err
				}

			}

			err = conn.AssumeClient(ctx, client, msg)
			if err != nil {
				return err
			}

			conn.Release()
		}

	case *definitions.Error:
		return msg.Err

	default:
		buf, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		log.Errorf("Message type %T not supported: %s", msg, string(buf))
	}

	return nil

}
