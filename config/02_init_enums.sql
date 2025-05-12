INSERT INTO cities
    (uuid, name)
VALUES 
    ('9e5e8171-3453-4e89-b09f-3883b5b1594c', 'Москва'),
    ('965459bf-863d-4a9f-931d-de0d9ada6f69', 'Санкт-Петербург'),
    ('487454ed-a9f3-49ad-b6f5-17ac0aedb5df', 'Казань')
ON CONFLICT (uuid) DO NOTHING;

INSERT INTO product_types 
    (uuid, name)
VALUES 
    ('acbb48fb-964b-4ccb-9295-b007bd001c2f', 'электроника'),
    ('b89618e4-3e2d-44f5-a700-7262fa5aa39e', 'одежда'),
    ('0f06c601-e37c-4b3e-a752-3d8e8f44f9eb', 'обувь')
ON CONFLICT (uuid) DO NOTHING;

INSERT INTO roles 
    (uuid, name)
VALUES 
    ('a516b6f9-8d5a-419c-9494-95bddf98e151', 'employee'),
    ('699d1b52-09cc-4b6a-ae8f-9986a1035218', 'moderator')
ON CONFLICT (uuid) DO NOTHING;

INSERT INTO statuses 
    (uuid, name)
VALUES 
    ('24b24a54-9398-46f7-b2a0-74c07b1bb77d', 'in_progress'),
    ('1370bddc-51d0-45c0-8399-122364d20512', 'close')
ON CONFLICT (uuid) DO NOTHING;