ALTER TABLE gui_config
    RENAME TO settings;

ALTER TABLE settings
    ADD COLUMN single_instance BOOL NOT NULL DEFAULT TRUE;
ALTER TABLE settings
    ADD COLUMN connect_timeout INT NOT NULL DEFAULT 10;
ALTER TABLE settings
    ADD COLUMN request_timeout INT NOT NULL DEFAULT 30;
ALTER TABLE settings
    ADD COLUMN non_blocking_connection BOOL NOT NULL DEFAULT TRUE;
ALTER TABLE settings
    ADD COLUMN sort_methods_by_name BOOL NOT NULL DEFAULT TRUE;
ALTER TABLE settings
    ADD COLUMN max_loop_depth INT NOT NULL DEFAULT 10;