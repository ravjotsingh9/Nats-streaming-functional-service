package main

import (
	"net/http"
)

/*
func onLogCreated(m event.MeowCreatedMessage) {

		// Index meow for searching
		meow := schema.Meow{
			ID:        m.ID,
			Body:      m.Body,
			CreatedAt: m.CreatedAt,
		}
		if err := search.InsertMeow(context.Background(), meow); err != nil {
			log.Println(err)
		}

}
*/
func getLogsHandler(w http.ResponseWriter, r *http.Request) {

}
