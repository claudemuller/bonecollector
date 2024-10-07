package engine

import (
	"fmt"
	"log/slog"

	"github.com/veandco/go-sdl2/sdl"

	assetstore "sdl/pkg/asset-store"
	eventbus "sdl/pkg/event-bus"
)

type graphics struct {
	title    string
	width    int32
	height   int32
	window   *sdl.Window
	renderer *sdl.Renderer
}

type Engine struct {
	Running    bool
	graphics   graphics
	eventbus   eventbus.EventBus
	assetStore assetstore.AssetStore
}

const (
	fps               = 60
	millisecsPerFrame = 1000 / fps
)

var rect = sdl.Rect{50, 50, 200, 100}

func New(title string, width, height int32) (Engine, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return Engine{}, fmt.Errorf("error initialising SDL: %v", err)
	}

	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		return Engine{}, fmt.Errorf("error creating window: %v", err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return Engine{}, fmt.Errorf("error creating renderer: %v", err)
	}

	e := Engine{
		Running: true,
		graphics: graphics{
			title:    title,
			width:    width,
			height:   height,
			window:   window,
			renderer: renderer,
		},
		eventbus:   eventbus.New(),
		assetStore: assetstore.AssetStore{},
	}

	e.assetStore.AddTexture(e.graphics.renderer, "tilemap", "./assets/tilemaps/jungle.png")

	e.eventbus.On(eventbus.EventPlayerMoveHorz, func(args ...interface{}) {
		for _, v := range args {
			val, ok := v.(int32)
			if !ok {
				slog.Error("player move", "error", "failed to assert int32")
				return
			}
			rect.X += val
		}
	})
	e.eventbus.On(eventbus.EventPlayerMoveVert, func(args ...interface{}) {
		for _, v := range args {
			val, ok := v.(int32)
			if !ok {
				slog.Error("player move", "error", "failed to assert int32")
				return
			}
			rect.Y += val
		}
	})

	return e, nil
}

func (e *Engine) Run() {
	slog.Info("starting game")
	for e.Running {
		e.processInput()
		e.eventbus.Process()
		e.update()
		e.render()
	}
}

func (e *Engine) processInput() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch ev := event.(type) {
		case *sdl.QuitEvent:
			e.Running = false
			return

		case *sdl.KeyboardEvent:
			if ev.Keysym.Sym == sdl.K_ESCAPE {
				e.Running = false
				return
			}

			if ev.Keysym.Sym == sdl.K_UP {
				e.eventbus.Emit(eventbus.Event{
					Type: eventbus.EventPlayerMoveVert,
					Args: int32(-1),
				})
				return
			}
			if ev.Keysym.Sym == sdl.K_DOWN {
				e.eventbus.Emit(eventbus.Event{
					Type: eventbus.EventPlayerMoveVert,
					Args: int32(1),
				})
				return
			}
			if ev.Keysym.Sym == sdl.K_LEFT {
				e.eventbus.Emit(eventbus.Event{
					Type: eventbus.EventPlayerMoveHorz,
					Args: int32(-1),
				})
				return
			}
			if ev.Keysym.Sym == sdl.K_RIGHT {
				e.eventbus.Emit(eventbus.Event{
					Type: eventbus.EventPlayerMoveHorz,
					Args: int32(1),
				})
				return
			}
		}
	}
}

var prevFrameMS uint64

func (e *Engine) update() {
	timeToWait := millisecsPerFrame - sdl.GetTicks64() - prevFrameMS
	if timeToWait > 0 && timeToWait <= millisecsPerFrame {
		sdl.Delay(uint32(timeToWait))
	}

	_ = (sdl.GetTicks64() - prevFrameMS) / 1000.0
	prevFrameMS = sdl.GetTicks64()
}

func (e *Engine) render() {
	err := e.graphics.renderer.SetDrawColor(21, 21, 21, 255)
	if err != nil {
		slog.Error("render", "error", err)
	}
	err = e.graphics.renderer.Clear()
	if err != nil {
		slog.Error("render clear", "error", err)
	}

	err = e.graphics.renderer.SetDrawColor(255, 0, 0, 255)
	if err != nil {
		slog.Error("render", "error", err)
	}

	err = e.graphics.renderer.FillRect(&rect)
	if err != nil {
		slog.Error("render", "error", err)
	}

	e.graphics.renderer.Present()
}

func (e *Engine) Destroy() {
	err := e.graphics.renderer.Destroy()
	if err != nil {
		slog.Error("engine destroy", "error", err)
	}

	err = e.graphics.window.Destroy()
	if err != nil {
		slog.Error("engine destroy", "error", err)
	}

	sdl.Quit()
}
