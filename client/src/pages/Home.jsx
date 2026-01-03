import { Link } from 'react-router-dom';
import './Home.css';

const Home = () => {
  return (
    <div className="home-page">
      <section className="hero">
        <div className="hero-content">
          <h1 className="hero-title">Welcome to Magic Stream Movies</h1>
          <p className="hero-subtitle">
            Discover and explore amazing movies from around the world
          </p>
          <div className="hero-actions">
            <Link to="/movies" className="btn btn-primary">
              Browse Movies
            </Link>
            <Link to="/register" className="btn btn-secondary">
              Get Started
            </Link>
          </div>
        </div>
      </section>

      <section className="features">
        <div className="container">
          <h2 className="section-title">Features</h2>
          <div className="features-grid">
            <div className="feature-card">
              <div className="feature-icon">🎬</div>
              <h3>Huge Collection</h3>
              <p>Browse through thousands of movies across all genres</p>
            </div>
            <div className="feature-card">
              <div className="feature-icon">⭐</div>
              <h3>Ratings & Reviews</h3>
              <p>Get honest ratings and detailed reviews from our admins</p>
            </div>
            <div className="feature-card">
              <div className="feature-icon">🎥</div>
              <h3>Watch Trailers</h3>
              <p>Preview movies with high-quality trailers before watching</p>
            </div>
            <div className="feature-card">
              <div className="feature-icon">👤</div>
              <h3>Personalized</h3>
              <p>Set your favorite genres and get personalized recommendations</p>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
};

export default Home;
