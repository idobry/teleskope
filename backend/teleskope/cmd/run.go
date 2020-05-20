/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/idobry/teleskope/backend/teleskope/controller"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func serve() {
	h := controller.NewHub()
	go h.Run()

	r := mux.NewRouter()
	corsObj:=handlers.AllowedOrigins([]string{"*"})
	r.HandleFunc("/list/ns", func(w http.ResponseWriter, r *http.Request) {
		controller.GetNamespaces(h, w, r)
	}).Methods("GET")
	r.HandleFunc("/list3/dep/{ns}", controller.GetDeployments).Methods("GET")
	r.HandleFunc("/dep/{ns}/{dep}", controller.GetDeployment).Methods("GET")
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		controller.StreamUpdateds(h, w, r)
	})

	go controller.StreamDeployments(h)

	fmt.Println("ListenAndServe...")
	err := http.ListenAndServe(":" + os.Getenv("PORT"), handlers.CORS(corsObj)(r))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func init() {
	rootCmd.AddCommand(runCmd)
}
