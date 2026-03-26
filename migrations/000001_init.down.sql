DROP INDEX IF EXISTS idx_devices_hostname_trgm;
DROP INDEX IF EXISTS idx_devices_active_deleted_id;
DROP INDEX IF EXISTS idx_devices_is_active;

DROP TABLE IF EXISTS devices;