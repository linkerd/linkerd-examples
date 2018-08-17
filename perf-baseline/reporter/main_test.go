package main

import "testing"

type reportExpected struct {
	report statsReport
	output string
}

func TestPublish(t *testing.T) {
	t.Run("Publishes expected stat output", func(t *testing.T) {

		expectations := []reportExpected{
			reportExpected{
				report: statsReport{
					"testproto1": map[string]*stats{
						"linkerd2-test": &stats{
							sr: .9,
							rr: 321,
							latencies: map[float64]uint64{
								0.5:   321,
								0.75:  322,
								0.9:   543,
								0.95:  654,
								0.99:  987,
								0.999: 9870,
							},
							mem: 987654321,
							cpu: 0.5,
						},
					},
				},
				output: `Stats report:
Protocol: testproto1:
  linkerd2-test
    Success rate:     90.00%
    Request rate:     321
    p50 latency (us): 321
    p75 latency (us): 322
    p90 latency (us): 543
    p95 latency (us): 654
    p99 latency (us): 987
    p99 latency (us): 9870
    Memory (bytes):   987654321
    CPU (cores):      0.500
`,
			},
		}

		for _, exp := range expectations {
			if exp.report.String() != exp.output {
				t.Fatalf("Expected: %+v\nGot: %+v", exp.output, exp.report.String())
			}
		}
	})
}
