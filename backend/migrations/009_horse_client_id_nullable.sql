-- Make client_id nullable on horses so owners can create horses without a client
ALTER TABLE horses ALTER COLUMN client_id DROP NOT NULL;
ALTER TABLE horses DROP CONSTRAINT IF EXISTS horses_client_id_fkey;
ALTER TABLE horses ADD CONSTRAINT horses_client_id_fkey FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE SET NULL;
