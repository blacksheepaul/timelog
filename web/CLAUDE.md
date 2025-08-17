# CLAUDE.md - TimeLog Web Frontend

This file provides guidance to Claude Code when working with the Vue 3 frontend of the TimeLog application.

## Frontend Overview

A modern Vue 3 + TypeScript + Tailwind CSS frontend for the TimeLog application, featuring a clean and responsive user interface for time tracking with structured tag management.

## Technology Stack

- **Framework**: Vue 3.4+ with Composition API
- **Language**: TypeScript 5.3+
- **Build Tool**: Vite 5.0+
- **Styling**: Tailwind CSS 3.4+
- **Icons**: Heroicons Vue 2.0+
- **HTTP Client**: Axios 1.6+
- **Date Utilities**: date-fns 3.0+

## Features

- ✅ Create, read, update, and delete time logs
- ✅ Structured tag selection with color-coded labels
- ✅ Real-time duration calculation
- ✅ Modern responsive UI with Tailwind CSS
- ✅ TypeScript for complete type safety
- ✅ Form validation and error handling
- ✅ Toast notification system
- ✅ API error handling with user feedback

## Development Setup

### Prerequisites

- Node.js 18+ and npm
- Go backend running on port 8083

### Installation

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Type check without emitting
npm run type-check
```

The development server runs on http://localhost:3000 and proxies API requests to the Go backend on port 8083.

### Building for Production

```bash
# Build for production
npm run build

# Preview production build
npm run preview
```

Built files go to `dist/` directory and are served by the Go backend at production time.

## Usage Scenarios

**Development Mode:**
1. Start Go backend: `./main` (runs on port 8083)
2. Start frontend dev server: `npm run dev` (runs on port 3000)
3. Access at http://localhost:3000

**Production Mode:**
1. Build frontend: `npm run build`
2. Start backend: `./main` (serves both API and frontend on port 8083)
3. Access at http://localhost:8083

## API Integration

The frontend communicates with the Go backend through a unified API structure:

### TimeLog Endpoints
- `GET /api/timelogs` - List all time logs with tag information
- `POST /api/timelogs` - Create new time log
- `GET /api/timelogs/:id` - Get specific time log
- `PUT /api/timelogs/:id` - Update time log
- `DELETE /api/timelogs/:id` - Delete time log

### Tag Endpoints
- `GET /api/tags` - List all available tags
- `POST /api/tags` - Create new tag (future admin feature)
- `PUT /api/tags/:id` - Update tag (future admin feature)
- `DELETE /api/tags/:id` - Delete tag (future admin feature)

### API Response Format
All endpoints return data in a consistent format:
```typescript
{
  data: T,           // The actual response data
  message: string,   // Success/error message
  status: number     // HTTP status code
}
```

## Project Structure

```
src/
├── api/              # API client functions
│   └── index.ts      # HTTP client setup and API endpoints
├── components/       # Vue components
│   ├── TimeLogList.vue    # Display list of time logs
│   └── TimeLogForm.vue    # Create/edit time log form
├── types/            # TypeScript type definitions
│   └── index.ts      # Interface definitions for TimeLog, Tag, etc.
├── utils/            # Utility functions
│   └── date.ts       # Date formatting and duration calculation
├── App.vue           # Main application component
├── main.ts           # Application entry point
└── style.css         # Global Tailwind CSS imports
```

## Key Components

### TimeLogList.vue
- Displays paginated list of time logs
- Shows color-coded tags with hover descriptions
- Calculates and displays duration for each entry
- Handles edit/delete actions

### TimeLogForm.vue
- Create/edit form for time logs
- Tag selection dropdown with colors
- DateTime pickers for start/end times
- Form validation and submission

### App.vue
- Main application layout
- Manages global state (time logs, tags, loading states)
- Handles API calls and error management
- Toast notification system

## TypeScript Types

### Core Interfaces
```typescript
interface Tag {
  id: number
  name: string
  color: string        // Hex color code
  description: string
  created_at: string
  updated_at: string
}

interface TimeLog {
  id: number
  start_time: string
  end_time?: string | null
  tag_id: number
  tag: Tag            // Populated tag object
  remarks: string
  created_at: string
  updated_at: string
}
```

## Styling Guidelines

- Uses Tailwind CSS utility-first approach
- Color scheme matches tag colors from backend
- Responsive design with mobile-first approach
- Consistent spacing and typography
- Hover states and smooth transitions

## Development Notes

- All API calls use the unified response format
- Form validation happens client-side with TypeScript
- Error handling provides user-friendly messages
- Duration calculation updates in real-time
- Tag selection is restricted to predefined options
- All timestamps are handled in ISO format with proper timezone conversion