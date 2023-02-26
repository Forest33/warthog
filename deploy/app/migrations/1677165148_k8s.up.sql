ALTER TABLE settings
    ADD COLUMN k8s_request_timeout INT NOT NULL DEFAULT 30;
