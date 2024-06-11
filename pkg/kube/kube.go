package kube

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type Context struct {
	Name string
}

type Client interface {
	Service(context Context, namespace string, name string) (Service, error)
	Forwarder(context Context, cservice Service, options ForwarderOptions) (Forwarder, error)
}

type client struct {
	cliendcmdConfig *clientcmdapi.Config
	k8sClients      map[string]*k8sClient
}

type k8sClient struct {
	clientset    *kubernetes.Clientset
	clientConfig clientcmd.ClientConfig
}

func NewClient() (Client, error) {
	var err error
	c := &client{
		k8sClients: make(map[string]*k8sClient),
	}

	if envKubeConfig := os.Getenv("KUBECONFIG"); envKubeConfig != "" {
		c.cliendcmdConfig = clientcmdapi.NewConfig()
		for _, kubeConfigPath := range strings.Split(envKubeConfig, ":") {
			cfg, err := clientcmd.LoadFromFile(kubeConfigPath)
			if err != nil {
				log.Printf("cannot load kube config from file %s: %s\n", kubeConfigPath, err.Error())
				continue
			}

			for key, cluster := range cfg.Clusters {
				if _, exists := c.cliendcmdConfig.Clusters[key]; exists {
					log.Printf("found cluster %s in %s, but a cluster with the same name was already loaded from another file, skipping\n", key, kubeConfigPath)
				} else {
					c.cliendcmdConfig.Clusters[key] = cluster
				}
			}

			for key, ctx := range cfg.Contexts {
				if _, exists := c.cliendcmdConfig.Contexts[key]; exists {
					log.Printf("found context %s in %s, but a context with the same name was already loaded from another file, skipping\n", key, kubeConfigPath)
				} else {
					c.cliendcmdConfig.Contexts[key] = ctx
				}
			}

			for key, authInfo := range cfg.AuthInfos {
				if _, exists := c.cliendcmdConfig.AuthInfos[key]; exists {
					log.Printf("found authInfo %s in %s, but a context with the same name was already loaded from another file, skipping\n", key, kubeConfigPath)
				} else {
					c.cliendcmdConfig.AuthInfos[key] = authInfo
				}
			}

			if cfg.CurrentContext != "" {
				c.cliendcmdConfig.CurrentContext = cfg.CurrentContext
			}
		}
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal("cannot find user's home dir, hence cannot find kubeconfig:", err)
		}
		kubeConfigPath := filepath.Join(home, ".kube", "config")
		c.cliendcmdConfig, err = clientcmd.LoadFromFile(kubeConfigPath)
		if err != nil {
			log.Fatalf("cannot load kubeconfig from default location %s: %s\n", kubeConfigPath, err)
		}
	}

	c.k8sClients[""], err = c.getKubeClient("") // current
	if err != nil {
		log.Println("error creating current k8s clientset")
		log.Fatalln(err)
	}

	return c, nil
}

func (c *client) getKubeClient(name string) (*k8sClient, error) {
	if cs, found := c.k8sClients[name]; found {
		return cs, nil
	}

	override := &clientcmd.ConfigOverrides{}
	if name != "" {
		override.CurrentContext = name
	}

	clientconfig := clientcmd.NewDefaultClientConfig(*c.cliendcmdConfig, override)
	restConfig, err := clientconfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot create rest config for context '%s': %w", name, err)
	}

	clientConfigset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot create client config from rest config for context '%s': %w", name, err)
	}

	client := &k8sClient{
		clientset:    clientConfigset,
		clientConfig: clientconfig,
	}

	return client, err
}
