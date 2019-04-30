package trns

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	srv "github.com/hectormao/facele/internal/srv"
	srvImpl "github.com/hectormao/facele/internal/srv/impl"
	"github.com/hectormao/facele/pkg/cfg"
	repo "github.com/hectormao/facele/pkg/repo/impl"
)

type cargaFacturaRequest struct {
	Empresa       string `json:"empresa"`
	NombreArchivo string `json:"nombre_archivo"`
	Contenido     []byte `json:"contenido"`
}

type cargaFacturaResponse struct {
	Id  string `json:"id"`
	Err string `json:"error,omitempty"`
}

func crearCargarFacturaEndpoint(cargaFacturaSrv srv.CargaFacturaSrv) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(cargaFacturaRequest)
		id, err := cargaFacturaSrv.Cargar(req.Empresa, req.NombreArchivo, req.Contenido)
		if err != nil {
			return cargaFacturaResponse{"", err.Error()}, err
		}
		return cargaFacturaResponse{id, ""}, nil
	}
}

func decodeCargaFacturaRequest(_ context.Context, r *http.Request) (interface{}, error) {

	r.ParseMultipartForm(32 << 20)

	empresa := r.Header.Get("X-Empresa")
	log.Printf("%v", empresa)

	archivo, handler, err := r.FormFile("archivo")
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	defer archivo.Close()
	log.Printf("%v", handler.Header)
	data, err := ioutil.ReadAll(archivo)
	if err != nil {
		log.Printf("%v\n", err)
		return nil, err
	}

	request := cargaFacturaRequest{
		Empresa:       empresa,
		NombreArchivo: handler.Filename,
		Contenido:     data,
	}

	return request, nil
}

func encodeHelloWorldResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func StartRestAPI(config cfg.FaceleConfigType) {
	log.Printf("Starting Web Server")
	cargaFacturaRepo := repo.FacturaRepoImpl{}
	colaFacturaRepo := repo.ColaRepoImpl{}
	cargaFacturaSrv := srvImpl.CargaFacturaSrvImpl{
		cargaFacturaRepo,
		colaFacturaRepo,
	}

	cargaFacturaHandler := httptransport.NewServer(
		crearCargarFacturaEndpoint(cargaFacturaSrv),
		decodeCargaFacturaRequest,
		encodeHelloWorldResponse,
	)

	router := mux.NewRouter()
	router.Path(config.WebServer.Path).Methods(config.WebServer.Method).Handler(cargaFacturaHandler)

	listenDir := fmt.Sprintf(":%d", config.WebServer.Port)

	log.Fatal(http.ListenAndServe(listenDir, router))

}
