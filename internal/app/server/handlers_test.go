package server

import (
	"database/sql"
	"net/http"
	"testing"

	"go.uber.org/zap"
	"project.com/internal/config"
	"project.com/internal/storage"
)

func Test_app_Ping(t *testing.T) {
	type fields struct {
		Storage *storage.MemStorage
		Logger  *zap.Logger
		Config  *config.ServerConfig
		DB      *sql.DB
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			mc := &app{
				Storage: tt.fields.Storage,
				Logger:  tt.fields.Logger,
				Config:  tt.fields.Config,
				DB:      tt.fields.DB,
			}
			mc.Ping(tt.args.w, tt.args.r)
		})
	}
}

func Test_app_HandlePostRequest(t *testing.T) {
	type fields struct {
		Storage *storage.MemStorage
		Logger  *zap.Logger
		Config  *config.ServerConfig
		DB      *sql.DB
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			mc := &app{
				Storage: tt.fields.Storage,
				Logger:  tt.fields.Logger,
				Config:  tt.fields.Config,
				DB:      tt.fields.DB,
			}
			mc.HandlePostRequest(tt.args.w, tt.args.r)
		})
	}
}
