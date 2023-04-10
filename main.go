package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var (
	DefaultClashUrl = "https://w.x.y.z/config.yaml"
)

type ClashProxy struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Server   string `yaml:"server"`
	Port     int    `yaml:"port"`
	Cipher   string `yaml:"cipher"`
	Password string `yaml:"password"`
	Udp      bool   `yaml:"udp"`
}

type ClashConfig struct {
	Proxies []*ClashProxy `yaml:"proxies"`
}

type ShadowsocksServer struct {
	Method   string `json:"method"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Server   string `json:"server"`
}

type ShadowsocksConfig struct {
	Locals  []interface{}        `json:"locals"`
	Servers []*ShadowsocksServer `json:"servers"`
}

func main() {

	var (
		clashConfig     ClashConfig
		clashConfigData []byte
		clashConfigUrl  string
		err             error
		outputFile      string
		resp            *http.Response
		ssConfig        ShadowsocksConfig
		ssConfigData    []byte
		ssConfigFile    string
		ssService       string
	)

	flag.StringVar(&clashConfigUrl, "c", DefaultClashUrl, "clash config url")
	flag.StringVar(&ssConfigFile, "s", "shadowsocks.json", "shadowsocks config")
	flag.StringVar(&outputFile, "o", "merged.json", "output config")
	flag.StringVar(&ssService, "service", "shadowsocks@client.service", "shadowsocks service")
	flag.Parse()

	// Read Clash config from network URL
	log.Println("[clash] Downloading config ...")
	if resp, err = http.Get(clashConfigUrl); err != nil {
		log.Println("Failed to read Clash config from URL")
		panic(err)
		return
	}
	defer resp.Body.Close()

	log.Println("[clash] Reading config ...")
	if clashConfigData, err = io.ReadAll(resp.Body); err != nil {
		log.Println("Failed to read Clash config response body")
		panic(err)
		return
	}

	// Parse Clash config YAML
	log.Println("[clash] Parsing config ...")
	if err = yaml.Unmarshal(clashConfigData, &clashConfig); err != nil {
		log.Println("Failed to parse Clash config YAML")
		panic(err)
		return
	}

	// Read Shadowsocks config from local file
	log.Println("[shadowsocks] Reading config ...")
	if ssConfigData, err = os.ReadFile(ssConfigFile); err != nil {
		log.Println("Failed to read Shadowsocks config file")
		panic(err)
		return
	}

	// Parse Shadowsocks config JSON
	log.Println("[shadowsocks] Parsing config ...")
	if err = json.Unmarshal(ssConfigData, &ssConfig); err != nil {
		log.Println("Failed to parse Shadowsocks config JSON")
		panic(err)
		return
	}

	// Merge Clash proxies to Shadowsocks servers
	ssConfig.Servers = make([]*ShadowsocksServer, 0, len(clashConfig.Proxies))
	for _, proxy := range clashConfig.Proxies {
		if proxy.Type == "ss" {
			log.Println("[shadowsocks] add outbound server:", proxy.Name)
			server := &ShadowsocksServer{
				Method:   proxy.Cipher,
				Password: proxy.Password,
				Port:     proxy.Port,
				Server:   proxy.Server,
			}
			ssConfig.Servers = append(ssConfig.Servers, server)
		}
	}

	// Serialize Shadowsocks config JSON
	log.Println("[shadowsocks] serializing new config ...")
	ssConfigData, err = json.MarshalIndent(ssConfig, "", "    ")
	if err != nil {
		log.Println("Failed to serialize Shadowsocks config JSON")
		panic(err)
		return
	}

	// Write Shadowsocks config to file
	log.Println("[shadowsocks] writting new config ...")
	if err = os.WriteFile(outputFile, ssConfigData, 0644); err != nil {
		log.Println("Failed to write Shadowsocks config file")
		panic(err)
		return
	}

	// Restart Shadowsocks service using systemctl command
	log.Println("Restarting service ...")
	if err = exec.Command("systemctl", "restart", ssService).Run(); err != nil {
		log.Println("Error restarting Shadowsocks service")
		panic(err)
		return
	}

	fmt.Println("Successfully updated Shadowsocks config with Clash proxies!")
}
