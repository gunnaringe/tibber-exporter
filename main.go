package main

import (
	"context"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/machinebox/graphql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"time"
)

var opts struct {
	Token    string `short:"t" long:"token" description:"Authorization token" required:"true" env:"TIBBER_TOKEN"`
	Endpoint string `short:"e" long:"endpoint" description:"Endpoint" required:"false" env:"TIBBER_ENDPOINT" default:"https://api.tibber.com/v1-beta/gql"`
	Listen   string `short:"l" long:"listen" description:"Address and port to listen" env:"TIBBER_LISTEN" default:":9501"`
	Verbose  []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
}

var client *graphql.Client

var (
	namespace     = "tibber"
	requestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "requests_total",
		Help:      "Number of requests done by exporter",
	})
	lastScrapeDuration = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "scrape_duration_seconds",
		Help:      "Scrape duration of last scrape",
	})
	priceInfoCurrentTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "priceinfo",
			Help:      "...",
		},
		[]string{
			"home_id",
			"type",
			"currency",
		})
)

type Response struct {
	Viewer struct {
		Name  string `json:"name"`
		Homes []struct {
			Id       string `json:"id"`
			TimeZone string `json:"timeZone"`
			Address  struct {
				Address1   string `json:"address1"`
				PostalCode string `json:"postalCode"`
				City       string `json:"city"`
			} `json:"address"`
			Owner struct {
				FirstName   string `json:"firstName"`
				LastName    string `json:"lastName"`
				ContactInfo struct {
					Email  string      `json:"email"`
					Mobile interface{} `json:"mobile"`
				} `json:"contactInfo"`
			} `json:"owner"`
			CurrentSubscription struct {
				PriceInfo struct {
					Current struct {
						Total    float64   `json:"total"`
						Energy   float64   `json:"energy"`
						Tax      float64   `json:"tax"`
						Currency string    `json:"currency"`
						StartsAt time.Time `json:"startsAt"`
					} `json:"current"`
				} `json:"priceInfo"`
			} `json:"currentSubscription"`
			Consumption struct {
				Nodes []struct {
					From            time.Time   `json:"from"`
					To              time.Time   `json:"to"`
					TotalCost       interface{} `json:"totalCost"`
					UnitCost        interface{} `json:"unitCost"`
					UnitPrice       float64     `json:"unitPrice"`
					UnitPriceVAT    float64     `json:"unitPriceVAT"`
					Consumption     interface{} `json:"consumption"`
					ConsumptionUnit string      `json:"consumptionUnit"`
				} `json:"nodes"`
			} `json:"consumption"`
		} `json:"homes"`
	} `json:"viewer"`
}

func updatePrometheus(response Response, scrapeDuration float64) {
	requestsTotal.Inc()
	lastScrapeDuration.Set(scrapeDuration)

	for _, home := range response.Viewer.Homes {
		currentSubscription := home.CurrentSubscription
		priceInfoCurrentTotal.WithLabelValues(home.Id, "energy", currentSubscription.PriceInfo.Current.Currency).Set(currentSubscription.PriceInfo.Current.Energy)
		priceInfoCurrentTotal.WithLabelValues(home.Id, "tax", currentSubscription.PriceInfo.Current.Currency).Set(currentSubscription.PriceInfo.Current.Tax)
	}

}

func scrape() Response {
	req := graphql.NewRequest(`
	    {
	      viewer {
	        name
	        homes {
              id
	          timeZone
	          address {
	            address1
	            postalCode
	            city
	          }
	          owner {
	            firstName
	            lastName
	            contactInfo {
	              email
	              mobile
	            }
	          }
	          currentSubscription{
	            priceInfo{
	              current{
	                total
	                energy
	                tax
                    currency
	                startsAt
	              }
	            }
	          }
	          consumption(resolution: HOURLY, last: 1) {
	            nodes {
	              from
	              to
	              totalCost
	              unitCost
	              unitPrice
	              unitPriceVAT
	              consumption
	              consumptionUnit
	            }
	          }
	        }
	      }
	    }
	    `)
	req.Header.Set("Authorization", opts.Token)
	ctx := context.Background()
	var response Response
	if err := client.Run(ctx, req, &response); err != nil {
		log.Fatal(err)
	}
	return response
}

type updater struct {
	h http.Handler
}

func (c updater) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	response := scrape()
	scrapeDuration := time.Since(start).Seconds()
	updatePrometheus(response, scrapeDuration)
	log.Println("Update OK")
	c.h.ServeHTTP(w, r)
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	client = graphql.NewClient(opts.Endpoint)
	if len(opts.Verbose) > 0 {
		client.Log = func(s string) { log.Println(s) }
	}

	fmt.Println("Server is listening on " + opts.Listen)
	http.Handle("/", promhttp.Handler())
	fmt.Println(http.ListenAndServe(opts.Listen, updater{http.DefaultServeMux}))
}
