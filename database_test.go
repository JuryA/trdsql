package trdsql

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func TestConnect(t *testing.T) {
	type args struct {
		driver string
		dsn    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "testSuccess",
			args:    args{driver: "sqlite3", dsn: ""},
			wantErr: false,
		},
		{
			name:    "testFail",
			args:    args{driver: "sqlite2", dsn: ""},
			wantErr: true,
		},
		{
			name:    "testPostgres",
			args:    args{driver: "postgres", dsn: "dbname=trdsql_test"},
			wantErr: false,
		},
		{
			name:    "testMysql",
			args:    args{driver: "mysql", dsn: "root@/trdsql_test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Connect(tt.args.driver, tt.args.dsn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDB_Disconnect(t *testing.T) {
	type args struct {
		driver string
		dsn    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "testSuccess",
			args:    args{driver: "sqlite3", dsn: ""},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Connect(tt.args.driver, tt.args.dsn)
			if err != nil {
				t.Fatal(err)
			}
			if err := db.Disconnect(); (err != nil) != tt.wantErr {
				t.Errorf("DB.Disconnect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_CreateTable(t *testing.T) {
	type fields struct {
		driver string
		dsn    string
	}
	type args struct {
		tableName   string
		names       []string
		types       []string
		isTemporary bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "testSuccess",
			fields: fields{driver: "sqlite3", dsn: ""},
			args: args{
				tableName:   "test",
				names:       []string{"a", "b"},
				types:       []string{"text", "text"},
				isTemporary: true,
			},
			wantErr: false,
		},
		{
			name:   "testSuccess2",
			fields: fields{driver: "sqlite3", dsn: ""},
			args: args{
				tableName:   "test",
				names:       []string{"c1"},
				types:       []string{"text"},
				isTemporary: false,
			},
			wantErr: false,
		},
		{
			name:   "testFail",
			fields: fields{driver: "sqlite3", dsn: ""},
			args: args{
				tableName:   "test",
				names:       []string{},
				types:       []string{},
				isTemporary: true,
			},
			wantErr: true,
		},
		{
			name:   "testFail2",
			fields: fields{driver: "sqlite3", dsn: ""},
			args: args{
				tableName:   "test",
				names:       []string{"c1"},
				types:       []string{},
				isTemporary: true,
			},
			wantErr: true,
		},
		{
			name:   "testFail3",
			fields: fields{driver: "sqlite3", dsn: ""},
			args: args{
				tableName:   "test",
				names:       []string{"c1"},
				types:       []string{},
				isTemporary: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Connect(tt.fields.driver, tt.fields.dsn)
			if err != nil {
				t.Fatal(err)
			}
			db.Tx, err = db.Begin()
			if err != nil {
				t.Fatal(err)
			}
			if err := db.CreateTable(tt.args.tableName, tt.args.names, tt.args.types, tt.args.isTemporary); (err != nil) != tt.wantErr {
				t.Errorf("DB.CreateTable() error = %v, wantErr %v", err, tt.wantErr)
			}
			err = db.Tx.Commit()
			if err != nil {
				t.Fatal(err)
			}
			err = db.Disconnect()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDB_Select(t *testing.T) {
	type fields struct {
		driver string
		dsn    string
	}
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "testErr",
			fields:  fields{driver: "sqlite3", dsn: ""},
			args:    args{query: ""},
			wantErr: true,
		},
		{
			name:    "testErr2",
			fields:  fields{driver: "sqlite3", dsn: ""},
			args:    args{query: "SELEC * FROM test"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Connect(tt.fields.driver, tt.fields.dsn)
			if err != nil {
				t.Fatal(err)
			}
			db.Tx, err = db.Begin()
			if err != nil {
				t.Fatal(err)
			}
			_, err = db.Select(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.Select() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			err = db.Tx.Commit()
			if err != nil {
				t.Fatal(err)
			}
			err = db.Disconnect()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDB_Import(t *testing.T) {
	type fields struct {
		driver string
		dsn    string
	}
	type args struct {
		tableName   string
		columnNames []string
		reader      Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "testErr",
			fields: fields{driver: "sqlite3", dsn: ""},
			args: args{
				tableName:   "test",
				columnNames: []string{"c1"},
				reader:      nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Connect(tt.fields.driver, tt.fields.dsn)
			if err != nil {
				t.Fatal(err)
			}
			db.Tx, err = db.Begin()
			if err != nil {
				t.Fatal(err)
			}
			if err := db.Import(tt.args.tableName, tt.args.columnNames, tt.args.reader); (err != nil) != tt.wantErr {
				t.Errorf("DB.Import() error = %v, wantErr %v", err, tt.wantErr)
			}
			err = db.Tx.Commit()
			if err != nil {
				t.Fatal(err)
			}
			err = db.Disconnect()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDB_EscapeName(t *testing.T) {
	type fields struct {
		escape string
	}
	type args struct {
		oldName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "test1",
			fields: fields{escape: "`"},
			args:   args{oldName: "test"},
			want:   "`test`",
		},
		{
			name:   "test2",
			fields: fields{escape: "\""},
			args:   args{oldName: "test"},
			want:   "\"test\"",
		},
		{
			name:   "test3",
			fields: fields{escape: "`"},
			args:   args{oldName: "`test`"},
			want:   "`test`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				escape: tt.fields.escape,
			}
			if got := db.EscapeName(tt.args.oldName); got != tt.want {
				t.Errorf("DB.EscapeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_RewriteSQL(t *testing.T) {
	type fields struct {
		rewritten []string
	}
	type args struct {
		query   string
		oldName string
		newName string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantRewrite string
	}{
		{
			name:        "test1",
			fields:      fields{rewritten: []string{}},
			args:        args{query: "SELECT * FROM test", oldName: "test", newName: "`test`"},
			wantRewrite: "SELECT * FROM `test`",
		},
		{
			name:        "test2",
			fields:      fields{rewritten: []string{"`test`"}},
			args:        args{query: "SELECT * FROM `test`", oldName: "test", newName: "`test`"},
			wantRewrite: "SELECT * FROM `test`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				rewritten: tt.fields.rewritten,
			}
			if gotRewrite := db.RewriteSQL(tt.args.query, tt.args.oldName, tt.args.newName); gotRewrite != tt.wantRewrite {
				t.Errorf("DB.RewriteSQL() = %v, want %v", gotRewrite, tt.wantRewrite)
			}
		})
	}
}
