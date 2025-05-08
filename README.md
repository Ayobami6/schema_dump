# AI ORM Schema Transformer

**AI ORM Schema Transformer** is a simple CLI tool that reads your SQL database schema directly from a database URL and automatically generates ORM model definitions for popular backend frameworks and languages — all powered by AI.

Supports:
- **Go** with GORM
- **Rust** with Diesel
- **TypeScript** with TypeORM (NestJS)
- **Python** with SQLAlchemy (FastAPI) and Django ORM
- **Java** with Spring Boot (JDBC)

---

## 🚀 Project Overview

Modern backend developers often waste valuable time manually converting database schemas to ORM models for different stacks. This CLI tool streamlines that process by using AI to read the structure of your tables and produce ORM code for your chosen backend language and framework.

This tool is great for:
- Rapid prototyping
- Cross-language development
- AI-assisted boilerplate generation

---

## 🛠️ What It Does

Given a database URL and a specific table, this tool will:
1. Connect to your database.
2. Read the schema for the specified table.
3. Generate ORM-compatible model definitions using AI.
4. Output the model to your desired format/language.

---

## 📦 Installation
```sh
wget https://raw.githubusercontent.com/Ayobami6/schema_dump/refs/heads/master/install.sh
```

```sh
chmod 711 install.sh
./install.sh
```

## 📼 Demo


https://github.com/user-attachments/assets/e75ec66a-a487-4ba3-b03b-63cf712376bd



## Commands
- dump-schema: Dumps SQL schema from a live database to a file Flags: --db (database type), --url (connection URL), --table (optional table name)
- list-tables: Lists all tables in the database and outputs to a JSON file Flags: --db (database type), --url (connection URL)
- transform: Transforms SQL schema to ORM models for various languages Flags: --db (database type), --url (connection URL), --table (table name), --lang (target language)
Supported languages: py, ts, java, rs, go

## Usage
Usage examples:

Dump schema for a particular table from a PostgreSQL database to a `schema.sql` in th current directory:
```
schema dump-schema --db postgres --url "postgresql://user:password@localhost:5432/dbname" --table users
```
Dump schema for all tables from a PostgreSQL database to a `schema.sql` in th current directory:
```
schema dump-schema --db postgres --url "postgresql://user:password@localhost:5432/dbname"
```

List all tables in a database to a `tables.json` in the current directory:
```
schema list-tables --db postgres --url "postgresql://user:password@localhost:5432/dbname"
```
Transform a table schema to an ORM model:
```
schema transform --db postgres --url "postgresql://user:password@localhost:5432/dbname" --table users --lang py
```
All commands require the --db and --url flags. The transform command additionally requires --table and --lang flags.