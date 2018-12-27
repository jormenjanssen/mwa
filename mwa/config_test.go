package main

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"
	"time"
)


func TestJsonConfigFromReader(t *testing.T) {

	jsonOk,cfgOk := InMemoryJson()
	jsonF,_ := InMemoryCorruptPropertiesJson()
	jsonC,_  := InMemoryCorruptedJsonFile()

	cfgDefault := Config{}


	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{name: "DeserializationOk", args: args{r: jsonOk}, want:cfgOk, wantErr:false},
		{name: "DeserializationValidationFailure", args: args{r: jsonF}, want:cfgDefault, wantErr:true},
		{name: "DeserializationJsonFailure", args: args{r: jsonC}, want:cfgDefault, wantErr:true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JsonConfigFromReader(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonConfigFromReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonConfigFromReader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func InMemoryJson() (io.Reader, Config) {

	d := 60 * time.Second

	bs := bytes.NewBufferString("")
	cfg := Config{
		Host: "192.168.10.240",
		RecoveryTime: d.String(),
	    AutoIpv4GatewayHost: true,
	    MonitorOnly:true}

	fmt.Fprintln(bs, "{")
	fmt.Fprintln(bs, fmt.Sprintf("\"RecoveryTime\":\"%v\",", d.String()))
	fmt.Fprintln(bs, fmt.Sprintf("\"AutoIpv4GatewayHost\":%v,", cfg.AutoIpv4GatewayHost))
	fmt.Fprintln(bs, fmt.Sprintf("\"MonitorOnly\":%v,", cfg.MonitorOnly))
	fmt.Fprintln(bs, fmt.Sprintf("\"Host\":\"%v\"", cfg.Host))
	fmt.Fprintln(bs, "}")


	return bs, cfg
}

func InMemoryCorruptPropertiesJson() (io.Reader, Config) {

	d := 60 * time.Second

	bs := bytes.NewBufferString("")
	cfg := Config{
		Host: "192.168.10.240",
		RecoveryTime: d.String(),
		AutoIpv4GatewayHost: true,
		MonitorOnly:true}

	fmt.Fprintln(bs, "{")
	fmt.Fprintln(bs, fmt.Sprintf("\"Recoveryttime\":\"%v\",", d.String()))
	fmt.Fprintln(bs, fmt.Sprintf("\"A5utoIpv44GatewayHost\":%v,", cfg.AutoIpv4GatewayHost))
	fmt.Fprintln(bs, fmt.Sprintf("\"MonitorrOnly\":%v,", cfg.MonitorOnly))
	fmt.Fprintln(bs, fmt.Sprintf("\"Hostt\":\"%v\"", cfg.Host))
	fmt.Fprintln(bs, "}")


	return bs, cfg
}

func InMemoryCorruptedJsonFile() (io.Reader, Config) {

	d := 60 * time.Second

	bs := bytes.NewBufferString("")
	cfg := Config{
		Host: "192.168.10.240",
		RecoveryTime: d.String(),
		AutoIpv4GatewayHost: true,
		MonitorOnly:true}

	fmt.Fprintln(bs, "{{")
	fmt.Fprintln(bs, fmt.Sprintf("\"Recoveryttime\":\"%v\",,,,", d.String()))
	fmt.Fprintln(bs, fmt.Sprintf("\"A5utoIpv44GatewayHost\":%v,,,", cfg.AutoIpv4GatewayHost))
	fmt.Fprintln(bs, fmt.Sprintf("\"MonitorrOnly\"::::%v,", cfg.MonitorOnly))
	fmt.Fprintln(bs, fmt.Sprintf("\"Hostt\":\"%v\"", cfg.Host))
	fmt.Fprintln(bs, ",}")


	return bs, cfg
}


