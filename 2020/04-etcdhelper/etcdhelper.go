package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"

	"k8s.io/api/core/v1"
	"k8s.io/kubectl/pkg/scheme"

	"k8s.io/apimachinery/pkg/runtime/serializer/protobuf"

	jsonserializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

func main() {
	var endpoint, keyFile, certFile, caFile string
	flag.StringVar(&endpoint, "endpoint", "https://127.0.0.1:2379", "etcd endpoint.")
	flag.StringVar(&keyFile, "key", "", "TLS client key.")
	flag.StringVar(&certFile, "cert", "", "TLS client certificate.")
	flag.StringVar(&caFile, "cacert", "", "Server TLS CA certificate.")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprint(os.Stderr, "ERROR: you need to specify action: dump or ls [<key>] or get <key>\n")
		os.Exit(1)
	}
	if flag.Arg(0) == "get" && flag.NArg() == 1 {
		fmt.Fprint(os.Stderr, "ERROR: you need to specify <key> for get operation\n")
		os.Exit(1)
	}
	if flag.Arg(0) == "dump" && flag.NArg() != 1 {
		fmt.Fprint(os.Stderr, "ERROR: you cannot specify positional arguments with dump\n")
		os.Exit(1)
	}
	if flag.Arg(0) == "change-service-cidr" && !regexp.MustCompile(`^(([0-9]{1,3}.){3}([0-9])(\/[0-9]{1,2})?)$`).MatchString(flag.Arg(1)) {
		if flag.Arg(1) == "" {
			fmt.Fprint(os.Stderr, "ERROR: you need to specify Service CIDR\n")
		} else if !regexp.MustCompile(`^(([0-9]{1,3}.){3}([0-9])(\/[0-9]{1,2})?)$`).MatchString(flag.Arg(1)) {
			fmt.Fprint(os.Stderr, "ERROR: you need to specify correct Service CIDR\nExample: 10.0.0.0 or 10.0.0.0/16\n")
		}
		os.Exit(1)
	}
	if flag.Arg(0) == "change-pod-cidr" && !regexp.MustCompile(`^(([0-9]{1,3}.){3}([0-9])(\/[0-9]{1,2})?)$`).MatchString(flag.Arg(1)) {
		if flag.Arg(1) == "" {
			fmt.Fprint(os.Stderr, "ERROR: you need to specify Pod CIDR\n")
		} else if !regexp.MustCompile(`^(([0-9]{1,3}.){3}([0-9])(\/[0-9]{1,2})?)$`).MatchString(flag.Arg(1)) {
			fmt.Fprint(os.Stderr, "ERROR: you need to specify correct Pod CIDR\nExample: 10.0.0.0 or 10.0.0.0/16\n")
		}
		os.Exit(1)
	}
	if flag.Arg(0) == "change-monitors-list" {
		if flag.Arg(1) == "" || flag.Arg(2) == "" {
			fmt.Fprint(os.Stderr, "ERROR: you have to specify both: PV name and list of comma-separated monitor IP-addresses\n")
			os.Exit(1)
		}
		if !regexp.MustCompile(`^pvc-[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`).MatchString(flag.Arg(1)) {
			fmt.Fprint(os.Stderr, "ERROR: invalid PV name. Ex: pvc-dd3afe18-1bae-411c-9a1d-df129847cb62")
			os.Exit(1)
		}
		if !regexp.MustCompile(`^(((\d+\.){3}\d+):\d+,?)+$`).MatchString(flag.Arg(2)) {
			fmt.Fprint(os.Stderr, "ERROR: invalid IP-address list. Ex: 1.1.1.1:6789,2.2.2.2:6789,3.3.3.3:6789\n")
			os.Exit(1)
		}
	}

	action := flag.Arg(0)
	key := ""
	if flag.NArg() > 1 {
		key = flag.Arg(1)
	}

	var tlsConfig *tls.Config
	if len(certFile) != 0 || len(keyFile) != 0 || len(caFile) != 0 {
		tlsInfo := transport.TLSInfo{
			CertFile:      certFile,
			KeyFile:       keyFile,
			TrustedCAFile: caFile,
		}
		var err error
		tlsConfig, err = tlsInfo.ClientConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: unable to create client config: %v\n", err)
			os.Exit(1)
		}
	}

	config := clientv3.Config{
		Endpoints:   []string{endpoint},
		TLS:         tlsConfig,
		DialTimeout: 5 * time.Second,
	}
	client, err := clientv3.New(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: unable to connect to etcd: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	switch action {
	case "ls":
		_, err = listKeys(client, key)
	case "get":
		err = getKey(client, key)
	case "change-service-cidr":
		err = changeServiceCIDR(client, flag.Arg(1))
	case "change-pod-cidr":
		err = changePodCIDR(client, flag.Arg(1))
	case "change-monitors-list":
		err = changeMonitorsList(client, flag.Arg(1), flag.Arg(2))
	case "dump":
		err = dump(client)
	default:
		fmt.Fprintf(os.Stderr, "ERROR: invalid action: %s\n", action)
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s-ing %s: %v\n", action, key, err)
		os.Exit(1)
	}
}

