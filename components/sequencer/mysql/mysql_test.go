package mysql

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/pkg/log"
	"testing"
)

const (
	MySQLUrl = "localhost:xxxxx"
	Value    = 1
	Key      = "sequenceKey"
	Size     = 50
)

func TestMySQLSequencer_Init(t *testing.T) {

	var MySQLUrl = MySQLUrl

	comp := NewMySQLSequencer(log.DefaultLogger)

	cfg := sequencer.Configuration{
		BiggerThan: nil,
		Properties: make(map[string]string),
	}

	cfg.Properties["MySQLHost"] = MySQLUrl

	err := comp.Init(cfg)

	assert.Error(t, err)
}

func TestMySQLSequencer_GetNextId(t *testing.T) {

	//var MySQLUrl = "localhost:3306"
	//
	comp := NewMySQLSequencer(log.DefaultLogger)

	//cfg := sequencer.Configuration{
	//	BiggerThan: nil,
	//	Properties: make(map[string]string),
	//}

	//cfg.Properties["MySQLHost"] = MySQLUrl
	//_ = comp.Init(cfg)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"sequencer_key", "sequencer_value"}).AddRow([]driver.Value{Key, Value}...)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	properties := make(map[string]string)

	req := &sequencer.GetNextIdRequest{Key: Key, Options: sequencer.SequencerOptions{AutoIncrement: sequencer.STRONG}, Metadata: properties}

	_, err1 := comp.GetNextId(req, db)
	if err1 != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NoError(t, err1)
}

func TestMySQLSequencer_GetSegment(t *testing.T) {

	comp := NewMySQLSequencer(log.DefaultLogger)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"sequencer_key", "sequencer_value"}).AddRow([]driver.Value{Key, Value}...)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	properties := make(map[string]string)

	req := &sequencer.GetSegmentRequest{Size: Size, Key: Key, Options: sequencer.SequencerOptions{AutoIncrement: sequencer.STRONG}, Metadata: properties}

	_, _, err1 := comp.GetSegment(req, db)
	if err1 != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NoError(t, err1)
}

func TestMySQLSequencer_Close(t *testing.T) {

	var MySQLUrl = MySQLUrl

	comp := NewMySQLSequencer(log.DefaultLogger)

	cfg := sequencer.Configuration{
		BiggerThan: nil,
		Properties: make(map[string]string),
	}

	cfg.Properties["MySQLHost"] = MySQLUrl

	_ = comp.Init(cfg)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	mock.ExpectBegin()
	mock.ExpectCommit()

	comp.Close(db)
}
