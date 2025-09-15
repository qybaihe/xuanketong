# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## High-level Architecture

This project is a "选课通系统" (Course Selection System) with a Go backend and a Vue 3 frontend. It also includes the AdminLTE Bootstrap 5 admin dashboard template.

- **Backend**: A RESTful API written in Go, located in the `backend/` directory. It handles business logic, database interactions, and user authentication.
- **Frontend**: A Vue 3 application built with Vite, located in the `frontend/` directory. This is the main user-facing application for students.
- **Admin Panel**: The `AdminLTE/` directory contains a separate frontend for administrators, built with Bootstrap 5.

## Backend (Go)

The backend is a standard Go application.

- **To run the backend development server**:
  ```bash
  cd backend
  go run main.go
  ```

## Frontend (Vue + Vite)

The frontend is a Vue 3 application using Vite. All commands should be run from the `frontend/` directory.

- **Install dependencies**:
  ```bash
  npm install
  ```
- **Run development server**:
  ```bash
  npm run dev
  ```
- **Build for production**:
  ```bash
  npm run build
  ```
- **Run unit tests**:
  ```bash
  npm run test:unit
  ```
- **Run end-to-end tests**:
  ```bash
  npx playwright install # (run only once)
  npm run build
  npm run test:e2e
  ```
- **Lint files**:
  ```bash
  npm run lint
  ```

## Admin Panel (AdminLTE)

The admin panel is based on the AdminLTE template. All commands should be run from the `AdminLTE/` directory.

- **Install dependencies**:
  ```bash
  npm install
  ```
- **Run development server**:
  ```bash
  npm start
  ```
- **Build for production**:
  ```bash
  npm run production
  ```
- **Lint files**:
  ```bash
  npm run lint
  ```