func changeServiceCIDR(client *clientv3.Client, cidr string) error {
	// создаем десериализатор
	decoder := scheme.Codecs.UniversalDeserializer()

	re := regexp.MustCompile(`^([0-9]*.[0-9]*)`)
	newServiceCIDR := regexp.MustCompile(`(.[0-9].[0-9/]*$)`).ReplaceAllString(cidr, "")

	// получаем список ключей сервисов
	svcKeyList, _ := listKeys(client, "/registry/services/specs")

	// в цикле проходим по списку, заменяем ClusterIP (если он есть) и сохраняем изменения в etcd
	for _, svcKey := range svcKeyList {
		// по ключу достаём значение из etcd
		resp, err := clientv3.NewKV(client).Get(context.Background(), svcKey)
		if err != nil {
			fmt.Printf("get key %s %s\n", svcKey, err)
		}

		// преобразуем полученные данные в go-объект
		obj, _, _ := decoder.Decode(resp.Kvs[0].Value, nil, nil)

		// полученный объект является интерфейсом, преобразуем его в *v1.Service
		svc := obj.(*v1.Service)
		if svc.Spec.ClusterIP != "None" {
			// с помощью regexp заменяем первые два байта адреса
			newClusterIP := re.ReplaceAllString(svc.Spec.ClusterIP, newServiceCIDR)
			fmt.Printf("%s %s > %s\n", svcKey, svc.Spec.ClusterIP, newClusterIP)

			// присваеваем новый адрес
			svc.Spec.ClusterIP = newClusterIP

			// создаем сериализатор
			protoSerializer := protobuf.NewSerializer(scheme.Scheme, scheme.Scheme)
			newObj := new(bytes.Buffer)
			// преобразуем go-объект в protobuf
			protoSerializer.Encode(obj, newObj)

			_, err = clientv3.NewKV(client).Put(context.Background(), svcKey, newObj.String())
			if err != nil {
				fmt.Printf("put to key %s %s\n", svcKey, err)
			}
		}
	}

	return nil
}

func changePodCIDR(client *clientv3.Client, cidr string) error {
	decoder := scheme.Codecs.UniversalDeserializer()

	re := regexp.MustCompile(`^([0-9]*.[0-9]*)`)
	newPodCIDR := regexp.MustCompile(`(.[0-9].[0-9/]*$)`).ReplaceAllString(cidr, "")

	nodesKeyList, _ := listKeys(client, "/registry/minions")
	for _, nodeKey := range nodesKeyList {
		resp, err := clientv3.NewKV(client).Get(context.Background(), nodeKey)
		if err != nil {
			fmt.Printf("get key %s %s\n", nodeKey, err)
		}

		obj, _, _ := decoder.Decode(resp.Kvs[0].Value, nil, nil)
		node := obj.(*v1.Node)

		newPodCIDR := re.ReplaceAllString(node.Spec.PodCIDR, newPodCIDR)
		fmt.Printf("%s %s > %s\n", nodeKey, node.Spec.PodCIDR, newPodCIDR)

		node.Spec.PodCIDR = newPodCIDR
		node.Spec.PodCIDRs[0] = newPodCIDR

		protoSerializer := protobuf.NewSerializer(scheme.Scheme, scheme.Scheme)
		newObj := new(bytes.Buffer)
		protoSerializer.Encode(obj, newObj)

		_, err = clientv3.NewKV(client).Put(context.Background(), nodeKey, newObj.String())
		if err != nil {
			fmt.Printf("put to key %s %s\n", nodeKey, err)
		}
	}

	return nil
}

