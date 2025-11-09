INSERT INTO audit.log_levels (name) VALUES
    ('DEBUG'), ('INFO'), ('WARNING'), ('ERROR')
ON CONFLICT DO NOTHING;

INSERT INTO audit.modules (name) VALUES
    ('flag_sender'), ('config_loader'), ('exploit_runner'), ('jacfarm-api')
ON CONFLICT DO NOTHING;