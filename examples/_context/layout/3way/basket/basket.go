package basket

import (
	. "github.com/go-on/html/ht"
	"github.com/go-on/html/types"
	"gopkg.in/go-on/wrap-contrib.v2/examples/context/layout/3way/model"
	"net/http"
)

type Basket struct{ User *model.User }

func (b *Basket) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	DIV(
		types.Class("basket-state"),
		H1("Basket of "+b.User.Name),
		P("Your basket is currently empty"),
	).WriteTo(w)
}
