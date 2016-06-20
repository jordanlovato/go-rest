package main

type ConfigMap struct {
	Database map[string]string
}

type Log struct {
	Date      string
	Firstname string
	Lastname  string
	Type      string
}

type ok interface {
	OK() error
}

func (l *Log) OK() error {
	// Basic validation
	if len(l.Firstname) == 0 {
		return ErrRequired("Firstname")
	}

	if len(l.Lastname) == 0 {
		return ErrRequired("Lastname")
	}

	return nil
}
