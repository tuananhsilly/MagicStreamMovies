# Magic Stream Movies - Frontend

React frontend application for Magic Stream Movies platform, built with Vite.

## 🚀 Features

- **User Authentication**: Login and registration with JWT token management
- **Movie Browsing**: Browse and search through movie collections
- **Movie Details**: View detailed information, ratings, and trailers
- **Admin Functions**: Add new movies (admin only)
- **Protected Routes**: Secure routes requiring authentication
- **Responsive Design**: Mobile-friendly dark theme UI

## 📋 Prerequisites

- Node.js (v18 or higher)
- npm or yarn
- Backend Go server running on port 8080

## 🛠️ Installation

```bash
# Navigate to client directory
cd client

# Install dependencies
npm install

# Start development server
npm run dev
```

The app will run on `http://localhost:3000`

## 🔧 Backend Connection

Make sure the Go backend server is running on `http://localhost:8080` before starting the frontend.

## 📁 Project Structure

```
client/
├── src/
│   ├── components/       # Reusable UI components
│   │   ├── Navbar.jsx
│   │   ├── MovieCard.jsx
│   │   └── ProtectedRoute.jsx
│   ├── pages/           # Page components
│   │   ├── Home.jsx
│   │   ├── Login.jsx
│   │   ├── Register.jsx
│   │   ├── Movies.jsx
│   │   ├── MovieDetail.jsx
│   │   └── AddMovie.jsx
│   ├── services/        # API services
│   │   └── api.js
│   ├── context/         # React context
│   │   └── AuthContext.jsx
│   └── utils/           # Utility functions
├── public/              # Static assets
└── index.html
```

## 🎯 Available Routes

- `/` - Home page
- `/login` - User login
- `/register` - User registration
- `/movies` - Browse all movies
- `/movie/:imdbId` - Movie details (protected)
- `/add-movie` - Add new movie (admin only, protected)

## 🔐 Authentication

The app uses JWT tokens for authentication:
- Access token is stored in `localStorage`
- Token is automatically added to API requests
- Protected routes redirect to login if not authenticated
- Admin routes require `ADMIN` role

## 🎨 UI Design

- **Dark Theme**: Netflix-inspired design
- **Color Scheme**: 
  - Background: #0a0a0a, #141414
  - Primary: #e50914 (Red)
  - Text: White/Gray
- **Responsive**: Mobile-first design

## 📡 API Endpoints Used

### Unprotected
- `GET /movies` - Get all movies
- `POST /register` - Register new user
- `POST /login` - Login user

### Protected (Requires Token)
- `GET /movie/:imdb_id` - Get single movie
- `POST /addmovie` - Add new movie (Admin only)

## 🏗️ Build for Production

```bash
npm run build
```

The build files will be in the `dist` folder.

## 📝 Environment Variables

Create a `.env` file if needed (currently using hardcoded backend URL):

```env
VITE_API_URL=http://localhost:8080
```

## 🐛 Troubleshooting

**CORS Issues**: Make sure the Go backend allows CORS from `http://localhost:3000`

**401 Errors**: Clear localStorage and login again

**Connection Failed**: Ensure backend server is running on port 8080

## 👥 User Roles

- **USER**: Can browse and view movies
- **ADMIN**: Can add new movies

## 🎬 Usage

1. **Register** a new account
2. **Select** your favorite genres
3. **Login** with your credentials
4. **Browse** movies on the home page
5. **Click** on a movie to view details and trailer
6. **Admin**: Use "Add Movie" to add new content

## 🔄 Development

The app uses React Router for navigation and Axios for API calls. State management is handled by Context API for authentication.

Hot reload is enabled in development mode.

---

Built with ❤️ using React + Vite
