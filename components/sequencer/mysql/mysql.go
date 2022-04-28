package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"mosn.io/layotto/components/pkg/utils"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/pkg/log"
)

type MySQLSequencer struct {
	metadata   utils.MySQLMetadata
	biggerThan map[string]int64
	logger     log.ErrorLogger
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewMySQLSequencer(logger log.ErrorLogger) *MySQLSequencer {
	s := &MySQLSequencer{
		logger: logger,
	}

	return s
}

func (e *MySQLSequencer) Init(config sequencer.Configuration) error {

	db := &sql.DB{}

	m, err := utils.ParseMySQLMetadata(config.Properties, db)

	if err != nil {
		return err
	}
	e.metadata = m
	e.biggerThan = config.BiggerThan

	if err = utils.NewMySQLClient(e.metadata); err != nil {
		return err
	}
	e.ctx, e.cancel = context.WithCancel(context.Background())

	if len(e.biggerThan) > 0 {
		var Key string
		var Value int64
		for k, bt := range e.biggerThan {
			if bt <= 0 {
				continue
			} else {
				m.Db.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE sequencer_key = ?", m.TableName), k).Scan(&Key, &Value)
				if Value < bt {
					return fmt.Errorf("MySQL sequencer error: can not satisfy biggerThan guarantee.key: %s,key in MySQL: %s", k, Key)
				}
			}
		}
	}
	return nil
}

func (e *MySQLSequencer) GetNextId(req *sequencer.GetNextIdRequest, db *sql.DB) (*sequencer.GetNextIdResponse, error) {

	metadata, err := utils.ParseMySQLMetadata(req.Metadata, db)

	var Key string
	var Value int64
	if err != nil {
		return nil, err
	}

	begin, err := metadata.Db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		switch err {
		case nil:
			err = begin.Commit()
		default:
			begin.Rollback()
		}
	}()

	err = begin.QueryRow("SELECT sequencer_key,sequencer_value FROM ? WHERE sequencer_key = ?", metadata.TableName, req.Key).Scan(&Key, &Value)
	if err != nil {
		return nil, err
	}
	Value += 1

	_, err1 := begin.Exec("UPDATE ? SET sequencer_value = ? WHERE sequencer_key = ?", metadata.TableName, Value, req.Key)
	if err1 != nil {
		return nil, err1
	}

	if _, err2 := begin.Exec("INSERT INTO ?(sequencer_key, sequencer_value) VALUES(?,?)", metadata.TableName, req.Key, Value+1); err2 != nil {
		return nil, err2
	}

	//if Value == 0 {
	//	Value += 1
	//	//insert, inerr := metadata.Db.Exec(fmt.Sprintf(`INSERT INTO %s VALUES(?,?)`, metadata.TableName), req.Key, Value)
	//
	//	if _, err1 := begin.Exec("UPDATE ? SET sequencer_value = ? WHERE sequencer_key = ?", "test", Value, Key); err1 != nil {
	//		return nil, err1
	//	}
	//
	//	//if inerr != nil {
	//	//	fmt.Println("insert into table fail", inerr.Error())
	//	//	return nil, inerr
	//	//}
	//	////fmt.Sprint("insert success", insert)
	//	//
	//	//if uperr != nil {
	//	//	fmt.Println("fetch row affected failed:", uperr.Error())
	//	//	return nil, uperr
	//	//}
	//	//fmt.Println("update recors number", update)
	//} else {
	//	Value += 1
	//	if _, err2 := begin.Exec("INSERT INTO test(sequencer_key, sequencer_value) VALUES(?,?)", "test1", Value); err2 != nil {
	//		return nil, err2
	//	}
	//	////update, uperr := metadata.Db.Exec(fmt.Sprintf(`UPDATE %s SET sequencer_value = ? WHERE sequencer_key = ?`, metadata.TableName), Value, req.Key)
	//	//update, uperr := metadata.Db.Exec(`UPDATE test SET sequencer_value = ? WHERE sequencer_key = ?`, 2, "test1")
	//	//if uperr != nil {
	//	//	fmt.Println("fetch row affected failed:", uperr.Error())
	//	//	return nil, uperr
	//	//}
	//	//fmt.Println("update recors number", update)
	//}

	e.Close(metadata.Db)

	return &sequencer.GetNextIdResponse{
		NextId: Value,
	}, nil
}

func (e *MySQLSequencer) GetSegment(req *sequencer.GetSegmentRequest, db *sql.DB) (support bool, result *sequencer.GetSegmentResponse, err error) {

	if req.Size == 0 {
		return true, nil, nil
	}

	metadata, err := utils.ParseMySQLMetadata(req.Metadata, db)

	var Key string
	var Value int64
	if err != nil {
		return false, nil, err
	}

	begin, err := metadata.Db.Begin()
	if err != nil {
		return false, nil, err
	}

	defer func() {
		switch err {
		case nil:
			err = begin.Commit()
		default:
			begin.Rollback()
		}
	}()

	err = begin.QueryRow("SELECT sequencer_key,sequencer_value FROM ? WHERE sequencer_key = ?", metadata.TableName, req.Key).Scan(&Key, &Value)
	if err != nil {
		return false, nil, err
	}
	Value += int64(req.Size)

	_, err1 := begin.Exec("UPDATE ? SET sequencer_value = ? WHERE sequencer_key = ?", metadata.TableName, Value, req.Key)
	if err1 != nil {
		return false, nil, err1
	}

	if _, err2 := begin.Exec("INSERT INTO ?(sequencer_key, sequencer_value) VALUES(?,?)", metadata.TableName, req.Key, Value+1); err2 != nil {
		return false, nil, err2
	}

	e.Close(metadata.Db)

	return false, &sequencer.GetSegmentResponse{
		From: Value - int64(req.Size) + 1,
		To:   Value,
	}, nil
}

func (e *MySQLSequencer) Close(db *sql.DB) error {

	return db.Close()
}
