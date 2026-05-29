package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/Zyko0/go-sdl3/sdl"

	"vds.io/goml/api"
	"vds.io/goml/siren"
)

type Config struct {
	Siren  string
	Volume float32
}

func main() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		userConfigDir = "/etc"
	}
	configPath := path.Join(userConfigDir, "goml.toml")

	log.Printf("Reading config %s\n", configPath)
	appConfig := Config{
		Siren:  "siren.wav",
		Volume: 5.0,
	}
	toml.DecodeFile(configPath, &appConfig)
	log.Printf("Config is %v\n", appConfig)

	defer binsdl.Load().Unload()
	defer sdl.Quit()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		log.Fatalf(err.Error())
	}

	alarm, err := siren.New(appConfig.Siren, appConfig.Volume)
	if err != nil {
		log.Fatalf(err.Error())
	}

	apiServer := api.NewServer(alarm)
	rc, err := apiServer.StartListening()
	if err != nil {
		log.Fatalf(err.Error())
	}

	select {
	case s := <-sigChan:
		fmt.Println("OS Signal", s)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		apiServer.Stop(ctx)
	case listenResult := <-rc:
		fmt.Println("API Server Listen Result", listenResult)
	}
}
