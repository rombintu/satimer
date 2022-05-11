package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gen2brain/beeep"
	daemon "github.com/sevlyar/go-daemon"
)

type Satimer struct {
	Message  string
	Timer    *time.Timer
	Daemoner *daemon.Context
}

func NewSatimer(message string) *Satimer {
	return &Satimer{
		Message: message,
		Daemoner: &daemon.Context{
			PidFileName: "satimer.pid",
			PidFilePerm: 0644,
			WorkDir:     ".",
			Umask:       027,
		},
	}
}

func (s *Satimer) Notify() error {
	if err := beeep.Notify(
		"Satimer",
		s.Message,
		"/usr/share/icons/Papirus/24x24/apps/alarm-timer.svg",
	); err != nil {
		return err
	}
	return nil
}

func (s *Satimer) StartTimer() {
	<-s.Timer.C
	if err := s.Notify(); err != nil {
		log.Fatalf("v+%", err)
	}
	fmt.Printf("Satimer: %s", s.Message)
}

func (s *Satimer) Daemon() {
	d, err := s.Daemoner.Reborn()
	if err != nil {
		log.Fatalf("v+%", err)
	}
	if d != nil {
		return
	}
	defer s.Daemoner.Release()

	s.StartTimer()
}

func main() {
	n := flag.String("n", "", "Time duration ex: 10h10m10s, 5s, 12m")
	message := flag.String("message", "Time's up!", "Message for notify")
	flag.Parse()
	if *n == "" {
		log.Println(flag.ErrHelp.Error())
		os.Exit(0)
	}
	parseDuration, err := time.ParseDuration(*n)
	if err != nil {
		log.Fatalf("v+%", err)
	}
	st := NewSatimer(*message)
	st.Timer = time.NewTimer(parseDuration)
	st.Daemon()
}
