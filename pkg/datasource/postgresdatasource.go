package datasource

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/bladedancer/envoyxds/pkg/apimgmt"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/jinzhu/copier"
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
		rows, err := ds.db.Query("SELECT id, tenant_name, base_path, auth, cred, swagger FROM proxy WHERE tenant_name = $1", tenant)
		if err != nil {
			log.Fatal(err)
		}

		proxies := []*apimgmt.Proxy{}
		for rows.Next() {
			proxyRow := new(_ProxyRow)
			err = rows.Scan(&proxyRow.ID, &proxyRow.Name, &proxyRow.BasePath, &proxyRow.Authorization, &proxyRow.Credential, &proxyRow.Swagger)
			if err != nil {
				log.Fatal(err)
			}

			feBasePath, authorizations := getFrontendDetails(proxyRow)
			host, port, beBasePath, tls, beAuthorization, beCredential := getBackendDetails(proxyRow)

			// Create the proxy
			proxy := &apimgmt.Proxy{
				Name: strconv.Itoa(proxyRow.ID),
				Frontend: &apimgmt.Frontend{
					BasePath:      feBasePath,
					Authorization: authorizations,
				},
				Backend: &apimgmt.Backend{
					Host:          host,
					Path:          beBasePath,
					Port:          port,
					TLS:           tls,
					Authorization: beAuthorization,
					Credential:    beCredential,
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

func getFrontendDetails(proxyRow *_ProxyRow) (string, []apimgmt.Authorization) {
	// TODO - support multiple frontend auth profiles.
	authorization := getFrontendAuthorization(proxyRow.Authorization)
	return proxyRow.BasePath, []apimgmt.Authorization{authorization}
}

func getBackendDetails(proxyRow *_ProxyRow) (string, uint32, string, bool, apimgmt.Authorization, map[string]string) {
	// Brittle but it's a POC
	tls := false
	port := uint32(80)
	host := "localhost"
	basePath := ""
	var authorization apimgmt.Authorization
	authorization = &apimgmt.PassthroughAuthorization{}

	if proxyRow.Swagger.Schemes != nil && len(proxyRow.Swagger.Schemes) > 0 {
		for _, scheme := range proxyRow.Swagger.Schemes {
			if strings.EqualFold(scheme, "https") {
				tls = true
				port = 443
				break
			}
		}
	}

	if len(proxyRow.Swagger.Host) > 0 {
		hostAndPort := strings.Split(proxyRow.Swagger.Host, ":")
		host = hostAndPort[0]
		if len(hostAndPort) == 2 {
			p, err := strconv.Atoi(hostAndPort[1])
			if err != nil {
				log.Fatal(err)
			}
			port = uint32(p)
		}
	}

	if len(proxyRow.Swagger.BasePath) > 0 {
		basePath = proxyRow.Swagger.BasePath
	}

	if len(proxyRow.Swagger.SwaggerProps.SecurityDefinitions) > 0 {
		// This is a huge hack. We don't have route per operations in the PoC
		// and we don't support AND/OR on the security scheme.....it's just a PoC.
		var scheme *spec.SecurityScheme
		for _, s := range proxyRow.Swagger.SwaggerProps.SecurityDefinitions {
			// Just grabbing the first security scheme in the definition as the default....
			if s.Type != "oauth2" {
				scheme = s
				break
			}
		}

		// If there is a global security definition then use it....kinda, again not
		// supporting AND/OR....it's all very dodgy.
		if len(proxyRow.Swagger.SwaggerProps.Security) > 0 {
			// Really taking short cuts with the spec here...but it's just a PoC,
			// No AND/OR support - taking the one scheme and applying across the board
			for security := range proxyRow.Swagger.SwaggerProps.Security[0] {
				s := proxyRow.Swagger.SwaggerProps.SecurityDefinitions[security]
				// Just grabbing the first security scheme
				if s.Type != "oauth2" {
					scheme = s
					break
				}
			}
		}

		if scheme != nil {
			switch scheme.Type {
			case "apiKey":
				authorization = &apimgmt.APIKeyAuthorization{
					Name:     scheme.Name,
					Location: scheme.In,
				}
			case "basic":
				authorization = &apimgmt.HTTPAuthorization{
					Scheme: "basic",
				}
			case "oauth2":
				log.Error("Didn't implement oauth2 for backends, try next one....")
			}
		}
	}

	// Backend Credential
	beCredential := make(map[string]string)
	if proxyRow.Credential != nil && len(proxyRow.Credential) > 0 {
		for key, value := range proxyRow.Credential {
			beCredential[key] = value.(string)
		}
	}

	return host, port, basePath, tls, authorization, beCredential
}

// getFrontendAuthorization Convert the auth details to an Authorization
func getFrontendAuthorization(authAttrs _Attrs) apimgmt.Authorization {
	if authAttrs == nil && len(authAttrs) == 0 {
		return &apimgmt.PassthroughAuthorization{}
	}

	var auth apimgmt.Authorization
	switch (apimgmt.AuthorizationType)(authAttrs["type"].(string)) {
	case apimgmt.AuthorizationTypePassthrough:
		auth = &apimgmt.PassthroughAuthorization{}
	case apimgmt.AuthorizationTypeAPIKey:
		auth = &apimgmt.APIKeyAuthorization{
			Name:     authAttrs["name"].(string),
			Location: authAttrs["in"].(string),
		}
	case apimgmt.AuthorizationTypeHTTP:
		auth = &apimgmt.HTTPAuthorization{
			Scheme: authAttrs["scheme"].(string),
		}
	}

	return auth
}

type _ProxyRow struct {
	ID            int
	Name          string
	BasePath      string
	Authorization _Attrs
	Credential    _Attrs
	Swagger       *_Swagger
	Created       time.Time
	Updated       time.Time
}

type _Attrs map[string]interface{}

func (a _Attrs) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *_Attrs) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(b, &a)
}

type _Swagger spec.Swagger

func (s *_Swagger) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data := value.([]byte)
	doc, err := loads.Analyzed(json.RawMessage(data), "")
	if err != nil {
		return err
	}
	copier.Copy(s, doc.Spec())
	return err
}
