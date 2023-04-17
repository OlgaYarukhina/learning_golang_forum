CREATE TABLE categoryPostRelation (
    id          INTEGER PRIMARY KEY,
    post_id     INTEGER REFERENCES post (id),
    category_id INTEGER REFERENCES category (id)
);