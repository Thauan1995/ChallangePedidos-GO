package rest

import (
	"challange/pedido"
	"challange/utils"
	"challange/utils/log"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func OrdemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		InsereOrdem(w, r)
		return
	}

	utils.RespondWithError(w, http.StatusMethodNotAllowed, 0, "Método não permitido")
	return
}

func InsereOrdem(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warningf(c, "Erro ao receber corpo da requisição")
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao receber o corpo da requisição")
		return
	}

	var order []pedido.Ordem
	if err = json.Unmarshal(body, &order); err != nil {
		log.Warningf(c, "Falha ao realizar unmarshal da requisição")
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Falha ao realizar unmarshal da requisição")
		return
	}

	if err := pedido.Save_order_to_db(c, order); err != nil {
		log.Warningf(c, "Erro na criação da ordem")
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro na criação da ordem")
		return
	}

	utils.RespondWithJSON(w, http.StatusAccepted, order)
}
