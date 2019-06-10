package statsd

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strconv"

	"github.com/faja/tornimo-agent/pkg/metrics"
)

var (
	fieldSeparator = []byte("|")
	valueSeparator = []byte(":")
)

func parsePacket(packet *packet, metricSamples []*metrics.MetricSample) []*metrics.MetricSample {

	for {
		message := nextMessage(&packet.messages)
		if message == nil {
			break
		}

		sample, err := parseMetricMessage(message)
		if err != nil {
			// TODO log and stats
			log.Println(err)
			continue
		}

		metricSamples = append(metricSamples, sample)
	}

	return metricSamples
}

func parseMetricMessage(message []byte) (*metrics.MetricSample, error) {
	// metric.name:42|g
	// metric.name:42|c
	// TODO add support to other statsd types
	// currently supported 'c' and 'g'

	// TODO add all the checks
	separatorCount := bytes.Count(message, fieldSeparator)
	if separatorCount < 1 || separatorCount > 3 {
		return nil, fmt.Errorf("invalid field format for %q", message)
	}

	rawNameAndValue, remainder := nextField(message, fieldSeparator)
	rawName, rawValue := nextField(rawNameAndValue, valueSeparator)
	// TODO add tag support, for now ignoring all other stuff after type
	rawType, _ := nextField(remainder, fieldSeparator)

	if len(rawName) == 0 || len(rawValue) == 0 || len(rawType) == 0 {
		return nil, fmt.Errorf("invalid metric message format, empty metric name, value or type")
	}

	var sType metrics.MetricType
	switch string(rawType) {
	case "g":
		sType = metrics.GaugeType
	case "c":
		sType = metrics.CountType
	default:
		return nil, fmt.Errorf("invalid metric message format, not supported metric type '%s'", string(rawType))
	}

	floatValue, err := strconv.ParseFloat(string(rawValue), 64)
	if err != nil {
		// TODO nicer error
		return nil, err
	}

	s := &metrics.MetricSample{
		Name:      string(rawName),
		Value:     floatValue,
		Timestamp: 0,
		Mtype:     sType,
	}
	return s, nil
}

func nextMessage(messages *[]byte) []byte {
	if len(*messages) == 0 {
		return nil
	}

	advance, token, err := bufio.ScanLines(*messages, true)
	if err != nil || len(token) == 0 {
		// TODO log? return not nil?
		return nil
	}

	*messages = (*messages)[advance:]
	return token
}

func nextField(slice, sep []byte) ([]byte, []byte) {
	i := bytes.Index(slice, sep)
	if i == -1 {
		return slice, nil
	}
	return slice[:i], slice[i+1:]
}
