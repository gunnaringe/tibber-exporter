package main

import (
	"context"
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
	info = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "home_info",
			Help:      "Home information",
		},
		[]string{
			"home_id",
			"time_zone",
			"address1",
			"address2",
			"address3",
			"city",
			"postal_code",
			"country",
			"latitude",
			"longitude",
		})
	priceInfoCurrentTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "current_energy_price",
			Help:      "Current energy price",
		},
		[]string{
			"home_id",
			"type",
			"currency",
		})
	lastPeriodConsumption = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "previous_hour_consumption_watt_hour",
			Help:      "Total Watt hours used last hourly period",
		},
		[]string{
			"home_id",
		})
	lastPeriodTotalCost = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "previous_hour_total_cost",
			Help:      "Total cost last hourly period",
		},
		[]string{
			"home_id",
			"currency",
		})
	lastPeriodPrice = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "previous_hour_energy_price",
			Help:      "Energy price last hourly period",
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
				Address2   string `json:"address2"`
				Address3   string `json:"address3"`
				City       string `json:"city"`
				PostalCode string `json:"postalCode"`
				Country    string `json:"country"`
				Latitude   string `json:"latitude"`
				Longitude  string `json:"longitude"`
			} `json:"address"`
			Owner struct {
				FirstName   string `json:"firstName"`
				LastName    string `json:"lastName"`
				ContactInfo struct {
					Email  string `json:"email"`
					Mobile string `json:"mobile"`
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
						Level    string    `json:"level"`
					} `json:"current"`
				} `json:"priceInfo"`
			} `json:"currentSubscription"`
			Consumption struct {
				Nodes []struct {
					From            time.Time `json:"from"`
					To              time.Time `json:"to"`
					TotalCost       float64   `json:"totalCost"`
					UnitCost        float64   `json:"unitCost"`
					UnitPrice       float64   `json:"unitPrice"`
					UnitPriceVAT    float64   `json:"unitPriceVAT"`
					Consumption     float64   `json:"consumption"`
					ConsumptionUnit string    `json:"consumptionUnit"`
					Currency        string    `json:"currency"`
				} `json:"nodes"`
			} `json:"consumption"`
		} `json:"homes"`
	} `json:"viewer"`
}

func updatePrometheus(response Response, scrapeDuration float64) {
	requestsTotal.Inc()
	lastScrapeDuration.Set(scrapeDuration)
	for _, home := range response.Viewer.Homes {
		info.WithLabelValues(
			home.Id,
			home.TimeZone,
			home.Address.Address1,
			home.Address.Address2,
			home.Address.Address3,
			home.Address.City,
			home.Address.PostalCode,
			home.Address.Country,
			home.Address.Latitude,
			home.Address.Longitude).Set(1)

		currentSubscription := home.CurrentSubscription
		priceInfoCurrentTotal.WithLabelValues(home.Id, "energy", currentSubscription.PriceInfo.Current.Currency).Set(currentSubscription.PriceInfo.Current.Energy)
		priceInfoCurrentTotal.WithLabelValues(home.Id, "tax", currentSubscription.PriceInfo.Current.Currency).Set(currentSubscription.PriceInfo.Current.Tax)

		if len(home.Consumption.Nodes) > 0 {
			node := home.Consumption.Nodes[0]
			if node.ConsumptionUnit != "kWh" {
				log.Printf("Expected unit kWh, but got %s\n", node.ConsumptionUnit)
			}
			// Use base unit; Wh instead of kWh
			lastPeriodConsumption.WithLabelValues(home.Id).Set(node.Consumption / 1000)
			lastPeriodPrice.WithLabelValues(home.Id, "energy", node.Currency).Set(node.UnitPrice)
			lastPeriodPrice.WithLabelValues(home.Id, "vat", node.Currency).Set(node.UnitPriceVAT)
			lastPeriodTotalCost.WithLabelValues(home.Id, node.Currency).Set(node.TotalCost)
		}
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
		address2
		address3
        city
        postalCode
		country
		latitude
		longitude
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
            level
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
          currency
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

	log.Println("Server is listening on " + opts.Listen)
	http.Handle("/", promhttp.Handler())
	log.Println(http.ListenAndServe(opts.Listen, updater{http.DefaultServeMux}))
}
