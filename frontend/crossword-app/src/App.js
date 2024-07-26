import React, { useState, useEffect } from 'react';
import axios from 'axios';
import CrosswordGrid from './components/CrosswordGrid/CrosswordGrid';
import './App.css';

function App() {
  const [crosswordData, setCrosswordData] = useState(null);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.post('http://localhost:8000/generate?category=javascript');
        setCrosswordData(response.data);
      } catch (err) {
        console.error("Error fetching crossword data:", err);
        setError("Failed to load crossword data. Please try again later.");
      }
    };

    fetchData();
  }, []);

  if (error) {
    return <div className="error">{error}</div>;
  }

  if (!crosswordData) {
    return <div>Loading...</div>;
  }

  return (
    <div className="App">
      <h1>Python Crossword Puzzle</h1>
      <CrosswordGrid
        gridData={{
          grid_x_dim: crosswordData.grid_y_dim,
          grid_y_dim: crosswordData.grid_x_dim
        }}
        characterPositions={crosswordData.characterPositions}
        words={crosswordData.words}
      />
      <div className="clues">
        <h2>Clues</h2>
        <div className="clues-container">
          <div className="across-clues">
            <h3>Across</h3>
            <ul>
              {crosswordData.words.filter(word => word.direction === 0).map((word, index) => (
                <li key={index}>
                  {word.number}: {word.clue}
                </li>
              ))}
            </ul>
          </div>
          <div className="down-clues">
            <h3>Down</h3>
            <ul>
              {crosswordData.words.filter(word => word.direction === 1).map((word, index) => (
                <li key={index}>
                  {word.number}: {word.clue}
                </li>
              ))}
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;