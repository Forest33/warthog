ALTER TABLE settings
    DROP COLUMN single_instance;
ALTER TABLE settings
    DROP COLUMN connect_timeout;
ALTER TABLE settings
    DROP COLUMN request_timeout;
ALTER TABLE settings
    DROP COLUMN non_blocking_connection;
ALTER TABLE settings
    DROP COLUMN sort_methods_by_name;
ALTER TABLE settings
    DROP COLUMN max_loop_depth;

ALTER TABLE settings RENAME TO gui_config;