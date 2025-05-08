# AI ORM Schema Transformer

**AI ORM Schema Transformer** is a simple CLI tool that reads your SQL database schema directly from a database URL and automatically generates ORM model definitions for popular backend frameworks and languages — all powered by AI.

Supports:
- **Go** with GORM
- **Rust** with Diesel
- **TypeScript** with TypeORM (NestJS)
- **Python** with SQLAlchemy (FastAPI)
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

