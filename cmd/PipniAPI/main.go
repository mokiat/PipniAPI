package main

import (
	"log/slog"
	"os"

	"github.com/mokiat/PipniAPI/internal"
	"github.com/mokiat/PipniAPI/resources"
	glapp "github.com/mokiat/lacking-native/app"
	glui "github.com/mokiat/lacking-native/ui"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/util/resource"
)

func main() {
	slog.Info("Started")
	if err := runApplication(); err != nil {
		slog.Error("Crashed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("Stopped")
}

func runApplication() error {
	locator := ui.WrappedLocator(resource.NewFSLocator(resources.UI))

	uiController := ui.NewController(
		locator,
		glui.NewShaderCollection(),
		internal.BootstrapApplication,
	)

	cfg := glapp.NewConfig("Pipni API", 1024, 768)
	cfg.SetMaximized(true)
	cfg.SetMinSize(1024, 768)
	cfg.SetVSync(true)
	cfg.SetIcon("images/icon.png")
	cfg.SetLocator(locator)
	return glapp.Run(cfg, uiController)
}
