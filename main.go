package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	err := doCi()
	if err != nil {
		fmt.Println(err)
	}
}

func doCi() error {
	ctx := context.Background()

	// create a Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	// get the projects source directory
	src := client.Host().Workdir()

	// initialize new container from npm image
	container := client.Container().From("node")

  // load env variables
  container = container.WithEnvVariable("DOTENV_KEY", "dotenv://:key_828fe6f34bf06d690ca21eab39b2efc80e5466f412058cc156352ec0d5556a87@dotenv.org/vault/.env.vault?environment=ci")

	// mount source directory to /src
  container = container.WithMountedDirectory("/src", src).WithWorkdir("/src")

	// execute npm install
	container = container.Exec(dagger.ContainerExecOpts{
		Args: []string{"npm", "install"},
	})

	// execute build command
	container = container.Exec(dagger.ContainerExecOpts{
		Args: []string{"npm", "run", "build"},
	})

	// get build output
	build, err := container.Stdout().Contents(ctx)
	if err != nil {
		return err
	}
	// print output to console
	fmt.Println(build)

	return nil
}
