// /src/App.tsx
import { useState } from 'react';
import LeadList from './components/LeadList';

interface Lead {
  [key: string]: any;
}

function App() {
  const [categoryID, setCategoryID] =
    useState<string>('padaria');
  const [zipcodeID, setZipcodeID] =
    useState<string>('01001000');
  const [radius, setRadius] = useState<number>(3000);
  const [maxResults, setMaxResults] = useState<number>(5);
  const [message, setMessage] = useState<string>('');
  const [leads, setLeads] = useState<Lead[]>([]);

  const BACKEND_URL_Search_GOOGLE =
    'http://192.168.0.9:8082';
  const BACKEND_URL_API = 'http://192.168.0.9:8085';

  const handleStartSearch = async () => {
    try {
      const url = `${BACKEND_URL_Search_GOOGLE}/start-search?category_id=${categoryID}&zipcode_id=${zipcodeID}&radius=${radius}&max_results=${maxResults}`;
      const response = await fetch(url);
      if (response.ok) {
        const text = await response.text();
        setMessage(text);
      } else {
        setMessage('Erro ao iniciar a busca');
      }
    } catch (error) {
      console.error(error);
      setMessage('Erro de conexão com o backend');
    }
  };

  const handleGetLeads = async () => {
    try {
      const url = `${BACKEND_URL_API}/list-leads`;
      const response = await fetch(url);
      console.log('leads list response', response);
      if (response.ok) {
        const data: Lead[] = await response.json();
        console.log('leads list data', data);
        setLeads(data);
      } else {
        setMessage('Erro ao buscar leads');
      }
    } catch (error) {
      console.error(error);
      setMessage('Erro de conexão ao buscar leads');
    }
  };

  return (
    <div style={{ padding: '1rem' }}>
      <h1>Buscar Leads no Google Places</h1>
      <div style={{ marginBottom: '1rem' }}>
        <label>Category ID:&nbsp;</label>
        <input
          type="text"
          value={categoryID}
          onChange={e => setCategoryID(e.target.value)}
        />
      </div>
      <div style={{ marginBottom: '1rem' }}>
        <label>Zipcode ID:&nbsp;</label>
        <input
          type="text"
          value={zipcodeID}
          onChange={e => setZipcodeID(e.target.value)}
        />
      </div>
      <div style={{ marginBottom: '1rem' }}>
        <label>Radius (m):&nbsp;</label>
        <input
          type="number"
          value={radius}
          onChange={e => setRadius(Number(e.target.value))}
        />
      </div>
      <div style={{ marginBottom: '1rem' }}>
        <label>Max Results:&nbsp;</label>
        <input
          type="number"
          value={maxResults}
          onChange={e =>
            setMaxResults(Number(e.target.value))
          }
        />
      </div>

      <button onClick={handleStartSearch}>
        Iniciar Busca
      </button>
      <p>{message}</p>

      <hr />
      <h2>Leads Salvos</h2>
      <button onClick={handleGetLeads}>Listar Leads</button>

      <LeadList leads={leads} />
    </div>
  );
}

export default App;
