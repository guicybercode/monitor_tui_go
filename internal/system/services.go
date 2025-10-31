package system

import (
	"context"
	"errors"

	"github.com/coreos/go-systemd/v22/dbus"
)

type ServiceInfo struct {
	Name        string
	Description string
	State       string
	ActiveState string
	SubState    string
	LoadState   string
}

func GetServices() ([]ServiceInfo, error) {
	conn, err := dbus.NewSystemConnectionContext(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	units, err := conn.ListUnitsContext(context.Background())
	if err != nil {
		return nil, err
	}

	var services []ServiceInfo
	for _, unit := range units {
		if len(unit.Name) == 0 {
			continue
		}

		services = append(services, ServiceInfo{
			Name:        unit.Name,
			Description: unit.Description,
			State:       string(unit.ActiveState),
			ActiveState: string(unit.ActiveState),
			SubState:    unit.SubState,
			LoadState:   string(unit.LoadState),
		})
	}

	return services, nil
}

func StartService(name string) error {
	conn, err := dbus.NewSystemConnectionContext(context.Background())
	if err != nil {
		return err
	}
	defer conn.Close()

	ch := make(chan string)
	_, err = conn.StartUnitContext(context.Background(), name, "replace", ch)
	if err != nil {
		return err
	}

	result := <-ch
	if result != "done" {
		return errors.New("failed to start service: " + result)
	}

	return nil
}

func StopService(name string) error {
	conn, err := dbus.NewSystemConnectionContext(context.Background())
	if err != nil {
		return err
	}
	defer conn.Close()

	ch := make(chan string)
	_, err = conn.StopUnitContext(context.Background(), name, "replace", ch)
	if err != nil {
		return err
	}

	result := <-ch
	if result != "done" {
		return errors.New("failed to stop service: " + result)
	}

	return nil
}

func RestartService(name string) error {
	conn, err := dbus.NewSystemConnectionContext(context.Background())
	if err != nil {
		return err
	}
	defer conn.Close()

	ch := make(chan string)
	_, err = conn.RestartUnitContext(context.Background(), name, "replace", ch)
	if err != nil {
		return err
	}

	result := <-ch
	if result != "done" {
		return errors.New("failed to restart service: " + result)
	}

	return nil
}
