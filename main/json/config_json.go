package json

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg json -path Main,Json

import (
	"context"
	"io"
	"os"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/platform"
	"v2ray.com/core/common/signal"

	"github.com/gogo/protobuf/proto"
	"v2ray.com/ext/tools/conf/serial"
	"bytes"
)

type logWriter struct{}

func (*logWriter) Write(b []byte) (int, error) {
	n, err := os.Stderr.Write(b)
	if err == nil {
		os.Stderr.WriteString(platform.LineSeparator())
	}
	return n, err
}

func jsonToProto(input io.Reader) (*core.Config, error) {
	pbConfig, err := serial.LoadJSONConfig(input)
	if err != nil {
		os.Stderr.WriteString("failed to parse json config: " + err.Error())
		os.Exit(-1)
	}

	bytesConfig, err := proto.Marshal(pbConfig)
	if err != nil {
		os.Stderr.WriteString("failed to marshal proto config: " + err.Error())
		os.Exit(-1)
	}

	stdoutReader := bytes.NewReader(bytesConfig)

	var config *core.Config

	loadTask := signal.ExecuteAsync(func() error {
		c, err := core.LoadConfig("protobuf", "", stdoutReader)
		if err != nil {
			return err
		}
		config = c
		return nil
	})

	if err := signal.ErrorOrFinish1(context.Background(), loadTask); err != nil {
		return nil, err
	}

	return config, nil
}

func init() {
	common.Must(core.RegisterConfigLoader(&core.ConfigFormat{
		Name:      "JSON",
		Extension: []string{"json"},
		Loader: func(input io.Reader) (*core.Config, error) {
			config, err := jsonToProto(input)
			if err != nil {
				return nil, newError("failed to execute v2ctl to convert config file.").Base(err).AtWarning()
			}
			return config, nil
		},
	}))
}
