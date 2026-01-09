# Cake Order Tracker

A lightweight, robust order management system built with Go and SQLite. This project focuses on a clean monolithic architecture, transactional data integrity, and a high-performance, single-binary deployment model.

## Key Features

* **Order Lifecycle Management:** Full CRUD operations for products and customer orders with SQL transactions to ensure data consistency.
* **Live Asynchronous Search:** Real-time order filtering by date range utilizing **HTMX** for partial DOM updates and reduced server payload.
* **Custom Authentication:** Secure session-based authentication featuring salted password hashing and middleware-based route protection.
* **Embedded Assets:** Single-binary deployment utilizing Go `embed` for HTML templates, static assets (HTMX, FontAwesome), and SQL migrations.
* **Interactive Client-Side Logic:** Dynamic order composition with real-time total price calculation using vanilla JavaScript and HTML5 templates.

## Tech Stack

* **Backend:** Go (Golang)
* **Database:** SQLite
* **Frontend:** HTMX, Go HTML Templates, CSS, JavaScript
* **Tooling:** Versioned SQL Migrations, Go Toolchain

## Architectural Highlights

### 1. Database Persistence

The system implements a versioned migration system. The schema handles relationships between orders and products via join tables, ensuring integrity through the `database/sql` package's transaction support.

### 2. Separation of Concerns

* **/models**: Pure data structures for cross-package communication.
* **/store**: Encapsulated data persistence layer containing all SQL logic and migrations.
* **/server**: HTTP routing, handler logic, and multi-entry template rendering.
* **/auth**: Dedicated security layer providing middleware for session validation.

## Installation and Running

1. **Run the application:**
```bash
go run main.go
```
