package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
  "os"

	"gopkg.in/fogleman/gg.v1"
)

var DatamallAccountKey = os.Getenv("DATAMALL_ACCOUNT_KEY")

type DatamallNextBus struct {
  EstimatedArrival time.Time
}

type DatamallService struct {
  ServiceNo string
  Operator string
  NextBus *DatamallNextBus
  NextBus2 *DatamallNextBus
  NextBus3 *DatamallNextBus
}

type DatamallResponse struct {
  BusStopCode string
  Services []DatamallService
}

type Service struct {
  Label string
  BusStopCode string
  ServiceNo string
  Timings []time.Time
}

type Data struct {
  Services []Service
  LastUpdated time.Time
}

func main() {
  if DatamallAccountKey == "" {
    panic(errors.New("DATAMALL_ACCOUNT_KEY env is required"))
  }

  services := make([]Service, 6)
  services[0] = Service{
    Label: "ssc",
    BusStopCode: "58271",
    ServiceNo: "859",
  }
  services[1] = Service{
    Label: "assyafaah",
    BusStopCode: "58279",
    ServiceNo: "859",
  }
  services[2] = Service{
    Label: "sunplaza",
    BusStopCode: "58271",
    ServiceNo: "962",
  }
  services[3] = Service{
    Label: "woodlands",
    BusStopCode: "58279",
    ServiceNo: "962",
  }
  services[4] = Service{
    Label: "ktph",
    BusStopCode: "58381",
    ServiceNo: "858",
  }
  services[5] = Service{
    Label: "tampines",
    BusStopCode: "58581",
    ServiceNo: "969",
  }
  data := Data{
    Services: services,
  }

  http.HandleFunc("/", makeHandler(&data))
  log.Println("Starting to listening on :8090")
  http.ListenAndServe(":8090", nil)
}

func makeHandler(data *Data) func(http.ResponseWriter, *http.Request) {
  return func(w http.ResponseWriter, req *http.Request) {
    log.Printf("[%s] GET /", time.Now().Format(time.RFC3339))
    updateBusArrivals(data)
    log.Printf("Data: %v\n", data)
    render(data, w)
  }
}

func updateBusArrivals(data *Data) {
  for i := 0; i < len(data.Services); i++ {
    updatedService, err := getBusArrival(data.Services[i])
    if err != nil {
      data.Services[i].Timings = make([]time.Time, 0)
      continue
    }

    data.Services[i] = updatedService
  }

  data.LastUpdated = time.Now()
}

func getBusArrival(service Service) (Service, error) {
  res, err := getDatamallBusArrival(service.BusStopCode, service.ServiceNo)
  if err != nil {
    return Service{}, err
  }

  service.Timings = make([]time.Time, 0)
  if res.Services[0].NextBus != nil {
    service.Timings = append(service.Timings, res.Services[0].NextBus.EstimatedArrival)
  }
  if res.Services[0].NextBus2 != nil {
    service.Timings = append(service.Timings, res.Services[0].NextBus2.EstimatedArrival)
  }
  if res.Services[0].NextBus3 != nil {
    service.Timings = append(service.Timings, res.Services[0].NextBus3.EstimatedArrival)
  }

  return service, nil
}

func getDatamallBusArrival(busStopCode, serviceNo string) (*DatamallResponse, error) {
  client := &http.Client{}
  req, _ := http.NewRequest("GET", fmt.Sprintf("http://datamall2.mytransport.sg/ltaodataservice/BusArrivalv2?BusStopCode=%s&ServiceNo=%s", busStopCode, serviceNo), nil)
  req.Header.Add("AccountKey", DatamallAccountKey)
  req.Header.Add("Accept", "application/json")
  res, err := client.Do(req)
  if err != nil {
    return nil, err
  }
  defer res.Body.Close()

  jsonBody := DatamallResponse{}
  err = json.NewDecoder(res.Body).Decode(&jsonBody)
  if err != nil {
    return nil, err
  }

  if jsonBody.BusStopCode == "" {
    return nil, errors.New(fmt.Sprintf("Unable to find bus stop code: %s", busStopCode))
  }

  if len(jsonBody.Services) == 0 {
    return nil, errors.New(fmt.Sprintf("Unable to find service %s at bus stop code %s", serviceNo, busStopCode))
  }
  
  return &jsonBody, nil
}

func render(data *Data, w http.ResponseWriter) {
    dc := gg.NewContext(600, 800)
    dc.DrawRectangle(0, 0, 600, 800)
    dc.SetRGB(0, 0, 0)
    dc.Fill()
  
    dc.SetRGB(1, 1, 1)
    x := float64(150)
    y := float64(128)
    for i := 0; i < len(data.Services) && i < 6; i++ {
      service := data.Services[i]
      timing1 := "-"
      if len(service.Timings) >= 1 {
        timing1 = service.Timings[0].Format("3:04")
      }
      timing2 := "-"
      if len(service.Timings) >= 2 {
        timing2 = service.Timings[1].Format("3:04")
      }
      timing3 := "-"
      if len(service.Timings) >= 3 {
        timing3 = service.Timings[2].Format("3:04")
      }

      drawBusTimings(dc, service.Label, service.ServiceNo, timing1, timing2, timing3, x, y)

      x += 150 * 2
      if x > 450 {
        y += 128 * 2
        x = 150
      }
    }

    loadFontFace(dc, "./arial.ttf", 16);
    dc.DrawString(fmt.Sprintf("Last updated: %s", data.LastUpdated.Format("Mon, 2 Jan 2006 03:04:05 PM")), 16, 780) 

    w.Header().Add("Content-Type", "image/png")
    dc.EncodePNG(w)
}

func drawBusTimings(dc *gg.Context, location, busNumber, timing1, timing2, timing3 string, x, y float64) {
    smallFontSize := float64(24)
    largeFontSize := float64(96)
    loadFontFace(dc, "./arial.ttf", smallFontSize);
    dc.DrawStringAnchored(location, x, y - (largeFontSize / 2 + smallFontSize), 0.5, 1) 
    loadFontFace(dc, "./arial.ttf", largeFontSize);
    dc.DrawStringAnchored(busNumber, x, y, 0.5, 0.5) 
    loadFontFace(dc, "./arial.ttf", smallFontSize);
    dc.DrawStringAnchored(fmt.Sprintf("%s   %s   %s", timing1, timing2, timing3), x, y + (largeFontSize / 2 + smallFontSize + (smallFontSize / 4)), 0.5, 0) 
}

func loadFontFace(dc *gg.Context, fontName string, size float64) {
    if err := dc.LoadFontFace(fontName, size); err != nil {
      panic(err)
    }
}

