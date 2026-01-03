import { Link } from 'react-router-dom';
import './MovieCard.css';

const MovieCard = ({ movie }) => {
  const getRankingStars = (value) => {
    return '⭐'.repeat(value);
  };

  return (
    <Link to={`/movie/${movie.imdb_id}`} className="movie-card">
      <div className="movie-card-image">
        <img src={movie.poster_path} alt={movie.title} />
      </div>
      <div className="movie-card-info">
        <h3 className="movie-title">{movie.title}</h3>
        <div className="movie-rating">
          <span className="stars">{getRankingStars(movie.ranking.ranking_value)}</span>
          <span className="rating-name">{movie.ranking.ranking_name}</span>
        </div>
        <div className="movie-genres">
          {movie.genre.slice(0, 3).map((g, index) => (
            <span key={index} className="genre-tag">{g.genre_name}</span>
          ))}
        </div>
      </div>
    </Link>
  );
};

export default MovieCard;
