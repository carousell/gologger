package logger

import (
	"os"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogger_Infod(t *testing.T) {
	log := MustNamed("test")
	defer log.Infod()("Test log infod should take about 3s")
	time.Sleep(3 * time.Second)
}

func TestWarnd(t *testing.T) {
	log := MustNew()

	t.Run("log duration", func(t *testing.T) {
		defer log.Warnd(10 * time.Second)("this should not be displayed at all")
		time.Sleep(300 * time.Millisecond)
	})

	t.Run("log warn", func(t *testing.T) {
		defer log.Warnd(10 * time.Millisecond)("this should be displayed in warn level")
		time.Sleep(30 * time.Millisecond)
	})

	t.Run("log auto", func(t *testing.T) {
		defer log.Autod(50 * time.Millisecond)("this should be logged in warn")
		defer log.Autod(300 * time.Millisecond)("this should be logged in debug")
		time.Sleep(100 * time.Millisecond)
	})
}

func TestCaller(t *testing.T) {
	os.Setenv("LOG_LEVEL", "info")
	log := MustNamed("test-caller")
	defer log.Sync()
	log.Info("1")
	log.Warn("2")

	uw := log.Unwrap()
	uw.Info("3")

	log2 := Wrap(uw)
	log2.Info("4")

	zaplog, _ := NewZap("test-caller")
	zaplog.Info("5")

	w := Wrap(zaplog)
	w.Warn("6")
}

func TestHook(t *testing.T) {
	os.Setenv("LOG_LEVEL", "info")
	log := MustNamed("test-hook").AddHook(func(e zapcore.Entry) error {
		if e.Level >= zap.WarnLevel {
			t.Logf("warn ne: %s - %s", e.Level, e.Message)
		}
		return nil
	})
	defer log.Sync()
	log.Info("1")
	log.Warn("2")

	uw := log.Unwrap()
	uw.Warn("3")
	uw.Info("4")
}

func Benchmark_ZapRawMessage(b *testing.B) {
	os.Setenv("LOG_LEVEL", "debug")
	log := MustNamed("benlogger").Unwrap()
	defer log.Sync()
	for i := 0; i < b.N; i++ {
		log.Infof("benchmarking")
	}
}

func Benchmark_LoggerRawMessage(b *testing.B) {
	os.Setenv("LOG_LEVEL", "debug")
	log := MustNamed("benlogger")
	defer log.Sync()
	for i := 0; i < b.N; i++ {
		log.Infof("benchmarking")
	}
}

func Benchmark_ZapLogNoOutput(b *testing.B) {
	os.Setenv("LOG_LEVEL", "warn")
	log := MustNamed("benzap").Unwrap()
	defer log.Sync()
	prefix := "run"
	for i := 0; i < b.N; i++ {
		log.Infof("%s benchmarking", prefix)
	}
}

func Benchmark_LoggerNoOutput(b *testing.B) {
	os.Setenv("LOG_LEVEL", "warn")
	log := MustNamed("benlogger").With("run")
	defer log.Sync()
	for i := 0; i < b.N; i++ {
		log.Infof("benchmarking")
	}
}

func Benchmark_ZapWithArguments(b *testing.B) {
	os.Setenv("LOG_LEVEL", "warn")
	log := MustNamed("benzap").Unwrap()
	defer log.Sync()
	prefix := "run"
	for i := 0; i < b.N; i++ {
		log.Infof("%s benchmarking: %d - %d", prefix, b.N, i)
	}
}

func Benchmark_LoggerWithArguments(b *testing.B) {
	os.Setenv("LOG_LEVEL", "warn")
	log := MustNamed("benlogger").With("run")
	defer log.Sync()
	for i := 0; i < b.N; i++ {
		log.Infof("benchmarking: %d - %d", b.N, i)
	}
}

func Benchmark_ZapToConsole(b *testing.B) {
	os.Setenv("LOG_LEVEL", "info")
	log := MustNamed("benzap").Unwrap()
	defer log.Sync()
	prefix := "run"
	for i := 0; i < b.N; i++ {
		log.Infof("%s benchmarking: %d - %d", prefix, b.N, i)
	}
}

func Benchmark_LoggerToConsole(b *testing.B) {
	os.Setenv("LOG_LEVEL", "info")
	log := MustNamed("benlogger").With("run")
	defer log.Sync()
	for i := 0; i < b.N; i++ {
		log.Infof("benchmarking: %d - %d", b.N, i)
	}
}
