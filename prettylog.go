package prettylog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/Jamlie/colors"
)

const (
	timeFormat     = "[15:04:05]"
	greyColor      = 245
	greenColor     = 40
	orangeColor    = 202
	redColor       = 196
	purplishColor  = 57
	yellowishColor = 178
)

var (
	messageColor = colors.New(colors.WhiteFg)
	attrsColor   = colors.NewCustomId(greyColor)

	infoColor    = colors.NewCustomId(greenColor)
	warningColor = colors.NewCustomId(orangeColor)
	errorColor   = colors.NewCustomId(redColor)
	debugColor   = colors.NewCustomId(purplishColor)

	timeColor = colors.NewCustomId(yellowishColor)
)

type Handler struct {
	h slog.Handler
	b *bytes.Buffer
	m *sync.Mutex
}

func NewHandler(opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	b := bytes.NewBuffer(nil)
	return &Handler{
		h: slog.NewJSONHandler(b, &slog.HandlerOptions{
			AddSource:   opts.AddSource,
			Level:       opts.Level,
			ReplaceAttr: suppressDefaults(opts.ReplaceAttr),
		}),
		b: b,
		m: &sync.Mutex{},
	}
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{
		h: h.h.WithAttrs(attrs),
		b: h.b,
		m: h.m,
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		h: h.h.WithGroup(name),
		b: h.b,
		m: h.m,
	}
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = debugColor.String(level)
	case slog.LevelInfo:
		level = infoColor.String(level)
	case slog.LevelWarn:
		level = warningColor.String(level)
	case slog.LevelError:
		level = errorColor.String(level)
	}

	attrs, err := h.computeAttrs(ctx, r)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(attrs, "", "  ")
	if err != nil {
		return fmt.Errorf("error when marshaling attrs: %w", err)
	}

	fmt.Println(
		timeColor.String(r.Time.Format(timeFormat)),
		level,
		messageColor.String(r.Message),
		attrsColor.String(string(bytes)),
	)

	return nil
}

func (h *Handler) computeAttrs(
	ctx context.Context,
	r slog.Record,
) (map[string]any, error) {
	h.m.Lock()
	defer func() {
		h.b.Reset()
		h.m.Unlock()
	}()

	if err := h.h.Handle(ctx, r); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle: %w", err)
	}

	var attrs map[string]any
	err := json.Unmarshal(h.b.Bytes(), &attrs)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshaling inner handler's Handle result: %w", err)
	}

	return attrs, nil
}

func suppressDefaults(
	next func([]string, slog.Attr) slog.Attr,
) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey ||
			a.Key == slog.LevelKey ||
			a.Key == slog.MessageKey {
			return slog.Attr{}
		}

		if next == nil {
			return a
		}

		return next(groups, a)
	}
}
