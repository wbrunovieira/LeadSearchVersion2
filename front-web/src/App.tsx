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
    <div className="bg-slate-100 text-black p-4 max-w-4xl mx-auto my-8 rounded-lg shadow-md">
      <h1 className="text-3xl font-bold text-blue-600 mb-6">
        Buscar Leads no Google Places
      </h1>

      <div className="mb-4 bg-white p-4 rounded-md">
        <label className="block text-gray-700 font-medium mb-2">
          Category ID:
        </label>
        <input
          type="text"
          value={categoryID}
          onChange={e => setCategoryID(e.target.value)}
          className="w-full p-2 border border-gray-300 rounded text-black"
        />
      </div>

      <div className="mb-4 bg-white p-4 rounded-md">
        <label className="block text-gray-700 font-medium mb-2">
          Zipcode ID:
        </label>
        <input
          type="text"
          value={zipcodeID}
          onChange={e => setZipcodeID(e.target.value)}
          className="w-full p-2 border border-gray-300 rounded text-black"
        />
      </div>

      <div className="mb-4 bg-white p-4 rounded-md">
        <label className="block text-gray-700 font-medium mb-2">
          Radius (m):
        </label>
        <input
          type="number"
          value={radius}
          onChange={e => setRadius(Number(e.target.value))}
          className="w-full p-2 border border-gray-300 rounded text-black"
        />
      </div>

      <div className="mb-4 bg-white p-4 rounded-md">
        <label className="block text-gray-700 font-medium mb-2">
          Max Results:
        </label>
        <input
          type="number"
          value={maxResults}
          onChange={e =>
            setMaxResults(Number(e.target.value))
          }
          className="w-full p-2 border border-gray-300 rounded text-black"
        />
      </div>

      <button
        onClick={handleStartSearch}
        className="bg-blue-500 hover:bg-blue-600 text-white py-2 px-6 rounded-md font-semibold"
      >
        Iniciar Busca
      </button>

      <p className="my-4 p-2 bg-blue-50 text-blue-800 rounded-md">
        {message}
      </p>

      <hr className="my-8 border-gray-300" />

      <h2 className="text-2xl font-bold text-blue-600 mb-4">
        Leads Salvos
      </h2>

      <button
        onClick={handleGetLeads}
        className="bg-green-500 hover:bg-green-600 text-white py-2 px-6 rounded-md font-semibold mb-4"
      >
        Listar Leads
      </button>

      <LeadList leads={leads} />
    </div>
  );
}

export default App;
