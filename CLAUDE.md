# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## High-level Architecture

This project is a "选课通系统" (Course Selection System) - a full-stack web application with a Go backend, Vue 3 frontend, and AdminLTE admin dashboard. It's designed as a course management platform where students can browse courses, rate them, and leave comments, while administrators manage courses and users.

**Architecture Overview:**
- **Backend**: Go RESTful API using Gin framework with GORM for database operations
- **Frontend**: Vue 3 + TypeScript application using Vite build system and Composition API
- **Admin Panel**: AdminLTE 4.0 (Bootstrap 5) for administrative operations
- **Database**: SQLite for development (backend/test.db), with MySQL support for production

**Key Features:**
- JWT-based user authentication
- Course rating system (1-5 stars)
- Course comments and reviews
- Multi-layer filtering (grade, semester, subject)
- Glassmorphism UI design patterns
- Responsive design across all interfaces

## Development Commands

### Backend (Go)
All backend commands should be run from the `backend/` directory.

- **Install dependencies**:
  ```bash
  go mod tidy
  ```

- **Run development server**:
  ```bash
  go run main.go
  ```
  *Runs on http://localhost:8080 with automatic route registration*

- **Database seeding** (for testing):
  ```bash
  # After starting the server, visit:
  curl http://localhost:8080/api/v1/seed
  ```

- **Database migration** (automatic):
  *GORM AutoMigrate handles database schema on startup for models: User, Course, Rating, Comment*

### Frontend (Vue + TypeScript)
All frontend commands should be run from the `frontend/` directory.

- **Install dependencies**:
  ```bash
  npm install
  ```

- **Run development server**:
  ```bash
  npm run dev
  ```
  *Runs on http://localhost:5173 with hot reload*

- **Build for production**:
  ```bash
  npm run build
  ```

- **Type checking**:
  ```bash
  npm run type-check
  ```

- **Lint and format**:
  ```bash
  npm run lint          # ESLint with auto-fix
  npm run format        # Prettier formatting
  ```

- **Testing**:
  ```bash
  npm run test:unit     # Vitest unit tests
  npx playwright install # Install Playwright (once)
  npm run test:e2e      # End-to-end tests
  ```

### Admin Panel (AdminLTE)
All admin panel commands should be run from the `AdminLTE/` directory.

- **Install dependencies**:
  ```bash
  npm install
  ```

- **Run development server**:
  ```bash
  npm start
  ```
  *Runs on http://localhost:3000 with live reload*

- **Build for production**:
  ```bash
  npm run production
  ```

- **Code quality**:
  ```bash
  npm run lint          # Comprehensive linting (JS, CSS, docs)
  ```

## Architecture Details

### Backend Structure
- **Controllers**: Handle HTTP requests and business logic
  - `auth_controller.go`: JWT authentication (register/login)
  - `course_controller.go`: Course CRUD with rating calculations
  - `rating_controller.go`: Course rating operations
  - `comment_controller.go`: Comment management
  - `admin_controller.go`: Administrative operations
  - `user_controller.go`: User management

- **Models**: GORM data models (User, Course, Rating, Comment)
- **Routes**: API route definitions grouped by functionality
- **Config**: Database connection and configuration

**Key API Endpoints:**
- `GET /api/v1/courses` - Get courses with ratings (supports filters: grade, semester, subject)
- `POST /api/v1/ratings` - Submit course rating
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration
- `GET /api/v1/admin/*` - Admin-only endpoints

### Frontend Structure
- **Services**: API communication layer in `src/services/api.ts`
  - `courseService`: Course data fetching with rating integration
  - `ratingService`: Rating submission and user rating tracking

- **Views**: Page components using Vue 3 Composition API
  - `HomeView.vue`: Main course listing with filters
  - `CourseDetailView.vue`: Individual course details
  - `CourseRateView.vue`: Course rating interface
  - `CourseManagement.vue`: Admin course management
  - `AuthView.vue`: Authentication forms

- **State Management**: Pinia stores (if implemented)
- **Routing**: Vue Router 4 with lazy-loaded routes

### Key Design Patterns
- **API Service Layer**: Centralized API communication with axios interceptors for JWT auth
- **Rating Integration**: Course list API includes average ratings and total counts
- **Glassmorphism UI**: Modern UI with backdrop-filter and transparency effects
- **Type Safety**: Full TypeScript integration with strict type checking

### Database Schema
The system uses four main models:
- **User**: Authentication and profile information
- **Course**: Course metadata (name, teacher, credits, etc.)
- **Rating**: 1-5 star ratings by users for courses
- **Comment**: Text comments by users on courses

### Development Workflow
1. **Backend First**: Start Go server (`go run main.go`) - handles database migrations automatically
2. **Frontend Second**: Start Vue dev server (`npm run dev`) - connects to backend API
3. **Admin Panel**: Start AdminLTE (`npm start`) - separate admin interface
4. **Testing**: Use the seed endpoint (`/api/v1/seed`) to populate test data

**Port Configuration:**
- Backend API: http://localhost:8080
- Frontend App: http://localhost:5173  
- Admin Panel: http://localhost:3000

**Environment Notes:**
- Node.js 20+ required for frontend
- Go 1.24+ required for backend
- Development uses SQLite, production supports MySQL
- CORS is configured for all origins in development
