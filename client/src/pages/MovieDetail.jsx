import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { movieAPI } from '../services/api';
import './MovieDetail.css';

const MovieDetail = () => {
  const { imdbId } = useParams();
  const [movie, setMovie] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchMovie();
  }, [imdbId]);

  const fetchMovie = async () => {
    try {
      const response = await movieAPI.getMovie(imdbId);
      setMovie(response.data);
      setLoading(false);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to fetch movie details');
      setLoading(false);
    }
  };

  if (loading) {
    return <div className="loading-container">Loading movie details...</div>;
  }

  if (error) {
    return <div className="error-container">{error}</div>;
  }

  if (!movie) {
    return <div className="error-container">Movie not found</div>;
  }

  const getRankingStars = (value) => {
    return '⭐'.repeat(value);
  };

  return (
    <div className="movie-detail-page">
      <div className="movie-detail-container">
        <div className="movie-poster">
          <img src={movie.poster_path} alt={movie.title} />
        </div>

        <div className="movie-info">
          <h1 className="movie-title">{movie.title}</h1>
          
          <div className="movie-rating-large">
            <span className="stars">{getRankingStars(movie.ranking.ranking_value)}</span>
            <span className="rating-text">{movie.ranking.ranking_name}</span>
          </div>

          <div className="movie-genres-list">
            {movie.genre.map((g, index) => (
              <span key={index} className="genre-tag-large">{g.genre_name}</span>
            ))}
          </div>

          <div className="movie-review">
            <h2>Review</h2>
            <p>{movie.admin_review || 'No review available'}</p>
          </div>

          {movie.youtube_id && (
            <div className="movie-trailer">
              <h2>Trailer</h2>
              <div className="video-wrapper">
                <iframe
                  src={`https://www.youtube.com/embed/${movie.youtube_id}`}
                  title="Movie Trailer"
                  frameBorder="0"
                  allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                  allowFullScreen
                ></iframe>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default MovieDetail;
