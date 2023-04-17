CREATE TABLE like (
    id      INTEGER PRIMARY KEY,
    post_id INTEGER REFERENCES post (id),
    user_id INTEGER REFERENCES user (id),
    type_of TEXT
);