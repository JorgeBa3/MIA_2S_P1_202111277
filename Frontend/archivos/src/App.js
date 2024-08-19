import React, { useState } from 'react';
import './App.css';

function App() {
  const [name, setName] = useState('');
  const [response, setResponse] = useState('');

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

  const handleSubmit = async () => {
    try {
      const response = await fetch('http://localhost:8080/greet', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name: name }),
      });

      const data = await response.json();
      // Unir los mensajes en una sola cadena para mostrar
      setResponse(data.messages.join('\n'));
    } catch (error) {
      console.error('Error:', error);
    }
  };

  return (
    <div className="container">
      <h1>Proyecto 1</h1>
      <textarea
        value={name}
        onChange={handleNameChange}
        placeholder="Escribe o carga un archivo"
      />
      <input type="file" onChange={handleFileChange} />
      <button onClick={handleSubmit}>Ejecutar</button>
      <textarea
        value={response}
        readOnly
        placeholder="Aquí se mostrará la información"
      />
    </div>
  );
}

export default App;
