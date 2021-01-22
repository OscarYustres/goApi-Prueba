package main

import(
      "time"
      "github.com/gorilla/mux"
      "net/http"
      "log"
      "encoding/json"
      "strconv"
)
type Tickets struct {
  ID string `json:"id,omitempty"`
  Usuario *Usuario `json:"usuario"`
  FechaCreat time.Time `json:"fecha_creacion"`
  FechaUpd time.Time `json:"fecha_actualizacion"`
  Status bool `json:"estado"`
}

type Usuario struct {
  Nombre string `json:"nombre,omitempty"`
  Apellido string `json:"apellido,omitempty"`
}

var ticketStore = make(map[string]Tickets)
var id int

// GetNoteHandler  - GET  /api/note
func GetTicketsHandler( w http.ResponseWriter, r *http.Request)  {
  var ticket []Tickets
   for _, v := range ticketStore {
       ticket = append(ticket, v)
   }
    w.Header().Set("Content-Type", "application/json")
    j, err := json.Marshal(ticket)
    if err != nil {
      panic(err)
    }
    w.WriteHeader(200)
    w.Write(j)
}
func GetTickeHandler(w http.ResponseWriter, r *http.Request){
  params :=  mux.Vars(r)
  for _, item := range ticketStore {
    if item .ID == params["id"] {
      json.NewEncoder(w).Encode(item)
      return
    }
  }
  json.NewEncoder(w).Encode(&Tickets{})
}
// PostNoteHandler  - POST  /api/note
func PostTicketsHandler(w http.ResponseWriter, r *http.Request)  {
  var ticket Tickets
  err := json.NewDecoder(r.Body).Decode(&ticket)
  if err != nil {
         panic(err)
    }
    ticket.FechaCreat = time.Now()
    id++
    k := strconv.Itoa(id)
    ticketStore[k] = ticket

    w.Header().Set("Content-Type", "application/json")
    j, err := json.Marshal(ticket)
    if err != nil {
      panic(err)
    }
    w.WriteHeader(201)
    w.Write(j)

}

// PutNoteHandler  - Put  /api/note
func PutTicketsHandler( w http.ResponseWriter, r *http.Request)  {
  vars := mux.Vars(r)
  k := vars["id"]
  var ticketUpdate Tickets
  err := json.NewDecoder(r.Body).Decode(&ticketUpdate)
  if err != nil {
    panic(err)
    }
    if ticket, ok := ticketStore[k]; ok {
            ticketUpdate.FechaUpd = ticket.FechaUpd
            delete(ticketStore, k)
            ticketStore[k] = ticketUpdate
    }else {
          log.Printf("No encontramos el id %s", k)
    }
          w.WriteHeader(204)
}
// DELETENoteHandler  - DELETE  /api/note
func DeleteTicketsHandler(w http.ResponseWriter, r *http.Request)  {
  vars := mux.Vars(r)
  k := vars["id"]
  if _, ok := ticketStore[k]; ok {

            delete(ticketStore, k)

    }else {
          log.Printf("No encontramos el id %s", k)
    }
          w.WriteHeader(204)
}



func main()  {
   r := mux.NewRouter().StrictSlash(false)
   r.HandleFunc("/api/tickets", GetTicketsHandler).Methods("GET")
   r.HandleFunc("/api/tickets/{id}", GetTickeHandler).Methods("GET")
   r.HandleFunc("/api/tickets/{id}", PostTicketsHandler).Methods("POST")
   r.HandleFunc("/api/tickets/{id}", PutTicketsHandler).Methods("PUT")
   r.HandleFunc("/api/tickets/{id}", DeleteTicketsHandler).Methods("DELETE")

   server := &http.Server{
     Addr:           ":8080",
     Handler:        r,
     ReadTimeout:    10 * time.Second,
     WriteTimeout:   10 * time.Second,
     MaxHeaderBytes: 1 << 20,
     }
     log.Println("Listening http://localhost:8080 ...")
     server.ListenAndServe()

}
