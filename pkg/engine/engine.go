package engine

import (
	"fmt"
	"log/slog"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

	assetstore "bonecollector/pkg/asset-store"
	"bonecollector/pkg/ecs"
	"bonecollector/pkg/entity"
	eventbus "bonecollector/pkg/event-bus"
	"bonecollector/pkg/system"
)

type graphics struct {
	title    string
	width    int32
	height   int32
	camera   sdl.Rect
	window   *sdl.Window
	renderer *sdl.Renderer
}

type Engine struct {
	graphics      graphics
	eventbus      eventbus.EventBus
	assetStore    *assetstore.AssetStore
	entityManager ecs.EntityManager
	isRunning     bool
	isDebug       bool
	isPaused      bool
	mapWidth      int
	mapHeight     int
}

const (
	fps               = 60
	millisecsPerFrame = 1000 / fps
)

func New(title string, winWidth, winHeight int32) (Engine, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return Engine{}, fmt.Errorf("error initialising SDL: %v", err)
	}

	if err := ttf.Init(); err != nil {
		return Engine{}, fmt.Errorf("error initialising TTF: %v", err)
	}

	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		return Engine{}, fmt.Errorf("error creating window: %v", err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return Engine{}, fmt.Errorf("error creating renderer: %v", err)
	}

	e := Engine{
		isRunning: true,
		graphics: graphics{
			title:    title,
			width:    winWidth,
			height:   winHeight,
			window:   window,
			renderer: renderer,
			camera: sdl.Rect{
				W: winWidth,
				H: winHeight,
			},
		},
		eventbus:      eventbus.New(),
		assetStore:    assetstore.New(),
		entityManager: ecs.New(),
	}

	return e, nil
}

func (e *Engine) Run() {
	slog.Info("starting game")

	e.setup()

	for e.isRunning {
		e.processInput()
		e.eventbus.Process()
		e.update()
		e.render()
	}
}

func (e *Engine) setup() {
	// Load textures
	if err := e.assetStore.AddTexture(e.graphics.renderer, "test", "assets/tank.png"); err != nil {
		slog.Error("error loading texture", "asset-id", "test", "path", "assets/tank.png")
	}

	// Load Fonts
	// if err := e.assetStore.AddFont( "test", "assets/tank.png", 20); err != nil {
	// 	slog.Error("error loading font", "asset-id", "test", "path", "assets/tank.png");
	// }

	// Add movement system
	e.entityManager.AddSystem(system.MovementSystem)
	e.entityManager.AddSystem(system.RenderSystem)

	// Load entities
	var ent entity.Entity
	e.entityManager.TagEntity(&ent, "tank")
	e.entityManager.GroupEntity(&ent, "player")
	e.entityManager.AddEntityTransformComponent(&ent, 50, 50)

	// e.assetStore.AddTexture(e.graphics.renderer, "tilemap", "./assets/tilemaps/jungle.png")
	//
	// e.eventbus.On(eventbus.EventPlayerMoveHorz, func(args ...interface{}) {
	// 	for _, v := range args {
	// 		val, ok := v.(int32)
	// 		if !ok {
	// 			slog.Error("player move", "error", "failed to assert int32")
	// 			return
	// 		}
	// 		rect.X += val
	// 	}
	// })
	// e.eventbus.On(eventbus.EventPlayerMoveVert, func(args ...interface{}) {
	// 	for _, v := range args {
	// 		val, ok := v.(int32)
	// 		if !ok {
	// 			slog.Error("player move", "error", "failed to assert int32")
	// 			return
	// 		}
	// 		rect.Y += val
	// 	}
	// })

	// Inin Lua

	// Load level..
}

func (e *Engine) processInput() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch ev := event.(type) {
		case *sdl.QuitEvent:
			e.isRunning = false

		case *sdl.KeyboardEvent:
			if ev.Keysym.Sym == sdl.K_ESCAPE {
				e.isRunning = false
			}

			if ev.Keysym.Sym == sdl.K_F1 {
				if e.isDebug {
					e.isDebug = false
				} else {
					e.isDebug = true
				}
			}

			// if ev.Keysym.Sym == sdl.K_KP_P {
			// 	e.isPaused != e.isPaused
			// }

			e.eventbus.Emit(eventbus.Event{
				Type: eventbus.EventPlayerMove,
				Args: ev.Keysym.Sym,
			})
		}
	}
}

var prevFrameMS uint64

func (e *Engine) update() {
	timeToWait := millisecsPerFrame - sdl.GetTicks64() - prevFrameMS
	if timeToWait > 0 && timeToWait <= millisecsPerFrame {
		sdl.Delay(uint32(timeToWait))
	}

	delta := (sdl.GetTicks64() - prevFrameMS) / 1000.0
	prevFrameMS = sdl.GetTicks64()

	if e.isPaused {
		return
	}

	// e.entityManager.GetSystem().SubscribeToEvents(e.eventbus)

	e.entityManager.Update(delta)
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

	e.entityManager.Render(e.graphics.renderer, e.graphics.camera, e.assetStore)

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

	ttf.Quit()
	sdl.Quit()
}
