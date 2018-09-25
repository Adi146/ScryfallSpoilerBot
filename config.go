package main

import (
	"SkryfallSpoilerBot/Messenger"
	"encoding/json"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SetTypes      []string   `yaml:"set_types"`
	SetBlacklist  []string   `yaml:"set_blacklist"`
	CheckInterval int        `yaml:"check_interval"`
	Messengers    Messengers `yaml:"messengers"`
}

type Messengers struct {
	Pushbullets []*Messenger.Pushbullet `yaml:"pushbullet"`
}

func (config *Config) Init() {
	for _, pb := range config.Messengers.Pushbullets {
		pb.Init()
	}
}

func NewConfig(file string) (*Config, error) {
	var config Config
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return &config, err
	}

	err = yaml.Unmarshal(dat, &config)
	config.Init()
	return &config, nil
}

func (messengers *Messengers) Get() chan Messenger.Messenger {
	ch := make(chan Messenger.Messenger)
	go func() {
		for _, pb := range messengers.Pushbullets {
			ch <- pb
		}
		close(ch)
	}()
	return ch
}

func ReadSpoiledCards(file string) (*map[string][]string, error) {
	spoiledCards := make(map[string][]string)

	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return &spoiledCards, err
	}

	err = json.Unmarshal(dat, &spoiledCards)
	return &spoiledCards, err
}

func SaveSpoiledCards(file string, spoiledCards *map[string][]string) error {
	dat, err := json.Marshal(spoiledCards)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, dat, 0644)
}
