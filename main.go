package main

import (
	"os"
	"log"
	"os/exec"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"webhook/model"
	"webhook/helper"
)

var config model.Config
var configFile string

func main() {
	args := os.Args

	//if we have a "real" argument we take this as conf path to the config file
	//you can specify the config.json on the first argument when running this app via bash
	if (len(args) > 1) {
		configFile = args[1]
	} else {
		configFile = "config.json"
	}

	//load config
	config = model.LoadConfig(configFile)

	//open log file
	writer, err := os.OpenFile(config.Logfile, os.O_RDWR | os.O_APPEND | os.O_CREATE, 0666)
	helper.PanicIf(err)

	//close logfile on exit
	defer func() {
		writer.Close()
	}()

	//setting logging output
	log.SetOutput(writer)

	//setting handler
	http.HandleFunc("/", hookHandler)

	address := config.Host + ":" + strconv.FormatInt(config.Port, 10)

	log.Println("Listening on " + address)

	//starting server
	err = http.ListenAndServe(address, nil)
	if (err != nil) {
		log.Println(err)
	}
}

func hookHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	event := r.Header.Get("X-Gitlab-Event")
	var hook model.Webhook

	//read request body
	var data, err = ioutil.ReadAll(r.Body)
	helper.PanicIf(err, "while reading request")

	//unmarshal request body
	err = json.Unmarshal(data, &hook)
	helper.PanicIf(err, "while unmarshaling request")

	log.Println("----------------------------------------")
	log.Printf("Incoming Request From Repository %s", hook.Repository.Name)

	for _, repo := range config.Repositories {
		// match repo url and event
		if repo.Url == hook.Repository.GitSshUrl && (repo.Event == hook.ObjectKind || repo.Event == hook.EventName) {
			log.Printf("Received %s event", event)
			log.Printf("Action: %s", repo.Event)

			for _, filter := range repo.Filters {
				// match the ref event
				if filter.Ref == hook.Ref {
					log.Printf("reff: %s", hook.Ref)

					// run before deploy
					runCommand(config.Deploy.Before, repo, filter)

					log.Println("Command: ")
					for _, cmd := range repo.Commands {
						runCommand(cmd, repo, filter)
					}

					// run after deploy
					runCommand(config.Deploy.After, repo, filter)
				}
			}
		}
	}

	log.Printf("End of Request")
	log.Printf("----------------------------------------")
	log.Printf("")

	json.NewEncoder(w).Encode("OK")
}

func runCommand(cmd string, repo model.ConfigRepository, filter model.ConfigFilter) {
	if cmd != "" {
		var command = exec.Command(cmd, filter.Path, repo.Url, filter.Branch)
		out, err := command.Output()
		if (err != nil) {
			log.Println(err)
		} else {
			log.Printf("- Executed: %s", cmd)
			log.Printf("- Output: %s", string(out))
		}
	}
}