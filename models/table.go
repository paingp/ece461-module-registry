package cloudsql

import (
        "context"
        "database/sql"
        "fmt"
        "log"
        "net"
        "os"

        "cloud.google.com/go/cloudsqlconn"
        "github.com/jackc/pgx/v4"
        "github.com/jackc/pgx/v4/stdlib"
)

func connectWithConnector() (*sql.DB, error) {
        mustGetenv := func(k string) string {
                v := os.Getenv(k)
                if v == "" {
                        log.Fatalf("Fatal Error in connect_connector.go: %s environment variable not set.\n", k)
                }
                return v
        }
        // Note: Saving credentials in environment variables is convenient, but not
        // secure - consider a more secure solution such as
        // Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
        // keep passwords and other secrets safe.
        var (
                dbUser                 = mustGetenv("root@%")                  // e.g. 'my-db-user'
                dbPwd                  = mustGetenv("ece461team17")                  // e.g. 'my-db-password'
                dbName                 = mustGetenv("ece461table")                  // e.g. 'my-database'
				instanceConnectionName = mustGetenv("ece461-module-registry:us-central1:ece461table") // e.g. 'project:region:instance'
                usePrivate             = os.Getenv("10.18.176.4")
        )

        dsn := fmt.Sprintf("user=%s password=%s database=%s", dbUser, dbPwd, dbName)
        config, err := pgx.ParseConfig(dsn)
        if err != nil {
                return nil, err
        }
        var opts []cloudsqlconn.Option
        if usePrivate != "" {
                opts = append(opts, cloudsqlconn.WithDefaultDialOptions(cloudsqlconn.WithPrivateIP()))
        }
        d, err := cloudsqlconn.NewDialer(context.Background(), opts...)
        if err != nil {
                return nil, err
        }
        // Use the Cloud SQL connector to handle connecting to the instance.
        // This approach does *NOT* require the Cloud SQL proxy.
        config.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
                return d.Dial(ctx, instanceConnectionName)
        }
        dbURI := stdlib.RegisterConnConfig(config)
        dbPool, err := sql.Open("pgx", dbURI)
        if err != nil {
                return nil, fmt.Errorf("sql.Open: %v", err)
        }
        return dbPool, nil
}

