-- Create the default Poprako bot user directly in database.
INSERT INTO "t_user" (
    "id",
    "qid",
    "nickname",
    "password_hash",
    "is_super_admin"
) VALUES (
    'user-00000000-0000-0000-0000-000000000002',
    '666666',
    '白杨子Bot',
    '$2a$10$eEEkAsc7h3jdkOyjahdH6OX20w/dHKdGVaH7MNREkh54O57v.E2y2',
    TRUE
) ON CONFLICT (id) DO NOTHING;

-- Create the default Poprako bot member in the initial team.
INSERT INTO "t_member" (
    "id",
    "user_id",
    "team_id",
    "user_nickname",
    "assigned_admin_at"
) VALUES (
    'member-00000000-0000-0000-0000-000000000002',
    'user-00000000-0000-0000-0000-000000000002',
    'team-00000000-0000-0000-0000-000000000001',
    '白杨子Bot',
    NOW()
) ON CONFLICT (id) DO NOTHING;
