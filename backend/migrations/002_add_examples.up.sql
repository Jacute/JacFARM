INSERT INTO exploits (id, name, type, is_local, executable_path, requirements_path)
VALUES (
    '6e7266cc-2405-417f-9410-a028f39cbedd',
    'Example Exploit',
    'python',
    1,
    './exploits/6e7266cc-2405-417f-9410-a028f39cbedd/example_exploit1.py',
    './exploits/6e7266cc-2405-417f-9410-a028f39cbedd/requirements.txt'
), (
    '2890535c-f720-4704-aca5-4e96f2a31b1b',
    'Example Exploit 2',
    'binary',
    1,
    './exploits/2890535c-f720-4704-aca5-4e96f2a31b1b/exploit2',
    NULL
), (
    'cf66fdab-3d46-4b15-9e6b-deb7a18b85c6',
    'Example Exploit 3',
    'binary',
    1,
    './exploits/cf66fdab-3d46-4b15-9e6b-deb7a18b85c6/main',
    NULL
);