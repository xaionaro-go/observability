package observability_test

import (
	"bytes"
	"fmt"
	"runtime"
	"testing"

	"github.com/facebookincubator/go-belt/tool/logger"
	xlogrus "github.com/facebookincubator/go-belt/tool/logger/implementation/logrus"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/xaionaro-go/observability"
)

func TestSecretsFilter(t *testing.T) {
	ll := xlogrus.DefaultLogrusLogger()
	ll.Formatter.(*logrus.TextFormatter).TimestampFormat = "unit-test"
	ll.Formatter.(*logrus.TextFormatter).CallerPrettyfier = func(f *runtime.Frame) (function string, file string) {
		return "", ""
	}
	l := xlogrus.New(ll).WithLevel(logger.LevelTrace).WithPreHooks(
		observability.StructFieldSecretsFilter{},
	)

	t.Run("error", func(t *testing.T) {
		e0 := fmt.Errorf("some error")
		e1 := fmt.Errorf("the parent error: %w", e0)
		var buf bytes.Buffer
		ll.SetOutput(&buf)
		l.Debugf("%v", e1)
		require.Equal(t, `time=unit-test level=debug msg="the parent error: some error"`+"\n", buf.String())
	})
}
