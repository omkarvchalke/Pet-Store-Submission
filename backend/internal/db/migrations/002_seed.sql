INSERT INTO stores (slug, name)
VALUES
  ('happy-paws', 'Happy Paws'),
  ('frog-heaven', 'Frog Heaven')
ON CONFLICT (slug) DO NOTHING;

INSERT INTO users (username, password_hash, role, store_id)
SELECT 'merchant_happy', '$2a$10$4jNZm4RVEZ53rvePBeJbwuycrPshqrndcrm0thkSjQNgnKxF1QXty', 'merchant', s.id
FROM stores s WHERE s.slug = 'happy-paws'
ON CONFLICT (username) DO NOTHING;

INSERT INTO users (username, password_hash, role, store_id)
SELECT 'merchant_frog', '$2a$10$4jNZm4RVEZ53rvePBeJbwuycrPshqrndcrm0thkSjQNgnKxF1QXty', 'merchant', s.id
FROM stores s WHERE s.slug = 'frog-heaven'
ON CONFLICT (username) DO NOTHING;

INSERT INTO users (username, password_hash, role, store_id)
VALUES
  ('customer_happy', '$2a$10$Pikbb.QJacgydqWFFJH5cejCncVNOMUIEX88jm3t48KYne2OSt.Z.', 'customer', NULL),
  ('customer_frog', '$2a$10$Pikbb.QJacgydqWFFJH5cejCncVNOMUIEX88jm3t48KYne2OSt.Z.', 'customer', NULL)
ON CONFLICT (username) DO NOTHING;
