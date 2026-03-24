package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"kubernetes/config"
	"kubernetes/dlogger"
	"kubernetes/helpers"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type rules struct {
	prefix string
	port   int
}

type PodList struct {
	Items []struct {
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
		Status struct {
			Phase string `json:"phase"`
		} `json:"status"`
	} `json:"items"`
}

func main() {
	// Init logger
	// TODO: the log file should be mapped based on how the pgm is called(by go run or as binary)
	logFile, err := dlogger.InitLogger("./data/kubernetes.log", slog.LevelInfo)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	//Edit the config file and exit
	edit := flag.Bool("edit", false, "Edit the config file")
	flag.Parse()
	if *edit {
		err := config.EditConfigFile()
		if err != nil {
			slog.Error("Error editing config file", "err", err)
		}
		return
	}

	// Init config
	config.Init()
	configData, err := config.ReadConfigFile()
	if err != nil {
		slog.Error("The config file does not have correct format")
		return
	}

	namespace := configData.Namespace
	portForwardRules := configData.PortForwardRules

	slog.Info(fmt.Sprintf("Starting automatic port-forwarding using namespace %s", namespace))

	kubeBinary := "kubectl"
	podArgs := []string{"get", "pods", "-n", namespace, "-o", "json"}
	slog.Debug(fmt.Sprint(podArgs))

	podsOut, err := helpers.RunCommand(kubeBinary, podArgs...)

	if err != nil {
		slog.Error("Unable to get list of pods")
		slog.Error(string(podsOut))
		return
	}

	var podList PodList
	err = json.Unmarshal(podsOut, &podList)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		slog.Info("Shutting Down")
		cancel()
	}()

	for _, rule := range portForwardRules {
		// Run in background
		for _, pod := range podList.Items {
			slog.Info(fmt.Sprintf("Pod: %s", pod.Metadata.Name))
			if strings.HasPrefix(pod.Metadata.Name, rule.Prefix) {
				pfArgs := []string{"port-forward", "-n", namespace, fmt.Sprintf("pod/%s", pod.Metadata.Name), fmt.Sprintf("%s:%s", rule.Port, rule.Port)}

				helpers.RunPortForward(ctx, kubeBinary, pfArgs...)

				slog.Info(fmt.Sprintf("Port-forwarding started for %s on port %s", pod.Metadata.Name, rule.Port))
			}
		}

		slog.Info(fmt.Sprintf("Forwarding local port %s to %s...\n", rule.Port, rule.Port))

	}

	<-ctx.Done()
	slog.Info("Script exited")
}
