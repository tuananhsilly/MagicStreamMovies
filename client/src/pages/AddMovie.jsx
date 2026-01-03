import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { movieAPI } from '../services/api';
import './AddMovie.css';

const AddMovie = () => {
  const [formData, setFormData] = useState({
    imdb_id: '',
    title: '',
    poster_path: '',
    youtube_id: '',
    admin_review: '',
    ranking: {
      ranking_value: 3,
      ranking_name: 'ok',
    },
    genre: [],
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const availableGenres = [
    { id: 1, name: 'Action' },
    { id: 2, name: 'Comedy' },
    { id: 3, name: 'Drama' },
    { id: 4, name: 'Horror' },
    { id: 5, name: 'Romance' },
    { id: 6, name: 'Sci-Fi' },
    { id: 7, name: 'Thriller' },
    { id: 8, name: 'Animation' },
  ];

  const rankingOptions = [
    { value: 1, name: 'excellent' },
    { value: 2, name: 'good' },
    { value: 3, name: 'ok' },
    { value: 4, name: 'bad' },
    { value: 5, name: 'terrible' },
  ];

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleRankingChange = (e) => {
    const value = parseInt(e.target.value);
    const rankingName = rankingOptions.find(r => r.value === value)?.name || 'ok';
    
    setFormData({
      ...formData,
      ranking: {
        ranking_value: value,
        ranking_name: rankingName,
      },
    });
  };

  const handleGenreToggle = (genre) => {
    const isSelected = formData.genre.some(g => g.genre_id === genre.id);
    
    if (isSelected) {
      setFormData({
        ...formData,
        genre: formData.genre.filter(g => g.genre_id !== genre.id),
      });
    } else {
      setFormData({
        ...formData,
        genre: [
          ...formData.genre,
          { genre_id: genre.id, genre_name: genre.name }
        ],
      });
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    if (formData.genre.length === 0) {
      setError('Please select at least one genre');
      return;
    }

    setLoading(true);

    try {
      await movieAPI.addMovie(formData);
      navigate('/movies');
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to add movie');
      setLoading(false);
    }
  };

  return (
    <div className="add-movie-page">
      <div className="container">
        <h1 className="page-title">Add New Movie</h1>
        
        <form onSubmit={handleSubmit} className="add-movie-form">
          {error && <div className="error-message">{error}</div>}

          <div className="form-group">
            <label htmlFor="imdb_id">IMDB ID *</label>
            <input
              type="text"
              id="imdb_id"
              name="imdb_id"
              value={formData.imdb_id}
              onChange={handleChange}
              placeholder="e.g., tt1234567"
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="title">Title *</label>
            <input
              type="text"
              id="title"
              name="title"
              value={formData.title}
              onChange={handleChange}
              placeholder="Movie title"
              required
              minLength="2"
            />
          </div>

          <div className="form-group">
            <label htmlFor="poster_path">Poster URL *</label>
            <input
              type="url"
              id="poster_path"
              name="poster_path"
              value={formData.poster_path}
              onChange={handleChange}
              placeholder="https://example.com/poster.jpg"
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="youtube_id">YouTube Video ID *</label>
            <input
              type="text"
              id="youtube_id"
              name="youtube_id"
              value={formData.youtube_id}
              onChange={handleChange}
              placeholder="e.g., dQw4w9WgXcQ"
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="ranking">Rating *</label>
            <select
              id="ranking"
              name="ranking"
              value={formData.ranking.ranking_value}
              onChange={handleRankingChange}
            >
              {rankingOptions.map((option) => (
                <option key={option.value} value={option.value}>
                  {'⭐'.repeat(option.value)} - {option.name}
                </option>
              ))}
            </select>
          </div>

          <div className="form-group">
            <label>Genres *</label>
            <div className="genre-selector">
              {availableGenres.map((genre) => (
                <button
                  key={genre.id}
                  type="button"
                  className={`genre-btn ${
                    formData.genre.some(g => g.genre_id === genre.id) ? 'selected' : ''
                  }`}
                  onClick={() => handleGenreToggle(genre)}
                >
                  {genre.name}
                </button>
              ))}
            </div>
          </div>

          <div className="form-group">
            <label htmlFor="admin_review">Review</label>
            <textarea
              id="admin_review"
              name="admin_review"
              value={formData.admin_review}
              onChange={handleChange}
              placeholder="Write your review here..."
              rows="6"
            />
          </div>

          <button type="submit" className="btn-submit" disabled={loading}>
            {loading ? 'Adding Movie...' : 'Add Movie'}
          </button>
        </form>
      </div>
    </div>
  );
};

export default AddMovie;
