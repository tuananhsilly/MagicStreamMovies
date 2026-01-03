import { useState, useEffect } from 'react';
import { movieAPI } from '../services/api';
import MovieCard from '../components/MovieCard';
import './Movies.css';

const Movies = () => {
  const [movies, setMovies] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchMovies();
  }, []);

  const fetchMovies = async () => {
    try {
      const response = await movieAPI.getMovies();
      setMovies(response.data || []);
      setLoading(false);
    } catch (err) {
      setError('Failed to fetch movies. Make sure the backend is running on http://localhost:8080');
      setMovies([]);
      setLoading(false);
    }
  };

  if (loading) {
    return <div className="loading-container">Loading movies...</div>;
  }

  if (error) {
    return <div className="error-container">{error}</div>;
  }

  return (
    <div className="movies-page">
      <div className="container">
        <h1 className="page-title">All Movies</h1>
        <div className="movies-grid">
          {movies && movies.length > 0 ? (
            movies.map((movie) => (
              <MovieCard key={movie.imdb_id} movie={movie} />
            ))
          ) : (
            <p className="no-movies">No movies available</p>
          )}
        </div>
      </div>
    </div>
  );
};

export default Movies;
