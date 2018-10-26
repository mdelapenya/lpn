package cmd

import (
	"net/http"
	"strings"

	docker "github.com/mdelapenya/lpn/docker"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(webCmd)
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Runs a web server with an user interface for LPN",
	Long: `Runs a web server with an user interface to configure the execution of LPN, allowing
	running it on a server`,
	Run: func(cmd *cobra.Command, args []string) {

		startWebApp()
	},
}

func startWebApp() {
	router := gin.Default()

	router.Use(static.Serve("/", static.LocalFile("./webui", true)))

	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	api.GET("/instances", instanceHandler)
	api.POST("/instance/remove/:instanceID", removeInstance)
	api.POST("/instance/stop/:instanceID", stopInstance)

	router.Run(":3000")
}

type containerAction func(string) error

func applyToInstance(c *gin.Context, function containerAction) {
	lpnInstances := getInstances()

	instanceID := c.Param("instanceID")

	for i := 0; i < len(lpnInstances); i++ {
		instance := lpnInstances[i]

		name := strings.Replace(instance.Name, "/", "", -1)

		if name == instanceID {
			function(name)

			lpnInstances = append(lpnInstances[:i], lpnInstances[i+1:]...)

			break
		}
	}

	c.JSON(http.StatusOK, &lpnInstances)
}

func getInstances() []docker.ContainerInstance {
	instances, _ := docker.PsFilter("lpn")

	lpnInstances := []docker.ContainerInstance{}

	for _, instance := range instances {
		lpnInstances = append(lpnInstances, docker.ContainerInstance{
			ID:     instance.ID,
			Name:   instance.Names[0],
			Status: instance.Status,
		})
	}

	return lpnInstances
}

// instanceHandler retrieves a list of all lpn instances that are running in the host
func instanceHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, getInstances())
}

// removeInstance removes a particular instance
func removeInstance(c *gin.Context) {
	applyToInstance(c, docker.RemoveDockerContainer)
}

// stopInstance removes a particular instance
func stopInstance(c *gin.Context) {
	applyToInstance(c, docker.StopDockerContainer)
}
