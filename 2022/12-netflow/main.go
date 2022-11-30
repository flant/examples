package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/exp/slices"

	//
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	// import various formatters
	"github.com/netsampler/goflow2/format"
	_ "github.com/netsampler/goflow2/format/json"

	//_ "github.com/netsampler/goflow2/format/protobuf"
	//_ "github.com/netsampler/goflow2/format/text"

	// import various transports
	"github.com/netsampler/goflow2/transport"
	//_ "github.com/netsampler/goflow2/transport/file"
	//_ "github.com/netsampler/goflow2/transport/kafka"

	"github.com/netsampler/goflow2/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	version    = ""
	buildinfos = ""
	AppVersion = "GoFlow2 " + version + " " + buildinfos

	ReusePort       = flag.Bool("reuseport", false, "Enable so_reuseport")
	ListenAddresses = flag.String("listen", "netflow://:2055", "listen addresses")

	Workers  = flag.Int("workers", 1, "Number of workers per collector")
	LogLevel = flag.String("loglevel", "info", "Log level")
	LogFmt   = flag.String("logfmt", "json", "Log formatter")

	MetricsAddr = flag.String("metrics.addr", ":8080", "Metrics address")
	MetricsPath = flag.String("metrics.path", "/metrics", "Metrics path")

	Version = flag.Bool("v", false, "Print version")

	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file. If not defined ServiceAccount will be used")

	ServiceConnections = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "netflow_connection_count",
			Help: "metric_help about ",
		},
		[]string{"srchost", "dsthost", "srcnamespace", "dstnamespace"},
	)

	ServiceBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "netflow_connection_bytes",
			Help: "metric_help about ",
		},
		[]string{"srchost", "dsthost", "srcnamespace", "dstnamespace"},
	)

	svc_list       = make(map[string]pod_info)
	svc_list_mutex = sync.RWMutex{}
	dev_namespaces = []string{"dev", "test", "test01", "test02", "test03", "test04", "test05", "test-99"}
)

type pod_info struct {
	name      string
	namespace string
}

func (d *PromDriver) Init(context.Context) error {
	d.q = make(chan bool, 1)
	return nil
}

type PromDriver struct {
	lock *sync.RWMutex
	q    chan bool
}

func (d *PromDriver) Prepare() error {
	return nil
}

func (d *PromDriver) Send(key, data []byte) error {
	//d.lock.RLock()
	log.Debugf(string(data))

	var dat map[string]interface{}
	err := json.Unmarshal(data, &dat)
	if err != nil {
		log.Error(err)
	} else {
		svc_list_mutex.RLock()
		srchost := svc_list[dat["SrcAddr"].(string)].name
		dsthost := svc_list[dat["DstAddr"].(string)].name
		srcnamespace := svc_list[dat["SrcAddr"].(string)].namespace
		dstnamespace := svc_list[dat["DstAddr"].(string)].namespace
		svc_list_mutex.RUnlock()
		if srchost == "" {
			srchost = "unknown"
		}
		if dsthost == "" {
			dsthost = "unknown"
		}
		if srcnamespace == "" {
			srcnamespace = "unknown"
		}
		if dstnamespace == "" {
			dstnamespace = "unknown"
		}
		if slices.Contains(dev_namespaces, srcnamespace) || strings.Contains(srcnamespace, "review-") {
			srcnamespace = "dev-aggregated"
		}
		if slices.Contains(dev_namespaces, dstnamespace) || strings.Contains(dstnamespace, "review-") {
			dstnamespace = "dev-aggregated"
		}

		log.Debugf("%s %s %s %s \n", srchost, dsthost, srcnamespace, dstnamespace)

		ServiceConnections.With(
			prometheus.Labels{
				"srchost":      string(srchost),
				"dsthost":      string(dsthost),
				"srcnamespace": string(srcnamespace),
				"dstnamespace": string(dstnamespace),
			}).
			Inc()

		ServiceBytes.With(
			prometheus.Labels{
				"srchost":      string(srchost),
				"dsthost":      string(dsthost),
				"srcnamespace": string(srcnamespace),
				"dstnamespace": string(dstnamespace),
			}).
			Add(dat["Bytes"].(float64))
	}
	//d.lock.RUnlock()
	return err
}

func (d *PromDriver) Close(context.Context) error {
	return nil
}

