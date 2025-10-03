INSERT INTO audit.log_levels (name) VALUES
    ('INFO'), ('WARNING'), ('ERROR')
ON CONFLICT DO NOTHING;

INSERT INTO audit.services (name) VALUES
    ('flag_sender'), ('config_loader'), ('exploit_runner'), ('jacfarm-api')
ON CONFLICT DO NOTHING;