// elostats

package elostats

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type SQLMetric struct {
	Query      string `toml:"query"`
	MetricName string `toml:"metric_name"`
}

type ELOData struct {
	SampleFrequency config.Duration `toml:"sample_frequency"`
	DBHost          string          `toml:"DBHost"`
	DBName          string          `toml:"DBName"`
	DBUser          string          `toml:"DBUser"`
	DBPassword      string          `toml:"DBPassword"`
	SQLMetrics      []SQLMetric     `toml:"sql_metrics"`
	ctx             context.Context
	cancel          context.CancelFunc

	Log telegraf.Logger `toml:"-"`
}

func init() {
	inputs.Add("elostats", func() telegraf.Input {
		return &ELOData{
			SampleFrequency: config.Duration(5 * time.Second),
		}
	})
}

func (r *ELOData) Init() error {
	return nil
}

func (r *ELOData) SampleConfig() string {
	r.Log.Infof("SampleConfig called")
	return `
	[[inputs.elostats]]
	# Sample frequency
	sample_frequency = "5s"
	# MSSQL Database Host
	DBHost = "localhost"
	# MSSQL Database Name
	DBName = "database"
	# MSSQL Database User
	DBUser = "user"
	# MSSQL Database Password
	DBPassword = "password"
	
	# List of SQL queries and associated metric names
	[[inputs.elostats.sql_metrics]]
	query = "SELECT count(*) FROM dbo.workflowactivedoc where wf_nodeid = 0"
	metric_name = "WFCount"

	[[inputs.elostats.sql_metrics]]
	query = "SELECT max(objid) FROM dbo.objekte"
	metric_name = "MaxObjID"
`
}

func (r *ELOData) Description() string {
	return "Connects to MSSQL database, executes queries, and sends the results to InfluxDB"
}

func (r *ELOData) Gather(a telegraf.Accumulator) error {
	ticker := time.NewTicker(time.Duration(r.SampleFrequency))
	defer ticker.Stop()

	for range ticker.C {
		for _, sqlMetric := range r.SQLMetrics {
			value, err := r.executeQuery(sqlMetric.Query)
			if err != nil {
				return fmt.Errorf("Error executing query: %s", err)
			}

			r.sendMetric(a, value, sqlMetric.MetricName)
		}
	}

	return nil
}

func (r *ELOData) executeQuery(query string) (int, error) {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s",
		r.DBHost, r.DBUser, r.DBPassword, r.DBName)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return 0, fmt.Errorf("Error opening database connection: %s", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var value int
	err = db.QueryRowContext(ctx, query).Scan(&value)
	if err != nil {
		return 0, fmt.Errorf("Error executing query: %s", err)
	}

	return value, nil
}

func (r *ELOData) sendMetric(a telegraf.Accumulator, value int, metricName string) {
	// Send the result to InfluxDB
	tags := map[string]string{"metric": metricName}
	fields := map[string]interface{}{"value": value}
	now := time.Now()
	a.AddFields("elostats", fields, tags, now)
}

func (r *ELOData) Stop() {
	r.cancel()
}
