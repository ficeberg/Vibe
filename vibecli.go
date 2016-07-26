package main

import (
	"./vibe/controllers"
	"./vibe/models"
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"github.com/urfave/cli"
	"os"
	"strings"
)

func main() {

	app := cli.NewApp()
	app.Name = "vibecli"
	app.Version = "0.0.1"
	app.Usage = "Vibe command line interface"
	app.EnableBashCompletion = true

	var force bool
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "force, f",
			Usage:       "force option",
			Destination: &force,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "settings",
			Aliases: []string{"s"},
			Usage:   "Setting operations",
			Subcommands: []cli.Command{
				{
					Name:  "show",
					Usage: "show configurations",
					Action: func(c *cli.Context) error {
						conf := models.Config{}.Init()
						if c.Args().First() != "" {
							cmap := structs.Map(conf)
							for k, v := range cmap {
								if strings.ToLower(k) == strings.ToLower(c.Args().First()) {
									fmt.Printf("%+v\n", v)
								}
							}
						} else {
							fmt.Printf("%+v\n", conf)
						}

						return nil
					},
				},
				{
					Name:  "env",
					Usage: "show environment variables",
					Action: func(c *cli.Context) error {
						var env []string
						env = os.Environ()

						for index, value := range env {
							name := strings.Split(value, "=")
							fmt.Printf("[%d] %s : %v\n", index, name[0], name[1])
						}

						return nil
					},
				},
			},
		},
		{
			Name:    "user",
			Aliases: []string{"u"},
			Usage:   "User operations",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "add a new user.{email} {username} {password} {role}",
					Action: func(c *cli.Context) error {

						u := new(controllers.User)
						u.Email = c.Args().Get(0)
						u.Username = c.Args().Get(1)
						u.Password = c.Args().Get(2)
						u.Role = c.Args().Get(3)

						if err := u.Create(); err != nil {
							switch err.Error() {
							case "E11000":
								fmt.Println("user already exist.")
							default:
								fmt.Println(err.Error())
							}
						} else {
							fmt.Println("user " + c.Args().First() + " added")
							uj, _ := json.Marshal(u)
							fmt.Printf(string(uj))
						}
						return nil
					},
				},
				{
					Name:  "sudo",
					Usage: "Sudo to someone",
					Action: func(c *cli.Context) error {
						u := new(controllers.User)
						u.Username = c.Args().Get(0)
						token, err := u.GenerateToken("cli", "", -1)
						if err != nil {
							fmt.Println(err)
							return nil
						}
						fmt.Println(token)

						return nil
					},
				},
				{
					Name:  "show",
					Usage: "show user detail in JSON",
					Action: func(c *cli.Context) error {
						u := controllers.User{Username: c.Args().Get(0)}
						if err := u.Get(); err != nil {
							fmt.Printf("Invalid username " + c.Args().Get(0) + ". " + err.Error())
							return nil
						}
						uj, _ := json.Marshal(u)
						fmt.Printf(string(uj))
						return nil
					},
				},
				{
					Name:  "del",
					Usage: "delete an existing user by username",
					Action: func(c *cli.Context) error {
						u := new(controllers.User)
						u.Username = c.Args().Get(0)
						if err := u.Delete(); err != nil {
							fmt.Println(err)
							return nil
						}
						fmt.Println("user " + c.Args().First() + " deleted")
						return nil
					},
				},
				{
					Name:  "update",
					Usage: "update an existing user",
					Action: func(c *cli.Context) error {
						fmt.Println("update user: ", c.Args().First())
						return nil
					},
				},
				{
					Name:  "role",
					Usage: "change role for an existing user",
					Action: func(c *cli.Context) error {
						u := controllers.User{Username: c.Args().Get(0)}
						if err := u.Get(); err != nil {
							fmt.Printf("Invalid username " + c.Args().Get(0) + ". " + err.Error())
							return nil
						}
						ogRole := u.Role
						u.Role = c.Args().Get(1)
						if err := u.Update(); err != nil {
							fmt.Printf("Unale to update user " + c.Args().Get(0) + " with role " + c.Args().Get(1) + ". " + err.Error())
							return nil
						}
						if err := u.Get(); err != nil {
							fmt.Printf(err.Error())
							return nil
						}

						fmt.Println("User " + c.Args().First() + " role has change from " + ogRole + " to " + u.Role)
						return nil
					},
				},
				{
					Name:  "enable",
					Usage: "enable user",
					Action: func(c *cli.Context) error {
						u := controllers.User{Username: c.Args().Get(0), IsDisabled: false}
						if err := u.Update(); err != nil {
							fmt.Printf("Unale to update user " + c.Args().Get(0) + ". " + err.Error())
							return nil
						}

						fmt.Println("user " + c.Args().First() + " has been enabled")
						return nil
					},
				},
				{
					Name:  "disable",
					Usage: "disable user",
					Action: func(c *cli.Context) error {
						u := controllers.User{Username: c.Args().Get(0), IsDisabled: true}
						if err := u.Update(); err != nil {
							fmt.Printf("Unale to update user " + c.Args().Get(0) + ". " + err.Error())
							return nil
						}

						fmt.Println("user " + c.Args().First() + " has been disabled")
						return nil
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
