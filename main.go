package main

import (
	"io/ioutil"
	"log"
	"sort"
	"strings"

	"github.com/spf13/viper"
)

type mailbox struct {
	Name     string
	Priority int
}

func main() {
	viper.SetConfigName("mbsort")
	viper.AddConfigPath("$HOME/.mutt")

	viper.SetDefault("input", "$HOME/.mutt/mailboxes_raw")
	viper.SetDefault("output", "$HOME/.mutt/mailboxes")
	viper.SetDefault("priorities", []string{})
	viper.SetDefault("defaultPrioriy", 1000)
	viper.SetDefault("debug", true)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Unable to read config file: %s", err)
	}

	// read reused configuration
	debug := viper.GetBool("debug")
	priorities := viper.GetStringSlice("priorities")
	defaultPriority := viper.GetInt("defaultPriority")
	inputFile, err := expandPath(viper.GetString("input"))
	if err != nil {
		log.Fatalf("Unable to expand input file: %s", err)
	}
	outputFile, err := expandPath(viper.GetString("output"))
	if err != nil {
		log.Fatalf("Unable to expand output file: %s", err)
	}

	if debug {
		log.Printf("Viper is using file %s", viper.ConfigFileUsed())
		log.Printf("Priorities are: %s", strings.Join(priorities, " "))

		viper.Debug()
	}

	// read file and parse it
	fr, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %s", err)
	}

	// mm hold mailboxes and its priorityes
	mm := []*mailbox{}

	for _, i := range strings.Split(string(fr), " ") {
		mb := &mailbox{
			Name:     i,
			Priority: getPosition(priorities, strings.Replace(i, `"`, "", -1)),
		}

		// default metric
		if mb.Priority == 0 {
			mb.Priority = defaultPriority
		}

		if debug {
			log.Printf("Adding %s with priority %d", mb.Name, mb.Priority)
		}
		mm = append(mm, mb)
	}

	// sort boxes
	sort.SliceStable(mm, func(i int, j int) bool {
		return mm[i].Priority > mm[j].Priority && mm[i].Name != "mailboxes"
	})

	// export to srt
	st := []string{}
	for _, i := range mm {
		st = append(st, i.Name)
	}

	err = ioutil.WriteFile(outputFile, []byte(strings.Join(st, " ")), 0644)
	if err != nil {
		log.Fatalf("Failed writing output file: %s", err)
	}
}
