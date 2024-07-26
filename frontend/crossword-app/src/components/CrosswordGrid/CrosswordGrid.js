import React, { useState } from 'react';
import './CrosswordGrid.css';

const CrosswordGrid = ({ gridData, characterPositions, words }) => {
  const [highlightedWord, setHighlightedWord] = useState(null);

  const isPartOfWord = (x, y, word) => {
    if (word.direction === 1) { // Across
      return y === word.startY && x >= word.startX && x < word.startX + word.length;
    } else { // Down
      return x === word.startX && y >= word.startY && y < word.startY + word.length;
    }
  };

  const getWordsAtPosition = (x, y) => {
    return words.filter(word => word.startX === x && word.startY === y);
  };

  const renderCell = (y, x) => {
    const isActive = characterPositions.some(pos => pos.x === x && pos.y === y);
    const cellWords = getWordsAtPosition(x, y);
    const isHighlighted = highlightedWord && isPartOfWord(x, y, highlightedWord);

    const handleMouseEnter = () => {
      if (cellWords.length > 0) setHighlightedWord(cellWords[0]);
    };

    const handleMouseLeave = () => {
      setHighlightedWord(null);
    };

    const handleClick = () => {
      if (cellWords.length > 1) {
        // Toggle between the two words on click
        setHighlightedWord(highlightedWord === cellWords[0] ? cellWords[1] : cellWords[0]);
      }
    };

    return (
      <div
        key={`${x}-${y}`}
        className={`cell ${isActive ? 'active' : ''} ${isHighlighted ? 'highlighted' : ''}`}
        onMouseEnter={handleMouseEnter}
        onMouseLeave={handleMouseLeave}
        onClick={handleClick}
      >
        {cellWords.length > 0 && (
          <div className="cell-numbers">
            {cellWords.map((word, index) => (
              <span key={index} className={`cell-number ${word.direction === 0 ? 'across' : 'down'}`}>
                {word.number}
              </span>
            ))}
          </div>
        )}
        {isActive && <input type="text" maxLength="1" className="cell-input" />}
      </div>
    );
  };

  const renderGrid = () => {
    const grid = [];
    for (let y = 0; y < gridData.grid_y_dim; y++) {
      for (let x = 0; x < gridData.grid_x_dim; x++) {
        grid.push(renderCell(x, y));
      }
    }
    return grid;
  };

  return (
    <div
      className="crossword-grid"
      style={{
        gridTemplateColumns: `repeat(${gridData.grid_x_dim}, 1fr)`,
        gridTemplateRows: `repeat(${gridData.grid_y_dim}, 1fr)`
      }}
    >
      {renderGrid()}
    </div>
  );
};

export default CrosswordGrid;