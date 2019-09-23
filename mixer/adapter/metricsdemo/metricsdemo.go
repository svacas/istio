// Copyright 2018 Istio Authors
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

// nolint:lll
// Generates the metricsdemo adapter's resource yaml. It contains the adapter's configuration, name, 
// supported template names (metric in this case), and whether it is session or no-session based.
//go:generate $GOPATH/src/istio.io/istio/bin/mixer_codegen.sh -a mixer/adapter/metricsdemo/config/config.proto -x "-s=false -n metricsdemo -t metric"

package metricsdemo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc"

	"istio.io/api/mixer/adapter/model/v1beta1"
	policy "istio.io/api/policy/v1beta1"
	"istio.io/istio/mixer/template/metric"
	"istio.io/istio/pkg/log"
)

const GatewayUrl = "http://127.0.0.1:9999/mule/smanalytics/applications/service-mesh"

type (
	// Server is basic server interface
	Server interface {
		Addr() string
		Close() error
		Run(shutdown chan error)
	}

	// MetricsDemo supports metric template.
	MetricsDemo struct {
		listener net.Listener
		server   *grpc.Server
		client	 *http.Client
	}
)

var _ metric.HandleMetricServiceServer = &MetricsDemo{}

// HandleMetric records metric entries
func (s *MetricsDemo) HandleMetric(ctx context.Context, r *metric.HandleMetricRequest) (*v1beta1.ReportResult, error) {
	log.Infof("received request %v\n", *r)

    values := map[string]string{}
	var sb strings.Builder
	for _, m := range instances(r.Instances) {
		sb.WriteString(fmt.Sprintf("Name: %v\nValue: %v\nDimensions:\n", m.Name, m.Value))
		for k, v := range m.Dimensions {
			sb.WriteString(fmt.Sprintf("\t%v: %v\n", k, v))
			values[k] = fmt.Sprintf("%v", v)
		}
		sb.WriteString("---------------------------------\n")
	}
	log.Infof(fmt.Sprintf("HandleMetric invoked with:\n  Instances: %s\n", sb.String()))

    jsonValue, _ := json.Marshal(values)
    resp, err := http.Post(GatewayUrl, "application/json", bytes.NewBuffer(jsonValue))
	if err == nil {
		defer cleanupResponse(resp)
	}
	if err != nil {
		log.Errorf("Error connecting to API Gateway %v", err)
	} else {
		log.Infof("Mule Agent response: " + resp.Status)
	}
	return &v1beta1.ReportResult{}, nil
}

func cleanupResponse(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, err := io.Copy(ioutil.Discard, resp.Body)
		if err != nil {
			log.Infof("Error consuming body: %v", err)
		}
		err = resp.Body.Close()
		if err != nil {
			log.Infof("Error closing body: %v", err)
		}
	}
}

func instances(in []*metric.InstanceMsg) []*metric.Instance {
	out := make([]*metric.Instance, 0, len(in))
	for _, inst := range in {
		out = append(out, &metric.Instance{
			Name:       inst.Name,
			Value:      decodeValue(inst.Value.GetValue()),
			Dimensions: decodeDimensions(inst.Dimensions),
		})
	}
	return out
}

func decodeDimensions(in map[string]*policy.Value) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = decodeValue(v.GetValue())
	}
	return out
}

func decodeValue(in interface{}) interface{} {
	switch t := in.(type) {
	case *policy.Value_StringValue:
		return t.StringValue
	case *policy.Value_Int64Value:
		return t.Int64Value
	case *policy.Value_DoubleValue:
		return t.DoubleValue
	case *policy.Value_IpAddressValue:
		return t.IpAddressValue.Value
	case *policy.Value_TimestampValue:
		return t.TimestampValue.Value.String()
	default:
		return fmt.Sprintf("%v", in)
	}
}

// Addr returns the listening address of the server
func (s *MetricsDemo) Addr() string {
	return s.listener.Addr().String()
}

// Run starts the server run
func (s *MetricsDemo) Run(shutdown chan error) {
	shutdown <- s.server.Serve(s.listener)
}

// Close gracefully shuts down the server; used for testing
func (s *MetricsDemo) Close() error {
	if s.server != nil {
		s.server.GracefulStop()
	}

	if s.listener != nil {
		_ = s.listener.Close()
	}

	return nil
}

// NewMetricsDemo creates a new IBP adapter that listens at provided port.
func NewMetricsDemo(addr string) (Server, error) {
	if addr == "" {
		addr = "0"
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", addr))
	if err != nil {
		return nil, fmt.Errorf("unable to listen on socket: %v", err)
	}
	s := &MetricsDemo{
		listener: listener,
		client:   &http.Client{Transport: customTransport(), Timeout: time.Second * 10},
	}
	fmt.Printf("listening on \"%v\"\n", s.Addr())
	s.server = grpc.NewServer()
	metric.RegisterHandleMetricServiceServer(s.server, s)
	return s, nil
}

func customTransport() *http.Transport {
    // Customize the Transport to have larger connection pool
    defaultRoundTripper := http.DefaultTransport
    defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
    if !ok {
        panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
    }
    defaultTransport := *defaultTransportPointer // dereference it to get a copy of the struct that the pointer points to
    defaultTransport.MaxIdleConns = 100
    defaultTransport.MaxIdleConnsPerHost = 100
    return &defaultTransport
}
