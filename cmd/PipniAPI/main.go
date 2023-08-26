package main

import (
	"os"

	gameui "github.com/mokiat/PipniAPI/internal/ui"
	"github.com/mokiat/PipniAPI/resources"
	glapp "github.com/mokiat/lacking-gl/app"
	glui "github.com/mokiat/lacking-gl/ui"
	"github.com/mokiat/lacking/log"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/util/resource"
)

func main() {
	log.Info("Started")
	if err := runApplication(); err != nil {
		log.Error("Crashed: %v", err)
		os.Exit(1)
	}
	log.Info("Stopped")
}

func runApplication() error {
	locator := ui.WrappedLocator(resource.NewFSLocator(resources.UI))

	uiController := ui.NewController(
		locator,
		glui.NewShaderCollection(),
		gameui.BootstrapApplication,
	)

	cfg := glapp.NewConfig("Pipni API", 1024, 768)
	cfg.SetMaximized(true)
	cfg.SetMinSize(800, 600)
	cfg.SetVSync(true)
	cfg.SetIcon("images/icon.png")
	cfg.SetLocator(locator)
	return glapp.Run(cfg, uiController)
}
