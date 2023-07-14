package env

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Manager struct {
	m map[string]string
}

func NewManager() *Manager {
	return &Manager{
		m: getEnvMap(),
	}
}

func (m *Manager) GetServiceIds() []string {
	services := make([]string, 0)
	r := regexp.MustCompile(getRegEx(TypeTpl))
	for k := range m.m {
		if r.Match([]byte(k)) {
			matches := r.FindStringSubmatch(k)
			services = append(services, matches[1])
		}
	}
	return services
}

func (m *Manager) GetGlobalEnv(key string) string {
	return m.m[key]
}

func (m *Manager) GetGlobalEnvWithDefault(key, defaultVal string) string {
	val := m.m[key]
	if val != "" {
		return val
	}
	return defaultVal
}

func (m *Manager) GetGlobalEnvWithNilCheck(key string) string {
	val := m.m[key]
	if val != "" {
		return val
	}
	panic(errors.New(fmt.Sprintf("'%s' must not be nil", key)))
}

func (m *Manager) GetEnv(tpl, id string) string {
	return m.GetGlobalEnv(fmt.Sprintf(tpl, id))
}

func (m *Manager) GetEnvWithDefault(tpl, id, defaultVal string) string {
	return m.GetGlobalEnvWithDefault(fmt.Sprintf(tpl, id), defaultVal)
}

func (m *Manager) GetEnvWithNilCheck(tpl, id string) string {
	return m.GetGlobalEnvWithNilCheck(fmt.Sprintf(tpl, id))
}

func getRegEx(tpl string) string {
	return fmt.Sprintf(tpl, ServiceIdRegEx)
}

func getEnvMap() map[string]string {
	m := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if strings.HasPrefix(pair[0], Prefix) {
			m[pair[0]] = pair[1]
		}
	}
	return m
}
