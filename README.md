ggVGy to copy all on vim

## Create database

CREATE DATABASE bandroom_cms;

## Create tables

CREATE TABLE courses (
id UUID PRIMARY KEY,
title TEXT NOT NULL,
description TEXT,
created_at TIMESTAMP DEFAULT now(),
updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE lessons (
id UUID PRIMARY KEY,
course_id UUID REFERENCES courses(id),
unit INT NOT NULL,
title TEXT NOT NULL,
description TEXT,
difficulty TEXT,
xp_reward INT,
crowns_reward INT,
created_at TIMESTAMP DEFAULT now(),
updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE questions (
id UUID PRIMARY KEY,
lesson_id UUID REFERENCES lessons(id),
question_text TEXT NOT NULL,
question_type TEXT NOT NULL,
image_url TEXT,
audio_url TEXT,
answer TEXT NOT NULL,
explanation TEXT,
created_at TIMESTAMP DEFAULT now(),
updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE options (
id UUID PRIMARY KEY,
question_id UUID REFERENCES questions(id),
option_text TEXT NOT NULL,
created_at TIMESTAMP DEFAULT now(),
updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE tags (
id UUID PRIMARY KEY,
name TEXT UNIQUE NOT NULL,
created_at TIMESTAMP DEFAULT now(),
updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE question_tags (
question_id UUID REFERENCES questions(id),
tag_id UUID REFERENCES tags(id),
PRIMARY KEY (question_id, tag_id),
created_at TIMESTAMP DEFAULT now()
);

git rm --cached .env

# Testing without token

curl -X POST http://localhost:8080/api/courses \
 -H "Content-Type: application/json" \
 -d '{"title": "Test Course", "description": "Sample"}'

# Testing with token

curl -X POST http://localhost:8080/api/courses \
 -H "Content-Type: application/json" \
 -H "Authorization: Bearer <your_token_here>" \
 -d '{"title": "Test Course", "description": "Sample"}'

# gcloud auth login

gcloud auth application-default login

hello
