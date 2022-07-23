package cfg

import (
	"os"
	"testing"
)

type Env struct {
	Key string
	Val string
}

type Config struct {
	String string  `env:"str"`
	Int    int     `env:"int"`
	Bool   bool    `env:"bool"`
	Int64  int64   `env:"int64"`
	Config Config2 `env:"-"`
}

type Config2 struct {
	Int   *int `env:"int2"`
	Bool  bool `env:"bool2"`
	Skip  bool
	Skip2 bool `env:"-"`
}

func TestLoadFromEnv(t *testing.T) {
	type args struct {
		config any
	}

	tests := []struct {
		name    string
		args    args
		Envs    []Env
		wantErr bool
	}{
		{
			name:    "test correct string env",
			args:    args{config: &Config{}},
			Envs:    []Env{{Key: "str", Val: "test"}},
			wantErr: false,
		},
		{
			name:    "test correct int env",
			args:    args{config: &Config{}},
			Envs:    []Env{{Key: "ini", Val: "123"}},
			wantErr: false,
		},
		{
			name:    "test correct bool env",
			args:    args{config: &Config{}},
			Envs:    []Env{{Key: "bool", Val: "true"}},
			wantErr: false,
		},
		{
			name:    "test wrong int env",
			args:    args{config: &Config{}},
			Envs:    []Env{{Key: "int", Val: "text"}},
			wantErr: true,
		},
		{
			name:    "test wrong bool env",
			args:    args{config: &Config{}},
			Envs:    []Env{{Key: "bool", Val: "truue"}},
			wantErr: true,
		},
		{
			name:    "test wrong type int64 env",
			args:    args{config: &Config{}},
			Envs:    []Env{{Key: "int64", Val: "2112"}},
			wantErr: true,
		},
		{
			name:    "test point type env",
			args:    args{config: &Config{Config: Config2{}}},
			Envs:    []Env{{Key: "int2", Val: "123"}},
			wantErr: false,
		},
		{
			name:    "test error in child struct env",
			args:    args{config: &Config{Config: Config2{}}},
			Envs:    []Env{{Key: "int2", Val: "text"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				for _, env := range tt.Envs {
					defer os.Unsetenv(env.Key)
					_ = os.Setenv(env.Key, env.Val)
				}
				if err := LoadFromEnv(tt.args.config); (err != nil) != tt.wantErr {
					t.Errorf("LoadFromEnv() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}
