package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/MiguelMachado-dev/ent-go-basic/ent"
	"github.com/MiguelMachado-dev/ent-go-basic/ent/user"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := "./db/main.db"

	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		fmt.Println("Creating database file...")
		err = os.MkdirAll("./db", os.ModePerm)

		if err != nil {
			log.Fatal(err)
		}

		file, err := os.Create(dbPath)

		if err != nil {
			log.Fatal(err)
		}

		file.Close()
	}

	client, err := ent.Open("sqlite3", dbPath+"?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	ctx := context.Background()
	// Run the auto migration tool.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	if _, err = CreateUser(ctx, client); err != nil {
		log.Fatal(err)
	}
	if _, err = QueryUser(ctx, client); err != nil {
		log.Fatal(err)
	}
}

func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", u)
	return u, nil
}

func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Query().
		Where(user.Name("a8m")).
		// `Only` fails if no user found,
		// or more than 1 user returned.
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}
	log.Println("user returned: ", u)
	return u, nil
}
