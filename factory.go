package gormer

var instance *Manager

func SetInstance(m *Manager) {
	instance = m
}

func Factory() *Manager {
	return instance
}

func Gormer() *Manager {
	return Factory()
}

func Connect(name ...string) *DB {
	return Factory().Connect(name...)
}
