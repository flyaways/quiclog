package quiclog

import "testing"

func Test_writer_Write(t *testing.T) {
	type fields struct {
		url string
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantN   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &writer{
				url: tt.fields.url,
			}
			gotN, err := w.Write(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("writer.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("writer.Write() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}
