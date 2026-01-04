import {useState, useEffect} from 'react';
import axiosClient from '../../api/axiosConfig';
import Movies from '../movies/Movies';

const Home = () => {
    const [movies, setMovies] = useState([]);
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState('Loading movies...');

    //when the components hooks up to the DOM, we want to fetch the movies from the server
    useEffect(() => {
        const fetchMovies = async () => {
            setLoading(true);
            //display a loading indicator to the user
            setMessage('Loading movies...');
            try{
                const response = await axiosClient.get('/movies');
                setMovies(response.data);
                if(response.data.length === 0){
                    setMessage('No movies found');
                }
            }catch(error){
                console.error('Error fetching movies:', error);
                setMessage('Failed to load movies');
            }finally{
                setLoading(false);
            }
        }
        fetchMovies();
    }, [])
    return (
        <>
            {loading ? (
                <h2>Loading...</h2>
            ):
                <Movies movies={movies} message={message} />
            }
        </>
    )
}

export default Home;