func changeMonitorsList(client *clientv3.Client, pvName, list string) error {
	decoder := scheme.Codecs.UniversalDeserializer()

	pvKey := fmt.Sprintf("/registry/persistentvolumes/%s", pvName)

	resp, err := clientv3.NewKV(client).Get(context.Background(), pvKey)
	if err != nil {
		fmt.Printf("get key %s %s\n", pvKey, err)
	}

	obj, _, _ := decoder.Decode(resp.Kvs[0].Value, nil, nil)

	pv := obj.(*v1.PersistentVolume)

	monitors := strings.Split(strings.Trim(list, ","), ",")
	pv.Spec.RBD.CephMonitors = monitors

	protoSerializer := protobuf.NewSerializer(scheme.Scheme, scheme.Scheme)
	newObj := new(bytes.Buffer)
	protoSerializer.Encode(obj, newObj)

	_, err = clientv3.NewKV(client).Put(context.Background(), pvKey, newObj.String())
	if err != nil {
		fmt.Printf("put to key %s %s\n", pvKey, err)
	}

	return nil
}

func listKeys(client *clientv3.Client, key string) ([]string, error) {
	var resp *clientv3.GetResponse
	var err error
	if len(key) == 0 {
		resp, err = clientv3.NewKV(client).Get(context.Background(), "/", clientv3.WithFromKey(), clientv3.WithKeysOnly())
	} else {
		resp, err = clientv3.NewKV(client).Get(context.Background(), key, clientv3.WithPrefix(), clientv3.WithKeysOnly())
	}
	if err != nil {
		return []string{""}, err
	}

	keys := []string{}

	for _, kv := range resp.Kvs {
		fmt.Println(string(kv.Key))
		keys = append(keys, string(kv.Key))
	}

	return keys, err
}

func getKey(client *clientv3.Client, key string) error {
	resp, err := clientv3.NewKV(client).Get(context.Background(), key)
	if err != nil {
		return err
	}

	decoder := scheme.Codecs.UniversalDeserializer()
	encoder := jsonserializer.NewSerializer(jsonserializer.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, true)

	for _, kv := range resp.Kvs {
		obj, gvk, err := decoder.Decode(kv.Value, nil, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARN: unable to decode %s: %v\n", kv.Key, err)
			continue
		}
		fmt.Println(gvk)
		err = encoder.Encode(obj, os.Stdout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARN: unable to encode %s: %v\n", kv.Key, err)
			continue
		}
	}

	return nil
}

func dump(client *clientv3.Client) error {
	response, err := clientv3.NewKV(client).Get(context.Background(), "/", clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	if err != nil {
		return err
	}

	kvData := []etcd3kv{}
	decoder := scheme.Codecs.UniversalDeserializer()
	encoder := jsonserializer.NewSerializer(jsonserializer.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, false)
	objJSON := &bytes.Buffer{}

	for _, kv := range response.Kvs {
		obj, _, err := decoder.Decode(kv.Value, nil, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARN: error decoding value %q: %v\n", string(kv.Value), err)
			continue
		}
		objJSON.Reset()
		if err := encoder.Encode(obj, objJSON); err != nil {
			fmt.Fprintf(os.Stderr, "WARN: error encoding object %#v as JSON: %v", obj, err)
			continue
		}
		kvData = append(
			kvData,
			etcd3kv{
				Key:            string(kv.Key),
				Value:          string(objJSON.Bytes()),
				CreateRevision: kv.CreateRevision,
				ModRevision:    kv.ModRevision,
				Version:        kv.Version,
				Lease:          kv.Lease,
			},
		)
	}

	jsonData, err := json.MarshalIndent(kvData, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonData))

	return nil
}

type etcd3kv struct {
	Key            string `json:"key,omitempty"`
	Value          string `json:"value,omitempty"`
	CreateRevision int64  `json:"create_revision,omitempty"`
	ModRevision    int64  `json:"mod_revision,omitempty"`
	Version        int64  `json:"version,omitempty"`
	Lease          int64  `json:"lease,omitempty"`
}
