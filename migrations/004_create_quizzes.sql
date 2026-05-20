CREATE TABLE quizzes (
    id UUID PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    created_by UUID NOT NULL
);

CREATE TABLE questions (
    id UUID PRIMARY KEY,
    quiz_id UUID NOT NULL REFERENCES quizzes(id),
    text VARCHAR(100) NOT NULL,
    order_num INTEGER NOT NULL,
    correct_option_id UUID NOT NULL REFERENCES options(id)
);

CREATE TABLE options (
    id UUID PRIMARY KEY,
    question_id UUID NOT NULL REFERENCES questions(id),
    text VARCHAR(100) NOT NULL,
    is_correct BOOLEAN NOT NULL
);