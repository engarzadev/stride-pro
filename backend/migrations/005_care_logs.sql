CREATE TABLE IF NOT EXISTS care_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    horse_id UUID NOT NULL REFERENCES horses(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    category VARCHAR(100) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_care_logs_user_id ON care_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_care_logs_horse_id ON care_logs(horse_id);
CREATE INDEX IF NOT EXISTS idx_care_logs_date ON care_logs(date);
