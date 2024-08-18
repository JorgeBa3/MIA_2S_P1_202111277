import React, { useState } from 'react';
import './App.css';

function App() {
  const [name, setName] = useState('');
  const [fileContent, setFileContent] = useState('');

  const handleNameChange = (event) => {
    setName(event.target.value);
  };

  const handleFileChange = (event) => {
    const file = event.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (e) => {
        setName(e.target.result);
      };
      reader.readAsText(file);
    }
  };

  const handleSubmit = () => {
    setFileContent(name);
  };

  return (
    <div className="container">
      <h1>Proyecto 1</h1>
      <textarea
        value={name}
        onChange={handleNameChange}
        placeholder="Escribe tu nombre o carga un archivo"
      />
      <input 
        type="file" 
        onChange={handleFileChange} 
      />
      <button onClick={handleSubmit}>Ejecutar</button>
      <textarea
        value={`La información es: ${fileContent}`}
        readOnly
        placeholder="Aquí se mostrará la información"
      />
    </div>
  );
}

export default App;
