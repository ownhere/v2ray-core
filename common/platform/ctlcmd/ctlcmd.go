package ctlcmd

import (
	"io"
	"os"

	"v2ray.com/core/common/buf"

	"github.com/gogo/protobuf/proto"
	"v2ray.com/ext/tools/conf/serial"
	"bytes"
)

//go:generate errorgen

func Run(args []string, input io.Reader) (buf.MultiBuffer, error) {
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

	outBuffer, err := buf.ReadAllToMultiBuffer(stdoutReader)

	return outBuffer.MultiBuffer, nil
}
