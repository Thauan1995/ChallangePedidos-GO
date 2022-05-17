package pedido

import (
	"challange/utils/consts"
	"challange/utils/log"
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/datastore"
)

const (
	KindOrdens = "Ordens"
)

type Pedido struct {
	IDPedido   int64
	ID         int64
	Date       time.Time
	TotalValue int64
	Status     int64
	Channel    string
	Customer   Customer
	Invoice    Invoice
	Delivery   []Delivery
	Items      []Items
}

type Customer struct {
	Name         string
	Document     string
	Address      string
	Number       int64
	Complement   string
	Neighborhood string
	City         string
	State        string
}

type Invoice struct {
	Number    int64
	Serie     int64
	AccessKey int64
	IssueDate time.Time
	Value     float64
}

type Delivery struct {
	Transporter    string
	TrackingNumber string
	Address        string
	Number         int64
	Complement     string
	Neighborhood   string
	City           string
	State          string
}

type Items struct {
	ID       string
	Sku      string
	Name     string
	Quantity float64
	Price    float64
}

type Ordem struct {
	ID          string `datastore:"-"`
	ExternalID  int64
	StoreID     string
	DataCriacao time.Time
	IsPicked    bool
	IsStockout  bool
}

type Inventory struct {
	ID          string
	ProductID   string
	Sku         string
	DataCriacao time.Time
	Quantity    int64
}

//Valida dados e salva no datastore
func Save_order_to_db(c context.Context, ordem []Ordem) error {
	for _, v := range ordem {
		if v.ID == "" {
			v.DataCriacao = time.Now()
		}

		if v.ExternalID == 0 {
			return fmt.Errorf("Campo ExternalID invalido")
		}

		if v.StoreID == "" {
			return fmt.Errorf("Campo StoreID invalido")
		}
	}

	return PutMultOrdens(c, ordem)

}

func PutOrder(c context.Context, ordem *Ordem) error {
	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Falha ao conectar-se com o Datastore: %v", err)
		return err
	}
	defer datastoreClient.Close()

	key := datastore.NameKey(KindOrdens, ordem.ID, nil)
	key, err = datastoreClient.Put(c, key, ordem)
	if err != nil {
		log.Warningf(c, "Erro ao atualizar publicação")
		return err
	}
	ordem.ID = key.Name
	return nil
}

func PutMultOrdens(c context.Context, ordens []Ordem) error {
	if len(ordens) == 0 {
		return nil
	}
	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Falha ao conectar-se com o Datastore: %v", err)
		return err
	}
	defer datastoreClient.Close()

	keys := make([]*datastore.Key, 0, len(ordens))
	for i := range ordens {
		keys = append(keys, datastore.NameKey(KindOrdens, ordens[i].ID, nil))
	}
	keys, err = datastoreClient.PutMulti(c, keys, ordens)
	if err != nil {
		log.Warningf(c, "Erro ao inserir Multi Ordens: %v", err)
		return err
	}
	return nil
}
