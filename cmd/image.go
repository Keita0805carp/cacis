package cmd

import (
  "fmt"

  "github.com/keita0805carp/cacis/cacis"

  "github.com/spf13/cobra"
)

var (
  dockerfileDir string
  imageRef       string

  imageCmd = &cobra.Command{
    Use: "image",
    Short: "Manage image",
    Run: imageCommand,
  }
)

func imageCommand(cmd *cobra.Command, args []string) {
  if err := imageAction(); err != nil {
    Exit(err, 1)
  }
}

func imageAction() (err error) {
  fmt.Println("image command")
  if dockerfileDir != "" && imageRef != "" {
    fmt.Println("from dockerfile")
    cacis.BuildFromDockerfile(dockerfileDir, imageRef)
  } else if imageRef != "" {
    fmt.Println("from imageRef")
    cacis.PullAndPush(imageRef)
  } else if dockerfileDir != "" && imageRef == "" {
    fmt.Println("Error: need imageRef")
    Exit(err, 1)
  } else {
    fmt.Println("Please set option")
    Exit(err, 1)
  }
  return nil
}

func init() {
  RootCmd.AddCommand(imageCmd)
  imageCmd.Flags().StringVarP(&dockerfileDir,  "dir", "d", "", "docker build(tag)->push")
  imageCmd.Flags().StringVarP(&imageRef, "image", "i", "", "docker pull->tag->push")
}
