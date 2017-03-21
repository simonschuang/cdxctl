package cmd

import (
	"fmt"
	"os"
        "log"
	"database/sql"
	"text/tabwriter"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(listDeviceCommand)
}

var listDeviceCommand = &cobra.Command{
        Use:   "list-device",
        Short: "list network devices",
        RunE: func(cmd *cobra.Command, args []string) error {
		// Open db for read/write scan result
		db, err := sql.Open("sqlite3", "/tmp/cdxctl.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		sqlStmt := `create table if not exists node(
			if_name varchar(50),
			vlan_id integer default 0,
			ipv4 varchar(50),
			ipv4b bolb,
			mac varchar(50),
			hostname varchar(50),
			vendor varchar(200),
			state varchar(10),
			last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			unique(vlan_id, ipv4, mac) on conflict replace
		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlStmt)
			return err
		}
		// Print output
		w := tabwriter.NewWriter(os.Stdout, 5, 0, 2, ' ', 0)
		fmt.Fprintln(w, "IP\tMAC\tSTATE\tVENDOR")

		rows, err := db.Query("select ipv4, mac, vendor, state from node order by ipv4b asc")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var ip string
			var mac string
			var vendor string
			var state string
			err = rows.Scan(&ip, &mac, &vendor, &state)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", ip, mac, state, vendor)
		}
		w.Flush()
		return nil
        },
}