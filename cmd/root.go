// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	vaultApi "github.com/hashicorp/vault/api"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uniseraph/go-spring/info"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

var cfgFile string

var logFile string

var logConfig zap.Config

var rootConfig struct{
	confPrefix  string
	serviceName string
	//logPath     string

	logFile     string
	grpcHealthProbeEnabled bool
	grpcHealthProbeHome string
	grpcHealthProbeBinary string
	consulClient *api.Client
}

func LogConfig() *zap.Config {
	return &logConfig
}
func ServiceName() string {
	return rootConfig.serviceName
}
func ConsulAddr() string {
	return viper.GetString("consul.address")
}
func EnableGrpcHealthProbe() bool{
	return rootConfig.grpcHealthProbeEnabled
}
func GrpcHealthProbeHome() string{
	return rootConfig.grpcHealthProbeHome
}
func GrpcHealthProbeBinary() string{
	return rootConfig.grpcHealthProbeBinary
}
func ConsulClient () *api.Client {
	return rootConfig.consulClient
}
// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "grpc-push",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRun: func(cmd *cobra.Command, args []string) {


		rootConfig.confPrefix = viper.GetString("conf.prefix")
		rootConfig.serviceName = viper.GetString("service.name")
		rootConfig.grpcHealthProbeHome = viper.GetString("grpc.health.probe.home")
		rootConfig.grpcHealthProbeEnabled = viper.GetBool("grpc.health.probe.enabled")
		rootConfig.grpcHealthProbeBinary = viper.GetString("grpc.health.probe.binary")

		config := api.DefaultConfig()
		config.Address = ConsulAddr()

		client, err := api.NewClient(config)
		if err != nil {
			panic(err)
		}
		rootConfig.consulClient = client
		if viper.GetBool("consul.enabled")==true {
			var configPaths []string

			if viper.GetBool("consul.context.enabled") {
				configPaths= []string{
					rootConfig.confPrefix + "/" + viper.GetString("consul.context"),
					rootConfig.confPrefix + "/" + rootConfig.serviceName}
			}else{
				configPaths = []string{
					rootConfig.confPrefix + "/" + rootConfig.serviceName}
			}


			if err := parseConsulConfig(client, configPaths); err != nil {
				panic(err)
			}
		}


		if viper.GetBool("vault.enabled") {
			vaultClient, err := vaultApi.NewClient(&vaultApi.Config{
				Address: viper.GetString("vault.address"),
			})

			if err != nil {
				panic(err)
			}
			vaultClient.SetToken(viper.GetString("vault.token"))
			if err := parseVaultConfig(vaultClient, []string{"secret/application",
				"secret/"+ rootConfig.serviceName}); err != nil {
				panic(err)
			}
		}



		logConfig= zap.NewProductionConfig()
		logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		logConfig.DisableStacktrace = true
		logConfig.OutputPaths=[]string{viper.GetString("log.file")}
		logConfig.ErrorOutputPaths=[]string{viper.GetString("log.file")}
		logConfig.Encoding = viper.GetString("log.encoding")
		
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/."+info.Target+".yaml)")
	RootCmd.PersistentFlags().String("conf.prefix", "config", "consul config prefix ,default: config")
	RootCmd.PersistentFlags().String("service.name", info.Target, "service name in consul")
	RootCmd.PersistentFlags().String("log.file", "/dev/stdout", "log file, default is stdout")
	RootCmd.PersistentFlags().String("log.encoding", "console", "json/console")

	RootCmd.PersistentFlags().String("consul.address", "127.0.0.1:8500", "consul server addr")
	RootCmd.PersistentFlags().Bool("grpc.health.probe.enabled", false, "enable grpc-health-probe ")
	RootCmd.PersistentFlags().String("grpc.health.probe.home", "~/go/bin", "grpc-health-probe home")
	RootCmd.PersistentFlags().String("grpc.health.probe.binary", "grpc_health_probe", "grpc-health-probe binary name")

	RootCmd.PersistentFlags().Bool("consul.enabled", true, "enable consul ")
	RootCmd.PersistentFlags().Bool("consul.context.enabled", true, "enable consul config context ")
	RootCmd.PersistentFlags().String("consul.context", "application", "consul config context ")
	RootCmd.PersistentFlags().Bool("vault.enabled", false, "enable vault ")
	RootCmd.PersistentFlags().String("vault.address", "http://127.0.0.1:8200", "vault server addr")
	RootCmd.PersistentFlags().String("vault.token", "root", "root")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	viper.BindPFlags(RootCmd.PersistentFlags())
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".grpc-push" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName("."+info.Target)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
func parseConsulConfig(client *api.Client, prefixs []string) error {

	viper.SetConfigType("json")
	m := make(map[string]string, 100)
	for _, key := range prefixs {
		pairs, _, err := client.KV().List(key, &api.QueryOptions{})

		if err != nil {
			return err
		}

		for _, pair := range pairs {
			m[pair.Key[len(key)+1:]] = string(pair.Value)
		}
	}

	buf, _ := json.Marshal(m)
	viper.MergeConfig(strings.NewReader(string(buf)))

	return nil

}

func parseVaultConfig(client *vaultApi.Client, prefixs []string) error {

	viper.SetConfigType("json")
	m := make(map[string]string, 100)
	for _, key := range prefixs {

		secret, err := client.Logical().Read(key)
		if err != nil {
			return err
		}

		for k, v := range secret.Data {
			logrus.WithField("key", k).WithField("value", v).Info("put into the viper from vault")
			m[k] = v.(string)
		}
	}

	buf, _ := json.Marshal(m)
	viper.MergeConfig(strings.NewReader(string(buf)))

	return nil

}
