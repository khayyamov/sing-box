package main

import (
	"context"
	"github.com/sagernet/sing-box/api/constant"
	"os"
	"os/user"
	"strconv"
	"time"

	_ "github.com/sagernet/sing-box/include"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing/service/filemanager"

	"github.com/spf13/cobra"
)

var (
	globalCtx         context.Context
	configPaths       []string
	configDirectories []string
	workingDir        string
	disableColor      bool
)

var mainCommand = &cobra.Command{
	Use:              "sing-box",
	PersistentPreRun: preRun,
}

func init() {

	mainCommand.PersistentFlags().StringVarP(&constant.ApiHost, "api-host", "", "", "set api host")
	mainCommand.PersistentFlags().StringVarP(&constant.ApiPort, "api-port", "", "", "set api port")
	mainCommand.PersistentFlags().StringVarP(&constant.DbHost, "mysql-host", "", "", "set mysql host")
	mainCommand.PersistentFlags().StringVarP(&constant.DbPort, "mysql-port", "", "", "set mysql port")
	mainCommand.PersistentFlags().StringVarP(&constant.DbUsername, "mysql-user", "", "", "set mysql username")
	mainCommand.PersistentFlags().StringVarP(&constant.DbPassword, "mysql-pass", "", "", "set mysql password")
	mainCommand.PersistentFlags().StringVarP(&constant.DbName, "mysql-name", "", constant.DbName, "set mysql name default users_db")
	mainCommand.PersistentFlags().BoolVarP(&constant.DbEnable, "mysql-enable", "", false, "enable mysql db default false")

	mainCommand.PersistentFlags().StringArrayVarP(&configPaths, "config", "c", nil, "set configuration file path")
	mainCommand.PersistentFlags().StringArrayVarP(&configDirectories, "config-directory", "C", nil, "set configuration directory path")
	mainCommand.PersistentFlags().StringVarP(&workingDir, "directory", "D", "", "set working directory")
	mainCommand.PersistentFlags().BoolVarP(&disableColor, "disable-color", "", false, "disable color output")

}

func main() {
	if err := mainCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}

func preRun(cmd *cobra.Command, args []string) {
	globalCtx = context.Background()
	sudoUser := os.Getenv("SUDO_USER")
	sudoUID, _ := strconv.Atoi(os.Getenv("SUDO_UID"))
	sudoGID, _ := strconv.Atoi(os.Getenv("SUDO_GID"))
	if sudoUID == 0 && sudoGID == 0 && sudoUser != "" {
		sudoUserObject, _ := user.Lookup(sudoUser)
		if sudoUserObject != nil {
			sudoUID, _ = strconv.Atoi(sudoUserObject.Uid)
			sudoGID, _ = strconv.Atoi(sudoUserObject.Gid)
		}
	}
	if sudoUID > 0 && sudoGID > 0 {
		globalCtx = filemanager.WithDefault(globalCtx, "", "", sudoUID, sudoGID)
	}
	if disableColor {
		log.SetStdLogger(log.NewDefaultFactory(context.Background(), log.Formatter{BaseTime: time.Now(), DisableColors: true}, os.Stderr, "", nil, false).Logger())
	}
	if workingDir != "" {
		_, err := os.Stat(workingDir)
		if err != nil {
			filemanager.MkdirAll(globalCtx, workingDir, 0o777)
		}
		err = os.Chdir(workingDir)
		if err != nil {
			log.Fatal(err)
		}
	}
	if len(configPaths) == 0 && len(configDirectories) == 0 {
		configPaths = append(configPaths, "config.json")
	}
}
