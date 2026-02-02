package services

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type Config struct {
	Id          int
	T           *testing.T
	Seq         *atomic.Int64
	Options     application.ServiceOptions
	StartupErr  bool
	ShutdownErr bool
}

func Configure[T any, P interface {
	*T
	Configure(Config)
}](srv P, c Config) application.Service {
	srv.Configure(c)
	return application.NewServiceWithOptions(srv, c.Options)
}

type Error struct {
	Id int
}

func (e *Error) Error() string {
	return fmt.Sprintf("service #%d mock failure", e.Id)
}

type Startupper struct {
	Config
	startup int64
}

func (s *Startupper) Configure(c Config) {
	s.Config = c
}

func (s *Startupper) Id() int {
	return s.Config.Id
}

func (s *Startupper) StartupSeq() int64 {
	return s.startup
}

func (s *Startupper) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	if s.startup != 0 {
		s.T.Errorf("Double startup for service #%d: first at seq=%d, then at seq=%d", s.Id(), s.startup, s.Seq.Load())
		return nil
	}

	s.startup = s.Seq.Add(1)

	if diff := cmp.Diff(s.Options, options); diff != "" {
		s.T.Errorf("Options mismatch for service #%d (-want +got):\n%s", s.Id(), diff)
	}

	if s.StartupErr {
		return &Error{Id: s.Id()}
	} else {
		return nil
	}
}

type Shutdowner struct {
	Config
	shutdown int64
}

func (s *Shutdowner) Configure(c Config) {
	s.Config = c
}

func (s *Shutdowner) Id() int {
	return s.Config.Id
}

func (s *Shutdowner) ShutdownSeq() int64 {
	return s.shutdown
}

func (s *Shutdowner) ServiceShutdown() error {
	if s.shutdown != 0 {
		s.T.Errorf("Double shutdown for service #%d: first at seq=%d, then at seq=%d", s.Id(), s.shutdown, s.Seq.Load())
		return nil
	}

	s.shutdown = s.Seq.Add(1)

	if s.ShutdownErr {
		return &Error{Id: s.Id()}
	} else {
		return nil
	}
}

type StartupShutdowner struct {
	Config
	startup  int64
	shutdown int64
	ctx      context.Context
}

func (s *StartupShutdowner) Configure(c Config) {
	s.Config = c
}

func (s *StartupShutdowner) Id() int {
	return s.Config.Id
}

func (s *StartupShutdowner) StartupSeq() int64 {
	return s.startup
}

func (s *StartupShutdowner) ShutdownSeq() int64 {
	return s.shutdown
}

func (s *StartupShutdowner) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	if s.startup != 0 {
		s.T.Errorf("Double startup for service #%d: first at seq=%d, then at seq=%d", s.Id(), s.startup, s.Seq.Load())
		return nil
	}

	s.startup = s.Seq.Add(1)
	s.ctx = ctx

	if diff := cmp.Diff(s.Options, options); diff != "" {
		s.T.Errorf("Options mismatch for service #%d (-want +got):\n%s", s.Id(), diff)
	}

	if s.StartupErr {
		return &Error{Id: s.Id()}
	} else {
		return nil
	}
}

func (s *StartupShutdowner) ServiceShutdown() error {
	if s.shutdown != 0 {
		s.T.Errorf("Double shutdown for service #%d: first at seq=%d, then at seq=%d", s.Id(), s.shutdown, s.Seq.Load())
		return nil
	}

	s.shutdown = s.Seq.Add(1)

	select {
	case <-s.ctx.Done():
	default:
		s.T.Errorf("Service #%d shut down before context cancellation", s.Id())
	}

	if s.ShutdownErr {
		return &Error{Id: s.Id()}
	} else {
		return nil
	}
}