func delPod(pod *v1.Pod) {
	if pod.Status.PodIP != "" {
		log.Infof("pod deleted %s %s %s\n", pod.ObjectMeta.Namespace, pod.Name, pod.Status.PodIP)
		svc_list_mutex.Lock()
		delete(svc_list, pod.Status.PodIP)
		svc_list_mutex.Unlock()
	}
}
func newPod(pod *v1.Pod) {
	if pod.Status.PodIP != "" {
		log.Infof("pod added %s %s %s\n", pod.ObjectMeta.Namespace, pod.Name, pod.Status.PodIP)
		app := pod.Name
		if pod.ObjectMeta.Labels["app"] != "" {
			app = pod.ObjectMeta.Labels["app"]
		} else if pod.ObjectMeta.Labels["name"] != "" {
			app = pod.ObjectMeta.Labels["name"]
		} else if pod.ObjectMeta.Labels["k8s-app"] != "" {
			app = pod.ObjectMeta.Labels["k8s-app"]
		} else if pod.ObjectMeta.Labels["app.kubernetes.io/name"] != "" {
			app = pod.ObjectMeta.Labels["app.kubernetes.io/name"]
		} else if pod.ObjectMeta.Labels["component"] != "" {
			app = pod.ObjectMeta.Labels["component"]
		} else if pod.ObjectMeta.Labels["job-name"] != "" {
			app = pod.ObjectMeta.Labels["job-name"]
		}
		svc_list_mutex.Lock()
		svc_list[pod.Status.PodIP] = pod_info{app, pod.ObjectMeta.Namespace}
		svc_list_mutex.Unlock()
	}
}

func init() {
	d := &PromDriver{
		lock: &sync.RWMutex{},
	}
	transport.RegisterTransportDriver("prom", d)
}

func httpServer() {
	http.Handle(*MetricsPath, promhttp.Handler())
	log.Fatal(http.ListenAndServe(*MetricsAddr, nil))
}

func main() {
	flag.Parse()

	if *Version {
		fmt.Println(AppVersion)
		os.Exit(0)
	}

	lvl, _ := log.ParseLevel(*LogLevel)
	log.SetLevel(lvl)

	var config utils.ProducerConfig

	ctx := context.Background()

	formatter, err := format.FindFormat(ctx, "json")
	if err != nil {
		log.Fatal(err)
	}

	transporter, err := transport.FindTransport(ctx, "prom")
	if err != nil {
		log.Fatal(err)
	}
	defer transporter.Close(ctx)

	switch *LogFmt {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	}

	prometheus.MustRegister(ServiceConnections)
	prometheus.MustRegister(ServiceBytes)

	log.Info("Starting GoFlow2")

	var kconfig *rest.Config

	if len(*kubeconfig) > 0 {
		kconfig, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	} else {
		kconfig, err = rest.InClusterConfig()
	}

	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(kconfig)
	if err != nil {
		log.Fatal(err)
	}

	watchlist := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(),
		string(v1.ResourcePods),
		v1.NamespaceAll,
		fields.Everything(),
	)
	_, controller := cache.NewInformer( // also take a look at NewSharedIndexInformer
		watchlist,
		&v1.Pod{},
		0, //Duration is int64
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pod, ok := obj.(*v1.Pod)
				if ok {
					newPod(pod)
				}
			},
			DeleteFunc: func(obj interface{}) {
				pod, ok := obj.(*v1.Pod)
				if ok {
					delPod(pod)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				pod, ok := newObj.(*v1.Pod)
				if ok {
					newPod(pod)
				}
			},
		},
	)
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(stop)

	go httpServer()

	wg := &sync.WaitGroup{}

	for _, listenAddress := range strings.Split(*ListenAddresses, ",") {
		wg.Add(1)
		go func(listenAddress string) {
			defer wg.Done()
			listenAddrUrl, err := url.Parse(listenAddress)
			if err != nil {
				log.Fatal(err)
			}

			hostname := listenAddrUrl.Hostname()
			port, err := strconv.ParseUint(listenAddrUrl.Port(), 10, 64)
			if err != nil {
				log.Errorf("Port %s could not be converted to integer", listenAddrUrl.Port())
				return
			}

			logFields := log.Fields{
				"scheme":   listenAddrUrl.Scheme,
				"hostname": hostname,
				"port":     port,
			}

			log.WithFields(logFields).Info("Starting collection")

			if listenAddrUrl.Scheme == "netflow" {
				sNF := &utils.StateNetFlow{
					Format:    formatter,
					Transport: transporter,
					Logger:    log.StandardLogger(),
					Config:    config,
				}
				err = sNF.FlowRoutine(*Workers, hostname, int(port), *ReusePort)
			} else {
				log.Errorf("scheme %s does not exist", listenAddrUrl.Scheme)
				return
			}

			if err != nil {
				log.WithFields(logFields).Fatal(err)
			}

		}(listenAddress)

	}
	wg.Wait()
}
