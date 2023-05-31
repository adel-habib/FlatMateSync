
-- This table represents the users in the application. 
-- Each user can be part of multiple flats. A user might have logged in using OIDC SSO
-- For users who sign up with OIDC the password_hash field remains null and the username column is the `preferred_username` defined in "OpenID Connect Core 1.0 incorporating errata set 1"
CREATE TABLE users (
  username VARCHAR(255) PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255),
  oidc_id VARCHAR(255),
  oidc_provider VARCHAR(255),
  deleted_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Flats table: 
-- This table represents the shared flats. Each flat can have multiple users.
CREATE TABLE flats (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  deleted_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- UserFlats table: 
-- This is a junction table which connects users with flats. It also contains user balance in the flat. 
-- The balance can be positive or negative depending on the user's total payments and shares in the flat.
CREATE TABLE user_flats (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id),
  flat_id INTEGER NOT NULL REFERENCES flats(id),
  is_admin BOOLEAN NOT NULL DEFAULT FALSE,
  balance FLOAT NOT NULL DEFAULT 0,
  deleted_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Payments table: 
-- This table represents payments made by the users in a flat. 
-- Each payment is associated with a flat and a user who made the payment.
CREATE TABLE payments (
  id SERIAL PRIMARY KEY,
  amount FLOAT NOT NULL,
  description TEXT,
  payer_id INTEGER NOT NULL REFERENCES users(id),
  flat_id INTEGER NOT NULL REFERENCES flats(id),
  payment_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- PaymentParticipants table: 
-- This is a junction table which connects payments with users. 
-- Each entry indicates a user's share in a particular payment.
CREATE TABLE payment_participants (
  id SERIAL PRIMARY KEY,
  payment_id INTEGER NOT NULL REFERENCES payments(id),
  participant_id INTEGER NOT NULL REFERENCES users(id),
  deleted_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
