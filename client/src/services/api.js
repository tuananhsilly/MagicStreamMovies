import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080';

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add token to headers
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Token expired or invalid
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Movie API calls
export const movieAPI = {
  // Get all movies (unprotected)
  getMovies: () => api.get('/movies'),
  
  // Get single movie (protected)
  getMovie: (imdbId) => api.get(`/movie/${imdbId}`),
  
  // Add new movie (protected - admin only)
  addMovie: (movieData) => api.post('/addmovie', movieData),
};

// Authentication API calls
export const authAPI = {
  // Register new user
  register: (userData) => api.post('/register', userData),
  
  // Login user
  login: (credentials) => api.post('/login', credentials),
};

export default api;
