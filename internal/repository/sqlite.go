package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gohome/internal/core"
	"slices"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDeviceRepository struct {
	db *sql.DB
}

func SetupSQLiteDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func NewSQLiteDeviceRepository(db *sql.DB) (*SQLiteDeviceRepository, error) {
	query := `
	CREATE TABLE IF NOT EXISTS devices (
		id TEXT PRIMARY KEY,
		address TEXT,
		address_type TEXT,
		name TEXT,
		protocol TEXT,
		adapter_ids TEXT,
		created_at DATETIME,
		capabilities TEXT,
		last_updated DATETIME
	);
	CREATE INDEX IF NOT EXISTS idx_device_address ON devices(address, address_type);
	`
	_, err := db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("failed to create devices table: %w", err)
	}

	return &SQLiteDeviceRepository{db: db}, nil
}

func (r *SQLiteDeviceRepository) Save(device *core.Device) error {
	adapterIDsJson, err := json.Marshal(device.AdapterIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal adapter_ids: %w", err)
	}

	capabilitiesJson, err := json.Marshal(device.Capabilities)
	if err != nil {
		return fmt.Errorf("failed to marshal capabilities: %w", err)
	}

	query := `
	INSERT OR REPLACE INTO devices 
	(id, address, address_type, name, protocol, adapter_ids, created_at, capabilities, last_updated)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.Exec(query,
		device.ID,
		device.Address,
		device.AddressType,
		device.Name,
		device.Protocol,
		string(adapterIDsJson),
		device.CreatedAt,
		string(capabilitiesJson),
		device.LastUpdated,
	)

	if err != nil {
		return fmt.Errorf("failed to save device: %w", err)
	}

	return nil
}

func (r *SQLiteDeviceRepository) FindByID(id string) (*core.Device, error) {
	query := `SELECT id, address, address_type, name, protocol, adapter_ids, created_at, capabilities, last_updated FROM devices WHERE id = ?`

	row := r.db.QueryRow(query, id)
	return r.scanDevice(row)
}

func (r *SQLiteDeviceRepository) FindAll() ([]*core.Device, error) {
	query := `SELECT id, address, address_type, name, protocol, adapter_ids, created_at, capabilities, last_updated FROM devices`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []*core.Device
	for rows.Next() {
		device, err := r.scanDevice(rows)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func (r *SQLiteDeviceRepository) FindByAddress(address string, addressType core.AddressType) (*core.Device, error) {
	query := `SELECT id, address, address_type, name, protocol, adapter_ids, created_at, capabilities, last_updated FROM devices WHERE address = ? AND address_type = ?`

	row := r.db.QueryRow(query, address, addressType)
	return r.scanDevice(row)
}

func (r *SQLiteDeviceRepository) LinkAdapter(deviceID, adapterID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	device, err := r.FindByID(deviceID)
	if err != nil {
		return err
	}
	if device == nil {
		return nil // Device not found
	}

	if slices.Contains(device.AdapterIDs, adapterID) {
		return nil
	}

	device.AdapterIDs = append(device.AdapterIDs, adapterID)

	adapterIDsJson, _ := json.Marshal(device.AdapterIDs)

	_, err = tx.Exec(`UPDATE devices SET adapter_ids = ? WHERE id = ?`, string(adapterIDsJson), deviceID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *SQLiteDeviceRepository) UnlinkAdapter(deviceID, adapterID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	device, err := r.FindByID(deviceID)
	if err != nil {
		return err
	}
	if device == nil {
		return nil
	}

	device.AdapterIDs = slices.DeleteFunc(device.AdapterIDs, func(e string) bool {
		return e == adapterID
	})

	adapterIDsJson, _ := json.Marshal(device.AdapterIDs)

	_, err = tx.Exec(`UPDATE devices SET adapter_ids = ? WHERE id = ?`, string(adapterIDsJson), deviceID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

type Scanner interface {
	Scan(dest ...any) error
}

func (r *SQLiteDeviceRepository) scanDevice(row Scanner) (*core.Device, error) {
	var d core.Device
	var adapterIDsJson []byte
	var capabilitiesJson []byte
	var addressType string

	err := row.Scan(
		&d.ID,
		&d.Address,
		&addressType,
		&d.Name,
		&d.Protocol,
		&adapterIDsJson,
		&d.CreatedAt,
		&capabilitiesJson,
		&d.LastUpdated,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	d.AddressType = core.AddressType(addressType)

	if len(adapterIDsJson) > 0 {
		if err := json.Unmarshal(adapterIDsJson, &d.AdapterIDs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal adapter_ids: %w", err)
		}
	}

	if len(capabilitiesJson) > 0 {
		if err := json.Unmarshal(capabilitiesJson, &d.Capabilities); err != nil {
			return nil, fmt.Errorf("failed to unmarshal capabilities: %w", err)
		}
	}

	if d.Capabilities == nil {
		d.Capabilities = make(map[core.CapabilityType]*core.Capability)
	}

	return &d, nil
}
