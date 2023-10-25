package internal

import "testing"

func TestSendDataToServer(t *testing.T) {
	type args struct {
		metrics   []*Metrics
		serverURL string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendDataToServer(tt.args.metrics, tt.args.serverURL); (err != nil) != tt.wantErr {
				t.Errorf("SendDataToServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
