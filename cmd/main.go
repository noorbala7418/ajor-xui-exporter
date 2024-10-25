package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/noorbala7418/ajor-xui-exporter/pkg/xray"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var scrapeTime int

var (
	totalUsers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ajor_xui_total_users_count",
		Help: "Total Users Count",
	})

	usersEnable = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ajor_xui_user_enable",
		Help: "User enable status",
	}, []string{
		"name",
	})

	usersAdminEnable = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ajor_xui_user_admin_state",
		Help: "User Admin State",
	}, []string{
		"name",
	})

	usersTotalTraffic = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ajor_xui_user_total_traffic",
		Help: "User Total Traffic",
	}, []string{
		"name",
	})

	usersRemainTraffic = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ajor_xui_user_remain_traffic",
		Help: "User Remain Traffic",
	}, []string{
		"name",
	})

	usersUploadTraffic = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ajor_xui_user_upload_traffic",
		Help: "User Upload Traffic",
	}, []string{
		"name",
	})

	usersDownloadTraffic = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ajor_xui_user_download_traffic",
		Help: "User Download Traffic",
	}, []string{
		"name",
	})
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only logrus the warning severity or above.
	switch os.Getenv("APP_LOG_MODE") {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
	checkEnvs()
}

func main() {
	logrus.Info("start Ajor XUI exporter")
	logrus.Info("scrape every ", scrapeTime, " minutes.")
	c := gocron.NewScheduler(time.UTC)
	c.Every(scrapeTime).Seconds().Do(checkStatus)
	c.StartAsync()

	logrus.Info("cron started.")

	// expose /metrics for prometheus to gather
	http.Handle("/metrics", promhttp.Handler())

	// http server
	logrus.Info("Starting server at port 9153")
	logrus.Info("Server started. Metrics are available at 0.0.0.0:9153/metrics")
	if err := http.ListenAndServe(":9153", nil); err != nil {
		logrus.Fatal("Failed to start metrics server. ", err)
	}
}

func checkStatus() {
	userList := xray.GetAllClients()

	logrus.Info("userlist received. Length is: ", len(userList))
	totalUsers.Set(float64(len(userList)))

	logrus.Info("set metrics started.")
	for _, user := range userList {
		usersEnable.WithLabelValues(user.Name).Set(convertBooleanToFloat(user.Enable))
		usersAdminEnable.WithLabelValues(user.Name).Set(convertBooleanToFloat(user.AdminEnabled))
		usersTotalTraffic.WithLabelValues(user.Name).Set(float64(user.TotalTraffic))
		usersDownloadTraffic.WithLabelValues(user.Name).Set(float64(user.DownloadTraffic))
		usersUploadTraffic.WithLabelValues(user.Name).Set(float64(user.UploadTraffic))
		usersRemainTraffic.WithLabelValues(user.Name).Set(float64(user.RemainTraffic))
	}
	logrus.Info("set metrics done.")
}

func convertBooleanToFloat(result bool) float64 {
	if result {
		return 1
	}
	return 0
}

// checkEnvs Checks environment variables and if one variable does not exist, Then it will Kill application.
func checkEnvs() {
	if os.Getenv("XPANEL_URL") == "" {
		logrus.Error("env variable $XPANEL_URL is not defined")
		os.Exit(1)
	}

	if os.Getenv("XPANEL_USERNAME") == "" {
		logrus.Error("env variable $XPANEL_USERNAME is not defined")
		os.Exit(1)
	}

	if os.Getenv("XPANEL_PASSWORD") == "" {
		logrus.Error("env variable $XPANEL_PASSWORD is not defined")
		os.Exit(1)
	}
	if os.Getenv("APP_SCRAPE_TIME") == "" {
		logrus.Warning("Env APP_SCRAPE_TIME is not defined. Default value is 30 seconds.")
		scrapeTime = 30
	} else {
		scrapeTime, _ = strconv.Atoi(os.Getenv("APP_SCRAPE_TIME"))
	}
}
