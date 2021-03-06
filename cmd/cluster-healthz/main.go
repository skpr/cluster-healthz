package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/skpr/cluster-healthz/internal/checks"
)

var (
	cliConfig = kingpin.Flag("config", "Path to the Kubernetes config file").Short('c').String()
	cliToken  = kingpin.Flag("token", "Token used for authentication").Required().Envar("CLUSTER_HEALTHZ_TOKEN").String()
)

func main() {
	kingpin.Parse()

	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {

		if c.Query("token") != *cliToken {
			c.String(http.StatusForbidden, "Access denied")
			return
		}

		config, err := clientcmd.BuildConfigFromFlags("", *cliConfig)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			c.Error(err)
			return
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			c.Error(err)
			return
		}

		issues, err := checks.IssueList(clientset)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			c.Error(err)
			return
		}

		if len(issues) > 0 {
			c.JSON(http.StatusInternalServerError, issues)

			for _, issue := range issues {
				log.With("issue", issue.Issue).With("name", issue.Name).With("namespace", issue.Namespace).Error(issue.Description)
			}

			return
		}

		c.String(http.StatusOK, "Healthy!")
	})

	err := r.Run()
	if err != nil {
		panic(err)
	}
}
