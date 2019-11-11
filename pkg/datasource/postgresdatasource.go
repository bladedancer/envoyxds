package datasource

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/bladedancer/envoyxds/pkg/apimgmt"
	_ "github.com/lib/pq" // Need postgres
)

// PostgresDatasource is a tenant datasource that reads tenant information from the database.
type PostgresDatasource struct {
	tenants       []*apimgmt.Tenant
	tenantsByName map[string]*apimgmt.Tenant
	updateChan    chan []*apimgmt.Tenant
	db            *sql.DB
	timestamp     time.Time
}

// MakePostgresDatasource Initialize the postgres datasource
func MakePostgresDatasource() *PostgresDatasource {
	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	ds := &PostgresDatasource{
		updateChan: make(chan []*apimgmt.Tenant),
		db:         db,
	}

	ds.tenants, ds.timestamp = ds.loadTenants()
	ds.tenantsByName = make(map[string]*apimgmt.Tenant)

	for _, tenant := range ds.tenants {
		ds.tenantsByName[tenant.Name] = tenant
	}

	ds.startPoll()
	return ds
}

// GetTenants Get the tenant data and a channel hat is updated on change
func (ds *PostgresDatasource) GetTenants() ([]*apimgmt.Tenant, chan []*apimgmt.Tenant) {
	return ds.tenants, ds.updateChan
}

// GetTenant finds the tenant with the specified name.
func (ds *PostgresDatasource) GetTenant(name string) *apimgmt.Tenant {
	return ds.tenantsByName[name]
}

// UpsertTenant update existing tenant in the store or add if not already present.
func (ds *PostgresDatasource) UpsertTenant(tenant *apimgmt.Tenant) {
	panic("Upsert not implemented for Postgres")
}

func (ds *PostgresDatasource) startPoll() {
	log.Infof("Polling db every %d seconds", config.DatabasePoll)
	tick := time.NewTicker(time.Duration(config.DatabasePoll) * time.Second)

	go func() {
		for {
			select {
			case <-tick.C:
				ds.tenants, ds.timestamp = ds.loadTenants()
				for _, tenant := range ds.tenants {
					ds.tenantsByName[tenant.Name] = tenant
				}
				ds.updateChan <- ds.tenants
			}
		}
	}()
}

// loadTenantUpdates only loads the updates to the te
func (ds *PostgresDatasource) loadTenants() ([]*apimgmt.Tenant, time.Time) {
	// Just for demo - scan the table for the updated tenants, then load the updated ones....just makes
	// it quicker to PoC without having to merge and worry about threading
	var err error
	now := time.Now()
	tenantRows, err := ds.db.Query("SELECT DISTINCT(tenant_name), CURRENT_TIMESTAMP as now FROM proxy WHERE last_updated >= $1 AND last_updated < $2", ds.timestamp, now)
	if err != nil {
		log.Fatal(err)
	}

	var updatedTenants []string
	var timestamp time.Time

	for tenantRows.Next() {
		var tenantName string
		err = tenantRows.Scan(&tenantName, &timestamp)
		updatedTenants = append(updatedTenants, tenantName)
	}
	tenantRows.Close()

	apiTenants := []*apimgmt.Tenant{}

	for _, tenant := range updatedTenants {
		log.Infof("Loading tenant [%s] updates", tenant)
		rows, err := ds.db.Query("SELECT id, tenant_name, base_path, swagger FROM proxy WHERE tenant_name = $1", tenant)
		if err != nil {
			log.Fatal(err)
		}

		proxies := []*apimgmt.Proxy{}
		for rows.Next() {
			proxyRow := new(_ProxyRow)
			err = rows.Scan(&proxyRow.ID, &proxyRow.Name, &proxyRow.BasePath, &proxyRow.Swagger)
			if err != nil {
				log.Fatal(err)
			}

			host, port, basePath, tls := getBackendDetails(proxyRow)

			// Create the proxy
			proxy := &apimgmt.Proxy{
				Name: strconv.Itoa(proxyRow.ID),
				Frontend: &apimgmt.Frontend{
					BasePath: proxyRow.BasePath,
				},
				Backend: &apimgmt.Backend{
					Host: host,
					Path: basePath,
					Port: port,
					TLS:  tls,
				},
			}

			proxies = append(proxies, proxy)
		}
		rows.Close()

		apiTenants = append(apiTenants, &apimgmt.Tenant{
			Name:    tenant,
			Proxies: proxies,
		})
	}

	return apiTenants, now
}

func getBackendDetails(proxyRow *_ProxyRow) (string, uint32, string, bool) {
	// Brittle but it's a POC
	schemes := []interface{}{}
	tls := false
	port := uint32(80)
	host := "localhost"
	basePath := ""

	if proxyRow.Swagger["schemes"] != nil {
		schemes = proxyRow.Swagger["schemes"].([]interface{})
	}

	if len(schemes) > 0 {
		for _, scheme := range schemes {
			if strings.EqualFold(scheme.(string), "https") {
				tls = true
				port = 443
				break
			}
		}
	}

	if proxyRow.Swagger["host"] != nil {
		hostAndPort := strings.Split(proxyRow.Swagger["host"].(string), ":")
		host = hostAndPort[0]
		if len(hostAndPort) == 2 {
			p, err := strconv.Atoi(hostAndPort[1])
			if err != nil {
				log.Fatal(err)
			}
			port = uint32(p)
		}
	}

	if proxyRow.Swagger["basePath"] != nil {
		basePath = proxyRow.Swagger["basePath"].(string)
	}

	return host, port, basePath, tls
}

type _ProxyRow struct {
	ID       int
	Name     string
	BasePath string
	Swagger  _Attrs
	Created  time.Time
	Updated  time.Time
}

type _Attrs map[string]interface{}

func (a _Attrs) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *_Attrs) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
