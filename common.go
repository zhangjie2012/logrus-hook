package logrushook

import (
	"fmt"
	"net"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack"
)

// LogWashFunc log struct to bytes
type LogWashFunc func(appName string, t time.Time, metadata logrus.Fields, caller *runtime.Frame, level logrus.Level, message string) []byte

type DefaultLogS struct {
	Time     time.Time     `msgpack:"@time"`
	MetaData logrus.Fields `msgpack:"@metadata"` // map[string]interface{}
	Ip       string        `msgpack:"@ip"`
	Level    string        `msgpack:"@level"`
	Caller   string        `msgpack:"@caller"`
	Message  string        `msgpack:"@message"`
}

// https://stackoverflow.com/a/37382208/802815
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

var instanceIp string

func getIp() string {
	if instanceIp == "" {
		instanceIp = getOutboundIP()
	}
	return instanceIp
}

// DefaultLogWashFunc log filter and serialization
//   - use msgpack serialize log
//   - level > 'DEBUG'
//   - set current instance ip
func DefaultLogWashFunc(appName string, t time.Time, metadata logrus.Fields, caller *runtime.Frame, level logrus.Level, message string) []byte {
	if level > logrus.DebugLevel {
		return nil
	}

	caller_ := ""
	if caller != nil {
		caller_ = fmt.Sprintf("%s:%d", filepath.Base(caller.File), caller.Line)
	}

	l := DefaultLogS{
		Time:     t,
		MetaData: metadata,
		Ip:       getIp(),
		Level:    level.String(),
		Caller:   caller_,
		Message:  message,
	}

	bs, err := msgpack.Marshal(&l)
	if err != nil {
		return nil
	}

	return bs
}
