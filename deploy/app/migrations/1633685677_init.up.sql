CREATE TABLE gui_config
(
    window_width  INT NOT NULL       DEFAULT 1024,
    window_height INT NOT NULL       DEFAULT 768,
    window_x      INT NOT NULL       DEFAULT 50,
    window_y      INT NOT NULL       DEFAULT 50,
    created_at    DATETIME DEFAULT (datetime('now', 'localtime')) NOT NULL,
    updated_at    DATETIME DEFAULT (datetime('now', 'localtime')) NOT NULL
);

INSERT INTO gui_config (window_width, window_height, window_x, window_y)
VALUES (1024, 768, 50, 50);

CREATE TABLE workspace
(
    id         INTEGER PRIMARY KEY,
    parent_id  INTEGER                              NULL,
    has_child  BOOL                                 NOT NULL DEFAULT FALSE,
    type       TEXT CHECK (type IN ('f', 's', 'r')) NOT NULL,
    title      TEXT,
    data       TEXT                                 NULL,
    sort       INTEGER                              NOT NULL DEFAULT 0,
    expanded   BOOL                                 NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DATETIME                            DEFAULT (datetime('now', 'localtime')) NOT NULL,
    updated_at TIMESTAMP DATETIME                            DEFAULT (datetime('now', 'localtime')) NOT NULL
);