CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    priority VARCHAR(20) NOT NULL DEFAULT 'Medium',
    CHECK(status IN('pending', 'in_progress', 'completed', 'cancelled')),
    CHECK(priority IN('Low', 'Medium', 'High')),
    due_date TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX i_tasks_user_id ON tasks(user_id);

CREATE TABLE schedule_slots(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    subject TEXT NOT NULL,
    teacher TEXT NOT NULL,
    room VARCHAR(20) NOT NULL,
    day_of_week INTEGER NOT NULL,
    CHECK(day_of_week BETWEEN 1 AND 7),
    start_time VARCHAR(20) NOT NULL,
    end_time VARCHAR(20) NOT NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'lecture',
    CHECK(type IN('lecture', 'practice', 'lab')),
    parity VARCHAR(20) NOT NULL DEFAULT 'even',
    CHECK(parity IN('even', 'odd', 'both')),
    is_cancelled BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX i_schedule_user_day ON schedule_slots(user_id, day_of_week);

CREATE TABLE proposed_changes(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    slot_id INTEGER REFERENCES schedule_slots(id) ON DELETE CASCADE,
    action VARCHAR(20) NOT NULL DEFAULT 'add_slot',
    CHECK(action IN('cancel', 'reschedule_time', 'change_room', 'add_slot')),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    CHECK(status IN('pending', 'approved', 'rejected')),
    raw_text TEXT,
    new_start_time VARCHAR(20) ,
    new_end_time VARCHAR(20) ,
    new_room TEXT,
    created_add TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX i_change_user_day ON proposed_changes(user_id, status);