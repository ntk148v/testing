package main

import (
	"fmt"
	"log"

	"github.com/cristalhq/aconfig"
)

type MyConfig struct {
	HTTPPort int `default:"1111" usage:"just a number"`
	Auth     struct {
		User string `default:"def-user" usage:"your user"`
		Pass string `default:"def-pass" usage:"make it strong"`
	}
}

func main() {
	var cfg MyConfig
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		SkipDefaults: true,
		SkipFiles:    true,
		SkipEnv:      true,
		SkipFlags:    true,
		Files:        []string{"/var/opt/myapp/config.json"},
		EnvPrefix:    "APP",
		FlagPrefix:   "app",
	})

	if err := loader.Load(); err != nil {
		log.Panic(err)
	}

	fmt.Println("* Simple:")
	fmt.Printf("HTTPPort:  %v\n", cfg.HTTPPort)
	fmt.Printf("Auth.User: %q\n", cfg.Auth.User)
	fmt.Printf("Auth.Pass: %q\n", cfg.Auth.Pass)

	fmt.Println("* WalkFields")
	loader = aconfig.LoaderFor(&cfg, aconfig.Config{
		SkipFiles: true,
		SkipEnv:   true,
		SkipFlags: true,
	})

	loader.WalkFields(func(f aconfig.Field) bool {
		fmt.Printf("%v: %q %q %q %q\n", f.Name(), f.Tag("env"), f.Tag("flag"), f.Tag("default"), f.Tag("usage"))
		return true
	})

	if err := loader.Load(); err != nil {
		log.Panic(err)
	}

	fmt.Println("* Default")
	fmt.Printf("HTTPPort:  %v\n", cfg.HTTPPort)
	fmt.Printf("Auth.User: %v\n", cfg.Auth.User)
	fmt.Printf("Auth.Pass: %v\n", cfg.Auth.Pass)
}
